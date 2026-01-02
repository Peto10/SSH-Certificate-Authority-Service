package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Peto10/SSH-like-Certificate-Authority-Service/internal/api"
	"github.com/Peto10/SSH-like-Certificate-Authority-Service/internal/server"
	"github.com/joho/godotenv"
)

const (
	defaultServerHostName = ":8443"
	defaultURL			  = "https://localhost" + defaultServerHostName
	certsDir              = "./certs/"
	privateKeyFile        = certsDir + "localhost+2-key.pem"
	publicCertFile        = certsDir + "localhost+2.pem"
	envFilePath           = ".env"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	godotenv.Load(envFilePath)
	keys := parseStaticKeys("CA_ACCESS_TOKEN")
	if len(keys) == 0 {
		logger.Warn("No static keys parsed from environment variable")
	}

	c := api.NewController(logger, keys)

	c.Log.Info("service starting", "URL", defaultURL)

	s := server.NewServer(c, defaultServerHostName)

	err := s.ListenAndServeTLS(publicCertFile, privateKeyFile)
	if err != nil {
		c.Log.Error("Error with listening and serving port", "error", err)
		os.Exit(1)
	}
}

func parseStaticKeys(envVarName string) map[string]struct{} {
	keysStr := os.Getenv(envVarName)
	keys := make(map[string]struct{})
	if keysStr == "" {
		return keys
	}
	for _, key := range strings.Split(keysStr, ",") {
		keys[key] = struct{}{}
	}
	return keys
}
