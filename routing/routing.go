package routing

import (
	"net/http"
	"user-service/health"
	"user-service/users"
)

func ConfigureMainRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", health.Status)

	mux.HandleFunc("GET /users/{id}", users.GetByID)
	mux.HandleFunc("POST /users/{id}", users.Create)
}
