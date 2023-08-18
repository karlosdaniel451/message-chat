package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karlosdaniel451/go-rest-api-template/api/controller"
)

func setupTaskRouter(app *fiber.App, restController *controller.TaskController) {
	app.Post("/tasks", restController.Create)
	app.Get("/tasks", restController.GetAll)
	app.Get("/tasks/:id", restController.GetById)
	app.Delete("/tasks/:id", restController.Delete)
}
