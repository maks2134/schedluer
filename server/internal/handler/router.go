package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"schedluer/internal/service"
)

// Router настраивает все маршруты API
type Router struct {
	scheduleHandler *ScheduleHandler
	groupHandler    *GroupHandler
	employeeHandler *EmployeeHandler
}

// NewRouter создает новый роутер
func NewRouter(
	scheduleService service.ScheduleService,
	groupService service.GroupService,
	employeeService service.EmployeeService,
	logger *logrus.Logger,
) *Router {
	return &Router{
		scheduleHandler: NewScheduleHandler(scheduleService, logger),
		groupHandler:    NewGroupHandler(groupService, logger),
		employeeHandler: NewEmployeeHandler(employeeService, logger),
	}
}

// SetupRoutes настраивает все маршруты
func (r *Router) SetupRoutes(engine *gin.Engine) {
	api := engine.Group("/api/v1")

	// Расписание
	schedule := api.Group("/schedule")
	{
		schedule.GET("/group/:groupNumber", r.scheduleHandler.GetGroupSchedule)
		schedule.POST("/group/:groupNumber/refresh", r.scheduleHandler.RefreshGroupSchedule)
		schedule.GET("/employee/:urlId", r.scheduleHandler.GetEmployeeSchedule)
		schedule.POST("/employee/:urlId/refresh", r.scheduleHandler.RefreshEmployeeSchedule)
	}

	// Группы
	groups := api.Group("/groups")
	{
		groups.GET("", r.groupHandler.GetAllGroups)
		groups.GET("/:groupNumber", r.groupHandler.GetGroupByNumber)
		groups.POST("/refresh", r.groupHandler.RefreshGroups)
	}

	// Преподаватели
	employees := api.Group("/employees")
	{
		employees.GET("", r.employeeHandler.GetAllEmployees)
		employees.GET("/:urlId", r.employeeHandler.GetEmployeeByURLID)
		employees.POST("/refresh", r.employeeHandler.RefreshEmployees)
	}
}
