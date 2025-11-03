package shoppinglist

import (
	"strings"
)

type Ingredient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type ConvertedIngredient struct {
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Unit        string  `json:"unit"`
	BaseAmount  float64 `json:"base_amount"`
	BaseUnit    string  `json:"base_unit"`
	IsConverted bool    `json:"is_converted"`
}

type SummedIngredient struct {
	Name        string  `json:"name"`
	TotalAmount float64 `json:"total_amount"`
	Unit        string  `json:"unit"`
}

// Conversion factors for different unit categories
var massConversions = map[string]float64{
	"g": 1, "gram": 1, "grams": 1,
	"kilogram": 1000, "kg": 1000, "kilograms": 1000,
	"teaspoon": 4, "tsp": 4, "teaspoons": 4,
	"tablespoon": 15, "tbs": 15, "tblsp": 15, "tbsp": 15, "tablespoons": 15,
	"scoop": 30, "scoops": 30,
	"pinch": 1, "dash": 1, "pinches": 1, "dashes": 1,
	"lb": 450, "lbs": 450, "pound": 450, "pounds": 450,
	"oz": 30, "ounce": 30, "ounces": 30,
}

var volumeConversions = map[string]float64{
	"ml": 1, "milliliter": 1, "milliliters": 1,
	"l": 1000, "liter": 1000, "litre": 1000, "liters": 1000, "litres": 1000,
	"cup": 250, "cups": 250,
}

var countableUnits = map[string]bool{
	"can": true, "cans": true, "tin": true, "tins": true,
	"packet": true, "pack": true, "package": true, "packages": true,
	"bottle": true, "bottles": true, "jar": true, "jars": true,
	"piece": true, "pieces": true, "item": true, "items": true,
}

// NormalizeUnit standardizes unit names and handles plural forms
func NormalizeUnit(unit string) string {
	unit = strings.ToLower(strings.TrimSpace(unit))

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
	case "items":
		return "item"
	case "pieces":
		return "piece"
	default:
		return unit
	}
}

// ParseIngredientString parses a raw ingredient string using your existing SplitLeadingNumberDecimal function
func ParseIngredientString(ingredientStr string) (Ingredient, error) {
	amount, nameAndUnit, err := SplitLeadingNumberDecimal(ingredientStr)
	if err != nil {
		return Ingredient{}, err
	}

	// Extract unit and name from the remaining string
	unit, name := extractUnitAndName(nameAndUnit)

	return Ingredient{
		Name:   strings.TrimSpace(name),
		Amount: amount,
		Unit:   unit,
	}, nil
}

// extractUnitAndName separates unit from the ingredient name
func extractUnitAndName(nameAndUnit string) (string, string) {
	words := strings.Fields(nameAndUnit)
	if len(words) == 0 {
		return "", ""
	}

	// Check if the first word is a known unit
	if len(words) > 0 {
		possibleUnit := NormalizeUnit(words[0])
		if IsConvertibleUnit(possibleUnit) {
			// First word is a unit, return it and the rest as name
			return possibleUnit, strings.Join(words[1:], " ")
		}
	}

	// No recognized unit found, return empty unit and entire string as name
	return "", nameAndUnit
}

// IsConvertibleUnit checks if a unit can be converted to a base unit
func IsConvertibleUnit(unit string) bool {
	normalizedUnit := NormalizeUnit(unit)

	if countableUnits[normalizedUnit] {
		return true
	}
	if _, exists := massConversions[normalizedUnit]; exists {
		return true
	}
	if _, exists := volumeConversions[normalizedUnit]; exists {
		return true
	}
	return false
}

// ConvertToBaseUnit converts an ingredient to its base unit for summation
func ConvertToBaseUnit(ingredient Ingredient) ConvertedIngredient {
	normalizedUnit := NormalizeUnit(ingredient.Unit)

	converted := ConvertedIngredient{
		Name:        ingredient.Name,
		Amount:      ingredient.Amount,
		Unit:        ingredient.Unit,
		BaseAmount:  ingredient.Amount,
		BaseUnit:    normalizedUnit,
		IsConverted: false,
	}

	// Check if it's a countable unit (no conversion needed, just normalization)
	if countableUnits[normalizedUnit] {
		converted.BaseUnit = normalizedUnit
		converted.IsConverted = true
		return converted
	}

	// Check mass conversions
	if factor, exists := massConversions[normalizedUnit]; exists {
		converted.BaseAmount = ingredient.Amount * factor
		converted.BaseUnit = "gram"
		converted.IsConverted = true
		return converted
	}

	// Check volume conversions
	if factor, exists := volumeConversions[normalizedUnit]; exists {
		converted.BaseAmount = ingredient.Amount * factor
		converted.BaseUnit = "milliliter"
		converted.IsConverted = true
		return converted
	}

	// Unknown unit - cannot convert, return as is
	return converted
}

// ConvertMultipleIngredients converts a slice of ingredients to their base units
func ConvertMultipleIngredients(ingredients []Ingredient) []ConvertedIngredient {
	converted := make([]ConvertedIngredient, len(ingredients))
	for i, ing := range ingredients {
		converted[i] = ConvertToBaseUnit(ing)
	}
	return converted
}

// ParseAndConvertMultiple parses raw ingredient strings and converts them to base units
func ParseAndConvertMultiple(ingredientStrs []string) ([]ConvertedIngredient, error) {
	var converted []ConvertedIngredient

	for _, str := range ingredientStrs {
		ingredient, err := ParseIngredientString(str)
		if err != nil {
			return nil, err
		}
		converted = append(converted, ConvertToBaseUnit(ingredient))
	}

	return converted, nil
}

// GetConversionFactor returns the conversion factor for a given unit
func GetConversionFactor(unit string) (float64, string, bool) {
	normalizedUnit := NormalizeUnit(unit)

	// Check countable units (factor is 1, base unit is the normalized unit)
	if countableUnits[normalizedUnit] {
		return 1, normalizedUnit, true
	}

	// Check mass conversions
	if factor, exists := massConversions[normalizedUnit]; exists {
		return factor, "gram", true
	}

	// Check volume conversions
	if factor, exists := volumeConversions[normalizedUnit]; exists {
		return factor, "milliliter", true
	}

	// Unknown unit
	return 1, normalizedUnit, false
}

// SumIngredients sums up all ingredients by name and base unit
func SumIngredients(ingredients []Ingredient) map[string]SummedIngredient {
	// First convert all ingredients to base units
	convertedIngredients := ConvertMultipleIngredients(ingredients)

	// Create a map to store summed ingredients (key: "name|baseUnit")
	summedMap := make(map[string]SummedIngredient)

	for _, convIng := range convertedIngredients {
		// Create a composite key from name and base unit
		compositeKey := strings.ToLower(convIng.Name) + "|" + convIng.BaseUnit

		if existing, exists := summedMap[compositeKey]; exists {
			// Add to existing ingredient
			existing.TotalAmount += convIng.BaseAmount
			summedMap[compositeKey] = existing
		} else {
			// Create new entry
			summedMap[compositeKey] = SummedIngredient{
				Name:        convIng.Name,
				TotalAmount: convIng.BaseAmount,
				Unit:        convIng.BaseUnit,
			}
		}
	}

	return summedMap
}

// SumAndConvertIngredients sums ingredients and converts them to human-readable units
func SumAndConvertIngredients(ingredients []Ingredient) map[string]SummedIngredient {
	summedMap := SumIngredients(ingredients)

	// Convert to human-readable units
	for key, summedIng := range summedMap {
		humanAmount, humanUnit := toHumanReadableUnit(summedIng.TotalAmount, summedIng.Unit)
		summedIng.TotalAmount = humanAmount
		summedIng.Unit = humanUnit
		summedMap[key] = summedIng
	}

	return summedMap
}

// SumRawIngredientStrings parses and sums raw ingredient strings in one operation
func SumRawIngredientStrings(ingredientStrs []string) (map[string]SummedIngredient, error) {
	var ingredients []Ingredient

	for _, str := range ingredientStrs {
		ingredient, err := ParseIngredientString(str)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}

	return SumAndConvertIngredients(ingredients), nil
}

// toHumanReadableUnit converts base units to more human-readable formats
func toHumanReadableUnit(amount float64, baseUnit string) (float64, string) {
	switch baseUnit {
	case "gram":
		if amount >= 1000 {
			return amount / 1000, "kg"
		}
		return amount, "g"
	case "milliliter":
		if amount >= 1000 {
			return amount / 1000, "L"
		}
		return amount, "ml"
	default:
		// For countable and unknown units, return as is
		return amount, baseUnit
	}
}
