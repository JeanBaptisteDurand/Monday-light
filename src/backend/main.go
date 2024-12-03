package main

import (
	"monday-light/db"
	"monday-light/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db.InitDB()
	defer db.DB.Close()

	// Set up the Gin router
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./frontend/static")

	// Define routes
	r.GET("/", handlers.ShowDashboard)
	r.GET("/project/:id", handlers.ShowProject)
	r.POST("/project", handlers.CreateProject)
	r.POST("/project/:id/category", handlers.AddCategory)
	r.POST("/task", handlers.CreateTask)

	// Start the server
	r.Run(":8080")
}
