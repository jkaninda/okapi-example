package services

import (
	"fmt"

	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/middlewares"
	"github.com/jkaninda/okapi-example/models"
)

type AuthService struct{}

// ******************** AuthService *****************

func (bc *AuthService) Login(c *okapi.Context) error {
	authRequest := &models.AuthRequest{}
	err := c.Bind(authRequest)
	if err != nil {
		return c.ErrorBadRequest(models.ErrorResponse("Bad Request", err))
	}
	// Validate the authRequest and generate a JWT token
	authResponse, err := middlewares.Login(authRequest)
	if err != nil {
		return c.ErrorUnauthorized(models.ErrorResponse("Invalid username or password", err))
	}
	return c.OK(models.SuccessResponse("Welcome back", authResponse))
}
func (bc *AuthService) WhoAmI(c *okapi.Context) error {
	// Get User Information from the context, shared by the JWT middleware using forwardClaims
	email := c.GetString("email")
	if email == "" {
		return c.AbortUnauthorized("Unauthorized", fmt.Errorf("user not authenticated"))
	}

	c.Response().Header().Set("X-Okapi-User", email)
	c.Response().Header().Set("X-Okapi-User-Name", c.GetString("name"))
	c.Response().Header().Set("X-Okapi-Role", c.GetString("role"))
	// Respond with the current user information
	return c.OK(models.SuccessResponse("Ok", models.UserInfo{
		Email: email,
		Role:  c.GetString("role"),
		Name:  c.GetString("name")}))
}
