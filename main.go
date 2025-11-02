package main

import (
	"BACKEND_SEJUTA_BERITA/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	router := routes.SetupRouter()
	router.Run(":8070")
}