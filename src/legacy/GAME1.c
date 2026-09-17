#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__system__parsecmd.h"
#include "common__log.h"
#include "common__object__armrlook.h"
#include "common__object__weaplook.h"
#include "common__system__team.h"

#include "client__gui__servopts__playrlst.h"
#include "client__system__ctrlevnt.h"
#include "client__video__draw_common.h"
#include "common__binfile.h"
#include "common__magic__speltree.h"
#include "common__net_list.h"
#include "common__random.h"
#include "common__strman.h"

#include "client__gui__window.h"
#include "client__io__win95__focus.h"
#include "common/fs/nox_fs.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_5d4594_527656;
extern uint32_t dword_5d4594_251712;
extern uint32_t dword_5d4594_251708;
extern uint32_t dword_5d4594_251716;
extern uint32_t dword_5d4594_10984;
extern uint32_t dword_5d4594_251720;
extern uint32_t dword_5d4594_251744;
extern uint32_t dword_5d4594_3484;
extern void* dword_5d4594_251560;
extern uint32_t nox_tile_def_cnt;
extern uint32_t dword_5d4594_251572;
extern uint32_t nox_player_netCode_85319C;

int nox_server_gameSettingsUpdated; // If you define it as 1-byte bool, the game will crash

extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern int ptr_5D4594_2650668_cap;

uint32_t nox_tile_def_cnt = 0;
nox_tileDef_t nox_tile_defs_arr[176] = {0};

//----- (00409470) --------------------------------------------------------
int nox_xxx_parseString_409470(FILE* a1, uint8_t* a2) {
	uint8_t* v2;          // ebx
	int v3;               // ecx
	int v4;               // edi
	int v5;               // ebp
	int v6;               // eax
	uint16_t CharType[2]; // [esp+10h] [ebp-4h]

	v2 = a2;
	v3 = 0;
	*(uint32_t*)CharType = 0;
	v4 = 1;
	do {
		while (1) {
			v5 = v3;
			nox_binfile_fread_408E40((char*)CharType, 1, 1, a1);
			if (nox_binfile_lastErr_409370(a1) == -1) {
				return 0;
			}
			v3 = *(uint32_t*)CharType;
			v6 = isspace(CharType[0]);
			if (v6) {
				break;
			}
			v4 = 0;
			if (v3 != '/' || v5 != '/') {
				*v2++ = v3;
			} else {
				nox_binfile_skipLine_409520(a1);
				v2 = a2;
				v3 = *(uint32_t*)CharType;
				v4 = 1;
			}
		}
	} while (v4);
	*v2 = 0;
	return 1;
}

//----- (00409A70) --------------------------------------------------------
int sub_409A70(short a1) {
	int result;        // eax
	unsigned char* v2; // ecx

	result = 0;
	v2 = getMemAt(0x587000, 4704);
	while (*(uint32_t*)v2 != (a1 & 0x17F0)) {
		v2 += 4;
		++result;
		if ((int)v2 >= (int)getMemAt(0x587000, 4728)) {
			return 0;
		}
	}
	return result;
}

//----- (00409B30) --------------------------------------------------------
char* nox_server_currentMapGetFilename_409B30() { return (char*)getMemAt(0x5D4594, 2598188); }

//----- (00409B40) --------------------------------------------------------
char* nox_xxx_mapGetMapName_409B40() { return (char*)getMemAt(0x85B3FC, 36); }

//----- (00409B50) --------------------------------------------------------
unsigned int sub_409B50(const char* a1) {
	unsigned int result; // eax

	result = strlen(a1) + 1;
	memcpy(getMemAt(0x5D4594, 3452), a1, result);
	return result;
}

//----- (00409B80) --------------------------------------------------------
char* sub_409B80() { return (char*)getMemAt(0x5D4594, 3452); }

//----- (00409D70) --------------------------------------------------------
char* nox_xxx_gameSetMapPath_409D70(char* a1) {
	char* result;  // eax
	char* v2;      // eax
	signed int v3; // esi

	result = (char*)nox_strcmpi((const char*)getMemAt(0x5D4594, 2598188), a1);
	if (result) {
		strncpy((char*)getMemAt(0x5D4594, 2598188), a1, 0x50u);
		*getMemU8Ptr(0x5D4594, 2598267) = 0;
		v2 = strrchr(a1, 92);
		if (v2) {
			v3 = strlen(v2 + 1) - 4;
			result = strncpy((char*)getMemAt(0x85B3FC, 36), v2 + 1, v3);
		} else {
			v3 = strlen((const char*)getMemAt(0x5D4594, 2598188)) - 4;
			if (v3 < 0) {
				v3 = 0;
			}
			result = strncpy((char*)getMemAt(0x85B3FC, 36), (const char*)getMemAt(0x5D4594, 2598188), v3);
		}
		*getMemU8Ptr(0x85B3FC, 36 + v3) = 0;
		nox_server_gameSettingsUpdated = 1;
	}
	return result;
}


































//----- (0040A770) --------------------------------------------------------
int sub_40A770() {
	int v0;   // edi
	char* v1; // esi
	int v2;   // eax
	int v3;   // ecx
	char* i;  // eax

	v0 = 0;
	if (!nox_xxx_CheckGameplayFlags_417DA0(4)) {
		for (i = nox_common_playerInfoGetFirst_416EA0(); i; i = nox_common_playerInfoGetNext_416EE0((int)i)) {
			if (!(i[3680] & 1) &&
				(i[2064] != 31 || !nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING))) {
				++v0;
			}
		}
		return v0;
	}
	v1 = nox_server_teamFirst_418B10();
	if (!v1) {
		return v0;
	}
	do {
		v2 = nox_xxx_getFirstPlayerUnit_4DA7C0();
		if (v2) {
			while (1) {
				v3 = *(uint32_t*)(*(uint32_t*)(v2 + 748) + 276);
				if (!(*(uint8_t*)(v3 + 3680) & 1) &&
					(*(uint8_t*)(v3 + 2064) != 31 ||
					 !nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING))) {
					break;
				}
				v2 = nox_xxx_getNextPlayerUnit_4DA7F0(v2);
				if (!v2) {
					goto LABEL_10;
				}
			}
			++v0;
		}
	LABEL_10:
		v1 = nox_server_teamNext_418B60((int)v1);
	} while (v1);
	return v0;
}

//----- (0040A830) --------------------------------------------------------
int nox_xxx_countNonEliminatedPlayersInTeam_40A830(nox_team_t* a1p) {
	int a1 = a1p;
	int v1; // edi
	int v2; // esi
	int v3; // eax

	v1 = 0;
	v2 = nox_xxx_getFirstPlayerUnit_4DA7C0();
	if (!v2) {
		return 0;
	}
	do {
		if (nox_xxx_teamCompare2_419180(v2 + 48, *(uint8_t*)(a1 + 57))) {
			v3 = *(uint32_t*)(*(uint32_t*)(v2 + 748) + 276);
			if (!(*(uint8_t*)(v3 + 3680) & 1) &&
				(*(uint8_t*)(v3 + 2064) != 31 ||
				 !nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING))) {
				++v1;
			}
		}
		v2 = nox_xxx_getNextPlayerUnit_4DA7F0(v2);
	} while (v2);
	return v1;
}

//----- (0040A8A0) --------------------------------------------------------
int nox_xxx_gamePlayIsAnyPlayers_40A8A0() {
	int v0;     // ebx
	char* v1;   // edi
	int v2;     // esi
	int v3;     // eax
	int result; // eax
	char* i;    // eax
	int v6;     // ecx

	v0 = 0;
	if (!nox_xxx_CheckGameplayFlags_417DA0(4)) {
		for (i = nox_common_playerInfoGetFirst_416EA0(); i; i = nox_common_playerInfoGetNext_416EE0((int)i)) {
			if (*((uint32_t*)i + 514)) {
				v6 = *((uint32_t*)i + 920);
				if (!(v6 & 1) || v6 & 0x20) {
					++v0;
				}
			}
		}
		result = v0 < 1;
		LOBYTE(result) = v0 > 1;
		return result;
	}
	v1 = nox_server_teamFirst_418B10();
	if (!v1) {
		result = v0 < 1;
		LOBYTE(result) = v0 > 1;
		return result;
	}
	do {
		v2 = nox_xxx_getFirstPlayerUnit_4DA7C0();
		if (v2) {
			while (1) {
				if (nox_xxx_teamCompare2_419180(v2 + 48, v1[57])) {
					v3 = *(uint32_t*)(*(uint32_t*)(*(uint32_t*)(v2 + 748) + 276) + 3680);
					if (!(v3 & 1) || v3 & 0x20) {
						break;
					}
				}
				v2 = nox_xxx_getNextPlayerUnit_4DA7F0(v2);
				if (!v2) {
					goto LABEL_10;
				}
			}
			++v0;
		}
	LABEL_10:
		v1 = nox_server_teamNext_418B60((int)v1);
	} while (v1);
	return v0 > 1;
}

//----- (0040A970) --------------------------------------------------------
void sub_40A970() {
	char* i;  // esi
	int v1;   // eax
	float v3; // [esp+0h] [ebp-8h]
	float v4; // [esp+0h] [ebp-8h]

	*getMemU32Ptr(0x5D4594, 3520) = gameFrame();
	*getMemU32Ptr(0x5D4594, 3536) = 0;
	v3 = nox_xxx_gamedataGetFloat_419D40("SuddenDeathPlayerThreshold");
	*getMemU32Ptr(0x5D4594, 3476) = nox_float2int(v3);
	v4 = nox_xxx_gamedataGetFloat_419D40("SuddenDeathLifeTime");
	*getMemU32Ptr(0x5D4594, 1392) = nox_float2int(v4);
	for (i = nox_common_playerInfoGetFirst_416EA0(); i; i = nox_common_playerInfoGetNext_416EE0((int)i)) {
		v1 = *((uint32_t*)i + 920);
		if (v1 & 0x100) {
			nox_xxx_playerUnsetStatus_417530((int)i, 256);
		}
	}
	nox_common_gameFlags_unset_40A540(0x4000000);
}

//----- (0040AA00) --------------------------------------------------------
int sub_40AA00() { return 20 * gameFPS() < (unsigned int)(gameFrame() - *getMemU32Ptr(0x5D4594, 3520)); }

//----- (0040AA20) --------------------------------------------------------
int sub_40AA20() { return *getMemU32Ptr(0x5D4594, 3536); }

//----- (0040AA30) --------------------------------------------------------
int sub_40AA30(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 3536) = a1;
	return result;
}

//----- (0040AA40) --------------------------------------------------------
int sub_40AA40() { return *getMemU32Ptr(0x5D4594, 3476); }

//----- (0040AA60) --------------------------------------------------------
int sub_40AA60(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 3508) = a1;
	return result;
}

//----- (0040AA70) --------------------------------------------------------
int sub_40AA70(nox_playerInfo* pl) {
	int a1 = pl;
	char* v1;   // edi
	int result; // eax
	int v3;     // eax
	int v4;     // esi
	int v5;     // ebx
	char* v6;   // eax

	v1 = sub_416640();
	if (!a1) {
		goto LABEL_31;
	}
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		return sub_40A770() < 6;
	}
	v3 = *(uint32_t*)(a1 + 3680);
	if (v3 & 0x100 && !*getMemU32Ptr(0x5D4594, 3508)) {
		return 0;
	}
	if (!sub_40A740() && !nox_common_gameFlags_check_40A5C0(0x8000)) {
		goto LABEL_31;
	}
	result = *(uint32_t*)(a1 + 2068);
	if (!result) {
		return result;
	}
	if (nox_server_teamByXxx_418AE0(*(uint32_t*)(a1 + 2068))) {
		goto LABEL_31;
	}
	v4 = (unsigned char)v1[52];
	if ((nox_common_gameFlags_check_40A5C0(96) ||
		 nox_common_gameFlags_check_40A5C0(16) && nox_xxx_CheckGameplayFlags_417DA0(4)) &&
		v4 > 2) {
		v4 = 2;
	}
	if ((unsigned char)sub_417DE0() >= v4) {
		return 0;
	}
	if (nox_common_gameFlags_check_40A5C0(96)) {
		v5 = (unsigned char)sub_417DE0();
		if (v5 >= sub_417DC0()) {
			return 0;
		}
	}
LABEL_31:
	if (nox_common_gameFlags_check_40A5C0(128)) {
		return 1;
	}
	if (!nox_common_gameFlags_check_40A5C0(1024)) {
		return 1;
	}
	v6 = nox_common_playerInfoGetFirst_416EA0();
	if (!v6) {
		return 1;
	}
	while (*((int*)v6 + 535) <= 0) {
		v6 = nox_common_playerInfoGetNext_416EE0((int)v6);
		if (!v6) {
			return 1;
		}
	}
	if (!sub_40AA00()) {
		return 1;
	}
	return 0;
}

//----- (0040E090) --------------------------------------------------------
void sub_40E090() { dword_5d4594_10984 = 0; }

//----- (00410360) --------------------------------------------------------
uint8_t* nox_xxx_doorAttachWall_410360(int a1, int a2, int a3) {
	uint8_t* result; // eax
	char v4;         // cl

	result = nox_xxx_wallCreateAt_410250(a2, a3);
	if (result) {
		v4 = result[4] | 0x10;
		*((uint32_t*)result + 7) = a1;
		result[4] = v4;
	}
	return result;
}

//----- (00410390) --------------------------------------------------------
uint32_t* sub_410390(int a1, int a2, int a3) {
	uint32_t* v3;     // esi
	uint32_t* result; // eax
	int v5;           // edx
	char v6;          // cl
	int v7[2];        // [esp+Ch] [ebp-8h]

	v3 = (uint32_t*)nox_xxx_wall_4105E0(a2, a3);
	if (v3 || (result = nox_xxx_wallCreateAt_410250(a2, a3), (v3 = result) != 0)) {
		v5 = *(uint32_t*)(a1 + 16);
		v7[0] = *(uint32_t*)(a1 + 12);
		v6 = *((uint8_t*)v3 + 4);
		v3[8] = a1;
		v7[1] = v5;
		*((uint8_t*)v3 + 4) = v6 | 0x10;
		result = nox_xxx_polygonIsPlayerInPolygon_4217B0((int2*)v7, 0);
		if (result || (result = sub_421990((int2*)v7, 10.0, 0)) != 0) {
			*((uint8_t*)v3 + 8) = *((uint8_t*)result + 130);
		} else {
			*((uint8_t*)v3 + 8) = 1;
		}
	}
	return result;
}

//----- (00410550) --------------------------------------------------------
int sub_410550(short a1) {
	int* v1; // eax

	v1 = (int*)nox_xxx_wallSecretGetFirstWall_410780();
	if (!v1) {
		return 0;
	}
	while (*(uint16_t*)(v1[3] + 10) != a1) {
		v1 = (int*)nox_xxx_wallSecretNext_410790(v1);
		if (!v1) {
			return 0;
		}
	}
	return v1[3];
}

//----- (00410730) --------------------------------------------------------
uint32_t* sub_410730() {
	uint32_t* result; // eax
	uint32_t* v1;     // esi

	result = dword_5d4594_251560;
	if (dword_5d4594_251560) {
		do {
			v1 = (uint32_t*)*result;
			free(result);
			result = v1;
		} while (v1);
		dword_5d4594_251560 = 0;
	} else {
		dword_5d4594_251560 = 0;
	}
	return result;
}

//----- (00410760) --------------------------------------------------------
uint32_t* nox_xxx_wallSecretBlock_410760(uint32_t* a1) {
	uint32_t* result; // eax

	result = a1;
	*a1 = dword_5d4594_251560;
	dword_5d4594_251560 = a1;
	return result;
}

//----- (00410780) --------------------------------------------------------
void* nox_xxx_wallSecretGetFirstWall_410780() { return dword_5d4594_251560; }

//----- (00410790) --------------------------------------------------------
int nox_xxx_wallSecretNext_410790(int* a1) {
	int result; // eax

	if (a1) {
		result = *a1;
	} else {
		result = 0;
	}
	return result;
}

//----- (004107A0) --------------------------------------------------------
int* sub_4107A0(void* lpMem) {
	int* result; // eax
	int* v2;     // esi

	result = dword_5d4594_251560;
	v2 = 0;
	if (dword_5d4594_251560) {
		while (result != lpMem) {
			v2 = result;
			result = (int*)nox_xxx_wallSecretNext_410790(result);
			if (!result) {
				return result;
			}
		}
		if (result == dword_5d4594_251560) {
			dword_5d4594_251560 = nox_xxx_wallSecretNext_410790(result);
		} else {
			*v2 = nox_xxx_wallSecretNext_410790(result);
		}
		free(lpMem);
	}
	return result;
}

//----- (00410F60) --------------------------------------------------------
int nox_xxx_tileAlloc_410F60_init() {
	ptr_5D4594_2650668 = calloc(ptr_5D4594_2650668_cap, sizeof(void*));
	if (!ptr_5D4594_2650668) {
		return 0;
	}

	for (int i = 0; i < ptr_5D4594_2650668_cap; i++) {
		ptr_5D4594_2650668[i] = (obj_5D4594_2650668_t*)calloc(ptr_5D4594_2650668_cap, sizeof(obj_5D4594_2650668_t));
		if (!ptr_5D4594_2650668[i]) {
			return 0;
		}
	}
	return 1;
}

//----- (00410FC0) --------------------------------------------------------
void nox_xxx_tileFree_410FC0_free() {
	for (int i = 0; i < ptr_5D4594_2650668_cap; i++) {
		obj_5D4594_2650668_t* ptr = ptr_5D4594_2650668[i];
		if (ptr) {
			free(ptr);
		}
	}
}

//----- (00411160) --------------------------------------------------------
int nox_xxx_tileNFromPoint_411160(float2* a1) {
	float v12 = (a1->field_0 + 11.5) * 0.021739131;
	float v13 = (a1->field_4 + 11.5) * 0.021739131;

	int i = nox_float2int(v12);
	int j = nox_float2int(v13);

	float v14 = a1->field_0 + 11.5;
	float v15 = a1->field_4 + 11.5;

	int v4 = nox_float2int(v14) % 46;
	int v5 = nox_float2int(v15) % 46;

	if (i <= 1 || i >= 127 || j <= 1 || j >= 127) {
		return -1;
	}

	int result = 0;
	int v16[2] = {0};
	if (v4 <= v5) {
		if (46 - v4 <= v5) {
			obj_5D4594_2650668_t* obj = &ptr_5D4594_2650668[i][j];
			result = obj->field_6;
			if (obj->field_10) {
				v16[0] = v4;
				v16[1] = v5 - 23;
				result = sub_411350(obj->field_10, v16, result);
			}
		} else {
			obj_5D4594_2650668_t* obj = &ptr_5D4594_2650668[i - 1][j];
			result = obj->field_1;
			if (obj->field_5) {
				v16[1] = v5;
				v16[0] = v4 + 23;
				result = sub_411350(obj->field_5, v16, result);
			}
		}
	} else if (46 - v4 <= v5) {
		obj_5D4594_2650668_t* obj = &ptr_5D4594_2650668[i][j];
		result = obj->field_1;
		if (obj->field_5) {
			v16[1] = v5;
			v16[0] = v4 - 23;
			result = sub_411350(obj->field_5, v16, result);
		}
	} else {
		obj_5D4594_2650668_t* obj = &ptr_5D4594_2650668[i][j - 1];
		result = obj->field_6;
		if (obj->field_10) {
			v16[0] = v4;
			v16[1] = v5 + 23;
			result = sub_411350(obj->field_10, v16, result);
		}
	}
	return result;
}

//----- (00411A90) --------------------------------------------------------
int nox_xxx_mapTileAllowTeleport_411A90(float2* a1) {
	if (*getMemIntPtr(0x587000, 26520) == -1) {
		for (int i = 0; i < 176; i++) {
			nox_tileDef_t* p = &nox_tile_defs_arr[i];
			if (!strcmp(&p->name[0], "WaterNoTeleport")) {
				*getMemU32Ptr(0x587000, 26516) = i;
			} else if (!strcmp(&p->name[0], "WaterDeepNoTeleport")) {
				*getMemU32Ptr(0x587000, 26520) = i;
			} else if (!strcmp(&p->name[0], "WaterShallowNoTeleport")) {
				*getMemU32Ptr(0x587000, 26524) = i;
			} else if (!strcmp(&p->name[0], "WaterSwampDeepNoTeleport")) {
				*getMemU32Ptr(0x587000, 26528) = i;
			} else if (!strcmp(&p->name[0], "WaterSwampShallowNoTeleport")) {
				*getMemU32Ptr(0x587000, 26532) = i;
			}
		}
	}
	int v3 = nox_xxx_tileNFromPoint_411160(a1);
	return v3 == *getMemU32Ptr(0x587000, 26516) || v3 == *getMemU32Ptr(0x587000, 26520) ||
		   v3 == *getMemU32Ptr(0x587000, 26524) || v3 == *getMemU32Ptr(0x587000, 26528) ||
		   v3 == *getMemU32Ptr(0x587000, 26532);
}

//----- (004133C0) --------------------------------------------------------
obj_412ae0_t* nox_xxx_modifNext_4133C0(obj_412ae0_t* a1) { return a1->field_34; }

//----- (004133D0) --------------------------------------------------------
int sub_4133D0(nox_object_t* a1p) {
	int a1 = a1p;
	int v1; // eax
	int v2; // esi
	int v3; // eax

	v1 = *getMemU32Ptr(0x5D4594, 251620);
	v2 = *(uint32_t*)(a1 + 692);
	if (!*getMemU32Ptr(0x5D4594, 251620)) {
		v3 = nox_xxx_modifGetIdByName_413290("Material7");
		v1 = nox_xxx_modifGetDescById_413330(v3);
		*getMemU32Ptr(0x5D4594, 251620) = v1;
	}
	return *(uint32_t*)(a1 + 8) & 0x13001000 && *(uint32_t*)(v2 + 4) == v1;
}

//----- (00413420) --------------------------------------------------------
int sub_413420(char a1) {
	unsigned char* v1; // esi
	int v2;            // ecx
	unsigned char* v3; // eax

	if (!*getMemU32Ptr(0x5D4594, 251624)) {
		v1 = getMemAt(0x587000, 27340);
		do {
			*(uint32_t*)v1 = nox_xxx_gLoadImg_42F970(*((const char**)v1 - 1));
			v1 += 20;
		} while ((int)v1 < (int)getMemAt(0x587000, 27460));
		*getMemU32Ptr(0x5D4594, 251624) = 1;
	}
	v2 = 0;
	v3 = getMemAt(0x587000, 27332);
	while (*v3 != a1) {
		v3 += 20;
		++v2;
		if ((int)v3 >= (int)getMemAt(0x587000, 27452)) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 27340 + 20 * v2);
}

//----- (004134D0) --------------------------------------------------------
int sub_4134D0() {
	dword_5d4594_251708 = 0;
	dword_5d4594_251712 = 0;
	dword_5d4594_251716 = 0;
	dword_5d4594_251720 = 0;
	return 0;
}

//----- (004137E0) --------------------------------------------------------
void sub_4137E0() {
	if (dword_5d4594_251720) {
		dword_5d4594_251720 = 0;
		sub_43DBE0();
	}
}

//----- (00413890) --------------------------------------------------------
char* sub_413890() {
	unsigned char* v0; // edi
	int v1;            // ecx
	short v2;          // dx
	unsigned char v3;  // al

	*getMemU8Ptr(0x5D4594, 251636) = 0;
	*getMemU8Ptr(0x5D4594, 251637) = 0;
	if (!getMemByte(0x5D4594, 251636)) {
		return 0;
	}
	v0 = getMemAt(0x5D4594, 251637 + strlen((const char*)getMemAt(0x5D4594, 251636)));
	v1 = *getMemU32Ptr(0x587000, 32316);
	v2 = *getMemU16Ptr(0x587000, 32320);
	*(uint32_t*)--v0 = *getMemU32Ptr(0x587000, 32312);
	v3 = getMemByte(0x587000, 32322);
	*((uint32_t*)v0 + 1) = v1;
	*((uint16_t*)v0 + 4) = v2;
	v0[10] = v3;
	return (char*)getMemAt(0x5D4594, 251636);
}

//----- (004138E0) --------------------------------------------------------
void sub_4138E0(int a1) {
	*getMemU32Ptr(0x5D4594, 251740) = nox_xxx_checkGameFlagPause_413A50();
	sub_413A00(1);
}

//----- (00413900) --------------------------------------------------------
void sub_413900(int a1) {
	if (!nox_video_inFadeTransition_44E0D0()) {
		if (!*getMemU32Ptr(0x5D4594, 251740)) {
			sub_413A00(0);
		}
	}
}

//----- (00413920) --------------------------------------------------------
int sub_413920() {
	sub_42EBB0(1u, sub_413900, 0, "Pause");
	sub_42EBB0(2u, sub_4138E0, 0, "Pause");
	dword_5d4594_251744 = 0;
	return 1;
}

//----- (004139B0) --------------------------------------------------------
int sub_4139B0() { return dword_5d4594_251744 != 0; }

//----- (00413A50) --------------------------------------------------------
int nox_xxx_checkGameFlagPause_413A50() { return nox_common_gameFlags_check_40A5C0(0x40000); }

//----- (00413E30) --------------------------------------------------------
void nox_xxx_gameLoopMemDump_413E30() {
	signed int v0;     // ebx
	int v1;            // esi
	unsigned char* v2; // eax
	int v3;            // edx
	signed int v4;     // edi
	const char* v5;    // ebp
	unsigned char* v6; // eax
	signed int v7;     // ecx
	int v8;            // ebx
	signed int v9;     // [esp+8h] [ebp-8h]

	v0 = 0;
	v1 = 0;
	v9 = 0;
	v2 = getMemAt(0x5D4594, 252364);
	do {
		*(uint32_t*)v2 = 0;
		v2 += 84;
	} while ((int)v2 < (int)getMemAt(0x5D4594, 338380));
	v3 = *getMemU32Ptr(0x5D4594, 338304);
	if (*getMemU32Ptr(0x5D4594, 338304)) {
		do {
			v4 = 0;
			if (v0 > 0) {
				v5 = (const char*)getMemAt(0x5D4594, 252284);
				while (strcmp(v5, (const char*)(v3 + 20))) {
					v0 = v9;
					++v4;
					v5 += 84;
					if (v4 >= v9) {
						goto LABEL_10;
					}
				}
				v0 = v9;
				*getMemU32Ptr(0x5D4594, 252364 + 84 * v4) += *(uint32_t*)(v3 + 16);
			}
		LABEL_10:
			if (v4 == v0) {
				*getMemU32Ptr(0x5D4594, 252364 + 84 * v4) = *(uint32_t*)(v3 + 16);
				strcpy((char*)getMemAt(0x5D4594, 252284 + 84 * v4), (const char*)(v3 + 20));
				v9 = ++v0;
			}
			v3 = *(uint32_t*)(v3 + 4);
		} while (v3);
		v1 = 0;
	}
	qsort(getMemAt(0x5D4594, 252284), v0, 0x54u, sub_413F60);
	if (v0 > 0) {
		v6 = getMemAt(0x5D4594, 252364);
		v7 = v0;
		do {
			v8 = *(uint32_t*)v6;
			v6 += 84;
			v1 += v8;
			--v7;
		} while (v7);
	}
}

//----- (00413F60) --------------------------------------------------------
int sub_413F60(const void* a1, const void* a2) { return *((uint32_t*)a1 + 20) - *((uint32_t*)a2 + 20); }

//----- (00414C90) --------------------------------------------------------
char nox_xxx_initSinCosTables_414C90() {
	long long v0;      // rax
	unsigned char* v1; // esi
	unsigned char* v2; // esi
	int v4;            // [esp+0h] [ebp-4h]
	int v5;            // [esp+0h] [ebp-4h]

	LOBYTE(v0) = getMemByte(0x5D4594, 371240);
	if (!getMemByte(0x5D4594, 371240)) {
		*getMemU8Ptr(0x5D4594, 371240) = 1;
		v4 = 0;
		v1 = getMemAt(0x85B3FC, 12260);
		do {
			*(uint32_t*)v1 = (long long)(sin((double)v4 * 0.0015339808) * *getMemDoublePtr(0x581450, 7200));
			v1 += 4;
			++v4;
		} while ((int)v1 < (int)getMemAt(0x85B3FC, 28644));
		v5 = 0;
		v2 = getMemAt(0x5D4594, 338472);
		do {
			v0 = (long long)(acos((double)v5 * 0.00024414062 - 1.0) * *getMemDoublePtr(0x581450, 7184));
			*(uint32_t*)v2 = v0;
			v2 += 4;
			++v5;
		} while ((int)v2 < (int)getMemAt(0x5D4594, 371240));
	}
	return v0;
}

//----- (00415960) --------------------------------------------------------
int sub_415960(wchar2_t* a1) {
	int v1;             // edi
	const wchar2_t** v2; // eax
	unsigned char* v3;  // esi
	int v4;             // ecx

	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 33392)) {
		return 0;
	}
	v2 = (const wchar2_t**)getMemAt(0x587000, 33392);
	v3 = getMemAt(0x587000, 33392);
	while (_nox_wcsicmp(a1, *v2)) {
		v4 = *((uint32_t*)v3 + 3);
		v3 += 12;
		++v1;
		v2 = (const wchar2_t**)v3;
		if (!v4) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 33400 + 12 * v1);
}

//----- (004159F0) --------------------------------------------------------
int sub_4159F0(int a1) {
	int v1;           // ecx
	unsigned char* i; // eax
	int v3;           // esi

	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 33392)) {
		return 0;
	}
	for (i = getMemAt(0x587000, 33392); *((uint32_t*)i + 2) != a1; i += 12) {
		v3 = *((uint32_t*)i + 3);
		++v1;
		if (!v3) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 33392 + 12 * v1);
}

//----- (00415BD0) --------------------------------------------------------
double sub_415BD0(int a1) {
	float* v1;     // eax
	double result; // st7

	if (*(uint32_t*)(a1 + 8) & 0x2000000 &&
		(v1 = (float*)nox_xxx_equipClothFindDefByTT_413270(*(unsigned short*)(a1 + 4))) != 0) {
		result = v1[16];
	} else {
		result = 0.0;
	}
	return result;
}

//----- (00415DA0) --------------------------------------------------------
int sub_415DA0(wchar2_t* a1) {
	int v1;             // edi
	const wchar2_t** v2; // eax
	unsigned char* v3;  // esi
	int v4;             // ecx

	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 35496)) {
		return 0;
	}
	v2 = (const wchar2_t**)getMemAt(0x587000, 35496);
	v3 = getMemAt(0x587000, 35496);
	while (_nox_wcsicmp(a1, *v2)) {
		v4 = *((uint32_t*)v3 + 3);
		v3 += 12;
		++v1;
		v2 = (const wchar2_t**)v3;
		if (!v4) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 35504 + 12 * v1);
}

//----- (00415E80) --------------------------------------------------------
int sub_415E80(int a1) {
	int v1;           // ecx
	unsigned char* i; // eax
	int v3;           // esi

	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 35496)) {
		return 0;
	}
	for (i = getMemAt(0x587000, 35496); *((uint32_t*)i + 2) != a1; i += 12) {
		v3 = *((uint32_t*)i + 3);
		++v1;
		if (!v3) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 35496 + 12 * v1);
}

// 4161E0: using guessed type int var_14[5];










// 416690: using guessed type char var_54[84];















//----- (00416E50) --------------------------------------------------------
char* nox_xxx_playerForceSendLessons_416E50(int a1) {
	char* result; // eax
	int* i;       // esi

	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int*)result; result; i = (int*)result) {
		i[534] = 0;
		i[535] = 0;
		if (a1) {
			if (i[514]) {
				nox_xxx_netReportLesson_4D8EF0(i[514]);
			}
		}
		result = nox_common_playerInfoGetNext_416EE0((int)i);
	}
	return result;
}

//----- (004170D0) --------------------------------------------------------
char* nox_xxx_playerByName_4170D0(wchar2_t* a1) {
	char* v1; // esi

	if (!a1) {
		return 0;
	}
	v1 = nox_common_playerInfoGetFirst_416EA0();
	if (!v1) {
		return 0;
	}
	while (_nox_wcsicmp((const wchar2_t*)v1 + 2352, a1)) {
		v1 = nox_common_playerInfoGetNext_416EE0((int)v1);
		if (!v1) {
			return 0;
		}
	}
	return v1;
}

//----- (004173D0) --------------------------------------------------------
int nox_xxx_playerMapTracksObj_4173D0(int a1, nox_object_t* a2p) {
	int a2 = a2p;
	int result; // eax
	int v3;     // ecx

	result = 0;
	if (a1 >= 0 && a1 < NOX_PLAYERINFO_MAX) {
		if (a2) {
			nox_playerInfo* pl = nox_common_playerInfoFromNumRaw(a1);
			v3 = pl->field_4580;
			if (v3) {
				while (*(uint32_t*)(v3 + 4) != a2) {
					v3 = *(uint32_t*)(v3 + 8);
					if (v3 == pl->field_4580 || !v3) {
						return result;
					}
				}
				result = 1;
			}
		}
	}
	return result;
}

//----- (00417470) --------------------------------------------------------
char* nox_xxx_netUnmarkMinimapSpec_417470(int a1, int a2) {
	char* result; // eax
	int i;        // esi

	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		nox_xxx_netUnmarkMinimapObj_417300(*(unsigned char*)(i + 2064), a1, a2);
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
	return result;
}

//----- (004174B0) --------------------------------------------------------
char* nox_xxx_netMarkMinimapForAll_4174B0(int a1, int a2) {
	char* result; // eax
	int i;        // esi

	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		nox_xxx_netMarkMinimapObject_417190(*(unsigned char*)(i + 2064), a1, a2);
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
	return result;
}

//----- (004174F0) --------------------------------------------------------
int nox_xxx_netNeedTimestampStatus_4174F0(nox_playerInfo* a1p, int a2) {
	int a1 = a1p;
	int result; // eax

	*(uint32_t*)(a1 + 3680) |= a2;
	result = nox_common_gameFlags_check_40A5C0(1);
	if (result) {
		if (a2 & 0x423) {
			result = nox_xxx_netReportPlayerStatus_417630(a1);
		}
	}
	return result;
}

//----- (00417530) --------------------------------------------------------
char nox_xxx_playerUnsetStatus_417530(nox_playerInfo* a1p, int a2) {
	int a1 = a1p;
	int v2;   // eax
	short v3; // ax

	*(uint32_t*)(a1 + 3680) &= ~a2;
	v2 = nox_common_gameFlags_check_40A5C0(1);
	if (v2) {
		if (a2 & 1) {
			v2 = nox_common_gameFlags_check_40A5C0(128);
			if (!v2) {
				v2 = nox_xxx_gamePlayIsAnyPlayers_40A8A0();
				if (v2) {
					v2 = sub_40A220();
					if (!v2) {
						v3 = nox_common_gameFlags_getVal_40A5B0();
						LOBYTE(v2) = sub_40A180(v3);
						if ((uint8_t)v2) {
							sub_40A250();
							LOBYTE(v2) = sub_40A1F0(1);
						}
					}
				}
			}
		}
		if (a2 & 0x423) {
			LOBYTE(v2) = nox_xxx_netReportPlayerStatus_417630(a1);
		}
	}
	return v2;
}
// 417577: variable 'v2' is possibly undefined

//----- (004175C0) --------------------------------------------------------
char* nox_xxx_sendAllClientStatus_4175C0(int a1) {
	char* result; // eax
	int i;        // esi
	int v3;       // [esp-18h] [ebp-24h]
	char v4[7];   // [esp+4h] [ebp-8h]

	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		v4[0] = 106;
		*(uint16_t*)&v4[1] = *(uint16_t*)(i + 2060);
		v3 = *(unsigned char*)(a1 + 2064);
		*(uint32_t*)&v4[3] = *(uint32_t*)(i + 3680) & 0x423;
		nox_xxx_netSendPacket1_4E5390(v3, (int)v4, 7, 0, 0);
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
	return result;
}

//----- (00417630) --------------------------------------------------------
int nox_xxx_netReportPlayerStatus_417630(nox_playerInfo* pl) {
	int a1 = pl;
	short v1;   // cx
	int v2;     // edx
	char v4[7]; // [esp+0h] [ebp-8h]

	v1 = *(uint16_t*)(a1 + 2060);
	v2 = *(uint32_t*)(a1 + 3680) & 0x423;
	v4[0] = 106;
	*(uint16_t*)&v4[1] = v1;
	*(uint32_t*)&v4[3] = v2;
	return nox_xxx_netSendPacket1_4E5390(255, (int)v4, 7, 0, 0);
}

//----- (00417680) --------------------------------------------------------
void nox_xxx_cliPlayerRespawn_417680(int a1, char a2) {
	int v2;       // esi
	int v3;       // eax
	int v4;       // ecx
	uint32_t* v5; // edi
	int v6;       // eax
	int v7;       // ecx
	uint32_t* v8; // edi
	int v9;       // eax
	int v10;      // edi
	int v11;      // eax
	int v12;      // edx
	char v13;     // dl
	int v14;      // eax
	char v15;     // al
	int v16;      // eax
	int v17;      // eax
	int v18;      // eax
	int v19;      // eax
	char v20;     // cl
	int v21;      // eax
	int v22;      // eax
	char v23;     // cl
	int v24;      // eax
	int v25;      // [esp-Ch] [ebp-18h]
	int v26;      // [esp-Ch] [ebp-18h]

	v2 = a1;
	if (!a1) {
		return;
	}
	if (!nox_common_gameFlags_check_40A5C0(1)) {
		*(uint32_t*)(v2 + 4) = 0;
	}
	v3 = v2 + 2328;
	v4 = 27;
	do {
		v5 = (uint32_t*)v3;
		*(uint32_t*)(v3 - 4) = 0;
		v3 += 24;
		*v5 = 0;
		--v4;
		v5[1] = 0;
		v5[2] = 0;
		v5[3] = 0;
	} while (v4);
	if (!nox_common_gameFlags_check_40A5C0(1)) {
		*(uint32_t*)v2 = 0;
	}
	v6 = v2 + 2976;
	v7 = 26;
	do {
		v8 = (uint32_t*)v6;
		*(uint32_t*)(v6 - 4) = 0;
		v6 += 24;
		*v8 = 0;
		--v7;
		v8[1] = 0;
		v8[2] = 0;
		v8[3] = 0;
	} while (v7);
	v9 = nox_xxx_modifGetIdByName_413290("UserColor1");
	v10 = nox_xxx_modifGetDescById_413330(v9);
	if (!v10) {
		return;
	}
	if (*(uint8_t*)(v2 + 2251) || nox_common_gameFlags_check_40A5C0(2048)) {
		LOBYTE(a1) = -1;
		v11 = nox_xxx_modifGetDescById_413330(*(uint32_t*)(v10 + 4) + *(unsigned char*)(v2 + 2269));
		v12 = *(unsigned char*)(v2 + 2270);
		BYTE1(a1) = *(uint8_t*)(v11 + 4);
		BYTE2(a1) = *((uint8_t*)nox_xxx_modifGetDescById_413330(*(uint32_t*)(v10 + 4) + v12) + 4);
		HIBYTE(a1) = -1;
		if (a2 & 1) {
			nox_xxx_clientEquipWeaponArmor_417AA0(82, *(uint32_t*)(v2 + 2060), 1024, (int)&a1);
		}
	}
	LOBYTE(a1) = -1;
	BYTE1(a1) = *((uint8_t*)nox_xxx_modifGetDescById_413330(*(uint32_t*)(v10 + 4) + *(unsigned char*)(v2 + 2268)) + 4);
	HIWORD(a1) = -1;
	if (a2 & 2) {
		nox_xxx_clientEquipWeaponArmor_417AA0(82, *(uint32_t*)(v2 + 2060), 4, (int)&a1);
	}
	v13 = *((uint8_t*)nox_xxx_modifGetDescById_413330(*(uint32_t*)(v10 + 4) + *(unsigned char*)(v2 + 2272)) + 4);
	v14 = *(unsigned char*)(v2 + 2271);
	LOBYTE(a1) = v13;
	BYTE1(a1) = *((uint8_t*)nox_xxx_modifGetDescById_413330(*(uint32_t*)(v10 + 4) + v14) + 4);
	HIWORD(a1) = -1;
	if (a2 & 4) {
		nox_xxx_clientEquipWeaponArmor_417AA0(82, *(uint32_t*)(v2 + 2060), 1, (int)&a1);
	}
	v15 = *(uint8_t*)(v2 + 2251);
	a1 = -1;
	if (v15 == 1) {
		if (nox_common_gameFlags_check_40A5C0(2048)) {
			if (a2 & 8) {
				v16 = nox_xxx_modifGetIdByName_413290("ArmorQuality1");
				LOBYTE(a1) = *((uint8_t*)nox_xxx_modifGetDescById_413330(v16) + 4);
				nox_xxx_clientEquipWeaponArmor_417AA0(80, *(uint32_t*)(v2 + 2060), 0x8000, (int)&a1);
			}
		} else if (nox_common_gameFlags_check_40A5C0(4096)) {
			a1 = -1;
			v17 = nox_xxx_modifGetIdByName_413290("Replenishment1");
			BYTE2(a1) = *((uint8_t*)nox_xxx_modifGetDescById_413330(v17) + 4);
			nox_xxx_clientEquipWeaponArmor_417AA0(80, *(uint32_t*)(v2 + 2060), 0x10000, (int)&a1);
		} else if (a2 & 0x10) {
			nox_xxx_clientEquipWeaponArmor_417AA0(79, *(uint32_t*)(v2 + 2060), 0x4000, (int)&a1);
		}
	}
	if (!*(uint8_t*)(v2 + 2251)) {
		if (nox_common_gameFlags_check_40A5C0(2048)) {
			if (a2 & 0x20) {
				v18 = nox_xxx_modifGetIdByName_413290("ArmorQuality1");
				LOBYTE(a1) = *((uint8_t*)nox_xxx_modifGetDescById_413330(v18) + 4);
				v19 = nox_xxx_modifGetIdByName_413290("Material1");
				v20 = *((uint8_t*)nox_xxx_modifGetDescById_413330(v19) + 4);
				v21 = *(uint32_t*)(v2 + 2060);
				LOBYTE(a1) = v20;
				nox_xxx_clientEquipWeaponArmor_417AA0(80, v21, 256, (int)&a1);
			}
		} else if (nox_common_gameFlags_check_40A5C0(4096)) {
			v25 = *(uint32_t*)(v2 + 2060);
			a1 = -1;
			nox_xxx_clientEquipWeaponArmor_417AA0(80, v25, 256, (int)&a1);
		} else {
			if (a2 & 0x40) {
				nox_xxx_clientEquipWeaponArmor_417AA0(80, *(uint32_t*)(v2 + 2060), 512, (int)&a1);
			}
			if (a2 < 0) {
				nox_xxx_clientEquipWeaponArmor_417AA0(79, *(uint32_t*)(v2 + 2060), 0x1000000, (int)&a1);
			}
		}
	}
	if (*(uint8_t*)(v2 + 2251) == 2) {
		if (nox_common_gameFlags_check_40A5C0(2048)) {
			if (a2 & 8) {
				v22 = nox_xxx_modifGetIdByName_413290("ArmorQuality1");
				v23 = *((uint8_t*)nox_xxx_modifGetDescById_413330(v22) + 4);
				v24 = *(uint32_t*)(v2 + 2060);
				LOBYTE(a1) = v23;
				nox_xxx_clientEquipWeaponArmor_417AA0(80, v24, 0x8000, (int)&a1);
			}
		} else if (nox_common_gameFlags_check_40A5C0(4096)) {
			v26 = *(uint32_t*)(v2 + 2060);
			a1 = -1;
			nox_xxx_clientEquipWeaponArmor_417AA0(80, v26, 4, (int)&a1);
		}
	}
}

//----- (00417AA0) --------------------------------------------------------
char* nox_xxx_clientEquipWeaponArmor_417AA0(char a1, int a2, int a3, int a4) {
	char* result; // eax
	char* v5;     // edx
	int v6;       // ecx
	int v7;       // esi
	char* v8;     // edi
	char* v9;     // edx
	int v10;      // ecx
	int v11;      // esi
	char* v12;    // edi

	result = nox_common_playerInfoGetByID_417040(a2);
	if (result) {
		if (a1 == 81 || a1 == 80) {
			v9 = result + 2324;
			*((uint32_t*)result + 1) |= a3;
			v10 = 0;
			while (*(uint32_t*)v9) {
				++v10;
				v9 += 24;
				if (v10 >= 27) {
					return result;
				}
			}
			*(uint32_t*)&result[24 * v10 + 2324] = a3;
			v11 = 0;
			v12 = &result[24 * v10 + 2328];
			do {
				result = (char*)nox_xxx_modifGetDescById_413330(*(unsigned char*)(v11 + a4));
				*(uint32_t*)v12 = result;
				++v11;
				v12 += 4;
			} while (v11 < 4);
		} else {
			v5 = result + 2972;
			*(uint32_t*)result |= a3;
			v6 = 0;
			while (*(uint32_t*)v5) {
				++v6;
				v5 += 24;
				if (v6 >= 26) {
					return result;
				}
			}
			*(uint32_t*)&result[24 * v6 + 2972] = a3;
			v7 = 0;
			v8 = &result[24 * v6 + 2976];
			do {
				result = (char*)nox_xxx_modifGetDescById_413330(*(unsigned char*)(v7 + a4));
				*(uint32_t*)v8 = result;
				++v7;
				v8 += 4;
			} while (v7 < 4);
		}
	}
	return result;
}

//----- (00417B80) --------------------------------------------------------
char* sub_417B80(char a1, int a2, int a3) {
	char* result; // eax
	int v4;       // ecx
	int v5;       // edx
	int v6;       // ecx
	char* i;      // edx
	int v8;       // edx
	int v9;       // ecx
	char* j;      // edx

	result = nox_common_playerInfoGetByID_417040(a2);
	if (result) {
		v4 = ~a3;
		if (a1 == 84) {
			v5 = v4 & *((uint32_t*)result + 1);
			v6 = 0;
			*((uint32_t*)result + 1) = v5;
			for (i = result + 2324; *(uint32_t*)i != a3; i += 24) {
				if (++v6 >= 27) {
					return result;
				}
			}
			*(uint32_t*)&result[24 * v6 + 2324] = 0;
		} else {
			v8 = v4 & *(uint32_t*)result;
			v9 = 0;
			*(uint32_t*)result = v8;
			for (j = result + 2972; *(uint32_t*)j != a3; j += 24) {
				if (++v9 >= 26) {
					return result;
				}
			}
			*(uint32_t*)&result[24 * v9 + 2972] = 0;
		}
	}
	return result;
}

// 417EA6: variable 'v0' is possibly undefined

// 417F37: variable 'v0' is possibly undefined
