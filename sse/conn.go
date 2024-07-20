package sse

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/abdotop/octopus"
	"github.com/google/uuid"
)

type (
	// SSEConfig définit la configuration pour le serveur SSE.
	ClientConfig struct {
		ID string // Identifiant unique du client

		// HeaderFields contient des valeurs d'en-tête personnalisées pour les réponses SSE.
		HeaderFields map[string]string

		// RetryInterval spécifie le temps en millisecondes avant une nouvelle tentative de connexion.
		RetryInterval time.Duration

		// ConnectionTimeout spécifie la durée après laquelle une connexion inactive doit être fermée.
		ConnectionTimeout time.Duration

		// EnableCORS indique si les en-têtes CORS doivent être ajoutés aux réponses.
		EnableCORS bool

		// CORSOrigin spécifie la valeur de l'en-tête Access-Control-Allow-Origin.
		CORSOrigin string

		// UseHTTPS indique si HTTPS doit être utilisé pour sécuriser les connexions.
		UseHTTPS bool

		// CompressionEnabled indique si la compression des données est activée.
		CompressionEnabled bool

		// LoggingEnabled indique si le logging est activé pour les événements SSE.
		LoggingEnabled bool
	}

	Conn struct {
		id      string
		writer  http.ResponseWriter
		flusher http.Flusher
		closeCh chan bool
		context context.Context
		closed  bool
		config  *ClientConfig
		mu      sync.Mutex // Mutex pour la gestion sûre de la fermeture
	}

	Event struct {
		eventType string
		conn      *Conn
	}
)

// // NewClientConfig crée une configuration spécifique pour un client avec des valeurs par défaut.
func (c *ClientConfig) Default() *ClientConfig {
	return &ClientConfig{
		ID: uuid.NewString(),
		HeaderFields: map[string]string{
			"Content-Type":  "text/event-stream",
			"Cache-Control": "no-cache",
			"Connection":    "keep-alive",
		},
		RetryInterval:      3000 * time.Millisecond,
		ConnectionTimeout:  5 * time.Minute,
		EnableCORS:         true,
		CORSOrigin:         "*",
		UseHTTPS:           false,
		CompressionEnabled: false,
		LoggingEnabled:     true,
	}
}

func ConnFrom(c *octopus.Ctx, conf *ClientConfig) (*Conn, error) {
	// Handle configuration override if one is provided
	if conf == nil {
		return nil, fmt.Errorf("no configuration provided")
	}

	// Retrieve the http.ResponseWriter from the context
	w_value, ok := c.Values.Get("response")
	if !ok {
		return nil, fmt.Errorf("failed to get Writer from context")
	}
	w, ok := w_value.(http.ResponseWriter)
	if !ok {
		return nil, fmt.Errorf("failed to assert Writer from context")
	}

	// Check if the ResponseWriter supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming unsupported")
	}

	r_value, ok := c.Values.Get("request")
	if !ok {
		return nil, fmt.Errorf("failed to get Writer from context")
	}
	r, ok := r_value.(*http.Request)
	if !ok {
		return nil, fmt.Errorf("failed to get Writer from context")
	}

	// Create and return the Conn instance
	conn := &Conn{
		id:      conf.ID,
		writer:  w,
		flusher: flusher,
		closeCh: make(chan bool),
		context: r.Context(),
		closed:  false,
		config:  conf,
		mu:      sync.Mutex{},
	}

	return conn, nil
}

// compress prend une chaîne de caractères et la compresse en utilisant gzip.
func compress(data string) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(data)); err != nil {
		gz.Close()
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Close marque la connexion comme fermée et ferme le canal closeCh.
func (c *Conn) Close() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		err = c.send("close", "Server closing connection")
		c.closed = true
		close(c.closeCh)
	}
	return err
}

// Done retourne un canal qui est fermé lorsque la connexion est fermée.
// Ce canal est fermé en réponse à la fermeture de closeCh ou du contexte.
func (c *Conn) Done() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		case <-c.closeCh:
			// La fermeture a été initiée par Close()
		case <-c.context.Done():
			c.mu.Lock()
			if !c.closed {
				c.closed = true
				close(c.closeCh) // Assurez-vous de fermer closeCh également ici
			}
			c.mu.Unlock()
			// La fermeture a été initiée par le contexte parent
		}
	}()
	return done
}

func (c *Conn) ID() string {
	return c.id
}

// SendEvent envoie un événement SSE avec un type et des données spécifiés, en respectant la configuration.
func (c *Conn) send(eventType, data string) error {
	if c.closed {
		return fmt.Errorf("connection is closed")
	}

	// Appliquer les en-têtes personnalisés
	for key, value := range c.config.HeaderFields {
		c.writer.Header().Set(key, value)
	}

	// Ajouter l'ID de la connexion dans les en-têtes HTTP
	// c.Writer.Header().Set("X-Connection-ID", c.ID)

	// Gérer CORS si activé
	if c.config.EnableCORS {
		c.writer.Header().Set("Access-Control-Allow-Origin", c.config.CORSOrigin)
	}

	// Compresser les données si la compression est activée
	if c.config.CompressionEnabled {
		compressedData, err := compress(data)
		if err != nil {
			return fmt.Errorf("compression error: %v", err)
		}
		data = string(compressedData)
		c.writer.Header().Set("Content-Encoding", "gzip")
	}

	// Préparer le message SSE
	message := fmt.Sprintf("event: %s\n", eventType)
	if c.config.RetryInterval > 0 {
		message += fmt.Sprintf("retry: %d\n", c.config.RetryInterval.Milliseconds())
	}
	message += fmt.Sprintf("data: %s\n\n", data)

	// Envoyer le message
	_, err := c.writer.Write([]byte(message))
	if err != nil {
		return err
	}

	// Flush si possible
	if flusher, ok := c.writer.(http.Flusher); ok {
		flusher.Flush()
	} else {
		return fmt.Errorf("failed to flush data")
	}

	// Logging si activé
	if c.config.LoggingEnabled {
		log.Printf("Sent SSE event: %s, Data: %s", eventType, data)
	}

	return nil
}

// SendJSON envoie des données JSON sur une connexion SSE.
func (c *Conn) SendJSON(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.send("", string(jsonData))
}

// SendText envoie des données textuelles sur une connexion SSE.
func (c *Conn) SendText(data string) error {
	return c.send("", data)
}

// SendEvent envoie un événement SSE avec un type et des données spécifiés, en respectant la configuration.
func (c *Conn) Event(eventType string) *Event {
	return &Event{
		eventType: eventType,
		conn:      c,
	}
}

// SendJSON envoie des données JSON sur une connexion SSE.
func (e *Event) SendJSON(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return e.conn.send(e.eventType, string(jsonData))
}

// SendText envoie des données textuelles sur une connexion SSE.
func (e *Event) SendText(data string) error {
	return e.conn.send(e.eventType, data)
}
