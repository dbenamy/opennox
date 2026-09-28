//go:build safe

package legacy

/*
#cgo CFLAGS: -g -O0
#cgo CFLAGS: -fsanitize=address
#cgo LDFLAGS: -fsanitize=address
*/
import "C"
import "github.com/opennox/opennox/v1/common/memmap"

const cgoSafe = true

func init() {
	memmap.SetRuntimeChecks(true)
}
