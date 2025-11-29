package container

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"schedluer/internal/config"
	"schedluer/internal/handler"
	"schedluer/internal/repository"
	"schedluer/internal/service"
	"schedluer/pkg/bsuir"
	"schedluer/pkg/database"
)

type Container struct {
	Config *config.Config

	MongoDB *database.MongoDB

	BSUIRClient *bsuir.Client

	ScheduleRepo repository.ScheduleRepository
	GroupRepo    repository.GroupRepository
	EmployeeRepo repository.EmployeeRepository

	ScheduleService service.ScheduleService
	GroupService    service.GroupService
	EmployeeService service.EmployeeService

	Router *handler.Router

	Logger *logrus.Logger
}

func NewContainer(cfg *config.Config) (*Container, error) {
	logger := logrus.StandardLogger()

	mongoConfig := database.Config{
		URI:      cfg.MongoDB.URI,
		Database: cfg.MongoDB.Database,
		Timeout:  30 * time.Second, // Увеличиваем таймаут для Atlas
	}

	mongoDB, err := database.NewMongoDB(mongoConfig)
	if err != nil {
		return nil, err
	}

	bsuirClient := bsuir.NewClient(&cfg.BSUIRAPI)

	scheduleRepo := repository.NewScheduleRepository(mongoDB.Database)
	groupRepo := repository.NewGroupRepository(mongoDB.Database)
	employeeRepo := repository.NewEmployeeRepository(mongoDB.Database)

	scheduleService := service.NewScheduleService(bsuirClient, scheduleRepo, logger)
	groupService := service.NewGroupService(bsuirClient, groupRepo, logger)
	employeeService := service.NewEmployeeService(bsuirClient, employeeRepo, logger)

	apiRouter := handler.NewRouter(scheduleService, groupService, employeeService, logger)

	return &Container{
		Config:          cfg,
		MongoDB:         mongoDB,
		BSUIRClient:     bsuirClient,
		ScheduleRepo:    scheduleRepo,
		GroupRepo:       groupRepo,
		EmployeeRepo:    employeeRepo,
		ScheduleService: scheduleService,
		GroupService:    groupService,
		EmployeeService: employeeService,
		Router:          apiRouter,
		Logger:          logger,
	}, nil
}

func (c *Container) Close() error {
	if c.MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return c.MongoDB.Close(ctx)
	}
	return nil
}

func (c *Container) GetDatabase() *mongo.Database {
	return c.MongoDB.Database
}
