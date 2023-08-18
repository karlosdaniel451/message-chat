package controller

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/karlosdaniel451/go-rest-api-template/domain/model"
	"github.com/karlosdaniel451/go-rest-api-template/errs"
	"github.com/karlosdaniel451/go-rest-api-template/usecase"
)

type TaskController struct {
	taskUseCase usecase.TaskUseCase
}

func NewTaskController(taskUseCase usecase.TaskUseCase) TaskController {
	return TaskController{taskUseCase: taskUseCase}
}

// Create a new Task.
// @Description Create a new Task.
// @Summary Create a new Task.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body model.Task true "Task"
// @Success 201 {object} model.Task
// @Router /tasks [post]
func (controller TaskController) Create(c *fiber.Ctx) error {
	var newTask model.Task

	err := c.BodyParser(&newTask)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"detail": "invalid task data: " + err.Error(),
		})
	}

	newTaskAllData, err := controller.taskUseCase.Create(&newTask)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newTaskAllData)
}

// Delete a Task.
// @Description Delete a Task by its id.
// @Summary Delete a Task and, in case there is no Task with the given ID,
// returns a 404 not found error.
// @Tags Tasks
// @Produce json
// @Param id path int true "Id of the Task be deleted"
// @Success 204
// @Failure 404
// @Router /tasks/{id} [delete]
func (controller TaskController) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "invalid type: id of task should be an integer greater than 0",
		})
	}

	err = controller.taskUseCase.DeleteById(uint(id))
	if err != nil {
		if errors.As(err, &errs.NotFoundError{}) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Get a Task by its id.
// @Description Get a Task by its id.
// @Summary Get a Task by its id.
// @Tags Tasks
// @Produce json
// @Success 200 {object} model.Task
// @Failure 404
// @Router /tasks/{id} [get]
func (controller TaskController) GetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "invalid type: id of task should be an integer greater than 0",
		})
	}

	task, err := controller.taskUseCase.GetById(uint(id))
	if err != nil {
		if errors.As(err, &errs.NotFoundError{}) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	return c.JSON(task)
}

// Get all Tasks.
// @Description Get all Tasks
// @Summary Get all Tasks.
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {array} model.Task
// @Router /tasks [get]
func (controller TaskController) GetAll(c *fiber.Ctx) error {
	tasks, err := controller.taskUseCase.GetAll()
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	return c.JSON(tasks)
}
