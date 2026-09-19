// For inet_addr
#ifdef _WIN32
#include <winsock.h>
#else
#include <arpa/inet.h>
#endif

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
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__system__parsecmd.h"
#include "common__net_list.h"
#include "common__system__settings.h"
#include "common__system__team.h"
#include "common__crypt.h"

#include "client__drawable__drawable.h"
#include "client__gui__gamewin__gamewin.h"
#include "client__gui__guiggovr.h"
#include "client__gui__guiquit.h"
#include "client__gui__window.h"
#include "client__shell__noxworld.h"
#include "client__shell__selchar.h"
#include "client__system__ctrlevnt.h"
#include "client__video__draw_common.h"

#include "client__gui__guicon.h"

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME5_2.h"
#include "client__draw__fx.h"
#include "client__gui__guiinv.h"
#include "client__gui__guimeter.h"
#include "client__gui__guishop.h"
#include "client__gui__guispell.h"
#include "client__gui__servopts__guiserv.h"
#include "client__gui__window.h"
#include "client__io__win95__focus.h"
#include "client__shell__optsback.h"
#include "client__system__ctrlevnt.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__binfile.h"
#include "common__log.h"
#include "common__magic__speltree.h"
#include "defs.h"
#include "input.h"
#include "input_common.h"
#include "operators.h"
#include "server__script__builtin.h"
#include "server__script__script.h"

#include <time.h>

extern uint32_t nox_xxx_useAudio_587000_80772;
extern uint32_t dword_5d4594_811904;
extern uint32_t dword_5d4594_805820;
extern uint32_t dword_5d4594_3807116;
extern uint32_t dword_5d4594_3807152;
extern uint32_t dword_5d4594_3807136;
extern uint32_t dword_5d4594_3807140;
extern uint32_t dword_5d4594_528252;
extern void* nox_alloc_screenParticles_806044;
extern uint32_t dword_5d4594_528256;
extern uint32_t nox_color_white_2523948;
extern uint32_t dword_8531A0_2576;

int nox_win_width = 0;
int nox_win_height = 0;


obj_5D4594_754088_t* ptr_5D4594_754088 = 0;
int ptr_5D4594_754088_cnt = 0;

obj_5D4594_754088_t* ptr_5D4594_754092 = 0;
int ptr_5D4594_754092_cnt = 0;


nox_screenParticle* nox_screenParticles_head = 0;
nox_screenParticle* dword_5d4594_806052 = 0;


void* dword_5d4594_805984 = 0;

//----- (0042A970) --------------------------------------------------------

// 42A970: using guessed type int var_400[256];

//----- (0042C330) --------------------------------------------------------

// 42CC50: using guessed type int sub_42CC50(uint32_t);

//----- (0042E7B0) --------------------------------------------------------
void nox_xxx_clientTalk_42E7B0(nox_drawable* a1p) {
	int a1 = a1p;
	int v1;   // esi
	short v2; // ax

	v1 = a1;
	if (a1 && (!dword_8531A0_2576 || !(*(uint8_t*)(dword_8531A0_2576 + 3680) & 3)) &&
		sub_478030() != 1 && nox_gui_xxx_check_446360() != 1) {
		v2 = *(uint16_t*)(v1 + 128);
		LOWORD(a1) = 464;
		HIWORD(a1) = v2;
		nox_netlist_addToMsgListCli_40EBC0(31, 0, &a1, 4);
	}
}

//----- (0042E810) --------------------------------------------------------
void nox_xxx_clientCollideOrUse_42E810(nox_drawable* a1p) {
	int a1 = a1p;
	int v1; // [esp-4h] [ebp-4h]

	if (a1 && (!dword_8531A0_2576 || !(*(uint8_t*)(dword_8531A0_2576 + 3680) & 3))) {
		v1 = a1;
		LOBYTE(a1) = 123;
		*(uint16_t*)((char*)&a1 + 1) = nox_xxx_netGetUnitCodeCli_578B00(v1);
		nox_netlist_addToMsgListCli_40EBC0(31, 0, &a1, 3);
	}
}

//----- (0042E850) --------------------------------------------------------
void nox_xxx_clientTrade_42E850(nox_drawable* a1p) {
	int a1 = a1p;
	int v1; // esi

	v1 = a1;
	if (a1 && (!dword_8531A0_2576 || !(*(uint8_t*)(dword_8531A0_2576 + 3680) & 3)) &&
		sub_47A260() != 1 && nox_gui_xxx_check_446360() != 1) {
		LOWORD(a1) = 5577;
		HIWORD(a1) = nox_xxx_netGetUnitCodeCli_578B00(v1);
		nox_netlist_addToMsgListCli_40EBC0(31, 0, &a1, 4);
	}
}

//----- (0042EB90) --------------------------------------------------------
int sub_42EB90(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 754052) = a1;
	return result;
}

//----- (0042EBA0) --------------------------------------------------------
int sub_42EBA0() { return *getMemU32Ptr(0x5D4594, 754052); }

//----- (0042EDC0) --------------------------------------------------------
void sub_42EDC0() {
	if (ptr_5D4594_754088) {
		free(ptr_5D4594_754088);
		ptr_5D4594_754088 = 0;
	}
	if (ptr_5D4594_754092) {
		free(ptr_5D4594_754092);
		ptr_5D4594_754092 = 0;
	}
}

//----- (00430AA0) --------------------------------------------------------
int sub_430AA0(int a1) {
	int result; // eax

	result = a1 - 1;
	if (a1 == 1) {
		dword_5d4594_805820 = 1;
		nox_xxx_useAudio_587000_80772 = 9;
	} else {
		result = a1 - 2;
		if (a1 == 2) {
			dword_5d4594_805820 = 2;
			nox_xxx_useAudio_587000_80772 = 13;
		} else {
			dword_5d4594_805820 = 0;
			nox_xxx_useAudio_587000_80772 = 5;
		}
	}
	return result;
}

//----- (00430AF0) --------------------------------------------------------
int nox_client_mousePriKey_430AF0() { return dword_5d4594_805820; }

//----- (00430B00) --------------------------------------------------------
int nox_xxx_cursor_430B00() { return nox_xxx_useAudio_587000_80772; }

//----- (00430B10) --------------------------------------------------------
void nox_client_setMousePos_430B10(int x, int y) { nox_client_changeMousePos_430A00(x, y, true); }

//----- (00431270) --------------------------------------------------------
void sub_431270() {
	if (dword_5d4594_805984) {
		sub_487680(dword_5d4594_805984);
		dword_5d4594_805984 = 0;
	}
}

//----- (00431290) --------------------------------------------------------
void sub_431290() {
	if (dword_5d4594_805984) {
		sub_487970(dword_5d4594_805984, -1);
	}
}

//----- (00431770) --------------------------------------------------------
char* nox_xxx_getHostInfoPtr_431770() { return (char*)getMemAt(0x5D4594, 807172); }

//----- (00431790) --------------------------------------------------------
char* nox_xxx_copyServerIPAndPort_431790(char* a1) {
	char* result; // eax

	result = a1;
	if (a1) {
		result = strncpy((char*)getMemAt(0x5D4594, 806060), a1, 0x17u);
	}
	return result;
}
// 4335F8: variable 'v12' is possibly undefined

//----- (00435040) --------------------------------------------------------
void sub_435040() {
	pixel8888 buf[256]; // [esp+4h] [ebp-400h]
	unsigned char* data;

	data = getMemAt(0x973F18, 3880);
	for (int i = 0; i < 256; ++i) {
		buf[i].field_0 = i;
		buf[i].field_1 = data[4 * i + 0];
		buf[i].field_2 = data[4 * i + 1];
		buf[i].field_3 = data[4 * i + 2];
	}
	sub_48C580(buf, 256);

	data = getMemAt(0x5D4594, 809604);
	for (int i = 0; i < 256; ++i) {
		data[4 * i + 0] = buf[i].field_1;
		data[4 * i + 1] = buf[i].field_2;
		data[4 * i + 2] = buf[i].field_3;
		data[4 * i + 3] = 0;
		*getMemU8Ptr(0x5D4594, 808304 + i) = buf[i].field_0;
	}
}

//----- (00435120) --------------------------------------------------------
void sub_435120(void* a1, void* a2) {
	char* result; // eax
	uint8_t* v3;  // ecx
	int v4;       // esi
	char v5;      // bl
	char* v6;     // eax
	char v7;      // bl
	char v8;      // bl

	result = a2;
	v3 = a1;
	v4 = 256;
	do {
		v5 = *result;
		v6 = result + 1;
		*v3 = v5;
		v7 = *v6++;
		v3[1] = v7;
		v8 = *v6;
		v3[3] = 4;
		v3[2] = v8;
		result = v6 + 1;
		v3 += 4;
		--v4;
	} while (v4);
}

//----- (00435150) --------------------------------------------------------
void sub_435150(uint8_t* a1, char* a2) {
	char* v2;        // ecx
	uint8_t* result; // eax
	int v4;          // esi
	char v5;         // dl
	uint8_t* v6;     // eax

	v2 = a2;
	result = a1;
	v4 = 256;
	do {
		v5 = *v2;
		v2 += 4;
		*result = v5;
		v6 = result + 1;
		*v6++ = *(v2 - 3);
		*v6 = *(v2 - 2);
		result = v6 + 1;
		--v4;
	} while (v4);
}

//----- (00435570) --------------------------------------------------------
long long nox_xxx_initTime_435570() {
	long long result; // rax

	result = nox_platform_get_ticks();
	*getMemU64Ptr(0x5D4594, 811908) = result;
	return result;
}

//----- (00435690) --------------------------------------------------------
uint32_t* sub_435690(uint32_t* a1) {
	uint32_t* result; // eax

	result = a1;
	*a1 = *getMemU32Ptr(0x5D4594, 811364);
	a1[1] = *getMemU32Ptr(0x5D4594, 811368);
	return result;
}

//----- (004356C0) --------------------------------------------------------
bool nox_client_drawable_testBuff_4356C0(nox_drawable* dr, char a2) {
	int a1 = dr;
	int result; // eax

	result = a1;
	if (a1) {
		result = ((1 << a2) & *(uint32_t*)(a1 + 124)) != 0;
	}
	return result;
}

//----- (00435700) --------------------------------------------------------
wchar2_t* sub_435700(wchar2_t* a1, int a2) {
	wchar2_t* result; // eax

	result = nox_wcscpy((wchar2_t*)getMemAt(0x5D4594, 811376), a1);
	*getMemU32Ptr(0x5D4594, 811060) = a2;
	return result;
}

//----- (004357A0) --------------------------------------------------------
int nox_xxx_cliToggleObsWindow_4357A0() {
	int result; // eax

	if (dword_8531A0_2576 && *(uint8_t*)(dword_8531A0_2576 + 3680) & 1) {
		result = nox_xxx_showObserverWindow_48CA70(0);
	} else {
		result = nox_xxx_showObserverWindow_48CA70(1);
	}
	return result;
}

//----- (00435F60) --------------------------------------------------------
int sub_435F60() {
	int result; // eax

	result = 1 - dword_5d4594_811904;
	dword_5d4594_811904 = 1 - dword_5d4594_811904;
	return result;
}

//----- (00436550) --------------------------------------------------------
int sub_436550() {
	int v0; // eax

	if (sub_459DA0() || nox_gui_xxx_check_446360() || sub_49CB40() || sub_49C810() || sub_446950() || sub_4706A0() ||
		nox_gui_console_flagXxx_451410()) {
		v0 = gameFrame();
	} else {
		v0 = gameFrame();
		if (gameFrame() != 2) {
			return gameFrame() - *getMemU32Ptr(0x5D4594, 811920) == 1;
		}
	}
	*getMemU32Ptr(0x5D4594, 811920) = v0;
	return 1;
}

//----- (00437100) --------------------------------------------------------
void sub_4706C0(int a1);
void sub_437100() {
	int result; // eax

	int flag = nox_client_getRenderGUI();
	if (*getMemU32Ptr(0x5D4594, 811064) != flag &&
		!nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING)) {
		*getMemU32Ptr(0x5D4594, 811064) = flag;
		sub_4721A0(flag);
		sub_460EA0(flag);
		nox_window_set_visible_unk5(flag);
		sub_45D500(flag);
		sub_455A00(flag);
		sub_455F10(flag);
		sub_4706C0(flag);
		if (!flag) {
			sub_478000();
		}
	}
}

//----- (004372B0) --------------------------------------------------------
int nox_xxx_playerAnimCheck_4372B0() {
	int v0;     // eax
	int result; // eax

	result = 1;
	if (*getMemU32Ptr(0x852978, 8)) {
		v0 = *(uint32_t*)(*getMemU32Ptr(0x852978, 8) + 276);
		if (v0 != 1 && v0 != 2 && v0 != 51) {
			result = 0;
		}
	}
	return result;
}

//----- (004372E0) --------------------------------------------------------
int nox_xxx_clientIsObserver_4372E0() {
	int result; // eax

	if (dword_8531A0_2576 && *(uint32_t*)(dword_8531A0_2576 + 2092) == 1) {
		result = (*(uint32_t*)(dword_8531A0_2576 + 3680) & 3) != 0;
	} else {
		result = 0;
	}
	return result;
}

//----- (00437320) --------------------------------------------------------


//----- (004375C0) --------------------------------------------------------


//----- (00437860) --------------------------------------------------------


//----- (004379C0) --------------------------------------------------------


//----- (00438330) --------------------------------------------------------


//----- (00438370) --------------------------------------------------------


//----- (00438480) --------------------------------------------------------


//----- (00438C80) --------------------------------------------------------


//----- (00438DD0) --------------------------------------------------------


//----- (00438E30) --------------------------------------------------------


//----- (00438EF0) --------------------------------------------------------


//----- (00439050) --------------------------------------------------------


// 439385: variable 'v2' is possibly undefined

//----- (00439450) --------------------------------------------------------


//----- (00439CC0) --------------------------------------------------------


//----- (00439D00) --------------------------------------------------------


//----- (00439D90) --------------------------------------------------------


// 43A3CE: variable 'v22' is possibly undefined

//----- (0043A920) --------------------------------------------------------


//----- (0043A9D0) --------------------------------------------------------


//----- (0043AA70) --------------------------------------------------------


//----- (0043AF30) --------------------------------------------------------


//----- (0043AF40) --------------------------------------------------------


//----- (0043AF80) --------------------------------------------------------


//----- (0043AF90) --------------------------------------------------------


//----- (0043AFA0) --------------------------------------------------------


//----- (0043B300) --------------------------------------------------------


//----- (0043B320) --------------------------------------------------------


//----- (0043B340) --------------------------------------------------------


//----- (0043B460) --------------------------------------------------------


//----- (0043B490) --------------------------------------------------------
