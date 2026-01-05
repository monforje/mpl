package model

import (
	"time"

	"github.com/google/uuid"
)

type Scan struct {
	UserID    uuid.UUID `json:"user_id"`
	FullName  string    `json:"full_name"`
	BirthDate time.Time `json:"birth_date"`
	Sex       string    `json:"sex"`

	Hemoglobin   float64 `json:"hemoglobin"`
	Erythrocytes float64 `json:"erythrocytes"`
	Hematocrit   float64 `json:"hematocrit"`

	MCV        float64 `json:"mcv"`
	Leukocytes float64 `json:"leukocytes"`

	Neutrophils float64 `json:"neutrophils"`
	Lymphocytes float64 `json:"lymphocytes"`
	Monocytes   float64 `json:"monocytes"`
	Eosinophils float64 `json:"eosinophils"`
	Basophils   float64 `json:"basophils"`

	Platelets float64 `json:"platelets"`
	MPV       float64 `json:"mpv"`
}
