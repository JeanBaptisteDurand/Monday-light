// file: handlers/user_handler.go
package handlers

import (
    "database/sql"
    "fmt"
    "math/rand"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "monday-light/db"
    "monday-light/models"
)

// Predefined colors for users
var predefinedColors = []string{
    "#FF5733", "#33FF57", "#3357FF", "#F3FF33",
    "#FF33F6", "#33FFF6", "#F633FF", "#FFC300",
    "#FF5733", "#DAF7A6",
}

// getUsedColors fetches all distinct colors used so far in the DB
func getUsedColors() ([]string, error) {
    rows, err := db.DB.Query("SELECT DISTINCT color FROM users WHERE color IS NOT NULL")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var usedColors []string
    for rows.Next() {
        var color string
        if err := rows.Scan(&color); err != nil {
            return nil, err
        }
        usedColors = append(usedColors, color)
    }
    return usedColors, nil
}

// getRandomUnusedColor picks a random color from predefinedColors that isn't used yet.
func getRandomUnusedColor() (string, error) {
    usedColors, err := getUsedColors()
    if err != nil {
        return "", err
    }

    usedSet := make(map[string]bool)
    for _, c := range usedColors {
        usedSet[c] = true
    }

    rand.Seed(time.Now().UnixNano())
    shuffled := rand.Perm(len(predefinedColors)) // shuffle indices

    for _, i := range shuffled {
        candidate := predefinedColors[i]
        if !usedSet[candidate] {
            return candidate, nil
        }
    }
    return "", fmt.Errorf("no colors available")
}

// ShowLogin => Render "login.html"
func ShowLogin(c *gin.Context) {
    errorParam := c.Query("error")
    c.HTML(http.StatusOK, "login.html", gin.H{
        "error": errorParam == "1",
    })
}

// ShowRegister => Render "register.html"
func ShowRegister(c *gin.Context) {
    errorParam := c.Query("error")
    c.HTML(http.StatusOK, "register.html", gin.H{
        "error": errorParam == "1",
    })
}

// Register handles new user registration, including the random color assignment.
func Register(c *gin.Context) {
    type Input struct {
        Username       string `form:"username" binding:"required"`
        Email          string `form:"email" binding:"required"`
        Password       string `form:"password" binding:"required"`
        DiscordID      string `form:"discord_id"`
        DiscordPseudo  string `form:"discord_pseudo"`
    }

    var input Input
    if err := c.ShouldBind(&input); err != nil {
        // If it's an HTMX request, return partial HTML error block
        if c.GetHeader("HX-Request") != "" {
            c.String(http.StatusBadRequest, `<div class="alert alert-danger">Données invalides.</div>`)
            return
        }
        c.Redirect(http.StatusFound, "/register?error=1")
        return
    }

    // Hash the password
    hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        if c.GetHeader("HX-Request") != "" {
            c.String(http.StatusInternalServerError, `<div class="alert alert-danger">Erreur interne.</div>`)
            return
        }
        c.Redirect(http.StatusFound, "/register?error=1")
        return
    }

    // Pick a random unused color
    color, err := getRandomUnusedColor()
    if err != nil {
        c.String(http.StatusInternalServerError, "No colors available")
        return
    }

    var userID int
    err = db.DB.QueryRow(`
        INSERT INTO users (username, email, password_hash, discord_id, discord_pseudo, color)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `, input.Username, input.Email, string(hashed), input.DiscordID, input.DiscordPseudo, color).
    Scan(&userID)
    if err != nil {
        // Possibly email or username already taken
        if c.GetHeader("HX-Request") != "" {
            c.String(http.StatusConflict, `<div class="alert alert-danger">Impossible de créer l'utilisateur (Email/Username déjà pris?).</div>`)
            return
        }
        c.Redirect(http.StatusFound, "/register?error=1")
        return
    }

    // On success, if HTMX => instruct client to redirect
    if c.GetHeader("HX-Request") != "" {
        c.Header("HX-Redirect", "/login")
        c.Status(http.StatusOK)
        return
    }

    // Normal request => redirect
    c.Redirect(http.StatusFound, "/login")
}

// Login checks user credentials, issues JWT, sets cookie.
func Login(c *gin.Context) {
    type LoginInput struct {
        Email    string `form:"email" binding:"required"`
        Password string `form:"password" binding:"required"`
    }

    var input LoginInput
    if err := c.ShouldBind(&input); err != nil {
        if c.GetHeader("HX-Request") != "" {
            c.String(http.StatusBadRequest, `<div class="alert alert-danger">Données invalides.</div>`)
            return
        }
        c.Redirect(http.StatusFound, "/login?error=1")
        return
    }

    var user models.User
    err := db.DB.QueryRow("SELECT id, password_hash FROM users WHERE email=$1", input.Email).
        Scan(&user.ID, &user.PasswordHash)
    if err != nil {
        if err == sql.ErrNoRows {
            // invalid credentials
            if c.GetHeader("HX-Request") != "" {
                c.String(http.StatusUnauthorized, `<div class="alert alert-danger">Email ou mot de passe invalide.</div>`)
                return
            }
            c.Redirect(http.StatusFound, "/login?error=1")
        } else {
            // internal error
            if c.GetHeader("HX-Request") != "" {
                c.String(http.StatusInternalServerError, `<div class="alert alert-danger">Erreur interne.</div>`)
                return
            }
            c.Redirect(http.StatusFound, "/login?error=1")
        }
        return
    }

    // Check password
    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
        // invalid password
        if c.GetHeader("HX-Request") != "" {
            c.String(http.StatusUnauthorized, `<div class="alert alert-danger">Email ou mot de passe invalide.</div>`)
            return
        }
        c.Redirect(http.StatusFound, "/login?error=1")
        return
    }

    // All good => generate JWT
    token, err := GenerateJWT(user.ID)
    if err != nil {
        if c.GetHeader("HX-Request") != "" {
            c.String(http.StatusInternalServerError, `<div class="alert alert-danger">Erreur interne.</div>`)
            return
        }
        c.Redirect(http.StatusFound, "/login?error=1")
        return
    }

    // Set JWT cookie
    c.SetCookie("token", token, 3600*24, "/", "", false, true)

    if c.GetHeader("HX-Request") != "" {
        c.Header("HX-Redirect", "/")
        c.Status(http.StatusOK)
        return
    }
    // Normal => redirect
    c.Redirect(http.StatusFound, "/")
}
