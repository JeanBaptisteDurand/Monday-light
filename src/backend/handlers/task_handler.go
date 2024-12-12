package handlers

import (
    "database/sql"
    "html/template"
    "log"
    "math"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    //"github.com/lib/pq"
    "monday-light/db"
    "monday-light/models"
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

    if name == "" {
        c.String(http.StatusBadRequest, "Nom de tâche requis")
        return
    }

    var estimatedTime int
    if estimatedTimeStr != "" {
        estimatedTime, _ = strconv.Atoi(estimatedTimeStr)
    }

    // Par défaut: status = "backlog"
    // L'utilisateur a demandé "la tâche a par défaut un status backlog"
    // (Dans votre code initial, c'était 'to_assign', ajustez selon vos besoins)
    err = db.DB.QueryRow(`
        INSERT INTO tasks (name, description, category, project_id, status, estimated_time, real_time)
        VALUES ($1, $2, $3, $4, 'backlog', $5, 0) RETURNING id
    `, name, description, category, projectID, estimatedTime).Scan(new(int))
    if err != nil {
        c.String(http.StatusInternalServerError, "Database error (create task)")
        return
    }

    ShowProject(c)
}

// GetTaskDetail returns the HTML for the task detail popup
func GetTaskDetail(c *gin.Context) {
    projectID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid project ID")
        return
    }

    taskID, err := strconv.Atoi(c.Param("task_id"))
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid task ID")
        return
    }

    var t models.Task
    query := `SELECT id, name, description, category, project_id, status, estimated_time, real_time, created_at, taken_from
              FROM tasks WHERE id=$1 AND project_id=$2`
    err = db.DB.QueryRow(query, taskID, projectID).Scan(&t.ID, &t.Name, &t.Description, &t.Category, &t.ProjectID, &t.Status, &t.EstimatedTime, &t.RealTime, &t.CreatedAt, &t.TakenFrom)
    if err != nil {
        if err == sql.ErrNoRows {
            c.String(http.StatusNotFound, "Tâche introuvable")
        } else {
            c.String(http.StatusInternalServerError, "Database error (get task)")
        }
        return
    }

    // Fetch assigned users
    t.AssignedUsers = getAssignedUsers(taskID)

    // Calculate real time if assigned or to_check
    if t.Status == "assigned" || t.Status == "to_check" {
        if !t.TakenFrom.IsZero() {
            diff := time.Since(t.TakenFrom).Minutes()
            t.RealTime = int(math.Floor(diff))
        }
    }

    // Compute progress percentage in code
    var progress int
    if t.EstimatedTime > 0 {
        progress = int(math.Floor(float64(t.RealTime) / float64(t.EstimatedTime) * 100))
        if progress > 100 {
            progress = 100
        }
        if progress < 0 {
            progress = 0
        }
    } else {
        progress = 0
    }

    data := gin.H{
        "Task":     t,
        "Progress": progress,
    }

    tmpl, err := template.ParseFiles("templates/task_detail.html")
    if err != nil {
        c.String(http.StatusInternalServerError, err.Error())
        return
    }

    c.Header("Content-Type", "text/html; charset=utf-8")
    tmpl.Execute(c.Writer, data)
}

// NextTaskStatus advances the task status along the workflow
func NextTaskStatus(c *gin.Context) {
    projectID, _ := strconv.Atoi(c.Param("id"))
    taskID, _ := strconv.Atoi(c.Param("task_id"))

    var status string
    err := db.DB.QueryRow("SELECT status FROM tasks WHERE id=$1 AND project_id=$2", taskID, projectID).Scan(&status)
    if err != nil {
        c.String(http.StatusNotFound, "Tâche introuvable")
        return
    }

    next := nextStatus(status)
    if next == "" {
        c.String(http.StatusBadRequest, "Pas de statut suivant")
        return
    }

    // If moving to assigned, set taken_from = NOW() if not set
    if next == "assigned" {
        _, err = db.DB.Exec("UPDATE tasks SET status=$1, taken_from=NOW() WHERE id=$2 AND project_id=$3", next, taskID, projectID)
    } else if next == "done" {
        // When moving to done, fix the real_time as final
        // real_time = difference (in minutes) between taken_from and NOW()
        // only if taken_from is not zero
        var takenFrom time.Time
        err = db.DB.QueryRow("SELECT taken_from FROM tasks WHERE id=$1", taskID).Scan(&takenFrom)
        if err == nil && !takenFrom.IsZero() {
            diff := time.Since(takenFrom).Minutes()
            realTime := int(math.Floor(diff))
            _, err = db.DB.Exec("UPDATE tasks SET status=$1, real_time=$2 WHERE id=$3 AND project_id=$4", next, realTime, taskID, projectID)
        } else {
            _, err = db.DB.Exec("UPDATE tasks SET status=$1 WHERE id=$2 AND project_id=$3", next, taskID, projectID)
        }
    } else {
        _, err = db.DB.Exec("UPDATE tasks SET status=$1 WHERE id=$2 AND project_id=$3", next, taskID, projectID)
    }

    if err != nil {
        c.String(http.StatusInternalServerError, "Erreur BDD changement de statut")
        return
    }

    // Return updated detail
    GetTaskDetail(c)
}

// AssignToSelf assigns the current user to a 'to_assign' task and moves it to 'assigned'
func AssignToSelf(c *gin.Context) {
    userIDVal, _ := c.Get("userID")
    userID := userIDVal.(int)

    projectID, _ := strconv.Atoi(c.Param("id"))
    taskID, _ := strconv.Atoi(c.Param("task_id"))

    var status string
    err := db.DB.QueryRow("SELECT status FROM tasks WHERE id=$1 AND project_id=$2", taskID, projectID).Scan(&status)
    if err != nil {
        c.String(http.StatusNotFound, "Tâche introuvable")
        return
    }

    if status != "to_assign" {
        c.String(http.StatusBadRequest, "Tâche non assignable à ce stade")
        return
    }

    // Insert user-task assignment
    _, err = db.DB.Exec("INSERT INTO user_tasks (user_id, task_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, taskID)
    if err != nil {
        c.String(http.StatusInternalServerError, "Erreur d'assignation")
        return
    }

    // Passer la tâche en "assigned" et prendre le temps
    _, err = db.DB.Exec("UPDATE tasks SET status='assigned', taken_from=NOW() WHERE id=$1 AND project_id=$2", taskID, projectID)
    if err != nil {
        c.String(http.StatusInternalServerError, "Erreur BDD lors du changement de statut")
        return
    }

    GetTaskDetail(c)
}

func getAssignedUsers(taskID int) []int {
    rows, err := db.DB.Query("SELECT user_id FROM user_tasks WHERE task_id=$1", taskID)
    if err != nil {
        log.Println("Error fetching assigned users:", err)
        return nil
    }
    defer rows.Close()

    var users []int
    for rows.Next() {
        var uid int
        if err := rows.Scan(&uid); err == nil {
            users = append(users, uid)
        }
    }
    return users
}

func nextStatus(current string) string {
    order := []string{"backlog", "to_assign", "assigned", "to_check", "done"}
    for i, s := range order {
        if s == current && i < len(order)-1 {
            return order[i+1]
        }
    }
    return ""
}
