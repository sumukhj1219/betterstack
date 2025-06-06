package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/sumukhj1219/betterstack/routers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println("Shutdown recieved successfully")
		cancel()
	}()

	
	routers.MonitorRouter(ctx, "https://hacktopus.tech")
}
