package sliceutil

import (
	"reflect"
	"sort"
)

// Compare will check if two slices are equal
// even if they aren't in the same order
// Inspired by github.com/stephanbaker white board sudo code
func Compare(s1, s2 interface{}) bool {
	if s1 == nil || s2 == nil {
		return false
	}

	// Convert slices to correct type
	slice1 := convertSliceToInterface(s1)
	slice2 := convertSliceToInterface(s2)
	if slice1 == nil || slice2 == nil {
		return false
	}

	if len(slice1) != len(slice2) {
		return false
	}

	// setup maps to store values and count of slices
	m1 := make(map[interface{}]int)
	m2 := make(map[interface{}]int)

	for i := 0; i < len(slice1); i++ {
		// Add each value to map and increment for each found
		m1[slice1[i]]++
		m2[slice2[i]]++
	}

	for key := range m1 {
		if m1[key] != m2[key] {
			return false
		}
	}

	return true
}

// OrderedCompare will check if two slices are equal, taking order into consideration.
func OrderedCompare(s1, s2 interface{}) bool {
	//If both are nil, they are equal
	if s1 == nil && s2 == nil {
		return true
	}

	//If only one is nil, they are not equal (!= represents XOR)
	if (s1 == nil) != (s2 == nil) {
		return false
	}

	// Convert slices to correct type
	slice1 := convertSliceToInterface(s1)
	slice2 := convertSliceToInterface(s2)

	//If both are nil, they are equal
	if slice1 == nil || slice2 == nil {
		return false
	}

	//If the lengths are different, the slices are not equal
	if len(slice1) != len(slice2) {
		return false
	}

	//Loop through and compare the slices at each index
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	//If nothing has failed up to this point, the slices are equal
	return true
}

// Contains checks if a slice contains an element
func Contains(s interface{}, e interface{}) bool {
	slice := convertSliceToInterface(s)

	for _, a := range slice {
		if a == e {
			return true
		}
	}
	return false
}

// ContainsString checks if a slice contains a string
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// convertSliceToInterface takes a slice passed in as an interface{}
// then converts the slice to a slice of interfaces
func convertSliceToInterface(s interface{}) (slice []interface{}) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Slice {
		return nil
	}

	length := v.Len()
	slice = make([]interface{}, length)
	for i := 0; i < length; i++ {
		slice[i] = v.Index(i).Interface()
	}

	return slice
}

// Reverse reverses slices of any type
func Reverse(s []interface{}) []interface{} {
	a := make([]interface{}, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

// FastContains offers a faster implementation of contains
func FastContains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
