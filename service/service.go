package service

import (
	"context"
	"fmt"
	"github.com/fatrbaby/distributed-toy/registry"
	"log"
	"net/http"
)

func Start(ctx context.Context, name, host, port string, handlerRegistrar func()) (context.Context, error) {
	handlerRegistrar()

	serviceName := registry.ServiceName(name)

	registrar := registry.Registrar{
		Name: serviceName,
		URL: fmt.Sprintf("http://%s:%s", host, port),
	}

	ctx = startServices(ctx, serviceName, host, port)
	err := registry.RegisterService(registrar)

	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startServices(ctx context.Context, name registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = ":" + port

	var url = fmt.Sprintf("http://%s:%s", host, port)

	go func() {
		log.Println(srv.ListenAndServe())

		if err := registry.ShutdownService(url); err != nil {
			log.Println(err)
		}

		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop. \n", name)
		var s string
		fmt.Scanln(&s)

		if err := registry.ShutdownService(url); err != nil {
			log.Println(err)
		}

		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
