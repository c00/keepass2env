package outputter

import "github.com/c00/keepass2env/config"

type Outputter interface {
	Output(entries []config.EntryWithSecret) error
}
