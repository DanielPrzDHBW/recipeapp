package main

import (
	"fmt"
	"os"

	"github.com/DanielPrzDHBW/recipeapp/api"

	"github.com/gin-gonic/gin"
)

const port = ":8080"

func main() {
	router := gin.Default()
	router.GET("/", api.LandingPage)
	fmt.Printf("Hello, %s!\n", os.Getenv("USER"))

	router.Run(port) // listen and serve on
}
