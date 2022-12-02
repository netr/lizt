package lizt

import (
	"bufio"
	"math/rand"
	"os"
)

func OpenFile(filename string) (*os.File, error) {
	return os.Open(filename)
}

func ReadFromFile(filename string) ([]string, error) {
	file, err := OpenFile(filename)
	if err != nil {
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func FileLineCount(filename string) (int, error) {
	file, err := OpenFile(filename)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}

func RepeatLines(lines []string, times int) []string {
	res := make([]string, len(lines)*times)
	idx := 0
	for i := 0; i < times; i++ {
		for j := i + 1; j < len(lines); j++ {
			res[idx] = lines[i]
			idx++
		}
	}
	return res
}

func GetDuplicateLines(lines []string) []string {
	var duplicates []string
	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines); j++ {
			if lines[i] == lines[j] {
				duplicates = append(duplicates, lines[i])
			}
		}
	}
	return duplicates
}

func WriteToFile(lines []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func Shuffle(lines []string) []string {
	res := make([]string, len(lines))
	perm := rand.Perm(len(lines))
	for i, v := range perm {
		res[v] = lines[i]
	}
	return res
}
