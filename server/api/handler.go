package api

import (
	"context"
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
		Data: LoginResponse{
			Id:       user.Id,
			Username: user.Username,
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
	idRoom, err := h.s.AcceptFriendRequest(c.Request().Context(), id)
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
		Data:    idRoom,
	})
}

func (h *Handler) AcceptFriendRequestTest(ctx context.Context, id string) (string, error) {
	idRoom, err := h.s.AcceptFriendRequest(ctx, id)
	if err != nil {
		return "", err
	}
	return idRoom, nil
}

func (h *Handler) GetListFriends(c echo.Context) error {
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")
	interactAt := c.QueryParam("time")
	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)
	offset := (pageInt - 1) * limitInt
	token := c.Request().Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	username, _, _ := ValidateToken(token)
	var data []Friend
	var err error
	if interactAt != "" {
		data, err = h.s.GetFriendsAfterTimestamp(c.Request().Context(), username, interactAt, limitInt, offset)
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				Response{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
					Data:    nil,
				})
		}
	} else {
		data, err = h.s.GetFriends(c.Request().Context(), username, limitInt, offset)
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				Response{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
					Data:    nil,
				})
		}
	}
	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "List friends",
		Data:    data,
	})
}

func (h *Handler) UpdateInteraction(c echo.Context) error {
	idRoom := c.Param("id")
	err := h.s.UpdateInteraction(c.Request().Context(), idRoom)
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

func (h *Handler) GetMessages(c echo.Context) error {
	idRoom := c.Param("id_room")
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")
	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)
	offset := (pageInt - 1) * limitInt
	idUser := c.Get("id")
	legit, err := h.s.CheckPermissionInRoom(c.Request().Context(), idUser.(string), idRoom)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}
	if !legit {
		return c.JSON(http.StatusUnauthorized,
			Response{
				Code:    http.StatusUnauthorized,
				Message: "You are not allowed to access this room",
				Data:    nil,
			})
	}
	messages, err := h.s.GetMessagesInRoom(c.Request().Context(), idRoom, limitInt, offset)
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
		Message: "Messages",
		Data:    messages,
	})
}

func (h *Handler) GetMessagesOlder(c echo.Context) error {
	idRoom := c.Param("id_room")
	idOldestMessage := c.Param("id_msg")
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")
	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)
	offset := (pageInt - 1) * limitInt
	idUser := c.Get("id")
	legit, err := h.s.CheckPermissionInRoom(c.Request().Context(), idUser.(string), idRoom)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
	}
	if !legit {
		return c.JSON(http.StatusUnauthorized,
			Response{
				Code:    http.StatusUnauthorized,
				Message: "You are not allowed to access this room",
				Data:    nil,
			})
	}
	messages, err := h.s.GetMessagesOlderThanID(c.Request().Context(), idRoom, idOldestMessage, limitInt, offset)
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
		Message: "Messages",
		Data:    messages,
	})
}

func (h *Handler) GetListFriendAndMessage(c echo.Context) error {
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")
	interactAt := c.QueryParam("time")
	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)
	offset := (pageInt - 1) * limitInt
	token := c.Request().Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	username, _, _ := ValidateToken(token)
	if interactAt == "" {
		interactAt = time.Now().Add(time.Hour).Format("2006-01-02 15:04:05")
	}
	fmt.Println(interactAt, username)
	data, err := h.s.GetListFriendAndMessage(c.Request().Context(), username, interactAt, limitInt, offset)
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
