package shoppinglist

import (
	"recipeapp/models"
	"sort"
	"strconv"
	"strings"
)

// IngredientConverter converts and sums ingredients from multiple recipes
type IngredientConverter struct {
	standardizedIngredients map[string]float64 // ingredient -> total amount in standard unit
	ingredientUnits         map[string]string  // ingredient -> standard unit
}

// NewIngredientConverter creates a new IngredientConverter
func NewIngredientConverter() *IngredientConverter {
	return &IngredientConverter{
		standardizedIngredients: make(map[string]float64),
		ingredientUnits:         make(map[string]string),
	}
}

// Conversion factors for mass units (to grams)
var massUnits = map[string]float64{
	"gram":        1,
	"g":           1,
	"kilogram":    1000,
	"kg":          1000,
	"teaspoon":    4,
	"tsp":         4,
	"teaspoons":   4,
	"tablespoon":  15,
	"tbs":         15,
	"tblsp":       15,
	"tbsp":        15,
	"tablespoons": 15,
	"scoop":       30,
	"pinch":       1,
	"dash":        1,
	"lb":          450,
	"lbs":         450,
	"pound":       450,
	"pounds":      450,
	"oz":          30,
	"ounces":      30,
	"ounce":       30,
}

// Conversion factors for volume units (to milliliters)
var volumeUnits = map[string]float64{
	"milliliter": 1,
	"ml":         1,
	"l":          1000,
	"liter":      1000,
	"litre":      1000,
	"litres":     1000,
	"cup":        250,
	"cups":       250,
}

// Standard unit names for non-standard units
var nonStandardUnits = map[string]string{
	"tin":     "Can",
	"can":     "Can",
	"cans":    "Can",
	"tins":    "Can",
	"packet":  "Packet",
	"pack":    "Packet",
	"package": "Packet",
	"bottle":  "Bottle",
	"jar":     "Jar",
}

// ConvertMeals processes multiple meals and returns the standardized shopping list as string array
func (ic *IngredientConverter) ConvertMeals(meals []models.Meal) []string {
	// Reset the maps for new conversion
	ic.standardizedIngredients = make(map[string]float64)
	ic.ingredientUnits = make(map[string]string)

	// Process each meal
	for _, meal := range meals {
		ic.processMeal(meal)
	}

	return ic.formatShoppingList()
}

// formatShoppingList formats the shopping list as "ingredient-measurement-unit" strings
func (ic *IngredientConverter) formatShoppingList() []string {
	var shoppingList []string

	// Get all ingredients and sort them for consistent output
	ingredients := make([]string, 0, len(ic.standardizedIngredients))
	for ingredient := range ic.standardizedIngredients {
		ingredients = append(ingredients, ingredient)
	}
	sort.Strings(ingredients)

	// Format each ingredient
	for _, ingredient := range ingredients {
		amount := ic.standardizedIngredients[ingredient]
		unit := ic.ingredientUnits[ingredient]

		// Format the amount nicely
		amountStr := formatAmount(amount)

		// Create the formatted string without spaces
		item := ingredient + "-" + amountStr + "-" + unit
		shoppingList = append(shoppingList, item)
	}

	return shoppingList
}

// formatAmount formats the amount nicely (removes .00 if whole number)
func formatAmount(amount float64) string {
	if amount == float64(int64(amount)) {
		return strconv.Itoa(int(amount))
	}
	return strconv.FormatFloat(amount, 'f', 2, 64)
}

// processMeal processes a single meal and adds its ingredients to the total
func (ic *IngredientConverter) processMeal(meal models.Meal) {
	// Extract all ingredients and measures from the meal
	ingredients := ic.extractIngredients(meal)
	measures := ic.extractMeasures(meal)

	// Process each ingredient-measure pair
	for i := 0; i < len(ingredients) && i < len(measures); i++ {
		ingredient := strings.TrimSpace(ingredients[i])
		measure := strings.TrimSpace(measures[i])

		// Skip empty ingredients
		if ingredient == "" {
			continue
		}

		ic.processIngredient(ingredient, measure)
	}
}

// extractIngredients extracts all non-empty ingredients from a meal
func (ic *IngredientConverter) extractIngredients(meal models.Meal) []string {
	ingredients := []string{
		meal.StrIngredient1, meal.StrIngredient2, meal.StrIngredient3, meal.StrIngredient4, meal.StrIngredient5,
		meal.StrIngredient6, meal.StrIngredient7, meal.StrIngredient8, meal.StrIngredient9, meal.StrIngredient10,
		meal.StrIngredient11, meal.StrIngredient12, meal.StrIngredient13, meal.StrIngredient14, meal.StrIngredient15,
		meal.StrIngredient16, meal.StrIngredient17, meal.StrIngredient18, meal.StrIngredient19, meal.StrIngredient20,
	}

	// Filter out empty ingredients
	var nonEmptyIngredients []string
	for _, ingredient := range ingredients {
		if strings.TrimSpace(ingredient) != "" {
			nonEmptyIngredients = append(nonEmptyIngredients, ingredient)
		}
	}

	return nonEmptyIngredients
}

// extractMeasures extracts all measures corresponding to non-empty ingredients
func (ic *IngredientConverter) extractMeasures(meal models.Meal) []string {
	measures := []string{
		meal.StrMeasure1, meal.StrMeasure2, meal.StrMeasure3, meal.StrMeasure4, meal.StrMeasure5,
		meal.StrMeasure6, meal.StrMeasure7, meal.StrMeasure8, meal.StrMeasure9, meal.StrMeasure10,
		meal.StrMeasure11, meal.StrMeasure12, meal.StrMeasure13, meal.StrMeasure14, meal.StrMeasure15,
		meal.StrMeasure16, meal.StrMeasure17, meal.StrMeasure18, meal.StrMeasure19, meal.StrMeasure20,
	}

	// Filter out measures for empty ingredients (we'll use the same indices as ingredients)
	var relevantMeasures []string
	ingredients := ic.extractIngredients(meal)

	for i := 0; i < len(ingredients) && i < len(measures); i++ {
		relevantMeasures = append(relevantMeasures, measures[i])
	}

	return relevantMeasures
}

// processIngredient processes a single ingredient and adds it to the total
func (ic *IngredientConverter) processIngredient(ingredient, measure string) {
	// Parse the measure using the shoppinglist utility
	amount, unit, err := SplitLeadingNumberDecimal(measure)
	if err != nil {
		// If we can't parse a number, assume amount = 1 and the whole string is the unit
		amount = 1
		unit = measure
	} else if unit == "" && amount > 0 {
		// If we have amount but no unit, assume it's a count (like "2 eggs")
		unit = "count"
	}

	// Normalize the unit for comparison
	normalizedUnit := strings.ToLower(strings.TrimSpace(unit))

	// Convert to standard unit
	standardizedAmount, standardUnit := ic.convertToStandardUnit(amount, normalizedUnit)

	// Update the total for this ingredient
	ic.standardizedIngredients[ingredient] += standardizedAmount
	ic.ingredientUnits[ingredient] = standardUnit
}

// convertToStandardUnit converts an amount and unit to the standard unit
func (ic *IngredientConverter) convertToStandardUnit(amount float64, unit string) (float64, string) {
	// Check if it's a mass unit
	if factor, exists := massUnits[unit]; exists {
		return amount * factor, "g"
	}

	// Check if it's a volume unit
	if factor, exists := volumeUnits[unit]; exists {
		return amount * factor, "ml"
	}

	// Check if it's a non-standard unit that we have a standard name for
	if standardUnit, exists := nonStandardUnits[unit]; exists {
		return amount, standardUnit
	}

	// If unit not found in our conversion tables, return as-is
	return amount, unit
}
