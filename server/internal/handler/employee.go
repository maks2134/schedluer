package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"schedluer/internal/service"

	_ "schedluer/internal/models" // для Swagger документации
)

type EmployeeHandler struct {
	employeeService service.EmployeeService
	logger          *logrus.Logger
}

func NewEmployeeHandler(employeeService service.EmployeeService, logger *logrus.Logger) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
		logger:          logger,
	}
}

// GetAllEmployees получает список всех преподавателей
// @Summary Получить список всех преподавателей
// @Description Получает список всех преподавателей из БГУИРа
// @Tags employees
// @Accept json
// @Produce json
// @Param useCache query bool false "Использовать кэш" default(true)
// @Success 200 {array} models.EmployeeListItem
// @Failure 500 {object} map[string]string
// @Router /api/v1/employees [get]
func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	useCache := c.DefaultQuery("useCache", "true") == "true"

	employees, err := h.employeeService.GetAllEmployees(c.Request.Context(), useCache)
	if err != nil {
		h.logger.Errorf("Failed to get employees: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}

// GetEmployeeByURLID получает преподавателя по URL ID
// @Summary Получить преподавателя по URL ID
// @Description Получает информацию о преподавателе по URL ID
// @Tags employees
// @Accept json
// @Produce json
// @Param urlId path string true "URL ID преподавателя"
// @Success 200 {object} models.EmployeeListItem
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/employees/{urlId} [get]
func (h *EmployeeHandler) GetEmployeeByURLID(c *gin.Context) {
	urlID := c.Param("urlId")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url id is required"})
		return
	}

	employee, err := h.employeeService.GetEmployeeByURLID(c.Request.Context(), urlID)
	if err != nil {
		h.logger.Errorf("Failed to get employee: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// RefreshEmployees обновляет список всех преподавателей
// @Summary Обновить список всех преподавателей
// @Description Принудительно обновляет список всех преподавателей из API БГУИРа
// @Tags employees
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/employees/refresh [post]
func (h *EmployeeHandler) RefreshEmployees(c *gin.Context) {
	if err := h.employeeService.RefreshEmployees(c.Request.Context()); err != nil {
		h.logger.Errorf("Failed to refresh employees: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employees refreshed successfully"})
}
