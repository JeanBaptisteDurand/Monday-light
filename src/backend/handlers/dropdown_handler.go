package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// ShowRecap => old logic: "Simple WIP page"
func ShowRecap(c *gin.Context) {
    // If you want "recap_content" subtemplate:
    data := gin.H{
        "Title":           "Recap",
        "Username":        "Unknown", // or fetch from DB
        "ContentTemplate": "recap_content", // define "recap_content" in recap.html
    }
    Render(c, data)
}

func Logout(c *gin.Context) {
    // Invalidate the cookie and redirect to login
    c.SetCookie("token", "", -1, "/", "", false, true)
    c.Redirect(http.StatusFound, "/login")
}
