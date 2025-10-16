package keyringoutput_test

import (
	"testing"

	"github.com/c00/keepass2env/config"
	"github.com/c00/keepass2env/keyringoutput"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestRunner_updateOutputFile(t *testing.T) {
	out := keyringoutput.KeyringOutput{
		Service: "keepass2keyring-test",
	}

	// set initial value
	err := keyring.Set(out.Service, "THING", "foo")
	require.NoError(t, err)

	entries := []config.EntryWithSecret{
		{Entry: config.Entry{EnvName: "THING"}, Secret: "some-pass"},
		{Entry: config.Entry{EnvName: "ANOTHER_THING"}, Secret: "second"},
	}

	err = out.Output(entries)
	require.NoError(t, err)

	//check new values
	for _, entry := range entries {
		got, err := keyring.Get(out.Service, entry.EnvName)
		require.NoError(t, err)
		assert.Equal(t, entry.Secret, got)
	}

}
