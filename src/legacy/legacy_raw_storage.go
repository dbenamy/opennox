package legacy

/*
#include "defs.h"
#include "server__script__internal.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// legacyGlobalStorage keeps the original field widths and dimensions in unmanaged C-heap memory.
// alloc.New uses calloc; this package-lifetime allocation is intentionally never freed.
type legacyGlobalStorage struct {
	nox_xxx_aClosewoodengat_587000_133480 unsafe.Pointer
	dword_5d4594_831236                   uint32
	dword_5d4594_1309720                  C.uint32_t
	dword_587000_127004                   unsafe.Pointer
	nox_alloc_tradeItems_2386496          unsafe.Pointer
	dword_587000_93164                    unsafe.Pointer
	nox_alloc_screenParticles_806044      unsafe.Pointer
	dword_587000_122852                   unsafe.Pointer
	dword_587000_81128                    unsafe.Pointer
	dword_5d4594_1090048                  *C.nox_window
	dword_5d4594_1090100                  C.uint32_t
	nox_alloc_monsterList_2386220         unsafe.Pointer
	nox_alloc_tradeSession_2386492        unsafe.Pointer
	nox_alloc_spawn_2386216               unsafe.Pointer
	nox_alloc_magicEnt_1569668            unsafe.Pointer
	dword_5d4594_1548532                  unsafe.Pointer
	array_5D4594_1049872                  [9]C.uint32_t
	nox_screenParticles_head              *C.nox_screenParticle
	dword_5d4594_806052                   *C.nox_screenParticle
	dword_5d4594_805984                   unsafe.Pointer
	nox_win_unk1                          *C.nox_window
	nox_draw_curDrawData_3799572          *C.nox_render_data_t
	dword_5d4594_830236                   unsafe.Pointer
	dword_5d4594_830232                   unsafe.Pointer
	nox_win_unk5                          *C.nox_window
	dword_5d4594_1062452                  *C.nox_window
	nox_client_inventory_grid_1050020     [C.NOX_INVENTORY_CELLS_MAX]C.nox_inventory_cell_t
	nox_pixbuffer_rows_3798784            **C.uint8_t
	nox_video_tileBuf_ptr_3798796         unsafe.Pointer
	nox_video_tileBuf_end_3798844         unsafe.Pointer
	nox_server_gameSettingsUpdated        C.int
	nox_wnd_xxx_1309740                   *C.nox_gui_animation
	dword_5d4594_1522616                  *C.nox_window
	dword_5d4594_1522620                  *C.nox_window
	dword_5d4594_1522624                  *C.nox_window
	dword_5d4594_1522628                  *C.nox_window
	nox_wnd_xxx_1522608                   *C.nox_gui_animation
	nox_gui_itemAmount_item_1319256       unsafe.Pointer
	nox_gui_itemAmount_dialog_1319228     unsafe.Pointer
	dword_5d4594_1321236                  *C.nox_window
	dword_5d4594_1321240                  *C.nox_window
	dword_5d4594_1321244                  *C.nox_window
	dword_5d4594_1321248                  *C.nox_window
	nox_common_maplist                    C.nox_list_item_t
}

var legacyGlobals, _ = alloc.New(legacyGlobalStorage{})

var byte_581450, _ = alloc.Make([]byte(nil), 23472)
var byte_587000, _ = alloc.Make([]byte(nil), 316820)
var byte_5D4594, _ = alloc.Make([]byte(nil), 2598284)
var byte_973CE0, _ = alloc.Make([]byte(nil), 568)
var byte_973F18, _ = alloc.Make([]byte(nil), 44881)
var byte_85B3FC, _ = alloc.Make([]byte(nil), 1029636)
var byte_852978, _ = alloc.Make([]byte(nil), 40)
var byte_973A20, _ = alloc.Make([]byte(nil), 704)
