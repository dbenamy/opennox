#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__system__parsecmd.h"
#include "client/audio/ail/compat_mss.h"

#include "client__gui__guicon.h"
#include "client__gui__guiquit.h"
#include "client__gui__window.h"
#include "client__network__cdecode.h"
#include "client__shell__mainmenu.h"
#include "client__shell__noxworld.h"
#include "client__shell__selchar.h"
#include "client__shell__selcolor.h"

#include "MixPatch.h"
#include "client__io__win95__focus.h"
#include "client__system__ctrlevnt.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__binfile.h"
#include "common__net_list.h"
#include "common__strman.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_5d4594_815748;
extern uint32_t dword_5d4594_816412;
extern uint32_t dword_5d4594_825736;
extern uint32_t dword_5d4594_825744;
extern uint32_t dword_5d4594_816372;
extern void* dword_587000_81128;
extern uint32_t dword_5d4594_816368;
extern uint32_t dword_587000_93156;
extern uint32_t dword_5d4594_816348;
void* dword_5d4594_830236 = 0;
void* dword_5d4594_830232 = 0;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_black_2650656;
extern uint32_t nox_color_orange_2614256;

extern int nox_win_width;
extern int nox_win_height;

uint32_t dword_5d4594_816364 = 0;
uint32_t dword_5d4594_816376 = 0;

//----- (0043B510) --------------------------------------------------------
void nox_client_gui_set_flag_815132(int v);
char* nox_client_getChatMap_49FF40(short* a1);

// 43B510: using guessed type char var_50[80];

//----- (0043B6D0) --------------------------------------------------------


//----- (0043BC10) --------------------------------------------------------


//----- (0043BC80) --------------------------------------------------------


//----- (0043BDB0) --------------------------------------------------------
int sub_43BDB0() { return *getMemU32Ptr(0x5D4594, 815092); }

//----- (0043C650) --------------------------------------------------------
int sub_43C650() {
	long long v0;    // kr00_8
	unsigned int v1; // eax
	int v2;          // ett
	int result;      // eax

	v0 = nox_platform_get_ticks();
	v1 = *getMemU32Ptr(0x5D4594, 815756);
	if (*getMemU64Ptr(0x5D4594, 815740)) {
		v2 = ((unsigned int)v0 < *getMemIntPtr(0x5D4594, 815740)) + *getMemU32Ptr(0x5D4594, 815744);
		*getMemU32Ptr(0x5D4594, 815220 + 8 * *getMemU32Ptr(0x5D4594, 815756)) = v0 - *getMemU32Ptr(0x5D4594, 815740);
		*getMemU32Ptr(0x5D4594, 815224 + 8 * v1) = HIDWORD(v0) - v2;
	} else {
		*getMemU32Ptr(0x5D4594, 815220 + 8 * *getMemU32Ptr(0x5D4594, 815756)) = v0;
		*getMemU32Ptr(0x5D4594, 815224 + 8 * v1) = HIDWORD(v0);
	}
	*getMemU64Ptr(0x5D4594, 815756) = (__PAIR64__(*getMemUintPtr(0x5D4594, 815760), v1) + 1) % 0x3C;
	result = dword_5d4594_815748 + 1;
	*getMemU64Ptr(0x5D4594, 815740) = v0;
	++dword_5d4594_815748;
	return result;
}

//----- (0043CEB0) --------------------------------------------------------
void sub_43CEB0() {
	unsigned int v1;           // esi
	unsigned int v2;           // edi
	unsigned int v3;           // ebx
	unsigned long long v4 = 0; // rax
	unsigned int v5;           // ecx
	int v6;                    // kr00_4
	unsigned int v7;           // kr08_4
	unsigned int v9;           // [esp+0h] [ebp-8h]

	int v0 = dword_5d4594_815748;
	if (*(int*)&dword_5d4594_815748 >= 60) {
		v0 = 60;
	}
	v1 = 0;
	v2 = 0;
	v3 = 0;
	if (!(v0 && v0 > 10)) {
		*getMemU32Ptr(0x587000, 91884) = 0;
		*getMemU32Ptr(0x587000, 91880) = 33;
		return;
	}
	v5 = 0;
	v9 = v0;
	do {
		do {
			v6 = *getMemU32Ptr(0x5D4594, 815220 + 8 * v5) + v2;
			v3 = (__PAIR64__(*getMemU32Ptr(0x5D4594, 815224 + 8 * v5), *getMemU32Ptr(0x5D4594, 815220 + 8 * v5)) +
				  __PAIR64__(v3, v2)) >>
				 32;
			v2 += *getMemU32Ptr(0x5D4594, 815220 + 8 * v5);
			v7 = v5 + 1;
			v1 = (__PAIR64__(v1, v5++) + 1) >> 32;
		} while (v1 < HIDWORD(v4));
		LODWORD(v4) = v9;
	} while (v1 <= HIDWORD(v4) && v7 < v9);
	long long v0a = __PAIR64__(v3, v6) / v4;
	*getMemU64Ptr(0x587000, 91880) = v0a;
}
// 43CFA0: variable 'v1' is possibly undefined

//----- (0043E1A0) --------------------------------------------------------
uint32_t* nox_xxx_gui_43E1A0(int a1) {
	uint32_t* result; // eax

	if (a1) {
		result = nox_window_new(0, 552, 0, 0, nox_win_width, nox_win_height, 0);
		dword_5d4594_816412 = result;
		result[14] = nox_color_black_2650656;
	} else {
		result = *(uint32_t**)&dword_5d4594_816412;
		if (dword_5d4594_816412) {
			result = (uint32_t*)nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_816412);
			dword_5d4594_816412 = 0;
		}
	}
	return result;
}

//----- (0043E8C0) --------------------------------------------------------
int sub_43E8C0(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 816408) = a1;
	return result;
}


//----- (00445450) --------------------------------------------------------
wchar2_t* sub_445450() {
	unsigned char* v0; // esi
	wchar2_t* result;   // eax

	v0 = getMemAt(0x5D4594, 824440);
	do {
		result = nox_wcscpy((wchar2_t*)v0 - 318, (const wchar2_t*)getMemAt(0x5D4594, 825740));
		*(uint32_t*)v0 = 0;
		v0[4] = 0;
		v0 += 644;
	} while ((int)v0 < (int)getMemAt(0x5D4594, 826372));
	dword_5d4594_825736 = 0;
	return result;
}

//----- (00445530) --------------------------------------------------------
int nox_xxx_drawMessageLines_445530() {
	int v1;             // ebp
	int v2;             // eax
	int v3;             // ebx
	int result;         // eax
	int v5;             // ebx
	unsigned short* v6; // ebx
	int v7;             // kr00_4
	unsigned char* v8;  // esi
	int v9;             // edi
	int v10;            // esi
	int v11;            // ecx
	int i;              // [esp+10h] [ebp-Ch]
	int v13;            // [esp+14h] [ebp-8h]
	int v14;            // [esp+18h] [ebp-4h]

	v13 = 0;
	nox_draw_viewport_t* rdr = nox_draw_getViewport_437250();
	v1 = 3 * rdr->height / 4 + rdr->y1 - 15;
	v2 = dword_5d4594_825736;
	for (i = dword_5d4594_825736;; v2 = i) {
		v3 = 161 * v2;
		result = gameFrame();
		v5 = 4 * v3;
		if (*getMemU32Ptr(0x5D4594, 824440 + v5) < gameFrame()) {
			break;
		}
		nox_xxx_drawSetTextColor_434390(nox_color_black_2650656);
		v6 = getMemU16Ptr(0x5D4594, 823804 + v5);
		nox_xxx_drawGetStringSize_43F840(0, v6, &v14, 0, 0);
		v7 = nox_win_width - v14;
		v8 = getMemAt(0x587000, 107848);
		v9 = (nox_win_width - v14) / 2;
		do {
			nox_xxx_drawString_43F6E0(0, (short*)v6, v9 + *(uint32_t*)v8, v1 + *((uint32_t*)v8 + 1));
			v8 += 8;
		} while ((int)v8 < (int)getMemAt(0x587000, 107880));
		v10 = v13;
		if (v13) {
			nox_xxx_drawSetTextColor_434390(*getMemIntPtr(0x5D4594, 2597996));
		} else {
			nox_xxx_drawSetTextColor_434390(nox_color_white_2523948);
		}
		nox_xxx_drawString_43F6E0(0, (short*)v6, v7 / 2, v1);
		v11 = -4 - nox_xxx_guiFontHeightMB_43F320(0);
		result = i;
		v1 += v11;
		if (i) {
			--i;
		} else {
			i = 2;
		}
		++v13;
		if (v10 + 1 >= 3) {
			break;
		}
	}
	return result;
}

//----- (004456E0) --------------------------------------------------------
int nox_xxx_guiChatMode_4456E0(int* a1) {
	int* v1; // esi
	int v2;  // edx
	int v4;  // [esp+4h] [ebp-4h]

	v1 = a1;
	nox_client_wndGetPosition_46AA60(a1, &a1, &v4);
	v2 = v1[25];
	a1 = (int*)((char*)a1 + v1[24]);
	nox_client_drawImageAt_47D2C0(v1[15], (int)a1, v2 + v4);
	return 1;
}

//----- (00445730) --------------------------------------------------------
int nox_xxx_guiChatShowHide_445730(int a1) { return nox_window_set_hidden(*(int*)&dword_5d4594_825744, a1); }

//----- (00445770) --------------------------------------------------------
int sub_445770() {
	int result; // eax

	result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_825744);
	dword_5d4594_825744 = 0;
	return result;
}


//----- (00445B20) --------------------------------------------------------
