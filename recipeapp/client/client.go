package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recipeapp/models"
	"recipeapp/serverError"
)

var url = "https://www.themealdb.com/api/json/v1/1/random.php"

type Response struct {
	Meals []models.Meal `json:"meals"`
}

// Function to fetch a single random recipe from the external API
func NewRecipe() (*Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, serverError.BadInternalApiCall
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil, serverError.BadInternalApiCall
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		fmt.Printf("Non-200 response: %v\n", resp.Status)
		return nil, serverError.BadInternalApiCall
	}

	r := Response{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, serverError.BadInternalApiCall
	}
	return &r, nil
}
