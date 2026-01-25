package main

import (
	"log/slog"
	"os"

	"github.com/Peto10/SSH-like-Certificate-Authority-Service/internal/api"
	"github.com/Peto10/SSH-like-Certificate-Authority-Service/internal/server"
	"golang.org/x/crypto/ssh"
)

const (
	defaultServerHostName = ":8443"
	defaultURL            = "https://localhost" + defaultServerHostName
	appSecretsDir		  = "/run/ca-service"
	privateHttpsKeyFile   = appSecretsDir + "/https/ca-service-local.key.pem"
	publicHttpsCertFile   = appSecretsDir + "/https/ca-service-local.cert.pem"
	privateKeyFile        = appSecretsDir + "/ssh/ca_key"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	allowedTokens, err := parseTokenPrincipals(os.Getenv("CA_ACCESS_TOKEN"))
	if err != nil {
		logger.Error("failed to parse allowed tokens from environment variable")
		os.Exit(1)
	}
	if len(allowedTokens) == 0 {
		logger.Warn("No static keys parsed from environment variable")
	}

	caKeyBytes, err := os.ReadFile(privateKeyFile)
	if err != nil {
		logger.Error("failed to read CA private key", "error", err)
		os.Exit(1)
	}

	caSigner, err := ssh.ParsePrivateKey(caKeyBytes)
	if err != nil {
		logger.Error("failed to parse CA private key", "error", err)
		os.Exit(1)
	}

	c := api.NewController(logger.With("component", "api"), allowedTokens, caSigner)
	c.Log.Info("service starting", "base URL", defaultURL)

	s := server.NewServer(c, defaultServerHostName)

	err = s.ListenAndServeTLS(publicHttpsCertFile, privateHttpsKeyFile)
	if err != nil {
		logger.Error("Error with listening and serving port", "error", err)
		os.Exit(1)
	}
}
