package middlewares

import (
	"fmt"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi-example/utils"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/models"
)

var (
	signingSecret = utils.GetSingingSecret()
	JWTAuth       = &okapi.JWTAuth{
		SigningSecret:    []byte(signingSecret),
		Audience:         "okapi.jkaninda.dev",
		Issuer:           "okapi.jkaninda.dev",
		TokenLookup:      "header:Authorization",
		ClaimsExpression: "Equals(`email_verified`, `true`) && OneOf(`user.role`, `admin`, `owner`,`user`) && Contains(`permissions`, `read`, `create`)",
		ForwardClaims: map[string]string{
			"email": "user.email",
			"role":  "user.role",
			"name":  "user.name",
		},
	}
	AdminJWTAuth = &okapi.JWTAuth{
		SigningSecret:    []byte(signingSecret),
		TokenLookup:      "header:Authorization",
		Audience:         "okapi.jkaninda.dev",
		Issuer:           "okapi.jkaninda.dev",
		ClaimsExpression: "Equals(`email_verified`, `true`) && Equals(`user.role`, `admin`) && Contains(`permissions`, `read`, `create`, `delete`, `update`)",
		ForwardClaims: map[string]string{
			"email": "user.email",
			"role":  "user.role",
			"name":  "user.name",
		},
		// CustomClaims claims validation function
		ValidateClaims: func(context okapi.Context, claims jwt.Claims) error {
			slog.Info("Validating JWT claims for role using custom function")
			// Simulate a custom claims validation
			if _, ok := claims.(jwt.Claims); ok {

			}
			return nil
		},
	}
	jwtClaims = jwt.MapClaims{
		"sub": "12345",
		"iss": "okapi.jkaninda.dev",
		"aud": "okapi.jkaninda.dev",
		"user": map[string]string{
			"name":  "",
			"role":  "",
			"email": "",
		},
		"email_verified": true,
		"permissions":    []string{"read", "create"},
		"exp":            time.Now().Add(2 * time.Hour).Unix(),
	}
	adminPermissions = []string{"read", "create", "delete", "update"}
)

func Login(authRequest *models.AuthRequest) (models.AuthResponse, error) {
	// This is where you would typically validate the user credentials against a database

	logger.Info("Login attempt", "username", authRequest.Username)
	// Simulate a login function that returns a JWT token
	if authRequest.Username != "admin" && authRequest.Password != "password" ||
		authRequest.Username != "user" && authRequest.Password != "password" {
		return models.AuthResponse{
			Success: false,
			Message: "Invalid username or password",
		}, fmt.Errorf("username or password is wrong")
	}

	if _, ok := jwtClaims["user"].(map[string]string); ok {
		jwtClaims["user"].(map[string]string)["name"] = strings.ToUpper(authRequest.Username)
		jwtClaims["user"].(map[string]string)["role"] = authRequest.Username
		jwtClaims["user"].(map[string]string)["email"] = authRequest.Username + "@example.com"
		if authRequest.Username == "admin" {
			jwtClaims["permissions"] = adminPermissions
		} else {
			jwtClaims["permissions"] = []string{"read", "create"}

		}

	}
	// Set the expiration time for the JWT token
	expireAt := 30 * time.Minute
	jwtClaims["exp"] = time.Now().Add(expireAt).Unix()

	token, err := okapi.GenerateJwtToken(JWTAuth.SigningSecret, jwtClaims, expireAt)
	if err != nil {

		return models.AuthResponse{
			Success: false,
			Message: "Invalid username or password",
		}, fmt.Errorf("failed to generate JWT token: %w", err)
	}
	return models.AuthResponse{
		Success:   true,
		Message:   "Welcome back " + authRequest.Username,
		Token:     token,
		ExpiresAt: time.Now().Add(expireAt).Unix(),
	}, nil

}
func CustomMiddleware(next okapi.HandleFunc) okapi.HandleFunc {
	return func(c okapi.Context) error {
		slog.Info("Custom middleware executed", "path", c.Request().URL.Path, "method", c.Request().Method)
		// You can add any custom logic here, such as logging, authentication, etc.
		// For example, let's log the request method and URL
		logger.Info("Request received", "method", c.Request().Method, "url", c.Request().URL.String())
		// Call the next handler in the chain
		if err := next(c); err != nil {
			// If an error occurs, log it and return a generic error response
			slog.Error("Error in custom middleware", "error", err)
			return c.JSON(http.StatusInternalServerError, okapi.M{"error": "Internal Server Error"})
		}
		return nil
	}
}
