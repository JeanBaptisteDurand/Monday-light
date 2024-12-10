package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"monday-light/db"
	"monday-light/models"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// RenderP chooses the template rendering strategy depending on HTMX request or not.
func RenderP(c *gin.Context, contentTemplate string, data gin.H) {
    if c.GetHeader("HX-Request") != "" {
        // HTMX request: parse only content template (and related partials)
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
        // Normal request: parse base layout + sidebar + partials
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

// CreateProject creates a new project with default categories.
func CreateProject(c *gin.Context) {
	name := c.PostForm("name")
	if strings.TrimSpace(name) == "" {
		log.Printf("Project name is empty.")
		c.String(http.StatusBadRequest, "Project name is required")
		return
	}

	// Default categories
	defaultCategories := []string{
		"high priority", "mid priority", "low priority",
		"urgent", "Frontend", "Backend",
		"RDV", "Communication", "Marketing", "Design",
	}

	var projectID int
	query := "INSERT INTO projects (name, categories) VALUES ($1, $2) RETURNING id"
	log.Printf("Executing query: %s with name=%s and default categories", query, name)
	err := db.DB.QueryRow(query, name, pq.Array(defaultCategories)).Scan(&projectID)
	if err != nil {
		log.Printf("Error creating project: %v", err)
		c.String(http.StatusInternalServerError, "Database error (create project)")
		return
	}

	log.Printf("Project created with ID %d and default categories", projectID)
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
        if err := rows.Scan(&p.ID, &p.Name, pq.Array(&categories)); err != nil {
            log.Printf("Error scanning project row: %v", err)
            c.String(http.StatusInternalServerError, "Error scanning projects")
            return
        }
        p.Categories = categories
        projects = append(projects, p)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error with rows iteration: %v", err)
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

func ShowNewProjectForm(c *gin.Context) {
	c.String(http.StatusOK, `
    <div class="modal fade" id="newProjectModal" tabindex="-1" aria-labelledby="newProjectModalLabel" aria-hidden="true" data-bs-backdrop="true" data-bs-keyboard="true">
        <div class="modal-dialog modal-dialog-centered">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="newProjectModalLabel">Créer un projet</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <form hx-post="/project" hx-target="#project-list-container" hx-swap="innerHTML">
                        <div class="mb-3">
                            <input type="text" name="name" class="form-control" placeholder="Nom du projet" required>
                        </div>
                        <div class="text-end">
                            <button type="submit" class="btn btn-primary">Créer</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
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

	query := "UPDATE projects SET categories = array_append(COALESCE(categories, '{}'), $1) WHERE id = $2"
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

// RenderCategories refreshes and renders the project categories partial.
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

	data := gin.H{"ID": projectID, "Categories": categories}
	RenderP(c, "project_categories.html", data)
}
