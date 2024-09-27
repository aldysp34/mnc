package model

import (
	"time"

	"github.com/google/uuid"
)

type TopUp struct {
	ID            string    `gorm:"primaryKey"`
	UserID        uuid.UUID `gorm:"index"`
	Amount        int       `json:"amount"`
	BalanceBefore int       `json:"balance_before"`
	BalanceAfter  int       `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}
