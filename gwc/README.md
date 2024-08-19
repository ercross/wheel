# gwc

gwc (short for `go wc`) is an implementation of the Linux commandline program `wc` in Golang.
A character is represented as a `rune` in Go. Hence, this implementation makes of rune extensively

### Features (as copied with little modifications from terminal command `man wc`)
wc – word, line, character, and byte count

The wc utility displays the number of lines, words, and bytes contained
in each input file, or standard input (if no file is specified) to the
standard output.  

- A line is defined as a string of characters delimited
  by a ⟨newline⟩ character.  Characters beyond the final ⟨newline⟩
  character will not be included in the line count.

- A character is the smallest unit of human-readable text in a string. 
  It represents a single symbol, such as a letter, number, punctuation mark, or special character including whitespaces

- A word is defined as a string of characters delimited by white space
  characters.  White space characters are the set of characters for which
  the iswspace(3) function returns true.

- If more than one input file is specified, 
  a line of cumulative counts for all the files is displayed on
  a separate line after the output for the last file.

### Additional Feature
- Prettier output

### Usage
1. Install in your local environment using `go install github.com/ercross/wheel/gwc`
2. Synopsis: `gwc [OPTIONS] [file ...]`
    where the following options are available:

   -c      The number of bytes in each input file is written to the standard
           output.  This will cancel out any prior usage of the -m option.

   -l      The number of lines in each input file is written to the standard output.

   -m      The number of characters in each input file is written to the
           standard output.  If the current locale does not support
           multibyte characters, this is equivalent to the -c option.  This
           will cancel out any prior usage of the -c option.

   -w      The number of words in each input file is written to the standard output.

### Note about usage
- When an option is specified, wc only reports the information requested by
that option.  The order of output always takes the form of line, word,
byte, and file name.  The default action is equivalent to specifying the
-c, -l and -w options.

- If no files are specified, the standard input is used and no file name is
displayed.  The prompt will accept input until receiving EOF, or [^D] in most environments.

- File or input should contain only UTF-8 encoded character set

### Limitations
- OS support (Non-Unix): `gwc` has not been tested on non-unix based OS (e.g., Windows) 
  and might not function as expected on such platform
- OS support (pre-OS X): This program may not work properly on older versions of Mac OS (before OS X) because
  of some character implementation differences in pre-OS X and modern OS X. 
  For example, pre-OS X used a single carriage return (CR) character to represent a newline,
  ASCII code \r, with a value of 0x0D while modern OS X use /n in alignment with Linux-based OS 

