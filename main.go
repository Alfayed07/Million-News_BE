package main

import (
	"BACKEND_SEJUTA_BERITA/routes"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Ensure uploads directory exists
	uploads := filepath.FromSlash("public/uploads")
	if err := os.MkdirAll(uploads, 0755); err != nil {
		log.Fatal(err)
	}

	router := routes.SetupRouter()
	router.Run(":8070")
}