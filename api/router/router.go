package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/karlosdaniel451/go-rest-api-template/api/controller"
)

func Setup(app *fiber.App, taskController *controller.TaskController) {
	setupSwagger(app, "/docs/*")
	setupTaskRouter(app, taskController)
}

func setupSwagger(app *fiber.App, path string) {
	app.Get("/docs/*", swagger.HandlerDefault)
}
