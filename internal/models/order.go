package models

import "time"

const (
	StatusAccepted = "accepted"
	StatusIssued   = "issued"
	StatusReturned = "returned"
	StatusArchived = "archived"
)

type OrderStatusChange struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Order struct {
	ID            string              `json:"id"`
	UserID        string              `json:"user_id"`
	StorageExpire time.Time           `json:"storage_deadline"`
	Status        string              `json:"status"` // "accepted", "issued", "returned"
	CreatedAt     time.Time           `json:"created_at"`
	IssuedAt      *time.Time          `json:"issued_at"`
	ReturnedAt    *time.Time          `json:"returned_at"`
	ArchivedAt    *time.Time          `json:"archived_at"`
	History       []OrderStatusChange `json:"history"`
	Weight        float64             `json:"weight"`
	Price         float64             `json:"price"`
	PackageType   string              `json:"package"`
}
