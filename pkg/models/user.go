package models

import "time"

type User struct {
	ID        int       `json:"id"`         // Unique identifier for the user
	Name      string    `json:"name"`       // Name of the user
	Email     string    `json:"email"`      // Email of the user
	CreatedAt time.Time `json:"created_at"` // Timestamp when the user was created
}
