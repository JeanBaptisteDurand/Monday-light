package handlers

import (
    "log"
    "net/http"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/lib/pq"
    "monday-light/db"
    "monday-light/models"
)

// ShowProject
func ShowProject(c *gin.Context) {
    userIDVal, _ := c.Get("userID")
    userID := userIDVal.(int)

    // old logging
    log.Printf("ShowProject: userID=%d", userID)

    var username string
    err := db.DB.QueryRow("SELECT username FROM users WHERE id=$1", userID).Scan(&username)
    if err != nil {
        log.Printf("Error fetching username: %v", err)
        c.String(http.StatusInternalServerError, "Database error (username)")
        return
    }

    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid project ID")
        return
    }
    log.Printf("ProjectID param: %d", projectID)

    var project models.Project
    var categories []string
    query := "SELECT id, name, categories FROM projects WHERE id=$1"
    log.Printf("Query: %s with projectID=%d", query, projectID)
    err = db.DB.QueryRow(query, projectID).Scan(&project.ID, &project.Name, pq.Array(&categories))
    if err != nil {
        log.Printf("Error fetching project: %v", err)
        c.String(http.StatusInternalServerError, "Database error (project)")
        return
    }
    project.Categories = categories

    rows, err := db.DB.Query("SELECT id, name, description, category, status, estimated_time, real_time FROM tasks WHERE project_id=$1", projectID)
    if err != nil {
        log.Printf("Error fetching tasks: %v", err)
        c.String(http.StatusInternalServerError, "Database error (tasks)")
        return
    }
    defer rows.Close()

    var backlogTasks, enCoursTasks, doneTasks []models.Task
    for rows.Next() {
        var t models.Task
        if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Category, &t.Status, &t.EstimatedTime, &t.RealTime); err != nil {
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
        "Title":          project.Name,
        "Username":       username,
        "Project":        project,
        "BacklogTasks":   backlogTasks,
        "EnCoursTasks":   enCoursTasks,
        "DoneTasks":      doneTasks,
        "ContentTemplate": "project_content",
    }

    Render(c, data)
}

// CreateProject + reload sidebar
func CreateProject(c *gin.Context) {
    name := c.PostForm("name")
    if strings.TrimSpace(name) == "" {
        log.Printf("Project name is empty.")
        c.String(http.StatusBadRequest, "Project name is required")
        return
    }

    defaultCategories := []string{
        "high priority", "mid priority", "low priority",
        "urgent", "Frontend", "Backend",
        "RDV", "Communication", "Marketing", "Design",
    }

    var projectID int
    query := "INSERT INTO projects (name, categories) VALUES ($1, $2) RETURNING id"
    log.Printf("Executing query: %s with name=%s", query, name)
    err := db.DB.QueryRow(query, name, pq.Array(defaultCategories)).Scan(&projectID)
    if err != nil {
        log.Printf("Error creating project: %v", err)
        c.String(http.StatusInternalServerError, "Database error (create project)")
        return
    }
    log.Printf("Project created with ID %d", projectID)

    RenderProjectList(c)
}

// RenderProjectList refreshes the "sidebar_projects" partial
func RenderProjectList(c *gin.Context) {
    query := "SELECT id, name, categories FROM projects"
    log.Printf("Executing query: %s", query)
    rows, err := db.DB.Query(query)
    if err != nil {
        log.Printf("Error: %v", err)
        c.String(http.StatusInternalServerError, "Database error (projects)")
        return
    }
    defer rows.Close()

    var projects []models.Project
    for rows.Next() {
        var p models.Project
        var cats []string
        if err := rows.Scan(&p.ID, &p.Name, pq.Array(&cats)); err != nil {
            log.Printf("Error scanning project row: %v", err)
            c.String(http.StatusInternalServerError, "Error scanning projects")
            return
        }
        p.Categories = cats
        projects = append(projects, p)
    }

    c.HTML(http.StatusOK, "sidebar_projects", gin.H{"Projects": projects})
}

func ShowNewProjectForm(c *gin.Context) {
    // old logic
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

// AddCategory => update categories, then re-render
func AddCategory(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid project ID")
        return
    }

    cat := c.PostForm("category_name")
    if strings.TrimSpace(cat) == "" {
        c.String(http.StatusBadRequest, "Category name is required")
        return
    }

    _, err = db.DB.Exec("UPDATE projects SET categories = array_append(COALESCE(categories, '{}'), $1) WHERE id=$2",
        cat, projectID)
    if err != nil {
        log.Printf("Error adding category: %v", err)
        c.String(http.StatusInternalServerError, "Database error (add category)")
        return
    }

    RenderCategories(c, projectID)
}

// RemoveCategory => remove, then re-render
func RemoveCategory(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid project ID")
        return
    }

    cat := c.PostForm("cat")
    if cat == "" {
        c.String(http.StatusBadRequest, "Category is required")
        return
    }

    _, err = db.DB.Exec("UPDATE projects SET categories = array_remove(categories, $1) WHERE id=$2",
        cat, projectID)
    if err != nil {
        log.Printf("Error removing category: %v", err)
        c.String(http.StatusInternalServerError, "Database error (remove category)")
        return
    }

    RenderCategories(c, projectID)
}

func RenderCategories(c *gin.Context, projectID int) {
    var cats []string
    err := db.DB.QueryRow("SELECT categories FROM projects WHERE id=$1", projectID).Scan(pq.Array(&cats))
    if err != nil {
        log.Printf("Error fetching categories: %v", err)
        c.String(http.StatusInternalServerError, "Database error (render categories)")
        return
    }

    data := gin.H{"ID": projectID, "Categories": cats}
    c.HTML(http.StatusOK, "project_categories", data)
}
