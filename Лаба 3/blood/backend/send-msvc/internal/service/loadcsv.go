package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"strings"

	"send/internal/core"
	"send/pkg/util"

	"github.com/google/uuid"
)

type LoadCSVService struct {
	producer core.Producer
}

func NewLoadCSVService(p core.Producer) *LoadCSVService {
	return &LoadCSVService{producer: p}
}

func (s *LoadCSVService) LoadCSV(file io.Reader, userID uuid.UUID) error {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	if len(records) < 2 {
		return nil
	}

	for i := 1; i < len(records); i++ {
		scan, err := util.CSVtoJSON(records[i])
		if err != nil {
			return err
		}

		scan.UserID = userID

		jsonData, err := json.Marshal(scan)
		if err != nil {
			return err
		}

		key := strings.TrimSpace(scan.FullName)
		if key == "" {
			key = "unknown"
		}

		err = s.producer.Publish(context.Background(), key, jsonData)
		if err != nil {
			return err
		}
	}

	return nil
}
