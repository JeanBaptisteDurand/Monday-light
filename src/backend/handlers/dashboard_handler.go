package handlers

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/lib/pq"
    "monday-light/db"
    "monday-light/models"
)

func ShowDashboard(c *gin.Context) {
    userIDVal, _ := c.Get("userID")
    userID := userIDVal.(int)

    // Get user's username
    var username string
    err := db.DB.QueryRow("SELECT username FROM users WHERE id=$1", userID).Scan(&username)
    if err != nil {
        log.Printf("Error fetching username for userID=%d: %v", userID, err)
        c.String(http.StatusInternalServerError, "Database error: failed to fetch username")
        return
    }

    rows, err := db.DB.Query("SELECT id, name, categories FROM projects")
    if err != nil {
        log.Printf("Error querying projects: %v", err)
        c.String(http.StatusInternalServerError, "Database error: failed to fetch projects")
        return
    }
    defer rows.Close()

    var projects []models.Project
    for rows.Next() {
        var project models.Project
        var categories pq.StringArray
        if err := rows.Scan(&project.ID, &project.Name, &categories); err != nil {
            log.Printf("Error scanning project row: %v", err)
            c.String(http.StatusInternalServerError, "Failed to scan projects")
            return
        }
        project.Categories = categories
        projects = append(projects, project)
    }
    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        c.String(http.StatusInternalServerError, "Error iterating over project rows")
        return
    }

    data := gin.H{
        "Title":           "Dashboard",
        "Username":        username,
        "Projects":        projects,
        "ContentTemplate": "dashboard_content", // snippet name in dashboard.html
    }

    Render(c, data)
}
