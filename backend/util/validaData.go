package util

import (
	"reflect"
)

// ValidateStruct checks if all fields in a struct are valid (non-zero values).
// Parameters:
// - s: the struct to validate, can be a pointer to a struct or a struct itself
// Returns: true if all fields are valid, false otherwise
func ValidateStruct(s interface{}) bool {
	val := reflect.ValueOf(s) // Get the reflect.Value of the input
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference the pointer if the input is a pointer
	}

	if val.Kind() != reflect.Struct {
		return false // Return false if the input is not a struct
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i) // Get the value of each field

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				return false // Return false if a string field is empty
			}
		case reflect.Int:
			if field.Int() == 0 {
				return false // Return false if an int field is zero
			}
		case reflect.Uint, reflect.Uint32:
			if field.Uint() == 0 {
				return false // Return false if an unsigned int field is zero
			}
		}
	}

	return true // Return true if all fields are valid
}
