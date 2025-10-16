package keyringoutput

import (
	"fmt"

	"github.com/c00/keepass2env/config"
	"github.com/c00/keepass2env/outputter"
	"github.com/zalando/go-keyring"
)

var _ outputter.Outputter = (*KeyringOutput)(nil)

type KeyringOutput struct {
	Service string
}

func (o *KeyringOutput) Output(entries []config.EntryWithSecret) error {
	for _, entry := range entries {
		err := keyring.Set(o.Service, entry.EnvName, entry.Secret)
		if err != nil {
			return fmt.Errorf("cannot save entry %v to keyring: %w", entry.EnvName, err)
		}
	}

	return nil
}
