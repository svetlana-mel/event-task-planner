package task

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	sl "github.com/svetlana-mel/event-task-planner/internal/lib/slog"

	"github.com/svetlana-mel/event-task-planner/internal/models"
	"github.com/svetlana-mel/event-task-planner/internal/repository"
)

type GetResponse struct {
	Status string       `json:"status"`
	Error  string       `json:"error,omitempty"`
	Data   *models.Task `json:"data,omitempty"`
}

type Handler struct {
	Repo   repository.PlannerRepository
	Logger *slog.Logger
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	op := "server.handlers.task.Get"
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	log := h.Logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	taskIDstr := chi.URLParam(r, "taskID")

	taskID, err := strconv.Atoi(taskIDstr)
	if err != nil {
		log.Error("error parse taskID url value", sl.AddErrorAtribute(err))
		w.WriteHeader(400)
		render.JSON(w, r, &GetResponse{
			Status: "Error",
			Error:  "error parse taskID, taskID is not the number",
		})
		return
	}

	task, err := h.Repo.GetTask(ctx, uint64(taskID))
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotExists) {
			log.Error("error GetTask: task not exists", sl.AddErrorAtribute(err))
			w.WriteHeader(404)
			render.JSON(w, r, &GetResponse{
				Status: "Error",
				Error:  "error task not exists",
			})
			return
		}
		log.Error("error GetTask", sl.AddErrorAtribute(err))
		w.WriteHeader(500)
		render.JSON(w, r, &GetResponse{
			Status: "Error",
			Error:  "error getting task",
		})
		return
	}

	w.WriteHeader(200)
	render.JSON(w, r, &GetResponse{
		Status: "OK",
		Data:   task,
	})
}
