package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/karlosdaniel451/go-rest-api-template/api/controller"
	"github.com/karlosdaniel451/go-rest-api-template/api/middleware"
	"github.com/karlosdaniel451/go-rest-api-template/api/router"
	"github.com/karlosdaniel451/go-rest-api-template/db"
	_ "github.com/karlosdaniel451/go-rest-api-template/docs"
	"github.com/karlosdaniel451/go-rest-api-template/repository"
	"github.com/karlosdaniel451/go-rest-api-template/usecase"
)

var port = os.Getenv("API_PORT")

// @title Go REST API Template
// @version 0.0.1
// @description Template for a RESTful web service in Go with Fiber.
func main() {
	var _ usecase.TaskUseCase = usecase.TaskUseCaseImpl{}
	var _ repository.TaskRepository = repository.TaskRepositoryDB{}

	app := fiber.New(fiber.Config{
		AppName:           "Simple Go RESTful API with Fiber and GORM",
		EnablePrintRoutes: true,
	})

	middleware.Setup(app)

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	taskRepository := repository.NewTaskRepositoryDB(db.DB)
	taskUseCase := usecase.NewTaskUseCaseImpl(taskRepository)
	taskController := controller.NewTaskController(taskUseCase)

	router.Setup(app, &taskController)

	log.Fatal(app.Listen(":" + port))
}
