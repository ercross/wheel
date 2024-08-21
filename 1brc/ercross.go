package main

// generateTestFile generates a text file containing 10 billion rows of weather station report.
//
// The following specification about sample data contained in fromFilename
// is copied from the official website of 1brc at
// https://1brc.dev/#rules-and-limits
//
// Input value ranges are as follows
//   - Station name: non null UTF-8 string of min length 1 character
//     and max length 100 bytes (i.e. this could be 100 one-byte characters, or 50 two-byte characters, etc.)
//   - Temperature value: non null double between -99.9 (inclusive) and 99.9 (inclusive),
//     always with one fractional digit
//   - There is a maximum of 10,000 unique station names.
func generateTestFile(filename string, fromFilename string) {

}

func main() {
	testFile := "generated_stations.txt"
	generateTestFile()
}
