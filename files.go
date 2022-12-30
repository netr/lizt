package lizt

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// OpenFile opens a file
func OpenFile(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("os.Open(): %s -> %w", filename, err)
	}
	return f, nil
}

// ReadFromFile reads a file into a slice of strings
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
		lines[idx] = strings.TrimSpace(scanner.Text())
		idx++
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("ReadFromFile(): %s -> %w", filename, scanner.Err())
	}

	return lines, nil
}

// DeleteFile takes a path and deletes the file
func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("os.Remove(): %s -> %w", path, err)
	}
	return nil
}

func FileToMap(pathname string) (map[string]struct{}, error) {
	lines, err := ReadFromFile(pathname)
	if err != nil {
		return nil, err
	}

	m := make(map[string]struct{}, len(lines))
	for _, l := range lines {
		m[l] = struct{}{}
	}
	return m, nil
}

// FileLineCount returns the number of lines in a file
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

// RepeatLines repeats lines in a file
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

// GetDuplicateLines returns duplicate lines in a file
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

// WriteToFile writes a string to a file
func WriteToFile(lines []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("os.Create(): %s -> %w", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	sb := &strings.Builder{}
	for _, line := range lines {
		sb.WriteString(line + "\n")
	}

	_, err = writer.WriteString(sb.String())
	if err != nil {
		return fmt.Errorf("writer.WriteString(): %s -> %w", filename, err)
	}

	if err = writer.Flush(); err != nil {
		return fmt.Errorf("writer.Flush(): %s -> %w", filename, err)
	}

	return nil
}

// Shuffle shuffles a slice of strings
func Shuffle(lines []string) []string {
	res := make([]string, len(lines))
	perm := rand.Perm(len(lines))
	for i, v := range perm {
		res[v] = lines[i]
	}
	return res
}

// ReadFileToString reads a file into a string
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

// DoesFileExist checks if a file exists
func DoesFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetCurrentDir returns the current directory
func GetCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errors.New("failed to get current directory")
	}
	return dir, nil
}
