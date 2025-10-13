package runner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func updateOutputFile(path string, entries []entryWithPass) error {
	lines, err := readFile(path)
	if err != nil {
		return fmt.Errorf("cannot read output file: %w", err)
	}

	// Start updating
	for _, entry := range entries {
		replaced := false
		for i, line := range lines {
			if strings.HasPrefix(line, entry.EnvName) {
				// Replace
				lines[i] = fmt.Sprintf("%v=%v", entry.EnvName, entry.password)
				replaced = true
			}
		}

		if !replaced {
			// Add
			lines = append(lines, fmt.Sprintf("%v=%v", entry.EnvName, entry.password))
		}
	}

	// Write out
	return writeFile(path, strings.Join(lines, "\n"))
}

func readFile(path string) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open or create file '%v': %w", path, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	return lines, nil
}

func writeFile(path string, contents string) error {
	return os.WriteFile(path, []byte(contents), 0666)
}
