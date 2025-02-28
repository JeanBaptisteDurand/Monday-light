package handlers

import (
    "database/sql"
    "net/http"
    "strings"
    "html"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "monday-light/db"
    "monday-light/models"
)

func ShowParam(c *gin.Context) {
    userID := c.GetInt("userID")

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
        "Title":           "Paramètres",
        "Username":        user.Username,
        "User":            user,
        "ContentTemplate": "param_content", // snippet in param.html
    }
    Render(c, data)
}

func ShowParamEdit(c *gin.Context) {
    field := c.Query("field")
    if field == "" {
        c.String(http.StatusBadRequest, "Field required")
        return
    }

    userID := c.GetInt("userID")

    var user models.User
    err := db.DB.QueryRow(`
        SELECT username, email, discord_id, discord_pseudo
        FROM users WHERE id=$1
    `, userID).Scan(&user.Username, &user.Email, &user.DiscordID, &user.DiscordPseudo)
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
        currentValue = ""
    default:
        c.String(http.StatusBadRequest, "Invalid field")
        return
    }

    // Render just the snippet "param_edit_field"
    c.HTML(http.StatusOK, "param_edit_field", gin.H{
        "Field":        field,
        "CurrentValue": currentValue,
    })
}

func UpdateParam(c *gin.Context) {
    userID := c.GetInt("userID")
    field := c.PostForm("field")
    if field == "" {
        c.String(http.StatusBadRequest, "Field required")
        return
    }

    if field == "password" {
        oldPass := c.PostForm("old_password")
        newPass := c.PostForm("new_password")
        confirm := c.PostForm("confirm_password")
        if newPass != confirm {
            c.String(http.StatusBadRequest, "Les mots de passe ne correspondent pas.")
            return
        }

        var currentHash string
        err := db.DB.QueryRow("SELECT password_hash FROM users WHERE id=$1", userID).Scan(&currentHash)
        if err != nil {
            c.String(http.StatusInternalServerError, "Failed to load user info")
            return
        }

        if bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(oldPass)) != nil {
            c.String(http.StatusUnauthorized, "Ancien mot de passe incorrect.")
            return
        }

        hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
        if err != nil {
            c.String(http.StatusInternalServerError, "Erreur interne, réessayez plus tard.")
            return
        }

        _, err = db.DB.Exec("UPDATE users SET password_hash=$1 WHERE id=$2", string(hashed), userID)
        if err != nil {
            c.String(http.StatusInternalServerError, "Échec de la mise à jour du mot de passe.")
            return
        }

        c.Header("HX-Trigger", "paramSuccess")
        renderUpdatedField(c, userID, field)
        return
    }

    // Other fields
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
        // reload entire page => update navbar
        c.Header("HX-Trigger", "pseudoSuccess")
        c.Status(http.StatusOK)
        return
    }

    c.Header("HX-Trigger", "paramSuccess")
    renderUpdatedField(c, userID, field)
}

func renderUpdatedField(c *gin.Context, userID int, field string) {
    var user models.User
    err := db.DB.QueryRow(`
        SELECT username, email, discord_id, discord_pseudo
        FROM users
        WHERE id=$1
    `, userID).Scan(&user.Username, &user.Email, &user.DiscordID, &user.DiscordPseudo)
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to load user info")
        return
    }

    // Use old logic with template.HTMLEscapeString if you want to ensure safety
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

    // If you want to escape it:
    //escapedValue := template.HTMLEscapeString(value)
    // But if you trust the stored data, you can leave it as is
    // We'll do it for safety:
    // We'll reintroduce the old approach:
    escapedValue := html.EscapeString(value)

    dd := `<dd class="col-sm-8 d-flex justify-content-between align-items-center" id="` + field + `-field">
            <span>` + escapedValue + `</span>
            <i class="bi bi-pencil-square text-info"
               style="cursor:pointer;"
               hx-get="/param/edit?field=` + field + `"
               hx-target="#` + field + `-field"
               hx-swap="outerHTML"></i>
        </dd>`

    c.Header("Content-Type", "text/html; charset=utf-8")
    c.String(http.StatusOK, dd)
}
