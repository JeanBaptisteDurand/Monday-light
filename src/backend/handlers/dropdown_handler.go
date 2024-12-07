package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func ShowRecap(c *gin.Context) {
    // Simple WIP page
    Render(c, "recap.html", gin.H{})
}

func Logout(c *gin.Context) {
    // Invalidate the cookie and redirect to login
    c.SetCookie("token", "", -1, "/", "", false, true)
    c.Redirect(http.StatusFound, "/login")
}
