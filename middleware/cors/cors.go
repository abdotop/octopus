package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/abdotop/octopus"
)

// Config defines the config for CORS middleware.
type Config struct {
	// allow origin
	AllowedOrigins []string
	// allow methods
	AllowedMethods []string
	// allow headers
	AllowedHeaders []string
	// allow credentials
	AllowCredentials bool
	// expose headers
	ExposedHeaders []string
	// max age
	MaxAge int
}

// New returns a new CORS middleware.
func New(config Config) octopus.HandlerFunc {
	// Defaults for config
	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"*"}
	}
	if len(config.AllowedMethods) == 0 {
		config.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if len(config.AllowedHeaders) == 0 {
		config.AllowedHeaders = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}
	}
	if len(config.ExposedHeaders) == 0 {
		config.ExposedHeaders = []string{}
	}
	if config.MaxAge == 0 {
		config.MaxAge = 86400 // 24 hours
	}

	// return middleware octopus handler func with config
	return func(c *octopus.Ctx) {
		w, ok := c.Values.Get("response")
		if ok {
			w := w.(http.ResponseWriter)
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(config.AllowedOrigins, ","))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ","))
			w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(config.AllowCredentials))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ","))
			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
			r, ok := c.Values.Get("request")
			if ok {
				r := r.(*http.Request)
				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}
			}
		}

		c.Next()
	}
}
