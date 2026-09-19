#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "GAME3_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME5_2.h"
#include "client__draw__animdraw.h"
#include "client__drawable__drawable.h"
#include "client__gui__guisumn.h"
#include "client__gui__servopts__advserv.h"
#include "common__net_list.h"
#include "common__system__team.h"

#include "client__gui__guiinput.h"
#include "client__gui__servopts__objlst.h"
#include "client__gui__servopts__spelllst.h"
#include "client__gui__tooltip.h"
#include "client__gui__window.h"
#include "client__shell__inputcfg__inputcfg.h"
#include "client__shell__mainmenu.h"

#include "client__draw__plasma.h"
#include "client__drawable__update__charmup.h"
#include "client__drawable__update__fireball.h"
#include "client__system__ctrlevnt.h"

#include "client__shell__optsback.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__strman.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_8531A0_2572;
extern uint32_t dword_8531A0_2576;
extern uint32_t dword_587000_180476;
extern uint32_t dword_5d4594_1319260;
extern uint32_t dword_5d4594_1321196;
extern uint32_t dword_5d4594_1313880;
extern uint32_t dword_5d4594_1320944;
extern uint32_t dword_5d4594_1320948;
extern uint32_t dword_587000_183460;
extern uint32_t dword_5d4594_1316412;
extern uint32_t dword_5d4594_1522968;
extern uint32_t dword_5d4594_1321520;
extern uint32_t dword_5d4594_1319268;
extern uint32_t dword_5d4594_1321024;
extern uint32_t dword_5d4594_1320936;
extern uint32_t dword_587000_183456;
extern uint32_t dword_5d4594_1320988;
extern uint32_t dword_5d4594_1319248;
extern uint32_t dword_587000_180480;
extern uint32_t dword_5d4594_1319236;
extern uint32_t dword_5d4594_1320972;
extern uint32_t dword_5d4594_1321208;
extern uint32_t dword_5d4594_1321800;
extern uint32_t dword_5d4594_1321224;
extern uint32_t dword_5d4594_1320932;
extern uint32_t dword_5d4594_1321032;
extern uint32_t dword_5d4594_1319264;
extern uint32_t dword_5d4594_1321044;
extern uint32_t dword_5d4594_1319232;
extern nox_window* dword_5d4594_1321236;
extern nox_window* dword_5d4594_1321240;
extern nox_window* dword_5d4594_1321248;
extern nox_window* dword_5d4594_1321244;
nox_window* dword_5d4594_1522616 = 0;
nox_window* dword_5d4594_1522620 = 0;
nox_window* dword_5d4594_1522624 = 0;
nox_window* dword_5d4594_1522628 = 0;
extern uint32_t dword_5d4594_1320992;
extern uint32_t dword_5d4594_1321204;
extern uint32_t dword_5d4594_1316408;
extern uint64_t qword_581450_9512;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_1320968;
extern uint32_t dword_5d4594_1522632;
extern uint32_t dword_5d4594_1321252;
extern uint32_t dword_5d4594_1522612;
extern uint32_t dword_5d4594_1522604;
extern uint32_t dword_5d4594_1321232;
extern uint32_t dword_5d4594_1320964;
extern uint32_t dword_5d4594_1321228;
extern uint32_t dword_5d4594_1321040;
extern uint32_t dword_5d4594_1320940;
extern uint32_t nox_player_netCode_85319C;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_blue_2650684;
extern uint32_t nox_color_yellow_2589772;
extern uint32_t nox_color_violet_2598268;
extern uint32_t nox_color_black_2650656;

nox_gui_animation* nox_wnd_xxx_1522608 = 0;

void* nox_gui_itemAmount_item_1319256 = 0;
void* nox_gui_itemAmount_dialog_1319228 = 0;

// 487CF0: using guessed type void  nullsub_10(uint32_t);

// 487CA0: using guessed type void  nullsub_9(uint32_t);

//----- (004BF010) --------------------------------------------------------


//----- (004BF7E0) --------------------------------------------------------


//----- (004BF9F0) --------------------------------------------------------

// 4BFA4C: variable 'v10' is possibly undefined
// 4BFA4C: variable 'v11' is possibly undefined

//----- (004BFAD0) --------------------------------------------------------


//----- (004BFB70) --------------------------------------------------------


//----- (004BFBB0) --------------------------------------------------------


//----- (004BFBF0) --------------------------------------------------------


//----- (004BFC70) --------------------------------------------------------


//----- (004BFC90) --------------------------------------------------------


//----- (004BFCD0) --------------------------------------------------------


//----- (004BFD10) --------------------------------------------------------


//----- (004BFD30) --------------------------------------------------------


//----- (004C3390) --------------------------------------------------------


//----- (004C3410) --------------------------------------------------------


//----- (004C3460) --------------------------------------------------------


//----- (004C34A0) --------------------------------------------------------


//----- (004CA540) --------------------------------------------------------
int nox_xxx_sprite_4CA540(uint32_t* a1, int a2) {
	uint32_t* v2; // esi
	double v3;    // st7
	double v4;    // st6
	double v5;    // st5
	double v6;    // st4
	int v7;       // eax
	int v8;       // edi
	int v9;       // eax
	float v11;    // [esp+0h] [ebp-14h]
	float v12;    // [esp+0h] [ebp-14h]
	float v13;    // [esp+10h] [ebp-4h]
	float v14;    // [esp+1Ch] [ebp+8h]

	v2 = (uint32_t*)a2;
	v3 = 0.0;
	v4 = *(float*)(a2 + 468);
	v5 = *(float*)(a2 + 472);
	v6 = 0.0;
	v7 = gameFrame() - *(uint32_t*)(a2 + 316) + 1;
	do {
		--v7;
		v13 = -(v5 * *(float*)(a2 + 476));
		v4 = v4 - v4 * *(float*)(a2 + 476);
		v5 = v5 + v13;
		v3 = v3 + v4;
		v6 = v6 + v5;
	} while (v7);
	v14 = v6;
	v11 = (double)(int)v2[81] + v3;
	v8 = nox_float2int(v11);
	v12 = (double)(int)v2[82] + v14;
	v9 = nox_float2int(v12);
	if (v8 > 0 && v9 > 0 && v8 < 5888 && v9 < 5888) {
		nox_xxx_updateSpritePosition_49AA90(v2, v8, v9);
		if (sub_4992B0(*a1 + v2[3] - a1[4], v2[4] + a1[1] - a1[5])) {
			return 1;
		}
	}
	nox_xxx_spriteDeleteStatic_45A4E0_drawable((int)v2);
	return 0;
}

//----- (004CA650) --------------------------------------------------------
int sub_4CA650(int a1, int a2) {
	int v2;             // esi
	int v3;             // eax
	int v4;             // ebx
	int v5;             // edi
	int v6;             // eax
	int v7;             // ecx
	int v8;             // eax
	int v9;             // ebp
	int v10;            // eax
	int v11;            // edi
	unsigned short v12; // dx
	int v13;            // ebp
	int v14;            // eax
	int result;         // eax
	int v16;            // [esp+10h] [ebp-10h]
	int v17;            // [esp+28h] [ebp+8h]

	v2 = a2;
	v3 = *(uint32_t*)(a2 + 16);
	v4 = *(unsigned short*)(a2 + 434) - v3;
	v5 = *(unsigned short*)(a2 + 432) - *(uint32_t*)(a2 + 12);
	v6 = sub_48C6B0(v5, v4);
	v7 = v6;
	v8 = *(unsigned char*)(a2 + 443);
	v17 = v8;
	++v7;
	v9 = *(uint32_t*)(v2 + 12);
	v10 = v5 * v8 / v7;
	v11 = *(uint32_t*)(v2 + 16);
	v16 = v9 + v10;
	v12 = *(uint16_t*)(v2 + 432);
	v13 = v9 - v12;
	v14 = v11 + v4 * v17 / v7;
	if (v7 <= 10 ||
		v13 * (v16 - v12) + (v11 - *(unsigned short*)(v2 + 434)) * (v14 - *(unsigned short*)(v2 + 434)) < 0) {
		nox_xxx_spriteDeleteStatic_45A4E0_drawable(v2);
		result = 0;
	} else {
		nox_xxx_updateSpritePosition_49AA90((uint32_t*)v2, v16, v14);
		result = 1;
	}
	return result;
}
// 4CA67E: variable 'v6' is possibly undefined


nox_window* dword_5d4594_1321236 = 0;
nox_window* dword_5d4594_1321240 = 0;
nox_window* dword_5d4594_1321244 = 0;
nox_window* dword_5d4594_1321248 = 0;
