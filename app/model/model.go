package model

type CardDetails struct {
	ID     string `sql:"card_id" json:"cardId"`
	Name   string `sql:"name" json:"name"`
	Secret string `sql:"secret"`
}

type CardInstance struct {
	ID     int
	CardID string
	Slots  int
}
