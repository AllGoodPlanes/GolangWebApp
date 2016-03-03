package main

import (
	"unicode"
)

var passwordMin int
var passwordMax int

func checkLength(value string, min, max int) bool {

	pwdLen := len(value)
	if min > 0 && pwdLen < min {
		return false
	}
	if max > 0 && pwdLen > max {
		return false
	}
	return true
}

var mustHave = []func(rune) bool{
	unicode.IsUpper,
	unicode.IsLower,
	//unicode.IsPunct,
	unicode.IsDigit,
}

func passwordCk(p string) bool {
	for _, testRune := range mustHave {
		found := false
		for _, r := range p {
			if testRune(r) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
