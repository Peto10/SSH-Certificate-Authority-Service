package api

import (
	"encoding/json"
	"net/http"
	"log/slog"
	"strings"
	"fmt"
	"golang.org/x/crypto/ssh"
)

type Controller struct {
	Log *slog.Logger
	ValidTokens map[string]struct{}
}

type requestBody struct {
	PublicKey string `json:"public_key"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type signResponse struct {
	Message    string `json:"message,omitempty"`
	SignedCert string `json:"signed_cert,omitempty"`
}

func NewController(logger *slog.Logger, validTokens map[string]struct{}) *Controller {
	return &Controller{Log: logger, ValidTokens: validTokens}
}

func (c *Controller) Sign(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authToken := r.Header.Get("Authorization")
	fmt.Println("Authorization token received:", authToken)
	isValid := c.validateAuthToken(authToken)
	if !isValid{
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "access token not valid"})
		return
	}

	var reqBody requestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "failed to decode request body"})
		return
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(reqBody.PublicKey))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "failed to parse public key"})
		return
	}

	if (pubKey.Type() != ssh.KeyAlgoED25519) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "only ed25519 keys are supported"})
		return
	}

	fmt.Println("Sign endpoint not yet implemented")	
	// TODO: implement signing logic
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(signResponse{Message: "Sign endpoint not yet implemented"})
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
