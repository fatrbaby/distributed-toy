package main

import (
	"context"
	"fmt"
	"github.com/fatrbaby/distributed-toy/logger"
	"github.com/fatrbaby/distributed-toy/service"
	"log"
)

func main() {
	logger.Run("./destination.log")
	host, port := "localhost", "4700"

	ctx, err := service.Start(context.Background(), "log.service", host, port, logger.RegisterHandlers)

	if err != nil {
		log.Fatalln(err)
	}

	<-ctx.Done()

	fmt.Println("shutdown services")
}
