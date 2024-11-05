package api

import (
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) InsertUser(ctx context.Context, username, password string) error {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserById(ctx context.Context, id string) (*User, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var user User
	err := r.db.QueryRow("SELECT username, password FROM users WHERE id = $1", id).
		Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var user User
	err := r.db.QueryRow("SELECT * FROM users WHERE username = $1", username).
		Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}
