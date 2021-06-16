package service

import (
	"context"
	"fmt"
	"github.com/fatrbaby/distributed-toy/registry"
	"log"
	"net/http"
	"net/url"
)

func Start(ctx context.Context, registrar registry.Registrar, handlerRegistrar func()) (context.Context, error) {
	handlerRegistrar()

	ctx = startServices(ctx, registrar)
	err := registry.RegisterService(registrar)

	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startServices(ctx context.Context, registrar registry.Registrar) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server

	parts, _ := url.Parse(registrar.URL)

	srv.Addr = ":" + parts.Port()

	go func() {
		log.Println(srv.ListenAndServe())

		if err := registry.ShutdownService(registrar.URL); err != nil {
			log.Println(err)
		}

		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop. \n", registrar.Name)
		var s string
		fmt.Scanln(&s)

		if err := registry.ShutdownService(registrar.URL); err != nil {
			log.Println(err)
		}

		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
