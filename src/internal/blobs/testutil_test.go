package blobs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// sourceCopy gives source-rewriting tests their own tree. Never point a test
// that calls Write, RewriteAccess or SplitBlob at the working checkout.
func sourceCopy(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.CopyFS(dir, os.DirFS("../..")))
	previous := Path()
	SetPath(dir)
	t.Cleanup(func() { SetPath(previous) })
}
