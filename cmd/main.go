package main

import (
	"net/http"

	"github.com/abdotop/octopus"
	"github.com/abdotop/octopus/middleware/adaptor"
)

func main() {
	app := octopus.New()

	app.Get("/", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	app.Get("/middleware", adaptor.HTTPHandler(next))

	app.Run(":8089")
}
