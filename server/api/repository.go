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
	err := r.db.QueryRow("SELECT users.id, users.username, users.password FROM users WHERE username = $1", username).
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
	_, err := r.db.Exec("UPDATE friends SET status = $1 WHERE friends.id = $2 ", status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ChangeFriendRequestStatusAndReturnId(ctx context.Context, tx *sql.Tx, id, status string) (string, string, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var userId, friendId string
	err := tx.QueryRow("UPDATE friends SET status = $1 WHERE friends.id = $2 RETURNING friends.id_user, friends.id_friend",
		status, id).Scan(&userId, &friendId)
	if err != nil {
		return "", "", err
	}
	return userId, friendId, nil
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
	rows, err := r.db.Query("SELECT "+
		"COALESCE(MAX(r.id), 0) AS id_room, "+
		"f.id AS id_friend, "+
		"frd.username AS friend_username, "+
		"frd.is_online AS is_online, "+
		"COALESCE(MAX(f.interaction_at), '2021-01-01 00:00:00') AS interaction_at "+
		"FROM "+
		"friends f "+
		"JOIN "+
		"users u ON u.id = f.id_user "+
		"JOIN users frd ON frd.id = f.id_friend "+
		"LEFT JOIN user_in_room uir1 ON uir1.id_user = u.id "+
		"LEFT JOIN user_in_room uir2 ON uir2.id_user = frd.id AND uir1.id_room = uir2.id_room "+
		"LEFT JOIN rooms r ON r.id = uir1.id_room AND r.id = uir2.id_room "+
		"WHERE "+
		"u.username = $1 "+
		"AND f.status = 'ACCEPTED' "+
		"GROUP BY "+
		"f.id, frd.username, frd.is_online "+
		"ORDER BY interaction_at DESC, f.id DESC "+
		"LIMIT $2 OFFSET $3",
		username, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friends []Friend
	for rows.Next() {
		var friend Friend
		err := rows.Scan(&friend.IdRoom, &friend.Id, &friend.Username, &friend.IsOnline, &friend.InteractionAt)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func (r *Repository) GetListFriendsAfterTimestamp(ctx context.Context, username string, interactAt string, limit, offset int) ([]Friend, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := r.db.Query("SELECT "+
		"COALESCE(MAX(r.id), 0) AS id_room, "+
		"f.id AS id_friend, "+
		"frd.username AS friend_username, "+
		"frd.is_online AS is_online, "+
		"COALESCE(MAX(f.interaction_at), '2021-01-01 00:00:00') AS interaction_at "+
		"FROM "+
		"friends f "+
		"JOIN "+
		"users u ON u.id = f.id_user "+
		"JOIN users frd ON frd.id = f.id_friend "+
		"LEFT JOIN user_in_room uir1 ON uir1.id_user = u.id "+
		"LEFT JOIN user_in_room uir2 ON uir2.id_user = frd.id AND uir1.id_room = uir2.id_room "+
		"LEFT JOIN rooms r ON r.id = uir1.id_room AND r.id = uir2.id_room "+
		"WHERE "+
		"u.username = $1 "+
		"AND f.status = 'ACCEPTED' "+
		"AND f.interaction_at < $2 "+
		"GROUP BY "+
		"f.id, frd.username, frd.is_online "+
		"ORDER BY interaction_at DESC, f.id DESC "+
		"LIMIT $3 OFFSET $4",
		username, interactAt, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friends []Friend
	for rows.Next() {
		var friend Friend
		err := rows.Scan(&friend.IdRoom, &friend.Id, &friend.Username, &friend.IsOnline, &friend.InteractionAt)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func (r *Repository) GetListFriendsById(ctx context.Context, id string) ([]Friend, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := r.db.Query("SELECT "+
		"COALESCE(MAX(r.id), 0) AS id_room, "+
		"f.id AS id_friend, "+
		"frd.username AS friend_username, "+
		"COALESCE(MAX(f.interaction_at), '2021-01-01 00:00:00') AS interaction_at "+
		"FROM "+
		"friends f "+
		"JOIN "+
		"users u ON u.id = f.id_user "+
		"JOIN users frd ON frd.id = f.id_friend "+
		"LEFT JOIN user_in_room uir1 ON uir1.id_user = u.id "+
		"LEFT JOIN user_in_room uir2 ON uir2.id_user = frd.id AND uir1.id_room = uir2.id_room "+
		"LEFT JOIN rooms r ON r.id = uir1.id_room AND r.id = uir2.id_room "+
		"WHERE "+
		"u.id = $1 "+
		"AND f.status = 'ACCEPTED' "+
		"GROUP BY "+
		"f.id, frd.username "+
		"ORDER BY interaction_at DESC, f.id DESC ",
		id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friends []Friend
	for rows.Next() {
		var friend Friend
		err := rows.Scan(&friend.IdRoom, &friend.Id, &friend.Username, &friend.InteractionAt)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func (r *Repository) CreateChatRoom(ctx context.Context, tx *sql.Tx, name string) (string, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var id string
	err := tx.QueryRow("INSERT INTO rooms (name) VALUES ($1) RETURNING id", name).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) AddUserToRoom(ctx context.Context, tx *sql.Tx, idUser, idRoom, idRole string) error {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := tx.Exec("INSERT INTO user_in_room (id_user, id_room, id_role) VALUES ($1, $2, $3)", idUser, idRoom, idRole)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateInteraction(ctx context.Context, idUser, idFriend string) error {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec("UPDATE friends SET interaction_at = $1 WHERE id_user = $2 AND id_friend = $3", time.Now(), idUser, idFriend)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) IsExistChatRoom(ctx context.Context, tx *sql.Tx, idUser1, idUser2 string) (string, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var id string
	err := tx.QueryRow("SELECT "+
		"r.id AS id_room "+
		"FROM rooms r "+
		"JOIN user_in_room u1 ON u1.id_room = r.id "+
		"JOIN user_in_room u2 ON u2.id_room = r.id "+
		"WHERE u1.id_user = $1 "+
		"AND u2.id_user = $2", idUser1, idUser2).Scan(&id)
	if err != nil || id == "" {
		return "", nil
	}
	return id, nil
}

func (r *Repository) GetMessagesInRoom(ctx context.Context, idRoom string, limit, offset int) ([]Messages, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := r.db.Query("SELECT "+
		"m.id, "+
		"m.id_sender, "+
		"u.username as username, "+
		"m.id_receiver, "+
		"m.content, "+
		"m.create_at "+
		"FROM messages m "+
		"JOIN users u ON u.id = m.id_sender "+
		"WHERE m.id_receiver = $1 "+
		"ORDER BY m.id DESC "+
		"LIMIT $2 OFFSET $3", idRoom, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Messages
	for rows.Next() {
		var message Messages
		err := rows.Scan(&message.Id, &message.IdSender, &message.UsernameSender, &message.IdReceiver, &message.Content, &message.CreateAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *Repository) GetMessagesOlderThanId(ctx context.Context, idMsg, idRoom string, limit, offset int) ([]Messages, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := r.db.Query("SELECT m.id, m.id_sender, u.username as username, m.id_receiver, m.content, m.create_at "+
		"FROM messages m "+
		"JOIN users u ON u.id = m.id_sender "+
		"WHERE m.id < $1 AND m.id_receiver = $2 "+
		"ORDER BY m.id DESC "+
		"LIMIT $3 OFFSET $4", idMsg, idRoom, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Messages
	for rows.Next() {
		var message Messages
		err := rows.Scan(&message.Id, &message.IdSender, &message.UsernameSender, &message.IdReceiver, &message.Content, &message.CreateAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *Repository) InsertMessage(ctx context.Context, idSender, idReceiver, content string, createAt time.Time) (string, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var id string
	err := r.db.QueryRow("INSERT INTO messages (id_sender, id_receiver, content, create_at) VALUES ($1, $2, $3, $4) RETURNING id",
		idSender, idReceiver, content, createAt).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) InsertMessageIntoRoom(ctx context.Context, idSender, idReceiver, content string, createAt time.Time) (string, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var id string
	err := r.db.QueryRow("INSERT INTO messages (id_sender, id_receiver, content, create_at) VALUES ($1, $2, $3, $4) RETURNING id",
		idSender, idReceiver, content, createAt).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) CheckPermission(ctx context.Context, idUser, idRoom string) (bool, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var id string
	err := r.db.QueryRow("SELECT id FROM user_in_room WHERE id_user = $1 AND id_room = $2", idUser, idRoom).Scan(&id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) ChangeOnlineStatus(ctx context.Context, id string, status bool) error {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec("UPDATE users SET is_online = $1 WHERE id = $2", status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) InsertMessageWithReadState(ctx context.Context, idSender, idRoom, content string, createAt time.Time) error {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec(
		"WITH new_message AS ( "+
			"INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ($4, $1, $2, $3) "+
			"RETURNING id AS message_id) "+
			"INSERT INTO message_read_status (id_message, id_receiver, is_read, read_at) "+
			"SELECT "+
			"new_message.message_id, "+
			"user_in_room.id_user, "+
			"CASE WHEN user_in_room.id_user = $1 THEN TRUE ELSE FALSE END AS is_read, "+
			"$4 "+
			"FROM new_message "+
			"JOIN user_in_room ON user_in_room.id_room = $2",
		idSender, idRoom, content, createAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetListFriendAndMessage(ctx context.Context, username, interactionAt string, limit, offset int) ([]FriendListResponse, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := r.db.Query(
		"WITH friends_with_room AS "+
			"(SELECT "+
			"COALESCE(MAX(r.id), 0) AS id_room, "+
			"u.id AS id_user, "+
			"f.id AS id_friend, "+
			"frd.username AS friend_username, "+
			"frd.is_online AS is_online, COALESCE(MAX(f.interaction_at), '2021-01-01 00:00:00') AS interaction_at "+
			"FROM "+
			"friends f "+
			"JOIN users u ON u.id = f.id_user "+
			"JOIN users frd ON frd.id = f.id_friend "+
			"LEFT JOIN user_in_room uir1 ON uir1.id_user = u.id "+
			"LEFT JOIN user_in_room uir2 ON uir2.id_user = frd.id AND uir1.id_room = uir2.id_room "+
			"LEFT JOIN rooms r ON r.id = uir1.id_room AND r.id = uir2.id_room "+
			"WHERE "+
			"u.username = $1 AND f.status = 'ACCEPTED' and f.interaction_at < $2 "+
			"GROUP BY "+
			"f.id, frd.username, frd.is_online, u.id "+
			"ORDER BY "+
			"interaction_at DESC "+
			"LIMIT $3 OFFSET $4) "+
			"SELECT "+
			"fwr.id_room, COALESCE(max(m.id), 0) AS id_message, fwr.friend_username, fwr.id_friend, fwr.is_online, "+
			"COALESCE(m.is_read, TRUE) AS is_read, fwr.interaction_at "+
			"FROM "+
			"friends_with_room fwr "+
			"LEFT JOIN LATERAL "+
			"(SELECT m.id, mrs.is_read FROM messages m "+
			"LEFT JOIN message_read_status mrs ON mrs.id_message = m.id AND mrs.id_receiver = fwr.id_user "+
			"WHERE "+
			"m.id_receiver = fwr.id_room "+
			"ORDER BY m.id DESC LIMIT 1) AS m ON TRUE "+
			"GROUP BY fwr.id_room, fwr.friend_username, fwr.id_friend, m.is_read, fwr.interaction_at, is_online",
		username, interactionAt, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friends []FriendListResponse
	for rows.Next() {
		var friend FriendListResponse
		err := rows.Scan(&friend.IdRoom, &friend.IdMessage, &friend.Username, &friend.Id, &friend.IsOnline, &friend.IsRead, &friend.InteractionAt)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}
