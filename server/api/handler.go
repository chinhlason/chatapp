package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
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
