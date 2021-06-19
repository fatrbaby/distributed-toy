package main

import (
	"context"
	"fmt"
	"github.com/fatrbaby/distributed-toy/calendar"
	"github.com/fatrbaby/distributed-toy/logging"
	"github.com/fatrbaby/distributed-toy/registry"
	"github.com/fatrbaby/distributed-toy/service"
	"log"
)

func main() {
	host := "http://localhost:7700"

	svc := registry.Service{
		Name:             registry.ServiceCalendar,
		URL:              host,
		RequiredServices: []registry.ServiceName{registry.ServiceLogging},
		UpdateURL:        host + "/services",
		HeartbeatURL:     host + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(), svc, calendar.RegisterHandler)

	if err != nil {
		log.Fatalln(err)
	}

	prov, err := registry.GetProvider(registry.ServiceLogging)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Logging service found at: %s\n", prov)
		logging.UseClientLogger(prov, svc.Name)
	}

	<-ctx.Done()

	fmt.Println("shutdown services")
}
