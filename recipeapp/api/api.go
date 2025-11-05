package api

import (
	"errors"
	"log"
	"recipeapp/client"
	"recipeapp/cookie"
	"recipeapp/database"
	"recipeapp/models"
	"recipeapp/serverError"
	"recipeapp/shoppinglist"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetRecipes(c *gin.Context) {
	id, err := uuid.Parse(cookie.GetCookie(c))
	if err != nil {
		log.Println(err)
	}
	db, err := database.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	recipes, err := database.GetRecipesFromDBByUUID(db, id)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "Internal server error"})
	}
	converter := shoppinglist.IngredientConverter{}
	shoppingList := converter.ConvertMeals(recipes)
	c.JSON(200, gin.H{
		"recipe":        recipes,
		"shopping_list": shoppingList,
	})
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
	converter := shoppinglist.IngredientConverter{}
	shoppingList := converter.ConvertMeals(recipes)
	db, err := database.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	id, err := database.CreateEntry(db, recipes)
	if err != nil {
		log.Fatal(err)
	}
	cookie.SetCookie(c, id.String())
	c.JSON(200, gin.H{
		"recipe":        recipes,
		"shopping_list": shoppingList,
	})
}
