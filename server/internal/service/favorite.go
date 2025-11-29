package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"schedluer/internal/models"
	"schedluer/internal/repository"
)

type FavoriteService interface {
	GetAllFavorites(ctx context.Context, userID string) ([]models.FavoriteGroup, error)
	AddFavorite(ctx context.Context, userID string, groupNumber string) error
	RemoveFavorite(ctx context.Context, userID string, groupNumber string) error
	IsFavorite(ctx context.Context, userID string, groupNumber string) (bool, error)
	GetFavoriteGroupNumbers(ctx context.Context, userID string) ([]string, error)
}

type favoriteService struct {
	favoriteRepo repository.FavoriteRepository
	logger       *logrus.Logger
}

func NewFavoriteService(favoriteRepo repository.FavoriteRepository, logger *logrus.Logger) FavoriteService {
	return &favoriteService{
		favoriteRepo: favoriteRepo,
		logger:       logger,
	}
}

func (s *favoriteService) GetAllFavorites(ctx context.Context, userID string) ([]models.FavoriteGroup, error) {
	return s.favoriteRepo.GetAll(ctx, userID)
}

func (s *favoriteService) AddFavorite(ctx context.Context, userID string, groupNumber string) error {
	favorite := &models.FavoriteGroup{
		ID:          primitive.NewObjectID(),
		GroupNumber: groupNumber,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.favoriteRepo.Add(ctx, favorite)
}

func (s *favoriteService) RemoveFavorite(ctx context.Context, userID string, groupNumber string) error {
	return s.favoriteRepo.Delete(ctx, userID, groupNumber)
}

func (s *favoriteService) IsFavorite(ctx context.Context, userID string, groupNumber string) (bool, error) {
	return s.favoriteRepo.IsFavorite(ctx, userID, groupNumber)
}

func (s *favoriteService) GetFavoriteGroupNumbers(ctx context.Context, userID string) ([]string, error) {
	favorites, err := s.favoriteRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorites: %w", err)
	}

	groupNumbers := make([]string, len(favorites))
	for i, fav := range favorites {
		groupNumbers[i] = fav.GroupNumber
	}

	return groupNumbers, nil
}
