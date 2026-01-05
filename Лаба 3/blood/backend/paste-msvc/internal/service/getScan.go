package service

import (
	"context"
	"paste/internal/core"
	"paste/internal/model"

	"github.com/google/uuid"
)

type ScanService struct {
	scanRepo core.ScanRepo
}

func NewScanService(scanRepo core.ScanRepo) *ScanService {
	return &ScanService{
		scanRepo: scanRepo,
	}
}

func (s *ScanService) CreateScan(ctx context.Context, scan *model.Scan) error {
	return s.scanRepo.Create(ctx, scan)
}

func (s *ScanService) GetScansByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Scan, error) {
	return s.scanRepo.GetByUserID(ctx, userID)
}
