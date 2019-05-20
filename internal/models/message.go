package models

import "time"

// Define Model of message
// Message of struct message
type Message struct {
	ID          uint64
	UserID      uint64
	ToUserID    uint64
	Subject     string
	Text        string
	Status      int8
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
