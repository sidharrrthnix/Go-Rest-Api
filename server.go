package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"http-api.com/middlewares"
	"http-api.com/respository/sqlconnect"
	"http-api.com/router"
	"http-api.com/utils"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using default/system environment variables")
	}

	// Load configuration using utilities
	serverConfig := utils.LoadServerConfig()

	// Connect to MariaDB on startup
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}
	defer db.Close()

	// Test query to verify connection works
	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		log.Printf("⚠️  Warning: Could not query database version: %v\n", err)
	} else {
		fmt.Printf("✓ MariaDB version: %s\n", version)
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := middlewares.NewRateLimiter(10, time.Minute)

	handler := utils.ApplyMiddlewares(router.NewRouter(),
		middlewares.Cors,
		rl.Middleware,
		middlewares.ResponseTimeMiddleware,
		middlewares.SecurityHeaders,
		middlewares.CompressionMiddleware,
	)

	server := &http.Server{
		Addr:      ":" + serverConfig.Port,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	fmt.Printf("🚀 Server is running on port: %s\n", serverConfig.Port)
	if err := server.ListenAndServeTLS(serverConfig.CertFile, serverConfig.KeyFile); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
