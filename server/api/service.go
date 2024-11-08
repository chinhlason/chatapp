package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	r  *Repository
	rd *redis.Client
}

func NewService(r *Repository, rd *redis.Client) *Service {
	return &Service{r, rd}
}

func (s *Service) CreateUser(ctx context.Context, username, password string) error {
	user, err := s.r.GetUserByUsername(ctx, username)
	if err != nil {
		fmt.Println("error getting user by username: ", err)
	}
	fmt.Println("user: ", user)
	if user.Username != "" {
		return errors.New("username already exists")
	}
	err = s.r.InsertUser(ctx, username, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.r.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) VerifyUser(ctx context.Context, username, password string) error {
	user, err := s.r.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user.Password != password {
		return errors.New("invalid password")
	}
	return nil
}

func (s *Service) GetFriendRequests(ctx context.Context, userId string) ([]FriendRequest, error) {
	friendRequests, err := s.r.GetFriendRequests(ctx, userId)
	if err != nil {
		return nil, err
	}
	return friendRequests, nil
}

func (s *Service) SetToRedis(ctx context.Context, key, value string) error {
	err := s.rd.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SentFriendRequest(ctx context.Context, userId, friendId string) error {
	err := s.r.FriendRequest(ctx, userId, friendId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) AcceptFriendRequest(ctx context.Context, id string) (string, error) {
	tx, err := s.r.db.Begin()
	if err != nil {
		return "", err
	}
	idUser, idFriend, err := s.r.ChangeFriendRequestStatusAndReturnId(ctx, tx, id, "ACCEPTED")
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}
	idRoom, err := s.r.IsExistChatRoom(ctx, tx, idUser, idFriend)
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}
	if idRoom == "" {
		idRoom, err = s.r.CreateChatRoom(ctx, tx, "new chat room")
		if err != nil {
			_ = tx.Rollback()
			return "", err
		}
	}
	err = s.r.AddUserToRoom(ctx, tx, idUser, idRoom, "2")
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}
	err = s.r.AddUserToRoom(ctx, tx, idFriend, idRoom, "2")
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}
	return idRoom, tx.Commit()
}

func (s *Service) RejectFriendRequest(ctx context.Context, id string) error {
	err := s.r.ChangeFriendRequestStatus(ctx, id, "REJECTED")
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetFriends(ctx context.Context, username string, limit, offset int) ([]Friend, error) {
	friends, err := s.r.GetListFriends(ctx, username, limit, offset)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func (s *Service) UpdateInteraction(ctx context.Context, idUser, idFriend string) error {
	err := s.r.UpdateInteraction(ctx, idUser, idFriend)
	if err != nil {
		return err
	}
	return nil
}

//func (s *Service) CreateNewChat(ctx context.Context, userId, friendId string) error {
//
//}
