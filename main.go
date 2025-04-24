package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	//Init DB in db.go
	db := InitializeDB()

	//Create fiber app instance
	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})

	//Define Routes in auth.go
	AuthHandlers(app.Group("/auth"), db)

	//Verify the JWT
	protected := app.Use(AuthMiddleware(db))

	//Define Book routes(Will require a valid JWT)
	BookHandlers(protected.Group("/book"), db)

	// Download the books as either CSV or JSON
	Download(protected.Group("/download"), db)

	//Start server on port 3000
	app.Listen(":3000")
}
