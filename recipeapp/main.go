package main

import (
	"log"
	"recipeapp/api"
	"recipeapp/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const port = ":8080"

var db *gorm.DB

func main() {
	db = initDB()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"} // restrict to local frontend
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Set-Cookie"}

	router := gin.Default()
	router.Use(cors.New(config))
	router.Use(static.Serve("/", static.LocalFile("./ui/recipeapp/out", true))) // Serving the frontend

	apiGroup := router.Group("/api") // API group for all API routes

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
