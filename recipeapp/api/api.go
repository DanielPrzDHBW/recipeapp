package api

import (
	"errors"
	"recipeapp/client"
	"recipeapp/models"
	"recipeapp/serverError"

	"github.com/gin-gonic/gin"
)

var recipes = []models.Meal{} // Placeholder for a database

// Placeholder to serve the landing page frontend
func LandingPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the Recipe API!",
	})
}

// Placeholder for future implementation of reading previous generated recipes from a database
func GetRecipes(c *gin.Context) {
	c.JSON(200, gin.H{
		"recipe": recipes,
	})
}

// Generates 7 new recipes and returning the as a JSON array
// TODO: Implement saving to the database
func NewRecipes(c *gin.Context) {
	recipes = []models.Meal{}
	for i := 0; i < 7; i++ {
		resp, err := client.NewRecipe()
		if err != nil {
			if errors.Is(err, serverError.BadInternalApiCall) {
				c.JSON(503, gin.H{
					"error": "Failed to fetch new recipe",
				})
				return
			}
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}
		// Filtering out unwanted categories
		if resp.Meals[0].StrCategory != "Dessert" && resp.Meals[0].StrCategory != "Side" && resp.Meals[0].StrCategory != "Miscellaneous" && resp.Meals[0].StrCategory != "Starter" {
			recipes = append(recipes, resp.Meals...)
		} else {
			i--
		}
	}
	c.JSON(200, gin.H{
		"recipe": recipes,
	})
}
