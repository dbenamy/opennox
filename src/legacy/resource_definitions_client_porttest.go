//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/binfile"
	"unsafe"
)

func PortTestResourceLightConstants() (*uint64, *uint64) {
	return (*uint64)(unsafe.Pointer(&qword_581450_9544)), (*uint64)(unsafe.Pointer(&qword_581450_9552))
}

func PortTestResourceClientParser(kind string, typ *client.ObjectType, f *binfile.MemFile, input unsafe.Pointer) bool {
	return resourceClientParser(kind, typ, f, (*byte)(input))
}
