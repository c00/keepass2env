package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/keepasssecrethelper/config"
	"github.com/tobischo/gokeepasslib/v3"
)

type Helper struct {
	Params HelperParams

	db *gokeepasslib.Database
}

type HelperParams struct {
	DatabasePassword string
	DatabasePath     string
	KeyfilePath      string
	OutputPath       string
	Config           config.Config
}

type entryWithPass struct {
	config.Entry
	password string
}

func (h *Helper) Run() error {
	// Open database
	err := h.openDatabase()
	if err != nil {
		return fmt.Errorf("cannot open database: %w", err)
	}

	// For each entry, get password
	entries := []entryWithPass{}
	for _, entry := range h.Params.Config.Entries {
		e := entryWithPass{Entry: entry}
		pass, err := h.getPassword(e.KeepassPath)
		if err != nil {
			return fmt.Errorf("cannot get password for entry '%v': %w", e.KeepassPath, err)
		}

		e.password = pass
		entries = append(entries, e)
	}

	// Find and replace in output
	err = updateOutputFile(h.Params.OutputPath, entries)
	if err != nil {
		return fmt.Errorf("cannot write to file '%v': %w", h.Params.OutputPath, err)
	}

	fmt.Printf("Written secrets to: %v\n", h.Params.OutputPath)

	return nil
}

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

func (h *Helper) openDatabase() error {
	if h.db != nil {
		return nil
	}

	file, err := os.Open(h.Params.DatabasePath)
	if err != nil {
		return fmt.Errorf("cannot open database file: %w", err)
	}

	defer file.Close()

	h.db = gokeepasslib.NewDatabase()

	if h.Params.KeyfilePath != "" {
		creds, err := gokeepasslib.NewPasswordAndKeyCredentials(h.Params.DatabasePassword, h.Params.KeyfilePath)
		if err != nil {
			return fmt.Errorf("cannot open database: %w", err)
		}
		h.db.Credentials = creds
	} else {
		h.db.Credentials = gokeepasslib.NewPasswordCredentials(h.Params.DatabasePassword)
	}

	_ = gokeepasslib.NewDecoder(file).Decode(h.db)

	h.db.UnlockProtectedEntries()

	return nil
}

func (h *Helper) getPassword(path string) (string, error) {
	// Navigate to the path
	parts := strings.Split(path, "/")

	for _, group := range h.db.Content.Root.Groups {
		pass, err := h.navigate(group, parts)
		if err == nil {
			return pass, nil
		}
	}

	return "", fmt.Errorf("path not found: %v", path)
}

func (h *Helper) navigate(group gokeepasslib.Group, parts []string) (string, error) {
	// Find the leaf node
	if len(parts) == 1 {
		for _, entry := range group.Entries {
			if entry.GetTitle() == parts[0] {
				return entry.GetPassword(), nil
			}
		}
		return "", fmt.Errorf("entry not found: %v", parts[0])
	}

	// Dig deeper
	for _, group := range group.Groups {
		if group.Name == parts[0] {
			return h.navigate(group, parts[1:])
		}
	}
	return "", fmt.Errorf("folder not found: %v", parts[0])
}
