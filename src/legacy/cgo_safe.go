//go:build safe

package legacy

import "github.com/opennox/opennox/v1/common/memmap"

const cgoSafe = true

func init() {
	memmap.SetRuntimeChecks(true)
}
