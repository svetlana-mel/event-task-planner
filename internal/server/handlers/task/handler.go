package task

import (
	"net/http"

	"github.com/svetlana-mel/event-task-planner/internal/repository"
)

type Handler struct {
	Repo repository.PlannerRepository
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {

}
