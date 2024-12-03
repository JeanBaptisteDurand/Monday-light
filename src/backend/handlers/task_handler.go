package handlers

import (
	"net/http"
	"monday-light/db"
	"monday-light/models"
	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var newTask models.Task
	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := db.DB.QueryRow(
		"INSERT INTO tasks (name, description, category, project_id, status, assigned_users, estimated_time, real_time) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		newTask.Name, newTask.Description, newTask.Category, newTask.ProjectID,
		newTask.Status, newTask.AssignedUsers, newTask.EstimatedTime, newTask.RealTime,
	).Scan(&newTask.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, newTask)
}
