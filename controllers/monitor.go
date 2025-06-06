package controllers

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sumukhj1219/betterstack/models"
	"github.com/sumukhj1219/betterstack/utils"
)

var (
	monitorLogs = []models.MonitorLogs{}
	logsMutex   = &sync.RWMutex{}
)

func GetMonitorLogs() []models.MonitorLogs {
	logsMutex.RLock()
	defer logsMutex.RUnlock()

	logsCopy := make([]models.MonitorLogs, len(monitorLogs))
	copy(logsCopy, monitorLogs)

	return logsCopy
}

func GetLogMutex() *sync.RWMutex {
	return logsMutex
}

func Monitor(ctx context.Context, url string) {
	ticker := time.NewTicker(10 * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped monitoring %s\n", url)
			return
		case <-ticker.C:
			currentTime := time.Now()
			resp, err := http.Get(url)

			if err != nil {
				fmt.Printf("Error checking %s: %v\n", url, err)
				payload := &utils.NewLogPayload{
					Url:    url,
					Time:   currentTime.Format(time.RFC3339),
					Status: fmt.Sprintf("Error: %v", err),
				}
				newLog := utils.NewLog(payload)

				logsMutex.Lock()
				monitorLogs = append(monitorLogs, *newLog)
				logsMutex.Unlock()
				continue
			}

			resp.Body.Close()
			statusCode := fmt.Sprintf("%d", resp.StatusCode)
			fmt.Printf("✅ %s is up [%d]\n", url, resp.StatusCode)
			payload := utils.NewLogPayload{
				Time:   currentTime.Format(time.RFC3339),
				Status: statusCode,
				Url:    url,
			}
			newLog := utils.NewLog(&payload)

			logsMutex.Lock()
			monitorLogs = append(monitorLogs, *newLog)
			logsMutex.Unlock()
		}
	}
}
