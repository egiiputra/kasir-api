package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.TransactionRepository
}

func NewReportService(repo *repositories.TransactionRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetReport(startDate, endDate time.Time) (*models.ReportResponse, error) {
	return s.repo.GetReport(startDate, endDate)
}
