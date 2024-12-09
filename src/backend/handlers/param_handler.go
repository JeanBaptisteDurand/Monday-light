package handlers

import (
    "database/sql"
    "html/template"
    "net/http"
    "strings"

    "monday-light/db"
    "monday-light/models"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

// ShowParam displays the user parameters
func ShowParam(c *gin.Context) {
    userIDVal, _ := c.Get("userID")
    userID := userIDVal.(int)

    var user models.User
    err := db.DB.QueryRow(`
        SELECT id, username, email, discord_id, discord_pseudo, color 
        FROM users 
        WHERE id=$1
    `, userID).Scan(&user.ID, &user.Username, &user.Email, &user.DiscordID, &user.DiscordPseudo, &user.Color)
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to load user info")
        return
    }

    data := gin.H{
        "User":     user,
        "Username": user.Username,
        "Title":    "Paramètres",
    }

    Render(c, "param.html", data)
}

// ShowParamEdit returns a form to edit a specific field
func ShowParamEdit(c *gin.Context) {
    field := c.Query("field")
    if field == "" {
        c.String(http.StatusBadRequest, "Field required")
        return
    }

    userIDVal, _ := c.Get("userID")
    userID := userIDVal.(int)

    var user models.User
    err := db.DB.QueryRow(`SELECT username, email, discord_id, discord_pseudo FROM users WHERE id=$1`, userID).
        Scan(&user.Username, &user.Email, &user.DiscordID, &user.DiscordPseudo)
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to load user info")
        return
    }

    var currentValue string
    switch field {
    case "username":
        currentValue = user.Username
    case "email":
        currentValue = user.Email
    case "discord_id":
        currentValue = user.DiscordID
    case "discord_pseudo":
        currentValue = user.DiscordPseudo
    case "password":
        // password has no current value
        currentValue = ""
    default:
        c.String(http.StatusBadRequest, "Invalid field")
        return
    }

    tmpl, err := template.ParseFiles("templates/param_edit_field.html")
    if err != nil {
        c.String(http.StatusInternalServerError, err.Error())
        return
    }
    data := gin.H{
        "Field":        field,
        "CurrentValue": currentValue,
    }

    c.Status(http.StatusOK)
    c.Header("Content-Type", "text/html; charset=utf-8")
    if err := tmpl.ExecuteTemplate(c.Writer, "content", data); err != nil {
        c.String(http.StatusInternalServerError, err.Error())
    }
}

// UpdateParam updates the specified user field
func UpdateParam(c *gin.Context) {
    userIDVal, _ := c.Get("userID")
    userID := userIDVal.(int)

    field := c.PostForm("field")
    if field == "" {
        c.String(http.StatusBadRequest, "Field required")
        return
    }

    if field == "password" {
        // Handle password change
        oldPassword := c.PostForm("old_password")
        newPassword := c.PostForm("new_password")
        confirmPassword := c.PostForm("confirm_password")

        if newPassword != confirmPassword {
            c.String(http.StatusBadRequest, "Les mots de passe ne correspondent pas.")
            return
        }

        // Fetch current password hash
        var currentHash string
        err := db.DB.QueryRow("SELECT password_hash FROM users WHERE id=$1", userID).Scan(&currentHash)
        if err != nil {
            c.String(http.StatusInternalServerError, "Failed to load user info")
            return
        }

        // Check old password
        if bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(oldPassword)) != nil {
            c.String(http.StatusUnauthorized, "Ancien mot de passe incorrect.")
            return
        }

        // Hash new password
        hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
        if err != nil {
            c.String(http.StatusInternalServerError, "Erreur interne, réessayez plus tard.")
            return
        }

        // Update in DB
        _, err = db.DB.Exec("UPDATE users SET password_hash=$1 WHERE id=$2", string(hashed), userID)
        if err != nil {
            c.String(http.StatusInternalServerError, "Échec de la mise à jour du mot de passe.")
            return
        }

        // Re-render the password field line
        renderUpdatedField(c, userID, field)
        return
    }

    // For other fields
    value := strings.TrimSpace(c.PostForm("value"))
    if value == "" {
        c.String(http.StatusBadRequest, "Value required")
        return
    }

    var query string
    switch field {
    case "username":
        query = "UPDATE users SET username=$1 WHERE id=$2"
    case "email":
        query = "UPDATE users SET email=$1 WHERE id=$2"
    case "discord_id":
        query = "UPDATE users SET discord_id=$1 WHERE id=$2"
    case "discord_pseudo":
        query = "UPDATE users SET discord_pseudo=$1 WHERE id=$2"
    default:
        c.String(http.StatusBadRequest, "Invalid field")
        return
    }

    _, err := db.DB.Exec(query, value, userID)
    if err != nil {
        if err == sql.ErrNoRows {
            c.String(http.StatusNotFound, "Utilisateur non trouvé.")
            return
        }
        c.String(http.StatusInternalServerError, "Échec de la mise à jour.")
        return
    }

    if field == "username" {
        // Reload entire page so that pseudo is updated in the navbar
        c.Header("HX-Trigger", "pseudoSuccess")
        c.Status(http.StatusOK)
        return
    }
    
    // For other fields, send a trigger on success
    c.Header("HX-Trigger", "paramSuccess")
    renderUpdatedField(c, userID, field)    
}

// renderUpdatedField re-renders a single field line after update
func renderUpdatedField(c *gin.Context, userID int, field string) {
    var user models.User
    err := db.DB.QueryRow(`SELECT username, email, discord_id, discord_pseudo FROM users WHERE id=$1`, userID).
        Scan(&user.Username, &user.Email, &user.DiscordID, &user.DiscordPseudo)
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to load user info")
        return
    }

    // Based on field, return a <dd> line just like original
    var value string
    switch field {
    case "username":
        value = user.Username
    case "email":
        value = user.Email
    case "discord_id":
        value = user.DiscordID
    case "discord_pseudo":
        value = user.DiscordPseudo
    case "password":
        value = "********"
    default:
        c.String(http.StatusBadRequest, "Invalid field")
        return
    }

    // Rebuild the line with edit icon
    // This should match the structure in param.html
    dd := `<dd class="col-sm-8 d-flex justify-content-between align-items-center" id="` + field + `-field">
            <span>` + template.HTMLEscapeString(value) + `</span>
            <i class="bi bi-pencil-square text-info"
               style="cursor:pointer;"
               hx-get="/param/edit?field=` + field + `"
               hx-target="#` + field + `-field"
               hx-swap="outerHTML"></i>
        </dd>`

    c.Header("Content-Type", "text/html; charset=utf-8")
    c.String(http.StatusOK, dd)
}
