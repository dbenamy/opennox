package blobs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatAccesses(t *testing.T) {
	sourceCopy(t)
	err := FormatAccesses()
	require.NoError(t, err)
}
