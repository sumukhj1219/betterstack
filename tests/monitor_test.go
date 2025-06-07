package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sumukhj1219/betterstack/controllers"
	"github.com/sumukhj1219/betterstack/models"
	"github.com/sumukhj1219/betterstack/routers"
)

func ginTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	dummyCtx, cancel := context.WithCancel(context.Background())

	defer cancel()

	return routers.SetupGinRouter(dummyCtx)
}

func TestPingEndPoint(t *testing.T) {
	router := ginTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d for /ping, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "pong" {
		t.Errorf("Expected body 'pong' for /ping, got '%s'", w.Body.String())
	}

}

func TestLogsEndPoint(t *testing.T) {
	router := ginTestRouter()

	testLogs := []models.MonitorLogs{
		{ID: "mock-1", Url: "http://mock1.com", Time: time.Now().Add(-5 * time.Minute).Format(time.RFC3339), Status: "up (200)"},
		{ID: "mock-2", Url: "http://mock2.com", Time: time.Now().Add(-1 * time.Minute).Format(time.RFC3339), Status: "down (500)"},
	}

	controllers.SetMonitorLogsForTest(testLogs)
	defer controllers.ResetMonitorLogsForTest()

	req, _ := http.NewRequest(http.MethodGet, "/logs", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d for /logs, got %d", http.StatusOK, w.Code)
	}

	var responseLogs []models.MonitorLogs
	err := json.Unmarshal(w.Body.Bytes(), &responseLogs)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body for /logs: %v", err)
	}

	if !reflect.DeepEqual(responseLogs, testLogs) {
		t.Errorf("Expected logs %v from /logs, got %v", testLogs, responseLogs)
	}
}
