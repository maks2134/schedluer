package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"schedluer/internal/service"
)

type ScheduleHandler struct {
	scheduleService service.ScheduleService
	logger          *logrus.Logger
}

func NewScheduleHandler(scheduleService service.ScheduleService, logger *logrus.Logger) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
		logger:          logger,
	}
}

// GetGroupSchedule получает расписание группы
// @Summary Получить расписание группы
// @Description Получает расписание группы по номеру группы
// @Tags schedule
// @Accept json
// @Produce json
// @Param groupNumber path string true "Номер группы"
// @Param useCache query bool false "Использовать кэш" default(true)
// @Success 200 {object} models.ScheduleResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/schedule/group/{groupNumber} [get]
func (h *ScheduleHandler) GetGroupSchedule(c *gin.Context) {
	groupNumber := c.Param("groupNumber")
	if groupNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group number is required"})
		return
	}

	useCache := c.DefaultQuery("useCache", "true") == "true"

	schedule, err := h.scheduleService.GetGroupSchedule(c.Request.Context(), groupNumber, useCache)
	if err != nil {
		h.logger.Errorf("Failed to get group schedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetEmployeeSchedule получает расписание преподавателя
// @Summary Получить расписание преподавателя
// @Description Получает расписание преподавателя по URL ID
// @Tags schedule
// @Accept json
// @Produce json
// @Param urlId path string true "URL ID преподавателя"
// @Param useCache query bool false "Использовать кэш" default(true)
// @Success 200 {object} models.ScheduleResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/schedule/employee/{urlId} [get]
func (h *ScheduleHandler) GetEmployeeSchedule(c *gin.Context) {
	urlID := c.Param("urlId")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url id is required"})
		return
	}

	useCache := c.DefaultQuery("useCache", "true") == "true"

	schedule, err := h.scheduleService.GetEmployeeSchedule(c.Request.Context(), urlID, useCache)
	if err != nil {
		h.logger.Errorf("Failed to get employee schedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// RefreshGroupSchedule обновляет расписание группы
// @Summary Обновить расписание группы
// @Description Принудительно обновляет расписание группы из API БГУИРа
// @Tags schedule
// @Accept json
// @Produce json
// @Param groupNumber path string true "Номер группы"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/schedule/group/{groupNumber}/refresh [post]
func (h *ScheduleHandler) RefreshGroupSchedule(c *gin.Context) {
	groupNumber := c.Param("groupNumber")
	if groupNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group number is required"})
		return
	}

	if err := h.scheduleService.RefreshGroupSchedule(c.Request.Context(), groupNumber); err != nil {
		h.logger.Errorf("Failed to refresh group schedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "schedule refreshed successfully"})
}

// RefreshEmployeeSchedule обновляет расписание преподавателя
// @Summary Обновить расписание преподавателя
// @Description Принудительно обновляет расписание преподавателя из API БГУИРа
// @Tags schedule
// @Accept json
// @Produce json
// @Param urlId path string true "URL ID преподавателя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/schedule/employee/{urlId}/refresh [post]
func (h *ScheduleHandler) RefreshEmployeeSchedule(c *gin.Context) {
	urlID := c.Param("urlId")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url id is required"})
		return
	}

	if err := h.scheduleService.RefreshEmployeeSchedule(c.Request.Context(), urlID); err != nil {
		h.logger.Errorf("Failed to refresh employee schedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "schedule refreshed successfully"})
}
