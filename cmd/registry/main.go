package main

import (
	"context"
	"fmt"
	"github.com/fatrbaby/distributed-toy/registry"
	"log"
	"net/http"
)

func main() {
	registry.SetupRegistry()

	http.Handle("/services", registry.Server{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Println("Registry service started. Press any key to stop\n")
		var input string
		fmt.Scanln(&input)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()

	fmt.Println("Shutting down registry service.")
}
