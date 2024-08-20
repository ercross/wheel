package main

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	filename := "./test-file.txt"
	err := generateTextFile(filename)
	if err != nil {
		t.Errorf("failed to generate file: %v", err)
	}

	args := []string{"-cw", filename}
	r, err := run(args)
	if err != nil {
		t.Errorf("run error: %v", err)
	}

	if r.numberOfWords != 10_000_000 {
		t.Errorf("wrong number of words: expected 10,000,000 got %d", r.numberOfWords)
	}
}

func TestCountLines(t *testing.T) {
	input := []byte("  Hello there,\n   World!\n  This is a test.\n \n ")
	count := countLines(input)
	if count != 4 {
		t.Errorf("count should be 4, got %d", count)
	}
}

func TestCountCharacters(t *testing.T) {
	t.Run("With Chinese Characters", func(t *testing.T) {
		input := []byte("Hello, 世界!")
		count, err := countCharacters(input)
		if err != nil {
			t.Error(err)
		}
		if count != 10 {
			t.Errorf("count should be 10, got %d", count)
		}
	})

	t.Run("With Emoji Characters", func(t *testing.T) {
		input := []byte("😊🌍🌟")
		count, err := countCharacters(input)
		if err != nil {
			t.Error(err)
		}
		if count != 3 {
			t.Errorf("count should be 3, got %d", count)
		}
	})
}

func BenchmarkCountWords(b *testing.B) {
	b.Run("SmallInput", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := countWords([]byte("  Hello,   world!  This is a test. "))
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("LargeInput", func(b *testing.B) {
		tenMillion := 10_000_000
		input := generateInput(tenMillion)
		for i := 0; i < b.N; i++ {
			_, err := countWords(input)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

}

func TestCountWords(t *testing.T) {
	input := []byte("  Hello,   world!  This is a test. ")
	count, err := countWords(input)
	if err != nil {
		t.Fatalf("countWords failed: %v", err)
	}
	if count != 6 {
		t.Errorf("Expected 6, got %d", count)
	}
}

func FuzzCountWords(f *testing.F) {
	tenMillion := 10_000_000
	f.Add(generateInput(tenMillion))
	f.Fuzz(func(t *testing.T, input []byte) {
		count, err := countWords(input)
		if err != nil {
			t.Fatalf("countWords failed: %v", err)
		}
		if count < 0 {
			t.Fatalf("countWords returned negative count")
		}
	})
}

func generateInput(wordCount int) []byte {
	var result []byte
	words := []string{"However", "to", "an", "extent", "NFC", "has", "always", "been",
		"optional", "technology", "rather", "than", "an", "essential"}

	for i := 0; i < wordCount; i++ {
		result = append(result, words[i%len(words)]...)
		result = append(result, ' ') // add whitespaces
	}
	return result
}

func generateTextFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the data to the file
	data := generateInput(10_000_000)
	_, err = file.WriteString(string(data))
	if err != nil {
		return err
	}

	return nil
}
