package handlers

import (
	"net/http"
	"monday-light/db"
	//"monday-light/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func CreateTask(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid project ID")
        return
    }

    name := c.PostForm("task_name")
    description := c.PostForm("task_description")
    category := c.PostForm("task_category") // Peut être vide
    estimatedTimeStr := c.PostForm("task_estimated_time")
    if strings.TrimSpace(name) == "" {
        c.String(http.StatusBadRequest, "Nom de tâche requis")
        return
    }

    var estimatedTime int
    if estimatedTimeStr != "" {
        estimatedTime, _ = strconv.Atoi(estimatedTimeStr)
    }

    // La tâche a par défaut un status "to_assign"
    err = db.DB.QueryRow(`
    INSERT INTO tasks (name, description, category, project_id, status, estimated_time, real_time)
    VALUES ($1, $2, $3, $4, 'to_assign', $5, 0) RETURNING id
    `, name, description, category, projectID, estimatedTime).Scan(new(int))
    if err != nil {
        c.String(http.StatusInternalServerError, "Database error (create task)")
        return
    }

    // Recharger la page du projet
    ShowProject(c)
}