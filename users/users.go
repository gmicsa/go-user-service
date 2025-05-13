package users

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"user-service/api"
)

type User struct {
	ID string `json:"id"`
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(User{ID: id})
	if err != nil {
		slog.Error("error encoding response", "error", err)
	}
}

func Create(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(23 * time.Millisecond)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	err := json.NewEncoder(w).Encode(api.Error{Message: "not yet implemented"})
	if err != nil {
		slog.Error("error encoding response", "error", err)
	}
}
