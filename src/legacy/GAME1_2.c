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
extern uint32_t dword_5d4594_815052;
extern uint32_t dword_5d4594_3807116;
extern uint32_t dword_5d4594_3807152;
extern uint32_t dword_5d4594_3807136;
extern uint32_t dword_5d4594_814548;
extern uint32_t dword_5d4594_3807140;
extern uint32_t nox_client_connError_814552;
extern uint32_t dword_5d4594_815056;
extern uint32_t dword_5d4594_741356;
extern uint32_t dword_5d4594_741364;
extern nox_gui_animation* nox_wnd_xxx_815040;
extern uint32_t dword_587000_87408;
extern uint32_t dword_5d4594_528252;
extern void* nox_alloc_screenParticles_806044;
extern uint32_t dword_5d4594_815044;
extern uint32_t nox_wol_server_result_cnt_815088;
extern uint32_t dword_5d4594_528256;
extern uint32_t dword_5d4594_815032;
extern uint32_t dword_5d4594_815020;
extern uint32_t dword_5d4594_815024;
extern uint32_t dword_5d4594_815028;
extern uint32_t dword_5d4594_814984;
extern uint32_t dword_5d4594_815016;
extern uint32_t nox_game_createOrJoin_815048;
extern uint32_t dword_587000_87412;
extern uint32_t dword_5d4594_815000;
extern uint32_t nox_wol_wnd_gameList_815012;
extern nox_window* dword_5d4594_815004;
extern nox_window* nox_wol_wnd_world_814980;
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

void* dword_5d4594_814624 = 0;

void* dword_5d4594_805984 = 0;











//----- (004282D0) --------------------------------------------------------


//----- (004282F0) --------------------------------------------------------


//----- (00428540) --------------------------------------------------------


//----- (004285C0) --------------------------------------------------------


//----- (00428810) --------------------------------------------------------


//----- (00428890) --------------------------------------------------------


//----- (004289D0) --------------------------------------------------------


//----- (00428B30) --------------------------------------------------------
int nox_server_mapRWObjectTOC_428B30() {
	int v1;            // eax
	unsigned short v2; // bp
	int v3;            // esi
	int v4;            // eax
	char* v5;            // esi
	int i;             // esi
	unsigned short v7; // ax
	int v8;            // [esp+4h] [ebp-110h]
	int v9;            // [esp+8h] [ebp-10Ch]
	int v10;           // [esp+Ch] [ebp-108h]
	int v11;           // [esp+10h] [ebp-104h]
	char v12[256];     // [esp+14h] [ebp-100h]

	v11 = 1;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v11, 2u);
	if ((short)v11 > 1) {
		return 0;
	}
	sub_42BFB0();
	if (nox_crypt_IsReadOnly()) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v9, 2u);
		for (i = 0; (unsigned short)i < (unsigned short)v9; ++i) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v10, 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v8, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v12, (unsigned char)v8);
			v12[(unsigned char)v8] = 0;
			if (!nox_common_gameFlags_check_40A5C0(2) || nox_common_gameFlags_check_40A5C0(1)) {
				v7 = nox_xxx_getNameId_4E3AA0(v12);
			} else {
				v7 = nox_xxx_getTTByNameSpriteMB_44CFC0(v12);
			}
			sub_42C310(v7, v10);
		}
		return 1;
	}
	sub_42BFE0();
	LOWORD(v1) = sub_42C300();
	v9 = v1;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v9, 2u);
	v2 = 0;
	if (!nox_xxx_unitDefGetCount_4E3AC0()) {
		return 1;
	}
	v3 = 0;
	do {
		LOWORD(v4) = sub_42C2E0(v3);
		v10 = v4;
		if ((uint16_t)v4) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v10, 2u);
			v5 = nox_xxx_getUnitNameByThingType_4E3A80(v3);
			LOBYTE(v8) = strlen(v5);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v8, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v5, (unsigned char)v8);
		}
		v3 = ++v2;
	} while (v2 < (unsigned int)nox_xxx_unitDefGetCount_4E3AC0());
	return 1;
}
// 428B85: variable 'v1' is possibly undefined
// 428BAB: variable 'v4' is possibly undefined
// 428B30: using guessed type char var_100[256];

//----- (00429200) --------------------------------------------------------
int nox_server_mapRWAmbientData_429200() {
	int result; // eax
	char* v1;   // esi
	int v2;     // [esp+0h] [ebp-10h]
	int v3[3];  // [esp+4h] [ebp-Ch]

	v2 = 1;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v2, 2u);
	if ((short)v2 < 1) {
		return 0;
	}
	if (nox_crypt_IsReadOnly()) {
		if (nox_crypt_IsReadOnly() == 1) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v3[0], 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v3[1], 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v3[2], 4u);
			sub_469B90(v3);
			if (nox_common_gameFlags_check_40A5C0(2097154)) {
				sub_4349C0(v3);
			}
		}
		result = 1;
	} else {
		v1 = nox_xxx_getAmbientColor_469BB0();
		nox_xxx_fileReadWrite_426AC0_file3_fread(v1, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(v1 + 4, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(v1 + 8, 4u);
		result = 1;
	}
	return result;
}

//----- (004292C0) --------------------------------------------------------
int nox_server_mapRWWindowWalls_4292C0(uint32_t* a1) {
	int result;   // eax
	uint32_t* v2; // edi
	char* v3;     // esi
	uint32_t* v4; // eax
	int v5;       // [esp+4h] [ebp-20h]
	int v6;       // [esp+8h] [ebp-1Ch]
	int2 v7;
	int4 v9; // [esp+14h] [ebp-10h]

	v5 = 2;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 2u);
	if ((short)v5 > 2) {
		return 0;
	}
	if (nox_crypt_IsReadOnly()) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741336), 2u);
		v6 = 0;
		if (*getMemU16Ptr(0x5D4594, 741336) > 0) {
			v2 = a1;
			do {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v7, 8u);
				if (a1) {
					v3 = nox_xxx_mapGetWallSize_426A70();
					sub_428170(a1, &v9);
					v7.field_0 += v9.field_0 / 23 - *(uint32_t*)v3;
					v7.field_4 += v9.field_4 / 23 - *((uint32_t*)v3 + 1);
				}
				if (nox_common_gameFlags_check_40A5C0(0x400000)) {
					v2 = 0;
					v4 = nox_xxx_cliWallGet_5042F0(v7.field_0, v7.field_4);
					if (v4) {
						v2 = (uint32_t*)*v4;
					}
				} else {
					v2 = (uint32_t*)nox_server_getWallAtGrid_410580(v7.field_0, v7.field_4);
				}
				if (v2) {
					*((uint8_t*)v2 + 4) |= 0x40u;
					if ((short)v5 < 2) {
						*((uint8_t*)v2 + 2) = 0;
					}
				}
				++v6;
			} while (v6 < *getMemI16Ptr(0x5D4594, 741336));
		}
		result = 1;
	} else {
		*getMemU16Ptr(0x5D4594, 741336) = 0;
		nox_xxx_wallForeachFn_410640(sub_429450, (int)a1);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741336), 2u);
		nox_xxx_wallForeachFn_410640(sub_4294B0, (int)a1);
		result = 1;
	}
	return result;
}

//----- (00429450) --------------------------------------------------------
void sub_429450(uint8_t* a1, uint32_t* a2) {
	int v2;  // eax
	int2 v3; // [esp+4h] [ebp-8h]

	if (!a2 || (v2 = (unsigned char)a1[6], v3.field_0 = 23 * (unsigned char)a1[5], v3.field_4 = 23 * v2,
				nox_xxx_wallMath_427F30(&v3, a2))) {
		if (a1[4] & 0x40) {
			++*getMemU16Ptr(0x5D4594, 741336);
		}
	}
}

//----- (004294B0) --------------------------------------------------------
void sub_4294B0(uint8_t* a1, uint32_t* a2) {
	int v2;  // eax
	int v3;  // edx
	int2 v4; // [esp+4h] [ebp-8h]

	if (!a2 || (v2 = (unsigned char)a1[6], v4.field_0 = 23 * (unsigned char)a1[5], v4.field_4 = 23 * v2,
				nox_xxx_wallMath_427F30(&v4, a2))) {
		if (a1[4] & 0x40) {
			v3 = (unsigned char)a1[6];
			v4.field_0 = (unsigned char)a1[5];
			v4.field_4 = v3;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v4, 8u);
		}
	}
}

//----- (00429520) --------------------------------------------------------
void nox_xxx_wallBreakableCounterClear_429520() { *getMemU32Ptr(0x5D4594, 741344) = 0; }

//----- (00429530) --------------------------------------------------------
int nox_server_mapRWDestructableWalls_429530(uint32_t* a1) {
	int result;   // eax
	uint32_t* v2; // edi
	char* v3;     // esi
	uint32_t* v4; // eax
	int v5;       // [esp+4h] [ebp-20h]
	int v6;       // [esp+8h] [ebp-1Ch]
	int2 v7;
	int4 v9; // [esp+14h] [ebp-10h]

	v5 = 1;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 2u);
	if ((short)v5 > 1) {
		return 0;
	}
	if (nox_crypt_IsReadOnly()) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741340), 2u);
		v6 = 0;
		if (*getMemU16Ptr(0x5D4594, 741340) > 0) {
			v2 = a1;
			do {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v7, 8u);
				if (a1) {
					v3 = nox_xxx_mapGetWallSize_426A70();
					sub_428170(a1, &v9);
					v7.field_0 += v9.field_0 / 23 - *(uint32_t*)v3;
					v7.field_4 += v9.field_4 / 23 - *((uint32_t*)v3 + 1);
				}
				if (nox_common_gameFlags_check_40A5C0(0x400000)) {
					v2 = 0;
					v4 = nox_xxx_cliWallGet_5042F0(v7.field_0, v7.field_4);
					if (v4) {
						v2 = (uint32_t*)*v4;
					}
				} else {
					v2 = (uint32_t*)nox_server_getWallAtGrid_410580(v7.field_0, v7.field_4);
				}
				if (v2) {
					*((uint8_t*)v2 + 4) |= 8u;
					*((uint16_t*)v2 + 5) = *getMemU16Ptr(0x5D4594, 741344);
					++*getMemU32Ptr(0x5D4594, 741344);
					if (!nox_common_gameFlags_check_40A5C0(0x400000)) {
						nox_xxx_wallBreackableListAdd_410840((int)v2);
					}
				}
				++v6;
			} while (v6 < *getMemI16Ptr(0x5D4594, 741340));
		}
		result = 1;
	} else {
		*getMemU16Ptr(0x5D4594, 741340) = 0;
		nox_xxx_wallForeachFn_410640(nox_xxx_wall_4296E0, (int)a1);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741340), 2u);
		nox_xxx_wallForeachFn_410640(sub_429740, (int)a1);
		result = 1;
	}
	return result;
}

//----- (004296E0) --------------------------------------------------------
void nox_xxx_wall_4296E0(uint8_t* a1, uint32_t* a2) {
	int v2;  // eax
	int2 v3; // [esp+4h] [ebp-8h]

	if (!a2 || (v2 = (unsigned char)a1[6], v3.field_0 = 23 * (unsigned char)a1[5], v3.field_4 = 23 * v2,
				nox_xxx_wallMath_427F30(&v3, a2))) {
		if (a1[4] & 8) {
			++*getMemU16Ptr(0x5D4594, 741340);
		}
	}
}

//----- (00429740) --------------------------------------------------------
void sub_429740(uint8_t* a1, uint32_t* a2) {
	int v2;  // eax
	int v3;  // edx
	int2 v4; // [esp+4h] [ebp-8h]

	if (!a2 || (v2 = (unsigned char)a1[6], v4.field_0 = 23 * (unsigned char)a1[5], v4.field_4 = 23 * v2,
				nox_xxx_wallMath_427F30(&v4, a2))) {
		if (a1[4] & 8) {
			v3 = (unsigned char)a1[6];
			v4.field_0 = (unsigned char)a1[5];
			v4.field_4 = v3;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v4, 8u);
		}
	}
}

//----- (004297B0) --------------------------------------------------------
void nox_xxx_wallSecretCounterClear_4297B0() { *getMemU32Ptr(0x5D4594, 741352) = 0; }

//----- (004297C0) --------------------------------------------------------
int nox_server_mapRWSecretWalls_4297C0(uint32_t* a1) {
	char* v2;    // esi
	int* v3;     // edi
	uint8_t* v4; // ebx
	char* v5;    // ebp
	int* v6;     // eax
	int v7;      // eax
	char v8;     // dl
	int v9;      // [esp+4h] [ebp-1Ch]
	int v10 = 0; // [esp+8h] [ebp-18h]
	int v11;     // [esp+Ch] [ebp-14h]
	int4 v12;    // [esp+10h] [ebp-10h]

	v9 = 2;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v9, 2u);
	if ((short)v9 > 2) {
		return 0;
	}
	if (!nox_crypt_IsReadOnly()) {
		*getMemU16Ptr(0x5D4594, 741348) = 0;
		nox_xxx_wallForeachFn_410640(sub_429A00, (int)a1);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741348), 2u);
		nox_xxx_wallForeachFn_410640(sub_429A60, (int)a1);
		return 1;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741348), 2u);
	v11 = 0;
	if (*getMemU16Ptr(0x5D4594, 741348) <= 0) {
		return 1;
	}
	while (1) {
		v2 = (char*)calloc(1u, 0x20u);
		v3 = (int*)(v2 + 4);
		nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 4, 8u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 16, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 20, 1u);
		v4 = v2 + 21;
		v2[21] = 0;
		if ((short)v9 >= 2) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 21, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 22, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 24, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 28, 4u);
		}
		if (a1) {
			v5 = nox_xxx_mapGetWallSize_426A70();
			sub_428170(a1, &v12);
			*v3 += v12.field_0 / 23 - *(uint32_t*)v5;
			*((uint32_t*)v2 + 2) += v12.field_4 / 23 - *((uint32_t*)v5 + 1);
		}
		if (!nox_common_gameFlags_check_40A5C0(0x400000)) {
			v10 = nox_server_getWallAtGrid_410580(*v3, *((uint32_t*)v2 + 2));
			v7 = v10;
		} else {
			v6 = nox_xxx_cliWallGet_5042F0(*v3, *((uint32_t*)v2 + 2));
			if (!v6) {
				v7 = 0;
			} else {
				v7 = *v6;
				v10 = v7;
			}
		}
		if (v7) {
			v8 = *(uint8_t*)(v7 + 4);
			*(uint32_t*)(v7 + 28) = v2;
			*(uint8_t*)(v7 + 4) = v8 | 4;
			*(uint16_t*)(v7 + 10) = *getMemU16Ptr(0x5D4594, 741352);
			++*getMemU32Ptr(0x5D4594, 741352);
			*((uint32_t*)v2 + 3) = v7;
			if (!*v4) {
				if (v2[20] & 8) {
					*((uint32_t*)v2 + 7) = -1;
					*v4 = 3;
					v2[22] = 23;
				} else {
					*((uint32_t*)v2 + 7) = 0;
					*v4 = 1;
					v2[22] = 0;
				}
			}
			if (!nox_common_gameFlags_check_40A5C0(0x400000)) {
				nox_xxx_wallSecretBlock_410760(v2);
			}
		} else {
			free(v2);
		}
		++v11;
		if (v11 >= *getMemI16Ptr(0x5D4594, 741348)) {
			return 1;
		}
	}
}

//----- (00429A00) --------------------------------------------------------
void sub_429A00(uint8_t* a1, uint32_t* a2) {
	int v2;  // eax
	int2 v3; // [esp+4h] [ebp-8h]

	if (!a2 || (v2 = (unsigned char)a1[6], v3.field_0 = 23 * (unsigned char)a1[5], v3.field_4 = 23 * v2,
				nox_xxx_wallMath_427F30(&v3, a2))) {
		if (a1[4] & 4) {
			++*getMemU16Ptr(0x5D4594, 741348);
		}
	}
}

//----- (00429A60) --------------------------------------------------------
void sub_429A60(int a1, uint32_t* a2) {
	int v2;      // eax
	int v3;      // edx
	uint8_t* v4; // esi
	int2 v5;     // [esp+4h] [ebp-8h]

	if (!a2 || (v2 = *(unsigned char*)(a1 + 6), v5.field_0 = 23 * *(unsigned char*)(a1 + 5), v5.field_4 = 23 * v2,
				nox_xxx_wallMath_427F30(&v5, a2))) {
		if (*(uint8_t*)(a1 + 4) & 4) {
			v3 = *(unsigned char*)(a1 + 6);
			v4 = *(uint8_t**)(a1 + 28);
			v5.field_0 = *(unsigned char*)(a1 + 5);
			v5.field_4 = v3;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 8u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 16, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 20, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 21, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 22, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 24, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 28, 4u);
		}
	}
}

//----- (00429B20) --------------------------------------------------------
int nox_server_mapRWWallMap_429B20(uint32_t* a1) {
	uint32_t* v2;                       // ebp
	int v3;                             // esi
	int v5;                             // edx
	int v6;                             // edx
	int v7;                             // ecx
	int v8;                             // eax
	int v9;                             // edi
	uint8_t* v10;                       // eax
	uint8_t* v11;                       // esi
	nox_player_polygon_check_data* v12; // eax
	int v13;                            // eax
	int v14;                            // eax
	char v15;                           // bl
	char v16;                           // bl
	unsigned char* v17;                 // esi
	int v18;                            // edi
	unsigned char* v19;                 // eax
	unsigned char* v20;                 // edi
	unsigned char* v21;                 // ebx
	char v22;                           // [esp+2h] [ebp-3Ah]
	char v23;                           // [esp+3h] [ebp-39h]
	int v24;                            // [esp+4h] [ebp-38h]
	int v25;                            // [esp+8h] [ebp-34h]
	int v26;                            // [esp+Ch] [ebp-30h]
	int v27;                            // [esp+10h] [ebp-2Ch]
	int v28;                            // [esp+14h] [ebp-28h]
	int v29;                            // [esp+18h] [ebp-24h]
	int v30;                            // [esp+1Ch] [ebp-20h]
	int v31;                            // [esp+20h] [ebp-1Ch]
	int2 v32;                           // [esp+24h] [ebp-18h]
	int4 v33;                           // [esp+2Ch] [ebp-10h]

	v31 = nox_xxx_wallGet_426A30();
	if (!getMemByte(0x5D4594, 741372)) {
		*getMemU8Ptr(0x5D4594, 741372) = nox_xxx_wallTileByName_410D60("MagicWallSystemUseOnly");
	}
	nox_xxx_wallSecretCounterClear_4297B0();
	nox_xxx_wallBreakableCounterClear_429520();
	v28 = 7;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v28, 2u);
	if ((short)v28 > 7) {
		return 0;
	}
	if ((short)v28 < 6) {
		return sub_42A150(v28, a1);
	}
	v2 = a1;
	if (!nox_crypt_IsReadOnly()) {
		if (a1) {
			sub_428170(a1, &v33);
			*getMemU32Ptr(0x5D4594, 741360) = v33.field_0 / 23;
			v3 = v33.field_8 / 23;
			dword_5d4594_741356 = v33.field_8 / 23;
			*getMemU32Ptr(0x5D4594, 741368) = v33.field_4 / 23;
			v5 = v33.field_C / 23;
			dword_5d4594_741364 = v5;
		} else {
			*getMemU32Ptr(0x5D4594, 741368) = 256;
			*getMemU32Ptr(0x5D4594, 741360) = 256;
			dword_5d4594_741364 = 0;
			dword_5d4594_741356 = 0;
			nox_xxx_wallForeachFn_410640(sub_42A0F0, 0);
			v3 = dword_5d4594_741356;
			v5 = dword_5d4594_741364;
		}
		v25 = v3 - *getMemU32Ptr(0x5D4594, 741360) + 1;
		v27 = v5 - *getMemU32Ptr(0x5D4594, 741368) + 1;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741360), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741368), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v25, 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v27, 4u);
	v26 = 0;
	v29 = 0;
	if (nox_crypt_IsReadOnly()) {
		if (v2) {
			sub_428170(v2, &v33);
			v13 = v33.field_0 / 23 - *getMemU32Ptr(0x5D4594, 741360);
			*getMemU32Ptr(0x5D4594, 741360) = v33.field_0 / 23;
			v29 = v13;
			v14 = v33.field_4 / 23 - *getMemU32Ptr(0x5D4594, 741368);
			*getMemU32Ptr(0x5D4594, 741368) = v33.field_4 / 23;
			v26 = v14;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v22, 1u);
		v15 = v22;
		if (v22 == -1) {
			return 1;
		}
		while (1) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 1u);
			LOBYTE(v30) = v29 + v15;
			LOBYTE(a1) = v26 + (uint8_t)a1;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v24, 1u);
			v16 = (unsigned char)v24 >> 7;
			LOBYTE(v24) = v24 & 0x7F;
			if (nox_common_gameFlags_check_40A5C0(0x400000)) {
				v17 = (unsigned char*)*sub_504290(v30, (char)a1);
			} else {
				v18 = (unsigned char)v30;
				v19 = (unsigned char*)nox_server_getWallAtGrid_410580((unsigned char)v30, (unsigned char)a1);
				v17 = v19;
				if (v19) {
					if (v31 & 1) {
						*v19 = nox_xxx_wall_42A6C0(*v19, v24);
					} else {
						*v19 = v24;
					}
					goto LABEL_46;
				}
				v17 = (unsigned char*)nox_xxx_wallCreateAt_410250(v18, (unsigned char)a1);
				if (!v17) {
					return 0;
				}
			}
			*v17 = v24;
		LABEL_46:
			if (v16) {
				v17[4] |= 0x80u;
			}
			v20 = v17 + 1;
			nox_xxx_fileReadWrite_426AC0_file3_fread(v17 + 1, 1u);
			v21 = v17 + 2;
			nox_xxx_fileReadWrite_426AC0_file3_fread(v17 + 2, 1u);
			if (v31 & 1 && *v21 >= nox_xxx_mapWallMaxVariation_410DD0(*v20, *v17, 0)) {
				*v21 = 0;
			}
			v17[7] = nox_xxx_mapWallGetHpByTile_410E20(*v20);
			if ((uint16_t)v28 == 6) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v24, 1u);
				v17[8] = 1;
				*((uint32_t*)v17 + 3) = 0;
			} else {
				nox_xxx_fileReadWrite_426AC0_file3_fread(v17 + 8, 1u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v24, 1u);
				*((uint32_t*)v17 + 3) = (unsigned char)v24;
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v22, 1u);
			v15 = v22;
			if (v22 == -1) {
				return 1;
			}
		}
	}
	v6 = *getMemU32Ptr(0x5D4594, 741368);
	v26 = *getMemU32Ptr(0x5D4594, 741368);
	if (*getMemU32Ptr(0x5D4594, 741368) <= *getMemIntPtr(0x5D4594, 741368) + v27) {
		v7 = v25;
		v8 = *getMemU32Ptr(0x5D4594, 741360);
		do {
			v9 = v8;
			v29 = v8;
			if (v8 <= v8 + v7) {
				do {
					v10 = (uint8_t*)nox_server_getWallAtGrid_410580(v9, v26);
					v11 = v10;
					if (v10) {
						if (v10[1] != getMemByte(0x5D4594, 741372)) {
							if (!v2 || (v32.field_0 = 23 * (unsigned char)v10[5],
										v32.field_4 = 23 * (unsigned char)v10[6], nox_xxx_wallMath_427F30(&v32, v2))) {
								nox_xxx_fileReadWrite_426AC0_file3_fread(v11 + 5, 1u);
								nox_xxx_fileReadWrite_426AC0_file3_fread(v11 + 6, 1u);
								if ((int8_t)v11[4] >= 0) {
									LOBYTE(v24) = *v11;
								} else {
									LOBYTE(v24) = *v11 | 0x80;
								}
								nox_xxx_fileReadWrite_426AC0_file3_fread(&v24, 1u);
								nox_xxx_fileReadWrite_426AC0_file3_fread(v11 + 1, 1u);
								nox_xxx_fileReadWrite_426AC0_file3_fread(v11 + 2, 1u);
								v32.field_0 = 23 * (unsigned char)v11[5] + 11;
								v32.field_4 = 23 * (unsigned char)v11[6] + 11;
								v12 = nox_xxx_polygonIsPlayerInPolygon_4217B0(&v32, 0);
								if (v12 || (v12 = (nox_player_polygon_check_data*)sub_421990(&v32, 10.0, 0)) != 0) {
									v23 = BYTE2(v12->field_0[32]);
								} else {
									v23 = 100;
								}
								nox_xxx_fileReadWrite_426AC0_file3_fread(&v23, 1u);
								if (nox_common_gameFlags_check_40A5C0(0x200000)) {
									LOBYTE(v24) = 0;
								} else {
									LOBYTE(v24) = v11[12];
								}
								nox_xxx_fileReadWrite_426AC0_file3_fread(&v24, 1u);
								v2 = a1;
								v9 = v29;
							}
						}
					}
					v8 = *getMemU32Ptr(0x5D4594, 741360);
					v7 = v25;
					v29 = ++v9;
				} while (v9 <= *getMemIntPtr(0x5D4594, 741360) + v25);
				v6 = *getMemU32Ptr(0x5D4594, 741368);
			}
			++v26;
		} while (v26 <= v6 + v27);
	}
	v22 = -1;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v22, 1u);
	return 1;
}

//----- (0042A0F0) --------------------------------------------------------
int sub_42A0F0(int a1) {
	int result; // eax

	if ((int)*(unsigned char*)(a1 + 5) < *getMemIntPtr(0x5D4594, 741360)) {
		*getMemU32Ptr(0x5D4594, 741360) = *(unsigned char*)(a1 + 5);
	}
	if ((int)*(unsigned char*)(a1 + 5) > *(int*)&dword_5d4594_741356) {
		dword_5d4594_741356 = *(unsigned char*)(a1 + 5);
	}
	if ((int)*(unsigned char*)(a1 + 6) < *getMemIntPtr(0x5D4594, 741368)) {
		*getMemU32Ptr(0x5D4594, 741368) = *(unsigned char*)(a1 + 6);
	}
	result = *(unsigned char*)(a1 + 6);
	if (result > *(int*)&dword_5d4594_741364) {
		dword_5d4594_741364 = *(unsigned char*)(a1 + 6);
	}
	return result;
}

//----- (0042A150) --------------------------------------------------------
int sub_42A150(short a1, uint32_t* a2) {
	int v2;                             // eax
	uint32_t* v3;                       // ebx
	int v4;                             // edi
	int v5;                             // esi
	int v6;                             // edx
	int v7;                             // edx
	int v8;                             // ebp
	int v9;                             // ecx
	int v10;                            // eax
	int v11;                            // edi
	int v12;                            // eax
	uint8_t* v13;                       // esi
	nox_player_polygon_check_data* v14; // eax
	int v16;                            // eax
	int v17;                            // ebp
	unsigned char v18;                  // bl
	char v19;                           // bl
	unsigned char** v20;                // eax
	unsigned char* v21;                 // esi
	unsigned char* v22;                 // eax
	unsigned char v23;                  // al
	short v24;                          // bx
	unsigned char* v25;                 // edi
	char v26;                           // [esp+13h] [ebp-2Dh]
	int v27;                            // [esp+14h] [ebp-2Ch]
	int v28;                            // [esp+18h] [ebp-28h]
	int v29;                            // [esp+1Ch] [ebp-24h]
	int v30;                            // [esp+20h] [ebp-20h]
	int v31;                            // [esp+24h] [ebp-1Ch]
	int2 v32;                           // [esp+28h] [ebp-18h]
	int4 v33;                           // [esp+30h] [ebp-10h]

	v2 = nox_xxx_wallGet_426A30();
	v3 = a2;
	v30 = v2;
	v4 = 0;
	if (!nox_crypt_IsReadOnly()) {
		if (a2) {
			sub_428170(a2, &v33);
			*getMemU32Ptr(0x5D4594, 741360) = v33.field_0 / 23;
			v5 = v33.field_8 / 23;
			dword_5d4594_741356 = v33.field_8 / 23;
			*getMemU32Ptr(0x5D4594, 741368) = v33.field_4 / 23;
			v6 = v33.field_C / 23;
			dword_5d4594_741364 = v33.field_C / 23;
		} else {
			*getMemU32Ptr(0x5D4594, 741368) = 256;
			*getMemU32Ptr(0x5D4594, 741360) = 256;
			dword_5d4594_741364 = 0;
			dword_5d4594_741356 = 0;
			nox_xxx_wallForeachFn_410640(sub_42A0F0, 0);
			v5 = dword_5d4594_741356;
			v6 = dword_5d4594_741364;
		}
		v28 = v5 - *getMemU32Ptr(0x5D4594, 741360) + 1;
		v29 = v6 - *getMemU32Ptr(0x5D4594, 741368) + 1;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741360), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 741368), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v28, 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v29, 4u);
	if (!nox_crypt_IsReadOnly()) {
		v7 = *getMemU32Ptr(0x5D4594, 741368);
		v8 = *getMemU32Ptr(0x5D4594, 741368);
		if (*getMemU32Ptr(0x5D4594, 741368) > *getMemIntPtr(0x5D4594, 741368) + v29) {
			return 1;
		}
		v9 = v28;
		v10 = *getMemU32Ptr(0x5D4594, 741360);
		while (1) {
			v11 = v10;
			if (v10 > v10 + v9) {
				goto LABEL_27;
			}
			do {
				v12 = nox_server_getWallAtGrid_410580(v11, v8);
				v13 = (uint8_t*)v12;
				if (v3) {
					if (!v12) {
						LOBYTE(v27) = -1;
						goto LABEL_19;
					}
					v32.field_0 = 23 * *(unsigned char*)(v12 + 5);
					v32.field_4 = 23 * *(unsigned char*)(v12 + 6);
					if (!nox_xxx_wallMath_427F30(&v32, v3)) {
						v13 = 0;
						LOBYTE(v27) = -1;
						goto LABEL_19;
					}
				}
				if (!v13) {
					LOBYTE(v27) = -1;
				} else if ((int)v13[4] >= 0) {
					LOBYTE(v27) = *v13;
				} else {
					LOBYTE(v27) = *v13 | 0x80;
				}
			LABEL_19:
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v27, 1u);
				if ((uint8_t)v27 != (uint8_t)-1) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(v13 + 1, 1u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v13 + 2, 1u);
					v32.field_0 = 23 * (unsigned char)v13[5] + 11;
					v32.field_4 = 23 * (unsigned char)v13[6] + 11;
					v14 = nox_xxx_polygonIsPlayerInPolygon_4217B0(&v32, 0);
					if (v14 || (v14 = (nox_player_polygon_check_data*)sub_421990(&v32, 10.0, 0)) != 0) {
						v26 = BYTE2(v14->field_0[32]);
					} else {
						v26 = 1;
					}
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v26, 1u);
					LOBYTE(a1) = v13[12];
					nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 1u);
				}
				v10 = *getMemU32Ptr(0x5D4594, 741360);
				v9 = v28;
				++v11;
			} while (v11 <= *getMemIntPtr(0x5D4594, 741360) + v28);
			v7 = *getMemU32Ptr(0x5D4594, 741368);
		LABEL_27:
			++v8;
			if (v8 > v7 + v29) {
				break;
			}
		}
		return 1;
	}
	if (v3) {
		sub_428170(v3, &v33);
		*getMemU32Ptr(0x5D4594, 741360) = v33.field_0 / 23;
		*getMemU32Ptr(0x5D4594, 741368) = v33.field_4 / 23;
	}
	v31 = 0;
	if (v29 < 0) {
		return 1;
	}
	v16 = v28;
	while (1) {
		v17 = 0;
		if (v16 >= 0) {
			while (1) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v27, 1u);
				if ((uint8_t)v27 != (uint8_t)-1) {
					v18 = v27;
					LOBYTE(v27) = v27 & 0x7F;
					v19 = v18 >> 7;
					if (!nox_common_gameFlags_check_40A5C0(0x400000)) {
						v22 = (unsigned char*)nox_server_getWallAtGrid_410580(v17 + *getMemU32Ptr(0x5D4594, 741360), v4 + *getMemU32Ptr(0x5D4594, 741368));
						v21 = v22;
						if (v22) {
							if (v30 & 1) {
								v23 = nox_xxx_wall_42A6C0(*v22, v27);
							} else {
								v23 = v27;
							}
						} else {
							v21 = (unsigned char*)nox_xxx_wallCreateAt_410250(v17 + *getMemU32Ptr(0x5D4594, 741360), v4 + *getMemU32Ptr(0x5D4594, 741368));
							if (!v21) {
								return 0;
							}
							v23 = v27;
						}
						*v21 = v23;
					} else {
						v20 = (unsigned char**)sub_504290(v17 + getMemByte(0x5D4594, 741360), v4 + getMemByte(0x5D4594, 741368));
						v21 = *v20;
						**v20 = v27;
					}
					if (v19) {
						v21[4] |= 0x80u;
					}
					v24 = a1;
					if (a1 < 2) {
						v21[1] = 0;
						v21[7] = nox_xxx_mapWallGetHpByTile_410E20(0);
						v21[8] = 1;
					} else {
						v25 = v21 + 1;
						nox_xxx_fileReadWrite_426AC0_file3_fread(v21 + 1, 1u);
						if (v24 >= 3) {
							nox_xxx_fileReadWrite_426AC0_file3_fread(v21 + 2, 1u);
						} else {
							sub_42A650(v21);
						}
						if (v30 & 1 && v21[2] >= nox_xxx_mapWallMaxVariation_410DD0(*v25, *v21, 0)) {
							v21[2] = 0;
						}
						v21[7] = nox_xxx_mapWallGetHpByTile_410E20(*v25);
						if (v24 < 4) {
							v21[8] = 1;
						} else {
							nox_xxx_fileReadWrite_426AC0_file3_fread(v21 + 8, 1u);
						}
						LOBYTE(a2) = 0;
						if (v24 >= 5) {
							nox_xxx_fileReadWrite_426AC0_file3_fread(&a2, 1u);
						}
						v4 = v31;
						*((uint32_t*)v21 + 3) = (unsigned char)a2;
					}
				}
				v16 = v28;
				++v17;
				if (v17 > v28) {
					break;
				}
			}
		}
		++v4;
		v31 = v4;
		if (v4 > v29) {
			break;
		}
	}
	return 1;
}

//----- (0042A650) --------------------------------------------------------
int sub_42A650(unsigned char* a1) {
	unsigned char v1; // cl
	int result;       // eax

	v1 = *a1;
	a1[2] = 0;
	if (!v1) {
		a1[2] = a1[5] % 3;
	}
	if (v1 == 1) {
		a1[2] = a1[5] % 3;
	}
	result = nox_xxx_getWallSprite_46A3B0(a1[1], v1, a1[2], (a1[4] >> 2) & 2);
	if (!result) {
		a1[2] = 0;
	}
	return result;
}

//----- (0042A6E0) --------------------------------------------------------
int nox_server_mapRWMapInfo_42A6E0() {
	int vers = 3;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&vers, 2);
	if (vers > 3) {
		return 0;
	}
	if (vers >= 1) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 2408), 64);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 2472), 512);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 2984), 16);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3000), 64);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3064), 64);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3128), 128);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3256), 128);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3384), 256);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3640), 128);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3768), 32);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3800), 4);
		if (vers == 2) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3804), 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3805), 1u);
		} else if (nox_crypt_IsReadOnly() == 1) {
			*getMemU8Ptr(0x973F18, 3804) = 2;
			*getMemU8Ptr(0x973F18, 3805) = 16;
		}
	}
	if (vers < 3) {
		*getMemU8Ptr(0x973F18, 3806) = getMemByte(0x5D4594, 741376);
		*getMemU8Ptr(0x973F18, 3838) = getMemByte(0x5D4594, 741380);
	} else {
		int v2 = strlen(getMemAt(0x973F18, 3806));
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v2, 1);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3806), v2);
		*getMemU8Ptr(0x973F18, 3806 + v2) = 0;

		v2 = strlen(getMemAt(0x973F18, 3838));
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v2, 1);
		nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x973F18, 3838), v2);
		*getMemU8Ptr(0x973F18, 3838 + v2) = 0;
	}
	return 1;
}

//----- (0042A8B0) --------------------------------------------------------


//----- (0042A970) --------------------------------------------------------

// 42A970: using guessed type int var_400[256];

//----- (0042AAA0) --------------------------------------------------------


//----- (0042ABF0) --------------------------------------------------------


//----- (0042AC50) --------------------------------------------------------


//----- (0042ADA0) --------------------------------------------------------


//----- (0042B810) --------------------------------------------------------


//----- (0042BCE0) --------------------------------------------------------


//----- (0042BD50) --------------------------------------------------------


//----- (0042BDC0) --------------------------------------------------------


//----- (0042BE30) --------------------------------------------------------


//----- (0042BEA0) --------------------------------------------------------


//----- (0042C330) --------------------------------------------------------

// 42CC50: using guessed type int sub_42CC50(uint32_t);

//----- (0042C360) --------------------------------------------------------


//----- (0042C480) --------------------------------------------------------


//----- (0042C770) --------------------------------------------------------


//----- (0042C7F0) --------------------------------------------------------


//----- (0042C820) --------------------------------------------------------


//----- (0042C8B0) --------------------------------------------------------


//----- (0042C8E0) --------------------------------------------------------


//----- (0042C910) --------------------------------------------------------


//----- (0042C9A0) --------------------------------------------------------


//----- (0042CA00) --------------------------------------------------------


//----- (0042CB20) --------------------------------------------------------


//----- (0042CB80) --------------------------------------------------------


//----- (0042CBF0) --------------------------------------------------------


//----- (0042CC70) --------------------------------------------------------


//----- (0042CCE0) --------------------------------------------------------


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
int sub_437320(int a1) {
	int v1;            // eax
	unsigned char* v2; // ecx
	int v3;            // ecx
	int v4;            // eax
	int result;        // eax

	v1 = 0;
	v2 = getMemAt(0x587000, 87484);
	do {
		if (*(uint32_t*)(a1 + 96) <= *(int*)v2) {
			break;
		}
		v2 += 4;
		++v1;
	} while ((int)v2 < (int)getMemAt(0x587000, 87496));
	if (v1 > 2) {
		v1 = 2;
	}
	v3 = *(uint32_t*)(a1 + 28) + 36;
	v4 = 16 * v1;
	if (*(int*)&dword_587000_87412 == -1) {
		*(uint32_t*)(*(uint32_t*)(a1 + 28) + 76) = *getMemU32Ptr(0x5D4594, 814568 + v4);
		*(uint32_t*)(v3 + 24) = *getMemU32Ptr(0x5D4594, 814564 + v4);
		result = *getMemU32Ptr(0x5D4594, 814564 + v4);
	} else {
		*(uint32_t*)(*(uint32_t*)(a1 + 28) + 76) = *getMemU32Ptr(0x5D4594, 814560 + v4);
		*(uint32_t*)(v3 + 24) = *getMemU32Ptr(0x5D4594, 814556 + v4);
		result = *getMemU32Ptr(0x5D4594, 814556 + v4);
	}
	*(uint32_t*)(v3 + 48) = result;
	return result;
}

//----- (004375C0) --------------------------------------------------------
void sub_4375C0(int a1) {
	if (nox_wol_server_result_cnt_815088) {
		sub_46AD20(*(uint32_t**)&nox_wol_wnd_world_814980, 10070, nox_wol_server_result_cnt_815088 + 10069, a1);
	}
}

//----- (00437860) --------------------------------------------------------
int sub_437860(int a1, int a2) {
	int result;        // eax
	unsigned char* v3; // ecx

	result = 0;
	v3 = getMemAt(0x587000, 87532);
	while (a1 <= *((short*)v3 - 2) || a1 >= *(short*)v3 || a2 <= *((short*)v3 - 1) || a2 >= *((short*)v3 + 1)) {
		v3 += 8;
		++result;
		if ((int)v3 >= (int)getMemAt(0x587000, 87564)) {
			return 0;
		}
	}
	return result;
}

//----- (004379C0) --------------------------------------------------------
void sub_4379C0() {
	if (dword_587000_87408 == 1) {
		nox_window_call_field_94(*(int*)&nox_wol_wnd_gameList_815012, 16399, 0, 0);
	}
}

//----- (00438330) --------------------------------------------------------
int sub_438330() {
	int (*v0)(void); // esi

	v0 = nox_wnd_xxx_815040->field_13;
	nox_gui_freeAnimation_43C570(nox_wnd_xxx_815040);
	if (!nox_common_gameFlags_check_40A5C0(0x10000000)) {
		nox_client_guiXxx_43A9D0();
	}
	if (v0) {
		v0();
	}
	return 1;
}

//----- (00438370) --------------------------------------------------------
int sub_438370() {
	if (nox_wnd_xxx_815040->state == NOX_GUI_ANIM_OUT_DONE) {
		return sub_438330();
	}
	nox_wnd_xxx_815040->state = NOX_GUI_ANIM_OUT;
	sub_43BE40(2);
	nox_xxx_clientPlaySoundSpecial_452D80(923, 100);
	return 1;
}

//----- (00438480) --------------------------------------------------------
int sub_438480() {
	nox_xxx_wndSetProc_46B2C0(*(int*)&nox_wol_wnd_gameList_815012, sub_439050);
	nox_xxx_wndSetWindowProc_46B300(*(int*)&nox_wol_wnd_gameList_815012, sub_438EF0);
	sub_46B120(*(uint32_t**)&dword_5d4594_815016, *(int*)&nox_wol_wnd_gameList_815012);
	nox_xxx_wndSetWindowProc_46B300(*(int*)&dword_5d4594_815016, sub_438EF0);
	sub_46B120(*(uint32_t**)&dword_5d4594_815020, *(int*)&nox_wol_wnd_gameList_815012);
	nox_xxx_wndSetWindowProc_46B300(*(int*)&dword_5d4594_815020, sub_438EF0);
	sub_46B120(*(uint32_t**)&dword_5d4594_815024, *(int*)&nox_wol_wnd_gameList_815012);
	nox_xxx_wndSetWindowProc_46B300(*(int*)&dword_5d4594_815024, sub_438EF0);
	sub_46B120(*(uint32_t**)&dword_5d4594_815028, *(int*)&nox_wol_wnd_gameList_815012);
	nox_xxx_wndSetWindowProc_46B300(*(int*)&dword_5d4594_815028, sub_438EF0);
	sub_46B120(*(uint32_t**)&dword_5d4594_815032, *(int*)&nox_wol_wnd_gameList_815012);
	nox_xxx_wndSetWindowProc_46B300(*(int*)&dword_5d4594_815032, sub_438EF0);
	**(uint32_t**)(*(uint32_t*)((uint32_t)dword_5d4594_815004 + 32) + 28) = 10035;
	**(uint32_t**)(*(uint32_t*)((uint32_t)dword_5d4594_815004 + 32) + 32) = 10036;
	**(uint32_t**)(*(uint32_t*)((uint32_t)dword_5d4594_815004 + 32) + 36) = 10032;
	nox_window_call_field_94(*(int*)&dword_5d4594_815016, 16408,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 28), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815020, 16408,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 28), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815024, 16408,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 28), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815028, 16408,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 28), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815032, 16408,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 28), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815016, 16409,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 32), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815020, 16409,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 32), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815024, 16409,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 32), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815028, 16409,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 32), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815032, 16409,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 32), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815016, 16410,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 36), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815020, 16410,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 36), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815024, 16410,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 36), 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_815028, 16410,
							 *(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 36), 0);
	return nox_window_call_field_94(*(int*)&dword_5d4594_815032, 16410,
									*(uint32_t*)(*(uint32_t*)(nox_wol_wnd_gameList_815012 + 32) + 36), 0);
}

//----- (00438C80) --------------------------------------------------------
int sub_438C80(int a1, int a2) {
	char v2[404]; // [esp+4h] [ebp-194h]

	nox_point mpos = nox_client_getMousePos_4309F0();
	if (!wndIsShown_nox_xxx_wndIsShown_46ACC0(*(int*)&dword_5d4594_815000)) {
		memcpy(v2, *(const void**)&dword_5d4594_815000, sizeof(v2));
		*(uint32_t*)&v2[16] -= 32;
		*(uint32_t*)&v2[20] -= 32;
		*(uint32_t*)&v2[8] += 64;
		*(uint32_t*)&v2[12] += 64;
		if (!dword_5d4594_815044 && !nox_xxx_wndPointInWnd_46AAB0(v2, mpos.x, mpos.y)) {
			nox_window_set_hidden(*(int*)&dword_5d4594_815000, 1);
			nox_window_call_field_94(*(int*)&nox_wol_wnd_gameList_815012, 16403, -1, 0);
			dword_5d4594_815056 = 0;
			nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_815000);
			nox_xxx_windowFocus_46B500(*(int*)&nox_wol_wnd_world_814980);
		}
	}
	if (sub_4A28B0() && !nox_xxx_wndPointInWnd_46AAB0(*(uint32_t**)getMemAt(0x5D4594, 815036), mpos.x, mpos.y)) {
		sub_4A2890();
		nox_xxx_windowFocus_46B500(*(int*)&nox_wol_wnd_world_814980);
	}
	if (nox_game_createOrJoin_815048 && sub_438DD0(mpos.x, mpos.y)) {
		nox_client_setCursorType_477610(9);
	} else if (!sub_44A4A0()) {
		nox_client_setCursorType_477610(0);
	}
	return 1;
}

//----- (00438DD0) --------------------------------------------------------
int sub_438DD0(unsigned int a1, unsigned int a2) {
	if (*(int*)&dword_587000_87412 == -1) {
		if (a1 > 0xD8 && a1 < 0x258 && a2 > 0x1B && a2 < 0x1C3) {
			return 1;
		}
	} else if (a1 > 0xE2 && a1 < 0x24E && a2 > 0x25 && a2 < 0x1B9) {
		return 1;
	}
	return 0;
}

//----- (00438E30) --------------------------------------------------------
int sub_438E30(uint32_t* a1, int a2) {
	uint32_t* v1; // esi
	int v2;       // edx
	int v3;       // esi
	short** v4;   // edi
	int v6;       // [esp+4h] [ebp-4h]

	v1 = a1;
	nox_client_wndGetPosition_46AA60(a1, &v6, &a1);
	v2 = v1[25];
	if (v1[9] & 6) {
		nox_client_drawImageAt_47D2C0(v1[19], v6 + v1[24], (int)a1 + v2);
	} else {
		nox_client_drawImageAt_47D2C0(v1[15], v6 + v1[24], (int)a1 + v2);
	}
	v3 = v1[100];
	if (!v3) {
		return 1;
	}
	do {
		if (!(*(uint8_t*)(v3 + 4) & 0x10) && *(uint32_t*)(v3 + 44) == 2048) {
			v4 = *(short***)(v3 + 32);
			nox_xxx_drawSetTextColor_434390(nox_color_white_2523948);
			nox_xxx_drawStringStyle_43F7B0(*(uint32_t*)(v3 + 236), *v4, v6 + *(uint32_t*)(v3 + 16), (int)a1 + *(uint32_t*)(v3 + 20));
		}
		v3 = *(uint32_t*)(v3 + 388);
	} while (v3);
	return 1;
}

//----- (00438EF0) --------------------------------------------------------
int sub_438EF0(uint32_t* a1, int a2, unsigned int a3, int a4) {
	int result; // eax
	int v5;     // esi
	int v6;     // esi

	if (a2 == 19) {
		v6 = a1[8];
		nox_window_call_field_94(*(int*)&nox_wol_wnd_gameList_815012, 16391, *(uint32_t*)(v6 + 28), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815016, 16391, *(uint32_t*)(v6 + 28), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815020, 16391, *(uint32_t*)(v6 + 28), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815024, 16391, *(uint32_t*)(v6 + 28), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815028, 16391, *(uint32_t*)(v6 + 28), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815032, 16391, *(uint32_t*)(v6 + 28), 0);
		result = 0;
	} else if (a2 == 20) {
		v5 = a1[8];
		nox_window_call_field_94(*(int*)&nox_wol_wnd_gameList_815012, 16391, *(uint32_t*)(v5 + 32), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815016, 16391, *(uint32_t*)(v5 + 32), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815020, 16391, *(uint32_t*)(v5 + 32), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815024, 16391, *(uint32_t*)(v5 + 32), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815028, 16391, *(uint32_t*)(v5 + 32), 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_815032, 16391, *(uint32_t*)(v5 + 32), 0);
		result = 0;
	} else {
		result = nox_xxx_wndListboxProcWithoutData10_4A28E0(a1, a2, a3, a4);
	}
	return result;
}

//----- (00439050) --------------------------------------------------------
int sub_439050(int a1, unsigned int a2, int* a3, unsigned int a4) {
	int v4; // edi
	int v5; // edi
	int v7; // eax

	if (a2 > 0x400F) {
		if (a2 == 16400) {
			v7 = nox_xxx_wndGetID_46B0A0(a3);
			if (v7 >= 10038 && v7 <= 10042) {
				nox_window_call_field_94(*(int*)&dword_5d4594_815016, 16403, a4, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_815020, 16403, a4, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_815024, 16403, a4, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_815028, 16403, a4, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_815032, 16403, a4, 0);
				if (a4 < *(int*)&nox_wol_server_result_cnt_815088) {
					nox_point pos = nox_client_getMousePos_4309F0();
					dword_5d4594_814624 = sub_4A04C0(a4);
					nox_client_gui_serverInfoBlockCheckExp_439370(&pos, dword_5d4594_814624);
				}
			}
		} else if (a2 == 16403 || a2 == 16412) {
			nox_window_call_field_94(*(int*)&dword_5d4594_815016, a2, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815020, a2, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815024, a2, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815028, a2, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815032, a2, (int)a3, 0);
		}
	} else if (a2 >= 0x400E) {
		nox_window_call_field_94(*(int*)&dword_5d4594_815016, a2, (int)a3, a4);
		nox_window_call_field_94(*(int*)&dword_5d4594_815020, a2, (int)a3, a4);
		nox_window_call_field_94(*(int*)&dword_5d4594_815024, a2, (int)a3, a4);
		nox_window_call_field_94(*(int*)&dword_5d4594_815028, a2, (int)a3, a4);
		nox_window_call_field_94(*(int*)&dword_5d4594_815032, a2, (int)a3, a4);
	} else {
		switch (a2) {
		case 0x17u:
			return 1;
		case 0x4000u:
			if (a3 == nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&nox_wol_wnd_gameList_815012, 10043) ||
				a3 == nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&nox_wol_wnd_gameList_815012, 10044)) {
				nox_window_call_field_94(*(int*)&dword_5d4594_815016, 0x4000, (int)a3, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_815020, 0x4000, (int)a3, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_815024, 0x4000, (int)a3, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_815028, 0x4000, (int)a3, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_815032, 0x4000, (int)a3, 0);
			}
			break;
		case 0x4009u:
			v4 = *(uint32_t*)(a1 + 32);
			nox_xxx_wndListboxProcPre_4A30D0(a1, 0x4009u, (wchar2_t*)a3, a4);
			v5 = sub_4A4800(v4);
			nox_window_call_field_94(*(int*)&dword_5d4594_815016, 16412, v5, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815020, 16412, v5, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815024, 16412, v5, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815028, 16412, v5, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815032, 16412, v5, 0);
			break;
		}
	}
	return nox_xxx_wndListboxProcPre_4A30D0(a1, a2, (wchar2_t*)a3, a4);
}

// 439385: variable 'v2' is possibly undefined

//----- (00439450) --------------------------------------------------------
uint32_t* sub_439450(int a1, int a2, uint32_t* a3) {
	uint32_t* result; // eax
	int v4;           // ecx

	result = a3;
	*a3 = a1 - 65;
	a3[1] = a2 - 20;
	if (a1 - 65 + 130 > 600) {
		*a3 = 470;
	}
	if (a2 - 20 + 120 > 451) {
		a3[1] = 331;
	}
	if (dword_587000_87408 == 1) {
		v4 = 55;
		if (a3[1] >= 55) {
			goto LABEL_10;
		}
	} else {
		v4 = 27;
		if (a3[1] >= 27) {
			goto LABEL_10;
		}
	}
	a3[1] = v4;
LABEL_10:
	if ((int)*a3 < 216) {
		*a3 = 216;
	}
	return result;
}

//----- (00439CC0) --------------------------------------------------------
char* sub_439CC0(int a1, char* a2) {
	size_t v2;    // esi
	char* result; // eax

	v2 = (size_t)&strstr((const char*)(a1 + 52), "'s_game")[-52 - a1];
	result = strncpy(a2, (const char*)(a1 + 52), v2);
	a2[v2] = 0;
	return result;
}

//----- (00439D00) --------------------------------------------------------
int sub_439D00(int* a1, int a2, unsigned int a3, int a4) {
	if (a2 == 5) {
		if (nox_xxx_wndGetID_46B0A0(a1) == 10020 && nox_game_createOrJoin_815048 == 1) {
			sub_439D90((unsigned short)a3, a3 >> 16);
			return 1;
		}
		return 0;
	}
	if (a2 != 21) {
		return 0;
	}
	if (a3 != 1) {
		if (a3 != 28 && a3 == 57) {
			nox_point mpos = nox_client_getMousePos_4309F0();
			nox_window_call_field_93(a1, 5, mpos.x | (mpos.y << 16), 0);
		}
		return 0;
	}
	if (a4 == 2) {
		sub_4373A0();
	}
	return 1;
}

//----- (00439D90) --------------------------------------------------------
int sub_439D90(unsigned int a1, unsigned int a2) {
	int result; // eax
	short v3;   // dx

	result = sub_438DD0(a1, a2);
	if (result) {
		v3 = a2 + *getMemU16Ptr(0x587000, 87530 + 8 * dword_587000_87412) - 27;
		*getMemU16Ptr(0x5D4594, 814916) = a1 + *getMemU16Ptr(0x587000, 87528 + 8 * dword_587000_87412) - 216;
		*getMemU16Ptr(0x5D4594, 814918) = v3;
		sub_43B460();
		if (sub_43BDB0() & 2) {
			nox_xxx_setQuest_4D6F60(1);
		}
		if (nox_xxx_isQuest_4D6F50()) {
			if (nox_client_countPlayerFiles04_4DC7D0()) {
				sub_4A7A70(1);
				nox_game_showSelChar_4A4DB0();
				return nox_client_setCursorType_477610(0);
			}
		} else if (nox_client_countPlayerFiles02_4DC630()) {
			sub_4A7A70(1);
			nox_game_showSelChar_4A4DB0();
			return nox_client_setCursorType_477610(0);
		}
		sub_4A7A70(0);
		nox_game_showSelClass_4A4840();
		nox_client_setCursorType_477610(0);
	}
	return result;
}

// 43A3CE: variable 'v22' is possibly undefined

//----- (0043A920) --------------------------------------------------------
int sub_43A920() {
	int result; // eax

	nox_xxx_windowFocus_46B500(*(int*)&nox_wol_wnd_world_814980);
	if (!sub_43BE30() || !*getMemU32Ptr(0x5D4594, 815084)) {
		sub_44A400();
	}
	result = sub_43AF90(0);
	dword_5d4594_815044 = 0;
	return result;
}

//----- (0043A9D0) --------------------------------------------------------
int nox_client_guiXxx_43A9D0() {
	nox_xxx_wndClearCaptureMain_46ADE0(*(int*)&dword_5d4594_814984);
	sub_489FB0();
	sub_4A2890();
	if (dword_5d4594_815000 && !*(uint32_t*)(dword_5d4594_815000 + 396)) {
		nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_815000);
		nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_815000);
		dword_5d4594_815000 = 0;
	}
	if (nox_wol_wnd_world_814980) {
		nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&nox_wol_wnd_world_814980);
		nox_wol_wnd_world_814980 = 0;
	}
	sub_43A920();
	nox_xxx_windowFocus_46B500(0);
	nox_wol_server_result_cnt_815088 = 0;
	sub_49FFA0(0);
	sub_554D10();
	nox_client_setCursorType_477610(0);
	return sub_43DE40(0);
}

//----- (0043AA70) --------------------------------------------------------
char* sub_43AA70() {
	char* v0;      // esi
	char* v1;      // ebx
	char* v2;      // eax
	char v3;       // al
	char v4;       // al
	char v5;       // cl
	char* v6;      // eax
	char* result;  // eax
	short v8;      // cx
	char v9[32];   // [esp+0h] [ebp-12Ch]
	char v10[268]; // [esp+20h] [ebp-10Ch]

	if (dword_5d4594_528252 && dword_5d4594_528256) {
		nox_xxx_networkLog_printf_413D30("RECON: Posting server to WOL");
	}
	nox_game_createOrJoin_815048 = 0;
	dword_5d4594_815052 = 1;
	v0 = nox_xxx_cliGamedataGet_416590(0);
	v1 = sub_416640();
	memcpy(v1 + 111, v0, 0x3Au);
	*(uint16_t*)(v1 + 163) = nox_common_gameFlags_getVal_40A5B0();
	*(uint32_t*)(v1 + 135) = -1;
	*(uint32_t*)(v1 + 139) = -1;
	*(uint32_t*)(v1 + 143) = -1;
	*(uint32_t*)(v1 + 147) = -1;
	*(uint32_t*)(v1 + 151) = -1;
	*(uint32_t*)(v1 + 155) = -1;
	*(uint32_t*)(v1 + 159) = -1;
	v2 = nox_xxx_serverOptionsGetServername_40A4C0();
	strncpy(v1 + 120, v2, 0xFu);
	strcpy(v1 + 111, nox_xxx_mapGetMapName_409B40());
	if (nox_xxx_isQuest_4D6F50()) {
		if (dword_5d4594_528256) {
			*(uint16_t*)(v1 + 165) = nox_game_getQuestStage_4E3CC0();
		} else {
			*(uint16_t*)(v1 + 165) = 1;
		}
	}
	v1[104] = nox_xxx_servGetPlrLimit_409FA0();
	v3 = nox_common_playerInfoCount_416F40();
	v1[103] = v3;
	if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING)) {
		v1[103] = v3 - 1;
		--v1[104];
	}
	v4 = sub_43BE50_get_video_mode_id();
	v5 = v1[102];
	*((uint32_t*)v1 + 12) = NOX_CLIENT_VERS_CODE;
	v1[102] = v5 & 0x80 | v4;
	*((uint32_t*)v1 + 11) = *getMemU32Ptr(0x5D4594, 814916);
	*(uint16_t*)(v1 + 109) = nox_xxx_servGetPort_40A430();
	nox_client_setServerConnectAddr_435720("localhost");
	if (0) {
		memset(v10, 0, sizeof(v10));
		v6 = sub_41FA40();
		nox_sprintf(v9, "%s%s", v6, getMemAt(0x587000, 90752));
		strcpy(&v10[52], v9);
		*(uint32_t*)v10 = sub_420100();
		*(uint32_t*)&v10[4] = 1;
		*(uint32_t*)&v10[8] = 32;
		*(uint32_t*)&v10[12] = 0;
		*(uint32_t*)&v10[16] = 0;
		*(uint32_t*)&v10[20] = 1;
		*(uint32_t*)&v10[24] = 1;
		*(uint32_t*)&v10[44] = 0;
		*(uint32_t*)&v10[28] = 0;
		*(uint32_t*)&v10[224] = NOX_CLIENT_VERS_CODE;
		*(uint32_t*)&v10[32] = *getMemU32Ptr(0x5D4594, 814916);
		v10[sub_425550(v1 + 100, &v10[69], 552) + 69] = 0;
	}
	result = nox_xxx_cliGamedataGet_416590(1);
	v8 = *((uint16_t*)result + 26) & 0xE90F;
	HIBYTE(v8) |= 1u;
	*((uint16_t*)result + 26) = v8;
	return result;
}

//----- (0043AF30) --------------------------------------------------------
int sub_43AF30() { return dword_5d4594_815052; }

//----- (0043AF40) --------------------------------------------------------
int sub_43AF40() { return nox_game_createOrJoin_815048; }

//----- (0043AF80) --------------------------------------------------------
int sub_43AF80() { return dword_5d4594_814548; }

//----- (0043AF90) --------------------------------------------------------
int sub_43AF90(int a1) {
	int result; // eax

	result = a1;
	dword_5d4594_814548 = a1;
	return result;
}

//----- (0043AFA0) --------------------------------------------------------
void nox_client_setConnError_43AFA0(int a1) {
	nox_client_connError_814552 = a1;
	sub_43AF90(2);
}

//----- (0043B300) --------------------------------------------------------
unsigned int nox_client_getServerAddr_43B300() {
	unsigned int result; // eax

	if (dword_5d4594_815056) {
		result = inet_addr((const char*)((uint32_t)dword_5d4594_814624 + 12));
	} else {
		result = 0;
	}
	return result;
}

//----- (0043B320) --------------------------------------------------------
int nox_client_getServerPort_43B320() { return dword_5d4594_815056 != 0 ? *getMemU32Ptr(0x5D4594, 814604) : 0; }

//----- (0043B340) --------------------------------------------------------
int sub_43B340() {
	int result; // eax

	if (dword_5d4594_815056) {
		result = *(unsigned short*)((uint32_t)dword_5d4594_814624 + 163);
	} else {
		result = 0;
	}
	return result;
}

//----- (0043B460) --------------------------------------------------------
int sub_43B460() {
	sub_438370();
	nox_wnd_xxx_815040->fnc_done_out = sub_43B490;
	nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_815000);
	return nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_814984, 0);
}

//----- (0043B490) --------------------------------------------------------
int sub_43B490() {
	if (nox_game_getStateCode_43BE10() == 1700) {
		return sub_438330();
	}
	nox_window_set_hidden(*(int*)&nox_wol_wnd_world_814980, 1);
	nox_window_set_hidden(*(int*)&dword_5d4594_815000, 1);
	nox_client_setCursorType_477610(0);
	return 1;
}
