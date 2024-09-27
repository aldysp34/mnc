// models/user.go
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`       // Unique identifier for the user
	FirstName   string    `json:"first_name"`                 // User's first name
	LastName    string    `json:"last_name"`                  // User's last name
	PhoneNumber string    `gorm:"unique" json:"phone_number"` // Unique phone number for the user
	Address     string    `json:"address"`                    // User's address
	Pin         string    `json:"pin"`                        // User's PIN for authentication
	Balance     int       `json:"balance"`                    // User's current balance
	CreatedAt   time.Time `json:"created_at"`                 // Timestamp of when the user was created
}

// NewUser creates a new user instance
func NewUser(firstName, lastName, phoneNumber, address, pin string) *User {
	return &User{
		ID:          uuid.New(),
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
		Address:     address,
		Pin:         pin,
		Balance:     0, // Initialize balance to 0
		CreatedAt:   time.Now(),
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
