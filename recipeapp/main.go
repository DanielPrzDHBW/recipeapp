package main

import (
	"log"
	"recipeapp/api"
	"recipeapp/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const port = ":8080"

var db *gorm.DB

func main() {
	db = initDB()

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

func initDB() *gorm.DB {
	dbNew, err := database.ConnectToSQLite()
	if err != nil {
		log.Fatal(err)
	}
	err = dbNew.AutoMigrate(&database.RecipesEntry{})
	if err != nil {
		log.Fatal(err)
	}
	database.SetDB(dbNew)
	return dbNew
}
