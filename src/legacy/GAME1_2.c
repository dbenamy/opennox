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

extern uint32_t dword_5d4594_3807116;
extern uint32_t dword_5d4594_3807152;
extern uint32_t dword_5d4594_3807136;
extern uint32_t dword_5d4594_3807140;
extern void* nox_alloc_screenParticles_806044;
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


//----- (0042E810) --------------------------------------------------------


//----- (0042E850) --------------------------------------------------------


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


//----- (00430AF0) --------------------------------------------------------


//----- (00430B00) --------------------------------------------------------


//----- (00430B10) --------------------------------------------------------


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


//----- (00435690) --------------------------------------------------------


//----- (004356C0) --------------------------------------------------------


//----- (00435700) --------------------------------------------------------


//----- (004357A0) --------------------------------------------------------


//----- (00435F60) --------------------------------------------------------


//----- (00436550) --------------------------------------------------------


//----- (00437100) --------------------------------------------------------
void sub_4706C0(int a1);


//----- (004372B0) --------------------------------------------------------


//----- (004372E0) --------------------------------------------------------


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
