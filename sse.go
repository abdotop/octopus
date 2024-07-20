package octopus

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sync"

// 	"github.com/google/uuid"
// )

// type sse struct {
// 	connections *value
// 	all_key     map[uuid.UUID]bool
// 	sync.RWMutex
// }

// func newSSE() *sse {
// 	return &sse{
// 		connections: new(value),
// 		all_key:     make(map[uuid.UUID]bool),
// 	}
// }

// func (sse *sse) Add(conn *Conn) {
// 	sse.Lock()
// 	defer sse.Unlock()
// 	conn.sse = sse
// 	sse.all_key[conn.id] = true
// 	sse.connections.Set(conn.id, conn)
// }

// func (sse *sse) CloseConn(id uuid.UUID) error {
// 	sse.Lock()
// 	defer sse.Unlock()
// 	// for index, v := range sse.all_key {
// 	if _, ok := sse.all_key[id]; ok {
// 		delete(sse.all_key, id)
// 		sse.connections.Delete(id)
// 	}
// 	return nil
// }

// func (sse *sse) GetAll() []*Conn {
// 	sse.RLock() // Utiliser RLock pour un accès en lecture
// 	defer sse.RUnlock()

// 	var allConns []*Conn
// 	for id := range sse.all_key {
// 		if conn, ok := sse.connections.Get(id); ok {
// 			allConns = append(allConns, conn.(*Conn))
// 		}
// 	}
// 	return allConns
// }

// func (sse *sse) GetConn(id uuid.UUID) (*Conn, error) {
// 	sse.RLock() // Utiliser RLock pour un accès en lecture
// 	defer sse.RUnlock()

// 	if conn, ok := sse.connections.Get(id); ok {
// 		return conn.(*Conn), nil
// 	}
// 	return nil, fmt.Errorf("no connection found with ID %s", id)
// }

// func NewSseConnFromCtx(c *Ctx) (*Conn, error) {

// 	w_value, ok := c.Values.Get("response")
// 	if !ok {
// 		return nil, fmt.Errorf("failed to get Writer from context")
// 	}
// 	w, ok := w_value.(http.ResponseWriter)
// 	if !ok {
// 		return nil, fmt.Errorf("failed to get Writer from context")
// 	}
// 	flusher, ok := w.(http.Flusher)
// 	if !ok {
// 		return nil, fmt.Errorf("streaming unsupported")
// 	}

// 	conn := &Conn{
// 		id:      uuid.New(),
// 		writer:  w,
// 		flusher: flusher,
// 		CloseCh: make(chan bool),
// 	}

// 	r_value, ok := c.Values.Get("request")
// 	if !ok {
// 		return nil, fmt.Errorf("failed to get Writer from context")
// 	}
// 	r, ok := r_value.(*http.Request)
// 	if !ok {
// 		return nil, fmt.Errorf("failed to get Writer from context")
// 	}

// 	app_value, ok := c.Values.Get("app")
// 	if !ok {
// 		return nil, fmt.Errorf("failed to get App from context")
// 	}
// 	app, ok := app_value.(*App)
// 	if !ok {
// 		return nil, fmt.Errorf("failed to get App from context")
// 	}
// 	conn.context = r.Context()

// 	app.sse_service.Add(conn)

// 	// Notify the client that the connection has been established
// 	fmt.Fprintf(w, ":connected\n\n")
// 	flusher.Flush()
// 	return conn, nil
// }

// type Conn struct {
// 	id      uuid.UUID
// 	writer  http.ResponseWriter
// 	flusher http.Flusher
// 	CloseCh chan bool
// 	context context.Context
// 	closed  bool
// 	sse     *sse
// }

// func (c *Conn) ID() uuid.UUID {
// 	return c.id
// }

// func (c *Conn) Close() {
// 	select {
// 	case <-c.CloseCh:
// 		// Le canal est déjà fermé
// 		return
// 	default:
// 		c.sse.CloseConn(c.id) // Appeler CloseConn sur sse avant de fermer le canal
// 		close(c.CloseCh)
// 		fmt.Fprintf(c.writer, "event: %s\n", "close") //
// 		c.closed = true
// 	}
// }

// func (c *Conn) Done() <-chan struct{} {
// 	if !c.closed {
// 		return c.context.Done()
// 	}
// 	return nil
// }

// func (c *Conn) SendToEvent(eventName string) *Conn {
// 	if !c.closed {
// 		fmt.Fprintf(c.writer, "event: %s\n", eventName)
// 	}
// 	return c
// }

// func (c *Conn) SendString(data string) error {
// 	if c.closed {
// 		return fmt.Errorf("cannot send data, connection is closed")
// 	}
// 	fmt.Fprintf(c.writer, "data: %s\n\n", data)
// 	c.flusher.Flush()
// 	return nil
// }

// func (c *Conn) SendJSON(data interface{}) error {
// 	if c.closed {
// 		return fmt.Errorf("cannot send JSON, connection is closed")
// 	}
// 	jsonStr, _ := json.Marshal(data)
// 	fmt.Fprintf(c.writer, "data: %s\n\n", jsonStr)
// 	c.flusher.Flush()
// 	return nil
// }
