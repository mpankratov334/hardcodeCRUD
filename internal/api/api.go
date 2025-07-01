package api

import (
	"TemplatestPGSQL/internal/api/middleware"
	"TemplatestPGSQL/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Routers struct {
	Service service.Service
}

func NewRouters(r *Routers, token string) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST, PUT, DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	apiGroup := app.Group("/v1", middleware.Authorization(token))

	apiGroup.Post("/tasks", r.Service.CreateTask)
	apiGroup.Post("/users", r.Service.CreateUser)
	apiGroup.Get("/tasks/all", r.Service.GetAllTasks)
	apiGroup.Get("/tasks/users/:id", r.Service.GetAllTasksByUserID)
	apiGroup.Delete("/tasks/:id", r.Service.DeleteTaskByID)
	apiGroup.Put("/tasks/:id/:status", r.Service.UpdateStatusByID)
	apiGroup.Get("tasks/users/:id/last", r.Service.GetLastTaskByUserID)
	apiGroup.Get("tasks/:id", r.Service.GetTaskByID)
	apiGroup.Get("tasks/users/:username", r.Service.GetTasksByUserName)
	return app
}
