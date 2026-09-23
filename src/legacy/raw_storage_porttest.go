//go:build porttest

package legacy

/*
#include "defs.h"
extern void* nox_xxx_aClosewoodengat_587000_133480;
extern uint32_t dword_5d4594_831236;
extern void* dword_5d4594_1309720;
extern void* dword_587000_127004;
extern void* nox_alloc_tradeItems_2386496;
extern void* dword_587000_93164;
extern void* nox_alloc_screenParticles_806044;
extern void* dword_587000_122852;
extern void* dword_587000_81128;
extern nox_window* dword_5d4594_1090048;
extern nox_window* dword_5d4594_1090100;
extern void* nox_alloc_monsterList_2386220;
extern void* nox_alloc_tradeSession_2386492;
extern void* nox_alloc_spawn_2386216;
extern void* nox_alloc_magicEnt_1569668;
extern void* dword_5d4594_1548532;
extern uint32_t array_5D4594_1049872[9];
extern nox_screenParticle* nox_screenParticles_head;
extern nox_screenParticle* dword_5d4594_806052;
extern void* dword_5d4594_805984;
extern nox_window* nox_win_unk1;
extern nox_render_data_t* nox_draw_curDrawData_3799572;
extern void* dword_5d4594_830236;
extern void* dword_5d4594_830232;
extern nox_window* nox_win_unk5;
extern nox_window* dword_5d4594_1062452;
extern nox_inventory_cell_t nox_client_inventory_grid_1050020[NOX_INVENTORY_CELLS_MAX];
extern uint8_t** nox_pixbuffer_rows_3798784;
extern void* nox_video_tileBuf_ptr_3798796;
extern void* nox_video_tileBuf_end_3798844;
extern int nox_server_gameSettingsUpdated;
extern nox_gui_animation* nox_wnd_xxx_1309740;
extern nox_window* dword_5d4594_1522616;
extern nox_window* dword_5d4594_1522620;
extern nox_window* dword_5d4594_1522624;
extern nox_window* dword_5d4594_1522628;
extern nox_gui_animation* nox_wnd_xxx_1522608;
extern void* nox_gui_itemAmount_item_1319256;
extern void* nox_gui_itemAmount_dialog_1319228;
extern nox_window* dword_5d4594_1321236;
extern nox_window* dword_5d4594_1321240;
extern nox_window* dword_5d4594_1321244;
extern nox_window* dword_5d4594_1321248;
extern nox_list_item_t nox_common_maplist;
extern unsigned char byte_581450[23472];
extern unsigned char byte_587000[316820];
extern unsigned char byte_5D4594[2598284];
extern unsigned char byte_973CE0[568];
extern unsigned char byte_973F18[44881];
extern unsigned char byte_85B3FC[1029636];
extern unsigned char byte_852978[40];
extern unsigned char byte_973A20[704];
*/
import "C"

import "unsafe"

type rawStorageContract struct {
	name                            string
	ptr                             unsafe.Pointer
	size, wantSize, align, blobBase uintptr
	alias                           func() unsafe.Pointer
	readWord                        func() uint32
	writeWord                       func(uint32)
	writePointer                    func(unsafe.Pointer)
}

var rawStorageContracts = []rawStorageContract{
	{name: "nox_xxx_aClosewoodengat_587000_133480", ptr: unsafe.Pointer(&C.nox_xxx_aClosewoodengat_587000_133480), size: unsafe.Sizeof(C.nox_xxx_aClosewoodengat_587000_133480), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_xxx_aClosewoodengat_587000_133480) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_xxx_aClosewoodengat_587000_133480))) }, writePointer: func(p unsafe.Pointer) { C.nox_xxx_aClosewoodengat_587000_133480 = (unsafe.Pointer)(p) }},
	{name: "dword_5d4594_831236", ptr: unsafe.Pointer(&C.dword_5d4594_831236), size: unsafe.Sizeof(C.dword_5d4594_831236), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_831236) }, readWord: func() uint32 { return uint32(C.dword_5d4594_831236) }, writeWord: func(v uint32) { C.dword_5d4594_831236 = C.uint32_t(v) }},
	{name: "dword_5d4594_1309720", ptr: unsafe.Pointer(&C.dword_5d4594_1309720), size: unsafe.Sizeof(C.dword_5d4594_1309720), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1309720) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1309720))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1309720 = (unsafe.Pointer)(p) }, writeWord: func(v uint32) { *(*uint32)(unsafe.Pointer(&C.dword_5d4594_1309720)) = v }},
	{name: "dword_587000_127004", ptr: unsafe.Pointer(&C.dword_587000_127004), size: unsafe.Sizeof(C.dword_587000_127004), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_587000_127004) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_587000_127004))) }, writePointer: func(p unsafe.Pointer) { C.dword_587000_127004 = (unsafe.Pointer)(p) }},
	{name: "nox_alloc_tradeItems_2386496", ptr: unsafe.Pointer(&C.nox_alloc_tradeItems_2386496), size: unsafe.Sizeof(C.nox_alloc_tradeItems_2386496), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_alloc_tradeItems_2386496) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_alloc_tradeItems_2386496))) }, writePointer: func(p unsafe.Pointer) { C.nox_alloc_tradeItems_2386496 = (unsafe.Pointer)(p) }},
	{name: "dword_587000_93164", ptr: unsafe.Pointer(&C.dword_587000_93164), size: unsafe.Sizeof(C.dword_587000_93164), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_587000_93164) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_587000_93164))) }, writePointer: func(p unsafe.Pointer) { C.dword_587000_93164 = (unsafe.Pointer)(p) }},
	{name: "nox_alloc_screenParticles_806044", ptr: unsafe.Pointer(&C.nox_alloc_screenParticles_806044), size: unsafe.Sizeof(C.nox_alloc_screenParticles_806044), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_alloc_screenParticles_806044) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_alloc_screenParticles_806044))) }, writePointer: func(p unsafe.Pointer) { C.nox_alloc_screenParticles_806044 = (unsafe.Pointer)(p) }},
	{name: "dword_587000_122852", ptr: unsafe.Pointer(&C.dword_587000_122852), size: unsafe.Sizeof(C.dword_587000_122852), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_587000_122852) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_587000_122852))) }, writePointer: func(p unsafe.Pointer) { C.dword_587000_122852 = (unsafe.Pointer)(p) }},
	{name: "dword_587000_81128", ptr: unsafe.Pointer(&C.dword_587000_81128), size: unsafe.Sizeof(C.dword_587000_81128), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_587000_81128) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_587000_81128))) }, writePointer: func(p unsafe.Pointer) { C.dword_587000_81128 = (unsafe.Pointer)(p) }},
	{name: "dword_5d4594_1090048", ptr: unsafe.Pointer(&C.dword_5d4594_1090048), size: unsafe.Sizeof(C.dword_5d4594_1090048), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1090048) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1090048))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1090048 = (*C.nox_window)(p) }},
	{name: "dword_5d4594_1090100", ptr: unsafe.Pointer(&C.dword_5d4594_1090100), size: unsafe.Sizeof(C.dword_5d4594_1090100), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1090100) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1090100))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1090100 = (*C.nox_window)(p) }, writeWord: func(v uint32) { *(*uint32)(unsafe.Pointer(&C.dword_5d4594_1090100)) = v }},
	{name: "nox_alloc_monsterList_2386220", ptr: unsafe.Pointer(&C.nox_alloc_monsterList_2386220), size: unsafe.Sizeof(C.nox_alloc_monsterList_2386220), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_alloc_monsterList_2386220) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_alloc_monsterList_2386220))) }, writePointer: func(p unsafe.Pointer) { C.nox_alloc_monsterList_2386220 = (unsafe.Pointer)(p) }},
	{name: "nox_alloc_tradeSession_2386492", ptr: unsafe.Pointer(&C.nox_alloc_tradeSession_2386492), size: unsafe.Sizeof(C.nox_alloc_tradeSession_2386492), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_alloc_tradeSession_2386492) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_alloc_tradeSession_2386492))) }, writePointer: func(p unsafe.Pointer) { C.nox_alloc_tradeSession_2386492 = (unsafe.Pointer)(p) }},
	{name: "nox_alloc_spawn_2386216", ptr: unsafe.Pointer(&C.nox_alloc_spawn_2386216), size: unsafe.Sizeof(C.nox_alloc_spawn_2386216), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_alloc_spawn_2386216) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_alloc_spawn_2386216))) }, writePointer: func(p unsafe.Pointer) { C.nox_alloc_spawn_2386216 = (unsafe.Pointer)(p) }},
	{name: "nox_alloc_magicEnt_1569668", ptr: unsafe.Pointer(&C.nox_alloc_magicEnt_1569668), size: unsafe.Sizeof(C.nox_alloc_magicEnt_1569668), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_alloc_magicEnt_1569668) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_alloc_magicEnt_1569668))) }, writePointer: func(p unsafe.Pointer) { C.nox_alloc_magicEnt_1569668 = (unsafe.Pointer)(p) }},
	{name: "dword_5d4594_1548532", ptr: unsafe.Pointer(&C.dword_5d4594_1548532), size: unsafe.Sizeof(C.dword_5d4594_1548532), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1548532) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1548532))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1548532 = (unsafe.Pointer)(p) }},
	{name: "array_5D4594_1049872", ptr: unsafe.Pointer(&C.array_5D4594_1049872), size: unsafe.Sizeof(C.array_5D4594_1049872), wantSize: 36, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.array_5D4594_1049872) }},
	{name: "nox_screenParticles_head", ptr: unsafe.Pointer(&C.nox_screenParticles_head), size: unsafe.Sizeof(C.nox_screenParticles_head), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_screenParticles_head) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_screenParticles_head))) }, writePointer: func(p unsafe.Pointer) { C.nox_screenParticles_head = (*C.nox_screenParticle)(p) }},
	{name: "dword_5d4594_806052", ptr: unsafe.Pointer(&C.dword_5d4594_806052), size: unsafe.Sizeof(C.dword_5d4594_806052), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_806052) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_806052))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_806052 = (*C.nox_screenParticle)(p) }},
	{name: "dword_5d4594_805984", ptr: unsafe.Pointer(&C.dword_5d4594_805984), size: unsafe.Sizeof(C.dword_5d4594_805984), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_805984) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_805984))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_805984 = (unsafe.Pointer)(p) }},
	{name: "nox_win_unk1", ptr: unsafe.Pointer(&C.nox_win_unk1), size: unsafe.Sizeof(C.nox_win_unk1), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_win_unk1) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_win_unk1))) }, writePointer: func(p unsafe.Pointer) { C.nox_win_unk1 = (*C.nox_window)(p) }},
	{name: "nox_draw_curDrawData_3799572", ptr: unsafe.Pointer(&C.nox_draw_curDrawData_3799572), size: unsafe.Sizeof(C.nox_draw_curDrawData_3799572), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_draw_curDrawData_3799572) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_draw_curDrawData_3799572))) }, writePointer: func(p unsafe.Pointer) { C.nox_draw_curDrawData_3799572 = (*C.nox_render_data_t)(p) }},
	{name: "dword_5d4594_830236", ptr: unsafe.Pointer(&C.dword_5d4594_830236), size: unsafe.Sizeof(C.dword_5d4594_830236), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_830236) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_830236))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_830236 = (unsafe.Pointer)(p) }},
	{name: "dword_5d4594_830232", ptr: unsafe.Pointer(&C.dword_5d4594_830232), size: unsafe.Sizeof(C.dword_5d4594_830232), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_830232) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_830232))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_830232 = (unsafe.Pointer)(p) }},
	{name: "nox_win_unk5", ptr: unsafe.Pointer(&C.nox_win_unk5), size: unsafe.Sizeof(C.nox_win_unk5), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_win_unk5) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_win_unk5))) }, writePointer: func(p unsafe.Pointer) { C.nox_win_unk5 = (*C.nox_window)(p) }},
	{name: "dword_5d4594_1062452", ptr: unsafe.Pointer(&C.dword_5d4594_1062452), size: unsafe.Sizeof(C.dword_5d4594_1062452), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1062452) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1062452))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1062452 = (*C.nox_window)(p) }},
	{name: "nox_client_inventory_grid_1050020", ptr: unsafe.Pointer(&C.nox_client_inventory_grid_1050020), size: unsafe.Sizeof(C.nox_client_inventory_grid_1050020), wantSize: 12432, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_client_inventory_grid_1050020) }},
	{name: "nox_pixbuffer_rows_3798784", ptr: unsafe.Pointer(&C.nox_pixbuffer_rows_3798784), size: unsafe.Sizeof(C.nox_pixbuffer_rows_3798784), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_pixbuffer_rows_3798784) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_pixbuffer_rows_3798784))) }, writePointer: func(p unsafe.Pointer) { C.nox_pixbuffer_rows_3798784 = (**C.uint8_t)(p) }},
	{name: "nox_video_tileBuf_ptr_3798796", ptr: unsafe.Pointer(&C.nox_video_tileBuf_ptr_3798796), size: unsafe.Sizeof(C.nox_video_tileBuf_ptr_3798796), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_video_tileBuf_ptr_3798796) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_video_tileBuf_ptr_3798796))) }, writePointer: func(p unsafe.Pointer) { C.nox_video_tileBuf_ptr_3798796 = (unsafe.Pointer)(p) }},
	{name: "nox_video_tileBuf_end_3798844", ptr: unsafe.Pointer(&C.nox_video_tileBuf_end_3798844), size: unsafe.Sizeof(C.nox_video_tileBuf_end_3798844), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_video_tileBuf_end_3798844) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_video_tileBuf_end_3798844))) }, writePointer: func(p unsafe.Pointer) { C.nox_video_tileBuf_end_3798844 = (unsafe.Pointer)(p) }},
	{name: "nox_server_gameSettingsUpdated", ptr: unsafe.Pointer(&C.nox_server_gameSettingsUpdated), size: unsafe.Sizeof(C.nox_server_gameSettingsUpdated), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_server_gameSettingsUpdated) }, readWord: func() uint32 { return uint32(C.nox_server_gameSettingsUpdated) }, writeWord: func(v uint32) { C.nox_server_gameSettingsUpdated = C.int(v) }},
	{name: "nox_wnd_xxx_1309740", ptr: unsafe.Pointer(&C.nox_wnd_xxx_1309740), size: unsafe.Sizeof(C.nox_wnd_xxx_1309740), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_wnd_xxx_1309740) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_wnd_xxx_1309740))) }, writePointer: func(p unsafe.Pointer) { C.nox_wnd_xxx_1309740 = (*C.nox_gui_animation)(p) }},
	{name: "dword_5d4594_1522616", ptr: unsafe.Pointer(&C.dword_5d4594_1522616), size: unsafe.Sizeof(C.dword_5d4594_1522616), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1522616) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1522616))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1522616 = (*C.nox_window)(p) }},
	{name: "dword_5d4594_1522620", ptr: unsafe.Pointer(&C.dword_5d4594_1522620), size: unsafe.Sizeof(C.dword_5d4594_1522620), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1522620) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1522620))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1522620 = (*C.nox_window)(p) }},
	{name: "dword_5d4594_1522624", ptr: unsafe.Pointer(&C.dword_5d4594_1522624), size: unsafe.Sizeof(C.dword_5d4594_1522624), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1522624) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1522624))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1522624 = (*C.nox_window)(p) }},
	{name: "dword_5d4594_1522628", ptr: unsafe.Pointer(&C.dword_5d4594_1522628), size: unsafe.Sizeof(C.dword_5d4594_1522628), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1522628) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1522628))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1522628 = (*C.nox_window)(p) }},
	{name: "nox_wnd_xxx_1522608", ptr: unsafe.Pointer(&C.nox_wnd_xxx_1522608), size: unsafe.Sizeof(C.nox_wnd_xxx_1522608), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_wnd_xxx_1522608) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_wnd_xxx_1522608))) }, writePointer: func(p unsafe.Pointer) { C.nox_wnd_xxx_1522608 = (*C.nox_gui_animation)(p) }},
	{name: "nox_gui_itemAmount_item_1319256", ptr: unsafe.Pointer(&C.nox_gui_itemAmount_item_1319256), size: unsafe.Sizeof(C.nox_gui_itemAmount_item_1319256), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_gui_itemAmount_item_1319256) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_gui_itemAmount_item_1319256))) }, writePointer: func(p unsafe.Pointer) { C.nox_gui_itemAmount_item_1319256 = (unsafe.Pointer)(p) }},
	{name: "nox_gui_itemAmount_dialog_1319228", ptr: unsafe.Pointer(&C.nox_gui_itemAmount_dialog_1319228), size: unsafe.Sizeof(C.nox_gui_itemAmount_dialog_1319228), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_gui_itemAmount_dialog_1319228) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.nox_gui_itemAmount_dialog_1319228))) }, writePointer: func(p unsafe.Pointer) { C.nox_gui_itemAmount_dialog_1319228 = (unsafe.Pointer)(p) }},
	{name: "dword_5d4594_1321236", ptr: unsafe.Pointer(&C.dword_5d4594_1321236), size: unsafe.Sizeof(C.dword_5d4594_1321236), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1321236) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1321236))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1321236 = (*C.nox_window)(p) }},
	{name: "dword_5d4594_1321240", ptr: unsafe.Pointer(&C.dword_5d4594_1321240), size: unsafe.Sizeof(C.dword_5d4594_1321240), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1321240) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1321240))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1321240 = (*C.nox_window)(p) }},
	{name: "dword_5d4594_1321244", ptr: unsafe.Pointer(&C.dword_5d4594_1321244), size: unsafe.Sizeof(C.dword_5d4594_1321244), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1321244) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1321244))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1321244 = (*C.nox_window)(p) }},
	{name: "dword_5d4594_1321248", ptr: unsafe.Pointer(&C.dword_5d4594_1321248), size: unsafe.Sizeof(C.dword_5d4594_1321248), wantSize: 4, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.dword_5d4594_1321248) }, readWord: func() uint32 { return uint32(uintptr(unsafe.Pointer(C.dword_5d4594_1321248))) }, writePointer: func(p unsafe.Pointer) { C.dword_5d4594_1321248 = (*C.nox_window)(p) }},
	{name: "nox_common_maplist", ptr: unsafe.Pointer(&C.nox_common_maplist), size: unsafe.Sizeof(C.nox_common_maplist), wantSize: 12, align: 4, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.nox_common_maplist) }},
	{name: "byte_581450", ptr: unsafe.Pointer(&C.byte_581450[0]), size: unsafe.Sizeof(C.byte_581450), wantSize: 23472, align: 4, blobBase: 0x581450, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.byte_581450[0]) }},
	{name: "byte_587000", ptr: unsafe.Pointer(&C.byte_587000[0]), size: unsafe.Sizeof(C.byte_587000), wantSize: 316820, align: 4, blobBase: 0x587000, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.byte_587000[0]) }},
	{name: "byte_5D4594", ptr: unsafe.Pointer(&C.byte_5D4594[0]), size: unsafe.Sizeof(C.byte_5D4594), wantSize: 2598284, align: 4, blobBase: 0x5d4594, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.byte_5D4594[0]) }},
	{name: "byte_973CE0", ptr: unsafe.Pointer(&C.byte_973CE0[0]), size: unsafe.Sizeof(C.byte_973CE0), wantSize: 568, align: 4, blobBase: 0x973ce0, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.byte_973CE0[0]) }},
	{name: "byte_973F18", ptr: unsafe.Pointer(&C.byte_973F18[0]), size: unsafe.Sizeof(C.byte_973F18), wantSize: 44881, align: 4, blobBase: 0x973f18, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.byte_973F18[0]) }},
	{name: "byte_85B3FC", ptr: unsafe.Pointer(&C.byte_85B3FC[0]), size: unsafe.Sizeof(C.byte_85B3FC), wantSize: 1029636, align: 4, blobBase: 0x85b3fc, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.byte_85B3FC[0]) }},
	{name: "byte_852978", ptr: unsafe.Pointer(&C.byte_852978[0]), size: unsafe.Sizeof(C.byte_852978), wantSize: 40, align: 4, blobBase: 0x852978, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.byte_852978[0]) }},
	{name: "byte_973A20", ptr: unsafe.Pointer(&C.byte_973A20[0]), size: unsafe.Sizeof(C.byte_973A20), wantSize: 704, align: 4, blobBase: 0x973a20, alias: func() unsafe.Pointer { return unsafe.Pointer(&C.byte_973A20[0]) }},
}

// Record static initialization before package init functions configure the game.
var rawStorageInitialZero = func() []bool {
	out := make([]bool, len(rawStorageContracts))
	for i, c := range rawStorageContracts {
		out[i] = true
		for _, v := range unsafe.Slice((*byte)(c.ptr), c.size) {
			if v != 0 {
				out[i] = false
				break
			}
		}
	}
	return out
}()

// Typed aliases supplement the byte contract with the inventory/list ABI.
func rawStorageLayout() []uintptr {
	var cell C.nox_inventory_cell_t
	var list C.nox_list_item_t
	return []uintptr{unsafe.Sizeof(cell), C.NOX_INVENTORY_CELLS_MAX, unsafe.Offsetof(cell.field_0), unsafe.Offsetof(cell.field_4), unsafe.Offsetof(cell.data_4), unsafe.Offsetof(cell.field_128), unsafe.Offsetof(cell.field_132), unsafe.Offsetof(cell.field_136), unsafe.Offsetof(cell.field_140), unsafe.Offsetof(cell.field_141), unsafe.Offsetof(cell.field_142), unsafe.Offsetof(cell.field_143), unsafe.Offsetof(cell.field_144), unsafe.Sizeof(list), unsafe.Offsetof(list.field_0), unsafe.Offsetof(list.field_1), unsafe.Offsetof(list.field_2)}
}

func rawStorageTypedEquipment(index int, v uint32) uint32 {
	C.array_5D4594_1049872[index] = C.uint32_t(v)
	return uint32(C.array_5D4594_1049872[index])
}

func rawStorageTypedInventory(index int, p unsafe.Pointer, code, tail uint32, count byte) {
	c := &C.nox_client_inventory_grid_1050020[index]
	c.field_0 = (*C.nox_drawable)(p)
	c.field_4 = C.uint32_t(code)
	c.field_140 = C.uint8_t(count)
	c.field_144 = C.uint32_t(tail)
}

func rawStorageTypedList(p unsafe.Pointer) {
	C.nox_common_maplist.field_0 = (*C.nox_list_item_t)(p)
	C.nox_common_maplist.field_1 = (*C.nox_list_item_t)(p)
	C.nox_common_maplist.field_2 = (*C.nox_list_item_t)(p)
}
