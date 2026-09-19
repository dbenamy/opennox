package blobs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadMemmap(t *testing.T) {
	sourceCopy(t)
	m, err := ReadMemmap()
	require.NoError(t, err)
	// The port intentionally removes C owners. Validate preserved mappings rather
	// than requiring the historical number of variables to remain in the source.
	require.NotEmpty(t, m.Vars)
	require.NoError(t, m.Write())
	again, err := ReadMemmap()
	require.NoError(t, err)
	require.ElementsMatch(t, m.Vars, again.Vars)
}

func TestReadMemmapFixture(t *testing.T) {
	previous := Path()
	SetPath(t.TempDir())
	t.Cleanup(func() { SetPath(previous) })
	path := Path(memmapGo2)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	const source = `package nox
var mappings = []struct {
 Blob, Off, Size uintptr
 Name string
}{
 {0x2000, 12, 1, "last"}, // 0x200C
 // {0x1000, 8, 2, "disabled"}, // TODO: retain this mapping
 {0x2000, 4, 8, "wide"},
 {0x1000, 0, 4, "first"},
}
`
	require.NoError(t, os.WriteFile(path, []byte(source), 0600))
	want := []Var{
		{Blob: 0x1000, Off: 0, Size: 4, Name: "first"},
		{Blob: 0x1000, Off: 8, Size: 2, Name: "disabled", Disabled: true, Comment: "TODO: retain this mapping"},
		{Blob: 0x2000, Off: 4, Size: 8, Name: "wide"},
		{Blob: 0x2000, Off: 12, Size: 1, Name: "last"},
	}
	m, err := ReadMemmap()
	require.NoError(t, err)
	require.Equal(t, want, m.Vars)
	require.NoError(t, m.Write())
	again, err := ReadMemmap()
	require.NoError(t, err)
	require.Equal(t, want, again.Vars)
}
