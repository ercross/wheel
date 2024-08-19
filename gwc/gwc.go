package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

var errInvalidInput = fmt.Errorf("input contains non-utf8 encoded character")

type output struct {
	WordCount      int
	LineCount      int
	CharacterCount int
	ByteCount      int
}

func executeBasedOnFileSize() {
	// todo: for large filesize, you can't keep all the bytes in memory at once
	// so have a limit, say if file is greater than 250mb, load bytes in chunks
}

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

func countLines(input []byte) int {
	return 0
}

func countCharacters(input []byte) int {
	return 0
}

func countBytes(input []byte) int {
	return 0
}
