package app

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGraph builds both generated graphs and exercises their managed cleanup.
func TestGraph(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APP_CONFIG", "../../config/config.example.yml")
	t.Setenv("APP_DATABASE__PATH", filepath.Join(tmp, "test.db"))
	t.Setenv("APP_LOG__OUTPUT", "file")
	t.Setenv("APP_LOG__PATH", filepath.Join(tmp, "test.log"))

	application, cleanupApp, err := InitializeApp("test")
	require.NoError(t, err)
	require.NotNil(t, application)
	require.NoError(t, application.migrator.Up(t.Context()))

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	res := httptest.NewRecorder()
	application.server.Handler.ServeHTTP(res, req)
	require.Equal(t, http.StatusOK, res.Code)
	require.Contains(t, res.Body.String(), `"version": "test"`)
	require.Contains(t, res.Body.String(), `"/users/{id}"`)

	require.NoError(t, cleanupApp())
	require.NoError(t, cleanupApp())

	management, cleanupCLI, err := InitializeCLI()
	require.NoError(t, err)
	require.NotNil(t, management)
	require.NoError(t, cleanupCLI())
	require.NoError(t, cleanupCLI())
}
