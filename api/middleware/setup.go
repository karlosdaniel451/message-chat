package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup(app *fiber.App) {
	loggerConfig := logger.Config{
		Next: nil,
		Done: nil,
		Format: "[${time}] - ${ip}:${port} - ${ua} - \"${method} ${path}\" - ${status} " +
			"- ${latency}\n",
		TimeFormat:    time.RFC3339Nano,
		TimeZone:      "UTC",
		TimeInterval:  500 * time.Millisecond,
		Output:        os.Stdout,
		DisableColors: true,
	}

	app.Use(logger.New(loggerConfig))
}
