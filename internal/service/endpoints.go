package service

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// createTaskHandler receive new task structure in request body and add it to storage
func (service *TodoService) createTaskHandler(c *fiber.Ctx) error {
	task := Task{}

	err := c.BodyParser(&task)
	if err != nil {
		return c.Status(400).SendString("Invalid task structure")
	}

	if service.storage == nil {
		return c.Status(500).SendString("Storage is not initialized")
	}

	err = service.storage.AddTask(task)
	if err != nil {
		return c.Status(500).SendString("Error while adding task")
	}
	return c.SendStatus(200)
}

// getTasksHandler returns all tasks in json format
func (service *TodoService) getTasksHandler(c *fiber.Ctx) error {
	tasks, err := service.storage.GetTasks()
	if err != nil {
		return c.Status(500).SendString("Error while getting tasks")
	}
	return c.JSON(tasks)
}

// PutTaskHandler updates task by id. Task structure should be in request body in json format
func (service *TodoService) putTaskHandler(c *fiber.Ctx) error {
	stringId := c.Params("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		return c.Status(400).SendString("Invalid id")
	}

	task := Task{}
	err = c.BodyParser(task)
	if err != nil {
		return c.Status(400).SendString("Invalid task structure")
	}
	task.ID = id

	err = service.storage.UpdateTask(task)
	if err != nil {
		return c.Status(404).SendString("Task not found")
	}
	return c.SendStatus(200)
}

// DeleteTaskHandler deletes task by id
func (service *TodoService) deleteTaskHandler(c *fiber.Ctx) error {
	stringId := c.Params("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		return c.Status(400).SendString("Invalid id")
	}
	err = service.storage.DeleteTask(id)
	if err != nil {
		return c.Status(404).SendString("Task not found")
	}

	return c.SendStatus(200)
}
