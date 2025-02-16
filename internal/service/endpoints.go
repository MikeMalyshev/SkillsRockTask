package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// createTaskHandler receive new task structure in request body and add it to storage
func (service *TodoService) createTaskHandler(c *fiber.Ctx) error {
	task := Task{}

	FuncErr := func(err error) error {
		return fmt.Errorf("createTaskHandler: %v", err)
	}

	err := c.BodyParser(&task)
	if err != nil {
		c.Status(400).SendString("Invalid task structure")
		return FuncErr(err)
	}

	if service.storage == nil {
		c.Status(500).SendString("Storage is not initialized")
		return FuncErr(err)
	}

	err = service.storage.AddTask(task)
	if err != nil {
		c.Status(500).SendString("Error while adding task")
		return FuncErr(err)
	}
	return c.SendStatus(200)
}

// getTasksHandler returns all tasks in json format
func (service *TodoService) getTasksHandler(c *fiber.Ctx) error {
	tasks, err := service.storage.GetTasks()
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Error while getting tasks: " + err.Error())
	}
	return c.JSON(tasks)
}

// PutTaskHandler updates task by id. Task structure should be in request body in json format
func (service *TodoService) putTaskHandler(c *fiber.Ctx) error {
	FuncErr := func(err error) error {
		return fmt.Errorf("putTaskHandler: %v", err)
	}

	stringId := c.Params("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		c.Status(400).SendString("Invalid id")
		return FuncErr(err)
	}

	task := Task{}
	err = c.BodyParser(&task)
	if err != nil {
		c.Status(400).SendString("Invalid task structure")
		return FuncErr(err)
	}
	task.ID = id

	err = service.storage.UpdateTask(task)
	if err != nil {
		c.Status(404).SendString("Task not found")
		return FuncErr(err)

	}
	return c.SendStatus(200)
}

// DeleteTaskHandler deletes task by id
func (service *TodoService) deleteTaskHandler(c *fiber.Ctx) error {
	FuncErr := func(err error) error {
		return fmt.Errorf("deleteTaskHandler: %v", err)
	}

	stringId := c.Params("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		c.Status(400).SendString("Invalid id")
		return FuncErr(err)
	}
	err = service.storage.DeleteTask(id)
	if err != nil {
		c.Status(404).SendString("Task not found")
		return FuncErr(err)
	}
	return c.SendStatus(200)
}
