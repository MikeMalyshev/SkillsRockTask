package service

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type TodoService struct {
	fiberApp *fiber.App
	storage  Storage
}

type Storage interface {
	AddTask(t Task) error
	GetTasks() ([]Task, error)
	UpdateTask(t Task) error
	DeleteTask(id int) error
}

type Task struct {
	ID          int        `json:"id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	CreatedAt   *time.Time `json:"createAt"`
	UpdatedAt   *time.Time `json:"updateAt"`
}

func New(stor Storage) *TodoService {
	srvc := &TodoService{
		fiberApp: fiber.New(),
		storage:  stor,
	}

	srvc.fiberApp.Post("/tasks", srvc.createTaskHandler)
	srvc.fiberApp.Get("/tasks", srvc.getTasksHandler)
	srvc.fiberApp.Put("/tasks/:id", srvc.putTaskHandler)
	srvc.fiberApp.Delete("/tasks/:id", srvc.deleteTaskHandler)

	return srvc
}

func (service *TodoService) Start() error {
	return service.fiberApp.Listen(":3000")
}
