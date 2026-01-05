package handler

import (
	"log"
	"send/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type LoadCSVHandler struct {
	loadCSVService *service.LoadCSVService
}

func NewLoadCSVHandler(loadCSVService *service.LoadCSVService) *LoadCSVHandler {
	return &LoadCSVHandler{
		loadCSVService: loadCSVService,
	}
}

func (h *LoadCSVHandler) LoadCSV(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return c.JSON(401, map[string]string{"error": "пользователь не авторизован"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "невалидный user_id"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, map[string]string{"error": "файл не найден"})
	}
	log.Println("Получен файл для загрузки:", file.Filename)

	src, err := file.Open()
	if err != nil {
		return c.JSON(500, map[string]string{"error": "не удалось открыть файл"})
	}
	defer src.Close()

	if err := h.loadCSVService.LoadCSV(src, userID); err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	log.Println("Файл успешно обработан:", file.Filename)

	return c.JSON(200, map[string]string{"message": "файл успешно обработан"})
}
