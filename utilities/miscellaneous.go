package utilities

import (
	"regexp"
	"strings"
)

var reAlphaNum = regexp.MustCompile(`[\w\d]`) // matches all alphanumeric characters

// TrimString trims a string (usually the novel name) to only lowercase alphanumeric characters for easier mapping
// (since novel websites sometimes move around capitalization and punctuation).
func TrimString(in string) string {
	in = reAlphaNum.ReplaceAllString(in, "")
	in = strings.ToLower(in)

	return in
}

// Len2DSlice gets the total number of elements within a 2D slice.
func Len2DSlice(in [][]interface{}) int {
	var length int
	for _, arr := range in {
		for range arr {
			length++
		}
	}

	return length
}
