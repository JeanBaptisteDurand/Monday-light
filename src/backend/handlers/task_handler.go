package handlers

import (
    "net/http"
    "monday-light/models"
    "github.com/gin-gonic/gin"
)

var tasks []models.Task

func CreateTask(c *gin.Context) {
    var newTask models.Task
    if err := c.BindJSON(&newTask); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    newTask.ID = len(tasks) + 1
    newTask.Status = "to_assign"
    tasks = append(tasks, newTask)
    c.JSON(http.StatusCreated, newTask)
}
