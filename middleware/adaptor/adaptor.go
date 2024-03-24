package adaptor

import (
	"net/http"

	"github.com/abdotop/octopus"
)

func HTTPHandler(h http.Handler) octopus.HandlerFunc {
	return func(c *octopus.Ctx) {
		r, rok := c.Values.Get("request")
		w, wok := c.Values.Get("response")
		if rok && wok {
			r := r.(*http.Request)
			w := w.(http.ResponseWriter)
			h.ServeHTTP(w, r)
			c.Next()
		}
	}
}

func HTTPHandlerFunc(h http.HandlerFunc) octopus.HandlerFunc {
	return func(c *octopus.Ctx) {
		r, rok := c.Values.Get("request")
		w, wok := c.Values.Get("response")
		if rok && wok {
			r := r.(*http.Request)
			w := w.(http.ResponseWriter)
			h(w, r)
			c.Next()
		}
	}
}

// func HTTPMiddleware(h http.Handler) octopus.MiddlewareFunc {
// 	return func(next octopus.HandlerFunc) octopus.HandlerFunc {
// 		return func(c *octopus.Ctx) {
// 			r, rok := c.Values.Get("request")
// 			w, wok := c.Values.Get("response")
// 			if !rok || !wok {
// 				r := r.(*http.Request)
// 				w := w.(http.ResponseWriter)
// 				h.ServeHTTP(w, r)
// 				next(c)
// 			}
// 		}
// 	}
// }

func OctopusHandler(h octopus.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := octopus.NewCtx()
		c.Values.Set("request", r)
		c.Values.Set("response", w)
		h(c)
	})
}

func OctopusHandlerFunc(h octopus.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := octopus.NewCtx()
		c.Values.Set("request", r)
		c.Values.Set("response", w)
		h(c)
	}
}

// func OctopusApp(app *octopus.app) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 	}
// }
