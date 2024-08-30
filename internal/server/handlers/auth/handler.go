package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	sl "github.com/svetlana-mel/event-task-planner/internal/lib/slog"
	auth_service "github.com/svetlana-mel/event-task-planner/internal/services/auth"
)

type AuthResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Token  string `json:"token,omitempty"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SignupRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Handler struct {
	Auth   auth_service.Auth
	Logger *slog.Logger
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "server.handlers.auth.Login"
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	log := h.Logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	var requestBody LoginRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Error("error decode request body", sl.AddErrorAtribute(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, &AuthResponse{
			Status: "Error",
			Error:  "bad request",
		})
		return
	}

	token, err := h.Auth.Login(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		log.Error("error login user", sl.AddErrorAtribute(err))
		if errors.Is(err, auth_service.ErrUserNotExists) || errors.Is(err, auth_service.ErrWrongPassword) {
			// inform frontend about type of client error
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &AuthResponse{
				Status: "Error",
				Error:  err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, &AuthResponse{
			Status: "Error",
			Error:  "error login user",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, &AuthResponse{
		Status: "OK",
		Token:  token,
	})
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	const op = "server.handlers.auth.Signup"
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	log := h.Logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	var requestBody SignupRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Error("error decode request body", sl.AddErrorAtribute(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, &AuthResponse{
			Status: "Error",
			Error:  "bad request",
		})
		return
	}

	if requestBody.Username == "" {
		requestBody.Username = "New user"
	}

	_, token, err := h.Auth.SignUp(
		ctx,
		requestBody.Username,
		requestBody.Email,
		requestBody.Password,
	)
	if err != nil {
		log.Error("error signup user", sl.AddErrorAtribute(err))
		if errors.Is(err, auth_service.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &AuthResponse{
				Status: "Error",
				Error:  err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, &AuthResponse{
			Status: "Error",
			Error:  "error signup user",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, &AuthResponse{
		Status: "OK",
		Token:  token,
	})
}
