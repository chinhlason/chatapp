package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Run(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello World")
}

func (h *Handler) RegisterUser(c echo.Context) error {
	var req RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
				Data:    nil,
			})
	}
	err := h.s.CreateUser(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}
	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "User created successfully",
		Data:    nil,
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
				Data:    nil,
			})
	}
	err := h.s.VerifyUser(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}
	user, _ := h.s.GetUserByUsername(c.Request().Context(), req.Username)
	token, err := GenJwt(req.Username, strconv.FormatInt(user.Id, 10))
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
	}
	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Login successful",
		Data: Token{
			Username: req.Username,
			Token:    token,
		},
	})
}

func (h *Handler) Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	username, _, _ := ValidateToken(token)
	err := h.s.SetToRedis(c.Request().Context(), username, token)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}

	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Logout successful",
		Data:    nil,
	})
}

func (h *Handler) GetFriendRequests(c echo.Context) error {
	userId := c.Param("userId")
	friendRequests, err := h.s.GetFriendRequests(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}
	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Friend requests",
		Data:    friendRequests,
	})
}

func (h *Handler) SendFriendRequest(c echo.Context) error {
	idUser := c.Get("id")
	idStr := idUser.(string)
	idFriend := c.Param("id")
	err := h.s.SentFriendRequest(c.Request().Context(), idStr, idFriend)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}
	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Friend request sent",
		Data:    time.Now(),
	})
}

func (h *Handler) AcceptFriendRequest(c echo.Context) error {
	id := c.Param("id")
	err := h.s.AcceptFriendRequest(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}
	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Friend request accepted",
		Data:    time.Now(),
	})
}

func (h *Handler) GetListFriends(c echo.Context) error {
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")
	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)
	offset := (pageInt - 1) * limitInt
	token := c.Request().Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	username, _, _ := ValidateToken(token)
	data, err := h.s.GetFriends(c.Request().Context(), username, limitInt, offset)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}
	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "List friends",
		Data:    data,
	})
}

func (h *Handler) UpdateInteraction(c echo.Context) error {
	idUser := c.Get("id")
	fmt.Println(idUser)
	idStr := idUser.(string)
	idFr := c.Param("id")
	err := h.s.UpdateInteraction(c.Request().Context(), idStr, idFr)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}
	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Interaction updated",
		Data:    time.Now(),
	})
}
