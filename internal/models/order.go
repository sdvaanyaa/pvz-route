package models

import "time"

const (
	StatusAccepted = "accepted"
	StatusIssued   = "issued"
	StatusReturned = "returned"
)

type Order struct {
	ID            string     `json:"id"`
	UserID        string     `json:"receipt_id"`
	StorageExpire time.Time  `json:"storage_deadline"`
	Status        string     `json:"status"` // "accepted", "issued", "returned"
	CreatedAt     time.Time  `json:"created_at"`
	IssuedAt      *time.Time `json:"issued_at"`
}
