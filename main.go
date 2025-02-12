package main

import (
	"github.com/MikeMalyshev/SkillRocks/internal/postgres"
	"github.com/MikeMalyshev/SkillRocks/internal/service"
)

func main() {
	bd := postgres.New()

	service := service.New(bd)
	service.Start()
}
