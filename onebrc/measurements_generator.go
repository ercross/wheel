package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// generateTestFile generates a text file containing 10 billion rows of weather station report.
//
// The following specification about sample data contained in fromFilename
// is copied from the official website of 1brc at
// https://1brc.dev/#rules-and-limits
//
// Input value ranges are as follows
//   - Station name: non-null UTF-8 string of min length 1 character
//     and max length 100 bytes (i.e. this could be 100 one-byte characters, or 50 two-byte characters, etc.)
//   - Temperature value: non-null double between -99.9 (inclusive) and 99.9 (inclusive),
//     always with one fractional digit
//   - There is a maximum of 10,000 unique station names.
func generateTestFile(filename string, fromFilename string) error {
	stations, err := uniqueStationNames(10_000, fromFilename)
	if err != nil {
		return err
	}

	_ = os.Remove(filename)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	mu := new(sync.Mutex)
	ch := make(chan string, 10)
	wg := new(sync.WaitGroup)

	numberOfWriters := 20
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < numberOfWriters; i++ {
		go writeMeasurements(ch, mu, ctx, file)
	}

	numberOfGenerators := runtime.GOMAXPROCS(runtime.NumCPU())
	totalRowsNeeded := 1_000_000_000
	rowsPerRoutine := totalRowsNeeded / numberOfGenerators
	for i := 0; i < numberOfGenerators; i++ {
		wg.Add(1)
		go generateStationTemperature(ch, wg, rowsPerRoutine, stations)
	}

	wg.Wait()
	cancel()
	return nil
}

func generateStationTemperature(ch chan string, wg *sync.WaitGroup, rows int, stations []string) {
	defer wg.Done()

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	minTemp := -99.9
	maxTemp := 99.9

	var (
		station     string
		temperature float64
	)

	for i := 0; i < rows; i++ {
		station = stations[r.Intn(len(stations))]
		temperature = minTemp + r.Float64()*(maxTemp-minTemp)
		ch <- fmt.Sprintf("%s;%f\n", station, temperature)
	}

}

func writeMeasurements(ch chan string, mu *sync.Mutex, ctx context.Context, file *os.File) {

	var builder strings.Builder
	var receivedCount int
	for {
		if receivedCount == 100 {
			mu.Lock()
			writer := bufio.NewWriter(file)
			_, err := writer.WriteString(builder.String())
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}

			receivedCount = 0
			builder.Reset()
			writer.Flush()
			mu.Unlock()
		}
		select {
		case msg := <-ch:
			builder.WriteString(msg)
			receivedCount++
		case <-ctx.Done():
			break
		}
	}
}

func uniqueStationNames(count int, fromFilename string) ([]string, error) {
	seedFile, err := os.Open(fromFilename)
	if err != nil {
		return nil, fmt.Errorf("error opening seed-file: %v", err)
	}
	defer seedFile.Close()
	reader := csv.NewReader(seedFile)
	reader.Comma = ';'
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading seed-file: %v", err)
	}

	m := make(map[string]struct{}, count)
	stations := make([]string, 0, count)
	for _, line := range lines {
		if len(m) == count {
			break
		}
		if _, ok := m[line[0]]; !ok {
			m[line[0]] = struct{}{}
			stations = append(stations, line[0])
		}
	}

	return stations, nil
}
