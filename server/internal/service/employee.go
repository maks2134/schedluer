package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"schedluer/internal/models"
	"schedluer/internal/repository"
	"schedluer/pkg/bsuir"
)

type EmployeeService interface {
	GetAllEmployees(ctx context.Context, useCache bool) ([]models.EmployeeListItem, error)
	GetEmployeeByURLID(ctx context.Context, urlID string) (*models.EmployeeListItem, error)
	RefreshEmployees(ctx context.Context) error
}

type employeeService struct {
	bsuirClient  *bsuir.Client
	employeeRepo repository.EmployeeRepository
	logger       *logrus.Logger
}

func NewEmployeeService(
	bsuirClient *bsuir.Client,
	employeeRepo repository.EmployeeRepository,
	logger *logrus.Logger,
) EmployeeService {
	return &employeeService{
		bsuirClient:  bsuirClient,
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func (s *employeeService) GetAllEmployees(ctx context.Context, useCache bool) ([]models.EmployeeListItem, error) {
	if useCache {
		stored, err := s.employeeRepo.GetAll(ctx)
		if err != nil {
			s.logger.Warnf("Failed to get employees from cache: %v", err)
		} else if len(stored) > 0 {
			result := make([]models.EmployeeListItem, len(stored))
			for i, e := range stored {
				result[i] = e.EmployeeData
			}
			return result, nil
		}
	}

	employees, err := s.bsuirClient.GetAllEmployees()
	if err != nil {
		return nil, fmt.Errorf("failed to get employees from BSUIR API: %w", err)
	}

	for _, e := range employees {
		stored := models.StoredEmployee{
			ID:             primitive.NewObjectID(),
			BSUIRID:        e.ID,
			URLID:          e.URLID,
			EmployeeData:   e,
			LastUpdateDate: time.Now().Format("02.01.2006"),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.employeeRepo.Update(ctx, &stored); err != nil {
			s.logger.Warnf("Failed to save employee %d to cache: %v", e.ID, err)
		}
	}

	return employees, nil
}

func (s *employeeService) GetEmployeeByURLID(ctx context.Context, urlID string) (*models.EmployeeListItem, error) {
	stored, err := s.employeeRepo.GetByURLID(ctx, urlID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}
	if stored != nil {
		return &stored.EmployeeData, nil
	}

	employees, err := s.GetAllEmployees(ctx, false)
	if err != nil {
		return nil, err
	}

	for _, e := range employees {
		if e.URLID == urlID {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("employee not found: %s", urlID)
}

func (s *employeeService) RefreshEmployees(ctx context.Context) error {
	employees, err := s.bsuirClient.GetAllEmployees()
	if err != nil {
		return fmt.Errorf("failed to get employees from BSUIR API: %w", err)
	}

	stored := make([]models.StoredEmployee, len(employees))
	for i, e := range employees {
		stored[i] = models.StoredEmployee{
			ID:             primitive.NewObjectID(),
			BSUIRID:        e.ID,
			URLID:          e.URLID,
			EmployeeData:   e,
			LastUpdateDate: time.Now().Format("02.01.2006"),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
	}

	for _, e := range stored {
		if err := s.employeeRepo.Update(ctx, &e); err != nil {
			s.logger.Warnf("Failed to update employee %d: %v", e.BSUIRID, err)
		}
	}

	return nil
}
