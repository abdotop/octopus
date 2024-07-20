package main

import (
	"net/http"
	"time"

	"github.com/abdotop/octopus"
	"github.com/abdotop/octopus/middleware/adaptor"
	"github.com/abdotop/octopus/middleware/cors"
	"github.com/abdotop/octopus/sse"
	// "github.com/abdotop/octopus/middleware/cor"
)

func main() {
	app := octopus.New()

	app.Use(cors.New(cors.Config{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
		ExposedHeaders:   []string{},
		MaxAge:           86400,
	}))
	app.Get("/", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	app.Get("/middleware", adaptor.HTTPHandler(next))

	app.Get("/sse", func(c *octopus.Ctx) {
		conn, err := sse.ConnFrom(c, new(sse.ClientConfig).Default())
		if err != nil {
			c.Status(octopus.StatusInternalServerError).JSON(octopus.Map{"error": err.Error()})
			return
		}

		store, err := c.AppStore()
		if err != nil {
			c.Status(octopus.StatusInternalServerError).JSON(octopus.Map{"error": err.Error()})
			return
		}

		store.Set("sse:"+conn.ID(), conn)

		notify := conn.Done()
		for {
			select {
			case <-notify:
				println("Close the connection when the client disconnects")
				store.Delete(conn.ID())
				return // Close the connection when the client disconnects
			default:
				// Send an event
				conn.SendJSON(octopus.Map{
					"id":   conn.ID(),
					"time": time.Now(),
				})
			}

			time.Sleep(5 * time.Second) // Simulate some delay
		}
	})

	app.Post("/getsse", func(c *octopus.Ctx) {
		type res struct {
			ID string
		}
		r := new(res)
		if c.BodyParser(r) != nil {
			c.Status(octopus.StatusBadRequest)
		}

		store, err := c.AppStore()
		if err != nil {
			c.Status(octopus.StatusInternalServerError).JSON(octopus.Map{"error": err.Error()})
			return
		}

		value, ok := store.Get("sse:" + r.ID)

		if !ok {
			c.Status(octopus.StatusInternalServerError).JSON(octopus.Map{"error": "no connection found with ID " + r.ID})
			return
		}
		conn, ok := value.(*sse.Conn)
		if !ok {
			c.Status(octopus.StatusInternalServerError).JSON(octopus.Map{"error": "failed to convert the value to *sse.Conn"})
			return
		}
		conn.Event("myEventType").SendText("Get test ok")
	})

	app.Post("/deletesse", func(c *octopus.Ctx) {
		type res struct {
			ID string
		}
		r := new(res)
		if c.BodyParser(r) != nil {
			c.Status(octopus.StatusBadRequest)
		}
		store, err := c.AppStore()
		if err != nil {
			c.Status(octopus.StatusInternalServerError).JSON(octopus.Map{"error": err.Error()})
			return
		}

		value, ok := store.Get("sse:" + r.ID)

		if !ok {
			c.Status(octopus.StatusInternalServerError).JSON(octopus.Map{"error": "no connection found with ID " + r.ID})
			return
		}
		conn, ok := value.(*sse.Conn)
		if !ok {
			c.Status(octopus.StatusInternalServerError).JSON(octopus.Map{"error": "failed to convert the value to *sse.Conn"})
			return
		}
		conn.Close()
	})

	app.Run(":8089")
}
