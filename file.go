package httpmatter

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

// makeFilePath makes a file path for a given namespace and file name
func makeFilePath(baseDir, namespace, fileName, extension string) string {
	return filepath.Join(baseDir, namespace, fileName+extension)
}

// openFile opens a file for reading
func openFile(filepath string) (*os.File, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// readFile reads a file and returns the frontmatter and content
func readFile(filepath string) (string, string, error) {
	file, err := openFile(filepath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()
	sections := []*bytes.Buffer{
		bytes.NewBufferString(""),
		bytes.NewBufferString(""),
	}
	index := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if isContentLine(line) {
			index = 1
		}
		sections[index].WriteString(line + "\n")
	}
	return sections[0].String(), sections[1].String(), nil
}

func isContentLine(line string) bool {
	validPrefixes := []string{
		"HTTP/", "GET ", "POST ", "PUT ", "DELETE ", "PATCH ", "HEAD ",
	}
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}
