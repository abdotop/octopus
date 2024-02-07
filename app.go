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

type HandlerFunc func(*Context)
type ErrorHandlerFunc func(*Context, int)

type Route struct {
	pattern  string
	handlers []HandlerFunc
	methods  map[string]bool
}

type App struct {
	routes []*Route
	// if an handler gets an error and  the app.OnErrorCode is called and the error code and the handler are passed as parameters to the app.OnErrorCode
	onErrorCode ErrorHandlerFunc
}

func New() *App {
	return &App{}
}

func (app *App) handle(pattern string, handlers []HandlerFunc, methods ...string) {
	methodsMap := make(map[string]bool)
	for _, method := range methods {
		methodsMap[method] = true
	}
	route := &Route{pattern: pattern, handlers: handlers, methods: methodsMap}
	app.routes = append(app.routes, route)
}

func (app *App) Use(handlers ...HandlerFunc) {
	for _, route := range app.routes {
		route.handlers = append(handlers, route.handlers...)
	}
}

// func (app *App) OnErrorCode(errorCode int, handler ErrorHandlerFunc) {
// 	app.onErrorCode = handler
// }
// ...

func (app *App) Static(path string, dir string) {
	fileServer := http.FileServer(http.Dir(dir))
	app.GET(path+"*", func(c *Context) {
		http.StripPrefix(path, fileServer).ServeHTTP(c.ResponseWriter, c.Request)
	})
}

func (app *App) GET(path string, handler ...HandlerFunc) {
	app.handle(path, handler, "GET")
}

func (app *App) PUT(path string, handler ...HandlerFunc) {
	app.handle(path, handler, "PUT")
}

func (app *App) POST(path string, handler ...HandlerFunc) {
	app.handle(path, handler, "POST")
}

func (app *App) DELETE(path string, handler ...HandlerFunc) {
	app.handle(path, handler, "DELETE")
}

func (app *App) PATCH(path string, handler ...HandlerFunc) {
	app.handle(path, handler, "PATCH")
}

func (app *App) OPTIONS(path string, handler ...HandlerFunc) {
	app.handle(path, handler, "OPTIONS")
}

func (app *App) HEAD(path string, handler ...HandlerFunc) {
	app.handle(path, handler, "HEAD")
}

func (app *App) Any(path string, handler ...HandlerFunc) {
	app.handle(path, handler, "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD")
}

func (app *App) OnErrorCode()

func (app *App) NotAllowed(c *Context) {
	http.Error(c.ResponseWriter, "405 Method not allowed", http.StatusMethodNotAllowed)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{ResponseWriter: w, Request: r}
	for _, route := range app.routes {
		if strings.HasSuffix(route.pattern, "*") {
			if strings.HasPrefix(r.URL.Path, strings.TrimSuffix(route.pattern, "*")) {
				if route.methods[r.Method] {
					c.handlers = route.handlers
					c.Next()
					return
				} else {
					app.NotAllowed(c)
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
					app.NotAllowed(c)
					return
				}
			}
		}
	}
	if app.onErrorCode != nil {
		app.onErrorCode(c, http.StatusNotFound)
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

func (app *App) Run(addr string) error {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(addr, app); err != nil {
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
