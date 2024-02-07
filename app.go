package octopus

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	// _ "github.com/mattn/go-sqlite3"
)

type HandlerFunc func(*Ctx)
type ErrorHandlerFunc func(*Ctx, int)

type Route struct {
	pattern  string
	handlers []HandlerFunc
	methods  map[string]bool
}

type app struct {
	routes []*Route
	// if an handler gets an error and  the a .OnErrorCode is called and the error code and the handler are passed as parameters to the a .OnErrorCode
	onErrorCode ErrorHandlerFunc
}

func New() *app {
	return &app{}
}

func (a *app) handle(pattern string, handlers []HandlerFunc, methods ...string) {
	methodsMap := make(map[string]bool)
	for _, method := range methods {
		methodsMap[method] = true
	}
	route := &Route{pattern: pattern, handlers: handlers, methods: methodsMap}
	a.routes = append(a.routes, route)
}

func (a *app) Use(handlers ...HandlerFunc) {
	for _, route := range a.routes {
		route.handlers = append(handlers, route.handlers...)
	}
}

// func (a *app) OnErrorCode(errorCode int, handler ErrorHandlerFunc) {
// 	a .onErrorCode = handler
// }
// ...

func (a *app) Static(path string, dir string) {
	fileServer := http.FileServer(http.Dir(dir))
	a.GET(path+"*", func(c *Ctx) {
		http.StripPrefix(path, fileServer).ServeHTTP(c.Response, c.Request)
	})
}

func (a *app) GET(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "GET")
}

func (a *app) PUT(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "PUT")
}

func (a *app) POST(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "POST")
}

func (a *app) DELETE(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "DELETE")
}

func (a *app) PATCH(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "PATCH")
}

func (a *app) OPTIONS(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "OPTIONS")
}

func (a *app) HEAD(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "HEAD")
}

func (a *app) Any(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD")
}

// func (a *app) OnErrorCode()

func (a *app) NotAllowed(c *Ctx) {
	http.Error(c.Response, "405 Method not allowed", http.StatusMethodNotAllowed)
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Ctx{Response: w, Request: r}
	for _, route := range a.routes {
		if strings.HasSuffix(route.pattern, "*") {
			if strings.HasPrefix(r.URL.Path, strings.TrimSuffix(route.pattern, "*")) {
				if route.methods[r.Method] {
					c.handlers = route.handlers
					c.Next()
					return
				} else {
					a.NotAllowed(c)
					return
				}
			}
		} else {
			if r.URL.Path == route.pattern {
				if route.methods[r.Method] {
					c.handlers = route.handlers
					c.Next()
					return
				} else {
					a.NotAllowed(c)
					return
				}
			}
		}
	}
	if a.onErrorCode != nil {
		a.onErrorCode(c, http.StatusNotFound)
	} else {
		http.NotFound(w, r)
	}
}

var wg sync.WaitGroup

func checkServer(addr string) {
	resp, err := http.Get("http://" + addr)
	if err != nil {
		log.Println("Server is not running")
	} else {
		defer resp.Body.Close()
		displayLaunchMessage(addr)
	}
}

func (a *app) Run(addr string) error {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(addr, a); err != nil {
			log.Fatal(err)
		}
	}()

	// Attendre que le serveur démarre
	time.Sleep(time.Second)

	// Vérifier si le serveur est en cours d'exécution
	checkServer(addr)

	wg.Wait()
	return nil
}

func displayLaunchMessage(addr string) {
	fmt.Println("*********************************************")
	fmt.Println("***************** Octopus *******************")
	fmt.Println("*********************************************")
	host, _ := os.Hostname()
	fmt.Printf("Hostname: %s\n", host)
	fmt.Printf("Listening on address: %s\n", addr)
	fmt.Println("*********************************************")
}
