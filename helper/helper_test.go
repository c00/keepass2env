package helper

import (
	"os"
	"testing"

	"github.com/keepasssecrethelper/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelper_openDatabase(t *testing.T) {

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

func TestHelper_getPassword(t *testing.T) {
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

func TestHelper_updateOutputFile(t *testing.T) {
	const output = "../assets/output.env"
	err := os.WriteFile(output, []byte(`ANOTHER_THING=first`), 0666)
	require.NoError(t, err)

	err = updateOutputFile(output, []entryWithPass{
		{Entry: config.Entry{EnvName: "THING"}, password: "some-pass"},
		{Entry: config.Entry{EnvName: "ANOTHER_THING"}, password: "second"},
	})
	require.NoError(t, err)

	got, err := os.ReadFile(output)
	require.NoError(t, err)

	expected := "ANOTHER_THING=second\nTHING=some-pass"

	assert.Equal(t, expected, string(got))
}
