package controllers

import (
	"encoding/json"

	model "github.com/Dparty/model/restaurant"
)

type cooperateMessage struct {
	Orders []model.Order
}

func (c cooperateMessage) ToJson() []byte {
	s, _ := json.Marshal(c)
	return s
}

var tables = make(map[string][]model.Order)
