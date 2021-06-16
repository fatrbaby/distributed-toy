package main

import (
	"context"
	"fmt"
	"github.com/fatrbaby/distributed-toy/calendar"
	"github.com/fatrbaby/distributed-toy/logger"
	"github.com/fatrbaby/distributed-toy/registry"
	"github.com/fatrbaby/distributed-toy/service"
	"log"
)

func main() {
	host := "http://localhost:7700"

	registrar := registry.Registrar{
		Name:             registry.ServiceCalendar,
		URL:              host,
		RequiredServices: []registry.ServiceName{registry.ServiceLogging},
		ServiceUpdateURL: fmt.Sprintf("%s/services", host),
	}

	ctx, err := service.Start(context.Background(), registrar, calendar.RegisterHandler)

	if err != nil {
		log.Fatalln(err)
	}

	if prov, err := registry.GetProvider("log.service"); err == nil {
		fmt.Printf("Logging service found at: %s\n", prov)
		logger.SetLogger(prov, registrar.Name)
	} else {
		fmt.Println(err)
	}

	<-ctx.Done()

	fmt.Println("shutdown services")
}
