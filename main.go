package main

import (
	"fmt"
	handler "go-rest-starter/handler/user"
	"go-rest-starter/middleware"
	"go-rest-starter/server"
	"net/http"
)

func main() {
	apiServer := server.NewServer(":8080")

	v1 := server.NewRouter()

	v1.Use(middleware.RequestLoggerMiddleware)

	v1.Get("/ping", []server.Middleware{}, func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Pong!")
		if err != nil {
			return err
		}
		return nil
	})

	v1.Post("/user", []server.Middleware{}, handler.HandleCreateUser)
	apiServer.Mount("/v1", v1)
	//apiServer.Register(v1)
	apiServer.Run()
}
