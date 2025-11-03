package ingredient

import (
	"fmt"
	"strings"
)

type Ingredient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type ConvertedIngredient struct {
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Unit       string  `json:"unit"`
	BaseUnit   string  `json:"base_unit"`
	BaseAmount float64 `json:"base_amount"`
}

// Conversion factors for different unit categories
var massConversions = map[string]float64{
	"g":           1,
	"gram":        1,
	"grams":       1,
	"kilogram":    1000,
	"kg":          1000,
	"kilograms":   1000,
	"teaspoon":    4,
	"tsp":         4,
	"teaspoons":   4,
	"tablespoon":  15,
	"tbs":         15,
	"tblsp":       15,
	"tbsp":        15,
	"tablespoons": 15,
	"scoop":       30,
	"scoops":      30,
	"pinch":       1,
	"dash":        1,
	"pinches":     1,
	"dashes":      1,
	"lb":          450,
	"lbs":         450,
	"pound":       450,
	"pounds":      450,
	"oz":          30,
	"ounce":       30,
	"ounces":      30,
}

var volumeConversions = map[string]float64{
	"ml":          1,
	"milliliter":  1,
	"milliliters": 1,
	"l":           1000,
	"liter":       1000,
	"litre":       1000,
	"liters":      1000,
	"litres":      1000,
	"cup":         250,
	"cups":        250,
}

// Units that don't need conversion (countable items)
var countableUnits = map[string]bool{
	"can": true, "cans": true, "tin": true, "tins": true,
	"packet": true, "pack": true, "package": true, "packages": true,
	"bottle": true, "bottles": true, "jar": true, "jars": true,
	"piece": true, "pieces": true, "item": true, "items": true,
}

func normalizeUnit(unit string) string {
	unit = strings.ToLower(strings.TrimSpace(unit))

	// Handle plural forms and common variations
	switch unit {
	case "tbs", "tblsp", "tbsp":
		return "tablespoon"
	case "tsp":
		return "teaspoon"
	case "kg", "kilograms":
		return "kilogram"
	case "grams":
		return "gram"
	case "lbs", "pounds":
		return "pound"
	case "oz", "ounces":
		return "ounce"
	case "l", "liters", "litres":
		return "liter"
	case "ml", "milliliters":
		return "milliliter"
	case "pack", "packages":
		return "packet"
	case "tins":
		return "tin"
	case "jars":
		return "jar"
	case "bottles":
		return "bottle"
	default:
		return unit
	}
}

func convertToBaseUnit(ingredient Ingredient) (ConvertedIngredient, bool) {
	normalizedUnit := normalizeUnit(ingredient.Unit)

	converted := ConvertedIngredient{
		Name:   ingredient.Name,
		Amount: ingredient.Amount,
		Unit:   ingredient.Unit,
	}

	// Check if it's a countable unit (no conversion needed)
	if countableUnits[normalizedUnit] {
		converted.BaseUnit = normalizedUnit
		converted.BaseAmount = ingredient.Amount
		return converted, true
	}

	// Check mass conversions
	if factor, exists := massConversions[normalizedUnit]; exists {
		converted.BaseUnit = "gram"
		converted.BaseAmount = ingredient.Amount * factor
		return converted, true
	}

	// Check volume conversions
	if factor, exists := volumeConversions[normalizedUnit]; exists {
		converted.BaseUnit = "milliliter"
		converted.BaseAmount = ingredient.Amount * factor
		return converted, true
	}

	// Unknown unit - return as is
	converted.BaseUnit = normalizedUnit
	converted.BaseAmount = ingredient.Amount
	return converted, false
}

type ShoppingListItem struct {
	Name          string  `json:"name"`
	TotalAmount   float64 `json:"total_amount"`
	Unit          string  `json:"unit"`
	OriginalUnits []struct {
		Amount float64 `json:"amount"`
		Unit   string  `json:"unit"`
	} `json:"original_units,omitempty"`
}

func SumIngredients(ingredients []Ingredient) []ShoppingListItem {
	// First, convert all ingredients to base units
	convertedIngredients := make([]ConvertedIngredient, 0, len(ingredients))

	for _, ing := range ingredients {
		converted, _ := convertToBaseUnit(ing)
		convertedIngredients = append(convertedIngredients, converted)
	}

	// Group by ingredient name and base unit using a composite key
	ingredientMap := make(map[string]*ShoppingListItem)

	for _, convIng := range convertedIngredients {
		// Create a composite key from name and base unit
		compositeKey := fmt.Sprintf("%s|%s", strings.ToLower(convIng.Name), convIng.BaseUnit)

		if item, exists := ingredientMap[compositeKey]; exists {
			item.TotalAmount += convIng.BaseAmount
			item.OriginalUnits = append(item.OriginalUnits, struct {
				Amount float64 `json:"amount"`
				Unit   string  `json:"unit"`
			}{
				Amount: convIng.Amount,
				Unit:   convIng.Unit,
			})
		} else {
			ingredientMap[compositeKey] = &ShoppingListItem{
				Name:        convIng.Name,
				TotalAmount: convIng.BaseAmount,
				Unit:        convIng.BaseUnit,
				OriginalUnits: []struct {
					Amount float64 `json:"amount"`
					Unit   string  `json:"unit"`
				}{
					{
						Amount: convIng.Amount,
						Unit:   convIng.Unit,
					},
				},
			}
		}
	}

	// Convert map to slice
	var shoppingList []ShoppingListItem
	for _, item := range ingredientMap {
		shoppingList = append(shoppingList, *item)
	}

	return shoppingList
}
