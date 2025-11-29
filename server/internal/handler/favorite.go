package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"schedluer/internal/models"
	"schedluer/internal/service"
)

type FavoriteHandler struct {
	favoriteService service.FavoriteService
	logger          *logrus.Logger
}

func NewFavoriteHandler(favoriteService service.FavoriteService, logger *logrus.Logger) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
		logger:          logger,
	}
}

// GetAllFavorites получает все избранные группы
// @Summary      Получить все избранные группы
// @Description  Возвращает список всех избранных групп пользователя
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        user_id  query     string  true  "ID пользователя (пока используем 'default')"
// @Success      200      {array}   models.FavoriteGroup
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /favorites [get]
func (h *FavoriteHandler) GetAllFavorites(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "default")

	favorites, err := h.favoriteService.GetAllFavorites(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get favorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}

	// Убеждаемся, что возвращаем массив, а не null
	if favorites == nil {
		favorites = []models.FavoriteGroup{}
	}

	h.logger.Debugf("Returning %d favorites for user %s", len(favorites), userID)
	c.JSON(http.StatusOK, favorites)
}

// AddFavorite добавляет группу в избранное
// @Summary      Добавить группу в избранное
// @Description  Добавляет группу в список избранных для пользователя
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        user_id       query     string  true  "ID пользователя (пока используем 'default')"
// @Param        group_number  path      string  true  "Номер группы"
// @Success      200           {object}  map[string]string
// @Failure      400           {object}  map[string]string
// @Failure      500           {object}  map[string]string
// @Router       /favorites/{groupNumber} [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "default")
	groupNumber := c.Param("groupNumber")

	if groupNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_number is required"})
		return
	}

	err := h.favoriteService.AddFavorite(c.Request.Context(), userID, groupNumber)
	if err != nil {
		h.logger.Errorf("Failed to add favorite: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group added to favorites", "group_number": groupNumber})
}

// RemoveFavorite удаляет группу из избранного
// @Summary      Удалить группу из избранного
// @Description  Удаляет группу из списка избранных для пользователя
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        user_id       query     string  true  "ID пользователя (пока используем 'default')"
// @Param        group_number  path      string  true  "Номер группы"
// @Success      200           {object}  map[string]string
// @Failure      400           {object}  map[string]string
// @Failure      500           {object}  map[string]string
// @Router       /favorites/{groupNumber} [delete]
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "default")
	groupNumber := c.Param("groupNumber")

	if groupNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_number is required"})
		return
	}

	err := h.favoriteService.RemoveFavorite(c.Request.Context(), userID, groupNumber)
	if err != nil {
		h.logger.Errorf("Failed to remove favorite: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group removed from favorites", "group_number": groupNumber})
}

// IsFavorite проверяет, является ли группа избранной
// @Summary      Проверить, является ли группа избранной
// @Description  Проверяет, добавлена ли группа в избранное для пользователя
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        user_id       query     string  true  "ID пользователя (пока используем 'default')"
// @Param        group_number  path      string  true  "Номер группы"
// @Success      200           {object}  map[string]bool
// @Failure      400           {object}  map[string]string
// @Failure      500           {object}  map[string]string
// @Router       /favorites/{groupNumber}/check [get]
func (h *FavoriteHandler) IsFavorite(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "default")
	groupNumber := c.Param("groupNumber")

	if groupNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_number is required"})
		return
	}

	isFav, err := h.favoriteService.IsFavorite(c.Request.Context(), userID, groupNumber)
	if err != nil {
		h.logger.Errorf("Failed to check favorite: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFav})
}

// SearchFavorites ищет избранные группы по запросу
// @Summary      Поиск избранных групп
// @Description  Ищет избранные группы пользователя по номеру группы
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        user_id  query     string  true  "ID пользователя (пока используем 'default')"
// @Param        query    query     string  true  "Поисковый запрос"
// @Success      200      {array}   models.FavoriteGroup
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /favorites/search [get]
func (h *FavoriteHandler) SearchFavorites(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "default")
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	favorites, err := h.favoriteService.SearchFavorites(c.Request.Context(), userID, query)
	if err != nil {
		h.logger.Errorf("Failed to search favorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search favorites"})
		return
	}

	if favorites == nil {
		favorites = []models.FavoriteGroup{}
	}

	h.logger.Debugf("Found %d favorites for user %s with query %s", len(favorites), userID, query)
	c.JSON(http.StatusOK, favorites)
}
