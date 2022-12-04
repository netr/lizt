package lizt

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func OpenFile(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("os.Open(): %s -> %w", filename, err)
	}
	return f, nil
}

func ReadFromFile(filename string) ([]string, error) {
	lc, _ := FileLineCount(filename)
	idx := 0
	lines := make([]string, lc)

	file, err := OpenFile(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines[idx] = scanner.Text()
		idx++
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("ReadFromFile(): %s -> %w", filename, scanner.Err())
	}

	return lines, nil
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

	if scanner.Err() != nil {
		return 0, fmt.Errorf("FileLineCount(): %s -> %w", filename, scanner.Err())
	}

	return count, nil
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
		return fmt.Errorf("os.Create(): %s -> %w", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("writer.WriteString(): %s -> %w", filename, err)
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

func ReadFileToString(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("os.Open(): %s -> %w", filename, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("file.Stat(): %s -> %w", filename, err)
	}

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return "", fmt.Errorf("file.Read(): %s -> %w", filename, err)
	}

	return string(bs), nil
}
func DoesFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
