package handlers

import (
    "net/http"
    "strconv" // For string-to-int conversion
    "monday-light/models"
    "github.com/gin-gonic/gin"
)

var projects []models.Project

func ShowDashboard(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{
        "Projects": projects,
    })
}

func ShowProject(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr) // Convert string to int
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    for _, project := range projects {
        if project.ID == id {
            c.HTML(http.StatusOK, "project.html", gin.H{
                "Project": project,
            })
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
}

func CreateProject(c *gin.Context) {
    var newProject models.Project
    if err := c.BindJSON(&newProject); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    newProject.ID = len(projects) + 1
    projects = append(projects, newProject)
    c.JSON(http.StatusCreated, newProject)
}

func AddCategory(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr) // Convert string to int
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    var category struct{ Name string }
    if err := c.BindJSON(&category); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    for i, project := range projects {
        if project.ID == id {
            projects[i].Categories = append(project.Categories, category.Name)
            c.JSON(http.StatusOK, project)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
}