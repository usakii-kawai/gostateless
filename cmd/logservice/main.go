package main

import (
	"context"
	log "distributed/log"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./distributed.log")
	host, port := "127.0.0.1", "6379"
	serviceAddr := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: "Log Service",
		ServiceURL:  serviceAddr,
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)

	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
