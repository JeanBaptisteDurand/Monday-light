package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"monday-light/db"
	"monday-light/models"
	"github.com/gin-gonic/gin"
)

// ShowProject displays a project page with its categories and tasks.
func ShowProject(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(int)

	var username string
	query := "SELECT username FROM users WHERE id = $1"
	log.Printf("Executing query: %s with userID=%d", query, userID)
	err := db.DB.QueryRow(query, userID).Scan(&username)
	if err != nil {
		log.Printf("Error fetching username: %v", err)
		c.String(http.StatusInternalServerError, "Database error (username)")
		return
	}

	idStr := c.Param("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid project ID: %s", idStr)
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	var project models.Project
	var categories []string
	query = "SELECT id, name, categories FROM projects WHERE id = $1"
	log.Printf("Executing query: %s with projectID=%d", query, projectID)
	err = db.DB.QueryRow(query, projectID).Scan(&project.ID, &project.Name, pq.Array(&categories))
	if err != nil {
		log.Printf("Error fetching project: %v", err)
		c.String(http.StatusInternalServerError, "Database error (project)")
		return
	}
	project.Categories = categories

	query = "SELECT id, name, description, category, status, estimated_time, real_time FROM tasks WHERE project_id=$1"
	log.Printf("Executing query: %s with projectID=%d", query, projectID)
	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		c.String(http.StatusInternalServerError, "Database error (tasks)")
		return
	}
	defer rows.Close()

	var backlogTasks, enCoursTasks, doneTasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Category, &t.Status, &t.EstimatedTime, &t.RealTime)
		if err != nil {
			log.Printf("Error scanning task row: %v", err)
			c.String(http.StatusInternalServerError, "Error scanning tasks")
			return
		}
		switch t.Status {
		case "backlog":
			backlogTasks = append(backlogTasks, t)
		case "done":
			doneTasks = append(doneTasks, t)
		default:
			enCoursTasks = append(enCoursTasks, t)
		}
	}

	data := gin.H{
		"Title":         project.Name,
		"Project":       project,
		"Username":      username,
		"EnCoursTasks":  enCoursTasks,
		"BacklogTasks":  backlogTasks,
		"DoneTasks":     doneTasks,
	}

	RenderP(c, "project.html", data)
}

// CreateProject creates a new project.
func CreateProject(c *gin.Context) {
	name := c.PostForm("name")
	if strings.TrimSpace(name) == "" {
		log.Printf("Project name is empty.")
		c.String(http.StatusBadRequest, "Project name is required")
		return
	}

	var projectID int
	query := "INSERT INTO projects (name) VALUES ($1) RETURNING id"
	log.Printf("Executing query: %s with name=%s", query, name)
	err := db.DB.QueryRow(query, name).Scan(&projectID)
	if err != nil {
		log.Printf("Error creating project: %v", err)
		c.String(http.StatusInternalServerError, "Database error (create project)")
		return
	}

	RenderProjectList(c)
}

// RenderProjectList refreshes and renders the project list for the sidebar.
func RenderProjectList(c *gin.Context) {
    query := "SELECT id, name, categories FROM projects"
    log.Printf("Executing query: %s", query)
    rows, err := db.DB.Query(query)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        c.String(http.StatusInternalServerError, "Database error (projects)")
        return
    }
    defer rows.Close()

    var projects []models.Project
    for rows.Next() {
        var p models.Project
        var categories []string
        // Adding detailed logging here
        if err := rows.Scan(&p.ID, &p.Name, pq.Array(&categories)); err != nil {
            log.Printf("Error scanning project row: %v", err) // Log the exact error
            log.Printf("Row Data -> ID: %v, Name: %v, Categories: %v", p.ID, p.Name, categories)
            c.String(http.StatusInternalServerError, "Error scanning projects")
            return
        }
        p.Categories = categories
        projects = append(projects, p)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error with rows iteration: %v", err) // Log if there's an error during iteration
        c.String(http.StatusInternalServerError, "Error iterating over projects")
        return
    }

    tmpl, err := template.ParseFiles("templates/sidebar_projects.html")
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        c.String(http.StatusInternalServerError, err.Error())
        return
    }

    data := gin.H{"Projects": projects}
    c.Status(http.StatusOK)
    c.Header("Content-Type", "text/html; charset=utf-8")
    if err := tmpl.ExecuteTemplate(c.Writer, "sidebar_projects", data); err != nil {
        log.Printf("Error executing template: %v", err)
        c.String(http.StatusInternalServerError, "Error rendering template")
    }
}

// ShowNewProjectForm displays a form for adding a new project.
func ShowNewProjectForm(c *gin.Context) {
	c.String(http.StatusOK, `
    <div class="px-3 py-2">
        <form hx-post="/project" hx-target="#project-list-container" hx-swap="innerHTML">
            <div class="input-group input-group-sm mb-2">
                <input type="text" name="name" class="form-control" placeholder="Nom du projet" required>
                <button class="btn btn-primary">Créer</button>
            </div>
        </form>
    </div>
    `)
}

// AddCategory adds a new category to a project.
func AddCategory(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid project ID: %s", c.Param("id"))
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	categoryName := c.PostForm("category_name")
	if strings.TrimSpace(categoryName) == "" {
		log.Printf("Category name is empty.")
		c.String(http.StatusBadRequest, "Category name is required")
		return
	}

	query := "UPDATE projects SET categories = array_append(categories, $1) WHERE id = $2"
	log.Printf("Executing query: %s with categoryName=%s, projectID=%d", query, categoryName, projectID)
	_, err = db.DB.Exec(query, categoryName, projectID)
	if err != nil {
		log.Printf("Error adding category: %v", err)
		c.String(http.StatusInternalServerError, "Database error (add category)")
		return
	}

	RenderCategories(c, projectID)
}

// RemoveCategory removes a category from a project.
func RemoveCategory(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid project ID: %s", c.Param("id"))
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	cat := c.PostForm("cat")
	if cat == "" {
		log.Printf("Category is empty.")
		c.String(http.StatusBadRequest, "Category is required")
		return
	}

	query := "UPDATE projects SET categories = array_remove(categories, $1) WHERE id = $2"
	log.Printf("Executing query: %s with category=%s, projectID=%d", query, cat, projectID)
	_, err = db.DB.Exec(query, cat, projectID)
	if err != nil {
		log.Printf("Error removing category: %v", err)
		c.String(http.StatusInternalServerError, "Database error (remove category)")
		return
	}

	RenderCategories(c, projectID)
}

// RenderCategories refreshes and renders the project categories.
func RenderCategories(c *gin.Context, projectID int) {
	var categories []string
	query := "SELECT categories FROM projects WHERE id=$1"
	log.Printf("Executing query: %s with projectID=%d", query, projectID)
	err := db.DB.QueryRow(query, projectID).Scan(pq.Array(&categories))
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		c.String(http.StatusInternalServerError, "Database error (render categories)")
		return
	}

	tmpl, err := template.ParseFiles("templates/project_categories.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := gin.H{"ProjectID": projectID, "Categories": categories}
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(c.Writer, "project_categories", data)
}
