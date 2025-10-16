package fileoutput_test

import (
	"os"
	"testing"

	"github.com/c00/keepass2env/config"
	"github.com/c00/keepass2env/fileoutput"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunner_updateOutputFile(t *testing.T) {
	const output = "../assets/output.env"
	err := os.WriteFile(output, []byte(`ANOTHER_THING=first`), 0666)
	require.NoError(t, err)

	out := fileoutput.FileOutput{
		Path: output,
	}

	err = out.Output([]config.EntryWithSecret{
		{Entry: config.Entry{EnvName: "THING"}, Secret: "some-pass"},
		{Entry: config.Entry{EnvName: "ANOTHER_THING"}, Secret: "second"},
	})
	require.NoError(t, err)

	got, err := os.ReadFile(output)
	require.NoError(t, err)

	expected := "ANOTHER_THING=second\nTHING=some-pass"

	assert.Equal(t, expected, string(got))
}
