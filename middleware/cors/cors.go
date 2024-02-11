package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/abdotop/octopus"
)

type Config struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	ExposedHeaders   []string
	MaxAge           int
}

func New(config Config) octopus.HandlerFunc {
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

	return func(c *octopus.Ctx) {
		c.Response.Header().Set("Access-Control-Allow-Origin", strings.Join(config.AllowedOrigins, ","))
		c.Response.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ","))
		c.Response.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ","))
		c.Response.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(config.AllowCredentials))
		c.Response.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ","))
		c.Response.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))

		if c.Request.Method == "OPTIONS" {
			c.Response.WriteHeader(http.StatusOK)
			return
		}

		c.Next()
	}
}
