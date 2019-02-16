package internal

import (
	"github.com/template/be/internal/controller"
	"github.com/template/be/internal/repository"
	"github.com/template/be/internal/usecase"
)

// GetUsecase for project
func GetUsecase(config *Config, ipAddress string) *Usecase {

	// DATABASE
	// db := initDatabase(config.Database)

	// REDIS
	// rds := initRedis(config.Redis)

	// NSQ PRODUCER
	// nsqProducer := initNSQProducer(config.NSQd.Endpoint)

	// ELASTICSEARCH
	// elasticSearch := initElastic(config.Elastic)

	// REPOSITORY
	systemRepo := repository.NewSystem(ipAddress)

	// CONTROLLER
	systemService := controller.NewSystem(systemRepo)

	// USECASE

	systemUsecase := usecase.NewSystem(systemService)

	ucase := &Usecase{
		System: systemUsecase,
	}

	return ucase
}
