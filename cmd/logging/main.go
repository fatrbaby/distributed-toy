package main

import (
	"context"
	"fmt"
	"github.com/fatrbaby/distributed-toy/logging"
	"github.com/fatrbaby/distributed-toy/registry"
	"github.com/fatrbaby/distributed-toy/service"
	"log"
)

func main() {
	logging.Run("./toy.log")

	url := "http://localhost:4700"

	svc := registry.Service{
		Name:             registry.ServiceLogging,
		URL:              url,
		RequiredServices: make([]registry.ServiceName, 0),
		UpdateURL:        fmt.Sprintf("%s/services", url),
	}

	ctx, err := service.Start(context.Background(), svc, logging.RegisterHandlers)

	if err != nil {
		log.Fatalln(err)
	}

	<-ctx.Done()

	fmt.Println("shutdown services")
}
