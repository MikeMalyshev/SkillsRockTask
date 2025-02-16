package main

import (
	"github.com/MikeMalyshev/SkillRocks/internal/postgres"
	"github.com/MikeMalyshev/SkillRocks/internal/service"
)

func main() {
	storage := postgres.New()
	if storage == nil {
		panic("failed to initialize database")
	}
	service := service.New(storage)
	service.Start()
}
