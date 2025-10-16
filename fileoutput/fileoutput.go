package fileoutput

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/c00/keepass2env/config"
	"github.com/c00/keepass2env/outputter"
	"github.com/c00/keepass2env/runner"
)

var _ outputter.Outputter = (*FileOutput)(nil)

type FileOutput struct {
	Path string
}

func (o *FileOutput) Output(entries []config.EntryWithSecret) error {
	expanded, err := runner.ExpandPath(o.Path)
	if err != nil {
		return fmt.Errorf("cannot expand path: %w", err)
	}
	o.Path = expanded

	lines, err := o.readFile()
	if err != nil {
		return fmt.Errorf("cannot read output file: %w", err)
	}

	// Start updating
	for _, entry := range entries {
		replaced := false
		for i, line := range lines {
			if strings.HasPrefix(line, entry.EnvName) {
				// Replace
				lines[i] = fmt.Sprintf("%v=%v", entry.EnvName, entry.Secret)
				replaced = true
			}
		}

		if !replaced {
			// Add
			lines = append(lines, fmt.Sprintf("%v=%v", entry.EnvName, entry.Secret))
		}
	}

	// Write out
	return writeFile(o.Path, strings.Join(lines, "\n"))
}

func (o *FileOutput) readFile() ([]string, error) {
	file, err := os.OpenFile(o.Path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open or create file '%v': %w", o.Path, err)
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
