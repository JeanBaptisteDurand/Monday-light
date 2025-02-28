package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Render is a reusable function that checks if the request is HTMX.
// - If HTMX, we only send the snippet data["ContentTemplate"].
// - Otherwise, we render "base" with data["ContentTemplate"] injected.
func Render(c *gin.Context, data gin.H) {
    snippetAny, ok := data["ContentTemplate"]
    if !ok {
        c.String(http.StatusInternalServerError, "ContentTemplate is missing from data")
        return
    }

    snippetName, ok := snippetAny.(string)
    if !ok {
        c.String(http.StatusInternalServerError, "ContentTemplate must be a string")
        return
    }

    if c.GetHeader("HX-Request") != "" {
        // HTMX request => snippet only
        c.HTML(http.StatusOK, snippetName, data)
    } else {
        // Normal request => base layout
        c.HTML(http.StatusOK, "base", data)
    }
}
