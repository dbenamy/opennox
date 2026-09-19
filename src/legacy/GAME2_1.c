#include "GAME2_1.h"
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "GAME3_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4_1.h"
#include "GAME5_2.h"
#include "client__draw__debugdraw.h"
#include "client__draw__staticdraw.h"
#include "client__drawable__drawable.h"
#include "client__gui__window.h"
#include "common__net_list.h"
#include "common__object__modifier.h"
#include "common__system__team.h"
#include "operators.h"

#include "client__gui__gamewin__gamewin.h"
#include "client__gui__gui_ctf.h"
#include "client__gui__guicon.h"
#include "client__gui__guiinput.h"
#include "client__gui__guiinv.h"
#include "client__gui__guijourn.h"
#include "client__gui__guimeter.h"
#include "client__gui__guiquit.h"
#include "client__gui__guirank.h"
#include "client__gui__guisave.h"
#include "client__gui__guispell.h"
#include "client__gui__guisumn.h"
#include "client__gui__guitrade.h"
#include "client__gui__tooltip.h"
#include "client__video__draw_common.h"

#include "client__system__ctrlevnt.h"
#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"
#include "input_common.h"

extern uint32_t dword_8531A0_2576;
extern uint32_t dword_8531A0_2572;
extern uint32_t dword_5d4594_1062552;
extern uint32_t dword_5d4594_1049844;
extern uint32_t dword_5d4594_1050008;
extern uint32_t dword_5d4594_1096272;
extern uint32_t dword_5d4594_1049516;
extern uint32_t dword_5d4594_1062520;
extern uint32_t dword_5d4594_1096276;
extern uint32_t dword_5d4594_1049976;
extern uint32_t dword_5d4594_1090284;
extern uint32_t dword_5d4594_1062484;
extern uint32_t dword_5d4594_1049992;
extern uint32_t dword_5d4594_1062556;
extern uint32_t dword_5d4594_1096264;
extern uint32_t dword_5d4594_1062564;
extern uint32_t dword_5d4594_1090280;
extern uint32_t dword_5d4594_1062560;
extern uint32_t dword_5d4594_1096280;
extern uint32_t dword_587000_145672;
extern uint32_t dword_5d4594_1064868;
extern uint32_t dword_5d4594_1049996;
extern uint32_t dword_5d4594_1062492;
extern uint32_t dword_5d4594_1062496;
extern uint32_t dword_5d4594_1096260;
extern uint32_t nox_client_gui_flag_1556112;
extern uint32_t dword_5d4594_1096284;
extern uint32_t dword_587000_145664;
extern uint32_t dword_5d4594_1096288;
extern uint32_t dword_5d4594_1047936;
extern uint32_t dword_5d4594_1062488;
extern uint32_t dword_5d4594_1062468;
extern uint32_t dword_5d4594_1064860;
extern uint32_t dword_5d4594_1096252;
extern uint32_t dword_5d4594_1064864;
extern uint32_t dword_587000_136184;
extern uint32_t dword_5d4594_1049532;
extern uint32_t dword_5d4594_1049484;
extern uint32_t dword_5d4594_1062516;
extern uint64_t qword_581450_9512;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_1090276;
extern uint32_t dword_5d4594_1062476;
extern uint32_t dword_5d4594_1047932;
extern uint32_t dword_5d4594_1049512;
extern uint32_t dword_5d4594_1062528;
extern uint32_t dword_5d4594_1062524;
extern uint32_t dword_5d4594_1064856;
extern uint32_t dword_5d4594_1049856;
extern uint32_t dword_5d4594_1049520;
extern uint32_t dword_5d4594_1049800_inventory_click_row_index;
extern uint32_t dword_5d4594_1062456;
extern uint32_t dword_5d4594_1063636;
extern uint32_t dword_5d4594_1049796_inventory_click_column_index;
extern uint32_t dword_5d4594_1049508;
extern nox_window* dword_5d4594_1090048;
extern uint32_t dword_5d4594_1049500;
extern uint32_t dword_5d4594_1062512;
extern uint32_t dword_5d4594_1049864;
extern uint32_t dword_5d4594_1062508;
extern uint32_t dword_5d4594_1049504;
extern uint32_t dword_5d4594_1090120;
extern uint32_t dword_5d4594_1063116;
extern uint32_t dword_5d4594_1062480;
extern uint32_t nox_player_netCode_85319C;
extern void* nox_xxx_aClosewoodengat_587000_133480;
extern int nox_win_width;
extern int nox_win_height;
extern uint32_t array_5D4594_1049872[9];

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_blue_2650684;
extern uint32_t nox_color_yellow_2589772;
extern uint32_t nox_color_violet_2598268;
extern uint32_t nox_color_black_2650656;

nox_window* nox_win_unk5 = 0;
nox_window* dword_5d4594_1062452 = 0;


obj_5D4594_2650668_t** ptr_5D4594_2650668 = 0;
const int ptr_5D4594_2650668_cap = 128;

nox_inventory_cell_t nox_client_inventory_grid_1050020[NOX_INVENTORY_CELLS_MAX] = {0};

//----- (00461460) --------------------------------------------------------
void nox_xxx_playerInitColors_461460(nox_playerInfo* pl) {
	int a1 = pl;
	int v1;     // eax
	char v2;    // cl
	char v3;    // dl
	int v4;     // eax
	char v5;    // cl
	char v6;    // dl
	int v7;     // eax
	char v8;    // cl
	char v9;    // dl
	int v10;    // eax
	char v11;   // cl
	char v12;   // dl

	v1 = nox_color_rgb_4344A0(*(uint8_t*)(a1 + 2253), *(uint8_t*)(a1 + 2254), *(uint8_t*)(a1 + 2255));
	v2 = *(uint8_t*)(a1 + 2257);
	v3 = *(uint8_t*)(a1 + 2256);
	*(uint32_t*)(a1 + 2296) = v1;
	v4 = nox_color_rgb_4344A0(v3, v2, *(uint8_t*)(a1 + 2258));
	v5 = *(uint8_t*)(a1 + 2260);
	v6 = *(uint8_t*)(a1 + 2259);
	*(uint32_t*)(a1 + 2292) = v4;
	v7 = nox_color_rgb_4344A0(v6, v5, *(uint8_t*)(a1 + 2261));
	v8 = *(uint8_t*)(a1 + 2263);
	v9 = *(uint8_t*)(a1 + 2262);
	*(uint32_t*)(a1 + 2300) = v7;
	v10 = nox_color_rgb_4344A0(v9, v8, *(uint8_t*)(a1 + 2264));
	v11 = *(uint8_t*)(a1 + 2266);
	v12 = *(uint8_t*)(a1 + 2265);
	*(uint32_t*)(a1 + 2304) = v10;
	*(uint32_t*)(a1 + 2308) = nox_color_rgb_4344A0(v12, v11, *(uint8_t*)(a1 + 2267));
	*(uint32_t*)(a1 + 2312) = nox_color_white_2523948;
}

//----- (00461520) --------------------------------------------------------
char* sub_461520() {
	char* result; // eax
	int i;        // esi

	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		nox_xxx_playerInitColors_461460(i);
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
	return result;
}

//----- (00469FA0) --------------------------------------------------------


//----- (0046A430) --------------------------------------------------------


//----- (0046A4A0) --------------------------------------------------------


//----- (0046A4B0) --------------------------------------------------------


//----- (0046A5D0) --------------------------------------------------------


//----- (0046A6A0) --------------------------------------------------------


//----- (0046A730) --------------------------------------------------------


//----- (0046A7E0) --------------------------------------------------------


//----- (0046A820) --------------------------------------------------------


//----- (0046A860) --------------------------------------------------------


//----- (00473920) --------------------------------------------------------
void nox_xxx____setargv_11_473920() { *getMemU32Ptr(0x5D4594, 1096520) = 1; }

//----- (00473930) --------------------------------------------------------
char* sub_473930() {
	char* result; // eax

	*getMemU32Ptr(0x5D4594, 1096456) = nox_xxx_gLoadAnim_42FA20("ConfusedBirdies");
	result = nox_xxx_gLoadAnim_42FA20("SphericalShieldAnim");
	*getMemU32Ptr(0x5D4594, 1096460) = result;
	return result;
}

//----- (00473960) --------------------------------------------------------
int sub_473960() {
	int result; // eax

	result = 0;
	*getMemU32Ptr(0x5D4594, 1096456) = 0;
	*getMemU32Ptr(0x5D4594, 1096460) = 0;
	return result;
}
