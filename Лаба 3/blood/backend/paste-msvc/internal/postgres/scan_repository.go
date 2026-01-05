package postgres

import (
	"context"
	"fmt"
	"paste/internal/model"

	"github.com/google/uuid"
)

type ScanRepository struct {
	pg *Postgres
}

func NewScanRepository(pg *Postgres) *ScanRepository {
	return &ScanRepository{pg: pg}
}

func (r *ScanRepository) Create(ctx context.Context, scan *model.Scan) error {
	query := `
		INSERT INTO scans (
			id, user_id, created_at,
			full_name, birth_date, sex,
			hemoglobin, erythrocytes, hematocrit, mcv, leukocytes,
			neutrophils, lymphocytes, monocytes, eosinophils, basophils,
			platelets, mpv
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)
	`

	_, err := r.pg.Pool.Exec(ctx, query,
		scan.ID, scan.UserID, scan.CreatedAt,
		scan.FullName, scan.BirthDate, scan.Sex,
		scan.Hemoglobin, scan.Erythrocytes, scan.Hematocrit, scan.MCV, scan.Leukocytes,
		scan.Neutrophils, scan.Lymphocytes, scan.Monocytes, scan.Eosinophils, scan.Basophils,
		scan.Platelets, scan.MPV,
	)
	if err != nil {
		return fmt.Errorf("не удалось вставить scan в БД: %w", err)
	}

	return nil
}

func (r *ScanRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Scan, error) {
	query := `
		SELECT 
			id, user_id, created_at,
			full_name, birth_date, sex,
			hemoglobin, erythrocytes, hematocrit, mcv, leukocytes,
			neutrophils, lymphocytes, monocytes, eosinophils, basophils,
			platelets, mpv
		FROM scans
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить сканы: %w", err)
	}
	defer rows.Close()

	var scans []*model.Scan
	for rows.Next() {
		scan := &model.Scan{}
		err := rows.Scan(
			&scan.ID, &scan.UserID, &scan.CreatedAt,
			&scan.FullName, &scan.BirthDate, &scan.Sex,
			&scan.Hemoglobin, &scan.Erythrocytes, &scan.Hematocrit, &scan.MCV, &scan.Leukocytes,
			&scan.Neutrophils, &scan.Lymphocytes, &scan.Monocytes, &scan.Eosinophils, &scan.Basophils,
			&scan.Platelets, &scan.MPV,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось сканировать строку: %w", err)
		}
		scans = append(scans, scan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return scans, nil
}
