package runner

import (
	"fmt"
	"os"
	"strings"

	"github.com/c00/keepass2env/config"
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
	Entries          []config.Entry
}

type entryWithPass struct {
	config.Entry
	password string
}

func (p *HelperParams) expandPaths() error {
	path, err := ExpandPath(p.DatabasePath)
	if err != nil {
		return fmt.Errorf("cannot expand path '%v': %w", p.DatabasePath, err)
	}
	p.DatabasePath = path

	path, err = ExpandPath(p.OutputPath)
	if err != nil {
		return fmt.Errorf("cannot expand path '%v': %w", p.OutputPath, err)
	}
	p.OutputPath = path

	path, err = ExpandPath(p.KeyfilePath)
	if err != nil {
		return fmt.Errorf("cannot expand path '%v': %w", p.KeyfilePath, err)
	}
	p.KeyfilePath = path

	return nil
}

func (h *Helper) Run() error {
	// Expand paths
	err := h.Params.expandPaths()
	if err != nil {
		return fmt.Errorf("cannot expand paths: %w", err)
	}

	// Open database
	err = h.openDatabase()
	if err != nil {
		return fmt.Errorf("cannot open database: %w", err)
	}

	// For each entry, get password
	entries := []entryWithPass{}
	for _, entry := range h.Params.Entries {
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

	err = gokeepasslib.NewDecoder(file).Decode(h.db)
	if err != nil {
		return fmt.Errorf("cannot decode database: %w", err)
	}

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
