package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"monday-light/db"
	"monday-light/models"
	"github.com/gin-gonic/gin"
)

// Helper function to render templates
func Render(c *gin.Context, contentTemplate string, data gin.H) {
	if c.GetHeader("HX-Request") != "" {
		// Render only the content template for HTMX requests
		tmpl, err := template.ParseFiles("templates/" + contentTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(c.Writer, "content", data); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	} else {
		// Render base.html with the content template
		tmpl, err := template.ParseFiles("templates/base.html", "templates/" + contentTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(c.Writer, "base", data); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	}
}

func ShowDashboard(c *gin.Context) {
	// Load projects from the database
	rows, err := db.DB.Query("SELECT id, name, categories FROM projects")
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Categories); err != nil {
			c.String(http.StatusInternalServerError, "Failed to scan projects")
			return
		}
		projects = append(projects, project)
	}

	data := gin.H{
		"Title":    "Dashboard",
		"Projects": projects,
	}

	Render(c, "dashboard.html", data)
}

func ShowProject(c *gin.Context) {
	// Get project ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Fetch project from the database
	var project models.Project
	row := db.DB.QueryRow("SELECT id, name, categories FROM projects WHERE id = $1", id)
	if err := row.Scan(&project.ID, &project.Name, &project.Categories); err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	// Fetch tasks for the project (assuming you have tasks)
	// ...

	data := gin.H{
		"Title":   project.Name,
		"Project": project,
		// "Tasks":   tasks, // Include tasks if available
	}

	Render(c, "project.html", data)
}

func CreateProject(c *gin.Context) {
	var newProject models.Project
	if err := c.Bind(&newProject); err != nil {
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
		Name string `form:"name" json:"name" binding:"required"`
	}
	if err := c.Bind(&category); err != nil {
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

func CreateTask(c *gin.Context) {
	// Implement the logic to create a task
	// ...
}
