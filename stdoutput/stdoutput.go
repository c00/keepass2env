package stdoutput

import (
	"fmt"
	"io"

	"github.com/c00/keepass2env/config"
	"github.com/c00/keepass2env/outputter"
)

var _ outputter.Outputter = (*StdOutput)(nil)

type StdOutput struct {
	Writer io.Writer
}

func (o *StdOutput) Output(entries []config.EntryWithSecret) error {
	for _, entry := range entries {
		_, err := fmt.Fprintf(o.Writer, "%v=%v\n", entry.EnvName, entry.Secret)
		if err != nil {
			return fmt.Errorf("cannot write to output: %w", err)
		}
	}

	return nil
}
