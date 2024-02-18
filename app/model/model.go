package model

import "github.com/michgur/puncher/app/design"

type CardDetails struct {
	ID     string            `sql:"card_id" json:"cardID"`
	Name   string            `sql:"name" json:"cardName"`
	Secret string            `sql:"secret"`
	Design design.CardDesign `sql:"design" json:"design"`
}

type CardInstance struct {
	ID     int
	CardID string
	Slots  int
}
