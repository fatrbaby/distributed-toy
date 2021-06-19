package service

import (
	"context"
	"fmt"
	"github.com/fatrbaby/distributed-toy/registry"
	"log"
	"net/http"
	"net/url"
)

func Start(ctx context.Context, service registry.Service, handlerRegistrar func()) (context.Context, error) {
	handlerRegistrar()

	ctx = startServices(ctx, service)
	err := registry.RegisterService(service)

	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startServices(ctx context.Context, service registry.Service) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server

	parts, _ := url.Parse(service.URL)

	srv.Addr = ":" + parts.Port()

	go func() {
		log.Println(srv.ListenAndServe())

		if err := registry.ShutdownService(service.URL); err != nil {
			log.Println(err)
		}

		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop. \n", service.Name)
		var s string
		fmt.Scanln(&s)

		if err := registry.ShutdownService(service.URL); err != nil {
			log.Println(err)
		}

		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
