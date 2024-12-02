package main

import (
    "monday-light/handlers"
	"monday-light/models"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.LoadHTMLGlob("templates/*")
    r.Static("/static", "./frontend/static")

    r.GET("/", handlers.ShowDashboard)
    r.GET("/project/:id", handlers.ShowProject)
    r.POST("/project", handlers.CreateProject)
    r.POST("/project/:id/category", handlers.AddCategory)
    r.POST("/task", handlers.CreateTask)

    r.Run(":8080")
}
