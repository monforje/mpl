package handler

import (
	"log"
	"net/http"
	"paste/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ScanHandler struct {
	scanService *service.ScanService
}

func NewScanHandler(scanService *service.ScanService) *ScanHandler {
	return &ScanHandler{
		scanService: scanService,
	}
}

func (h *ScanHandler) GetMyScans(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "пользователь не авторизован"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "невалидный user_id"})
	}

	scans, err := h.scanService.GetScansByUserID(c.Request().Context(), userID)
	if err != nil {
		log.Printf("Ошибка при получении сканов: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось получить сканы"})
	}

	return c.JSON(http.StatusOK, scans)
}
