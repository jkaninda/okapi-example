package config

import (
	"fmt"
	"strings"

	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/session"
	"github.com/jkaninda/okapi-example/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	Database       DatabaseConfig
	Redis          RedisConfig
	Server         ServerConfig
	Cors           CorsConfig
	JWT            JWTConfig
	Log            LogConfig
	SessionManager *session.SessionManager
}

type DatabaseConfig struct {
	dbHost     string
	dbUser     string
	dbPassword string
	dbName     string
	dbPort     int
	dbSslMode  string
}

type RedisConfig struct {
	URL string
}

type ServerConfig struct {
	port        int
	tlsPort     int
	environment string
	enableDocs  bool
	tls         Tls
}
type Tls struct {
	Cert        string
	Key         string
	CA          string
	RequireAuth bool
}
type CorsConfig struct {
	AllowedOrigins []string
}

type JWTConfig struct {
	Secret string
}

type LogConfig struct {
	Level string
}

func New() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()
	return &Config{
		Database: DatabaseConfig{
			dbHost:     goutils.Env("DB_HOST", "localhost"),
			dbUser:     goutils.Env("DB_USER", ""),
			dbPassword: goutils.Env("DB_PASSWORD", ""),
			dbName:     goutils.Env("DB_NAME", "account"),
			dbPort:     goutils.EnvInt("DB_PORT", 5432),
			dbSslMode:  goutils.Env("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			URL: goutils.Env("REDIS_URL", "redis://localhost:6379/0"),
		},
		Server: ServerConfig{
			enableDocs:  goutils.EnvBool("ENABLE_DOCS", true),
			port:        goutils.EnvInt("PORT", 8080),
			tlsPort:     goutils.EnvInt("TLS_PORT", 8443),
			environment: goutils.Env("ENVIRONMENT", "development"),
			tls: Tls{
				Cert: goutils.Env("TLS_CERT_PATH", ""),
				Key:  goutils.Env("TLS_KEY_PATH", ""),
				// mTLS
				CA:          goutils.Env("TLS_CA_PATH", ""),
				RequireAuth: goutils.EnvBool("TLS_REQUIRE_AUTH", false),
			},
		},

		Cors: CorsConfig{
			AllowedOrigins: strings.Split(goutils.Env("CORS_ALLOWED_ORIGINS", "http://localhost:8080"), ","),
		},
		JWT: JWTConfig{
			Secret: goutils.Env("JWT_SECRET", "default-secret-key"),
		},
		Log: LogConfig{
			Level: goutils.Env("LOG_LEVEL", "info"),
		},
		SessionManager: session.New(),
	}
}
func (c *Config) validate() error {
	if c.Redis.URL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if c.Server.port == 0 {
		return fmt.Errorf("PORT is required")
	}
	return nil
}
func (c *Config) Initialize(app *okapi.Okapi) error {
	if err := c.validate(); err != nil {
		return err
	}

	// Init Doc
	if c.Server.enableDocs {
		app.WithOpenAPIDocs(okapi.OpenAPI{
			Title:   utils.AppName,
			Version: utils.AppVersion,
			License: okapi.License{
				Name: "MIT",
			},
			Contact: okapi.Contact{
				Name:  "Jonas Kaninda",
				Email: "me@jkaninda.dev",
				URL:   "https://github.com/jkaninda/okapi"},
			SecuritySchemes: okapi.SecuritySchemes{
				{
					Name:         "bearerAuth",
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
			},
		})
	}
	if len(c.Cors.AllowedOrigins) > 0 {
		app.WithCORS(okapi.Cors{AllowedOrigins: c.Cors.AllowedOrigins})
	}
	// Init TLS server
	if len(c.Server.tls.Cert) > 0 && len(c.Server.tls.Key) > 0 {
		tls, err := goutils.LoadTLSConfig(c.Server.tls.Cert, c.Server.tls.Key, c.Server.tls.CA, c.Server.tls.RequireAuth)
		if err != nil {
			panic(err)
		} else {
			tlsAdd := fmt.Sprintf(":%d", c.Server.tlsPort)
			// Add tls server
			app.With(okapi.WithTLSServer(tlsAdd, tls))
			// Or use a single server
			// app.With(okapi.WithTLS(tls))
		}
	}

	addr := fmt.Sprintf(":%d", c.Server.port)
	app.With(okapi.WithAddr(addr))
	return nil
}
