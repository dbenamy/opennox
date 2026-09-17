//go:build porttest

package legacy

import (
	"github.com/opennox/libs/ifs"
	"github.com/opennox/opennox/v1/internal/binfile"
	"unsafe"
)

func PortTestServerConfigAdmission(op string, index int, p, q unsafe.Pointer) unsafe.Pointer {
	switch op {
	case "allowed-add":
		return serverConfigAllowedAdd((*uint16)(p))
	case "blocked-add":
		return serverConfigBlockedAdd(int32(index), (*uint16)(p), (*byte)(q))
	case "allowed-remove":
		return serverConfigAllowedRemove(int32(index))
	case "blocked-remove":
		serverConfigBlockedRemove(int32(index))
		return nil
	case "allowed-first":
		return unsafe.Pointer(serverConfigAllowedFirst())
	case "allowed-next":
		return unsafe.Pointer(serverConfigAllowedNext((*serverConfigAllowed)(p)))
	case "blocked-first":
		return unsafe.Pointer(serverConfigBlockedFirst())
	case "blocked-next":
		return unsafe.Pointer(serverConfigBlockedNext((*serverConfigBlocked)(p)))
	case "expire":
		serverConfigExpire()
		return nil
	case "close":
		return serverConfigAdmissionClose()
	default:
		panic(op)
	}
}
func PortTestServerConfigFile(op, path string) int {
	switch op {
	case "read":
		return int(serverConfigFileRead(path))
	case "write":
		return int(serverConfigFileWrite(path))
	case "allowed", "blocked":
		file, err := ifs.Open(path)
		if err != nil {
			return -1
		}
		bf := binfile.NewTextFile(file)
		defer bf.Close()
		if op == "allowed" {
			return int(serverConfigFileReadAllowed(bf))
		}
		return int(serverConfigFileReadBlocked(bf))
	default:
		panic(op)
	}
}
func PortTestServerConfigRulePopulate(settings unsafe.Pointer) int {
	return serverConfigRulePopulate((*byte)(settings))
}
func PortTestServerConfigAdmissionInit() int { return int(serverConfigAdmissionInit()) }
