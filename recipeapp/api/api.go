package api

import "github.com/gin-gonic/gin"

func LandingPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the Recipe API!",
	})
}
