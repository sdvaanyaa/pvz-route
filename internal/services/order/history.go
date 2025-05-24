package order

import (
	"sort"
	"time"
)

type HistoryEntry struct {
	OrderID   string
	Status    string
	Timestamp time.Time
}

func (s *orderService) History() ([]*HistoryEntry, error) {
	orders, err := s.storage.GetOrders()
	if err != nil {
		return nil, err
	}

	var history []*HistoryEntry
	for _, order := range orders {
		for _, change := range order.History {
			history = append(history, &HistoryEntry{
				OrderID:   order.ID,
				Status:    change.Status,
				Timestamp: change.Timestamp,
			})
		}
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].Timestamp.After(history[j].Timestamp)
	})

	return history, nil
}
