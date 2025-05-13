package health

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Variables set at build time using -ldflags
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

type HealthStatus struct {
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	GitCommit string `json:"gitCommit"`
	Status    string `json:"status"`
}

func Status(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	status := HealthStatus{
		Version:   Version,
		BuildDate: BuildDate,
		GitCommit: Commit,
		Status:    "up",
	}
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		slog.Error("error encoding health status", "error", err)
	}
}
