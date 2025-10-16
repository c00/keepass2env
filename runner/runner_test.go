package runner

import (
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunner_openDatabase(t *testing.T) {

	h := Helper{
		Params: HelperParams{
			DatabasePath:     "../assets/with-keyfile.kdbx",
			DatabasePassword: "ilikebeans",
			KeyfilePath:      "../assets/keyfile.key",
		},
	}

	err := h.openDatabase()
	require.NoError(t, err)
}

func TestRunner_getPassword(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path    string
		want    string
		wantErr bool
	}{
		{
			name: "root node",
			path: "Test Entry 1",
			want: "spongebob1",
		},
		{
			name: "one deep",
			path: "Foo/Test Entry 2",
			want: "spongebob2",
		},
		{
			name: "two deep",
			path: "Foo/Bar/Test Entry 3",
			want: "spongebob3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Helper{
				Params: HelperParams{
					DatabasePath:     "../assets/with-keyfile.kdbx",
					DatabasePassword: "ilikebeans",
					KeyfilePath:      "../assets/keyfile.key",
				},
			}

			err := h.openDatabase()
			require.NoError(t, err)

			got, gotErr := h.getPassword(tt.path)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getPassword() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getPassword() succeeded unexpectedly")
			}

			if tt.want != got {
				t.Errorf("getPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunner_getAttribute(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path    string
		attr    string
		want    string
		wantErr bool
	}{
		{
			name: "root node",
			path: "Test Entry 1",
			want: "Here\nare\nsome\nvalues",
			attr: "some-attribute",
		},
		{
			name:    "root node",
			path:    "Test Entry 1",
			attr:    "not-existing",
			wantErr: true,
		},
		{
			name:    "root node",
			path:    "Not Existing",
			attr:    "some-attribute",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Helper{
				Params: HelperParams{
					DatabasePath:     "../assets/with-keyfile.kdbx",
					DatabasePassword: "ilikebeans",
					KeyfilePath:      "../assets/keyfile.key",
				},
			}

			err := h.openDatabase()
			require.NoError(t, err)

			got, gotErr := h.getAttribute(tt.path, tt.attr)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getAttribute() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getAttribute() succeeded unexpectedly")
			}

			if tt.want != got {
				t.Errorf("getAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelperParams_expandPaths(t *testing.T) {
	h := Helper{
		Params: HelperParams{
			DatabasePath: "~/thing.kdbx",
			KeyfilePath:  "",
		},
	}

	err := h.Params.expandPaths()
	require.NoError(t, err)

	usr, err := user.Current()
	require.NoError(t, err)

	assert.Equal(t, filepath.Join(usr.HomeDir, "thing.kdbx"), h.Params.DatabasePath)
	assert.Equal(t, "", h.Params.KeyfilePath)
}
