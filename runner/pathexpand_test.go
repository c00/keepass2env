package runner_test

import (
	"os/user"
	"path/filepath"
	"testing"

	"github.com/c00/keepass2env/runner"
	"github.com/stretchr/testify/require"
)

func TestExpandPath(t *testing.T) {
	usr, err := user.Current()
	require.NoError(t, err)
	homeDir := usr.HomeDir

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path    string
		want    string
		wantErr bool
	}{
		{name: "no tilde", path: "/home/coo/thing", want: "/home/coo/thing"},
		{name: "only tilde", path: "~", want: homeDir},
		{name: "home path", path: "~/things", want: filepath.Join(homeDir, "things")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := runner.ExpandPath(tt.path)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExpandPath() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExpandPath() succeeded unexpectedly")
			}

			if tt.want != got {
				t.Errorf("ExpandPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
