package main

import (
	"fmt"
	"regexp"
)

// https://gosamples.dev/remove-non-alphanumeric/

var nonAlphaNumricRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func clearString(str string) string {
	return nonAlphaNumricRegex.ReplaceAllString(str, "")
}

func main() {
	//str := "Test@%String#321gosamples.dev ـا ą ٦"
	fmt.Println(nonAlphaNumricRegex.FindString("124@341$go"))
}
