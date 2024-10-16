package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/m0hit-kumar/tambola/migrations"
	"github.com/m0hit-kumar/tambola/storage"
)

// test
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/test",r.Test)
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
	api.Post("/login", r.Login)
	api.Post("/signup", r.SignUp)
	api.Get("/ticketDesign/:id?", r.GeTTicketTemplate)
	api.Post("/create_ticketDesign", r.CreateTicketTemplate)
}
func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("No .env file found, using environment variables")
		}
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
	 
    app.Use(cors.New())
	r.SetupRoutes(app)
	app.Listen(":8080")

}
