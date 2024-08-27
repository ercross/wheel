package main

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	starts := time.Now()
	run()
	fmt.Printf("Took %v to complete", time.Since(starts))
}

func run() {
	lineChan := make(chan string, 100)
	stationStatsChan := make(chan stationStats, 30)
	chunkChan := make(chan map[string][]float32, 100)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go processChunk(chunkChan, stationStatsChan, wg)

	wg.Add(1)
	go aggregateResult(lineChan, chunkChan, wg)

	wg.Add(1)
	go mergeResult(stationStatsChan, wg)

	wg.Add(1)
	go readMeasurements("test.txt", lineChan, wg)

	wg.Wait()
}

type stats struct {
	mean float32
	min  float32
	max  float32
}

type stationStats struct {
	station string
	stats   stats
}

func readMeasurements(filename string, lineChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("error opening file: %v", err))
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string

	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				close(lineChan)
				break
			}
			fmt.Printf("error reading file: %v", err)
			continue
		}

		lineChan <- line
	}
}

func aggregateResult(lineChan chan string, chunkChan chan map[string][]float32, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		parts     []string
		temp      float64
		err       error
		lineCount int
		ok        bool
		line      string
	)
	maxChunkSize := 1000
	chunk := make(map[string][]float32, maxChunkSize)

	for {

		line, ok = <-lineChan
		if !ok {
			close(chunkChan)
			break
		}

		if lineCount == maxChunkSize {
			chunkChan <- maps.Clone(chunk)
			clear(chunk)
			continue
		}

		// parse line
		line = strings.ReplaceAll(line, "\n", "")
		parts = strings.Split(line, ";")
		if len(parts) != 2 {
			// drop invalid lines
			fmt.Printf("invalid measurement entry: %s", line)
			continue
		}
		temp, err = strconv.ParseFloat(parts[1], 32)
		if err != nil {
			// drop lines with invalid temperatures
			fmt.Printf("error parsing temperature value %s: %v", parts[1], err)
			continue
		}

		// insert new measurement
		temps, ok := chunk[parts[0]]
		if !ok {
			chunk[parts[0]] = []float32{float32(temp)}
		} else {
			chunk[parts[0]] = append(temps, float32(temp))
		}

		lineCount++
	}

}

func processChunk(chunkChan chan map[string][]float32, stationStatsChan chan stationStats, wg *sync.WaitGroup) {
	defer wg.Done()
	var chunk map[string][]float32
	var ok bool
	for {
		chunk, ok = <-chunkChan
		if !ok {
			close(stationStatsChan)
			break
		}
		for stationName, temperatures := range chunk {
			slices.SortFunc(temperatures, func(a, b float32) int {
				return cmp.Compare(a, b)
			})
			chunk[stationName] = temperatures

			stationStatsChan <- stationStats{
				station: stationName,
				stats: stats{
					mean: calculateMean(temperatures),
					min:  findMin(temperatures),
					max:  findMax(temperatures),
				},
			}
		}
	}
}

func mergeResult(stationStatsChan chan stationStats, wg *sync.WaitGroup) {
	defer wg.Done()
	result := make(map[string]stats)
	var stationStat stationStats
	var ok bool
	for {
		stationStat, ok = <-stationStatsChan
		if !ok {
			break
		}
		result[stationStat.station] = newStat(result[stationStat.station], stationStat.stats)

	}
	writeResult("result.txt", result)
}

func newStat(s1, s2 stats) stats {
	return stats{
		mean: calculateMean([]float32{s1.mean, s2.mean}),
		min:  findMin([]float32{s1.min, s2.min}),
		max:  findMax([]float32{s1.max, s2.max}),
	}
}

func writeResult(filename string, result map[string]stats) {
	_ = os.Remove(filename)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("error opening file: %v", err))
	}
	defer file.Close()

	var builder strings.Builder
	for k, v := range result {
		builder.WriteString(fmt.Sprintf("%s;%f;%f;%f\n", k, v.min, v.mean, v.max))
	}

	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString(builder.String())
}

func findMax(temps []float32) float32 {
	return temps[(len(temps) - 1)]
}

func calculateMean(temps []float32) float32 {
	total := float32(0.0)
	for _, v := range temps {
		total = total + v
	}

	return total / float32(len(temps))
}

func findMin(temps []float32) float32 {
	return temps[0]
}
