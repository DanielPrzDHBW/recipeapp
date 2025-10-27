package cookie

import (
	"net/http"

	"github.com/google/uuid"
)

// SetCookie handles requests to /set and sets a cookie in the client's browser
func SetCookie(w http.ResponseWriter, uuid uuid.UUID) {
	http.SetCookie(w, &http.Cookie{
		Name:  "recipes",
		Value: uuid.String(),
	})
}

// GetCookie handles requests to /get and tries to read the previously set cookie
func GetCookie(r *http.Request) (uuid.UUID, error) {
	c, err := r.Cookie("recipes")
	if err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(c.Value)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
