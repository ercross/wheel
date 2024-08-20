package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	programName                           = "gwc"
	printNumberOfBytes      flagCharacter = 'c'
	printNumberOfWords      flagCharacter = 'w'
	printNumberOfLines      flagCharacter = 'l'
	printNumberOfCharacters flagCharacter = 'm'
)

type flagCharacter rune

type result struct {
	numberOfBytes      int
	numberOfWords      int
	numberOfLines      int
	numberOfCharacters int
}

// command is of the form `gwc [OPTIONS] filepath...
type command struct {
	options   outputOptions
	filePaths []string
}

type outputOptions struct {
	printNumberOfBytes      bool
	printNumberOfWords      bool
	printNumberOfLines      bool
	printNumberOfCharacters bool
}

func (c command) process() (r result, err error) {
	var (
		characterCount, lineCount, wordCount, byteCount int
		raw                                             []byte
	)

	for _, file := range c.filePaths {
		var loopCC, loopLC, loopWC, loopBC int
		raw, err = os.ReadFile(file)
		if err != nil {
			return r, fmt.Errorf("error reading file %s: %w", file, err)
		}

		if !c.options.printNumberOfCharacters && !c.options.printNumberOfBytes &&
			!c.options.printNumberOfWords && !c.options.printNumberOfLines {
			loopCC, err = countCharacters(raw)
			if err != nil {
				return r, err
			}

			loopWC, err = countWords(raw)
			if err != nil {
				return r, err
			}

			loopLC = countLines(raw)
			loopBC = countBytes(raw)

		} else {
			if c.options.printNumberOfCharacters {
				loopCC, err = countCharacters(raw)
				if err != nil {
					return r, err
				}
			}

			if c.options.printNumberOfWords {
				loopWC, err = countWords(raw)
				if err != nil {
					return r, err
				}
			}

			if c.options.printNumberOfLines {
				loopLC = countLines(raw)
			}

			if c.options.printNumberOfBytes {
				loopBC = countBytes(raw)
			}
		}

		wordCount = wordCount + loopWC
		byteCount = byteCount + loopBC
		lineCount = lineCount + loopLC
		characterCount = characterCount + loopCC
	}

	return result{
		numberOfBytes:      byteCount,
		numberOfWords:      wordCount,
		numberOfLines:      lineCount,
		numberOfCharacters: characterCount,
	}, nil
}

func (r result) format(o outputOptions) string {

	if !o.printNumberOfWords && !o.printNumberOfCharacters &&
		!o.printNumberOfLines && !o.printNumberOfBytes {
		return fmt.Sprintf("words: %d\nlines: %d\ncharacters: %d\nbytes: %d",
			r.numberOfWords, r.numberOfLines, r.numberOfCharacters, r.numberOfBytes)
	}

	var builder strings.Builder
	if o.printNumberOfWords {
		builder.WriteString(fmt.Sprintf("words: %d\n", r.numberOfWords))
	}
	if o.printNumberOfCharacters {
		builder.WriteString(fmt.Sprintf("characters: %d\n", r.numberOfCharacters))
	}

	if o.printNumberOfLines {
		builder.WriteString(fmt.Sprintf("lines: %d\n", r.numberOfLines))
	}

	if o.printNumberOfBytes {
		builder.WriteString(fmt.Sprintf("bytes: %d\n", r.numberOfBytes))
	}

	return builder.String()
}

func parseArgs(args []string) (command, error) {

	options, err := parseFlagManually(args)
	if err != nil {
		return command{}, err
	}

	filePaths, err := extractFilePaths(args)
	if err != nil {
		return command{}, err
	}

	return command{
		options:   options,
		filePaths: filePaths,
	}, nil
}

// parseFlagManually does not use the flag package because
// flags may be passed as a combined string e.g., -mlc, -cl,
// or as a standalone -c -l
func parseFlagManually(args []string) (outputOptions, error) {
	var (
		flagPrintNumberOfBytes      bool
		flagPrintNumberOfWords      bool
		flagPrintNumberOfLines      bool
		flagPrintNumberOfCharacters bool
	)

	// parse flag as combination i.e., flags are written together
	for _, arg := range args {
		if isFlag(arg) {
			bytes := []byte(strings.ReplaceAll(arg, "-", ""))
			for len(bytes) > 0 {

				r, runeSize := utf8.DecodeRune(bytes)
				bytes = bytes[runeSize:]

				if unicode.IsSpace(r) {
					continue
				}
				switch r {
				case rune(printNumberOfBytes):
					flagPrintNumberOfBytes = true
				case rune(printNumberOfWords):
					flagPrintNumberOfWords = true
				case rune(printNumberOfLines):
					flagPrintNumberOfLines = true
				case rune(printNumberOfCharacters):
					flagPrintNumberOfCharacters = true
				default:
					return outputOptions{}, fmt.Errorf("unknown [OPTION] %s", string(r))
				}
			}
		} else {

			// don't process any flag that comes after filepath
			break
		}
	}

	return outputOptions{
		printNumberOfBytes:      flagPrintNumberOfBytes,
		printNumberOfWords:      flagPrintNumberOfWords,
		printNumberOfLines:      flagPrintNumberOfLines,
		printNumberOfCharacters: flagPrintNumberOfCharacters,
	}, nil
}

func extractFilePaths(args []string) ([]string, error) {
	var filePaths []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}

		if !fileExists(arg) {
			return nil, fmt.Errorf("invalid file path: (%s)", arg)
		}

		filePaths = append(filePaths, arg)
	}

	return filePaths, nil
}

func fileExists(filepath string) bool {
	stat, err := os.Stat(filepath)

	if err != nil || stat.IsDir() {
		return false
	}
	return true
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-") && len(arg) > 1
}
