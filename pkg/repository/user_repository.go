package repository

import (
	"ReviewerService/pkg/models"
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

// Get users with search and pagination
func (r *UserRepository) GetUsers(search string, page, limit int) ([]models.User, int, error) {
	var users []models.User
	offset := (page - 1) * limit

	// Count total users matching the search
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE name ILIKE $1 OR email ILIKE $1
	`
	err := r.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting users: %v", err)
	}

	// Fetch users with pagination
	query := `
		SELECT id, name, email, created_at
		FROM users
		WHERE name ILIKE $1 OR email ILIKE $1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning user: %v", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *UserRepository) GetUserEmailByID(userID int) (string, error) {
	var email string
	query := `SELECT email FROM users WHERE id = $1`
	err := r.DB.QueryRow(query, userID).Scan(&email)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found")
	}
	return email, err
}
