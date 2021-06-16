package main

import (
	"context"
	"fmt"
	"github.com/fatrbaby/distributed-toy/logger"
	"github.com/fatrbaby/distributed-toy/registry"
	"github.com/fatrbaby/distributed-toy/service"
	"log"
)

func main() {
	logger.Run("./toy.log")

	url := "http://localhost:4700"

	registrar := registry.Registrar{
		Name: registry.ServiceLogging,
		URL: url,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: fmt.Sprintf("%s/services", url),
	}

	ctx, err := service.Start(context.Background(), registrar, logger.RegisterHandlers)

	if err != nil {
		log.Fatalln(err)
	}

	<-ctx.Done()

	fmt.Println("shutdown services")
}
