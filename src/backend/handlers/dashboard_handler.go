package handlers

import (
	"html/template"
	"net/http"
	"log"

	"monday-light/db"
	"monday-light/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func RenderP(c *gin.Context, contentTemplate string, data gin.H) {
    if c.GetHeader("HX-Request") != "" {
        tmpl, err := template.ParseFiles(
            "templates/"+contentTemplate,
            "templates/project_categories.html",
        )
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
        tmpl, err := template.ParseFiles(
            "templates/base.html",
            "templates/"+contentTemplate,
            "templates/sidebar_projects.html",
            "templates/project_categories.html",
        )
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
		tmpl, err := template.ParseFiles("templates/base.html", "templates/" + contentTemplate, "templates/sidebar_projects.html")
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
    userIDVal, _ := c.Get("userID")
    userID := userIDVal.(int)

    // Get user's username
    var username string
    err := db.DB.QueryRow("SELECT username FROM users WHERE id = $1", userID).Scan(&username)
    if err != nil {
        c.String(http.StatusInternalServerError, "Database error: failed to fetch username")
        log.Printf("Error fetching username for userID=%d: %v", userID, err)
        return
    }

    // Fetch projects from the database
    rows, err := db.DB.Query("SELECT id, name, categories FROM projects")
    if err != nil {
        c.String(http.StatusInternalServerError, "Database error: failed to fetch projects")
        log.Printf("Error querying projects: %v", err)
        return
    }
    defer rows.Close()

    var projects []models.Project
    for rows.Next() {
        var project models.Project
        var categories pq.StringArray // Use pq.StringArray to scan the array field

        // Scan the project row
        if err := rows.Scan(&project.ID, &project.Name, &categories); err != nil {
            log.Printf("Error scanning project row: %v", err)
            log.Printf("Row data -> ID: %v, Name: %v, Categories: %v", project.ID, project.Name, categories)
            c.String(http.StatusInternalServerError, "Failed to scan projects")
            return
        }

        // Assign categories to the project and add to the list
        project.Categories = categories
        projects = append(projects, project)
    }

    // Check for iteration errors
    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        c.String(http.StatusInternalServerError, "Error iterating over project rows")
        return
    }

    // Prepare data for the template
    data := gin.H{
        "Title":    "Dashboard",
        "Projects": projects,
        "Username": username, // Include username for rendering
    }

    // Render the dashboard
    Render(c, "dashboard.html", data)
}
