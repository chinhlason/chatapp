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

//--------------------------USER REPO-------------------------------

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

func (r *Repository) FriendRequest(ctx context.Context, userId, friendId string) error {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec("INSERT INTO friends (id_user, id_friend, status) VALUES ($1, $2, $3)", userId, friendId, "PENDING")
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ChangeFriendRequestStatus(ctx context.Context, id, status string) error {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec("UPDATE friends SET status = $1 WHERE friends.id = $2", status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetFriendRequests(ctx context.Context, userId string) ([]FriendRequest, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := r.db.Query("select f.id as id, req.username as requester, f.status as status, rec.username as receiver "+
		"from friends f "+
		"join users req on req.id = f.id_friend "+
		"join users rec on rec.id = f.id_user "+
		"where f.id_user = $1 and f.status = 'PENDING' "+
		"order by f.id desc",
		userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friendReqs []FriendRequest
	for rows.Next() {
		var friend FriendRequest
		err := rows.Scan(&friend.Id, &friend.Requester, &friend.Status, &friend.Receiver)
		if err != nil {
			return nil, err
		}
		friendReqs = append(friendReqs, friend)
	}
	return friendReqs, nil
}

func (r *Repository) GetListFriends(ctx context.Context, username string, limit, offset int) ([]Friend, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := r.db.Query("select f.id as id, frd.username as username, f.interaction_at as interaction_at "+
		"from friends f "+
		"join users frd on frd.id = f.id "+
		"join users u on u.id = f.id_user "+
		"where u.username = $1 and f.status = 'ACCEPTED' "+
		"order by f.interaction_at desc limit $2 offset $3",
		username, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friends []Friend
	for rows.Next() {
		var friend Friend
		err := rows.Scan(&friend.Id, &friend.Username, &friend.InteractionAt)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func (r *Repository) CreateChatRoom(ctx context.Context, name string) (string, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var id string
	err := r.db.QueryRow("INSERT INTO rooms (name) VALUES ($1) RETURNING id", name).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) AddUserToRoom(ctx context.Context, idUser, idRoom, idRole string) error {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec("INSERT INTO user_in_room (id_user, id_room, id_role) VALUES ($1, $2, $3)", idUser, idRoom, idRole)
	if err != nil {
		return err
	}
	return nil
}
