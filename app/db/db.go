package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/michgur/puncher/app/model"
)

var sqlsPath = "./sqls"
var sqls = map[string]string{}
var db *sql.DB

func readSQLs() error {
	// iterate over ./sqls, store the contents of each file in a map
	// key: filename, value: contents

	fnames, err := os.ReadDir(sqlsPath)
	if err != nil {
		return err
	}

	for _, fname := range fnames {
		f, err := os.ReadFile(sqlsPath + "/" + fname.Name())
		if err != nil {
			return err
		}
		sqls[fname.Name()] = string(f)
	}

	return nil
}

func getAllCardInstances() (cardInstances []model.CardInstance, err error) {
	rows, err := db.Query(sqls["select-all-card-instances.sql"])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		println(cardInstances)
		var ci model.CardInstance
		rows.Scan(&ci.ID, &ci.CardID, &ci.Slots)
		cardInstances = append(cardInstances, ci)
	}
	println(cardInstances)

	return cardInstances, nil
}

func insertCardInstance(cardInstance model.CardInstance) error {
	_, err := db.Exec(sqls["insert-card-instance.sql"], cardInstance.CardID, cardInstance.Slots)
	return err
}

func init() {
	// fetch queries
	var err = readSQLs()
	// create the database
	db, err = sql.Open("sqlite3", "puncher.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	// Create CardInstances table
	_, err = db.Exec(sqls["create-card-instances.sql"])
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}
}

func Main() {
	for i := range 4 {
		err := insertCardInstance(model.CardInstance{CardID: i * 3, Slots: (i * 45480832) % 10})
		if err != nil {
			fmt.Println("Error inserting into table:", err)
			return
		}
	}

	rows, err := getAllCardInstances()
	if err != nil {
		fmt.Println("Error fetching from table:", err)
		return
	}

	fmt.Println("CardInstances:")
	for _, row := range rows {
		fmt.Println(row)
	}

	fmt.Println("SUCCESS!")
}
