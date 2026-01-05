package core

import (
	"context"
	"paste/internal/model"

	"github.com/google/uuid"
)

type ScanRepo interface {
	Create(ctx context.Context, scan *model.Scan) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Scan, error)
}
