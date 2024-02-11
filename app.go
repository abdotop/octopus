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

type App struct {
	routes           []*Route
	globalMiddleware []HandlerFunc
	// if an handler gets an error and  the a .OnErrorCode is called and the error code and the handler are passed as parameters to the a .OnErrorCode
	onErrorCode ErrorHandlerFunc
}

func New() *App {
	return &App{}
}

func (a *App) handle(pattern string, handlers []HandlerFunc, methods ...string) {
	methodsMap := make(map[string]bool)
	for _, method := range methods {
		methodsMap[method] = true
	}
	handlers = append(a.globalMiddleware, handlers...) // Ajoutez cette ligne
	route := &Route{pattern: pattern, handlers: handlers, methods: methodsMap}
	a.routes = append(a.routes, route)
}

func (a *App) Use(handlers ...HandlerFunc) {
	for _, route := range a.routes {
		route.handlers = append(handlers, route.handlers...)
	}
}

// func (a *App) OnErrorCode(errorCode int, handler ErrorHandlerFunc) {
// 	a .onErrorCode = handler
// }
// ...

func (a *App) Static(path string, dir string) {
	fileServer := http.FileServer(http.Dir(dir))
	a.GET(path+"*", func(c *Ctx) {
		http.StripPrefix(path, fileServer).ServeHTTP(c.Response, c.Request)
	})
}

func (a *App) GET(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "GET")
}

func (a *App) PUT(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "PUT")
}

func (a *App) POST(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "POST")
}

func (a *App) DELETE(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "DELETE")
}

func (a *App) PATCH(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "PATCH")
}

func (a *App) OPTIONS(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "OPTIONS")
}

func (a *App) HEAD(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "HEAD")
}

func (a *App) Any(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD")
}

// func (a *App) OnErrorCode()

func (a *App) NotAllowed(c *Ctx) {
	http.Error(c.Response, "405 Method not allowed", http.StatusMethodNotAllowed)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Ctx{Response: w, Request: r, handlers: nil, index: 0, Values: map[any]any{}, Context: r.Context()}

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

func (a *App) Run(addr string) error {
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
