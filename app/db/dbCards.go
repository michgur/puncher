package db

import (
	"fmt"

	"github.com/michgur/puncher/app/model"
	"github.com/pquerna/otp/totp"
)

func init() {
	if _, err := Exec("card-details-create.sql"); err != nil {
		panic(err)
	}
}

func generateCardSecret(details model.CardDetails) (string, error) {
	otp, err := totp.Generate(totp.GenerateOpts{
		Issuer:      details.ID,
		AccountName: details.Name,
	})
	if err != nil {
		return "", err
	}
	return otp.Secret(), nil
}

func NewCard(details model.CardDetails) error {
	if details.Secret == "" {
		secret, err := generateCardSecret(details)
		if err != nil {
			return err
		}
		details.Secret = secret
	}

	_, err := Exec("card-details-insert.sql", details.ID, details.Name, details.Secret)
	return err
}

func GetCardDetails(cardID string) (model.CardDetails, error) {
	var details model.CardDetails
	err := QueryRow("card-details-select.sql", cardID).Scan(&details.ID, &details.Name, &details.Secret)
	return details, err
}

func GetAllCardDetails() ([]model.CardDetails, error) {
	rows, err := Query("card-details-select-all.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []model.CardDetails
	for rows.Next() {
		var d model.CardDetails
		err = rows.Scan(&d.ID, &d.Name, &d.Secret)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		details = append(details, d)
	}

	return details, nil
}
