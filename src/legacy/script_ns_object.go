package legacy

/*
#include "defs.h"

int nox_xxx_playDialogFile_44D900(unsigned char* a1, int a2);
int nox_xxx_inventoryServPlace_4F36F0(nox_object_t* a1p, nox_object_t* a2p, int a3, int a4);
void nox_xxx_playerCanCarryItem_513B00(nox_object_t* a1p, nox_object_t* a2p);
void nox_xxx_unitAdjustHP_4EE460(nox_object_t* unit, int dv);
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_xxx_inventoryServPlace_4F36F0 func(obj, it *server.Object, a3, a4 int) bool
)

func Nox_xxx_getObjectByScrName_4DA4F0(name string) *server.Object {
	return objectLookupByName(name)
}
func Nox_server_scriptMoveTo_5123C0(a1 *server.Object, a2 *server.Waypoint) {
	scriptBindingMove(a1, a2)
}
func Nox_xxx_playerCanCarryItem_513B00(a1 *server.Object, a2 *server.Object) {
	C.nox_xxx_playerCanCarryItem_513B00(asObjectC(a1), asObjectC(a2))
}

//export nox_xxx_inventoryServPlace_4F36F0
func nox_xxx_inventoryServPlace_4F36F0(a1 *nox_object_t, a2 *nox_object_t, a3 int, a4 int) int {
	if Nox_xxx_inventoryServPlace_4F36F0(asObjectS(a1), asObjectS(a2), a3, a4) {
		return 1
	}
	return 0
}

func Sub_516D00(a1 *server.Object) {
	monsterControlRevive(a1)
}
func Nox_xxx_netSendChat_528AC0(a1 *server.Object, a2 string, a3 uint16) {
	gameplayTextChat(a1, gameplayTextUnits(alloc.InternCString16(a2)), a3)
}
