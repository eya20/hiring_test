package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/eya20/hiring_test/app/catalog"
	"github.com/eya20/hiring_test/app/database"
	"github.com/eya20/hiring_test/models"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize repositories
	prodRepo := models.NewProductsRepository(db)
	categoriesRepo := models.NewCategoriesRepository(db)

	// Initialize services
	catalogService := catalog.NewCatalogService(prodRepo)

	// Initialize handlers
	catalogHandler := catalog.NewCatalogHandler(catalogService)
	categoriesHandler := catalog.NewCategoriesHandler(categoriesRepo)

	// Set up routing
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", catalogHandler.GetCatalog)
	mux.HandleFunc("GET /catalog/{code}", catalogHandler.GetProductDetails)
	mux.HandleFunc("GET /categories", categoriesHandler.GetCategories)
	mux.HandleFunc("POST /categories", categoriesHandler.CreateCategory)

	// Set up the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: mux,
	}

	// Start the server
	go func() {
		log.Printf("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}

		log.Println("Server stopped gracefully")
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	srv.Shutdown(ctx)
	stop()
}
