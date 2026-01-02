package api

import (
	"net/http"
	"log/slog"
	"strings"
	"fmt"
)

type Controller struct {
	Log *slog.Logger
	ValidTokens map[string]struct{}
}

func NewController(logger *slog.Logger, validTokens map[string]struct{}) *Controller {
	return &Controller{Log: logger, ValidTokens: validTokens}
}

func (c *Controller) Sign(w http.ResponseWriter, r *http.Request) {
	// TODO parse Authorisation header
	c.Log.Info("Sign endpoint called")
	authToken := r.Header.Get("Authorization")
	fmt.Println("Authorization token received:", authToken)
	isValid := c.validateAuthToken(authToken)
	if !isValid{
		fmt.Println("Unauthorized access attempt to Sign endpoint")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println("Sign endpoint not yet implemented")
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) validateAuthToken(token string) bool {
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	} else {
		return false
	}
	fmt.Println("Validating token:", token)
	fmt.Println("Valid tokens are:", c.ValidTokens)
	_, exists := c.ValidTokens[token]
	return exists
}
