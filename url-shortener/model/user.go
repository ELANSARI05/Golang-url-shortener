package model

import (
	"database/sql"
	"errors"

	"url-shortener/db"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

// CreateUser inserts a new user into the database
func CreateUser(username, passwordHash string) error {
	_, err := db.DB.Exec(`
        INSERT INTO users (username, password_hash) VALUES (?, ?)
    `, username, passwordHash)

	return err
}

// GetUserByUsername returns a user if the username exists
func GetUserByUsername(username string) (*User, error) {
	row := db.DB.QueryRow(`
        SELECT id, username, password_hash FROM users WHERE username = ?
    `, username)

	user := &User{}
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
