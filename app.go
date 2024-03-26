package octopus

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type HandlerFunc func(*Ctx)

type app struct {
	sync.RWMutex
	w      sync.WaitGroup
	routes *routes

	globalMiddleware []HandlerFunc
	subApps          []*route
	errorHandlers    map[statusCode]HandlerFunc
}

func New() *app {
	return &app{
		subApps:          make([]*route, 0),
		routes:           new(routes),
		errorHandlers:    make(map[statusCode]HandlerFunc),
		w:                sync.WaitGroup{},
		globalMiddleware: make([]HandlerFunc, 0),
	}
}

func (a *app) handle(pattern string, handlers []HandlerFunc, methods ...string) {
	a.Lock()
	defer a.Unlock()
	for _, method := range methods {
		handlers = append(a.globalMiddleware, handlers...)
		a.routes.add(pattern, method, handlers...)
	}
}

// func (a *app) Mount(path string, subApp *app) {
// 	a.Lock()
// 	defer a.Unlock()

// 	route := &route{
// 		path: path,
// 		a:    subApp,
// 	}
// 	a.subApps = append(a.subApps, route)
// }

func (a *app) Static(path string, dir string) {
	fileServer := http.FileServer(http.Dir(dir))
	a.Get(path+"*", func(c *Ctx) {
		r, rok := c.Values.Get("request")
		w, wok := c.Values.Get("response")
		if rok && wok {
			r := r.(*http.Request)
			w := w.(http.ResponseWriter)
			http.StripPrefix(path, fileServer).ServeHTTP(w, r)
		}
	})
}

func (a *app) Use(handlers ...HandlerFunc) {
	a.Lock()
	defer a.Unlock()
	a.globalMiddleware = append(a.globalMiddleware, handlers...)
}

func (a *app) Group(path string, fn ...HandlerFunc) *route {
	r := new(route)
	r.globalMiddleware = fn
	r.path = path
	r.a = a
	return r
}

// func (a *app) mountSubApp() {
// 	a.Lock()
// 	defer a.Unlock()
// 	for _, subApp := range a.subApps {
// 		if subApp.a != nil {
// 			subApp.a.mountSubApp()
// 		}
// 		app := subApp.a
// 		p := subApp.path
// 		app.routes.rrange(func(path string, route *route) bool {
// 			for method, handlers := range route.data {
// 				a.routes.add(p+path, method, handlers...)
// 			}
// 			return true
// 		})
// 	}
// }

func (a *app) DELETE(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "DELETE")
}

func (a *app) Get(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "GET")
}

func (a *app) PUT(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "PUT")
}

func (a *app) Post(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "POST")
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

func (a *app) Method(method string, path string, handler ...HandlerFunc) {
	methods := strings.Split(method, " ")
	a.handle(path, handler, methods...)
}

func (a *app) OnErrorCode(code statusCode, f HandlerFunc) {
	a.Lock()
	defer a.Unlock()
	a.errorHandlers[code] = f
}

func (a *app) handleError(code statusCode, c *Ctx) {
	a.RLock()
	handler, exists := a.errorHandlers[code]
	a.RUnlock()
	if exists {
		handler(c)
	} else {
		func(c *Ctx) {
			message := statusMessages[code]
			if code == StatusNotFound {
				c.WriteString(string(message))
			}
		}(c)
	}
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Ctx{handlers: nil, index: 0, Values: new(value), Context: r.Context()}
	c.Values.Set("request", r)
	c.Values.Set("response", w)
	c.Values.Set("app", a)

	var routExist = false

	a.routes.rrange(func(path string, route *route) bool {
		hs, ok := route.methodExists(r.Method)
		if strings.HasSuffix(path, "*") {
			if strings.HasPrefix(r.URL.Path, strings.TrimSuffix(path, "*")) {
				if ok {
					c.handlers = hs
					c.Next()
					routExist = true
					return false
				} else {
					c.Status(StatusMethodNotAllowed)
					routExist = true
					return false
				}
			}
		} else {
			if path == r.URL.Path {
				if ok {
					c.handlers = hs
					c.Next()
					routExist = true
					return false
				} else {
					c.Status(StatusMethodNotAllowed)
					routExist = true
					return false
				}
			}
		}
		return true
	})

	if !routExist {
		c.Status(StatusNotFound)
	}
}

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
	// a.mountSubApp()

	a.w.Add(1)
	go func() {
		defer a.w.Done()
		if err := http.ListenAndServe(addr, a); err != nil {
			log.Fatal(err)
		}
	}()

	// Attendre que le serveur démarre
	time.Sleep(time.Second)

	// Vérifier si le serveur est en cours d'exécution
	checkServer(addr)

	a.w.Wait()
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
