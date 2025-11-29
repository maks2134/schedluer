package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"schedluer/internal/service"
)

type Router struct {
	scheduleHandler *ScheduleHandler
	groupHandler    *GroupHandler
	employeeHandler *EmployeeHandler
	favoriteHandler *FavoriteHandler
}

func NewRouter(scheduleService service.ScheduleService, groupService service.GroupService, employeeService service.EmployeeService, favoriteService service.FavoriteService, logger *logrus.Logger) *Router {
	return &Router{
		scheduleHandler: NewScheduleHandler(scheduleService, logger),
		groupHandler:    NewGroupHandler(groupService, logger),
		employeeHandler: NewEmployeeHandler(employeeService, logger),
		favoriteHandler: NewFavoriteHandler(favoriteService, logger),
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	api := engine.Group("/api/v1")

	schedule := api.Group("/schedule")
	{
		schedule.GET("/group/:groupNumber", r.scheduleHandler.GetGroupSchedule)
		schedule.POST("/group/:groupNumber/refresh", r.scheduleHandler.RefreshGroupSchedule)
		schedule.GET("/employee/:urlId", r.scheduleHandler.GetEmployeeSchedule)
		schedule.POST("/employee/:urlId/refresh", r.scheduleHandler.RefreshEmployeeSchedule)
	}

	groups := api.Group("/groups")
	{
		groups.GET("", r.groupHandler.GetAllGroups)
		groups.GET("/:groupNumber", r.groupHandler.GetGroupByNumber)
		groups.POST("/refresh", r.groupHandler.RefreshGroups)
	}

	employees := api.Group("/employees")
	{
		employees.GET("", r.employeeHandler.GetAllEmployees)
		employees.GET("/:urlId", r.employeeHandler.GetEmployeeByURLID)
		employees.POST("/refresh", r.employeeHandler.RefreshEmployees)
	}

	favorites := api.Group("/favorites")
	{
		favorites.GET("", r.favoriteHandler.GetAllFavorites)
		favorites.GET("/search", r.favoriteHandler.SearchFavorites)
		favorites.POST("/:groupNumber", r.favoriteHandler.AddFavorite)
		favorites.DELETE("/:groupNumber", r.favoriteHandler.RemoveFavorite)
		favorites.GET("/:groupNumber/check", r.favoriteHandler.IsFavorite)
	}
}
