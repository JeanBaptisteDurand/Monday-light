package main

import (
    "log"
    "os"

    "monday-light/db"
    "monday-light/handlers"

    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize the database
    db.InitDB()
    defer db.DB.Close()

    r := gin.Default()
    r.LoadHTMLGlob("templates/*")

    r.Static("/static", "./frontend/static")

    // Public routes
    r.GET("/login", handlers.ShowLogin)
    r.POST("/login", handlers.Login)
    r.GET("/register", handlers.ShowRegister)
    r.POST("/register", handlers.Register)

    // Protected routes
    authorized := r.Group("/")
    authorized.Use(handlers.AuthMiddleware())
    {
        authorized.GET("/", handlers.ShowDashboard)
        authorized.GET("/project/:id", handlers.ShowProject)
        authorized.POST("/project", handlers.CreateProject)
        authorized.POST("/project/:id/category", handlers.AddCategory)
        authorized.POST("/task", handlers.CreateTask)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting server on :%s", port)
    r.Run(":" + port)
}
