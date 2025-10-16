package runner

import (
	"fmt"
	"os"
	"strings"

	"github.com/c00/keepass2env/config"
	"github.com/c00/keepass2env/outputter"
	"github.com/tobischo/gokeepasslib/v3"
)

type Helper struct {
	Params HelperParams

	Output outputter.Outputter
	db     *gokeepasslib.Database
}

type HelperParams struct {
	DatabasePassword string
	DatabasePath     string
	KeyfilePath      string
	Entries          []config.Entry
}

func (p *HelperParams) expandPaths() error {
	path, err := ExpandPath(p.DatabasePath)
	if err != nil {
		return fmt.Errorf("cannot expand path '%v': %w", p.DatabasePath, err)
	}
	p.DatabasePath = path

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

	// For each entry, get password / attribute
	entries := []config.EntryWithSecret{}
	for _, entry := range h.Params.Entries {
		e := config.EntryWithSecret{Entry: entry}

		var secret string
		if e.Attribute == "" || e.Attribute == "password" {
			secret, err = h.getPassword(e.KeepassPath)
			if err != nil {
				return fmt.Errorf("cannot get password for entry '%v': %w", e.KeepassPath, err)
			}

		} else {
			// Find attribute
			secret, err = h.getAttribute(e.KeepassPath, e.Attribute)
			if err != nil {
				return fmt.Errorf("cannot get attribute for entry '%v': %w", e.KeepassPath, err)
			}
		}
		e.Secret = secret
		entries = append(entries, e)
	}

	// Find and replace in output
	err = h.Output.Output(entries)
	if err != nil {
		return fmt.Errorf("cannot output entries: %w", err)
	}

	fmt.Printf("Written secrets")

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
			return fmt.Errorf("cannot open key file: %w", err)
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
		entry, err := h.navigate(group, parts)
		if err == nil {
			return entry.GetPassword(), nil
		}
	}

	return "", fmt.Errorf("path not found: %v", path)
}

func (h *Helper) getAttribute(path, attribute string) (string, error) {
	// Navigate to the path
	parts := strings.Split(path, "/")

	for _, group := range h.db.Content.Root.Groups {
		entry, err := h.navigate(group, parts)
		if err == nil {
			valueData := entry.Get(attribute)
			if valueData == nil {
				return "", fmt.Errorf("attribute %v not found in %v", attribute, path)
			}
			return valueData.Value.Content, nil
		}
	}

	return "", fmt.Errorf("path not found: %v", path)
}

func (h *Helper) navigate(group gokeepasslib.Group, parts []string) (gokeepasslib.Entry, error) {
	// Find the leaf node
	if len(parts) == 1 {
		for _, entry := range group.Entries {
			if entry.GetTitle() == parts[0] {
				return entry, nil
			}
		}
		return gokeepasslib.Entry{}, fmt.Errorf("entry not found: %v", parts[0])
	}

	// Dig deeper
	for _, group := range group.Groups {
		if group.Name == parts[0] {
			return h.navigate(group, parts[1:])
		}
	}
	return gokeepasslib.Entry{}, fmt.Errorf("folder not found: %v", parts[0])
}
