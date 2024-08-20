package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

var errInvalidInput = fmt.Errorf("input contains non-utf8 encoded character")

// countWords counts the number of words in a slice of bytes,
// where a word is defined as sequences of characters delimited by whitespace.
//
// countWords return error if it encounters a character that is not UTF8 encoded
func countWords(input []byte) (int, error) {

	if len(input) == 0 {
		return 0, nil
	}

	// make a copy of input to avoid permanently modifying input
	copied := make([]byte, len(input))
	copy(copied, input)

	// inWord indicates the current iteration is still within a word
	inWord := false
	count := 0

	// Iterate over each rune in the byte slice
	for len(copied) > 0 {

		r, runeSize := utf8.DecodeRune(copied)

		// check for invalid rune
		if r == utf8.RuneError && runeSize == 1 {
			return 0, errInvalidInput
		}

		copied = copied[runeSize:]
		if unicode.IsSpace(r) {
			inWord = false
		} else if !inWord {
			// encountered a non-whitespace character and loop is not in a word
			inWord = true
			count++
		}
	}

	return count, nil
}

// countLines basically counts the number of unix newline character found in input.
// This implies that if input contains no other characters
// except the unix newline character, countLines returns a non-zero result
func countLines(input []byte) int {
	if len(input) == 0 {
		return 0
	}

	// make a copy of input to avoid permanently modifying input
	copied := make([]byte, len(input))
	copy(copied, input)

	count := 0
	var newlineCharacter byte = '\n'

	for _, c := range copied {
		if c == newlineCharacter {
			count++
		}
	}

	return count
}

// countCharacters counts the number of UTF-8 encoded characters
// (including but not limited to whitespaces, newline, tab, etc.) in input.
//
// countCharacters return error if it encounters a character that is not UTF8 encoded
func countCharacters(input []byte) (int, error) {

	if len(input) == 0 {
		return 0, nil
	}

	// make a copy of input to avoid permanently modifying input
	copied := make([]byte, len(input))
	copy(copied, input)

	count := 0

	for len(copied) > 0 {
		r, runeSize := utf8.DecodeRune(copied)
		// check for invalid rune
		if r == utf8.RuneError && runeSize == 1 {
			return 0, errInvalidInput
		}
		copied = copied[runeSize:]
		count++
	}

	return count, nil
}

func countBytes(input []byte) int {
	return len(input)
}
