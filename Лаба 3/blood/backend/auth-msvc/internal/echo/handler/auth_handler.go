package handler

import (
	"auth/internal/model"
	"auth/internal/service"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "невалидные данные"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	log.Printf("Регистрация пользователя: %s", req.Email)

	response, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		log.Printf("Ошибка регистрации: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	log.Printf("Пользователь успешно зарегистрирован: %s", req.Email)
	return c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "невалидные данные"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	log.Printf("Попытка входа: %s", req.Email)

	response, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		log.Printf("Ошибка входа: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	log.Printf("Пользователь успешно вошёл: %s", req.Email)
	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetMe(c echo.Context) error {
	userIDStr := c.Get("user_id")
	if userIDStr == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "пользователь не авторизован"})
	}

	userIDString, ok := userIDStr.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "невалидный user_id"})
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "невалидный user_id"})
	}

	user, err := h.authService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		log.Printf("Ошибка получения пользователя: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
