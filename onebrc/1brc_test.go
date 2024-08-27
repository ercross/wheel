package main

import (
	"os"
	"testing"
)

func TestGenerateTestFile(t *testing.T) {
	_ = os.Remove("test.txt")
	err := generateTestFile("test.txt", "weather_stations.csv", 1000)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkGenerateTestFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = os.Remove("test.txt")
		err := generateTestFile("test.txt", "weather_stations.csv", 100_000)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadMeasurements(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run()
	}
}
