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

type ScheduleService interface {
	GetGroupSchedule(ctx context.Context, groupNumber string, useCache bool) (*models.ScheduleResponse, error)
	GetEmployeeSchedule(ctx context.Context, urlID string, useCache bool) (*models.ScheduleResponse, error)
	RefreshGroupSchedule(ctx context.Context, groupNumber string) error
	RefreshEmployeeSchedule(ctx context.Context, urlID string) error
}

type scheduleService struct {
	bsuirClient  *bsuir.Client
	scheduleRepo repository.ScheduleRepository
	logger       *logrus.Logger
}

func NewScheduleService(bsuirClient *bsuir.Client, scheduleRepo repository.ScheduleRepository, logger *logrus.Logger) ScheduleService {
	return &scheduleService{
		bsuirClient:  bsuirClient,
		scheduleRepo: scheduleRepo,
		logger:       logger,
	}
}

func (s *scheduleService) GetGroupSchedule(ctx context.Context, groupNumber string, useCache bool) (*models.ScheduleResponse, error) {
	if useCache {
		stored, err := s.scheduleRepo.GetByGroupNumber(ctx, groupNumber)
		if err != nil {
			s.logger.Warnf("Failed to get schedule from cache: %v", err)
		} else if stored != nil {
			// Проверяем, не устарело ли расписание
			// Можно добавить логику проверки даты обновления
			return &stored.ScheduleData, nil
		}
	}

	// Получаем из API
	schedule, err := s.bsuirClient.GetGroupSchedule(groupNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule from BSUIR API: %w", err)
	}

	// Сохраняем в БД
	stored := &models.StoredSchedule{
		ID:             primitive.NewObjectID(),
		GroupNumber:    groupNumber,
		ScheduleData:   *schedule,
		LastUpdateDate: time.Now().Format("02.01.2006"),
	}

	if err := s.scheduleRepo.Update(ctx, stored); err != nil {
		s.logger.Warnf("Failed to save schedule to cache: %v", err)
	}

	return schedule, nil
}

func (s *scheduleService) GetEmployeeSchedule(ctx context.Context, urlID string, useCache bool) (*models.ScheduleResponse, error) {
	if useCache {
		stored, err := s.scheduleRepo.GetByEmployeeURLID(ctx, urlID)
		if err != nil {
			s.logger.Warnf("Failed to get schedule from cache: %v", err)
		} else if stored != nil {
			return &stored.ScheduleData, nil
		}
	}

	schedule, err := s.bsuirClient.GetEmployeeSchedule(urlID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule from BSUIR API: %w", err)
	}

	stored := &models.StoredSchedule{
		ID:             primitive.NewObjectID(),
		EmployeeURLID:  urlID,
		ScheduleData:   *schedule,
		LastUpdateDate: time.Now().Format("02.01.2006"),
	}

	if err := s.scheduleRepo.Update(ctx, stored); err != nil {
		s.logger.Warnf("Failed to save schedule to cache: %v", err)
	}

	return schedule, nil
}

func (s *scheduleService) RefreshGroupSchedule(ctx context.Context, groupNumber string) error {
	schedule, err := s.bsuirClient.GetGroupSchedule(groupNumber)
	if err != nil {
		return fmt.Errorf("failed to get schedule from BSUIR API: %w", err)
	}

	stored := &models.StoredSchedule{
		ID:             primitive.NewObjectID(),
		GroupNumber:    groupNumber,
		ScheduleData:   *schedule,
		LastUpdateDate: time.Now().Format("02.01.2006"),
	}

	return s.scheduleRepo.Update(ctx, stored)
}

func (s *scheduleService) RefreshEmployeeSchedule(ctx context.Context, urlID string) error {
	schedule, err := s.bsuirClient.GetEmployeeSchedule(urlID)
	if err != nil {
		return fmt.Errorf("failed to get schedule from BSUIR API: %w", err)
	}

	stored := &models.StoredSchedule{
		ID:             primitive.NewObjectID(),
		EmployeeURLID:  urlID,
		ScheduleData:   *schedule,
		LastUpdateDate: time.Now().Format("02.01.2006"),
	}

	return s.scheduleRepo.Update(ctx, stored)
}
