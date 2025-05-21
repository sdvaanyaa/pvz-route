package models

import "time"

type Return struct {
	ID         string    `json:"order_id"`
	ReturnDate time.Time `json:"return_date"`
	Reason     string    `json:"reason"`
}
