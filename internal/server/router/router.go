package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/svetlana-mel/event-task-planner/internal/server/handlers/auth"
	"github.com/svetlana-mel/event-task-planner/internal/server/handlers/event"
	"github.com/svetlana-mel/event-task-planner/internal/server/handlers/task"
)

// SetupRoutes setup data routes
// do not setup login and signup routes
func SetupRoutes(
	r chi.Router,
	eventHandler *event.Handler,
	taskHandler *task.Handler,
	authHandler *auth.Handler,
) {

	r.Route("/tasks", func(r chi.Router) {
		// r.Post("/", taskHandler.Create)
		// r.Get("/", taskHandler.GetAll)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Get("/", taskHandler.Get)
			// r.Delete("/", taskHandler.Delete)
			// r.Patch("/", taskHandler.Update)
		})
	})

	// r.Route("/events", func(r chi.Router) {
	// 	r.Post("/", eventHandler.Create)
	// 	r.Get("/", eventHandler.GetAll)
	// 	r.Route("/{eventID}", func(r chi.Router) {
	// 		r.Get("/", eventHandler.Get)
	// 		r.Delete("/", eventHandler.Delete)
	// 		r.Patch("/", eventHandler.Update)
	// 	})
	// })

	// r.Get("/today", eventHandler.GetAll)
}
