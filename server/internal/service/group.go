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

type GroupService interface {
	GetAllGroups(ctx context.Context, useCache bool) ([]models.StudentGroupListItem, error)
	GetGroupByNumber(ctx context.Context, groupNumber string) (*models.StudentGroupListItem, error)
	RefreshGroups(ctx context.Context) error
}

type groupService struct {
	bsuirClient *bsuir.Client
	groupRepo   repository.GroupRepository
	logger      *logrus.Logger
}

func NewGroupService(bsuirClient *bsuir.Client, groupRepo repository.GroupRepository, logger *logrus.Logger) GroupService {
	return &groupService{
		bsuirClient: bsuirClient,
		groupRepo:   groupRepo,
		logger:      logger,
	}
}

func (s *groupService) GetAllGroups(ctx context.Context, useCache bool) ([]models.StudentGroupListItem, error) {
	if useCache {
		stored, err := s.groupRepo.GetAll(ctx)
		if err != nil {
			s.logger.Warnf("Failed to get groups from cache: %v", err)
		} else if len(stored) > 0 {
			result := make([]models.StudentGroupListItem, len(stored))
			for i, g := range stored {
				result[i] = g.GroupData
			}
			return result, nil
		}
	}

	groups, err := s.bsuirClient.GetAllGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get groups from BSUIR API: %w", err)
	}

	for _, g := range groups {
		stored := models.StoredGroup{
			ID:             primitive.NewObjectID(),
			BSUIRID:        g.ID,
			GroupData:      g,
			LastUpdateDate: time.Now().Format("02.01.2006"),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.groupRepo.Update(ctx, &stored); err != nil {
			s.logger.Warnf("Failed to save group %d to cache: %v", g.ID, err)
		}
	}

	return groups, nil
}

func (s *groupService) GetGroupByNumber(ctx context.Context, groupNumber string) (*models.StudentGroupListItem, error) {
	stored, err := s.groupRepo.GetByNumber(ctx, groupNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	if stored != nil {
		return &stored.GroupData, nil
	}

	groups, err := s.GetAllGroups(ctx, false)
	if err != nil {
		return nil, err
	}

	for _, g := range groups {
		if g.Name == groupNumber {
			return &g, nil
		}
	}

	return nil, fmt.Errorf("group not found: %s", groupNumber)
}

func (s *groupService) RefreshGroups(ctx context.Context) error {
	groups, err := s.bsuirClient.GetAllGroups()
	if err != nil {
		return fmt.Errorf("failed to get groups from BSUIR API: %w", err)
	}

	stored := make([]models.StoredGroup, len(groups))
	for i, g := range groups {
		stored[i] = models.StoredGroup{
			ID:             primitive.NewObjectID(),
			BSUIRID:        g.ID,
			GroupData:      g,
			LastUpdateDate: time.Now().Format("02.01.2006"),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
	}

	for _, g := range stored {
		if err := s.groupRepo.Update(ctx, &g); err != nil {
			s.logger.Warnf("Failed to update group %d: %v", g.BSUIRID, err)
		}
	}

	return nil
}
