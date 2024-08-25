package main

import (
	"log"

	"fmt"

	"github.com/adarsh-jaiss/shipper/server"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*","http://localhost:3000","https://shipper-ui-gamma.vercel.app/build"},
		AllowHeaders: []string{"Origin, Content-Type, Accept, Methods"},
	}))

	app.Post("/build", server.BuildHandler)

	log.Fatal(app.Listen(":8080"))
}


