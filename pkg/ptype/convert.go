package ptype

import (
	"strconv"
	"time"
)

// PTime ...
func PTime(t time.Time) *time.Time {
	return &t
}

// TimeP ...
func TimeP(t *time.Time) time.Time {
	return *t
}

// PFloat64 converts a float64 pointer to float64
func PFloat64(num *float64) float64 {
	if num == nil {
		return 0
	}
	return *num
}

// Float64P converts a float to an float pointer
func Float64P(f float64) *float64 {
	return &f
}

// PFloat32 converts a float32 pointer to float32
func PFloat32(num *float32) float32 {
	if num == nil {
		return 0
	}
	return *num
}

// Float32P converts a float to an float pointer
func Float32P(f float32) *float32 {
	return &f
}

// IntP converts an int to int pointer
func IntP(i int) *int {
	return &i
}

// PInt converts a pointer int to int
func PInt(i *int) int {
	return *i
}

// StringP converts a string to a string pointer
func StringP(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

// PString converts a string pointer to a string
func PString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

// PBool ...
func PBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// BoolP converts bool to bool pointer
func BoolP(b bool) *bool {
	return &b
}

// FloatToString Converts float to a string
func FloatToString(num float64) string {
	return strconv.FormatFloat(num, 'G', -1, 64)
}
