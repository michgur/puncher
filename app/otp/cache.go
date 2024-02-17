package otp

// cache for OTPs using SQLite
// key: cardID, value: list of transaction IDs
// each transaction ID lives for 30 seconds
// when a transaction ID is used, it is removed from the list
// when a card is removed, all its transaction IDs are removed

import (
	"time"

	"github.com/michgur/puncher/app/db"
)

type TransactionKey struct {
	TransactionID int    `sql:"transaction_id"`
	CardID        string `sql:"card_id"`
	UnixTime      int64  `sql:"unix_time"`
}

func init() {
	// create the table
	_, err := db.Exec("otp-cache-create.sql")
	if err != nil {
		panic(err)
	}
}

func AddTransaction(cardID string) (TransactionKey, error) {
	unixTime := time.Now().Unix()
	res, err := db.Exec("otp-cache-insert.sql", cardID, unixTime)
	if err != nil {
		return TransactionKey{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return TransactionKey{}, err
	}
	return TransactionKey{int(id), cardID, unixTime}, nil
}

func GetTransactions(cardID string) ([]TransactionKey, error) {
	rows, err := db.Query("otp-cache-select.sql", cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []TransactionKey
	for rows.Next() {
		var tk TransactionKey
		rows.Scan(&tk.TransactionID, &tk.CardID, &tk.UnixTime)
		transactions = append(transactions, tk)
	}

	return transactions, nil
}

func RemoveTransaction(key TransactionKey) error {
	_, err := db.Exec("otp-cache-delete.sql", key.TransactionID)
	return err
}
