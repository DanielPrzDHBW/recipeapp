package cookie

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, token string) {
	// Set cookie
	token_lifespan := 30
	c.SetCookie("recipe_cookie", token, token_lifespan*60*60*24, "/", "localhost", false, true)

	cookie, err := c.Cookie("recipe_cookie")
	if err != nil {
		cookie = "NotSet"
	}

	fmt.Printf("Cookie value: %s \n", cookie)
}

func GetCookie(c *gin.Context) string {
	cookie, err := c.Cookie("recipe_cookie")
	if err != nil {
		fmt.Println("Error retrieving cookie:", err)
		return ""
	}
	return cookie
}
