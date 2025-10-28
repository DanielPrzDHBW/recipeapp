package cookie

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setCookieHandler(c *gin.Context, uuid uuid.UUID) {
	c.SetCookie("user_recipes", uuid.String(), 604800, "", "", false, true)
}

func getCookieHandler(c *gin.Context) (uuid.UUID, error) {
	cookie, err := c.Cookie("user_recipes")
	if err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(cookie)
	if err != nil {
		return id, err
	}
	return id, nil
}
