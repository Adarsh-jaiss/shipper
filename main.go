package main

import (
	"flag"
	"log"

	"fmt"

	"github.com/adarsh-jaiss/shipper/server"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	fmt.Println("Running API server...")
	app := fiber.New()
	app.Post("/build", server.BuildHandler)

	log.Fatal(app.Listen(":8080"))
}


