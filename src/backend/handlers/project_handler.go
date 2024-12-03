package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"monday-light/db"
	"monday-light/models"
	"github.com/gin-gonic/gin"
)

func ShowDashboard(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, categories FROM projects")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Categories); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan projects"})
			return
		}
		projects = append(projects, project)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Projects": projects,
	})
}

func ShowProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	row := db.DB.QueryRow("SELECT id, name, categories FROM projects WHERE id = $1", id)
	err = row.Scan(&project.ID, &project.Name, &project.Categories)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.HTML(http.StatusOK, "project.html", gin.H{
		"Project": project,
	})
}

func CreateProject(c *gin.Context) {
	var newProject models.Project
	if err := c.BindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := db.DB.QueryRow(
		"INSERT INTO projects (name, categories) VALUES ($1, $2) RETURNING id",
		newProject.Name, newProject.Categories,
	).Scan(&newProject.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, newProject)
}

func AddCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var category struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = db.DB.Exec(
		"UPDATE projects SET categories = array_append(categories, $1) WHERE id = $2",
		category.Name, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category added"})
}
