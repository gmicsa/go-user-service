package test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
	"user-service/app"
	"user-service/health"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	service := app.New()
	service.Start()
	defer service.Stop()

	status := getHealthResponse(t)

	assert.Equal(t, "up", status.Status)
	assert.NotEmpty(t, status.Version)
	assert.NotEmpty(t, status.BuildDate)
	assert.NotEmpty(t, status.GitCommit)
}

func getHealthResponse(t *testing.T) health.HealthStatus {
	body := getOKResponse(t, "http://localhost:8080/health")

	var status health.HealthStatus
	err := json.Unmarshal(body, &status)
	require.NoError(t, err)

	return status
}

func getOKResponse(t *testing.T, url string) []byte {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(url)
	require.NoError(t, err)

	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err)
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return body
}
