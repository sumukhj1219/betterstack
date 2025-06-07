package routers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/sumukhj1219/betterstack/controllers"
	"github.com/sumukhj1219/betterstack/utils"
)

func MonitorRouter(ctx context.Context, url string) {
	router := SetupGinRouter(ctx)

	go controllers.Monitor(ctx, url)

	go func() {
		printTicker := time.NewTicker(15 * time.Second)
		defer printTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Console log printer goroutine stopped.")
				return
			case <-printTicker.C:
				utils.PrintLogs(controllers.GetLogMutex(), controllers.GetMonitorLogs())
			}
		}

	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server ....")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	router.Run()
}
