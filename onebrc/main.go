package main

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	starts := time.Now()
	oneBillion := 1_000_000_000
	err := generateTestFile("./measurements.txt", "weather_stations.csv", oneBillion)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generation took:", time.Since(starts))
}

type aggregate struct {

	// if all data from the 1 billion row is loaded into stationTemp,
	// it is expected to grow up to a size of 40GiB at most
	//
	// string storage: up to 108bytes (i.e., string-header:8bytes data:100bytes)
	// float32 storage: data consumes 4bytes
	// []float32 slice header: 24bytes (i.e., 8bytes each for pointers to data, length, and capacity)
	// key memory: 108bytes * 10,000 unique station names = 1,080,000bytes (1.08MB)
	// An entry of value consumes: 1 million of (float32) + 24bytes(slice-header) = (1,000,000 * 4) + 24 = 4,000,024 bytes
	// Total memory consumption for 10,000 unique station names: 4,000,024 * 10,000 = 40,000,240,000 bytes = 40GiB
	// Max total memory consumption: 1.08MB + 40GiB = 40.00108 GB
	stationTemp map[string][]float32
	mu          *sync.RWMutex
}

func bufferReadFileStreamInput(filename string) (err error) {
	a := aggregate{
		// if there are 10,000 unique station names,
		// the average of it is 5000 for an initial capacity
		stationTemp: make(map[string][]float32, 5000),
		mu:          nil,
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var chunk string
	for {

		// todo try use bufio.Scanner
		chunk, err = reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("error reading file: %v", err)
		}
		parts := strings.Split(chunk, ";")
		temp, _ := strconv.ParseFloat(parts[1], 32)
		a.aggregateResult(parts[0], float32(temp))
	}
}

func (a *aggregate) aggregateResult(stationName string, temp float32) {
	a.mu.Lock()
	defer a.mu.Unlock()

	v, ok := a.stationTemp[stationName]
	if !ok {

		// for one billion rows and a maximum of 10,000 unique station names,
		// each station name can appear in an average of 1,000,000 rows
		a.stationTemp[stationName] = make([]float32, 0, 600_000)
		a.stationTemp[stationName] = append(a.stationTemp[stationName], temp)
	} else {
		a.stationTemp[stationName] = append(v, temp)
	}
}

func (a *aggregate) calculateResult() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var builder strings.Builder
	for k, v := range a.stationTemp {

		slices.SortFunc(v, func(a, b float32) int {
			return cmp.Compare(a, b)
		})

		mean := calculateMean(v)
		median := calculateMedian(v)
		builder.WriteString(fmt.Sprintf("%s;%f;%f;%f\n", k, mean, median, v[(len(v)-1)]))
	}

	return builder.String()
}

func calculateMean(temps []float32) float32 {
	total := float32(0.0)
	for _, v := range temps {
		total = total + v
	}

	return total / float32(len(temps))
}

func calculateMedian(temps []float32) float32 {
	// Calculate the middle index
	middleIndex := len(temps) / 2

	// If the length is odd, return the middle value
	if len(temps)%2 != 0 {
		return temps[middleIndex]
	}

	// If the length is even, return the average of the two middle values
	return (temps[middleIndex-1] + temps[middleIndex]) / 2
}
