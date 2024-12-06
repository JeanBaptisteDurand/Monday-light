package handlers

import (
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"

    "html/template"
    "monday-light/db"
    "monday-light/models"
)

// RenderStandalone loads and executes a single template without "content" block
func RenderStandalone(c *gin.Context, filename string, data gin.H) {
    tmpl, err := template.ParseFiles("templates/" + filename)
    if err != nil {
        c.String(http.StatusInternalServerError, err.Error())
        return
    }
    c.Status(http.StatusOK)
    c.Header("Content-Type", "text/html; charset=utf-8")
    if err := tmpl.Execute(c.Writer, data); err != nil {
        c.String(http.StatusInternalServerError, err.Error())
    }
}

func ShowRegister(c *gin.Context) {
    // If there's an error param, show the error alert
    errorParam := c.Query("error")
    RenderStandalone(c, "register.html", gin.H{
        "error": errorParam == "1",
    })
}

func ShowLogin(c *gin.Context) {
    errorParam := c.Query("error")
    RenderStandalone(c, "login.html", gin.H{
        "error": errorParam == "1",
    })
}

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
            c.Status(http.StatusBadRequest)
            c.Writer.Write([]byte(`<div class="alert alert-danger">Données invalides. Veuillez vérifier les champs.</div>`))
            return
        }
        // Otherwise, redirect or show a full page error
        c.Redirect(http.StatusFound, "/register?error=1")
        return
    }

    // Check password confirmation via JS on frontend, but if needed, can verify again here:
    // Not provided in form since "password_confirm" is not posted to server, 
    // front-end ensures they match.

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        if c.GetHeader("HX-Request") != "" {
            c.Status(http.StatusInternalServerError)
            c.Writer.Write([]byte(`<div class="alert alert-danger">Erreur interne. Réessayez plus tard.</div>`))
            return
        }
        c.Redirect(http.StatusFound, "/register?error=1")
        return
    }

    var userID int
    err = db.DB.QueryRow(`
        INSERT INTO users (username, email, password_hash, discord_id, discord_pseudo)
        VALUES ($1, $2, $3, $4, $5) RETURNING id
    `, input.Username, input.Email, string(hashedPassword), input.DiscordID, input.DiscordPseudo).Scan(&userID)
    if err != nil {
        if c.GetHeader("HX-Request") != "" {
            // Possibly email or username already taken
            c.Status(http.StatusConflict)
            c.Writer.Write([]byte(`<div class="alert alert-danger">Impossible de créer l'utilisateur. Vérifiez si l'email ou le nom d'utilisateur est déjà pris.</div>`))
            return
        }
        c.Redirect(http.StatusFound, "/register?error=1")
        return
    }

    // On success, if HTMX: redirect user
    if c.GetHeader("HX-Request") != "" {
        c.Header("HX-Redirect", "/login")
        c.Status(http.StatusOK)
        return
    }

    // normal request
    c.Redirect(http.StatusFound, "/login")
}

func Login(c *gin.Context) {
    type LoginInput struct {
        Email    string `form:"email" binding:"required"`
        Password string `form:"password" binding:"required"`
    }

    var input LoginInput
    if err := c.ShouldBind(&input); err != nil {
        if c.GetHeader("HX-Request") != "" {
            c.Status(http.StatusBadRequest)
            c.Writer.Write([]byte(`<div class="alert alert-danger">Données invalides. Veuillez vérifier les champs.</div>`))
            return
        }
        c.Redirect(http.StatusFound, "/login?error=1")
        return
    }

    var user models.User
    err := db.DB.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", input.Email).Scan(&user.ID, &user.PasswordHash)
    if err != nil {
        if err == sql.ErrNoRows {
            if c.GetHeader("HX-Request") != "" {
                c.Status(http.StatusUnauthorized)
                c.Writer.Write([]byte(`<div class="alert alert-danger">Email ou mot de passe invalide.</div>`))
                return
            }
            c.Redirect(http.StatusFound, "/login?error=1")
            return
        } else {
            if c.GetHeader("HX-Request") != "" {
                c.Status(http.StatusInternalServerError)
                c.Writer.Write([]byte(`<div class="alert alert-danger">Erreur interne, réessayez plus tard.</div>`))
                return
            }
            c.Redirect(http.StatusFound, "/login?error=1")
            return
        }
    }

    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
        if c.GetHeader("HX-Request") != "" {
            c.Status(http.StatusUnauthorized)
            c.Writer.Write([]byte(`<div class="alert alert-danger">Email ou mot de passe invalide.</div>`))
            return
        }
        c.Redirect(http.StatusFound, "/login?error=1")
        return
    }

    token, err := GenerateJWT(user.ID)
    if err != nil {
        if c.GetHeader("HX-Request") != "" {
            c.Status(http.StatusInternalServerError)
            c.Writer.Write([]byte(`<div class="alert alert-danger">Erreur interne, réessayez plus tard.</div>`))
            return
        }
        c.Redirect(http.StatusFound, "/login?error=1")
        return
    }

    c.SetCookie("token", token, 3600*24, "/", "", false, true)

    // On success, if HTMX: redirect to dashboard
    if c.GetHeader("HX-Request") != "" {
        c.Header("HX-Redirect", "/")
        c.Status(http.StatusOK)
        return
    }

    c.Redirect(http.StatusFound, "/")
}
