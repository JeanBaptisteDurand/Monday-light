package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Servir les fichiers statiques
	router.Static("/static", "./frontend")

	// Exemple d'API
	router.GET("/api/projects", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Liste des projets",
		})
	})

	router.Run(":8080") // Port configuré dans le .env
}
