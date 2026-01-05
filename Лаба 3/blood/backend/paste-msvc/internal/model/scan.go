package model

import (
	"time"

	"github.com/google/uuid"
)

type Scan struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	FullName  string    `json:"full_name" db:"full_name"`
	BirthDate time.Time `json:"birth_date" db:"birth_date"`
	Sex       string    `json:"sex" db:"sex"`

	Hemoglobin   float64 `json:"hemoglobin" db:"hemoglobin"`
	Erythrocytes float64 `json:"erythrocytes" db:"erythrocytes"`
	Hematocrit   float64 `json:"hematocrit" db:"hematocrit"`

	MCV        float64 `json:"mcv" db:"mcv"`
	Leukocytes float64 `json:"leukocytes" db:"leukocytes"`

	Neutrophils float64 `json:"neutrophils" db:"neutrophils"`
	Lymphocytes float64 `json:"lymphocytes" db:"lymphocytes"`
	Monocytes   float64 `json:"monocytes" db:"monocytes"`
	Eosinophils float64 `json:"eosinophils" db:"eosinophils"`
	Basophils   float64 `json:"basophils" db:"basophils"`

	Platelets float64 `json:"platelets" db:"platelets"`
	MPV       float64 `json:"mpv" db:"mpv"`
}
