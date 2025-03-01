package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/m0hit-kumar/tambola/migrations"
	"github.com/m0hit-kumar/tambola/storage"
)


func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/test",r.Test)
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
	api.Post("/login", r.Login)
	api.Post("/signup", r.SignUp)
	api.Get("/ticketDesign/:roomId?", r.GetTicketTemplate)
	api.Post("/create_ticketDesign", r.CreateTicketTemplate)
}

func main() {
	allowOrigins:="https://tambola-theta.vercel.app/"
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("No .env file found, using environment variables")
		}
		allowOrigins = "http://localhost:3000"
	}
	
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Failed to load database")
	}
	err = migrations.MigrateTables(db)
	if err != nil {
		log.Fatal("Failed to migrate books")
	}
	r := Repository{
		DB: db,
	}
	app := fiber.New()
	fmt.Print(allowOrigins)
	app.Use(cors.New(cors.Config{
		AllowOrigins:  allowOrigins,   // Your frontend URL
        AllowCredentials: true,                      // Important for cookies
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        ExposeHeaders:    "Set-Cookie",              // Important for cookies
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",            
     }))
	r.SetupRoutes(app)
	app.Listen(":8080")

}
