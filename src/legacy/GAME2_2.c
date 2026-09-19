#include <math.h>

#include "compat.h"

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
#include "GAME3_3.h"
#include "GAME4_1.h"
#include "GAME5_2.h"
#include "client__draw__debugdraw.h"
#include "client__draw__staticdraw.h"
#include "client__drawable__drawable.h"

#include "client__gui__gadgets__listbox.h"
#include "client__gui__guibook.h"
#include "client__gui__guicon.h"
#include "client__gui__guishop.h"
#include "client__gui__servopts__guiserv.h"
#include "client__gui__tooltip.h"
#include "client__gui__window.h"

#include "client__video__draw_common.h"

#include "common/fs/nox_fs.h"
#include "common__binfile.h"
#include "common__magic__speltree.h"
#include "common__net_list.h"
#include "common__strman.h"
#include "input.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_8531A0_2576;
extern uint32_t dword_5d4594_1096636;
extern uint32_t dword_5d4594_1123520;
extern uint32_t dword_5d4594_3804684;
extern uint32_t nox_xxx_waypointCounterMB_587000_154948;
extern uint32_t dword_5d4594_3798800;
extern uint64_t qword_581450_9552;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_3798808;
extern uint32_t dword_5d4594_1193384;
extern uint32_t dword_5d4594_1193360;
extern uint32_t dword_5d4594_3798836;
extern uint32_t dword_5d4594_251572;
extern uint32_t dword_5d4594_3798804;
extern uint32_t dword_5d4594_3798820;
extern uint32_t dword_5d4594_3798824;
extern uint32_t dword_5d4594_3798840;
extern void* dword_5d4594_1123524;
extern uint32_t dword_5d4594_1193380;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_yellow_2589772;

extern nox_render_data_t* nox_draw_curDrawData_3799572;


extern uint32_t nox_tile_def_cnt;
extern nox_tileDef_t nox_tile_defs_arr[176];

uint32_t dword_5d4594_1193156 = 0;
uint8_t** nox_pixbuffer_rows_3798784 = 0;

void* dword_5d4594_1096640 = 0;


void* nox_client_spriteUnderCursorXxx_1096644 = 0;
uint32_t nox_client_highResFloors_154952 = 1;
void* nox_video_tileBuf_ptr_3798796 = 0;
void* nox_video_tileBuf_end_3798844 = 0;

//----- (00476F40) --------------------------------------------------------
unsigned int nox_xxx_packetGetMarshall_476F40() {
	unsigned int result; // eax

	if (dword_5d4594_1096640) {
		result = nox_xxx_netGetUnitCodeCli_578B00(*(int*)&dword_5d4594_1096640);
	} else {
		result = 0;
	}
	return result;
}

//----- (00476FA0) --------------------------------------------------------
void nox_xxx_clientEnumHover_476FA0() {
	int4 v2; // [esp+10h] [ebp-10h]

	if (!*getMemU32Ptr(0x5D4594, 1096632)) {
		*getMemU32Ptr(0x5D4594, 1096632) = nox_xxx_getNameId_4E3AA0("Glyph");
	}
	nox_point mpos = nox_client_getMousePos_4309F0();
	sub_473970(&mpos, &mpos);
	dword_5d4594_1096640 = 0;
	nox_client_spriteUnderCursorXxx_1096644 = 0;
	*getMemU32Ptr(0x5D4594, 1096628) = 0;
	v2.field_0 = mpos.x - 96;
	v2.field_8 = mpos.x + 96;
	v2.field_C = mpos.y + 96;
	v2.field_4 = mpos.y - 96;
	dword_5d4594_1096636 = 0;
	nox_xxx_forEachSprite_49AB00(&v2, nox_xxx_clientOnCursorHover_477050, &mpos);
}

//----- (00477050) --------------------------------------------------------
int nox_xxx_client_4984B0_drawable(nox_drawable* dr);
void nox_xxx_clientOnCursorHover_477050(int arg0, int a2) {
	int v2;    // esi
	int v3;    // eax
	int v4;    // eax
	int v5;    // ecx
	char* v6;  // eax
	int v7;    // eax
	float v8;  // ebp
	int v9;    // edi
	int v10;   // ebx
	int v11;   // eax
	int v12;   // edx
	int v13;   // edi
	int v14;   // ebx
	float v15; // eax
	int v16;   // edx
	int v17;   // eax
	int v18;   // edi
	int v19;   // ebx
	int v20;   // eax
	int v21;   // ecx
	int v22;   // edi
	float v23; // [esp+0h] [ebp-24h]
	float v24; // [esp+0h] [ebp-24h]
	float v25; // [esp+0h] [ebp-24h]
	float v26; // [esp+0h] [ebp-24h]
	float2 a3; // [esp+14h] [ebp-10h]
	float2 a1; // [esp+1Ch] [ebp-8h]
	int v29;   // [esp+28h] [ebp+4h]

	if (!*getMemU32Ptr(0x5D4594, 1096648)) {
		*getMemU32Ptr(0x5D4594, 1096648) = nox_xxx_getTTByNameSpriteMB_44CFC0("Polyp");
	}
	v2 = arg0;
	if (arg0 == *getMemU32Ptr(0x852978, 8)) {
		return;
	}
	v3 = *(uint32_t*)(arg0 + 120);
	if (!((v3 & 0x8000) == 0 && (!nox_client_drawable_testBuff_4356C0(arg0, 0) ||
								 nox_client_drawable_testBuff_4356C0(*getMemIntPtr(0x852978, 8), 21)))) {
		return;
	}
	v4 = *(uint32_t*)(arg0 + 112);
	if (!(!(v4 & 2) || (v5 = *(uint32_t*)(arg0 + 116), !(v5 & 0x4000)))) {
		return;
	}
	if (!(v4 & 0x80400206 || *(uint32_t*)(arg0 + 108) == *getMemU32Ptr(0x5D4594, 1096648))) {
		return;
	}
	if (!nox_xxx_client_4984B0_drawable(arg0)) {
		return;
	}
	if (!(!(*(uint8_t*)(arg0 + 112) & 4) ||
		  (v6 = nox_common_playerInfoGetByID_417040(*(uint32_t*)(arg0 + 128))) != 0 && !(v6[3680] & 1))) {
		return;
	}
	v7 = *(uint32_t*)(arg0 + 112);
	if (!((!(v7 & 0x400000) || (*(uint8_t*)(arg0 + 116) & 0x80)) && (!(v7 & 2) || *(uint32_t*)(arg0 + 276) != 10))) {
		return;
	}
	v23 = (double)*(int*)(arg0 + 16) - *(float*)(arg0 + 100) - (double)*(short*)(arg0 + 104);
	v29 = nox_float2int(v23);
	v24 = (double)*(int*)(v2 + 16) - *(float*)(v2 + 96) - (double)*(short*)(v2 + 104);
	v8 = COERCE_FLOAT(nox_float2int(v24));
	a3.field_0 = v8;
	if (*(uint32_t*)(v2 + 44) == 2) {
		v25 = *(float*)(v2 + 48) * *(float*)(v2 + 48);
		LODWORD(a3.field_0) = nox_float2int(v25);
		v17 = nox_float2int(*(float*)(v2 + 48));
		v18 = *(uint32_t*)(v2 + 12);
		v19 = *(uint32_t*)(v2 + 12) - v17;
		v20 = v18 + nox_float2int(*(float*)(v2 + 48));
		v21 = *(uint32_t*)(a2 + 4);
		if (v21 <= SLODWORD(v8)) {
			v8 = *(float*)&v29;
			if (v21 >= v29) {
				if (*(int*)a2 <= v19 || *(int*)a2 >= v20) {
					return;
				}
				goto LABEL_38;
			}
		}
		v15 = a3.field_0;
		v16 = (*(uint32_t*)a2 - v18) * (*(uint32_t*)a2 - v18) + (v21 - LODWORD(v8)) * (v21 - LODWORD(v8));
	} else {
		if (*(uint32_t*)(v2 + 44) != 3) {
			return;
		}
		a1.field_0 = (double)*(int*)(v2 + 12);
		a1.field_4 = (double)SLODWORD(a3.field_0);
		a3.field_0 = (double)*(int*)a2;
		a3.field_4 = (double)*(int*)(a2 + 4);
		if (nox_xxx_map_57B850(&a1, (float*)(v2 + 44), &a3) ||
			(a1.field_4 = (double)v29, nox_xxx_map_57B850(&a1, (float*)(v2 + 44), &a3)) ||
			(v9 = *(uint32_t*)(v2 + 12) + nox_float2int(*(float*)(v2 + 72)),
			 v10 = v29 + nox_float2int(*(float*)(v2 + 76)), v11 = LODWORD(v8) + nox_float2int(*(float*)(v2 + 76)),
			 *(uint32_t*)a2 > v9) &&
				*(uint32_t*)a2 < *(int*)(v2 + 12) && (v12 = *(uint32_t*)(a2 + 4), v12 > v10) && v12 < v11) {
			goto LABEL_38;
		}
		v13 = *(uint32_t*)(v2 + 12) + nox_float2int(*(float*)(v2 + 80));
		v14 = v29 + nox_float2int(*(float*)(v2 + 84));
		LODWORD(v15) = LODWORD(v8) + nox_float2int(*(float*)(v2 + 84));
		if (*(int*)a2 < *(int*)(v2 + 12)) {
			return;
		}
		if (*(int*)a2 >= v13) {
			return;
		}
		v16 = *(uint32_t*)(a2 + 4);
		if (v16 <= v14) {
			return;
		}
	}
	if (v16 >= SLODWORD(v15)) {
		return;
	}
LABEL_38:
	v26 = (double)*(short*)(v2 + 104) + (double)*(int*)(v2 + 16) + *(float*)(v2 + 96);
	v22 = nox_float2int(v26);
	if (v22 > *getMemIntPtr(0x5D4594, 1096628)) {
		*getMemU32Ptr(0x5D4594, 1096628) = v22;
		dword_5d4594_1096640 = v2;
	}
	if (v2 != *getMemU32Ptr(0x852978, 8) && v22 > *(int*)&dword_5d4594_1096636 && nox_xxx_client_57B400(v2)) {
		if (dword_8531A0_2576 && *(uint8_t*)(dword_8531A0_2576 + 2251) == 1 &&
			*(uint32_t*)(v2 + 108) == *getMemU32Ptr(0x5D4594, 1096632)) {
			if (!nox_client_spriteUnderCursorXxx_1096644) {
				nox_client_spriteUnderCursorXxx_1096644 = v2;
				dword_5d4594_1096636 = 0;
			}
		} else {
			dword_5d4594_1096636 = v22;
			nox_client_spriteUnderCursorXxx_1096644 = v2;
		}
	}
}

//----- (00477600) --------------------------------------------------------
int nox_xxx_guiCursor_477600() { return *getMemU32Ptr(0x5D4594, 1096672); }

//----- (00479950) --------------------------------------------------------
int sub_479950() {
	void* v2; // [esp+0h] [ebp-4h]

	if (wndIsShown_nox_xxx_wndIsShown_46ACC0(*(int*)&dword_5d4594_1123524) == 1) {
		return 0;
	}
	*getMemU8Ptr(0x5D4594, 1123516) = 0;
	LOWORD(v2) = 720;
	BYTE2(v2) = 0;
	nox_netlist_addToMsgListCli_40EBC0(31, 0, &v2, 3);
	return 1;
}

//----- (004799A0) --------------------------------------------------------
int sub_4799A0() {
	int result;    // eax
	uint32_t* v1;  // edi
	uint32_t* v2;  // ebp
	uint32_t* v3;  // eax
	int v4;        // esi
	uint32_t* v5;  // ebx
	char* v6;      // eax
	uint32_t* v7;  // eax
	char* v8;      // [esp-18h] [ebp-1Ch]
	char* v9;      // [esp-14h] [ebp-18h]
	uint32_t* v10; // [esp+0h] [ebp-4h]

	*getMemU32Ptr(0x5D4594, 1107052) = nox_color_rgb_4344A0(240, 128, 64);
	result = nox_new_window_from_file("Dialog.wnd", nox_xxx_guiDialog_479B00);
	dword_5d4594_1123524 = result;
	if (result) {
		nox_xxx_wndSetWindowProc_46B300(result, sub_479BE0);
		nox_xxx_wndSetDrawFn_46B340(*(int*)&dword_5d4594_1123524, sub_479CB0);
		nox_gui_winSetFunc96_46B070(*(int*)&dword_5d4594_1123524, sub_479D00);
		v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3904);
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3903);
		v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3902);
		v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3901);
		v4 = (int)v3;
		v5 = (uint32_t*)v3[8];
		v9 = nox_xxx_gLoadImg_42F970("UISliderLit");
		v8 = nox_xxx_gLoadImg_42F970("UISliderLit");
		v6 = nox_xxx_gLoadImg_42F970("UISlider");
		sub_4B5700((int)v1, 0, 0, (int)v6, (int)v8, (int)v9);
		nox_xxx_wnd_46B280((int)v1, v4);
		nox_xxx_wnd_46B280((int)v2, v4);
		nox_xxx_wnd_46B280((int)v10, v4);
		v5[9] = v1;
		v5[7] = v2;
		v5[8] = v10;
		v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3906);
		nox_xxx_wndSetDrawFn_46B340((int)v7, sub_479C40);
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1123524, 0);
		nox_window_set_hidden(*(int*)&dword_5d4594_1123524, 1);
		dword_5d4594_1123520 = 0;
		result = 1;
	}
	return result;
}

//----- (00479B00) --------------------------------------------------------
int nox_xxx_guiDialog_479B00(int a1, int a2, int* a3, int a4) {
	int v3;     // esi
	int result; // eax

	if (a2 != 16391) {
		return 0;
	}
	v3 = nox_xxx_wndGetID_46B0A0(a3);
	if (sub_45D9B0()) {
		return 0;
	}
	nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
	switch (v3) {
	case 3906:
		sub_479950();
		result = 0;
		break;
	case 3907:
		nox_xxx_playDialogFile_44D900(*getMemIntPtr(0x5D4594, 1115312), 100);
		result = 0;
		break;
	case 3908:
		*getMemU8Ptr(0x5D4594, 1123516) = 1;
		LOWORD(a2) = 720;
		BYTE2(a2) = 1;
		nox_netlist_addToMsgListCli_40EBC0(31, 0, &a2, 3);
		result = 0;
		break;
	case 3909:
		*getMemU8Ptr(0x5D4594, 1123516) = 2;
		BYTE2(a2) = 2;
		LOWORD(a2) = 720;
		nox_netlist_addToMsgListCli_40EBC0(31, 0, &a2, 3);
		return 0;
	default:
		return 0;
	}
	return result;
}
// 479B4D: variable 'v4' is possibly undefined

//----- (00479BE0) --------------------------------------------------------
int sub_479BE0(uint32_t* a1, int a2, unsigned int a3, int a4) {
	switch (a2) {
	case 5:
	case 6:
	case 7:
	case 9:
	case 10:
	case 11:
	case 13:
	case 14:
	case 15:
		nox_xxx_wndPointInWnd_46AAB0(a1, (unsigned short)a3, a3 >> 16);
		break;
	default:
		return 1;
	}
	return 1;
}

//----- (00479C40) --------------------------------------------------------
int nox_xxx_wndButtonDrawNoImg_4A81D0(int a1, int a2);
int sub_479C40(uint32_t* a1, int a2) {
	char v2;   // bl
	int yTop;  // [esp+8h] [ebp-8h]
	int xLeft; // [esp+Ch] [ebp-4h]

	v2 = nox_xxx_bookGet_430B40_get_mouse_prev_seq();
	if (!sub_44D930() && (v2 & 0x7Fu) < 0x1E && v2 & 8) {
		nox_client_wndGetPosition_46AA60(a1, &xLeft, &yTop);
		sub_49CD30(xLeft, yTop, a1[2], a1[3] - 2, *getMemIntPtr(0x5D4594, 1107052), 4);
	}
	return nox_xxx_wndButtonDrawNoImg_4A81D0((int)a1, a2);
}

//----- (00479CB0) --------------------------------------------------------
int sub_479CB0(int a1, int a2) {
	int v2; // esi
	int v4; // [esp+4h] [ebp-8h]
	int v5; // [esp+8h] [ebp-4h]

	v2 = *(uint32_t*)(a2 + 24);
	nox_client_wndGetPosition_46AA60(*(uint32_t**)&dword_5d4594_1123524, &v4, &v5);
	nox_client_drawImageAt_47D2C0(v2, nox_win_width - NOX_DEFAULT_WIDTH, nox_win_height - NOX_DEFAULT_HEIGHT);
	return 1;
}

//----- (00479D00) --------------------------------------------------------
int sub_479D00() { return 1; }

//----- (00479D10) --------------------------------------------------------
int sub_479D10() {
	int result; // eax

	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1123524);
	result = 0;
	dword_5d4594_1123524 = 0;
	dword_5d4594_1123520 = 0;
	return result;
}

//----- (0047A260) --------------------------------------------------------
int sub_47A260() { return dword_5d4594_1123520; }

//----- (0047DBC0) --------------------------------------------------------
unsigned char sub_47DBC0() { return getMemByte(0x5D4594, 1193128); }

//----- (0047FCE0) --------------------------------------------------------
int sub_47FCE0(uint32_t* a1, int a2) {
	int v2;            // edx
	unsigned char* v3; // eax
	int v4;            // eax
	int v5;            // esi

	v2 = 0;
	if (*(int*)&dword_5d4594_3804684 > 0) {
		v3 = getMemAt(0x973F18, 6092);
		while (*((uint32_t*)v3 - 1) != a1[3] || *(uint32_t*)v3 != a1[2] || *((uint32_t*)v3 + 1) != a1[21] ||
			   *((uint32_t*)v3 + 2) != a1[26]) {
			++v2;
			v3 += 16;
			if (v2 >= *(int*)&dword_5d4594_3804684) {
				goto LABEL_8;
			}
		}
		return 1;
	}
LABEL_8:
	v4 = 16 * dword_5d4594_3804684;
	v5 = dword_5d4594_3804684 + 1;
	*getMemU32Ptr(0x973F18, 6088 + v4) = a1[3];
	*getMemU32Ptr(0x973F18, 6092 + v4) = a1[2];
	*getMemU32Ptr(0x973F18, 6096 + v4) = a1[21];
	*getMemU32Ptr(0x973F18, 6100 + v4) = a1[26];
	dword_5d4594_3804684 = v5;
	return 1;
}

// 4514E0: using guessed type void  nullsub_4(uint32_t, uint32_t, uint32_t, uint32_t);

//----- (00480220) --------------------------------------------------------
uint8_t* sub_480220(uint8_t* a1, uint8_t* a2) {
	unsigned int v2; // edx
	uint8_t* result; // eax

	result = a2;
	*a1 = 8 * *a2;
	LOWORD(v2) = *(uint16_t*)a2;
	a1[1] = (v2 >> 3) & 0xFC;
	a1[2] = a2[1] & 0xF8;
	return result;
}
// 480232: variable 'v2' is possibly undefined

//----- (00480250) --------------------------------------------------------
uint16_t* sub_480250(uint8_t* a1, uint16_t* a2) {
	uint16_t* result; // eax

	result = a2;
	*a2 = (*a1 >> 3) | (8 * (a1[1] & 0xFC | (32 * (a1[2] & 0xF8))));
	return result;
}

// 487CF0: using guessed type void  nullsub_10(uint32_t);

//----- (004896E0) --------------------------------------------------------
int sub_4896E0() {
	if (dword_5d4594_1193360) {
		free(*(void**)&dword_5d4594_1193360);
	}
	return 1;
}

//----- (00489870) --------------------------------------------------------
int sub_489870() {
	int v0;            // eax
	unsigned char* v1; // esi
	uint32_t* v2;      // eax
	const wchar2_t* v3; // eax
	unsigned int v4;   // eax
	int v5;            // edx
	char v6;           // cl
	int v7;            // eax

	v0 = 0;
	v1 = getMemAt(0x5D4594, 1193388 + 44 * v0);
	if (*getMemU32Ptr(0x5D4594, 1193372 + 4 * v0) == 2) {
		*(uint32_t*)v1 =
			(nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10028)->draw_data.field_0 >> 2) & 1;
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10031);
		v3 = (const wchar2_t*)nox_window_call_field_94((int)v2, 16413, 0, 0);
		*((uint32_t*)v1 + 4) = nox_wcstol(v3, 0, 10);
		*((uint32_t*)v1 + 1) =
			(nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10029)->draw_data.field_0 >> 2) & 1;
		v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10030)->draw_data.field_0;
		*((uint32_t*)v1 + 3) = 0;
		*((uint32_t*)v1 + 2) = (v4 >> 2) & 1;
		if (nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10015)->draw_data.field_0 & 4) {
			v5 = *((uint32_t*)v1 + 3);
			LOBYTE(v5) = v5 | 0x80;
			*((uint32_t*)v1 + 3) = v5;
			v6 = *((uint8_t*)(&nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10016)
								   ->draw_data.field_0));
			v7 = *((uint32_t*)v1 + 3);
			if (v6 & 4) {
				LOBYTE(v7) = v7 | 1;
			} else {
				LOBYTE(v7) = v7 | 2;
			}
			*((uint32_t*)v1 + 3) = v7;
		}
		*((uint32_t*)v1 + 5) =
			(nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10014)->draw_data.field_0 >> 2) & 1;
		*((uint32_t*)v1 + 10) =
			(nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10018)->draw_data.field_0 >> 2) & 1;
	}
	return nox_window_set_hidden(*(int*)&dword_5d4594_1193380, 1);
}

//----- (004899C0) --------------------------------------------------------
int nox_xxx_checkSomeFlagsOnJoin_4899C0(nox_gui_server_ent_t* srv) {
	int a1 = srv;
	int v1;            // eax
	int v2;            // edx
	int v3;            // eax
	unsigned char* v4; // ebp
	int v5;            // eax
	int v6;            // ebx
	unsigned int v7;   // eax
	char v9;           // al
	unsigned char v10; // cl
	unsigned char v11; // cl
	char v12;          // al
	int v13[15];       // [esp+Ch] [ebp-3Ch]
	unsigned char v14; // [esp+4Ch] [ebp+4h]
	unsigned char v15; // [esp+4Ch] [ebp+4h]

	v1 = 0;
	v2 = 11 * v1;
	v3 = *getMemU32Ptr(0x5D4594, 1193372 + 4 * v1);
	v4 = getMemAt(0x5D4594, 1193388 + 4 * v2);
	if (!v3) {
		return 1;
	}
	v5 = v3 - 1;
	if (v5) {
		if (v5 == 1) {
			v6 = a1;
			if (*(uint32_t*)v4) {
				v7 = *(uint32_t*)(a1 + 96);
				if (v7 > *((uint32_t*)v4 + 4) && v7 != 9999) {
					return 0;
				}
			}
			if (*((uint32_t*)v4 + 1) && *(uint8_t*)(a1 + 100) & 0x10) {
				return 0;
			}
			if (*((uint32_t*)v4 + 2) && *(uint8_t*)(a1 + 100) & 0x20) {
				return 0;
			}
			v9 = *(uint8_t*)(a1 + 102);
			if (v9 < 0 && *((uint32_t*)v4 + 3) > (v9 & 0x7F)) {
				return 0;
			}
			if (*((uint32_t*)v4 + 5)) {
				strcpy((char*)v13, (const char*)(a1 + 111));
				sub_57A1E0(v13, 0, 0, 5, *(uint16_t*)(a1 + 163));
				v10 = 0;
				v14 = 0;
				while (v13[v14 + 6] == *(uint32_t*)(4 * v14 + v6 + 135)) {
					v14 = ++v10;
					if (v10 >= 5u) {
						v11 = 0;
						v15 = 0;
						while (*((uint8_t*)&v13[11] + v15) == *(uint8_t*)(v15 + v6 + 155)) {
							v15 = ++v11;
							if (v11 >= 4u) {
								if (v13[12] == *(uint32_t*)(v6 + 159)) {
									goto LABEL_26;
								}
								return 0;
							}
						}
						return 0;
					}
				}
				return 0;
			}
		LABEL_26:
			if (*((uint32_t*)v4 + 10) && *(uint32_t*)(v6 + 48) != NOX_CLIENT_VERS_CODE) {
				return 0;
			}
		}
		return 1;
	}
	v12 = *(uint8_t*)(a1 + 100);
	if (v12 & 0x10) {
		return 0;
	}
	if (v12 & 0x20) {
		return 0;
	}
	return *(uint32_t*)(a1 + 48) == NOX_CLIENT_VERS_CODE;
}

//----- (00489B80) --------------------------------------------------------
uint32_t* sub_489B80(int a1) {
	uint32_t* result;  // eax
	int v2;            // ebx
	unsigned char* v3; // edi
	uint32_t* v4;      // eax
	uint32_t* v5;      // eax
	uint32_t* v6;      // eax
	uint32_t* v7;      // eax
	uint32_t* v8;      // eax
	uint32_t* v9;      // eax
	uint32_t* v10;     // esi
	int v11;           // ebx
	int v12;           // ebx
	wchar2_t v13[16];   // [esp+0h] [ebp-20h]

	result = nox_new_window_from_file("filter.wnd", nox_xxx_windowMplayFilterProc_489E70);
	dword_5d4594_1193380 = result;
	if (result) {
		dword_5d4594_1193384 = nox_xxx_wndGetChildByID_46B0C0(result, 10012);
		v2 = 0;
		v3 = getMemAt(0x5D4594, 1193388 + 44 * v2);
		sub_46B120(*(uint32_t**)&dword_5d4594_1193380, a1);
		sub_46B120(*(uint32_t**)&dword_5d4594_1193384, *(int*)&dword_5d4594_1193380);
		nox_xxx_wndSetProc_46B2C0(*(int*)&dword_5d4594_1193384, nox_xxx_windowMplayFilterProc_489E70);
		if (*(uint32_t*)v3) {
			v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10028);
			v4[9] |= 4u;
		}
		nox_swprintf(v13, L"%d", *((uint32_t*)v3 + 4));
		v5 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10031);
		nox_window_call_field_94((int)v5, 16414, (int)v13, -1);
		if (*((uint32_t*)v3 + 1)) {
			v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10029);
			v6[9] |= 4u;
		}
		if (*((uint32_t*)v3 + 2)) {
			v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10030);
			v7[9] |= 4u;
		}
		if (*((uint32_t*)v3 + 3)) {
			v8 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10015);
			v8[9] |= 4u;
		}
		if (v3[12] & 2) {
			v9 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10017);
		} else {
			v9 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10016);
		}
		v10 = v9;
		v9[9] |= 4u;
		if (*((uint32_t*)v3 + 5)) {
			v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10014);
			v10[9] |= 4u;
		}
		if (*((uint32_t*)v3 + 10)) {
			v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10018);
			v10[9] |= 4u;
		}
		v11 = *getMemU32Ptr(0x5D4594, 1193372 + 4 * v2);
		if (v11) {
			v12 = v11 - 1;
			if (v12) {
				if (v12 == 1) {
					v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10026);
					sub_489DC0();
				}
			} else {
				v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10025);
				nox_window_set_hidden(*(int*)&dword_5d4594_1193384, 1);
			}
		} else {
			v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10024);
			nox_window_set_hidden(*(int*)&dword_5d4594_1193384, 1);
		}
		v10[9] |= 4u;
		result = *(uint32_t**)&dword_5d4594_1193380;
	}
	return result;
}

//----- (00489DC0) --------------------------------------------------------
void sub_489DC0() {
	uint32_t* v0; // eax
	uint32_t* v1; // eax

	nox_window_set_hidden(*(int*)&dword_5d4594_1193384, 0);
	if (nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193384, 10028)->draw_data.field_0 & 4) {
		v0 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193384, 10031);
		nox_xxx_wnd_46ABB0((int)v0, 1);
	} else {
		v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193384, 10031);
		nox_xxx_wnd_46ABB0((int)v1, 0);
	}
	if (nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193384, 10015)->draw_data.field_0 & 4) {
		sub_46AD20(*(uint32_t**)&dword_5d4594_1193384, 10016, 10017, 1);
	} else {
		sub_46AD20(*(uint32_t**)&dword_5d4594_1193384, 10016, 10017, 0);
	}
}

//----- (00489E70) --------------------------------------------------------
int nox_xxx_windowMplayFilterProc_489E70(int a1, int a2, int* a3, int a4) {
	int v3;       // ebx
	int v4;       // esi
	int result;   // eax
	uint32_t* v6; // eax
	int v7;       // [esp-Ch] [ebp-10h]
	int v8;       // [esp-Ch] [ebp-10h]

	v3 = 0;
	if (a2 == 23) {
		return 1;
	}
	if (a2 != 16391) {
		return 0;
	}
	v4 = nox_xxx_wndGetID_46B0A0(a3);
	nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
	switch (v4) {
	case 10015:
		sub_46AD20(*(uint32_t**)&dword_5d4594_1193380, 10016, 10017, ((unsigned int)~a3[9] >> 2) & 1);
		return 0;
	case 10024:
		v7 = dword_5d4594_1193384;
		*getMemU32Ptr(0x5D4594, 1193372 + 4 * v3) = 0;
		nox_window_set_hidden(v7, 1);
		result = 0;
		break;
	case 10025:
		v8 = dword_5d4594_1193384;
		*getMemU32Ptr(0x5D4594, 1193372 + 4 * v3) = 1;
		nox_window_set_hidden(v8, 1);
		result = 0;
		break;
	case 10026:
		*getMemU32Ptr(0x5D4594, 1193372 + 4 * v3) = 2;
		sub_489DC0();
		result = 0;
		break;
	case 10028:
		v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10031);
		nox_xxx_wnd_46ABB0((int)v6, ((unsigned int)~a3[9] >> 2) & 1);
		result = 0;
		break;
	default:
		return 0;
	}
	return result;
}

//----- (00489FB0) --------------------------------------------------------
int sub_489FB0() {
	int result; // eax

	result = dword_5d4594_1193380;
	if (dword_5d4594_1193380) {
		sub_489870();
		nox_xxx_wndClearCaptureMain_46ADE0(*(int*)&dword_5d4594_1193380);
		result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1193380);
		dword_5d4594_1193380 = 0;
	}
	return result;
}

//----- (00489FF0) --------------------------------------------------------
int sub_489FF0(int a1, int a2, const void* a3) {
	int result; // eax

	*getMemU32Ptr(0x5D4594, 1193372 + 4 * a1) = a2;
	result = 11 * a1;
	memcpy(getMemAt(0x5D4594, 1193388 + 44 * a1), a3, 0x2Cu);
	return result;
}

//----- (0048A210) --------------------------------------------------------
int nox_xxx_setSomeFunc_48A210(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1193504) = a1;
	return result;
}
