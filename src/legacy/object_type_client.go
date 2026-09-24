package legacy

/*
#include "GAME2_3.h"
#include "client__draw__debugdraw.h"


*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
)

func init() {
	client.RegisterThingParse("LIGHTDIRECTION", resourceClientField("direction"))
	client.RegisterThingParse("LIGHTPENUMBRA", resourceClientField("penumbra"))
	client.RegisterThingParse("CLIENTUPDATE", resourceClientField("update"))
	client.RegisterThingParse("PRETTYIMAGE", resourceClientField("image"))
	client.ThingDrawDefault = C.nox_thing_debug_draw
}

type nox_thing = C.nox_thing

func nox_xxx_getTTByNameSpriteMB_44CFC0(cstr *C.char) int {
	id := GoString(cstr)
	return GetClient().Cli().Things.IndByID(id)
}

func nox_get_thing_name(i int) *C.char {
	t := GetClient().Cli().Things.TypeByInd(i)
	if t == nil {
		return nil
	}
	return (*C.char)(unsafe.Pointer(t.Name))
}

func nox_get_thing_pretty_name(i int) *wchar2_t {
	t := GetClient().Cli().Things.TypeByInd(i)
	if t == nil {
		return nil
	}
	return (*wchar2_t)(unsafe.Pointer(t.PrettyName))
}
