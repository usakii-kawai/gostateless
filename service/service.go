package service

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var once sync.Once

func Start(ctx context.Context, host, port string, reg registry.Registration, registerHandlersFunc func()) (context.Context, error) {
	registerHandlersFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		log.Println(srv.ListenAndServe())

		once.Do(func() {
			fmt.Println("try to in delete:", fmt.Sprintf("http://%s:%s", host, port))
			err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
			if err != nil {
				log.Println(err)
			}
		})

		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop.\n", serviceName)
		var s string
		fmt.Scanln(&s)

		once.Do(func() {
			fmt.Println("try to in deletes:", fmt.Sprintf("http://%s:%s", host, port))
			err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
			if err != nil {
				log.Println(err)
			}
		})

		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
