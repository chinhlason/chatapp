package api

import (
	"context"
	"errors"
	"fmt"
)

type Service struct {
	r *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r}
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
