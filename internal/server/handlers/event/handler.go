package event

import (
	"log/slog"
	"net/http"

	// "github.com/svetlana-mel/event-task-planner/internal/models"
	"github.com/svetlana-mel/event-task-planner/internal/repository"
)

type Handler struct {
	Repo   repository.PlannerRepository
	Logger *slog.Logger
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {

}
