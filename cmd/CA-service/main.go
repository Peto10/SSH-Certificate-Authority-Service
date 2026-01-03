package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Peto10/SSH-like-Certificate-Authority-Service/internal/api"
	"github.com/Peto10/SSH-like-Certificate-Authority-Service/internal/server"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
)

const (
	defaultServerHostName = ":8443"
	defaultURL			  = "https://localhost" + defaultServerHostName
	certsDir              = "./certs/"
	privateHttpsKeyFile   = certsDir + "https/localhost+2-key.pem"
	publicHttpsCertFile   = certsDir + "https/localhost+2.pem"
	privateKeyFile        = certsDir + "ca-key-pair/ca_key"
	// publicKeyFile         = certsDir + "ca-key-pair/ca_key.pub"
	envFile	              = ".env"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	godotenv.Load(envFile)
	allowedTokens := parseTokenPrincipals("CA_ACCESS_TOKEN")
	if allowedTokens == nil {
		logger.Error("failed to parse allowed tokens from environment variable")
		os.Exit(1)
	}
	if len(allowedTokens) == 0 {
		logger.Warn("No static keys parsed from environment variable")
	}

	caSigner, err := parseCAPrivateKey(privateKeyFile)
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

func parseTokenPrincipals(envVarName string) map[string][]string {
	envTokens := os.Getenv(envVarName)
    m := make(map[string][]string)
	if envTokens == "" {
		return m
	}

    entries := strings.Split(envTokens, ";")
    for _, entry := range entries {
        parts := strings.Split(entry, ":")
        if len(parts) != 2 {
            return nil
        }

        token := parts[0]
        principals := strings.Split(parts[1], ",")
        m[token] = principals
    }

    return m
}

func parseCAPrivateKey(filePath string) (ssh.Signer, error) {
	caKeyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	caSigner, err := ssh.ParsePrivateKey(caKeyBytes)
	if err != nil {
		return nil, err
	}

	return caSigner, nil
}
