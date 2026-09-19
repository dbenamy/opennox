//go:build porttest

package legacy

/*
#include "defs.h"
#include "memfile.h"
extern uint64_t qword_581450_9544;
extern uint64_t qword_581450_9552;
bool nox_parse_thing_light_dir(nox_thing*, nox_memfile*, char*);
bool nox_parse_thing_light_penumbra(nox_thing*, nox_memfile*, char*);
bool nox_parse_thing_client_update(nox_thing*, nox_memfile*, char*);
bool nox_parse_thing_pretty_image(nox_thing*, nox_memfile*, char*);
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/binfile"
	"unsafe"
)

func PortTestResourceLightConstants() (*uint64, *uint64) {
	return (*uint64)(unsafe.Pointer(&C.qword_581450_9544)), (*uint64)(unsafe.Pointer(&C.qword_581450_9552))
}
func PortTestResourceClientParser(kind string, typ *client.ObjectType, f *binfile.MemFile, input unsafe.Pointer) bool {
	var file *C.nox_memfile
	if f != nil {
		file = (*C.nox_memfile)(f.C())
	}
	obj := (*C.nox_thing)(typ.C())
	s := (*C.char)(input)
	switch kind {
	case "direction":
		return bool(C.nox_parse_thing_light_dir(obj, file, s))
	case "penumbra":
		return bool(C.nox_parse_thing_light_penumbra(obj, file, s))
	case "update":
		return bool(C.nox_parse_thing_client_update(obj, file, s))
	case "image":
		return bool(C.nox_parse_thing_pretty_image(obj, file, s))
	}
	panic(kind)
}
