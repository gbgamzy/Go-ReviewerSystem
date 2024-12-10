package repository

import (
	"database/sql"
	"fmt"
)

type UserRepository struct {
	DB *sql.DB
}

// Login or Signup a user
func (r *UserRepository) LoginOrSignup(email, name string) (int, error) {
	var userID int

	// Check if the user already exists
	query := `SELECT id FROM users WHERE email = $1`
	err := r.DB.QueryRow(query, email).Scan(&userID)

	if err == sql.ErrNoRows {
		// User does not exist, create a new user
		insertQuery := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
		err = r.DB.QueryRow(insertQuery, name, email).Scan(&userID)
		if err != nil {
			return 0, fmt.Errorf("error creating new user: %v", err)
		}
		return userID, nil
	} else if err != nil {
		return 0, fmt.Errorf("error checking user: %v", err)
	}

	// User exists
	return userID, nil
}
