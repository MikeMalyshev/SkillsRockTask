package service

import (
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
	ID          int    `json:"id" postgres:"id"`
	Title       string `json:"title" postgres:"title"`
	Description string `json:"description" postgres:"description"`
	Status      string `json:"status" postgres:"status"`
	CreateAt    string `json:"createAt" postgres:"createAt"`
	UpdateAt    string `json:"updateAt" postgres:"updateAt"`
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
