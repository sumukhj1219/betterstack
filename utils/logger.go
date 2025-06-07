package utils

import (
	"log"
	"sync"

	"github.com/sumukhj1219/betterstack/models"
)

type NewLogPayload struct {
	Time   string
	Url    string
	Status string
}

func NewLog(payload *NewLogPayload) *models.MonitorLogs {
	return &models.MonitorLogs{
		Time:   payload.Time,
		Url:    payload.Url,
		Status: payload.Status,
	}
}

func ReadLogs(logsMutex *sync.RWMutex, monitorLogs []models.MonitorLogs) []models.MonitorLogs {
	logsMutex.RLock()
	defer logsMutex.RUnlock()

	logsCopy := make([]models.MonitorLogs, len(monitorLogs))
	copy(logsCopy, monitorLogs)

	return logsCopy
}

func PrintLogs(logsMutex *sync.RWMutex, monitorLogs []models.MonitorLogs) {
	logs := ReadLogs(logsMutex, monitorLogs)

	log.Printf("\n--- Current Monitoring Logs (%d entries) ---\n", len(logs))
	if len(logs) == 0 {
		log.Println("No logs collected yet.")
		return
	}

	for i, logEntry := range logs {
		log.Printf("[%d] URL: %s, Time: %s, Status: %s\n", i+1, logEntry.Url, logEntry.Time, logEntry.Status)
	}
	log.Println("------------------------------------------")
}
