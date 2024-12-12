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

    gin.SetMode(gin.DebugMode)
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
        authorized.GET("/show-new-project-form", handlers.ShowNewProjectForm)

        authorized.POST("/project/:id/category", handlers.AddCategory)
        authorized.POST("/project/:id/category/remove", handlers.RemoveCategory)

        // Task-related routes
        authorized.POST("/project/:id/task", handlers.CreateTask)
        authorized.GET("/project/:id/task/:task_id", handlers.GetTaskDetail)
        authorized.POST("/project/:id/task/:task_id/next_status", handlers.NextTaskStatus)
        authorized.POST("/project/:id/task/:task_id/assign", handlers.AssignToSelf)

        authorized.GET("/recap", handlers.ShowRecap)
        authorized.GET("/param", handlers.ShowParam)
        authorized.GET("/logout", handlers.Logout)
        authorized.GET("/param/edit", handlers.ShowParamEdit)
        authorized.POST("/param/update", handlers.UpdateParam)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting server on :%s", port)
    r.Run(":" + port)
}
