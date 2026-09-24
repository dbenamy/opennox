package legacy

/*
#include "defs.h"
int sub_4DF9B0(void* a1, void* a2, void* a3, int a4);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/internal/netlist"
)

var (
	GetNetPlayerBufSize         func() int
	Nox_netlist_addToMsgListSrv func(ind ntype.PlayerInd, buf []byte) bool
)

//export nox_netlist_addToMsgListCli_40EBC0
func nox_netlist_addToMsgListCli_40EBC0(ind1, ind2 int, buf *C.uchar, sz int) int {
	return bool2int(GetServer().S().NetList.AddToMsgListCli(ntype.PlayerInd(ind1), netlist.Kind(ind2), unsafe.Slice((*byte)(unsafe.Pointer(buf)), sz)))
}

func Nox_xxx_netImportant_4E5770(a1 byte, a2 int) {
	reliableDeliver(a1, a2)
}
