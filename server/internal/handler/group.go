package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"schedluer/internal/service"
)

type GroupHandler struct {
	groupService service.GroupService
	logger       *logrus.Logger
}

func NewGroupHandler(groupService service.GroupService, logger *logrus.Logger) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
		logger:       logger,
	}
}

// GetAllGroups получает список всех групп
// @Summary Получить список всех групп
// @Description Получает список всех групп из БГУИРа
// @Tags groups
// @Accept json
// @Produce json
// @Param useCache query bool false "Использовать кэш" default(true)
// @Success 200 {array} models.StudentGroupListItem
// @Failure 500 {object} map[string]string
// @Router /api/v1/groups [get]
func (h *GroupHandler) GetAllGroups(c *gin.Context) {
	useCache := c.DefaultQuery("useCache", "true") == "true"

	groups, err := h.groupService.GetAllGroups(c.Request.Context(), useCache)
	if err != nil {
		h.logger.Errorf("Failed to get groups: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetGroupByNumber получает группу по номеру
// @Summary Получить группу по номеру
// @Description Получает информацию о группе по номеру
// @Tags groups
// @Accept json
// @Produce json
// @Param groupNumber path string true "Номер группы"
// @Success 200 {object} models.StudentGroupListItem
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/groups/{groupNumber} [get]
func (h *GroupHandler) GetGroupByNumber(c *gin.Context) {
	groupNumber := c.Param("groupNumber")
	if groupNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group number is required"})
		return
	}

	group, err := h.groupService.GetGroupByNumber(c.Request.Context(), groupNumber)
	if err != nil {
		h.logger.Errorf("Failed to get group: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// RefreshGroups обновляет список всех групп
// @Summary Обновить список всех групп
// @Description Принудительно обновляет список всех групп из API БГУИРа
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/groups/refresh [post]
func (h *GroupHandler) RefreshGroups(c *gin.Context) {
	if err := h.groupService.RefreshGroups(c.Request.Context()); err != nil {
		h.logger.Errorf("Failed to refresh groups: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "groups refreshed successfully"})
}
