package api

import (
	"errors"
	"log"
	"recipeapp/client"
	"recipeapp/cookie"
	"recipeapp/database"
	"recipeapp/models"
	"recipeapp/serverError"

	"github.com/gin-gonic/gin"
  "github.com/google/uuid"
)

var recipes = []models.Meal{}        // Placeholder for a database
var shoppingList = []string{"test0"} // Placeholder for a database

// Placeholder to serve the landing page frontend
func LandingPage(c *gin.Context) {
}

// Placeholder for future implementation of reading previous generated recipes from a database
func GetRecipes(c *gin.Context) {
	c.JSON(200, gin.H{
		"recipe":        recipes,
		"shopping_list": shoppingList,
	})
	id, err := uuid.Parse(cookie.GetCookie(c))
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	database.GetRecipesFromDBByUUID(db, id)
	// TODO: implement further usage of the recipes
}

// Generates 7 new recipes and returning the as a JSON array
func NewRecipes(c *gin.Context) {
	recipes := []models.Meal{}
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
	shoppingList = append(shoppingList, "Test Item 1", "Test Item 2", "Test Item 3") // Placeholder for shopping list generation
	c.JSON(200, gin.H{
		"recipe":        recipes,
		"shopping_list": shoppingList,
	})
	db, err := database.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	id, err := database.CreateEntry(db, recipes)
	if err != nil {
		log.Fatal(err)
	}
	cookie.SetCookie(c, id.String())
}
