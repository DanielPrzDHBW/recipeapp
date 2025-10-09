package main

import (
	"recipeapp/api"

	"github.com/gin-gonic/gin"
)

const port = ":8080"

func main() {
	router := gin.Default()
	router.GET("/", api.LandingPage) //Provides the frontend of the

	apiGroup := router.Group("/api") // API group for all API routes
	apiGroup.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	})

	apiGroup.GET("/recipes", api.GetRecipes)    // Get a list of saved Recipes from the database by the users cookies
	apiGroup.GET("/newrecipes", api.NewRecipes) // Get a list of new Recipes from the database by the users cookies

	router.Run(port) // listen and serve on
}
