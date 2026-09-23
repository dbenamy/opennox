#include "client__system__parsecmd.h"
#include "client__drawable__drawable.h"

void nullsub_68();
void* nox_xxx_aClosewoodengat_587000_133480 = 0;
uint32_t dword_5d4594_251572 = 0;
uint32_t dword_5d4594_831236 = 0;
uint32_t dword_5d4594_1309720 = 0;
uint32_t nox_xxx_lightningTarget_5d4594_2487908 = 0;
uint32_t dword_5d4594_2487556 = 0;
uint32_t nox_xxx_lightningTargetArrayIndex_5d4594_2487904 = 0;
void* dword_587000_127004 = 0;
uint32_t dword_587000_93156 = 0x1;
uint32_t dword_5d4594_1045432 = 0;
void* nox_alloc_tradeItems_2386496 = 0;
uint32_t dword_5d4594_816368 = 0;
void* dword_587000_93164 = 0;
void* nox_alloc_screenParticles_806044 = 0;
void* dword_587000_122852 = 0;
uint32_t dword_5d4594_1550916 = 0;
uint32_t dword_5d4594_2487540 = 0;
void* dword_587000_81128 = 0;
uint32_t dword_5d4594_2487560 = 0;
uint32_t dword_5d4594_2491616 = 0;
uint32_t dword_5d4594_2487248 = 0;
uint32_t dword_5d4594_2487532 = 0;
uint32_t dword_5d4594_588084 = 0;
uint32_t dword_5d4594_816372 = 0;
uint32_t nox_xxx_energyBoltTarget_5d4594_2487880 = 0;
nox_window* dword_5d4594_1090048 = 0;
uint32_t dword_5d4594_1090100 = 0;
uint32_t dword_5d4594_3835356 = 0;
uint32_t dword_587000_122848 = 0x1;
void* nox_alloc_monsterList_2386220 = 0;
uint32_t nox_xxx_lightningClosestTargetDistance_5d4594_2487912 = 0;
uint32_t dword_5d4594_1550912 = 0;
uint32_t dword_5d4594_2487620 = 0;
void* nox_alloc_tradeSession_2386492 = 0;
uint32_t dword_5d4594_3835352 = 0;
uint32_t dword_5d4594_2487624 = 0;
uint32_t dword_5d4594_2487576 = 0;
uint32_t dword_5d4594_3835348 = 0;
uint32_t dword_5d4594_2487652 = 0;
uint32_t dword_5d4594_3835388 = 0;
uint32_t dword_5d4594_2487676 = 0;
void* nox_alloc_spawn_2386216 = 0;
uint32_t dword_5d4594_2487672 = 0;
void* nox_alloc_magicEnt_1569668 = 0;
uint32_t dword_5d4594_2487564 = 0;
uint32_t dword_5d4594_2487584 = 0;
uint32_t nox_xxx_lightningOwner_5d4594_2487900 = 0;
uint32_t dword_5d4594_2487932 = 0;
uint32_t dword_5d4594_2487628 = 0;
uint32_t dword_5d4594_1045428 = 0;
uint32_t dword_5d4594_2487536 = 0;
uint32_t dword_587000_126996 = 0x1;
uint32_t dword_5d4594_2487568 = 0;
uint32_t dword_5d4594_2487580 = 0;
uint32_t dword_5d4594_3835392 = 0;
uint32_t dword_5d4594_3835372 = 0;
uint32_t dword_5d4594_2487884 = 0;
void* dword_5d4594_1548532 = 0;
uint32_t dword_5d4594_2489436 = 0;
uint32_t dword_5d4594_1045420 = 0;
uint32_t dword_5d4594_3835364 = 0;
uint32_t dword_5d4594_1549844 = 0;
uint32_t dword_5d4594_2487632 = 0;
uint32_t dword_5d4594_3835360 = 0;
uint32_t dword_5d4594_3835368 = 0;
uint32_t dword_5d4594_2487656 = 0;
uint32_t array_5D4594_1049872[9];

nox_screenParticle* nox_screenParticles_head = 0;
nox_screenParticle* dword_5d4594_806052 = 0;
void* dword_5d4594_805984 = 0;
nox_window* nox_win_unk1 = 0;
uint32_t dword_5d4594_831092 = 0;
nox_render_data_t* nox_draw_curDrawData_3799572 = 0;

void* dword_5d4594_830236 = 0;
void* dword_5d4594_830232 = 0;
uint32_t dword_5d4594_816376 = 0;
nox_window* nox_win_unk5 = 0;
nox_window* dword_5d4594_1062452 = 0;
nox_inventory_cell_t nox_client_inventory_grid_1050020[NOX_INVENTORY_CELLS_MAX] = {0};
uint8_t** nox_pixbuffer_rows_3798784 = 0;
void* nox_video_tileBuf_ptr_3798796 = 0;
void* nox_video_tileBuf_end_3798844 = 0;

int nox_server_gameSettingsUpdated; // ABI remains a four-byte int.
nox_gui_animation* nox_wnd_xxx_1309740 = 0;

// Shared storage retained after extension/listing translation-unit removal.
nox_window* dword_5d4594_1522616 = 0;
nox_window* dword_5d4594_1522620 = 0;
nox_window* dword_5d4594_1522624 = 0;
nox_window* dword_5d4594_1522628 = 0;
nox_gui_animation* nox_wnd_xxx_1522608 = 0;
void* nox_gui_itemAmount_item_1319256 = 0;
void* nox_gui_itemAmount_dialog_1319228 = 0;
nox_window* dword_5d4594_1321236 = 0;
nox_window* dword_5d4594_1321240 = 0;
nox_window* dword_5d4594_1321244 = 0;
nox_window* dword_5d4594_1321248 = 0;

// Shared catalog head; catalog algorithms are native.
nox_list_item_t nox_common_maplist = {0};
