package cookie

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, token string) {
	// Set cookie
	maxAge := int((30 * 24 * time.Hour).Seconds())

	// host-only cookie (empty domain), not Secure for local dev, HttpOnly true
	// always set the cookie (don't rely on presence in request)
	c.SetCookie("recipe_cookie", token, maxAge, "/", "", false, true)

	// debug: print Set-Cookie header(s) written to the response
	fmt.Printf("Set-Cookie header(s) written: %v\n", c.Writer.Header().Values("Set-Cookie"))

}

func GetCookie(c *gin.Context) string {
	cookie, err := c.Cookie("recipe_cookie")
	if err != nil {
		fmt.Println("Error retrieving cookie:", err)
		return ""
	}
	return cookie
}
