package db

import (
	"fmt"

	"github.com/michgur/puncher/app/design"
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
	var designJSON string
	err := QueryRow("card-details-select.sql", cardID).Scan(&details.ID, &details.Name, &details.Secret, &designJSON)
	if err != nil {
		return details, err
	}
	details.Design, err = design.ParseCardDesign(designJSON)
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
		var designJSON string
		err = rows.Scan(&d.ID, &d.Name, &d.Secret, &designJSON)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		// parse design as JSON
		d.Design, err = design.ParseCardDesign(designJSON)
		if err != nil {
			fmt.Println("Error parsing design:", err)
		}
		details = append(details, d)
	}

	return details, nil
}

func SetCardDesign(cardID string, design design.CardDesign) error {
	designJSON := design.ToJSON()
	_, err := Exec("card-details-update-design.sql", designJSON, cardID)
	return err
}
