#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__gui__guigen.h"
#include "common__binfile.h"
#include "common__crypt.h"
#include "common__net_list.h"
#include "common__random.h"
#include "operators.h"
#include "input.h"

#include "client__video__draw_common.h"
#include "common__magic__speltree.h"
#include "server__script__script.h"

#include <time.h>

extern uint32_t dword_5d4594_2487656;
extern uint32_t dword_5d4594_3835368;
extern uint32_t dword_5d4594_2487804;
extern uint32_t dword_5d4594_2487632;
extern uint32_t dword_5d4594_3835364;
extern uint32_t dword_5d4594_2487884;
extern uint32_t dword_5d4594_3835372;
extern uint32_t dword_5d4594_3835392;
extern uint32_t dword_5d4594_2487580;
extern uint32_t dword_5d4594_2487568;
extern uint32_t dword_5d4594_2487536;
extern uint32_t dword_5d4594_2487628;
extern uint32_t dword_5d4594_2487584;
extern uint32_t dword_5d4594_2487564;
extern uint32_t dword_5d4594_2487672;
extern uint32_t dword_5d4594_2487676;
extern uint32_t dword_5d4594_3835388;
extern uint32_t dword_5d4594_2487652;
extern uint32_t dword_5d4594_3835348;
extern uint32_t dword_5d4594_2487576;
extern uint32_t dword_5d4594_2487624;
extern uint32_t dword_5d4594_3835352;
extern uint32_t dword_5d4594_2487620;
extern uint32_t dword_5d4594_2487708;
extern uint32_t nox_xxx_energyBoltTarget_5d4594_2487880;
extern uint32_t dword_5d4594_2487532;
extern uint32_t dword_5d4594_2487248;

extern uint32_t dword_5d4594_2487560;
extern uint32_t dword_5d4594_2487540;
extern uint32_t dword_5d4594_2487712;
extern uint32_t dword_5d4594_2487524;
extern uint32_t dword_5d4594_2487556;
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;

float get_nox_xxx_warriorMaxHealth_587000_312784();
float get_nox_xxx_wizardMaxHealth_587000_312816();
float get_nox_xxx_conjurerMaxHealth_587000_312800();

float get_nox_xxx_warriorMaxMana_587000_312788();
float get_nox_xxx_wizardMaximumMana_587000_312820();
float get_nox_xxx_conjurerMaxMana_587000_312804();

extern nox_tileDef_t nox_tile_defs_arr[176];

//----- (0051DEA0) --------------------------------------------------------
int nox_xxx_mapCountWallsMB_51DEA0(int a1) {
	int result; // eax

	if ((int)*(unsigned char*)(a1 + 5) < *getMemIntPtr(0x5D4594, 2487252)) {
		*getMemU32Ptr(0x5D4594, 2487252) = *(unsigned char*)(a1 + 5);
	}
	result = *(unsigned char*)(a1 + 6);
	if (result < *getMemIntPtr(0x5D4594, 2487256)) {
		*getMemU32Ptr(0x5D4594, 2487256) = *(unsigned char*)(a1 + 6);
	}
	return result;
}

//----- (0051DED0) --------------------------------------------------------
int sub_51DED0() {
	int* v0;            // edi
	char* v1;           // eax
	float* v3;          // esi
	int v4;             // ebx
	int v5;             // eax

	v0 = (int*)sub_45A060();
	if (!v0) {
		return 1;
	}
	do {
		if (!sub_4E3AD0(v0[27]) && sub_4E3B80(v0[27])) {
			v1 = (char*)nox_get_thing_name(v0[27]);
			v3 = (float*)nox_xxx_newObjectByTypeID_4E3810(v1);
			v4 = *((uint32_t*)v3 + 9);
			v3[14] = (double)v0[3] + 0.5;
			v3[15] = (double)v0[4] + 0.5;
			v5 = v0[32];
			*((uint32_t*)v3 + 10) = v5;
			*((uint32_t*)v3 + 11) = v5;
			*((uint32_t*)v3 + 9) = v5;
			*((uint32_t*)v3 + 4) = v0[30];
			*((uint32_t*)v3 + 5) = v0[70];
			nox_xxx_xfer_saveObj_51DF90((int)v3);
			*((uint32_t*)v3 + 9) = v4;
			nox_xxx_objectFreeMem_4E38A0((int)v3);
		}
		v0 = (int*)nox_drawable_next_45A070((int)v0);
	} while (v0);
	return 1;
}

//----- (0051E010) --------------------------------------------------------
int nox_xxx_nxzCompressFile_57BDD0(char* a1, char* a2);
int nox_xxx_mapSaveMap_51E010(char* a1, int a2) {
	char* v2;         // edi
	unsigned char v3; // cl
	int result;       // eax
	int v5;           // esi
	int v7;           // [esp+10h] [ebp-804h]
	char v8[1024];    // [esp+14h] [ebp-800h]
	char Mem[1024];   // [esp+414h] [ebp-400h]

	v7 = -86050098;
	strcpy(Mem, a1);
	v8[0] = 0;
	strncat(v8, a1, 1024-1);
	v8[strlen(v8)-4] = 0;
	v2 = &v8[strlen(v8) + 1];
	v3 = getMemByte(0x587000, 253116);
	*(uint32_t*)--v2 = *getMemU32Ptr(0x587000, 253112);
	v2[4] = v3;
	result = nox_xxx_cryptOpen_426910(Mem, 0, 19);
	if (result) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v7, 4u);
		v5 = nox_xxx_cryptFlush_4268E0();
		if (nox_xxx_map_51E140()) {
			sub_4268F0(v5);
			nox_xxx_cryptClose_4269F0();
			if (!a2 || (result = nox_xxx_nxzCompressFile_57BDD0(Mem, (int)v8)) != 0) {
				result = 1;
			}
		} else {
			nox_xxx_cryptClose_4269F0();
			result = 0;
		}
	}
	return result;
}
// 51E0D5: variable 'v6' is possibly undefined

//----- (0051E140) --------------------------------------------------------
void nox_xxx_map_5004F0();
void nox_xxx_mapSetWallInGlobalDir0pr1_5004D0();
int nox_xxx_map_51E140() {
	int result; // eax
	char v2;    // [esp+1h] [ebp-1h]

	*getMemU32Ptr(0x5D4594, 2487252) = 256;
	*getMemU32Ptr(0x5D4594, 2487256) = 256;
	nox_xxx_wallForeachFn_410640(nox_xxx_mapCountWallsMB_51DEA0, 0);
	nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 2487252), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 2487256), 4u);
	nox_xxx_mapWall_426A80(getMemIntPtr(0x5D4594, 2487252));
	nox_xxx_mapSetWallInGlobalDir0pr1_5004D0();
	if (nox_xxx_mapWriteSectionsMB_426E20(0)) {
		nox_xxx_map_5004F0();
		v2 = 0;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v2, 1u);
		result = 1;
	} else {
		nox_xxx_cryptClose_4269F0();
		result = 0;
	}
	return result;
}

//----- (0051E1D0) --------------------------------------------------------
int nox_xxx_mapGenSpellIdByName_51E1D0(const char* a1) {
	unsigned int i; // esi
	char v3[60];    // [esp+8h] [ebp-78h]
	char v4[66];        // [esp+44h] [ebp-3Ch]

	if (strlen(a1) >= sizeof(v3)) {
		return 0;
	}
	strcpy(v3, a1);
	for (i = 0; i < strlen(v3); ++i) {
		v3[i] = toupper(v3[i]);
	}
	nox_sprintf(v4, "SPELL_%s", v3);
	return nox_xxx_spellNameToN_4243F0(v4);
}
// 51E1D0: using guessed type char var_78[60];

//----- (0051E260) --------------------------------------------------------
int nox_xxx_mapGenReadTheme_51E260(int* a1, int a2) {
	FILE* v2;     // eax
	FILE* v3;     // edi
	int v4;       // eax
	int v5;       // ebx
	char v7[256]; // [esp+18h] [ebp-100h]

	memset(a1, 0, 0x45Cu);
	nox_sprintf(v7, "mapgen/%s.thm", a2);
	a1[3] = 3;
	a1[5] = 3;
	a1[1] = 5;
	a1[2] = 2;
	a1[4] = 1;
	a1[6] = 40;
	a1[7] = 20;
	a1[8] = 10;
	a1[9] = 5;
	a1[10] = 20;
	a1[16] = 1157234688;
	a1[18] = 100;
	a1[19] = time(0);
	a1[14] = 1;
	a1[22] = 0;
	a1[11] = 0;
	a1[12] = 100;
	a1[13] = 100;
	a1[30] = 0;
	a1[38] = 0;
	a1[46] = 0;
	*getMemU32Ptr(0x5D4594, 2487520) = 0;
	v2 = nox_binfile_open_408CC0(v7, 0);
	v3 = v2;
	if (!v2 || !nox_binfile_cryptSet_408D40((int)v2, 1)) {
		return 0;
	}
	while (nox_xxx_mapGenReadLine_51E540(v3, getMemAt(0x5D4594, 2487264))) {
		if (!nox_strcmpi("ALGORITHM_DATA", (const char*)getMemAt(0x5D4594, 2487264))) {
			v4 = nox_xxx_genReadAlgData_51EBB0((int)a1, v3);
			goto LABEL_23;
		}
		if (!nox_strcmpi("EXIT", (const char*)getMemAt(0x5D4594, 2487264))) {
			v4 = nox_xxx_genReadExit_51F800((int)a1, v3);
			goto LABEL_23;
		}
		if (!nox_strcmpi("PREFAB", (const char*)getMemAt(0x5D4594, 2487264))) {
			v4 = nox_xxx_genReadPrefab_520BF0((int)a1, v3);
			goto LABEL_23;
		}
		if (!nox_strcmpi("AMBIENT_LIGHT", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (!nox_xxx_mapGenReadLine_51E540(v3, getMemAt(0x5D4594, 2487264)) ||
				(a1[134] = atoi((const char*)getMemAt(0x5D4594, 2487264)),
				 !nox_xxx_mapGenReadLine_51E540(v3, getMemAt(0x5D4594, 2487264))) ||
				(a1[135] = atoi((const char*)getMemAt(0x5D4594, 2487264)),
				 !nox_xxx_mapGenReadLine_51E540(v3, getMemAt(0x5D4594, 2487264)))) {
				v5 = 0;
				goto LABEL_27;
			}
			a1[136] = atoi((const char*)getMemAt(0x5D4594, 2487264));
			continue;
		}
		if (nox_strcmpi("SPELL_SET", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (nox_strcmpi("ARMOR_SET", (const char*)getMemAt(0x5D4594, 2487264))) {
				if (nox_strcmpi("WEAPON_SET", (const char*)getMemAt(0x5D4594, 2487264))) {
					if (nox_strcmpi("DECOR", (const char*)getMemAt(0x5D4594, 2487264))) {
						v5 = 0;
						goto LABEL_27;
					}
					v4 = nox_xxx_genReadDecor_51F9F0(a1, v3);
				} else {
					v4 = nox_xxx_genReadWeaponSet_51F030((int)a1, v3);
				}
			} else {
				v4 = nox_xxx_genReadArmorSet_51F640((int)a1, v3);
			}
		} else {
			v4 = nox_xxx_genReadSpellSet_51EFB0((int)a1, v3);
		}
	LABEL_23:
		v5 = v4;
		if (!v4) {
			goto LABEL_27;
		}
	}
	v5 = 1;
LABEL_27:
	nox_binfile_close_408D90(v3);
	if (v5 != 1 || nox_xxx_mapgenCheckSettings_520AD0(a1 + 22) && nox_xxx_mapgenCheckSettings_520AD0(a1 + 30) &&
					   (!a1[46] || nox_xxx_mapgenCheckSettings_520AD0(a1 + 46))) {
		return v5;
	}
	return 0;
}

//----- (0051E540) --------------------------------------------------------
int nox_xxx_mapGenReadLine_51E540(FILE* a1, uint8_t* a2) {
	int result; // eax

	result = sub_51E570(a1, a2);
	if (result) {
		result = sub_51E670(a1) != 0;
	}
	return result;
}

//----- (0051E570) --------------------------------------------------------
int sub_51E570(FILE* a1, uint8_t* a2) {
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
			if (*(uint32_t*)CharType == 10) {
				++*getMemU32Ptr(0x5D4594, 2487520);
			}
			v6 = isspace(CharType[0]);
			if (v6) {
				break;
			}
			v4 = 0;
			if (v3 != 47 || v5 != 47) {
				*v2++ = v3;
			} else {
				sub_51E630(a1);
				v2 = a2;
				v3 = *(uint32_t*)CharType;
				v4 = 1;
			}
		}
	} while (v4);
	*v2 = 0;
	return 1;
}

//----- (0051E630) --------------------------------------------------------
int sub_51E630(FILE* a1) {
	FILE* v1;   // esi
	int result; // eax

	v1 = a1;
	while (1) {
		LOBYTE(a1) = 0;
		nox_binfile_fread_408E40((char*)&a1, 1, 1, v1);
		result = nox_binfile_lastErr_409370(v1);
		if (result == -1) {
			break;
		}
		if ((uint8_t)a1 == 10) {
			++*getMemU32Ptr(0x5D4594, 2487520);
			return result;
		}
	}
	return result;
}

//----- (0051E670) --------------------------------------------------------
int sub_51E670(FILE* a1) {
	int result; // eax
	int v2;     // [esp+4h] [ebp-4h]

	if (nox_strcmpi("IF", (const char*)getMemAt(0x5D4594, 2487264))) {
		if (nox_strcmpi("ELSE", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (nox_strcmpi("ENDIF", (const char*)getMemAt(0x5D4594, 2487264))) {
				return 1;
			}
		} else {
			result = sub_51E720(a1);
			if (!result) {
				return result;
			}
		}
		return nox_xxx_mapGenReadLine_51E540(a1, getMemAt(0x5D4594, 2487264)) != 0;
	}
	result = sub_51E800((int)a1, &v2);
	if (result) {
		if (v2) {
			return nox_xxx_mapGenReadLine_51E540(a1, getMemAt(0x5D4594, 2487264)) != 0;
		}
		result = sub_51E780(a1);
		if (result) {
			return nox_xxx_mapGenReadLine_51E540(a1, getMemAt(0x5D4594, 2487264)) != 0;
		}
	}
	return result;
}

//----- (0051E720) --------------------------------------------------------
int sub_51E720(FILE* a1) {
	do {
		while (1) {
			if (!sub_51E570(a1, getMemAt(0x5D4594, 2487264))) {
				return 0;
			}
			if (nox_strcmpi("IF", (const char*)getMemAt(0x5D4594, 2487264))) {
				break;
			}
			if (!sub_51E720(a1)) {
				return 0;
			}
		}
	} while (nox_strcmpi("ENDIF", (const char*)getMemAt(0x5D4594, 2487264)));
	return 1;
}

//----- (0051E780) --------------------------------------------------------
int sub_51E780(FILE* a1) {
	do {
		while (1) {
			if (!sub_51E570(a1, getMemAt(0x5D4594, 2487264))) {
				return 0;
			}
			if (nox_strcmpi("IF", (const char*)getMemAt(0x5D4594, 2487264))) {
				break;
			}
			if (!sub_51E720(a1)) {
				return 0;
			}
		}
	} while (nox_strcmpi("ENDIF", (const char*)getMemAt(0x5D4594, 2487264)) &&
			 nox_strcmpi("ELSE", (const char*)getMemAt(0x5D4594, 2487264)));
	return 1;
}

//----- (0051E800) --------------------------------------------------------
int sub_51E800(int a1, uint32_t* a2) {
	FILE* v2;     // edi
	uint32_t* v3; // esi
	char* v4;     // eax
	int result;   // eax
	char* v6;     // eax
	char* v7;     // eax
	int v8;       // eax
	int v9;       // esi
	int v10;      // eax
	int v11;      // esi

	v2 = (FILE*)a1;
	if (!nox_xxx_mapGenReadLine_51E540((FILE*)a1, getMemAt(0x5D4594, 2487264))) {
		return 0;
	}
	if (!nox_strcmpi("WIZARD", (const char*)getMemAt(0x5D4594, 2487264))) {
		v3 = a2;
		*a2 = 0;
		v4 = nox_common_playerInfoGetFirst_416EA0();
		if (v4) {
			while (v4[2251] != 1) {
				v4 = nox_common_playerInfoGetNext_416EE0((int)v4);
				if (!v4) {
					return 1;
				}
			}
			*v3 = 1;
			return 1;
		}
		return 1;
	}
	if (!nox_strcmpi("CONJURER", (const char*)getMemAt(0x5D4594, 2487264))) {
		v3 = a2;
		*a2 = 0;
		v6 = nox_common_playerInfoGetFirst_416EA0();
		if (v6) {
			while (v6[2251] != 2) {
				v6 = nox_common_playerInfoGetNext_416EE0((int)v6);
				if (!v6) {
					return 1;
				}
			}
			*v3 = 1;
			return 1;
		}
		return 1;
	}
	if (!nox_strcmpi("WARRIOR", (const char*)getMemAt(0x5D4594, 2487264))) {
		v3 = a2;
		*a2 = 0;
		v7 = nox_common_playerInfoGetFirst_416EA0();
		if (v7) {
			while (v7[2251]) {
				v7 = nox_common_playerInfoGetNext_416EE0((int)v7);
				if (!v7) {
					return 1;
				}
			}
			*v3 = 1;
			return 1;
		}
		return 1;
	}
	if (!nox_strcmpi("EXPERIENCE_LEVEL", (const char*)getMemAt(0x5D4594, 2487264))) {
		if (sub_51EAF0((int)v2, &a1) && nox_xxx_mapGenReadLine_51E540(v2, getMemAt(0x5D4594, 2487264))) {
			v8 = atoi((const char*)getMemAt(0x5D4594, 2487264));
			switch (a1) {
			case 0:
				*a2 = v8 > 3;
				result = 1;
				break;
			case 1:
				*a2 = v8 < 3;
				result = 1;
				break;
			case 2:
				*a2 = v8 == 3;
				result = 1;
				break;
			case 3:
				*a2 = v8 != 3;
				result = 1;
				break;
			default:
				return 1;
			}
			return result;
		}
		return 0;
	}
	if (nox_strcmpi("NUMPLAYERS", (const char*)getMemAt(0x5D4594, 2487264))) {
		if (strchr((const char*)getMemAt(0x5D4594, 2487264), 37)) {
			v11 = atoi((const char*)getMemAt(0x5D4594, 2487264));
			*a2 = nox_xxx_mapGenRandFunc_526AC0(1, 100) <= v11;
			return 1;
		}
		return 0;
	}
	v9 = nox_common_playerInfoCount_416F40();
	if (!sub_51EAF0((int)v2, &a1) || !nox_xxx_mapGenReadLine_51E540(v2, getMemAt(0x5D4594, 2487264))) {
		return 0;
	}
	v10 = atoi((const char*)getMemAt(0x5D4594, 2487264));
	switch (a1) {
	case 0:
		*a2 = v9 < v10;
		result = 1;
		break;
	case 1:
		*a2 = v9 > v10;
		result = 1;
		break;
	case 2:
		*a2 = v9 == v10;
		result = 1;
		break;
	case 3:
		*a2 = v9 != v10;
		result = 1;
		break;
	default:
		return 1;
	}
	return result;
}

//----- (0051EAF0) --------------------------------------------------------
int sub_51EAF0(int a1, uint32_t* a2) {
	int result; // eax

	if (!nox_xxx_mapGenReadLine_51E540((FILE*)a1, getMemAt(0x5D4594, 2487264))) {
		return 0;
	}
	if (!nox_strcmpi("<", (const char*)getMemAt(0x5D4594, 2487264))) {
		*a2 = 0;
		return 1;
	}
	if (!nox_strcmpi(">", (const char*)getMemAt(0x5D4594, 2487264))) {
		result = 1;
		*a2 = 1;
		return result;
	}
	if (!nox_strcmpi("==", (const char*)getMemAt(0x5D4594, 2487264))) {
		*a2 = 2;
		return 1;
	}
	if (nox_strcmpi("!=", (const char*)getMemAt(0x5D4594, 2487264))) {
		return 0;
	}
	result = 1;
	*a2 = 3;
	return result;
}

//----- (0051EBB0) --------------------------------------------------------
int nox_xxx_genReadAlgData_51EBB0(int a1, FILE* a2) {
	int v2;            // edi
	const char** v3;   // eax
	unsigned char* v4; // esi
	int v5;            // ecx
	char v7;           // [esp+10h] [ebp-3Ch]

	while (nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
		if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
			return 1;
		}
		if (!nox_xxx_mapGenReadLine_51E540(a2, &v7)) {
			break;
		}
		if (nox_strcmpi("midHallLength", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (nox_strcmpi("hallLengthVariance", (const char*)getMemAt(0x5D4594, 2487264))) {
				if (nox_strcmpi("midHallWidth", (const char*)getMemAt(0x5D4594, 2487264))) {
					if (nox_strcmpi("hallWidthVariance", (const char*)getMemAt(0x5D4594, 2487264))) {
						if (nox_strcmpi("hallLimit", (const char*)getMemAt(0x5D4594, 2487264))) {
							if (nox_strcmpi("hallBranchRate", (const char*)getMemAt(0x5D4594, 2487264))) {
								if (nox_strcmpi("hallRoomRate", (const char*)getMemAt(0x5D4594, 2487264))) {
									if (nox_strcmpi("midRoomSize", (const char*)getMemAt(0x5D4594, 2487264))) {
										if (nox_strcmpi("roomVariance", (const char*)getMemAt(0x5D4594, 2487264))) {
											if (nox_strcmpi("irregularRoomRate",
														 (const char*)getMemAt(0x5D4594, 2487264))) {
												if (nox_strcmpi("mapSize", (const char*)getMemAt(0x5D4594, 2487264))) {
													if (nox_strcmpi("recurs"
																 "ionLim"
																 "it",
																 (const char*)getMemAt(0x5D4594, 2487264))) {
														if (nox_strcmpi("seed",
																	 (const char*)getMemAt(0x5D4594, 2487264))) {
															if (nox_strcmpi("emptyRoomRate",
																		 (const char*)getMemAt(0x5D4594, 2487264))) {
																if (nox_strcmpi("mergeRate", (const char*)getMemAt(
																							  0x5D4594, 2487264))) {
																	if (nox_strcmpi("useDoors", (const char*)getMemAt(
																								 0x5D4594, 2487264))) {
																		if (nox_strcmpi("debug", (const char*)getMemAt(
																								  0x5D4594, 2487264))) {
																			if (nox_strcmpi("adjacentPortalRate",
																						 (const char*)getMemAt(
																							 0x5D4594, 2487264))) {
																				if (!nox_strcmpi("skeleton",
																							  (const char*)getMemAt(
																								  0x5D4594, 2487264))) {
																					v2 = 0;
																					if (*getMemU32Ptr(0x587000,
																									  253244)) {
																						v3 = (const char**)getMemAt(
																							0x587000, 253244);
																						v4 = getMemAt(0x587000, 253244);
																						do {
																							if (!nox_strcmpi(*v3, &v7)) {
																								break;
																							}
																							v5 = *((uint32_t*)v4 + 1);
																							v4 += 4;
																							++v2;
																							v3 = (const char**)v4;
																						} while (v5);
																					}
																					*(uint32_t*)a1 =
																						*getMemU32Ptr(0x587000,
																									  253244 +
																										  4 * v2) != 0
																							? v2
																							: 0;
																				}
																			} else {
																				*(uint32_t*)(a1 + 52) = atoi(&v7);
																			}
																		} else if (nox_strcmpi(&v7, "true")) {
																			*(uint32_t*)(a1 + 60) = 0;
																		} else {
																			*(uint32_t*)(a1 + 60) = 1;
																		}
																	} else if (nox_strcmpi(&v7, "true")) {
																		*(uint32_t*)(a1 + 56) = 0;
																	} else {
																		*(uint32_t*)(a1 + 56) = 1;
																	}
																} else {
																	*(uint32_t*)(a1 + 48) = atoi(&v7);
																}
															} else {
																*(uint32_t*)(a1 + 44) = atoi(&v7);
															}
														} else {
															*(uint32_t*)(a1 + 76) = atoi(&v7);
														}
													} else {
														*(uint32_t*)(a1 + 72) = atoi(&v7);
													}
												} else {
													*(float*)(a1 + 64) = atof(&v7);
												}
											} else {
												*(uint32_t*)(a1 + 40) = atoi(&v7);
											}
										} else {
											*(uint32_t*)(a1 + 36) = atoi(&v7);
										}
									} else {
										*(uint32_t*)(a1 + 32) = atoi(&v7);
									}
								} else {
									*(uint32_t*)(a1 + 28) = atoi(&v7);
								}
							} else {
								*(uint32_t*)(a1 + 24) = atoi(&v7);
							}
						} else {
							*(uint32_t*)(a1 + 20) = atoi(&v7);
						}
					} else {
						*(uint32_t*)(a1 + 16) = atoi(&v7);
					}
				} else {
					*(uint32_t*)(a1 + 12) = atoi(&v7);
				}
			} else {
				*(uint32_t*)(a1 + 8) = atoi(&v7);
			}
		} else {
			*(uint32_t*)(a1 + 4) = atoi(&v7);
		}
	}
	return 0;
}

//----- (0051EFB0) --------------------------------------------------------
int nox_xxx_genReadSpellSet_51EFB0(int a1, FILE* a2) {
	int v2; // eax

	while (1) {
		if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
			return 0;
		}
		if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264)) || *(uint32_t*)(a1 + 1096) >= 137) {
			break;
		}
		v2 = nox_xxx_mapGenSpellIdByName_51E1D0((const char*)getMemAt(0x5D4594, 2487264));
		if (v2) {
			*(uint32_t*)(a1 + 4 * (*(uint32_t*)(a1 + 1096))++ + 548) = v2;
		}
	}
	return 1;
}

//----- (0051F030) --------------------------------------------------------
int nox_xxx_genReadWeaponSet_51F030(int a1, FILE* a2) {
	char* v2; // ebx
	int v3;   // ecx
	int v4;   // eax
	int i;    // eax

	dword_5d4594_2487524 = 0;
	while (1) {
		if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
			return 0;
		}
		if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
			break;
		}
		if (!nox_strcmpi("WEAPON", (const char*)getMemAt(0x5D4594, 2487264))) {
			v2 = (char*)calloc(1u, 0x9Cu);
			if (!v2 || !nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				return 0;
			}
			strcpy(v2 + 60, (const char*)getMemAt(0x5D4594, 2487264));
			if (nox_strcmpi("TEMPLATE", v2 + 60)) {
				if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
					return 0;
				}
				strcpy(v2, (const char*)getMemAt(0x5D4594, 2487264));
			}
			if (!sub_51F230((int)v2, a2)) {
				return 0;
			}
			if (nox_strcmpi("TEMPLATE", v2 + 60)) {
				*((uint32_t*)v2 + 38) = 0;
				v3 = *(uint32_t*)(a1 + 1100);
				if (v3) {
					for (i = *(uint32_t*)(v3 + 152); i; i = *(uint32_t*)(i + 152)) {
						v3 = i;
					}
					*(uint32_t*)(v3 + 152) = v2;
					++*(uint32_t*)(a1 + 1104);
				} else {
					v4 = *(uint32_t*)(a1 + 1104);
					*(uint32_t*)(a1 + 1100) = v2;
					*(uint32_t*)(a1 + 1104) = v4 + 1;
				}
			} else {
				dword_5d4594_2487524 = v2;
			}
		}
	}
	if (dword_5d4594_2487524) {
		nox_xxx_mapGenFreeStr_51F1F0(*(void**)&dword_5d4594_2487524);
	}
	dword_5d4594_2487524 = 0;
	return 1;
}

//----- (0051F1F0) --------------------------------------------------------
void nox_xxx_mapGenFreeStr_51F1F0(void* lpMem) {
	char* v1; // esi
	int v2;   // edi

	v1 = (char*)lpMem + 120;
	v2 = 4;
	do {
		if (*((uint32_t*)v1 + 4)) {
			free(*(void**)v1);
		}
		v1 += 4;
		--v2;
	} while (v2);
	free(lpMem);
}

//----- (0051F230) --------------------------------------------------------
int sub_51F230(int a1, FILE* a2) {
	int v2;           // ebx
	int v3;           // edx
	uint32_t* v4;     // ecx
	int v5;           // edi
	signed int* v6;   // edx
	int v7;           // esi
	int v8;           // edi
	signed int v9;    // ebx
	int v10;          // ebp
	uint32_t* v11;    // esi
	int v12;          // edx
	const char* v13;  // edi
	uint32_t* v14;    // ecx
	char* v15;        // ebp
	signed int v16;   // esi
	char* v17;        // edi
	int v18;          // esi
	int v19;          // ebp
	const char* v20;  // edx
	char* v21;        // eax
	signed int v22;   // esi
	char* v23;        // edi
	bool v24;         // zf
	char* v25;        // edx
	char* v26;        // ebp
	const char* v27;  // edi
	unsigned int v28; // ecx
	char v29;         // al
	const char* v30;  // esi
	char* v31;        // edi
	char* v32;        // edi
	const char* v33;  // esi
	int v34;          // ecx
	uint32_t* v35;    // eax
	void* v36;        // eax
	int v37;          // edx
	signed int v38;   // ebp
	const char* v39;  // edi
	char* v40;        // esi
	uint32_t* v41;    // ecx
	char* v42;        // eax
	int v43;          // edx
	char* v44;        // edi
	uint32_t* v46;    // [esp+10h] [ebp-12C38h]
	char* v47;        // [esp+14h] [ebp-12C34h]
	char* v48;        // [esp+18h] [ebp-12C30h]
	uint32_t* v49;    // [esp+1Ch] [ebp-12C2Ch]
	char* v50;        // [esp+20h] [ebp-12C28h]
	char* v51;        // [esp+24h] [ebp-12C24h]
	int* v52;         // [esp+28h] [ebp-12C20h]
	char* v53;        // [esp+2Ch] [ebp-12C1Ch]
	int v54;          // [esp+30h] [ebp-12C18h]
	int v55;          // [esp+34h] [ebp-12C14h]
	int v56;          // [esp+38h] [ebp-12C10h]
	int v57;          // [esp+3Ch] [ebp-12C0Ch]
	int v58;          // [esp+40h] [ebp-12C08h]
	int i;            // [esp+44h] [ebp-12C04h]
	char v60[15360];  // [esp+48h] [ebp-12C00h]
	char v61[61440];  // [esp+3C48h] [ebp-F000h]

	v55 = 0;
	v56 = 0;
	v57 = 0;
	v58 = 0;
	while (1) {
		if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
			return 0;
		}
		if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
			v50 = v61;
			v4 = (uint32_t*)(a1 + 120);
			v5 = 16 - a1;
			v6 = &v55;
			v7 = 120 - (uint32_t)&v55;
			v47 = 0;
			v52 = &v55;
			v49 = (uint32_t*)(a1 + 120);
			v54 = 16 - a1;
			for (i = 120 - (uint32_t)&v55;; v7 = i) {
				if (dword_5d4594_2487524) {
					v8 = (int)v4 + v5;
					v9 = *(uint32_t*)(v8 + dword_5d4594_2487524);
					if (v9 > 0) {
						v10 = *(uint32_t*)(v8 + dword_5d4594_2487524);
						v11 = (signed int*)((char*)v6 + v7 + dword_5d4594_2487524);
						v12 = 0;
						v46 = v11;
						do {
							v13 = (const char*)(v12 + *v46);
							v48 = &v60[v12];
							v12 += 60;
							strcpy(v48, v13);
							--v10;
						} while (v10);
					}
					v14 = (uint32_t*)*v52;
					if (*v52 > 0) {
						v15 = v50;
						v48 = v50;
						v46 = v14;
						v51 = &v60[60 * v9];
						do {
							if (*v15 == 45) {
								v16 = 0;
								if (v9 <= 0) {
									goto LABEL_43;
								}
								v17 = v60;
								while (nox_strcmpi(v17, v15 + 1)) {
									++v16;
									v17 += 60;
									if (v16 >= v9) {
										goto LABEL_43;
									}
								}
								if (v16 >= v9) {
									goto LABEL_43;
								}
								v18 = v16 + 1;
								if (v18 < v9) {
									v19 = v9 - v18;
									v20 = &v60[60 * v18];
									do {
										v53 = (char*)(v20 - 60);
										v20 += 60;
										strcpy(v53, v20);
										--v19;
									} while (v19);
									v15 = v48;
								}
								--v9;
								v21 = v51 - 60;
							} else {
								v22 = 0;
								if (v9 > 0) {
									v23 = v60;
									do {
										if (!nox_strcmpi(v23, v15)) {
											break;
										}
										++v22;
										v23 += 60;
									} while (v22 < v9);
								}
								if (v22 != v9) {
									goto LABEL_43;
								}
								++v9;
								v21 = strcpy(v51, v15) + 60;
							}
							v51 = v21;
						LABEL_43:
							v15 += 60;
							v24 = v46 == (uint32_t*)1;
							v48 = v15;
							v46 = (uint32_t*)((char*)v46 - 1);
						} while (!v24);
					}
				} else {
					v9 = *v6;
					v24 = *v6 == 0;
					if (*v6 <= 0) {
						goto LABEL_49;
					}
					v25 = v50;
					v26 = v60;
					v46 = (uint32_t*)v9;
					do {
						v27 = v25;
						v25 += 60;
						v28 = strlen(v27) + 1;
						v29 = v28;
						v30 = v27;
						v31 = v26;
						v26 += 60;
						v28 >>= 2;
						memcpy(v31, v30, 4 * v28);
						v33 = &v30[4 * v28];
						v32 = &v31[4 * v28];
						v34 = v29 & 3;
						v35 = (uint32_t*)((char*)v46 - 1);
						v24 = v46 == (uint32_t*)1;
						memcpy(v32, v33, v34);
						v46 = v35;
					} while (!v24);
				}
				v24 = v9 == 0;
			LABEL_49:
				if (!v24) {
					v36 = calloc(v9, 0x3Cu);
					*v49 = v36;
					if (!v36) {
						return 0;
					}
					if (v9 > 0) {
						v37 = 0;
						v38 = v9;
						do {
							v39 = &v60[v37];
							v40 = (char*)(v37 + *v49);
							v37 += 60;
							v53 = strcpy(v40, v39);
							--v38;
						} while (v38);
					}
				}
				v41 = v49;
				v42 = v47;
				v43 = (int)v52;
				v44 = v50;
				v49[4] = v9;
				v6 = (signed int*)(v43 + 4);
				v4 = v41 + 1;
				v47 = v42 + 1;
				v52 = v6;
				v49 = v4;
				v50 = v44 + 15360;
				if ((int)(v42 + 1) >= 4) {
					return 1;
				}
				v5 = v54;
			}
		}
		if (nox_strcmpi("EFFECTIVENESS", (const char*)getMemAt(0x5D4594, 2487264)) &&
			nox_strcmpi("QUALITY", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (!nox_strcmpi("MATERIAL", (const char*)getMemAt(0x5D4594, 2487264))) {
				v2 = 1;
				goto LABEL_13;
			}
			if (!nox_strcmpi("PRIMARY_ENCHANTMENT", (const char*)getMemAt(0x5D4594, 2487264))) {
				v2 = 2;
				goto LABEL_13;
			}
			if (!nox_strcmpi("SECONDARY_ENCHANTMENT", (const char*)getMemAt(0x5D4594, 2487264))) {
				v2 = 3;
				goto LABEL_13;
			}
			return 0;
		}
		v2 = 0;
	LABEL_13:
		while (1) {
			if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				return 0;
			}
			if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
				break;
			}
			v3 = *(&v55 + v2);
			v47 = &v61[60 * (v3 + (v2 << 8))];
			strcpy(v47, (const char*)getMemAt(0x5D4594, 2487264));
			*(&v55 + v2) = v3 + 1;
		}
	}
}
// 51F230: using guessed type char var_F000[61440];
// 51F230: using guessed type char var_12C00[15360];

//----- (0051F640) --------------------------------------------------------
int nox_xxx_genReadArmorSet_51F640(int a1, FILE* a2) {
	char* v2; // ebx
	int v3;   // ecx
	int v4;   // eax
	int i;    // eax

	dword_5d4594_2487524 = 0;
	while (1) {
		if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
			return 0;
		}
		if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
			break;
		}
		if (!nox_strcmpi("ARMOR", (const char*)getMemAt(0x5D4594, 2487264))) {
			v2 = (char*)calloc(1u, 0x9Cu);
			if (!v2 || !nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				return 0;
			}
			strcpy(v2 + 60, (const char*)getMemAt(0x5D4594, 2487264));
			if (nox_strcmpi("TEMPLATE", v2 + 60)) {
				if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
					return 0;
				}
				strcpy(v2, (const char*)getMemAt(0x5D4594, 2487264));
			}
			if (!sub_51F230((int)v2, a2)) {
				return 0;
			}
			if (nox_strcmpi("TEMPLATE", v2 + 60)) {
				*((uint32_t*)v2 + 38) = 0;
				v3 = *(uint32_t*)(a1 + 1108);
				if (v3) {
					for (i = *(uint32_t*)(v3 + 152); i; i = *(uint32_t*)(i + 152)) {
						v3 = i;
					}
					*(uint32_t*)(v3 + 152) = v2;
					++*(uint32_t*)(a1 + 1112);
				} else {
					v4 = *(uint32_t*)(a1 + 1112);
					*(uint32_t*)(a1 + 1108) = v2;
					*(uint32_t*)(a1 + 1112) = v4 + 1;
				}
			} else {
				dword_5d4594_2487524 = v2;
			}
		}
	}
	if (dword_5d4594_2487524) {
		nox_xxx_mapGenFreeStr_51F1F0(*(void**)&dword_5d4594_2487524);
	}
	dword_5d4594_2487524 = 0;
	return 1;
}

//----- (0051F800) --------------------------------------------------------
int nox_xxx_genReadExit_51F800(int a1, FILE* a2) {
	while (1) {
		if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
			return 0;
		}
		if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
			break;
		}
		if (nox_strcmpi("OBJECT", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (nox_strcmpi("LINKDATA", (const char*)getMemAt(0x5D4594, 2487264)) ||
				!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				return 0;
			}
			strcpy((char*)(a1 + 476), (const char*)getMemAt(0x5D4594, 2487264));
		} else {
			if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				return 0;
			}
			strcpy((char*)(a1 + 216 + (*(uint32_t*)(a1 + 472) << 6)), (const char*)getMemAt(0x5D4594, 2487264));
			if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				return 0;
			}
			if (nox_strcmpi("NORTH", (const char*)getMemAt(0x5D4594, 2487264))) {
				if (nox_strcmpi("SOUTH", (const char*)getMemAt(0x5D4594, 2487264))) {
					if (nox_strcmpi("EAST", (const char*)getMemAt(0x5D4594, 2487264))) {
						if (!nox_strcmpi("WEST", (const char*)getMemAt(0x5D4594, 2487264))) {
							*(uint32_t*)((*(uint32_t*)(a1 + 472) << 6) + a1 + 216 + 60) = 3;
						}
						++*(uint32_t*)(a1 + 472);
					} else {
						*(uint32_t*)(((*(uint32_t*)(a1 + 472))++ << 6) + a1 + 216 + 60) = 2;
					}
				} else {
					*(uint32_t*)(((*(uint32_t*)(a1 + 472))++ << 6) + a1 + 216 + 60) = 1;
				}
			} else {
				*(uint32_t*)(((*(uint32_t*)(a1 + 472))++ << 6) + a1 + 216 + 60) = 0;
			}
		}
	}
	return 1;
}

//----- (0051F9F0) --------------------------------------------------------
int nox_xxx_genReadDecor_51F9F0(uint32_t* a1, FILE* a2) {
	uint8_t* v2;       // eax
	int v3;            // edi
	FILE* v4;          // esi
	int v5;            // ebx
	int v6;            // ebp
	const char** v7;   // eax
	unsigned char* v8; // esi
	int v9;            // ecx
	int v10;           // ecx
	int result;        // eax
	int v12;           // ecx
	int v13;           // ecx
	int v14;           // ecx
	char* v15;         // [esp+10h] [ebp-4h]

	v15 = 0;
	v2 = calloc(1u, 0xE0u);
	v3 = (int)v2;
	if (!v2) {
		return 0;
	}
	v4 = a2;
	v2[64] = 0;
	*((uint32_t*)v2 + 18) = 1000;
	*((uint32_t*)v2 + 19) = 0;
	*((uint32_t*)v2 + 20) = 999999;
	if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
		return 0;
	}
	if (!nox_xxx_mapGenReadLine_51E540(a2, (uint8_t*)v3)) {
		return 0;
	}
	if (!nox_strcmpi("ROOM", (const char*)getMemAt(0x5D4594, 2487264))) {
		*(uint32_t*)(v3 + 60) = 0;
	} else if (!nox_strcmpi("HALL", (const char*)getMemAt(0x5D4594, 2487264))) {
		*(uint32_t*)(v3 + 60) = 1;
	} else if (!nox_strcmpi("TEMPLATE", (const char*)getMemAt(0x5D4594, 2487264))) {
		*(uint32_t*)(v3 + 60) = 2;
	} else if (!nox_strcmpi("BACKDROP", (const char*)getMemAt(0x5D4594, 2487264))) {
		*(uint32_t*)(v3 + 60) = 3;
	} else {
		return 0;
	}
	while (1) {
		if (!nox_xxx_mapGenReadLine_51E540(v4, getMemAt(0x5D4594, 2487264))) {
			return 0;
		}
		if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
			switch (*(uint32_t*)(v3 + 60)) {
			case 0:
				*(uint32_t*)(v3 + 220) = a1[22];
				v10 = a1[23];
				a1[22] = v3;
				a1[23] = v10 + 1;
				result = 1;
				break;
			case 1:
				*(uint32_t*)(v3 + 220) = a1[30];
				v12 = a1[31];
				a1[30] = v3;
				a1[31] = v12 + 1;
				result = 1;
				break;
			case 2:
				*(uint32_t*)(v3 + 220) = a1[38];
				v13 = a1[39];
				a1[38] = v3;
				a1[39] = v13 + 1;
				result = 1;
				break;
			case 3:
				*(uint32_t*)(v3 + 220) = a1[46];
				v14 = a1[47] + 1;
				a1[46] = v3;
				a1[47] = v14;
				result = 1;
				break;
			default:
				result = 1;
				break;
			}
			return result;
		}
		if (nox_strcmpi("WALL_FLOOR", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (nox_strcmpi("OCCUR_CONSTRAINT", (const char*)getMemAt(0x5D4594, 2487264))) {
				if (nox_strcmpi("OCCUR_LIMIT", (const char*)getMemAt(0x5D4594, 2487264))) {
					if (nox_strcmpi("MUST_OCCUR", (const char*)getMemAt(0x5D4594, 2487264))) {
						if (nox_strcmpi("FREQUENCY", (const char*)getMemAt(0x5D4594, 2487264))) {
							if (nox_strcmpi("DOOR", (const char*)getMemAt(0x5D4594, 2487264))) {
								if (nox_strcmpi("DOUBLE_DOOR", (const char*)getMemAt(0x5D4594, 2487264))) {
									if (nox_strcmpi("ROOM_SIZE_CONSTRAINT", (const char*)getMemAt(0x5D4594, 2487264))) {
										if (nox_strcmpi("COPY", (const char*)getMemAt(0x5D4594, 2487264))) {
											if (nox_strcmpi("DECOR_SET", (const char*)getMemAt(0x5D4594, 2487264))) {
												v5 = 0;
												v6 = 0;
												if (!*getMemU32Ptr(0x587000, 253200)) {
													return 0;
												}
												v7 = (const char**)getMemAt(0x587000, 253200);
												v8 = getMemAt(0x587000, 253200);
												do {
													if (!nox_strcmpi(*v7, (const char*)getMemAt(
																		   0x5D4594, 2487264))) {
														if (!v15 || !sub_51FEC0((int)v15, v6, a2)) {
															return 0;
														}
														v5 = 1;
													}
													v9 = *((uint32_t*)v8 + 1);
													v8 += 4;
													++v6;
													v7 = (const char**)v8;
												} while (v9);
												if (!v5) {
													return 0;
												}
												v4 = a2;
											} else if (!nox_xxx_genDecorReadDecorSet_51FFA0(v3, v4)) {
												return 0;
											}
										} else if (!nox_xxx_genDecorReadCopy_520660(a1, (const char*)v3, v4)) {
											return 0;
										}
									} else if (!nox_xxx_genDecorReadRoomSizeCon_5209F0(v3, v4)) {
										return 0;
									}
								} else if (!nox_xxx_genDecorReadDoubleDoor_520AB0(v3, v4)) {
									return 0;
								}
							} else if (!nox_xxx_genDecorReadDoor_520A90(v3, v4)) {
								return 0;
							}
						} else if (!nox_xxx_genDecorReadFrequency_520910(v3, v4)) {
							return 0;
						}
					} else {
						*(uint8_t*)(v3 + 67) = 1;
					}
				} else if (!nox_xxx_genDecorReadOccurLimit_5208D0(v3, v4)) {
					return 0;
				}
			} else if (!nox_xxx_genDecorReadOccurConstraint_520810(v3, v4)) {
				return 0;
			}
		} else {
			v15 = nox_xxx_genDecorReadWallFloor_51FE00(v3, v4);
			if (!v15) {
				return 0;
			}
		}
	}
	return 0;
}

//----- (0051FE00) --------------------------------------------------------
char* nox_xxx_genDecorReadWallFloor_51FE00(int a1, FILE* a2) {
	char* result; // eax
	char* v3;     // ebx
	int v4;       // ecx

	result = (char*)calloc(1u, 0x80u);
	v3 = result;
	if (result) {
		result = (char*)nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264));
		if (result) {
			strcpy(v3, (const char*)getMemAt(0x5D4594, 2487264));
			result = (char*)nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264));
			if (result) {
				strcpy(v3 + 60, (const char*)getMemAt(0x5D4594, 2487264));
				*((uint32_t*)v3 + 31) = *(uint32_t*)(a1 + 84);
				v4 = *(uint32_t*)(a1 + 88) + 1;
				*(uint32_t*)(a1 + 84) = v3;
				*(uint32_t*)(a1 + 88) = v4;
				result = v3;
			}
		}
	}
	return result;
}

//----- (0051FEC0) --------------------------------------------------------
int sub_51FEC0(int a1, int a2, FILE* a3) {
	int result; // eax
	int v4;     // ebx
	int v5;     // ecx
	int i;      // eax

	result = (int)calloc(1u, 0x80u);
	v4 = result;
	if (result) {
		*(uint32_t*)result = a2;
		result = nox_xxx_mapGenReadLine_51E540(a3, getMemAt(0x5D4594, 2487264));
		if (result) {
			strcpy((char*)(v4 + 4), (const char*)getMemAt(0x5D4594, 2487264));
			result = nox_xxx_mapGenReadLine_51E540(a3, getMemAt(0x5D4594, 2487264));
			if (result) {
				strcpy((char*)(v4 + 64), (const char*)getMemAt(0x5D4594, 2487264));
				v5 = *(uint32_t*)(a1 + 120);
				if (v5) {
					for (i = *(uint32_t*)(v5 + 124); i; i = *(uint32_t*)(i + 124)) {
						v5 = i;
					}
					*(uint32_t*)(v5 + 124) = v4;
					result = 1;
				} else {
					*(uint32_t*)(a1 + 120) = v4;
					result = 1;
				}
			}
		}
	}
	return result;
}

//----- (0051FFA0) --------------------------------------------------------
int nox_xxx_genDecorReadDecorSet_51FFA0(int a1, FILE* a2) {
	int v2;       // ebp
	int result;   // eax
	int* v4;      // esi
	uint32_t* v5; // eax
	uint32_t* v6; // eax
	uint32_t* v7; // eax
	int v8;       // eax
	int v9;       // eax
	int v10;      // ecx
	int i;        // eax
	int v12;      // [esp+10h] [ebp-4h]

	v2 = 0;
	result = (int)calloc(1u, 0x18u);
	v4 = (int*)result;
	v12 = result;
	if (!result) {
		return 0;
	}
	result = nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264));
	if (!result) {
		return 0;
	}
	*v4 = atoi((const char*)getMemAt(0x5D4594, 2487264));
	result = nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264));
	if (!result) {
		return 0;
	}
	v4[1] = atoi((const char*)getMemAt(0x5D4594, 2487264));
	while (1) {
		while (1) {
			if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				return 0;
			}
			if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
				v10 = *(uint32_t*)(a1 + 92);
				if (v10) {
					for (i = *(uint32_t*)(v10 + 20); i; i = *(uint32_t*)(i + 20)) {
						v10 = i;
					}
					*(uint32_t*)(v10 + 20) = v4;
				} else {
					*(uint32_t*)(a1 + 92) = v4;
				}
				v4[5] = 0;
				++*(uint32_t*)(a1 + 96);
				return 1;
			}
			if (nox_strcmpi("FOREACH", (const char*)getMemAt(0x5D4594, 2487264))) {
				break;
			}
			if (!v2) {
				return 0;
			}
			if (*(uint32_t*)v2 != 1) {
				return 0;
			}
			v5 = nox_xxx_gen_5205B0(a2);
			if (!v5) {
				return 0;
			}
			v5[2] = *(uint32_t*)(v2 + 84);
			*(uint32_t*)(v2 + 84) = v5;
		}
		if (!nox_strcmpi("CONTAINS", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (v2) {
				if (!*(uint32_t*)v2) {
					v6 = nox_xxx_gen_520380(a2);
					*(uint32_t*)(v2 + 80) = v6;
					if (v6) {
						continue;
					}
				}
			}
			return 0;
		}
		v7 = calloc(1u, 0x64u);
		v2 = (int)v7;
		if (!v7) {
			return 0;
		}
		*v7 = 0;
		if (*getMemU32Ptr(0x587000, 253216)) {
			do {
				if (!nox_strcmpi(*(const char**)getMemAt(0x587000, 253216 + 4 * *(uint32_t*)v2),
							  (const char*)getMemAt(0x5D4594, 2487264))) {
					break;
				}
				v8 = *(uint32_t*)v2 + 1;
				*(uint32_t*)v2 = v8;
			} while (*getMemU32Ptr(0x587000, 253216 + 4 * v8));
		}
		if (!*getMemU32Ptr(0x587000, 253216 + 4 * *(uint32_t*)v2)) {
			free((void*)v2);
			return 0;
		}
		switch (*(uint32_t*)v2) {
		case 0:
		case 1:
			if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				free((void*)v2);
				return 0;
			}
			strcpy((char*)(v2 + 4), (const char*)getMemAt(0x5D4594, 2487264));
			if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				free((void*)v2);
				return 0;
			}
			if (nox_strcmpi("DENSITY", (const char*)getMemAt(0x5D4594, 2487264))) {
				*(uint32_t*)(v2 + 64) = 0;
				*(uint32_t*)(v2 + 68) = atoi((const char*)getMemAt(0x5D4594, 2487264));
				if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
					free((void*)v2);
					return 0;
				}
				*(uint32_t*)(v2 + 72) = atoi((const char*)getMemAt(0x5D4594, 2487264));
				v4 = (int*)v12;
				goto LABEL_39;
			}
			*(uint32_t*)(v2 + 64) = 1;
			if (nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				*(float*)(v2 + 76) = atof((const char*)getMemAt(0x5D4594, 2487264));
				if (nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
					*(uint32_t*)(v2 + 68) = !strcmp("*", (const char*)getMemAt(0x5D4594, 2487264))
												? 0
												: atoi((const char*)getMemAt(0x5D4594, 2487264));
					if (nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
						if (!strcmp("*", (const char*)getMemAt(0x5D4594, 2487264))) {
							*(uint32_t*)(v2 + 72) = 999999;
						} else {
							*(uint32_t*)(v2 + 72) = atoi((const char*)getMemAt(0x5D4594, 2487264));
						}
						v4 = (int*)v12;
						goto LABEL_39;
					}
				}
			}
			free((void*)v2);
			return 0;
		case 3:
		case 4:
		case 5:
			if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
				free((void*)v2);
				return 0;
			}
			strcpy((char*)(v2 + 4), (const char*)getMemAt(0x5D4594, 2487264));
			v4 = (int*)v12;
			break;
		default:
			break;
		}
	LABEL_39:
		*(uint32_t*)(v2 + 88) = v4[2];
		v9 = v4[3] + 1;
		v4[2] = v2;
		v4[3] = v9;
	}
}

//----- (00520380) --------------------------------------------------------
uint32_t* nox_xxx_gen_520380(FILE* a1) {
	uint32_t* v1;      // ebx
	uint32_t* v2;      // ebp
	int v3;            // eax
	int v4;            // edi
	const char** v5;   // eax
	unsigned char* v6; // esi
	int v7;            // ecx
	int v8;            // eax
	uint32_t* v9;      // eax
	uint32_t* i;       // ecx
	uint32_t* v11;     // ecx
	uint32_t* v13;     // eax
	uint32_t* v14;     // esi
	int v15;           // [esp+10h] [ebp-Ch]
	int v16;           // [esp+14h] [ebp-8h]
	int v17;           // [esp+18h] [ebp-4h]

	v1 = 0;
	v17 = 0;
	v15 = 0;
	v16 = 100;
LABEL_2:
	v2 = calloc(1u, 0x80Cu);
	if (v2 && nox_xxx_mapGenReadLine_51E540(a1, getMemAt(0x5D4594, 2487264))) {
		if (!strcmp("*", (const char*)getMemAt(0x5D4594, 2487264))) {
			*v2 = -1;
			++v15;
		} else {
			v3 = atoi((const char*)getMemAt(0x5D4594, 2487264));
			*v2 = v3;
			v16 -= v3;
		}
		while (nox_xxx_mapGenReadLine_51E540(a1, getMemAt(0x5D4594, 2487264))) {
			if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
				v17 = 1;
			} else if (!nox_strcmpi("OR", (const char*)getMemAt(0x5D4594, 2487264))) {
				// nop
			} else {
				goto LABEL_A;
			}
			if (v1) {
				v9 = (uint32_t*)v1[514];
				for (i = v1; v9; v9 = (uint32_t*)v9[514]) {
					i = v9;
				}
				i[514] = v2;
			} else {
				v1 = v2;
			}
			if (v17) {
				if (v15) {
					v11 = v1;
					do {
						if ((int)*v11 < 0) {
							*v11 = v16 / v15 + 1;
						}
						v11 = (uint32_t*)v11[514];
					} while (v11);
				}
				return v1;
			}
			goto LABEL_2;
		LABEL_A:
			v4 = 0;
			if (*getMemU32Ptr(0x587000, 253216)) {
				v5 = (const char**)getMemAt(0x587000, 253216);
				v6 = getMemAt(0x587000, 253216);
				do {
					if (!nox_strcmpi(*v5, (const char*)getMemAt(0x5D4594, 2487264))) {
						break;
					}
					v7 = *((uint32_t*)v6 + 1);
					v6 += 4;
					++v4;
					v5 = (const char**)v6;
				} while (v7);
			}
			if (!*getMemU32Ptr(0x587000, 253216 + 4 * v4)) {
				free(v2);
				break;
			}
			v2[16 * v2[513] + 1] = v4;
			if (!nox_xxx_mapGenReadLine_51E540(a1, &v2[16 * v2[513] + 2])) {
				break;
			}
			v8 = v2[513] + 1;
			v2[513] = v8;
			if (v8 >= 32) {
				v2[513] = 31;
			}
		}
	}
	v13 = v1;
	if (v1) {
		do {
			v14 = (uint32_t*)v13[514];
			free(v13);
			v13 = v14;
		} while (v14);
	}
	return 0;
}

//----- (005205B0) --------------------------------------------------------
uint32_t* nox_xxx_gen_5205B0(FILE* a1) {
	uint32_t* result; // eax
	uint32_t* v2;     // esi
	uint32_t* v3;     // eax

	result = (uint32_t*)nox_xxx_mapGenReadLine_51E540(a1, getMemAt(0x5D4594, 2487264));
	if (result) {
		result = calloc(1u, 0xCu);
		v2 = result;
		if (result) {
			*result = nox_xxx_getNameId_4E3AA0((char*)getMemAt(0x5D4594, 2487264));
			if (nox_xxx_mapGenReadLine_51E540(a1, getMemAt(0x5D4594, 2487264))) {
				if (nox_strcmpi("CONTAINS", (const char*)getMemAt(0x5D4594, 2487264))) {
					free(v2);
					result = 0;
				} else {
					v3 = nox_xxx_gen_520380(a1);
					v2[1] = v3;
					if (v3) {
						result = v2;
					} else {
						free(v2);
						result = 0;
					}
				}
			} else {
				free(v2);
				result = 0;
			}
		}
	}
	return result;
}

//----- (00520660) --------------------------------------------------------
int nox_xxx_genDecorReadCopy_520660(uint32_t* a1, const char* a2, FILE* a3) {
	int result;     // eax
	int v4;         // eax
	int v5;         // eax
	const char* v6; // eax
	const char* v7; // ebx
	uint32_t* i;    // ebp
	uint32_t* v9;   // eax
	int v10;        // edx
	uint32_t* v11;  // ebp
	uint32_t* j;    // eax
	uint32_t* k;    // ebx
	uint32_t* v14;  // eax
	int v15;        // edx

	result = nox_xxx_mapGenReadLine_51E540(a3, getMemAt(0x5D4594, 2487264));
	if (!result) {
		return 0;
	}
	v4 = *((uint32_t*)a2 + 15);
	if (v4) {
		v5 = v4 - 1;
		if (v5) {
			if (v5 == 1) {
				v6 = (const char*)a1[38];
			} else {
				v6 = a2;
			}
		} else {
			v6 = (const char*)a1[30];
		}
	} else {
		v6 = (const char*)a1[22];
	}
	v7 = v6;
	if (v6) {
		while (strcmp(v7, (const char*)getMemAt(0x5D4594, 2487264))) {
			v7 = (const char*)*((uint32_t*)v7 + 55);
			if (!v7) {
				break;
			}
		}
	}
	if (!v7) {
		v7 = (const char*)a1[38];
		if (!v7) {
			return 0;
		}
		while (strcmp(v7, (const char*)getMemAt(0x5D4594, 2487264))) {
			v7 = (const char*)*((uint32_t*)v7 + 55);
			if (!v7) {
				return 0;
			}
		}
	}
	for (i = (uint32_t*)*((uint32_t*)v7 + 21); i; i = (uint32_t*)i[31]) {
		v9 = calloc(1u, 0x80u);
		if (!v9) {
			return 0;
		}
		memcpy(v9, i, 0x80u);
		v9[31] = *((uint32_t*)a2 + 21);
		v10 = *((uint32_t*)a2 + 22) + 1;
		*((uint32_t*)a2 + 21) = v9;
		*((uint32_t*)a2 + 22) = v10;
	}
	v11 = (uint32_t*)*((uint32_t*)a2 + 23);
	if (v11) {
		for (j = (uint32_t*)v11[5]; j; j = (uint32_t*)j[5]) {
			v11 = j;
		}
	}
	for (k = (uint32_t*)*((uint32_t*)v7 + 23); k; v11 = v14) {
		v14 = calloc(1u, 0x18u);
		if (!v14) {
			return 0;
		}
		v15 = 1;
		memcpy(v14, k, 0x18u);
		v14[4] = 1;
		v14[5] = 0;
		if (v11) {
			v11[5] = v14;
			v15 = *((uint32_t*)a2 + 24) + 1;
		} else {
			*((uint32_t*)a2 + 23) = v14;
		}
		*((uint32_t*)a2 + 24) = v15;
		k = (uint32_t*)k[5];
	}
	return 1;
}

//----- (00520810) --------------------------------------------------------
int nox_xxx_genDecorReadOccurConstraint_520810(int a1, FILE* a2) {
	char* v3;          // ebx
	int v4;            // edi
	const char** v5;   // eax
	unsigned char* v6; // esi
	int v7;            // ecx

	*(uint8_t*)(a1 + 64) = 0;
	if (!nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
		return 0;
	}
	if (nox_strcmpi("NONE", (const char*)getMemAt(0x5D4594, 2487264))) {
		v3 = strtok((char*)getMemAt(0x5D4594, 2487264), "+");
		if (v3) {
			while (1) {
				v4 = 0;
				if (*getMemU32Ptr(0x587000, 253144)) {
					v5 = (const char**)getMemAt(0x587000, 253144);
					v6 = getMemAt(0x587000, 253144);
					do {
						if (!nox_strcmpi(*v5, v3)) {
							break;
						}
						v7 = *((uint32_t*)v6 + 1);
						v6 += 4;
						++v4;
						v5 = (const char**)v6;
					} while (v7);
				}
				if (!*getMemU32Ptr(0x587000, 253144 + 4 * v4)) {
					break;
				}
				*(uint8_t*)(a1 + 64) |= 1 << v4;
				v3 = strtok(0, "+");
				if (!v3) {
					return 1;
				}
			}
			return 0;
		}
	}
	return 1;
}

//----- (005208D0) --------------------------------------------------------
int nox_xxx_genDecorReadOccurLimit_5208D0(int a1, FILE* a2) {
	int result; // eax

	result = nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264));
	if (result) {
		*(uint8_t*)(a1 + 65) = atoi((const char*)getMemAt(0x5D4594, 2487264));
		result = 1;
	}
	return result;
}

//----- (00520910) --------------------------------------------------------
int nox_xxx_genDecorReadFrequency_520910(int a1, FILE* a2) {
	int result; // eax

	result = nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264));
	if (result) {
		if (nox_strcmpi("COMMON", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (nox_strcmpi("UNCOMMON", (const char*)getMemAt(0x5D4594, 2487264))) {
				if (nox_strcmpi("RARE", (const char*)getMemAt(0x5D4594, 2487264))) {
					if (nox_strcmpi("VERY_RARE", (const char*)getMemAt(0x5D4594, 2487264))) {
						if (!nox_strcmpi("HARDLY_EVER", (const char*)getMemAt(0x5D4594, 2487264))) {
							*(uint32_t*)(a1 + 72) = 1;
						}
						result = 1;
					} else {
						result = 1;
						*(uint32_t*)(a1 + 72) = 10;
					}
				} else {
					*(uint32_t*)(a1 + 72) = 100;
					result = 1;
				}
			} else {
				result = 1;
				*(uint32_t*)(a1 + 72) = 500;
			}
		} else {
			result = 1;
			*(uint32_t*)(a1 + 72) = 1000;
		}
	}
	return result;
}

//----- (005209F0) --------------------------------------------------------
int nox_xxx_genDecorReadRoomSizeCon_5209F0(int a1, FILE* a2) {
	int result; // eax

	result = nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264));
	if (result) {
		*(uint32_t*)(a1 + 76) = nox_strcmpi("*", (const char*)getMemAt(0x5D4594, 2487264))
									? atoi((const char*)getMemAt(0x5D4594, 2487264))
									: 0;
		result = nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264));
		if (result) {
			if (nox_strcmpi("*", (const char*)getMemAt(0x5D4594, 2487264))) {
				*(uint32_t*)(a1 + 80) = atoi((const char*)getMemAt(0x5D4594, 2487264));
			} else {
				*(uint32_t*)(a1 + 80) = 999999;
			}
			result = 1;
		}
	}
	return result;
}

//----- (00520A90) --------------------------------------------------------
int nox_xxx_genDecorReadDoor_520A90(int a1, FILE* a2) {
	return nox_xxx_mapGenReadLine_51E540(a2, (uint8_t*)(a1 + 100)) != 0;
}

//----- (00520AB0) --------------------------------------------------------
int nox_xxx_genDecorReadDoubleDoor_520AB0(int a1, FILE* a2) {
	return nox_xxx_mapGenReadLine_51E540(a2, (uint8_t*)(a1 + 160)) != 0;
}

//----- (00520AD0) --------------------------------------------------------
int nox_xxx_mapgenCheckSettings_520AD0(int* a1) {
	int v1; // ecx

	v1 = *a1;
	a1[2] = 0;
	a1[3] = 0;
	a1[4] = 0;
	a1[5] = 0;
	a1[6] = 0;
	for (a1[7] = 0; v1; v1 = *(uint32_t*)(v1 + 220)) {
		if (!*(uint8_t*)(v1 + 64)) {
			a1[2] += *(uint32_t*)(v1 + 72);
			a1[3] += *(uint32_t*)(v1 + 72);
			a1[4] += *(uint32_t*)(v1 + 72);
			a1[5] += *(uint32_t*)(v1 + 72);
			a1[6] += *(uint32_t*)(v1 + 72);
			a1[7] += *(uint32_t*)(v1 + 72);
		}
		if (*(uint8_t*)(v1 + 64) & 1) {
			a1[2] += *(uint32_t*)(v1 + 72);
		}
		if (*(uint8_t*)(v1 + 64) & 2) {
			a1[3] += *(uint32_t*)(v1 + 72);
		}
		if (*(uint8_t*)(v1 + 64) & 4) {
			a1[4] += *(uint32_t*)(v1 + 72);
		}
		if (*(uint8_t*)(v1 + 64) & 8) {
			a1[5] += *(uint32_t*)(v1 + 72);
		}
		if (*(uint8_t*)(v1 + 64) & 0x10) {
			a1[6] += *(uint32_t*)(v1 + 72);
		}
		if (*(uint8_t*)(v1 + 64) & 0x20) {
			a1[7] += *(uint32_t*)(v1 + 72);
		}
	}
	if (!a1[2]) {
		return 0;
	}
	if (!a1[3]) {
		return 0;
	}
	if (!a1[4]) {
		return 0;
	}
	if (!a1[5]) {
		return 0;
	}
	if (a1[6]) {
		return a1[7] != 0;
	}
	return 0;
}

//----- (00520BF0) --------------------------------------------------------
int nox_xxx_genReadPrefab_520BF0(int a1, FILE* a2) {
	char* v2;     // esi
	int v3;       // ebx
	char* v4;     // eax
	uint32_t* v5; // eax

	v2 = 0;
	v3 = 0;
	while (nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264))) {
		if (!nox_strcmpi("END", (const char*)getMemAt(0x5D4594, 2487264))) {
			return 1;
		}
		if (nox_strcmpi("AREAMAP", (const char*)getMemAt(0x5D4594, 2487264))) {
			if (nox_strcmpi("MUST_OCCUR", (const char*)getMemAt(0x5D4594, 2487264))) {
				if (!nox_strcmpi("FOREACH", (const char*)getMemAt(0x5D4594, 2487264))) {
					if (!v2) {
						return 0;
					}
					v5 = nox_xxx_gen_5205B0(a2);
					if (!v5) {
						return 0;
					}
					v5[2] = *((uint32_t*)v2 + 38);
					*((uint32_t*)v2 + 38) = v5;
				}
			} else {
				v3 = 1;
			}
		} else {
			v4 = sub_520CE0(a1, a2);
			v2 = v4;
			if (!v4) {
				return 0;
			}
			if (v3) {
				*((uint32_t*)v4 + 18) = 1;
				v3 = 0;
			}
		}
	}
	return 0;
}

//----- (00520CE0) --------------------------------------------------------
char* sub_520CE0(int a1, FILE* a2) {
	char* result; // eax
	int v3;       // ecx

	result = (char*)nox_xxx_mapGenReadLine_51E540(a2, getMemAt(0x5D4594, 2487264));
	if (result) {
		result = (char*)calloc(1u, 0xA0u);
		if (result) {
			strcpy(result, (const char*)getMemAt(0x5D4594, 2487264));
			*((uint32_t*)result + 39) = *(uint32_t*)(a1 + 80);
			v3 = *(uint32_t*)(a1 + 84) + 1;
			*(uint32_t*)(a1 + 80) = result;
			*(uint32_t*)(a1 + 84) = v3;
		}
	}
	return result;
}

//----- (00520D50) --------------------------------------------------------
uint32_t* sub_520D50(uint32_t* a1) {
	uint32_t* v1;     // eax
	uint32_t* v2;     // esi
	uint32_t* v3;     // eax
	uint32_t* v4;     // esi
	uint32_t* v5;     // eax
	uint32_t* v6;     // esi
	uint32_t* v7;     // eax
	uint32_t* v8;     // esi
	uint32_t* result; // eax
	uint32_t* v10;    // esi

	v1 = (uint32_t*)a1[22];
	if (v1) {
		do {
			v2 = (uint32_t*)v1[55];
			free(v1);
			v1 = v2;
		} while (v2);
	}
	v3 = (uint32_t*)a1[30];
	if (v3) {
		do {
			v4 = (uint32_t*)v3[55];
			free(v3);
			v3 = v4;
		} while (v4);
	}
	v5 = (uint32_t*)a1[275];
	if (v5) {
		do {
			v6 = (uint32_t*)v5[38];
			nox_xxx_mapGenFreeStr_51F1F0(v5);
			v5 = v6;
		} while (v6);
	}
	v7 = (uint32_t*)a1[277];
	if (v7) {
		do {
			v8 = (uint32_t*)v7[38];
			nox_xxx_mapGenFreeStr_51F1F0(v7);
			v7 = v8;
		} while (v8);
	}
	result = (uint32_t*)a1[20];
	if (result) {
		do {
			v10 = (uint32_t*)result[39];
			free(result);
			result = v10;
		} while (v10);
	}
	return result;
}

//----- (00521C60) --------------------------------------------------------
int nox_xxx_mapgen_521C60(int a1, int a2) {
	int result; // eax
	int i;      // esi
	int v4;     // eax

	result = sub_504980();
	for (i = result; result; i = result) {
		v4 = a2;
		if (a2) {
			while (1) {
				if (*(unsigned short*)(i + 4) == *(uint32_t*)v4) {
					nox_xxx_mapgen_521FE0(a1, i, *(uint32_t**)(v4 + 4));
					break;
				}
				v4 = *(uint32_t*)(v4 + 8);
				if (!v4) {
					break;
				}
			}
		}
		result = sub_5049C0(i);
	}
	return result;
}

//----- (00521CB0) --------------------------------------------------------
int sub_521CB0(int a1, int a2, int a3, int a4) {
	float* v4;     // ebx
	int v5;        // ebp
	int v6;        // esi
	signed int v7; // edi
	signed int v8; // ebp
	int v9;        // eax
	int v10;       // eax
	double v11;    // st7
	int v12;       // eax
	float* v14;    // [esp-4h] [ebp-30h]
	float v15;     // [esp+10h] [ebp-1Ch]
	float v16;     // [esp+14h] [ebp-18h]
	char* v17;     // [esp+18h] [ebp-14h]
	int v18;       // [esp+1Ch] [ebp-10h]
	int v19;       // [esp+20h] [ebp-Ch]
	float2 v20;    // [esp+24h] [ebp-8h]
	int v22;       // [esp+34h] [ebp+8h]

	v4 = (float*)calloc(1u, 0x1Cu);
	if (!v4) {
		return 0;
	}
	v15 = sub_502E70(a4);
	v16 = sub_502EA0(a4);
	v5 = (long long)(v15 * 0.030743772);
	v18 = v5;
	v6 = a2;
	v19 = (long long)(v16 * 0.030743772);
	v7 = *(uint32_t*)(a2 + 12) - v5;
	v8 = *(uint32_t*)(a2 + 16) - v19;
	if (v7 >= 0 && v8 >= 0) {
		v17 = sub_526AA0(a4);
		v20.field_0 = (double)nox_xxx_mapGenRandFunc_526AC0(0, v7) * 32.526913 + *(float*)(a2 + 36);
		v20.field_4 = (double)nox_xxx_mapGenRandFunc_526AC0(0, v8) * 32.526913 + *(float*)(a2 + 40);
		v9 = *((uint32_t*)v17 + 15);
		if (v9 & 1) {
			v20.field_4 = *(float*)(a2 + 40);
		} else if (v9 & 2) {
			v20.field_4 = *(float*)(a2 + 48) - v16;
		}
		v10 = *((uint32_t*)v17 + 15);
		if (v10 & 4) {
			v20.field_0 = *(float*)(a2 + 44) - v15;
		} else if (v10 & 8) {
			v20.field_0 = *(float*)(a2 + 36);
		}
		if (v17[60] & 0x10) {
			v22 = (*(uint32_t*)(a2 + 16) - v19) / 2;
			v20.field_0 = (double)(int)((*(uint32_t*)(v6 + 12) - v18) / 2) * 32.526913 + *(float*)(v6 + 36);
			v20.field_4 = (double)v22 * 32.526913 + *(float*)(v6 + 40);
		}
		*(uint32_t*)v4 = 1;
		v4[1] = v20.field_0;
		v4[2] = v20.field_4;
		v4[3] = v20.field_0 + v15;
		v11 = v20.field_4 + v16;
		*((uint32_t*)v4 + 5) = a4;
		v4[4] = v11;
		v12 = sub_521EB0((float*)v6, v4);
		v14 = v4;
		if (!v12) {
			goto LABEL_18;
		}
		if (!sub_521F10(v6, v4)) {
			sub_502D70(a4);
			nox_xxx_mapgen_521C60(a1, *(uint32_t*)(a3 + 84));
			sub_503B30(&v20);
			v4[6] = *(float*)(v6 + 368);
			*(uint32_t*)(v6 + 368) = v4;
			return 1;
		}
	}
	v14 = v4;
LABEL_18:
	free(v14);
	return 0;
}

//----- (00521FE0) --------------------------------------------------------
void nox_xxx_mapgen_521FE0(int a1, int a2, uint32_t* a3) {
	uint32_t* v3;  // ebp
	signed int v4; // eax
	uint32_t* v5;  // edi
	char* v6;      // esi
	uint32_t* v7;  // eax
	int v8;        // [esp+10h] [ebp+Ch]

	v3 = a3;
	if (a3) {
		v4 = nox_xxx_mapGenRandFunc_526AC0(1, 100);
		while (1) {
			v4 -= *v3;
			if (v4 <= 0) {
				break;
			}
			v3 = (uint32_t*)v3[514];
			if (!v3) {
				return;
			}
		}
		if (v3) {
			v8 = 0;
			if (v3[513] > 0) {
				v5 = 0;
				v6 = (char*)(v3 + 2);
				do {
					switch (*((uint32_t*)v6 - 1)) {
					case 0:
						v7 = nox_xxx_newObjectByTypeID_4E3810(v6);
						v5 = v7;
						break;
					case 3:
						v7 = nox_xxx_mapGenMakeEnchantedItem_5221A0(v6, *(char**)(a1 + 1100), *(uint32_t*)(a1 + 1104));
						v5 = v7;
						break;
					case 4:
						v7 = nox_xxx_mapGenMakeEnchantedItem_5221A0(v6, *(char**)(a1 + 1108), *(uint32_t*)(a1 + 1112));
						v5 = v7;
						break;
					case 5:
						v7 = nox_xxx_mapGenMakeSpellbook_5220E0(a1, v6);
						v5 = v7;
						break;
					default:
						break;
					}
					if (v5) {
						sub_522300(a2, v5);
					}
					v6 += 64;
					++v8;
				} while (v8 < v3[513]);
			}
		}
	}
}

//----- (005220E0) --------------------------------------------------------
uint32_t* nox_xxx_mapGenMakeSpellbook_5220E0(int a1, const char* a2) {
	uint32_t* v3; // ebx
	char v4;      // [esp+10h] [ebp+4h]

	if (!*(uint32_t*)(a1 + 1096)) {
		return 0;
	}
	v3 = nox_xxx_newObjectByTypeID_4E3810("SpellBook");
	if (!v3) {
		return v3;
	}
	if (!strcmp("*", a2)) {
		v4 = *(uint8_t*)(a1 + 4 * nox_xxx_mapGenRandFunc_526AC0(0, *(uint32_t*)(a1 + 1096) - 1) + 548);
		nox_xxx_mapGenFinishSpellbook_527DB0((int)v3, v4);
		return v3;
	}
	v4 = nox_xxx_mapGenSpellIdByName_51E1D0(a2);
	if (v4) {
		nox_xxx_mapGenFinishSpellbook_527DB0((int)v3, v4);
		return v3;
	}
	nox_xxx_objectFreeMem_4E38A0((nox_object_t*)v3);
	return 0;
}

//----- (005221A0) --------------------------------------------------------
uint32_t* nox_xxx_mapGenMakeEnchantedItem_5221A0(char* a1, char* a2, int a3) {
	signed int v3;   // eax
	char* v4;        // esi
	int i;           // ecx
	char* v6;        // eax
	uint32_t* v7;    // ebp
	uint32_t* v8;    // esi
	int v9;          // ebx
	const char* v10; // edi
	int v11;         // eax
	char v13[20];    // [esp+10h] [ebp-14h]

	if (!strcmp("*", a1)) {
		v3 = nox_xxx_mapGenRandFunc_526AC0(0, a3 - 1);
		v4 = a2;
		for (i = 0; v4; ++i) {
			if (i == v3) {
				break;
			}
			v4 = (char*)*((uint32_t*)v4 + 38);
		}
		v6 = v4 + 60;
	} else {
		v4 = a2;
		if (a2) {
			while (1) {
				if (!nox_strcmpi(v4, a1)) {
					v6 = v4 + 60;
					break;
				}
				v4 = (char*)*((uint32_t*)v4 + 38);
				if (!v4) {
					break;
				}
			}
		}
		if (!v4) {
			v6 = a1;
		}
	}
	v7 = nox_xxx_newObjectByTypeID_4E3810(v6);
	if (v7 && v4) {
		v8 = v4 + 136;
		*(uint32_t*)v13 = 0;
		v9 = 0;
		*(uint32_t*)&v13[4] = 0;
		*(uint32_t*)&v13[8] = 0;
		*(uint32_t*)&v13[12] = 0;
		do {
			if (*v8) {
				if (nox_xxx_mapGenRandFunc_526AC0(1, 100) <= *getMemU32Ptr(0x587000, 254688 + v9)) {
					v10 = (const char*)(*(v8 - 4) + 60 * nox_xxx_mapGenRandFunc_526AC0(0, *v8 - 1));
					if (nox_strcmpi("none", v10)) {
						v11 = nox_xxx_modifGetIdByName_413290(v10);
						*(uint32_t*)&v13[v9] = nox_xxx_modifGetDescById_413330(v11);
					}
				}
			}
			v9 += 4;
			++v8;
		} while (v9 < 16);
		if (!*(uint32_t*)&v13[8]) {
			*(uint32_t*)&v13[12] = 0;
		}
		nox_xxx_modifSetItemAttrs_4E4990((int)v7, (int*)v13);
	}
	return v7;
}

//----- (00522300) --------------------------------------------------------
uint32_t* sub_522300(int a1, uint32_t* a2) {
	uint32_t* result; // eax
	int v3;           // edx

	result = a2;
	a2[125] = 0;
	a2[124] = *(uint32_t*)(a1 + 504);
	v3 = *(uint32_t*)(a1 + 504);
	if (v3) {
		*(uint32_t*)(v3 + 500) = a2;
	}
	*(uint32_t*)(a1 + 504) = a2;
	a2[123] = a1;
	return result;
}

//----- (00522340) --------------------------------------------------------
void nox_xxx_mapgen_522340(int a1, int a2) {
	int* i; // esi

	for (i = *(int**)(*(uint32_t*)(a2 + 372) + 92); i; i = (int*)i[5]) {
		nox_xxx_mapgen_522370(a1, a2, i);
	}
}

//----- (00522370) --------------------------------------------------------
void nox_xxx_mapgen_522370(int a1, int a2, int* a3) {
	int* v3;        // ebp
	signed int v4;  // ecx
	int v5;         // eax
	int v6;         // edi
	uint32_t* v7;   // ebx
	uint32_t* v8;   // esi
	int v9;         // eax
	int v10;        // eax
	signed int v11; // esi
	int v12;        // [esp+10h] [ebp-4h]
	signed int v13; // [esp+20h] [ebp+Ch]

	v3 = a3;
	if (a3[2]) {
		v4 = nox_xxx_mapGenRandFunc_526AC0(*a3, a3[1]);
		v5 = a3[3];
		v13 = v4;
		if (v4 > v5) {
			v13 = v5;
			v4 = v5;
		}
		v6 = v3[2];
		v7 = (uint32_t*)v6;
		v12 = 1;
		*(uint32_t*)(v6 + 92) = 0;
		*(uint32_t*)(v6 + 96) = 0;
		v8 = *(uint32_t**)(v3[2] + 88);
		if (v3[3] > 1) {
			do {
				if (nox_xxx_mapGenRandFunc_526AC0(1, 100) >= 50) {
					if (nox_xxx_mapGenRandFunc_526AC0(1, 100) >= 50) {
						v7[23] = v8;
						v8[24] = v7;
						v8[23] = 0;
						v7 = v8;
					} else {
						v8[23] = v7;
						v8[24] = v7[24];
						v10 = v7[24];
						if (v10) {
							*(uint32_t*)(v10 + 92) = v8;
						} else {
							v6 = (int)v8;
						}
						v7[24] = v8;
					}
				} else if (nox_xxx_mapGenRandFunc_526AC0(1, 100) >= 50) {
					v9 = *(uint32_t*)(v6 + 92);
					if (v9) {
						*(uint32_t*)(v9 + 96) = v8;
						v8[24] = v6;
						v8[23] = *(uint32_t*)(v6 + 92);
					} else {
						v8[24] = v6;
						v8[23] = *(uint32_t*)(v6 + 92);
						v7 = v8;
					}
					*(uint32_t*)(v6 + 92) = v8;
				} else {
					v8[23] = v6;
					v8[24] = 0;
					*(uint32_t*)(v6 + 96) = v8;
					v6 = (int)v8;
				}
				v8 = (uint32_t*)v8[22];
				++v12;
			} while (v12 < v3[3]);
			v4 = v13;
		}
		if (v4) {
			v11 = v4;
			do {
				nox_xxx_mapgen_5224B0(a1, a2, v6);
				v6 = *(uint32_t*)(v6 + 92);
				--v11;
			} while (v11);
		}
	}
}

//----- (005224B0) --------------------------------------------------------
float* nox_xxx_mapgen_5224B0(int a1, int a2, int a3) {
	int v3;         // esi
	float* v4;      // ebp
	signed int v5;  // edi
	signed int v6;  // eax
	float* result;  // eax
	int v8;         // eax
	float* v9;      // ebp
	int v10;        // ebx
	float v11[2];        // [esp+10h] [ebp-8h]
	int v12;        // [esp+20h] [ebp+8h]
	signed int v13; // [esp+24h] [ebp+Ch]

	v3 = a3;
	if (*(uint32_t*)(a3 + 64)) {
		v4 = (float*)a2;
		v5 = (long long)((double)(int)(*(uint32_t*)(a2 + 12) * *(uint32_t*)(a2 + 16)) * *(float*)(a3 + 76));
		v6 = *(uint32_t*)(a3 + 68);
		if (v5 < v6 || (v6 = *(uint32_t*)(a3 + 72), v5 > v6)) {
			v5 = v6;
		}
	} else {
		v4 = (float*)a2;
		v5 = nox_xxx_mapGenRandFunc_526AC0(*(uint32_t*)(a3 + 68), *(uint32_t*)(a3 + 72));
	}
	result = *(float**)a3;
	switch (*(uint32_t*)a3) {
	case 0:
		if (v5 > 0) {
			v12 = v5;
			do {
				v8 = nox_xxx_mapGenMakeMonsterInRoom_522810(v4, (char*)(a3 + 4));
				if (v8) {
					nox_xxx_mapgen_521FE0(a1, v8, *(uint32_t**)(a3 + 80));
				}
				result = (float*)--v12;
			} while (v12);
		}
		break;
	case 1:
		result = (float*)sub_5268F0((const char*)(a3 + 4));
		v9 = result;
		if ((int)result >= 0 && v5 > 0) {
			v13 = v5;
			do {
				v10 = 0;
				do {
					if (sub_521CB0(a1, a2, v3, (int)v9)) {
						break;
					}
					++v10;
				} while (v10 < 3);
				result = (float*)--v13;
			} while (v13);
		}
		break;
	case 2:
		result = (float*)sub_521C10((int)v4);
		break;
	case 3:
		result = (float*)sub_5226D0((int)v4, 0.89999998, (int)v11);
		if (result) {
			result = (float*)nox_xxx_mapGenMakeEnchantedItem_5221A0((char*)(a3 + 4), *(char**)(a1 + 1100),
																	*(uint32_t*)(a1 + 1104));
			if (result) {
				result = nox_xxx_mapGenMoveObject_527A10(result, v11);
			}
		}
		break;
	case 4:
		result = (float*)sub_5226D0((int)v4, 0.89999998, (int)v11);
		if (result) {
			result = (float*)nox_xxx_mapGenMakeEnchantedItem_5221A0((char*)(a3 + 4), *(char**)(a1 + 1108),
																	*(uint32_t*)(a1 + 1112));
			if (result) {
				result = nox_xxx_mapGenMoveObject_527A10(result, v11);
			}
		}
		break;
	case 5:
		result = (float*)sub_5226D0((int)v4, 0.89999998, (int)v11);
		if (result) {
			result = (float*)nox_xxx_mapGenMakeSpellbook_5220E0(a1, (const char*)(a3 + 4));
			if (result) {
				result = nox_xxx_mapGenMoveObject_527A10(result, v11);
			}
		}
		break;
	default:
		return result;
	}
	return result;
}

//----- (00522810) --------------------------------------------------------
int nox_xxx_mapGenMakeMonsterInRoom_522810(float* a1, char* a2) {
	int result; // eax
	float* v3;  // eax
	int v4;     // esi
	int v5;     // eax
	int v6;     // eax
	float v7;   // [esp+0h] [ebp-18h]
	float v8;   // [esp+4h] [ebp-14h]
	float2 v9;  // [esp+8h] [ebp-10h]
	float2 v11; // [esp+10h] [ebp-8h]

	v7 = (a1[11] + a1[9]) * 0.5;
	v8 = (a1[12] + a1[10]) * 0.5;
	result = sub_5226D0((int)a1, 0.94999999, &v9);
	if (result) {
		nox_xxx_mapGenGetObjID_527940(a2);
		v3 = nox_xxx_mapGenPlaceObj_5279B0(&v9);
		v4 = (int)v3;
		if (v3) {
			if ((uint8_t)v3[2] & 2) {
				v11.field_0 = v7 - v9.field_0;
				v11.field_4 = v8 - v9.field_4;
				v5 = nox_xxx_math_509ED0(&v11);
				v6 = nox_xxx_math_509EA0(v5);
				nox_xxx_mapGenOrientObj_527C60(v4, v6);
			}
		}
		result = v4;
	}
	return result;
}

//----- (00522C80) --------------------------------------------------------
float* sub_522C80(float* a1) {
	float* result; // eax

	result = sub_51D1A0((float2*)a1);
	if (!result) {
		result = (float*)sub_51D120(a1);
	}
	return result;
}

//----- (00522D30) --------------------------------------------------------
float* sub_522D30(int a1) {
	float* result; // eax
	float* i;      // esi
	int* v3;       // edi
	int v4;        // ebx
	int v5;        // ebp
	float** v6;    // ebx
	int v7;        // edi
	int v8;        // ebp
	int v9;        // [esp+4h] [ebp-14h]
	int v10;       // [esp+4h] [ebp-14h]
	float2 v11;    // [esp+8h] [ebp-10h]
	float2 v12;    // [esp+10h] [ebp-8h]

	sub_51D0F0(128);
	result = (float*)nox_xxx_mapGenGetTopRoom_521710();
	for (i = result; result; i = result) {
		nox_xxx_mapGenSetFlags_5235F0(156);
		if (*(uint32_t*)i == 1) {
			v3 = (int*)(i + 22);
			v4 = 0;
			v9 = (int)(i + 22);
			do {
				v5 = 0;
				if (*((uint8_t*)i + v4 + 216)) {
					do {
						if (nox_xxx_mapGenCheckRoomType_5238F0((int*)*v3) && sub_523920(*(uint32_t*)*v3) == v4) {
							sub_523CB0(*v3, (int)&v12);
							sub_522CA0((int)i, &v12.field_0);
							sub_523D30((float*)*v3, &v11.field_0);
							sub_522C80(&v12.field_0);
							sub_522C80(&v11.field_0);
							sub_51D3F0(&v12, &v11);
							if (!*(uint32_t*)(a1 + 60)) {
								sub_51D3F0(&v11, &v12);
							}
						}
						++v5;
						++v3;
					} while (v5 < *((unsigned char*)i + v4 + 216));
				}
				++v4;
				v3 = (int*)(v9 + 32);
				v9 += 32;
			} while (v4 < 4);
		} else {
			sub_523D30(i, &v12.field_0);
			v6 = (float**)(i + 22);
			v7 = 0;
			v10 = (int)(i + 22);
			do {
				if (!sub_521AA0(i, v7)) {
					v8 = 0;
					if (*((uint8_t*)i + v7 + 216)) {
						do {
							if (nox_xxx_mapGenCheckRoomType_5238F0((int*)*v6)) {
								sub_523D30(*v6, &v11.field_0);
							} else {
								sub_523C30((int)i, (int)&v11);
								sub_522CA0((int)*v6, &v11.field_0);
							}
							sub_522C80(&v12.field_0);
							sub_522C80(&v11.field_0);
							sub_51D3F0(&v12, &v11);
							if (!*(uint32_t*)(a1 + 60)) {
								sub_51D3F0(&v11, &v12);
							}
							++v8;
							++v6;
						} while (v8 < *((unsigned char*)i + v7 + 216));
					}
				}
				++v7;
				v6 = (float**)(v10 + 32);
				v10 += 32;
			} while (v7 < 4);
		}
		result = (float*)sub_521720((int)i);
	}
	return result;
}

//----- (00522F40) --------------------------------------------------------
unsigned char* nox_xxx_mapGenTryNextRoom_522F40(uint32_t* a1) {
	unsigned char* result; // eax
	unsigned char* i;      // edi
	int v3;                // eax
	int v4;                // ebp
	int v5;                // esi
	float2* v6;            // ebx
	float2* v7;            // [esp+8h] [ebp+4h]

	result = (unsigned char*)a1;
	if (!a1[15]) {
		result = (unsigned char*)nox_xxx_mapGenGetTopRoom_521710();
		for (i = result; result; i = result) {
			if (*(uint32_t*)i == 1) {
				nox_xxx_mapGenSetFlags_5235F0(156);
				v3 = i[352];
				v4 = 0;
				if (v3 > 0) {
					v7 = (float2*)(i + 224);
					do {
						v5 = 0;
						if (v3 > 0) {
							v6 = (float2*)(i + 224);
							do {
								if (v4 != v5) {
									sub_51D3F0(v7, v6);
								}
								++v5;
								++v6;
							} while (v5 < i[352]);
						}
						v3 = i[352];
						++v4;
						++v7;
					} while (v4 < v3);
				}
			}
			result = (unsigned char*)sub_521720((int)i);
		}
	}
	return result;
}

//----- (00522FF0) --------------------------------------------------------
int nox_xxx_netSendPointFx_522FF0(char a1, float2* a2) {
	short v2;    // ax
	float v3;    // edx
	char a2a[5]; // [esp+4h] [ebp-8h]

	a2a[0] = a1;
	v2 = nox_float2int(a2->field_0);
	v3 = a2->field_4;
	*(uint16_t*)&a2a[1] = v2;
	*(uint16_t*)&a2a[3] = nox_float2int(v3);
	return nox_xxx_netSendFxAllCli_523030(a2, a2a, 5);
}

//----- (00523030) --------------------------------------------------------
int nox_xxx_netSendFxAllCli_523030(float2* a1, const void* a2, int a3) {
	int result; // eax
	int i;      // esi
	int v5;     // ecx
	int v6;     // eax
	float v7;   // edx
	float v8;   // eax
	double v9;  // st7
	double v10; // st7
	double v11; // st6
	float v12;  // [esp+4h] [ebp-18h]
	float v13;  // [esp+Ch] [ebp-10h]
	float v14;  // [esp+10h] [ebp-Ch]
	float v15;  // [esp+18h] [ebp-4h]

	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		v5 = *(uint32_t*)(*(uint32_t*)(i + 748) + 276);
		if (*(uint8_t*)(v5 + 3680) & 3 && (v6 = *(uint32_t*)(v5 + 3628)) != 0) {
			v7 = *(float*)(v6 + 56);
			v8 = *(float*)(v6 + 60);
			v12 = v7;
		} else {
			v8 = *(float*)(i + 60);
			v12 = *(float*)(i + 56);
		}
		v9 = (double)*(unsigned short*)(v5 + 10);
		v13 = v12 - v9 - 50.0;
		v10 = v9 + v12 + 50.0;
		v11 = (double)*(unsigned short*)(v5 + 12);
		if (a1->field_0 > (double)v13 && v10 > a1->field_0) {
			v14 = v8 - v11 - 50.0;
			if (a1->field_4 > (double)v14) {
				v15 = v11 + v8 + 50.0;
				if (a1->field_4 < (double)v15) {
					nox_netlist_addToMsgListCli_40EBC0(*(unsigned char*)(v5 + 2064), 1, a2, a3);
				}
			}
		}
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (00523150) --------------------------------------------------------
int sub_523150(char a1, char a2, float* a3) {
	char v4[6]; // [esp+4h] [ebp-8h]

	v4[0] = a1;
	v4[1] = a2;
	*(uint16_t*)&v4[2] = nox_float2int(*a3);
	*(uint16_t*)&v4[4] = nox_float2int(a3[1]);
	return nox_xxx_netSendFxAllCli_523030((float2*)a3, v4, 6);
}

//----- (005231B0) --------------------------------------------------------
int nox_xxx_netSparkExplosionFx_5231B0(float* a1, char a2) {
	short v2;   // ax
	float v3;   // ecx
	char v5[6]; // [esp+4h] [ebp-8h]

	v5[0] = -109;
	v2 = nox_float2int(*a1);
	v3 = a1[1];
	*(uint16_t*)&v5[1] = v2;
	*(uint16_t*)&v5[3] = nox_float2int(v3);
	v5[5] = a2;
	return nox_xxx_netSendFxAllCli_523030((float2*)a1, v5, 6);
}

//----- (00523200) --------------------------------------------------------
void nox_xxx_sendGeneratorBreakFX_523200(float* a1, char a2) {
	short v2;   // ax
	float v3;   // ecx
	char v5[7]; // [esp+4h] [ebp-8h]

	v5[0] = -16;
	v5[1] = 25;
	v2 = nox_float2int(*a1);
	v3 = a1[1];
	*(uint16_t*)&v5[2] = v2;
	*(uint16_t*)&v5[4] = nox_float2int(v3);
	v5[6] = a2;
	nox_xxx_netSendFxAllCli_523030((float2*)a1, v5, 7);
}

//----- (00523270) --------------------------------------------------------
int nox_xxx_netSendVampFx_523270(char a1, short* a2, short a3) {
	short v3;          // dx
	unsigned short v4; // cx
	unsigned short v5; // ax
	float2 a1a;        // [esp+0h] [ebp-14h]
	char a2a[11];      // [esp+8h] [ebp-Ch]

	a2a[0] = a1;
	v3 = a2[2];
	*(uint16_t*)&a2a[1] = *a2;
	v4 = a2[4];
	v5 = a2[6];
	*(uint16_t*)&a2a[3] = v3;
	*(uint16_t*)&a2a[5] = v4;
	*(uint16_t*)&a2a[7] = v5;
	a1a.field_0 = (double)v4;
	*(uint16_t*)&a2a[9] = a3;
	a1a.field_4 = (double)v5;
	return nox_xxx_netSendFxAllCli_523030(&a1a, a2a, 11);
}

//----- (00523530) --------------------------------------------------------
int nox_xxx_netClientPredictLinear_523530(int a1) {
	short v1;      // ax
	double v2;     // st7
	long long v3;  // rax
	double v4;     // st7
	long long v5;  // rax
	double v6;     // st7
	short v7;      // cx
	long long v8;  // rax
	double v9;     // st7
	long long v10; // rax
	double v11;    // st7
	int result;    // eax
	int i;         // esi
	char v14[14];  // [esp+4h] [ebp-10h]

	v14[0] = -75;
	v1 = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)a1);
	v2 = *(float*)(a1 + 56);
	*(uint16_t*)&v14[1] = v1;
	*(uint16_t*)&v14[3] = *(uint16_t*)(a1 + 4);
	v3 = (long long)v2;
	v4 = *(float*)(a1 + 60);
	*(uint16_t*)&v14[5] = v3;
	v5 = (long long)v4;
	v6 = *(float*)(a1 + 112) * 16.0;
	v7 = *(uint16_t*)(a1 + 124);
	*(uint16_t*)&v14[7] = v5;
	*(uint16_t*)&v14[9] = v7;
	v8 = (long long)v6;
	v9 = *(float*)(a1 + 80) * 16.0;
	v14[11] = v8;
	v10 = (long long)v9;
	v11 = *(float*)(a1 + 84) * 16.0;
	v14[12] = v10;
	v14[13] = (long long)v11;
	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		nox_netlist_addToMsgListCli_40EBC0(*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064), 1, v14,
										   14);
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (005235F0) --------------------------------------------------------
void nox_xxx_mapGenSetFlags_5235F0(char a1) {
	int v3; // [esp+0h] [ebp-4h]

	if (!nox_common_gameFlags_check_40A5C0(0x200000)) {
		dword_5d4594_2487564 = nox_platform_get_ticks();
		if (dword_5d4594_2487568 > *(int*)&dword_5d4594_2487564) {
			dword_5d4594_2487568 = 0;
		}
		nox_input_pollEvents_4453A0();
		if (dword_5d4594_2487564 - dword_5d4594_2487568 > *getMemIntPtr(0x587000, 254948)) {
			*(uint16_t*)((char*)&v3 + 1) = *getMemU16Ptr(0x5D4594, 2487572);
			LOBYTE(v3) = a1;
			++*getMemU32Ptr(0x5D4594, 2487572);
			nox_xxx_mapGenClientText_4A9D00((unsigned char*)&v3);
			dword_5d4594_2487568 = dword_5d4594_2487564;
		}
	}
}

//----- (00523670) --------------------------------------------------------
int nox_xxx_netSendShieldFx_523670(int a1, float* a2) {
	char v2;    // al
	int v3;     // eax
	char v5[4]; // [esp+4h] [ebp-Ch]
	float2 v6;  // [esp+8h] [ebp-8h]

	v5[0] = -128;
	*(uint16_t*)&v5[1] = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)a1);
	if (a2) {
		v6.field_0 = *(float*)(a1 + 56) - *a2;
		v6.field_4 = *(float*)(a1 + 60) - a2[1];
		v3 = nox_xxx_math_509ED0(&v6);
		v2 = nox_xxx_math_509EA0(v3);
	} else {
		v2 = nox_xxx_math_509EA0(*(short*)(a1 + 124));
	}
	v5[3] = v2;
	return nox_xxx_netSendFxAllCli_523030((float2*)(a1 + 56), v5, 4);
}

//----- (005236F0) --------------------------------------------------------
int nox_xxx_sendSummonStartFX_5236F0(short a1, float* a2, char a3, short a4, short a5) {
	double v5;    // st7
	long long v6; // rax
	double v7;    // st7
	char v9[12];  // [esp+4h] [ebp-Ch]

	v9[0] = 126;
	*(uint16_t*)&v9[5] = a1;
	*(uint16_t*)&v9[7] = a4;
	v5 = *a2;
	v9[9] = a3;
	v6 = (long long)v5;
	v7 = a2[1];
	*(uint16_t*)&v9[1] = v6;
	*(uint16_t*)&v9[3] = (long long)v7;
	*(uint16_t*)&v9[10] = a5;
	return nox_xxx_netSendPacket0_4E5420(255, v9, 12, 0, 1);
}

//----- (00523760) --------------------------------------------------------
int nox_xxx_sendSummonCancelFX_523760(short a1) {
	char v3[3]; // [esp+0h] [ebp-4h]
	v3[0] = 127;
	*(uint16_t*)&v3[1] = a1;
	return nox_xxx_netSendPacket0_4E5420(255, v3, 3, 0, 1);
}

//----- (00523830) --------------------------------------------------------
void nox_xxx_sendGeneratorSpawnFX_523830(int4* a1, short a2) {
	double v2;    // st7
	short v3;     // cx
	short v4;     // dx
	double v5;    // st7
	short v6;     // cx
	float2 a1a;   // [esp+0h] [ebp-14h]
	char a2a[12]; // [esp+8h] [ebp-Ch]

	a2a[0] = -16;
	a2a[1] = 16;
	v2 = (double)a1->field_8;
	v3 = a1->field_0;
	*(uint16_t*)&a2a[4] = a1->field_4;
	v4 = a1->field_C;
	a1a.field_0 = v2;
	v5 = (double)a1->field_C;
	*(uint16_t*)&a2a[2] = v3;
	v6 = a1->field_8;
	*(uint16_t*)&a2a[8] = v4;
	*(uint16_t*)&a2a[6] = v6;
	a1a.field_4 = v5;
	*(uint16_t*)&a2a[10] = a2;
	nox_xxx_netSendFxAllCli_523030(&a1a, a2a, 12);
}

//----- (005238A0) --------------------------------------------------------
void nox_xxx_sendArrowTrapFX_5238A0(float* a1, char a2) {
	short v2;   // ax
	float v3;   // ecx
	char v5[6]; // [esp+4h] [ebp-8h]

	v5[0] = -95;
	v2 = nox_float2int16(*a1);
	v3 = a1[1];
	*(uint16_t*)&v5[1] = v2;
	*(uint16_t*)&v5[3] = nox_float2int16(v3);
	v5[5] = a2;
	nox_xxx_netSendFxAllCli_523030((float2*)a1, v5, 6);
}

int sub_526CA0(char* a1);

//----- (005259F0) --------------------------------------------------------
void sub_5259F0(int a1, int a2, float a3) {
	long double v3; // st7
	int v4;         // ebp
	int* v5;        // edi
	int v6;         // eax
	uint8_t* v7;    // esi
	int v8;         // ebx
	int v9;         // [esp+4h] [ebp+4h]
	int v10;        // [esp+8h] [ebp+8h]
	float v11;      // [esp+Ch] [ebp+Ch]

	nox_xxx_mapGenSetFlags_5235F0(156);
	if (a2) {
		switch (*(uint32_t*)a2) {
		case 1:
			v3 = sqrt(*(float*)(a2 + 28) * *(float*)(a2 + 28) + *(float*)(a2 + 32) * *(float*)(a2 + 32));
			break;
		case 2:
		case 3:
			v3 = *(float*)(a2 + 32);
			break;
		case 4:
		case 5:
			v3 = *(float*)(a2 + 28);
			break;
		default:
			v3 = a3;
			break;
		}
	} else {
		v3 = 0.0;
	}
	v4 = a1;
	v11 = v3 + a3;
	if (*(uint8_t*)(a1 + 220) != 1 || v11 < (double)*(float*)(a1 + 356)) {
		*(float*)(a1 + 356) = v11;
		v5 = (int*)(a1 + 88);
		v6 = -216 - a1;
		*(uint8_t*)(a1 + 220) = 1;
		v7 = (uint8_t*)(a1 + 216);
		v10 = a1 + 88;
		v9 = -216 - a1;
		do {
			v8 = 0;
			if (*v7) {
				do {
					if (*v5) {
						sub_5259F0(*v5, v4, v11);
					}
					++v8;
					++v5;
				} while (v8 < (unsigned char)*v7);
				v6 = v9;
			}
			++v7;
			v5 = (int*)(v10 + 32);
			v10 += 32;
		} while ((int)&v7[v6] < 4);
	}
}

//----- (00525AF0) --------------------------------------------------------
float* sub_525AF0(int a1) {
	int v1;        // esi
	int v2;        // edi
	float* i;      // eax
	float* result; // eax
	int v5;        // ecx
	int v6;        // eax

	v1 = 0;
	v2 = 0;
	dword_5d4594_2487580 = 0;
	dword_5d4594_2487576 = a1;
	sub_525BF0(a1);
	*(uint32_t*)(dword_5d4594_2487576 + 52) |= 1u;
	for (i = (float*)nox_xxx_mapGenGetTopRoom_521710(); i; i = (float*)sub_521720((int)i)) {
		++v1;
		i[90] = i[89] / *(float*)&dword_5d4594_2487580;
	}
	dword_5d4594_2487584 = 0;
	result = sub_525C90();
	v5 = a1;
	if (a1) {
		do {
			if (v1 > 1) {
				v6 = v2 / (v1 - 1);
			} else {
				v6 = 100;
			}
			if (v5 == a1) {
				*(uint32_t*)(v5 + 364) = 1;
			} else if (v6 >= 30) {
				if (v6 >= 60) {
					if (v6 >= 90) {
						*(uint32_t*)(v5 + 364) = 16;
					} else {
						*(uint32_t*)(v5 + 364) = 8;
					}
				} else {
					*(uint32_t*)(v5 + 364) = 4;
				}
			} else {
				*(uint32_t*)(v5 + 364) = 2;
			}
			v5 = *(uint32_t*)(v5 + 64);
			v2 += 100;
		} while (v5);
		result = *(float**)&dword_5d4594_2487576;
		*(uint32_t*)(dword_5d4594_2487576 + 364) = 32;
	} else {
		*(uint32_t*)(dword_5d4594_2487576 + 364) = 32;
	}
	return result;
}

//----- (00525BF0) --------------------------------------------------------
void sub_525BF0(int a1) {
	int v1;      // eax
	int v2;      // eax
	uint8_t* v3; // esi
	int v4;      // ebp
	int v5;      // ebx
	int* v6;     // edi
	int v7;      // [esp+4h] [ebp+4h]

	if (*(uint8_t*)(a1 + 220) != 2) {
		if (*(float*)&dword_5d4594_2487580 < (double)*(float*)(a1 + 356) && *(uint32_t*)a1 == 1) {
			v1 = *(uint32_t*)(a1 + 356);
			dword_5d4594_2487576 = a1;
			dword_5d4594_2487580 = v1;
		}
		v2 = -216 - a1;
		*(uint8_t*)(a1 + 220) = 2;
		v3 = (uint8_t*)(a1 + 216);
		v4 = a1 + 88;
		v7 = -216 - a1;
		do {
			v5 = 0;
			if (*v3) {
				v6 = (int*)v4;
				do {
					if (*v6) {
						sub_525BF0(*v6);
					}
					++v5;
					++v6;
				} while (v5 < (unsigned char)*v3);
				v2 = v7;
			}
			++v3;
			v4 += 32;
		} while ((int)&v3[v2] < 4);
	}
}

//----- (00525C90) --------------------------------------------------------
float* sub_525C90() {
	float* result; // eax
	float* i;      // edx
	int v2;        // ecx
	int v3;        // esi
	int v4;        // eax
	int v5;        // eax

	dword_5d4594_2487584 = 0;
	result = (float*)nox_xxx_mapGenGetTopRoom_521710();
	for (i = result; result; i = result) {
		v2 = dword_5d4594_2487584;
		if (dword_5d4594_2487584) {
			v3 = 0;
			while (1) {
				if (i[89] < (double)*(float*)(v2 + 356)) {
					v4 = *(uint32_t*)(v2 + 68);
					if (v4) {
						*(uint32_t*)(v4 + 64) = i;
					} else {
						dword_5d4594_2487584 = i;
					}
					v5 = *(uint32_t*)(v2 + 68);
					*((uint32_t*)i + 16) = v2;
					*((uint32_t*)i + 17) = v5;
					*(uint32_t*)(v2 + 68) = i;
					break;
				}
				v3 = v2;
				v2 = *(uint32_t*)(v2 + 64);
				if (!v2) {
					break;
				}
			}
			if (!v2) {
				*((uint32_t*)i + 17) = v3;
				i[16] = 0.0;
				*(uint32_t*)(v3 + 64) = i;
				goto LABEL_13;
			}
		} else {
			dword_5d4594_2487584 = i;
			i[16] = 0.0;
			i[17] = 0.0;
		}
	LABEL_13:
		result = (float*)sub_521720((int)i);
	}
	return result;
}

//----- (00525D20) --------------------------------------------------------
int nox_xxx_mapGen_InPrefab1_525D20(int a1) {
	int v1;      // edi
	uint32_t* i; // esi
	uint32_t* j; // esi

	v1 = 0;
	if (!dword_5d4594_2487656) {
		dword_5d4594_2487656 = nox_xxx_getNameId_4E3AA0("ExitNorthMarker");
		*getMemU32Ptr(0x5D4594, 2487660) = nox_xxx_getNameId_4E3AA0("ExitSouthMarker");
		*getMemU32Ptr(0x5D4594, 2487664) = nox_xxx_getNameId_4E3AA0("ExitEastMarker");
		*getMemU32Ptr(0x5D4594, 2487668) = nox_xxx_getNameId_4E3AA0("ExitWestMarker");
	}
	sub_525DF0(*(float*)&a1);
	for (i = *(uint32_t**)(a1 + 80); i; i = (uint32_t*)i[39]) {
		if (i[18]) {
			if (v1 == 5 || !nox_xxx_mapGenPrefabMkRoom_526100(a1, (int)i)) {
				return 0;
			}
			i[19] = 1;
			++v1;
		}
	}
	for (j = *(uint32_t**)(a1 + 80); j; j = (uint32_t*)j[39]) {
		if (!j[18]) {
			if (v1 == 5) {
				return 1;
			}
			if (!nox_xxx_mapGenPrefabMkRoom_526100(a1, (int)j)) {
				return 0;
			}
			j[19] = 1;
			++v1;
		}
	}
	return 1;
}

//----- (00525DF0) --------------------------------------------------------
void sub_525DF0(float a1) {
	double v1;         // st7
	int v2;            // eax
	int v3;            // ecx
	unsigned char* v4; // edx
	int v5;            // esi
	signed int v6;     // eax
	int v7;            // ecx
	double v8;         // st7
	signed int v9;     // eax
	double v10;        // st7
	double v11;        // st7
	double v12;        // st7
	double v13;        // st7
	float v14;         // [esp+8h] [ebp-4h]
	float v15;         // [esp+8h] [ebp-4h]
	float v16;         // [esp+10h] [ebp+4h]

	v1 = *(float*)(LODWORD(a1) + 64) * 0.5;
	*getMemU32Ptr(0x5D4594, 2487608) = 0;
	v2 = *(uint32_t*)(LODWORD(a1) + 84);
	v16 = v1;
	if (v2 >= 5) {
		v2 = 5;
		dword_5d4594_2487652 = 5;
	} else {
		dword_5d4594_2487652 = v2;
		if (v2 == 1) {
			*getMemU32Ptr(0x5D4594, 2487588) = 0;
			goto LABEL_10;
		}
		if (v2 == 2) {
			if (nox_xxx_mapGenRandFunc_526AC0(0, 100) >= 50) {
				v2 = dword_5d4594_2487652;
				*getMemU32Ptr(0x5D4594, 2487588) = 1;
				*getMemU32Ptr(0x5D4594, 2487592) = 0;
			} else {
				v2 = dword_5d4594_2487652;
				*getMemU32Ptr(0x5D4594, 2487588) = 0;
				*getMemU32Ptr(0x5D4594, 2487592) = 1;
			}
			goto LABEL_10;
		}
	}
	v3 = 0;
	if (v2 > 0) {
		v4 = getMemAt(0x5D4594, 2487588);
		do {
			*(uint32_t*)v4 = v3++;
			v4 += 4;
		} while (v3 < v2);
	}
	v5 = 0;
	if (v2 > 0) {
		do {
			v6 = nox_xxx_mapGenRandFunc_526AC0(0, v2 - 2);
			v7 = *getMemU32Ptr(0x5D4594, 2487588 + 4 * v6);
			*getMemU32Ptr(0x5D4594, 2487588 + 4 * v6) = *getMemU32Ptr(0x5D4594, 2487592 + 4 * v6);
			*getMemU32Ptr(0x5D4594, 2487592 + 4 * v6) = v7;
			v2 = dword_5d4594_2487652;
			++v5;
		} while (v5 < *(int*)&dword_5d4594_2487652);
	}
LABEL_10:
	switch (v2) {
	case 1:
		v14 = -v16;
		*getMemFloatPtr(0x5D4594, 2487612) = sub_526BC0(v14, v16);
		*getMemFloatPtr(0x5D4594, 2487616) = sub_526BC0(v14, v16);
		break;
	case 2:
		v8 = -v16;
		v15 = v8;
		if (nox_xxx_mapGenRandFunc_526AC0(0, 100) >= 50) {
			*getMemFloatPtr(0x5D4594, 2487612) = v8;
			*getMemFloatPtr(0x5D4594, 2487616) = sub_526BC0(v15, v16);
			*(float*)&dword_5d4594_2487620 = v16;
			*(float*)&dword_5d4594_2487624 = sub_526BC0(v15, v16);
		} else {
			*getMemFloatPtr(0x5D4594, 2487612) = sub_526BC0(v15, v16);
			*getMemFloatPtr(0x5D4594, 2487616) = v15;
			*(float*)&dword_5d4594_2487620 = sub_526BC0(v15, v16);
			*(float*)&dword_5d4594_2487624 = v16;
		}
		break;
	case 3:
		v9 = nox_xxx_mapGenRandFunc_526AC0(0, 100);
		if (v9 >= 25) {
			if (v9 >= 50) {
				v12 = -v16;
				if (v9 >= 75) {
					*getMemFloatPtr(0x5D4594, 2487612) = v12;
					*getMemFloatPtr(0x5D4594, 2487616) = v12;
					*(float*)&dword_5d4594_2487620 = v12;
					*(float*)&dword_5d4594_2487624 = v16;
					*(float*)&dword_5d4594_2487628 = v16;
				} else {
					*getMemFloatPtr(0x5D4594, 2487616) = v12;
					*(float*)&dword_5d4594_2487628 = v12;
					*getMemFloatPtr(0x5D4594, 2487612) = v16;
					*(float*)&dword_5d4594_2487620 = v16;
					*(float*)&dword_5d4594_2487624 = v16;
				}
				dword_5d4594_2487632 = 0;
			} else {
				v11 = -v16;
				*getMemFloatPtr(0x5D4594, 2487612) = v11;
				*(float*)&dword_5d4594_2487632 = v11;
				*getMemFloatPtr(0x5D4594, 2487616) = v16;
				*(float*)&dword_5d4594_2487620 = v16;
				*(float*)&dword_5d4594_2487624 = v16;
				dword_5d4594_2487628 = 0;
			}
		} else {
			v10 = -v16;
			*getMemFloatPtr(0x5D4594, 2487612) = v10;
			*getMemFloatPtr(0x5D4594, 2487616) = v10;
			*(float*)&dword_5d4594_2487624 = v10;
			*(float*)&dword_5d4594_2487620 = v16;
			dword_5d4594_2487628 = 0;
			*(float*)&dword_5d4594_2487632 = v16;
		}
		break;
	case 4:
	case 5:
		*getMemU32Ptr(0x5D4594, 2487644) = 0;
		v13 = -v16;
		*getMemFloatPtr(0x5D4594, 2487612) = v13;
		*getMemFloatPtr(0x5D4594, 2487616) = v13;
		*(float*)&dword_5d4594_2487624 = v13;
		*(float*)&dword_5d4594_2487628 = v13;
		*(float*)&dword_5d4594_2487620 = v16;
		*(float*)&dword_5d4594_2487632 = v16;
		*getMemFloatPtr(0x5D4594, 2487636) = v16;
		*getMemFloatPtr(0x5D4594, 2487640) = v16;
		*getMemU32Ptr(0x5D4594, 2487648) = 0;
		break;
	default:
		return;
	}
}

//----- (00526100) --------------------------------------------------------
int nox_xxx_mapGenPrefabMkRoom_526100(int a1, int a2) {
	int result;   // eax
	float* v3;    // eax
	int v4;       // edi
	uint32_t* v5; // eax
	int v6;       // edx
	int v7;       // ecx
	uint32_t* v8; // eax
	float2 v9;    // [esp+Ch] [ebp-8h]

	result = sub_5262F0(a1, a2);
	if (result) {
		v3 = nox_xxx_mapGenMakeRoomStruct_521940(
			4 * *(uint32_t*)(a2 + 144) + (unsigned long long)(long long)(*(float*)(a2 + 60) * 0.030743772 + 0.5),
			4 * *(uint32_t*)(a2 + 144) + (unsigned long long)(long long)(*(float*)(a2 + 64) * 0.030743772 + 0.5));
		*(uint32_t*)(a2 + 148) = v3;
		sub_526260((int)v3, &v9.field_0);
		nox_xxx_mapGenSetRoomPos_521880(*(uint32_t**)(a2 + 148), &v9);
		v4 = 1;
		v5 = (uint32_t*)sub_521200(*(uint32_t*)(a2 + 148));
		if (v5 && !sub_5212B0(*(uint32_t*)(a2 + 148), v5)) {
			v4 = 0;
		}
		if (sub_5217A0(a1, *(uint32_t*)(a2 + 148)) && v4) {
			nox_xxx_mapGenAddNewRoom_521730(*(uint32_t**)(a2 + 148));
			v6 = 4;
			*(uint32_t*)(*(uint32_t*)(a2 + 148) + 52) |= 2u;
			v7 = 2 * *(uint32_t*)(a2 + 144);
			v8 = (uint32_t*)(a2 + 80);
			do {
				if (v8[3]) {
					*v8 += v7 + *(uint32_t*)(*(uint32_t*)(a2 + 148) + 4);
					v8[1] += v7 + *(uint32_t*)(*(uint32_t*)(a2 + 148) + 8);
				}
				v8 += 4;
				--v6;
			} while (v6);
			result = 1;
		} else {
			sub_521A10(*(void**)(a2 + 148));
			*(uint32_t*)(a2 + 148) = 0;
			result = 0;
		}
	}
	return result;
}

//----- (00526260) --------------------------------------------------------
long long sub_526260(int a1, float* a2) {
	int v2;           // ecx
	long long result; // rax

	v2 = *getMemU32Ptr(0x5D4594, 2487588 + 4 * *getMemU32Ptr(0x5D4594, 2487608));
	*a2 = *getMemFloatPtr(0x5D4594, 2487612 + 8 * v2);
	a2[1] = *getMemFloatPtr(0x5D4594, 2487616 + 8 * v2);
	++*getMemU32Ptr(0x5D4594, 2487608);
	*a2 = *a2 - *(float*)(a1 + 28) * 0.5;
	a2[1] = a2[1] - *(float*)(a1 + 32) * 0.5;
	*a2 = (double)(int)(long long)(*a2 * 0.030743772) * 32.526913;
	result = (long long)(a2[1] * 0.030743772);
	a2[1] = (double)(int)result * 32.526913;
	return result;
}

//----- (005262F0) --------------------------------------------------------
int sub_5262F0(int a1, int a2) {
	int v2;        // eax
	double v3;     // st7
	int v4;        // ecx
	int i;         // edi
	int v6;        // eax
	int v7;        // eax
	int v8;        // eax
	int v9;        // ecx
	int v10;       // eax
	int v11;       // eax
	int v12;       // ecx
	int v13;       // eax
	int v14;       // eax
	int v15;       // ecx
	int v16;       // eax
	int v17;       // eax
	int v18;       // ecx
	int v19;       // edi
	int v20;       // edx
	int v21;       // ebx
	uint32_t* v22; // eax
	int result;    // eax
	int2 v24;      // [esp+Ch] [ebp-18h]
	float2 v26;    // [esp+14h] [ebp-10h]
	float2 v28;    // [esp+1Ch] [ebp-8h]

	v2 = sub_5029A0((char*)a2);
	*(uint32_t*)(a2 + 68) = v2;
	if (v2 == -1) {
		return 0;
	}
	sub_502D70(v2);
	v3 = sub_502E70(*(uint32_t*)(a2 + 68));
	v4 = *(uint32_t*)(a2 + 68);
	*(float*)(a2 + 60) = v3;
	*(float*)(a2 + 64) = sub_502EA0(v4);
	for (i = sub_504980(); i; i = sub_5049C0(i)) {
		sub_503EC0(i, &v26);
		v28.field_0 = v26.field_0 - 1.0;
		v28.field_4 = v26.field_4 - 1.0;
		nox_xxx_mapGenRoundFloatToPtr_520DF0(&v28, &v24);
		v6 = *(unsigned short*)(i + 4);
		if ((unsigned short)v6 == dword_5d4594_2487656) {
			if (*(uint32_t*)(a2 + 92)) {
				if (v24.field_0 < *(int*)(a2 + 80)) {
					*(uint32_t*)(a2 + 80) = v24.field_0;
				}
				v7 = *(uint32_t*)(a2 + 88) + 1;
				*(uint32_t*)(a2 + 84) = v24.field_4;
				*(uint32_t*)(a2 + 88) = v7;
			} else {
				v8 = v24.field_0;
				v9 = v24.field_4;
				*(uint32_t*)(a2 + 92) = 1;
				*(uint32_t*)(a2 + 80) = v8;
				*(uint32_t*)(a2 + 84) = v9;
				*(uint32_t*)(a2 + 88) = 1;
			}
		} else if (v6 == *getMemU32Ptr(0x5D4594, 2487660)) {
			if (*(uint32_t*)(a2 + 108)) {
				if (v24.field_0 < *(int*)(a2 + 96)) {
					*(uint32_t*)(a2 + 96) = v24.field_0;
				}
				v10 = *(uint32_t*)(a2 + 104) + 1;
				*(uint32_t*)(a2 + 100) = v24.field_4;
				*(uint32_t*)(a2 + 104) = v10;
			} else {
				v11 = v24.field_0;
				v12 = v24.field_4;
				*(uint32_t*)(a2 + 108) = 1;
				*(uint32_t*)(a2 + 96) = v11;
				*(uint32_t*)(a2 + 100) = v12;
				*(uint32_t*)(a2 + 104) = 1;
			}
		} else if (v6 == *getMemU32Ptr(0x5D4594, 2487668)) {
			if (*(uint32_t*)(a2 + 140)) {
				if (v24.field_4 < *(int*)(a2 + 132)) {
					*(uint32_t*)(a2 + 132) = v24.field_4;
				}
				v13 = *(uint32_t*)(a2 + 136) + 1;
				*(uint32_t*)(a2 + 128) = v24.field_0;
				*(uint32_t*)(a2 + 136) = v13;
			} else {
				v14 = v24.field_0;
				v15 = v24.field_4;
				*(uint32_t*)(a2 + 140) = 1;
				*(uint32_t*)(a2 + 128) = v14;
				*(uint32_t*)(a2 + 132) = v15;
				*(uint32_t*)(a2 + 136) = 1;
			}
		} else if (v6 == *getMemU32Ptr(0x5D4594, 2487664)) {
			if (*(uint32_t*)(a2 + 124)) {
				if (v24.field_4 < *(int*)(a2 + 116)) {
					*(uint32_t*)(a2 + 116) = v24.field_4;
				}
				v16 = *(uint32_t*)(a2 + 120) + 1;
				*(uint32_t*)(a2 + 112) = v24.field_0;
				*(uint32_t*)(a2 + 120) = v16;
			} else {
				v17 = v24.field_0;
				v18 = v24.field_4;
				*(uint32_t*)(a2 + 124) = 1;
				*(uint32_t*)(a2 + 112) = v17;
				*(uint32_t*)(a2 + 116) = v18;
				*(uint32_t*)(a2 + 120) = 1;
			}
		}
	}
	v19 = 4;
	v20 = *(uint32_t*)(a1 + 12) + *(uint32_t*)(a1 + 16);
	v21 = 0;
	v22 = (uint32_t*)(a2 + 88);
	do {
		if (v22[1]) {
			if (*v22 > *(int*)(a2 + 144)) {
				*(uint32_t*)(a2 + 144) = *v22;
			}
			++v21;
		}
		v22 += 4;
		--v19;
	} while (v19);
	if (*(uint32_t*)(a2 + 144) > v20 || v21 > 2 ||
		(*(uint32_t*)(a2 + 92) || *(uint32_t*)(a2 + 108)) && (*(uint32_t*)(a2 + 124) || *(uint32_t*)(a2 + 140))) {
		return 0;
	} else {
		return 1;
	}
}

//----- (00526550) --------------------------------------------------------
int sub_526550(int a1, int a2) {
	int v2;           // ecx
	int v3;           // ebp
	int v4;           // edi
	int v5;           // ebx
	int v6;           // eax
	int* v7;          // esi
	int v8;           // ecx
	int v9;           // eax
	int v10;          // eax
	int v11;          // edx
	int v12;          // eax
	int* v13;         // edi
	int v14;          // eax
	int v15;          // ebp
	int v16;          // esi
	int v17;          // edi
	int v18;          // ebx
	long long v19;    // rax
	int* v20;         // esi
	int v21;          // ecx
	int v23;          // [esp+10h] [ebp-38h]
	int v24;          // [esp+14h] [ebp-34h]
	int v25;          // [esp+18h] [ebp-30h]
	int v26[5] = {0}; // [esp+1Ch] [ebp-2Ch]
	int v27;          // [esp+30h] [ebp-18h]
	int v28[5];       // [esp+34h] [ebp-14h]
	int v29;          // [esp+4Ch] [ebp+4h]

	v2 = 16 * (a2 + 5);
	v3 = 0;
	v4 = 0;
	v5 = 0;
	v23 = *(uint32_t*)(v2 + a1);
	v6 = *(uint32_t*)(v2 + a1 + 4);
	v27 = 0;
	v24 = v6;
	v29 = 0;
	v7 = (int*)nox_xxx_mapGenGetTopRoom_521710();
	if (v7) {
		do {
			if (!nox_xxx_mapGenCheckRoomType_5238F0(v7)) {
				v8 = v7[1] + v7[3] / 2 - v23;
				v9 = v7[2] + v7[4] / 2 - v24;
				switch (a2) {
				case 0:
					if (v9 < 0) {
						break;
					}
					goto LABEL_A;
				case 1:
					if (v9 > 0) {
						break;
					}
					goto LABEL_A;
				case 2:
					if (v8 > 0) {
						break;
					}
					goto LABEL_A;
				case 3:
					if (v8 < 0) {
						break;
					}
					goto LABEL_A;
				default:
					break;
				}
				v10 = v8 * v8 + v9 * v9;
				if (v3 >= 24) {
					if (*(&v25 + v5) > v10) {
						*(&v27 + v5) = (int)v7;
						*(&v25 + v5) = v10;
					}
				} else {
					++v4;
					*(int*)((char*)&v27 + v3) = (int)v7;
					*(int*)((char*)&v25 + v3) = v10;
					v29 = v4;
					v3 += 4;
				}
				v11 = v25;
				v5 = 0;
				v12 = 1;
				if (v3 > 4) {
					v13 = v26;
					do {
						if (*v13 > v11) {
							v5 = v12;
							v11 = *v13;
						}
						++v12;
						++v13;
					} while (v12 < v29);
					v4 = v29;
				}
			}
		LABEL_A:
			v14 = sub_521720((int)v7);
			v7 = (int*)v14;
		} while (v14);
		if (v4 > 1) {
			v15 = v4 - 1;
			if (v4 - 1 > 0) {
				do {
					v16 = *(&v25 + v14);
					v17 = v26[v14];
					if (v16 > v17) {
						v18 = *(&v27 + v14);
						*(&v25 + v14) = v17;
						*(&v27 + v14) = v28[v14];
						v26[v14] = v16;
						v28[v14] = v18;
						v14 = -1;
					}
					++v14;
				} while (v14 < v15);
				v4 = v29;
			}
		}
	}
	v19 = 0;
	if (v4 > 0) {
		v20 = &v27;
		do {
			v21 = *v20;
			*(uint32_t*)(*v20 + 72) = 0;
			if (HIDWORD(v19)) {
				*(uint32_t*)(HIDWORD(v19) + 72) = v21;
			} else {
				LODWORD(v19) = v21;
			}
			++v20;
			--v4;
			HIDWORD(v19) = v21;
		} while (v4);
	}
	return v19;
}
// 526550: using guessed type int var_2C[5];
// 526550: using guessed type int var_14[5];

//----- (005266F0) --------------------------------------------------------
int nox_xxx_mapGen_InPrefab2_5266F0(int a1) {
	int v1;       // esi
	int v2;       // ecx
	float* v3;    // eax
	uint32_t* v4; // ebp
	int v5;       // ebx
	uint32_t* v6; // edi
	float2 v8;    // [esp+10h] [ebp-8h]

	v1 = *(uint32_t*)(a1 + 80);
	if (!v1) {
		return 1;
	}
	while (1) {
		if (*(uint32_t*)(v1 + 76)) {
			v2 = *(uint32_t*)(v1 + 148);
			v8.field_0 = *(float*)(v2 + 20);
			v8.field_4 = *(float*)(v2 + 24);
			v8.field_0 = (double)*(int*)(v1 + 144) * 65.053825 + v8.field_0;
			v8.field_4 = (double)*(int*)(v1 + 144) * 65.053825 + v8.field_4;
			sub_521760(*(uint32_t*)(v1 + 148));
			sub_521A10(*(void**)(v1 + 148));
			v3 = nox_xxx_mapGenMakeRoomStruct_521940((long long)(*(float*)(v1 + 60) * 0.030743772 + 0.5),
													 (long long)(*(float*)(v1 + 64) * 0.030743772 + 0.5));
			*(uint32_t*)(v1 + 148) = v3;
			nox_xxx_mapGenSetRoomPos_521880(v3, &v8);
			nox_xxx_mapGenAddNewRoom_521730(*(uint32_t**)(v1 + 148));
			v4 = (uint32_t*)(v1 + 92);
			v5 = 0;
			*(uint32_t*)(*(uint32_t*)(v1 + 148) + 52) |= 2u;
			while (1) {
				if (*v4) {
					v6 = (uint32_t*)sub_526550(v1, v5);
					if (!v6) {
						return 0;
					}
					while (!sub_54B2D0((int*)v1, v5, v6)) {
						v6 = (uint32_t*)v6[18];
						if (!v6) {
							return 0;
						}
					}
					if (!v6) {
						return 0;
					}
				}
				++v5;
				v4 += 4;
				if (v5 >= 4) {
					goto LABEL_11;
				}
			}
		}
	LABEL_11:
		v1 = *(uint32_t*)(v1 + 156);
		if (!v1) {
			return 1;
		}
	}
}

//----- (00526830) --------------------------------------------------------
int nox_xxx_mapGenPlacePrefabs_526830(int a1) {
	int* v1; // ebx
	int v2;  // esi
	int v3;  // edi
	int v4;  // eax

	v1 = *(int**)(a1 + 80);
	if (!v1) {
		return 1;
	}
	do {
		if (v1[19]) {
			sub_502D70(v1[17]);
			v2 = sub_504980();
			if (v2) {
				do {
					v3 = sub_5049C0(v2);
					v4 = *(unsigned short*)(v2 + 4);
					if ((unsigned short)v4 == dword_5d4594_2487656 || v4 == *getMemU32Ptr(0x5D4594, 2487660) ||
						v4 == *getMemU32Ptr(0x5D4594, 2487664) || v4 == *getMemU32Ptr(0x5D4594, 2487668)) {
						sub_504A10(v2);
					}
					v2 = v3;
				} while (v3);
			}
			nox_xxx_mapgen_521C60(a1, v1[38]);
			sub_503B30((float2*)(v1[37] + 20));
		}
		v1 = (int*)v1[39];
	} while (v1);
	return 1;
}

//----- (005268F0) --------------------------------------------------------
int sub_5268F0(const char* a1) {
	int v1;        // ebp
	const char* i; // edi

	v1 = 0;
	if (*(int*)&dword_5d4594_2487676 <= 0) {
		return -1;
	}
	for (i = *(const char**)&dword_5d4594_2487672; strcmp(i, a1); i += 64) {
		if (++v1 >= *(int*)&dword_5d4594_2487676) {
			return -1;
		}
	}
	return v1;
}

//----- (00526950) --------------------------------------------------------
char* sub_526950() {
	int v0;           // eax
	char* result;     // eax
	int v2;           // ebp
	char* v3;         // ebx
	char* v4;         // esi
	int v5;           // edi
	unsigned char v6; // al
	int v7;           // [esp+4h] [ebp-104h]
	char v8[256];     // [esp+8h] [ebp-100h]

	v0 = sub_502A20();
	dword_5d4594_2487676 = v0;
	result = (char*)calloc(v0, 0x40u);
	v2 = 0;
	dword_5d4594_2487672 = result;
	if (result) {
		v7 = 0;
		if (dword_5d4594_2487676 > 0) {
			while (1) {
				v3 = &result[v2];
				strcpy(v8, (const char*)sub_5029F0(v7));
				strcpy((char*)(dword_5d4594_2487672 + v2), v8);
				strtok(v8, "-");
				strtok(0, "-");
				v4 = strtok(0, "-");
				if (v4) {
					v5 = 0;
					if (getMemByte(0x587000, 255032)) {
						do {
							if (strchr(v4, (char)getMemByte(0x587000, 255032 + v5))) {
								*((uint32_t*)v3 + 15) |= 1 << v5;
							}
							v6 = getMemByte(0x587000, 255033 + v5++);
						} while (v6);
					}
				}
				result = (char*)(v7 + 1);
				v2 += 64;
				if (++v7 >= *(int*)&dword_5d4594_2487676) {
					break;
				}
				result = *(char**)&dword_5d4594_2487672;
			}
		}
	} else {
		dword_5d4594_2487676 = 0;
	}
	return result;
}

//----- (00526AA0) --------------------------------------------------------
char* sub_526AA0(int a1) { return (char*)(dword_5d4594_2487672 + (a1 << 6)); }

void* nox_objectTypeGetXfer(char* id);

void* nox_objectTypeGetXfer(char* id);

//----- (00527E50) --------------------------------------------------------
int nox_xxx_netUpdateObjectSpecial_527E50(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	uint32_t* a2 = a2p;
	unsigned int v2; // edi
	int v3;          // eax

	v2 = *(unsigned char*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
	if (a2 && v2 < 0x20) {
		v3 = a2[v2 + 140];
		if (!(v3 & 0xFFF0000)) {
			return 0;
		}
		if ((v3 & 0x10000) == 0x10000) {
			nox_xxx_netReportAnimFrame_4D81F0(v2, a2);
			a2[v2 + 140] &= 0xFFFEFFFF;
		}
		if ((a2[v2 + 140] & 0x20000) == 0x20000) {
			if (!nox_xxx_unitIsEnemyTo_5330C0(a1, (int)a2)) {
				nox_xxx_netReportUnitCurrentHP_4D8620(v2, a2);
			}
			a2[v2 + 140] &= 0xFFFDFFFF;
		}
		if ((a2[v2 + 140] & 0x40000) == 0x40000) {
			nox_xxx_netReportObjHidden_4D8FD0(v2, a2);
			a2[v2 + 140] &= 0xFFFBFFFF;
		}
		if ((a2[v2 + 140] & 0x80000) == 0x80000) {
			nox_xxx_netReportXStatus_4D8230(v2, a2);
			a2[v2 + 140] &= 0xFFF7FFFF;
		}
		if ((a2[v2 + 140] & 0x400000) == 0x400000) {
			nox_xxx_netReportUnitHeight_4D9020(v2, (int)a2);
			a2[v2 + 140] &= 0xFFBFFFFF;
		}
		if ((0x800000 & a2[v2 + 140]) == 0x800000) {
			nox_xxx_netReportEnchant_4D8F90(v2, a2);
			a2[v2 + 140] &= 0xFF7FFFFF;
		}
		if ((a2[v2 + 140] & 0x2000000) == 0x2000000) {
			nox_xxx_netReportTeamBase_4D92D0(v2, (int)a2);
			a2[v2 + 140] &= 0xFDFFFFFF;
		}
		if ((a2[v2 + 140] & 0x4000000) == 0x4000000) {
			nox_xxx_netSendReportNPC_4D93A0(v2, (int)a2);
			a2[v2 + 140] &= 0xFBFFFFFF;
		}
	}
	return 1;
}

//----- (00528030) --------------------------------------------------------
short sub_528030(int a1) {
	int v1;             // ebp
	int v2;             // esi
	unsigned short* v3; // edi
	int v4;             // ebx
	int v5;             // edi
	int v6;             // edx
	unsigned int v7;    // ebx
	int v8;             // eax
	short result;       // ax
	int v10;            // edx
	int v11;            // [esp+18h] [ebp+4h]

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 748);
	v3 = *(unsigned short**)(a1 + 556);
	v4 = *(uint32_t*)(v2 + 276);
	v11 = *(unsigned char*)(v4 + 2064);
	if (*(uint16_t*)(v2 + 10) == *v3) {
		v5 = gameFrame();
		v7 = gameFPS();
	} else if (abs(*v3 - *(unsigned short*)(v2 + 10)) >= v3[2] / 10 ||
		(v5 = gameFrame(), v6 = *(uint32_t*)(v4 + 2176), v7 = gameFPS(),
		 (unsigned int)(gameFrame() - v6) > (int)gameFPS() >> 2)) {
		nox_xxx_netSendPlrHealthToTeam_4D86E0(v11);
		v8 = *(uint32_t*)(v2 + 276);
		*(uint16_t*)(v2 + 10) = **(uint16_t**)(v1 + 556);
		*(uint32_t*)(v8 + 2176) = gameFrame();
		v5 = gameFrame();
		v7 = gameFPS();
	}
	result = *(uint16_t*)(v2 + 4);
	if (*(uint16_t*)(v2 + 6) != result) {
		result = 0;
		if (abs(*(unsigned short*)(v2 + 4) - *(unsigned short*)(v2 + 6)) >= *(unsigned short*)(v2 + 8) / 10 ||
			v5 - *(uint32_t*)(*(uint32_t*)(v2 + 276) + 2180) > v7 >> 2) {
			nox_xxx_netReportMana_4D8930(v11, v1);
			v10 = *(uint32_t*)(v2 + 276);
			*(uint16_t*)(v2 + 6) = *(uint16_t*)(v2 + 4);
			result = (unsigned short)gameFrame();
			*(uint32_t*)(v10 + 2180) = gameFrame();
		}
	}
	return result;
}

//----- (00528190) --------------------------------------------------------
int nox_xxx_checkIsKillable_528190(nox_object_t* a1p) {
	int a1 = a1p;
	uint16_t* v1; // eax
	bool v2;      // zf

	v1 = *(uint16_t**)(a1 + 556);
	if (!v1) {
		return 0;
	}
	if (*v1) {
		v2 = v1[2] == 0;
		if (v1[2]) {
			return 1;
		}
	} else {
		v2 = v1[2] == 0;
	}
	if (v2) {
		return 1;
	} else {
		return 0;
	}
}

//----- (005281D0) --------------------------------------------------------
int nox_xxx_frameCounterSetCopyToNextFrame_5281D0() {
	int result; // eax

	result = gameFrame() + 1;
	*getMemU32Ptr(0x5D4594, 2487684) = gameFrame() + 1;
	return result;
}

//----- (005281E0) --------------------------------------------------------
int nox_xxx_frameCounterSetCopy_5281E0() {
	int result; // eax

	result = gameFrame();
	*getMemU32Ptr(0x5D4594, 2487684) = gameFrame();
	return result;
}

//----- (005281F0) --------------------------------------------------------
int nox_xxx_unitCanSee_536FB0(nox_object_t* a1, nox_object_t* a2, int a3);
void nox_xxx_unitUpdateSightMB_5281F0(nox_object_t* a1p) {
	uint32_t a1 = a1p;
	uint32_t v1;   // edi
	int v2;     // eax
	int v3;     // ebp
	double v4;  // st7
	int v5;     // esi
	int* v6;    // ebx
	double v7;  // st7
	double v8;  // st6
	double v9;  // st7
	double v10; // st6
	int v11;    // eax
	int v12;    // eax
	int v13;    // esi
	int v14;    // eax
	int v15;    // eax
	int v16;    // esi
	int v17;    // [esp+10h] [ebp-10h]
	float v18;  // [esp+10h] [ebp-10h]
	int v19;    // [esp+14h] [ebp-Ch]
	float v20;  // [esp+18h] [ebp-8h]
	float v21;  // [esp+24h] [ebp+4h]

	v1 = a1;
	v17 = 0;
	v2 = *(uint32_t*)(a1 + 16);
	v3 = *(uint32_t*)(a1 + 748);
	if ((v2 & 0x8000) != 0 && !nox_xxx_unitIsZombie_534A40(a1)) {
		return;
	}
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		v4 = 640.0;
	} else {
		v4 = 250.0;
	}
	if (v4 >= *(float*)(v3 + 1312)) {
		v21 = v4;
	} else {
		v21 = *(float*)(v3 + 1312);
	}
	if (gameFrame() - *(uint32_t*)(v3 + 1212) <= (unsigned int)(2 * gameFPS())) {
		v19 = 0;
	} else {
		v19 = 1;
		*(uint32_t*)(v3 + 1212) = gameFrame();
	}
	v5 = 0;
	if (*(uint8_t*)(v3 + 1129)) {
		v6 = (int*)(v3 + 1132);
		do {
			if (*(uint32_t*)(*v6 + 16) & 0x8020 || !nox_xxx_unitCanSee_536FB0(v1, *v6, 0) ||
				(v7 = *(float*)(v1 + 56) - *(float*)(*v6 + 56),
				 v8 = *(float*)(v1 + 60) - *(float*)(*v6 + 60), v20 = (v21 + 30.0) * (v21 + 30.0),
				 v8 * v8 + v7 * v7 > v20) ||
				(v9 = *(float*)(v1 + 56) - *(float*)(v1 + 72),
				 v10 = *(float*)(v1 + 60) - *(float*)(v1 + 76), v10 * v10 + v9 * v9 > 1000.0) ||
				v19 && !nox_xxx_unitCanInteractWith_5370E0(v1, *v6, 0)) {
				nox_xxx_aiLostSight_528560(v1, v5--);
				v17 = 1;
				--v6;
			}
			++v5;
			++v6;
		} while (v5 < *(unsigned char*)(v3 + 1129));
	}
	v11 = *(uint32_t*)(v3 + 1196);
	if (v11 && nox_xxx_testUnitBuffs_4FF350(v11, 28)) {
		v17 = 1;
	}
	if ((!*(uint32_t*)(v3 + 1196) ||
		 gameFrame() - *(uint32_t*)(v3 + 1204) > (unsigned int)(2 * gameFPS())) &&
		(*(uint32_t*)(v3 + 1208) <= gameFrame() ||
		 gameFrame() == *getMemU32Ptr(0x5D4594, 2487684))) {
		nox_xxx_unitsGetInCircle_517F90((float2*)(v1 + 56), v21, nox_xxx_monsterUpdateSeenEnemies_5286D0, v1);
		*(uint32_t*)(v3 + 1204) = gameFrame();
		*(uint32_t*)(v3 + 1212) = gameFrame();
		v17 = 1;
	}
	if (v17) {
		v12 = *(uint32_t*)(v3 + 1196);
		if (v12) {
			v13 = *(uint32_t*)(v12 + 36);
		} else {
			v13 = 0;
		}
		sub_528610(v1);
		v14 = *(uint32_t*)(v3 + 1196);
		if (v14 && v13 && v13 != *(uint32_t*)(v14 + 36)) {
			*(uint32_t*)(v3 + 1200) = v13;
		}
	}
	if (*(uint32_t*)(v3 + 1204) == gameFrame()) {
		v15 = *(uint32_t*)(v3 + 1440);
		if (v15 & 0x400 || nox_common_gameFlags_check_40A5C0(0x2000) || *(uint32_t*)(v3 + 1196)) {
			*(uint32_t*)(v3 + 1208) = gameFrame() + nox_common_randomInt_415FA0(5, 10);
		} else {
			v16 = 5 * gameFPS();
			v18 = sub_5336D0(v1);
			*(float*)(v3 + 524) = v18;
			if (v18 < 0.0) {
				*(uint32_t*)(v3 + 1208) = v16 + gameFrame();
			} else if (v18 > (double)v21) {
				*(uint32_t*)(v3 + 1208) = (unsigned long long)(long long)((v18 - v21) * (double)v16 / (1000.0 - v21)) +
										  10 + gameFrame();
			} else {
				*(uint32_t*)(v3 + 1208) = nox_common_randomInt_415FA0(5, 10) + gameFrame();
			}
		}
	}
}

//----- (00528560) --------------------------------------------------------
int nox_xxx_aiLostSight_528560(int a1, int a2) {
	int v2;     // esi
	int v3;     // eax
	int* v4;    // edi
	int v5;     // eax
	int v6;     // eax
	int v7;     // edx
	int v8;     // ecx
	int result; // eax
	int v10;    // [esp-4h] [ebp-14h]

	v2 = *(uint32_t*)(a1 + 748);
	v3 = *(uint32_t*)(v2 + 4 * a2 + 1132);
	v4 = (int*)(v2 + 4 * a2 + 1132);
	v10 = *(uint32_t*)(v3 + 36);
	v5 = nox_xxx_getUnitName_4E39D0(v3);
	nox_ai_debug_printf_5341A0("%d: Lost sight of %s(#%d)\n", gameFrame(), v5, v10);
	nox_xxx_scriptCallByEventBlock_502490((int*)(v2 + 1296), *v4, a1, 15);
	v6 = *(uint32_t*)(v2 + 1196);
	if (*v4 == v6) {
		v7 = *(uint32_t*)(v6 + 36);
		*(uint32_t*)(v2 + 1196) = 0;
		*(uint32_t*)(v2 + 1200) = v7;
	}
	v8 = a2;
	LOBYTE(result) = *(uint8_t*)(v2 + 1129) - 1;
	*(uint8_t*)(v2 + 1129) = result;
	result = (unsigned char)result;
	if (a2 < (unsigned char)result) {
		result = v2 + 4 * a2 + 1132;
		do {
			++v8;
			*(uint32_t*)result = *(uint32_t*)(result + 4);
			result += 4;
		} while (v8 < *(unsigned char*)(v2 + 1129));
	}
	return result;
}

//----- (00528610) --------------------------------------------------------
void sub_528610(int a1) {
	int v1;    // ebx
	int v2;    // esi
	int v3;    // ebp
	char v4;   // al
	int* i;    // edi
	double v6; // st7
	double v7; // st6
	double v8; // st5
	float v9;  // [esp+14h] [ebp+4h]

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 748);
	v3 = 0;
	v9 = 100000000.0;
	v4 = *(uint8_t*)(v2 + 1129);
	*(uint32_t*)(v2 + 1196) = 0;
	if (v4) {
		for (i = (int*)(v2 + 1132); *i != *(uint32_t*)(v2 + 1216); ++i) {
			if (nox_xxx_unitIsEnemyTo_5330C0(v1, *i) && nox_xxx_checkIsKillable_528190(*i)) {
				v6 = *(float*)(*i + 56) - *(float*)(v1 + 56);
				v7 = *(float*)(*i + 60) - *(float*)(v1 + 60);
				v8 = v7 * v7 + v6 * v6;
				if (v8 < v9) {
					v9 = v8;
					*(uint32_t*)(v2 + 1196) = *i;
				}
			}
			if (++v3 >= *(unsigned char*)(v2 + 1129)) {
				return;
			}
		}
		*(uint32_t*)(v2 + 1196) = *(uint32_t*)(v2 + 1216);
	}
}

//----- (005286D0) --------------------------------------------------------
void nox_xxx_monsterUpdateSeenEnemies_5286D0(int a1, int a2) {
	int v2;    // esi
	int v3;    // ebx
	int v4;    // eax
	int v5;    // eax
	double v6; // st7
	double v7; // st6
	float* v8; // eax
	float v9;  // [esp+10h] [ebp-4h]
	float v10; // [esp+1Ch] [ebp+8h]

	v2 = a2;
	v3 = *(uint32_t*)(a2 + 748);
	if (a2 != a1) {
		if (*(uint8_t*)(a1 + 8) & 6) {
			if (!(*(uint32_t*)(a1 + 16) & 0x8020)) {
				v4 = *(uint32_t*)(v3 + 1440);
				if ((v4 & 0x400 || nox_xxx_unitIsEnemyTo_5330C0(a2, a1)) && !sub_528950(a2, a1)) {
					v5 = *(uint32_t*)(v3 + 1440);
					if (v5 & 0x100 ||
						(v6 = *(float*)(a1 + 56) - *(float*)(a2 + 56), v7 = *(float*)(a1 + 60) - *(float*)(a2 + 60),
						 v9 = v7, v8 = getMemFloatPtr(0x587000, 194136 + 8 * *(short*)(a2 + 124)),
						 v10 = sqrt(v7 * v9 + v6 * v6) + 0.001, v9 / v10 * v8[1] + v6 / v10 * *v8 >= 0.5)) {
						if (nox_xxx_unitCanInteractWith_5370E0(v2, a1, 0)) {
							nox_xxx_monsterVisionSeeEnemy_5287B0(v2, a1);
						}
					}
				}
			}
		}
	}
}

//----- (005287B0) --------------------------------------------------------
void nox_xxx_monsterVisionSeeEnemy_5287B0(int a1, int a2) {
	int v2;     // esi
	double v3;  // st7
	int v4;     // edi
	int v5;     // ebx
	int v6;     // ebp
	int v7;     // edx
	double v8;  // st6
	double v9;  // st5
	double v10; // st4
	int v11;    // ebx
	double v12;
	double v13;
	unsigned int v15; // eax
	int v16;          // eax
	int v17;          // eax
	float v18;        // [esp+14h] [ebp+4h]

	v2 = a1;
	v3 = 0.0;
	v4 = *(uint32_t*)(a1 + 748);
	v5 = 16;
	v6 = 0;
	if (*(uint8_t*)(v4 + 1129) == 16) {
		v7 = v4 + 1132;
		do {
			v8 = *(float*)(*(uint32_t*)v7 + 56) - *(float*)(v2 + 56);
			v9 = *(float*)(*(uint32_t*)v7 + 60) - *(float*)(v2 + 60);
			v10 = v9 * v9 + v8 * v8;
			if (v10 > v3) {
				v18 = v10;
				v3 = v18;
				v6 = *(uint32_t*)v7;
			}
			v7 += 4;
			--v5;
		} while (v5);
		v11 = a2;
		v12 = *(float*)(a2 + 56) - *(float*)(v2 + 56);
		v13 = *(float*)(a2 + 60) - *(float*)(v2 + 60);
		if (v3 <= v13 * v13 + v12 * v12) {
			return;
		}
		sub_528910(v2, v6);
	} else {
		v11 = a2;
	}
	*(uint32_t*)(v4 + 4 * *(unsigned char*)(v4 + 1129) + 1132) = v11;
	v15 = *(uint32_t*)(v4 + 536);
	++*(uint8_t*)(v4 + 1129);
	if (gameFrame() > v15) {
		if (nox_xxx_unitIsEnemyTo_5330C0(v2, v11)) {
			if (!nox_xxx_unitIsZombie_534A40(v2) || (v16 = *(uint32_t*)(v2 + 16), (v16 & 0x8000) == 0)) {
				v17 = nox_xxx_monsterGetSoundSet_424300(v2);
				if (v17) {
					nox_xxx_aud_501960(*(uint32_t*)(v17 + 68), v2, 0, 0);
				}
				*(uint32_t*)(v4 + 536) =
					gameFrame() + nox_common_randomInt_415FA0(2 * gameFPS(), 4 * gameFPS());
			}
		}
	}
	nox_xxx_scriptCallByEventBlock_502490((int*)(v4 + 1232), v11, v2, 14);
}

//----- (00528910) --------------------------------------------------------
int sub_528910(int a1, int a2) {
	int result;  // eax
	int v3;      // edx
	int v4;      // ecx
	uint32_t* i; // edx

	result = 0;
	v3 = *(uint32_t*)(a1 + 748);
	v4 = *(unsigned char*)(v3 + 1129);
	if (v4 > 0) {
		for (i = (uint32_t*)(v3 + 1132); *i != a2; ++i) {
			if (++result >= v4) {
				return result;
			}
		}
		result = nox_xxx_aiLostSight_528560(a1, result);
	}
	return result;
}

//----- (00528950) --------------------------------------------------------
int sub_528950(int a1, int a2) {
	int v2;      // ecx
	int v3;      // eax
	int v4;      // edx
	uint32_t* i; // ecx

	v2 = *(uint32_t*)(a1 + 748);
	v3 = 0;
	v4 = *(unsigned char*)(v2 + 1129);
	if (v4 <= 0) {
		return 0;
	}
	for (i = (uint32_t*)(v2 + 1132); *i != a2; ++i) {
		if (++v3 >= v4) {
			return 0;
		}
	}
	return 1;
}

//----- (00528990) --------------------------------------------------------
nox_object_t* nox_xxx_getFirstUpdatableObject_4DA8A0();
nox_object_t* nox_xxx_getNextUpdatableObject_4DA8B0(nox_object_t* obj);
int sub_528990(nox_object_t* a1) {
	int result; // eax
	int i;      // esi

	result = nox_xxx_getFirstUpdatableObject_4DA8A0();
	for (i = result; result; i = result) {
		if (*(uint8_t*)(i + 8) & 2) {
			if (!(*(uint8_t*)(i + 16) & 0x20)) {
				sub_528910(i, a1);
				sub_528610(i);
			}
		}
		result = nox_xxx_getNextUpdatableObject_4DA8B0(i);
	}
	return result;
}

//----- (005289D0) --------------------------------------------------------
void nox_xxx_netReportDestroyObject_5289D0(nox_object_t* a1p) {
	int a1 = a1p;
	char* result; // eax
	int i;        // edi
	int v4;       // [esp+0h] [ebp-4h]

	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		if ((1 << *(uint8_t*)(i + 2064)) & *(uint32_t*)(a1 + 148)) {
			LOBYTE(v4) = ((unsigned char)*(uint32_t*)(a1 + 20) >> 6) | 0x31;
			*(uint16_t*)((char*)&v4 + 1) = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)a1);
			nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(i + 2064), &v4, 3, 0, 1);
		}
		if (*(uint8_t*)(a1 + 8) & 6) {
			nox_xxx_netFriendAddRemove_4D97A0(*(unsigned char*)(i + 2064), (uint32_t*)a1, 0);
		}
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
}

//----- (00528A60) --------------------------------------------------------
int nox_xxx_netObjectOutOfSight_528A60(int a1, uint32_t* a2) {
	char v4[3]; // [esp+0h] [ebp-4h]
	v4[0] = 50;
	*(uint16_t*)&v4[1] = nox_xxx_netGetUnitCodeServ_578AC0(a2);
	return nox_xxx_netSendPacket0_4E5420(a1, v4, 3, 0, 1);
}

//----- (00528A90) --------------------------------------------------------
int nox_xxx_netObjectInShadows_528A90(int a1, uint32_t* a2) {
	char v4[3]; // [esp+0h] [ebp-4h]
	v4[0] = 51;
	*(uint16_t*)&v4[1] = nox_xxx_netGetUnitCodeServ_578AC0(a2);
	return nox_xxx_netSendPacket0_4E5420(a1, v4, 3, 0, 1);
}

//----- (00528BD0) --------------------------------------------------------
int nox_xxx_monsterCmdSend_528BD0(int unit, int source, const char* command, short a4) {
	short v4;      // ax
	double v5;     // st7
	long long v6;  // rax
	double v7;     // st7
	int result;    // eax
	int i;         // esi
	char v10[520]; // [esp+Ch] [ebp-208h]

	v10[0] = -88; // MSG_TEXT_MESSAGE
	v10[3] = 8;
	v4 = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)unit);
	v5 = *(float*)(unit + 56);
	*(uint16_t*)&v10[1] = v4;
	v6 = (long long)v5;
	v7 = *(float*)(unit + 60);
	*(uint16_t*)&v10[4] = v6;
	*(uint16_t*)&v10[6] = (long long)v7;
	*(uint16_t*)&v10[9] = a4;
	v10[8] = strlen(command) + 1;
	result = source;
	strcpy(&v10[11], command);
	if (source) {
		if (*(uint8_t*)(source + 8) & 4) { // if source is player / local player ?
			result = nox_netlist_addToMsgListCli_40EBC0(
				*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(source + 748) + 276) + 2064), 1, v10,
				(unsigned char)v10[8] + 11);
		}
	} else {
		result = nox_xxx_getFirstPlayerUnit_4DA7C0();
		for (i = result; result; i = result) {
			nox_netlist_addToMsgListCli_40EBC0(*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064), 1,
											   v10, (unsigned char)v10[8] + 11);
			result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
		}
	}
	return result;
}

//----- (00528D60) --------------------------------------------------------
int nox_xxx_destroyEveryChatMB_528D60() {
	int result; // eax
	int i;      // esi
	char v2[3]; // [esp+4h] [ebp-4h]

	v2[0] = -54;
	*(uint16_t*)&v2[1] = -8531;
	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		nox_netlist_addToMsgListCli_40EBC0(*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064), 1, v2,
										   3);
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (00528DB0) --------------------------------------------------------
int nox_xxx_XFerMonster_528DB0(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;            // esi
	int result;        // eax
	char* v3;          // edi
	char* v4;          // eax
	char* v5;          // eax
	char* v6;          // eax
	char* v7;          // eax
	char* v8;          // eax
	char* v9;          // eax
	char* v10;         // eax
	char* v11;         // eax
	char* v12;         // eax
	char* v13;         // edi
	int v14;           // eax
	short v15;         // ax
	int v16;           // edx
	uint32_t* v17;     // edi
	uint32_t* v18;     // eax
	int v19;           // ecx
	int v20;           // ebx
	char* v21;         // ebp
	int i;             // edi
	int v23;           // eax
	int* v24;          // ebp
	char* v25;         // ebx
	char* v26;         // ebx
	int v27;           // ebp
	unsigned char v28; // bl
	uint8_t* v29;      // edi
	int v30;           // eax
	int v31;           // eax
	int v32;           // edi
	int v33;           // ebx
	int j;             // eax
	int v35;           // esi
	int v36;           // ecx
	uint32_t* v37;     // eax
	uint32_t* v38;     // eax
	int v39;           // edi
	int* v40;          // ecx
	uint32_t* v41;     // edx
	int v42;           // esi
	int v43;           // eax
	int v44;           // [esp+10h] [ebp-128h]
	int v45;           // [esp+14h] [ebp-124h]
	int v46;           // [esp+18h] [ebp-120h]
	uint8_t* v47;      // [esp+1Ch] [ebp-11Ch]
	int v48;           // [esp+20h] [ebp-118h]
	int v49;           // [esp+24h] [ebp-114h]
	int v50;           // [esp+28h] [ebp-110h]
	int v51;           // [esp+2Ch] [ebp-10Ch]
	uint32_t v52[2];   // [esp+30h] [ebp-108h]
	char v53[256];     // [esp+38h] [ebp-100h]

	v1 = *(uint32_t*)(a1 + 748);
	v51 = *(uint32_t*)(a1 + 136);
	if (!*getMemU32Ptr(0x5D4594, 2487692)) {
		*getMemU32Ptr(0x5D4594, 2487692) = nox_xxx_getNameId_4E3AA0("Glyph");
	}
	v45 = 64;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 2u);
	if ((short)v45 > 64) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530((int*)a1, (short)v45);
	if (!result) {
		return result;
	}
	if (!nox_crypt_IsReadOnly()) {
		nox_xxx_xferIndexedDirection_509E20(*(short*)(a1 + 124), (int2*)v52);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(v52, 8u);
	if ((short)v45 >= 3) {
		v3 = *(char**)(a1 + 756);
		if (v3) {
			v4 = v3 + 640;
		} else {
			v4 = 0;
		}
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1232, v4);
		if (v3) {
			v5 = v3 + 896;
		} else {
			v5 = 0;
		}
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1264, v5);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1220), 2u);
		if (v3) {
			v6 = v3 + 768;
		} else {
			v6 = 0;
		}
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1224, v6);
		if ((short)v45 >= 31) {
			v7 = v3 ? v3 + 1024 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1240, v7);
			v8 = v3 ? v3 + 1152 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1248, v8);
			v9 = v3 ? v3 + 1280 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1256, v9);
			v10 = v3 ? v3 + 1408 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1272, v10);
			v11 = v3 ? v3 + 1536 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1280, v11);
			v12 = v3 ? v3 + 1664 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1288, v12);
			if ((short)v45 >= 52) {
				if (v3) {
					v13 = v3 + 1792;
				} else {
					v13 = 0;
				}
				nox_xxx_xferReadScriptHandler_4F5580(v1 + 1296, v13);
			}
		}
	} else {
		sub_4F5540(v1 + 1232);
		sub_4F5540(v1 + 1264);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1220), 2u);
		sub_4F5540(v1 + 1224);
	}
	v14 = nox_crypt_IsReadOnly();
	if (nox_crypt_IsReadOnly() != 1 ||
		(v15 = nox_xxx_xferDirectionToAngle_509E00(v52), *(uint16_t*)(a1 + 126) = v15, *(uint16_t*)(a1 + 124) = v15,
		 v14 = nox_crypt_IsReadOnly(), nox_crypt_IsReadOnly() != 1)) {
		if (!v14) {
			v46 = 0;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 4u);
		}
	} else if ((short)v45 >= 11) {
		v46 = 0;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 4u);
	}
	if ((short)v45 >= 31) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1332), 1u);
		if ((short)v45 < 51) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 2u);
			*(uint32_t*)(v1 + 1440) = (unsigned short)v46;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1440), 4u);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1352), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1336), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1344), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1312), 4u);
		if ((short)v45 < 33) {
			nox_xxx_cryptSeekCur_40E0A0(2);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1304), 4u);
		*(uint32_t*)(v1 + 1308) = *(uint32_t*)(v1 + 1304);
		if ((short)v45 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1360), 4u);
		}
		LOBYTE(v48) = strlen((const char*)(v1 + 1364));
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v48, 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1364), (unsigned char)v48);
		*(uint8_t*)((unsigned char)v48 + v1 + 1364) = 0;
		if ((short)v45 >= 34) {
			if (nox_crypt_IsReadOnly()) {
				memset((void*)(v1 + 1488), 0, 0x224u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 4u);
				for (i = 0; i < v44; ++i) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v48, 1u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v53, (unsigned char)v48);
					v53[(unsigned char)v48] = 0;
					v23 = nox_xxx_spellNameToN_4243F0(v53);
					nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 4 * v23 + 1488), 4u);
				}
			} else {
				v16 = 0;
				v17 = (uint32_t*)(v1 + 1488);
				v44 = 0;
				v18 = (uint32_t*)(v1 + 1488);
				v19 = 137;
				do {
					if (*v18) {
						++v16;
					}
					++v18;
					--v19;
				} while (v19);
				v44 = v16;
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 4u);
				v20 = 0;
				v47 = (uint8_t*)(v1 + 1488);
				do {
					if (*v17) {
						v21 = nox_xxx_spellNameByN_424870(v20);
						LOBYTE(v46) = strlen(v21);
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 1u);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v21, (unsigned char)v46);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v47, 4u);
					}
					++v20;
					v17 = v47 + 4;
					v47 += 4;
				} while (v20 < 137);
			}
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1488), 0x224u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1448) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1450) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1448), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1450), 2u);
		}
		if ((short)v45 <= 32) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1456) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1458) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1456), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1458), 2u);
		}
		if ((short)v45 <= 32) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1464) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1466) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1464), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1466), 2u);
		}
		if ((short)v45 <= 32) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1472) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1474) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1472), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1474), 2u);
		}
		if ((short)v45 <= 32) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1480) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1482) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1480), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1482), 2u);
		}
		if ((short)v45 > 32 || (nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u), (short)v45 >= 32)) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1316), 4u);
		}
		if ((short)v45 >= 33) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2040), 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1320), 4u);
			if ((short)v45 < 42) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 2u);
				if (!(uint16_t)v44) {
					*(uint8_t*)(v1 + 1445) = 1;
				}
			}
			if ((short)v45 < 53) {
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2044), 4u);
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2048), 4u);
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2052), 4u);
			} else {
				v24 = (int*)(v1 + 2044);
				v46 = 3;
				do {
					if (nox_crypt_IsReadOnly() == 1) {
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v53, (unsigned char)v44);
						v53[(unsigned char)v44] = 0;
						*v24 = nox_xxx_spellNameToN_4243F0(v53);
					} else {
						v25 = nox_xxx_spellNameByN_424870(*v24);
						LOBYTE(v44) = strlen(v25);
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v25, (unsigned char)v44);
					}
					++v24;
					--v46;
				} while (v46);
			}
		}
		if ((short)v45 >= 34) {
			if (nox_crypt_IsReadOnly()) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(v53, (unsigned char)v44);
				v53[(unsigned char)v44] = 0;
				*(uint32_t*)(v1 + 1360) = nox_xxx_actionNByNameMB_5345F0(v53);
			} else {
				v26 = sub_5345B0(*(uint32_t*)(v1 + 1360));
				LOBYTE(v44) = strlen(v26);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(v26, (unsigned char)v44);
			}
		}
	}
	if ((short)v45 >= 41) {
		result = nox_xxx_XFer_ActionData_529CE0(a1);
		if (!result) {
			return result;
		}
	}
	if ((short)v45 >= 42) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1445), 1u);
	}
	if ((short)v45 >= 43 && *(uint8_t*)(a1 + 12) & 8) {
		LOBYTE(v44) = 0;
		v27 = *(uint32_t*)(a1 + 692);
		if ((short)v45 >= 50) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v27 + 1716), 4u);
		}
		if ((short)v45 >= 61) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v27 + 1720), 4u);
		}
		if ((short)v45 >= 48) {
			LOBYTE(v47) = strlen((const char*)(v27 + 1684));
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v27 + 1684), (unsigned char)v47);
			*(uint8_t*)((unsigned char)v47 + v27 + 1684) = 0;
		}
		if (nox_crypt_IsReadOnly() == 1) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			if (!v27) {
				goto LABEL_137;
			}
			*(uint8_t*)v27 = v44;
		} else {
			if (v27) {
				LOBYTE(v44) = *(uint8_t*)v27;
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
		}
		if (v27) {
			LOBYTE(v44) = 0;
			if (*(uint8_t*)v27) {
				do {
					if (nox_crypt_IsReadOnly() == 1) {
						nox_xxx_XFer_ReadShopItem_52A840(v27 + 28 * (unsigned char)v44 + 4, (short)v45);
					} else {
						nox_xxx_XFer_WriteShopItem_52A5F0(v27 + 28 * (unsigned char)v44 + 4);
					}
					LOBYTE(v44) = v44 + 1;
				} while ((unsigned char)v44 < *(uint8_t*)v27);
			}
		}
	}
LABEL_137:
	if ((short)v45 >= 44) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v1, 4u);
	}
	if ((short)v45 >= 45) {
		v50 = *(uint32_t*)(a1 + 12) & 0x180;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v50, 4u);
		*(uint32_t*)(a1 + 12) |= v50;
	}
	if ((short)v45 >= 49) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(*(uint8_t**)(a1 + 556), 2u);
	}
	if ((short)v45 >= 51) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1348), 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1340), 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1444), 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2036), 1u);
	}
	if (*(uint8_t*)(a1 + 12) & 0x20 && (short)v45 >= 54) {
		v28 = 0;
		v29 = (uint8_t*)(v1 + 2076);
		LOBYTE(v46) = 0;
		do {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v29, 3u);
			if (nox_crypt_IsReadOnly() == 1) {
				nox_xxx_setNPCColor_4E4A90(a1, v46, (int)v29);
			}
			++v28;
			v29 += 3;
			LOBYTE(v46) = v28;
		} while (v28 < 6u);
	}
	if ((short)v45 >= 55 && *(uint8_t*)(a1 + 12) & 0x20) {
		nox_xxx_readNPCVoiceSet_52AD10(a1);
	}
	if (!((short)v45 < 62 || (result = nox_xxx_XFer_ReadMonsterBuffs_52AAB0((uint32_t*)a1)) != 0)) {
		return result;
	}
	if ((short)v45 >= 63 && *(uint32_t*)(a1 + 12) & 0x80000) {
		nox_xxx_readNPCVoiceSet_52AD10(a1);
	}
	if ((short)v45 >= 64) {
		LOBYTE(v47) = *(uint8_t*)(a1 + 540);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 1u);
		if (nox_crypt_IsReadOnly() != 1) {
			goto LABEL_171;
		}
		if (!(uint8_t)v47) {
			goto LABEL_164;
		}
		nox_xxx_setSomePoisonData_4EEA90(a1, (unsigned char)v47);
	}
	if (nox_crypt_IsReadOnly() != 1) {
		goto LABEL_171;
	}
LABEL_164:
	if (!*(uint8_t*)(v1 + 1445)) {
		v30 = *(uint32_t*)(a1 + 16);
		if ((v30 & 0x8000) != 0 && !nox_xxx_unitIsZombie_534A40(a1)) {
			*(uint32_t*)(a1 + 16) |= 0x40u;
		}
		goto LABEL_171;
	}
	if (nox_common_gameFlags_check_40A5C0(1)) {
		*(uint16_t*)(*(uint32_t*)(a1 + 556) + 4) = 0;
		**(uint16_t**)(a1 + 556) = 0;
	}
	if (nox_crypt_IsReadOnly() == 1) {
		v30 = *(uint32_t*)(a1 + 16);
		if ((v30 & 0x8000) != 0 && !nox_xxx_unitIsZombie_534A40(a1)) {
			*(uint32_t*)(a1 + 16) |= 0x40u;
		}
	}
LABEL_171:
	if (*(uint32_t*)(a1 + 136)) {
		if (nox_crypt_IsReadOnly() != 1) {
			*(uint32_t*)(a1 + 136) = v51;
			return 1;
		}
		result = nox_xxx_xfer_4F3E30(v45, a1, *(uint32_t*)(a1 + 136));
		if (!result) {
			return result;
		}
	}
	if (nox_crypt_IsReadOnly() == 1) {
		if (nox_common_gameFlags_check_40A5C0(0x200000) || !nox_xxx_gameIsSwitchToSolo_4DB240()) {
			nox_xxx_monsterOnSpawnSpellcaster_529BC0(a1);
		}
		if (nox_crypt_IsReadOnly() == 1) {
			if (*(uint8_t*)(a1 + 8) & 2) {
				v31 = *(uint32_t*)(a1 + 12);
				if (v31 & 0x2000) {
					v32 = 0;
					v33 = 0;
					for (j = nox_xxx_inventoryGetFirst_4E7980(a1); j;
						 j = nox_xxx_inventoryGetNext_4E7990(j)) {
						if (*(unsigned short*)(j + 4) == *getMemU32Ptr(0x5D4594, 2487692)) {
							v32 = 1;
						}
					}
					v35 = v1 + 2044;
					v36 = 3;
					v37 = (uint32_t*)v35;
					do {
						if (*v37) {
							++v33;
						}
						++v37;
						--v36;
					} while (v36);
					if (!v32 && v33) {
						v38 = nox_xxx_newObjectByTypeID_4E3810("Glyph");
						v46 = (int)v38;
						if (v38) {
							v39 = v38[173];
							if (v33 > 0) {
								v40 = (int*)v35;
								v41 = (uint32_t*)v38[173];
								v42 = v33;
								do {
									v43 = *v40;
									++v40;
									*v41 = v43;
									++v41;
									--v42;
								} while (v42);
								v38 = (uint32_t*)v46;
							}
							*(uint8_t*)(v39 + 20) = v33;
							*(uint32_t*)(v39 + 24) = 0;
							*(uint32_t*)(v39 + 28) = *(uint32_t*)(a1 + 56);
							*(uint32_t*)(v39 + 32) = *(uint32_t*)(a1 + 60);
						}
						nox_xxx_inventoryPutImpl_4F3070(a1, (int)v38, 1);
					}
				}
			}
		}
	}
	*(uint32_t*)(a1 + 136) = v51;
	return 1;
}
// 528DB0: using guessed type char var_100[256];

//----- (00529BC0) --------------------------------------------------------
void nox_xxx_monsterOnSpawnSpellcaster_529BC0(int a1) {
	int v1;       // esi
	uint32_t* v2; // eax
	int i;        // ecx
	int v4;       // edx
	int v5;       // edi
	int v6;       // edi
	int v7;       // ecx

	if (a1) {
		v1 = *(uint32_t*)(a1 + 748);
		v2 = nox_xxx_monsterDefByTT_517560(*(unsigned short*)(a1 + 4));
		if (v2) {
			if (!*(uint8_t*)(v1 + 1445)) {
				**(uint16_t**)(a1 + 556) = *((uint16_t*)v2 + 34);
				*(uint16_t*)(*(uint32_t*)(a1 + 556) + 4) = *((uint16_t*)v2 + 34);
				*(uint16_t*)(*(uint32_t*)(a1 + 556) + 2) = *((uint16_t*)v2 + 34);
			}
			if (*(uint8_t*)(v1 + 1444) == 1) {
				*(uint32_t*)(v1 + 1440) = v2[23];
			} else {
				for (i = 0; i < 22; ++i) {
					v4 = 1 << i;
					if (!((1 << i) & 0x19C40)) {
						v5 = *(uint32_t*)(v1 + 1440);
						if (v5 & v4 && !(v4 & v2[23])) {
							*(uint32_t*)(v1 + 1440) = v5 & ~v4;
						}
						v6 = *(uint32_t*)(v1 + 1440);
						if (!(v6 & v4) && v4 & v2[23]) {
							*(uint32_t*)(v1 + 1440) = v4 | v6;
						}
					}
				}
			}
			v7 = *(uint32_t*)(v1 + 1440);
			if (!(v7 & 0x20)) {
				BYTE1(v7) &= 0xE7u;
				*(uint32_t*)(v1 + 1440) = v7;
			}
			if (*(uint8_t*)(v1 + 1340) == 1) {
				*(uint32_t*)(v1 + 1336) = v2[20];
			}
			if (*(uint8_t*)(v1 + 1348) == 1) {
				*(uint32_t*)(v1 + 1344) = v2[21];
			}
			if (*(uint8_t*)(v1 + 2036) == 1) {
				nox_xxx_monsterAutoSpells_54C0C0(a1);
			}
		}
	}
}

//----- (00529CE0) --------------------------------------------------------
nox_waypoint_t* sub_579C80(unsigned int a1);
int nox_xxx_XFer_ActionData_529CE0(int a1) {
	int v1;             // ebp
	int v3;             // ebx
	bool v4;            // cc
	uint8_t** v5;       // esi
	int v6;             // esi
	bool v7;            // zf
	uint8_t* v8;        // esi
	int v9;             // ecx
	int v10;            // eax
	int v11;            // eax
	int v12;            // eax
	int v13;            // eax
	int v14;            // ebx
	unsigned char* v15; // edi
	int* v16;           // ebp
	int v17;            // eax
	int v18;            // [esp+4h] [ebp-118h]
	char v19;           // [esp+Bh] [ebp-111h]
	int v20;            // [esp+Ch] [ebp-110h]
	int v21;            // [esp+10h] [ebp-10Ch]
	int v22;            // [esp+14h] [ebp-108h]
	int v23;            // [esp+18h] [ebp-104h]
	char v24[256];      // [esp+1Ch] [ebp-100h]

	v1 = *(uint32_t*)(a1 + 748);
	v21 = 4;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v21, 2u);
	if ((short)v21 > 4) {
		return 0;
	}
	v19 = 1;
	if ((short)v21 < 2) {
		goto LABEL_67;
	}
	v19 = 0;
	if (nox_common_gameFlags_check_40A5C0(1) && !nox_common_gameFlags_check_40A5C0(0x400000)) {
		v19 = 1;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v19, 1u);
	if (!(v19 || (uint16_t)v21 == 1)) {
		return 1;
	}
LABEL_67:
	v23 = gameFrame();
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v23, 4u);
	v3 = gameFrame() - v23;
	v18 = 0;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 8), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 12), 8 * *(uint32_t*)(v1 + 8));
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 268), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 272), 8u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 280), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 284), 1u);
	if (nox_crypt_IsReadOnly() == 1) {
		nox_xxx_AssignIfGreater_52A420((int*)(v1 + 280), v3);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 296), 4u);
	v4 = *(int*)(v1 + 296) <= 0;
	v18 = 0;
	if (!v4) {
		v5 = (uint8_t**)(v1 + 300);
		do {
			if (nox_crypt_IsReadOnly()) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v22, 4u);
				*v5 = sub_579C80(v22);
			} else {
				nox_xxx_fileReadWrite_426AC0_file3_fread(*v5, 4u);
			}
			++v5;
			v4 = ++v18 < *(int*)(v1 + 296);
		} while (v4);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 364), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 368), 8u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 376), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 380), 8u);
	strcpy(v24, nox_xxx_getSndName_40AF80(*(uint32_t*)(v1 + 388)));
	LOBYTE(v20) = strlen(v24);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v20, 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(v24, (unsigned char)v20);
	v24[(unsigned char)v20] = 0;
	*(uint32_t*)(v1 + 388) = nox_xxx_utilFindSound_40AF50(v24);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 396), 8u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 404), 4u);
	if (nox_crypt_IsReadOnly() == 1) {
		nox_xxx_AssignIfGreater_52A420((int*)(v1 + 404), v3);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 481), 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 482), 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 483), 1u);
	if ((short)v21 < 3) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 496), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 500), 8u);
	if (nox_crypt_IsReadOnly() == 1) {
		nox_xxx_AssignIfGreater_52A420((int*)(v1 + 496), v3);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 536), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 540), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 536), v3);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 540), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 544), 1u);
	v6 = 0;
	if ((signed char)*(uint8_t*)(v1 + 544) >= 0) {
		v18 = v1 + 552;
		do {
			sub_52A440(a1, v18, v3);
			++v6;
			v18 += 24;
		} while (v6 <= *(char*)(v1 + 544));
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1129), 1u);
	v7 = *(uint8_t*)(v1 + 1129) == 0;
	v18 = 0;
	if (!v7) {
		v8 = (uint8_t*)(v1 + 1132);
		do {
			if (nox_crypt_IsReadOnly()) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(v8, 4u);
			} else {
				if (!*getMemU32Ptr(0x5D4594, 2487688)) {
					*getMemU32Ptr(0x5D4594, 2487688) = nox_xxx_getNameId_4E3AA0("NewPlayer");
				}
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(*(uint32_t*)v8 + 44), 4u);
			}
			v9 = *(unsigned char*)(v1 + 1129);
			v8 += 4;
			++v18;
		} while (v18 < v9);
	}
	if (nox_crypt_IsReadOnly()) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1196), 4u);
	} else {
		v10 = *(uint32_t*)(v1 + 1196);
		if (v10) {
			v22 = *(uint32_t*)(v10 + 44);
		} else {
			v22 = 0;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v22, 4u);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1204), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 1204), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2096), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2100), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2104), 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2105), 1u);
	LOBYTE(v20) = strlen((const char*)(v1 + 2106));
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v20, 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2106), (unsigned char)v20);
	*(uint8_t*)((unsigned char)v20 + v1 + 2106) = 0;
	if ((short)v21 < 4) {
		return 1;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 4), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 288), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 292), 4u);
	v11 = nox_server_getObjectFromNetCode_4ECCB0(*(uint32_t*)(v1 + 392));
	if (v11) {
		v18 = *(uint32_t*)(v11 + 44);
	} else {
		v18 = 0;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	if (nox_crypt_IsReadOnly() == 1) {
		*(uint32_t*)(v1 + 392) = v18;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 492), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 508), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 508), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 512), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 512), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 516), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 516), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 520), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 520), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 528), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 528), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 532), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 532), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 524), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 548), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 548), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1128), 1u);
	v12 = nox_server_getObjectFromNetCode_4ECCB0(*(uint32_t*)(v1 + 1200));
	if (v12) {
		v18 = *(uint32_t*)(v12 + 44);
	} else {
		v18 = 0;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	if (nox_crypt_IsReadOnly() == 1) {
		*(uint32_t*)(v1 + 1200) = v18;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1208), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 1208), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1212), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 1212), v3);
	v13 = *(uint32_t*)(v1 + 1216);
	v14 = 0;
	if (v13) {
		v18 = *(uint32_t*)(v13 + 44);
	} else {
		v18 = 0;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	if (nox_crypt_IsReadOnly() == 1) {
		*(uint32_t*)(v1 + 1216) = v18;
	}
	v15 = (unsigned char*)(v1 + 2172);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2172), 1u);
	if (*(uint8_t*)(v1 + 2172)) {
		v16 = (int*)(v1 + 2140);
		do {
			v17 = nox_server_getObjectFromNetCode_4ECCB0(*v16);
			if (v17) {
				v18 = *(uint32_t*)(v17 + 44);
			} else {
				v18 = 0;
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
			if (nox_crypt_IsReadOnly() == 1) {
				*v16 = v18;
			}
			++v14;
			++v16;
		} while (v14 < *v15);
	}
	return 1;
}
// 529CE0: using guessed type char var_100[256];

//----- (0052A420) --------------------------------------------------------
int nox_xxx_AssignIfGreater_52A420(int* a1, int a2) {
	int result; // eax

	result = a2 + *a1;
	if (result >= 1) {
		*a1 = result;
	} else {
		*a1 = 1;
	}
	return result;
}

//----- (0052A440) --------------------------------------------------------
int sub_52A440(int a1, int a2, int a3) {
	int v3;        // eax
	int v4;        // edi
	int* v5;       // esi
	int result;    // eax
	size_t v7;     // [esp-4h] [ebp-120h]
	int v8;        // [esp+10h] [ebp-10Ch]
	int v9;        // [esp+14h] [ebp-108h]
	int v10;       // [esp+18h] [ebp-104h]
	char v11[256]; // [esp+1Ch] [ebp-100h]

	strcpy(v11, sub_534650(*(uint32_t*)a2));
	LOBYTE(v10) = strlen(v11);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v10, 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(v11, (unsigned char)v10);
	v11[(unsigned char)v10] = 0;
	v3 = nox_xxx_actionByName_534670(v11);
	*(uint32_t*)a2 = v3;
	LOBYTE(v8) = getMemByte(0x587000, 255604 + 16 * v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v8, 1u);
	v4 = 0;
	if ((unsigned char)v8 > 0u) {
		v5 = (int*)(a2 + 4);
		while (1) {
			result = *getMemU32Ptr(0x587000, 255608 + 4 * (v4 + 4 * *(uint32_t*)a2));
			switch (result) {
			case 0:
				v7 = 8;
				nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
				break;
			case 1:
				if (nox_crypt_IsReadOnly()) {
					v7 = 4;
					nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
					break;
				}
				if (*v5) {
					nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(*v5 + 44), 4u);
				} else {
					v9 = 0;
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v9, 4u);
				}
				break;
			case 2:
				if (nox_crypt_IsReadOnly()) {
					v7 = 4;
					nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
					break;
				}
				if (*v5) {
					nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)*v5, 4u);
				} else {
					v9 = 0;
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v9, 4u);
				}
				break;
			case 3:
			case 4:
			case 6:
				v7 = 4;
				nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
				break;
			case 5:
				nox_xxx_fileReadWrite_426AC0_file3_fread(v5, 4u);
				if (nox_crypt_IsReadOnly() == 1) {
					nox_xxx_AssignIfGreater_52A420(v5, a3);
				}
				break;
			case 7:
				v7 = 1;
				nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
				break;
			default:
				return result;
			}
			++v4;
			v5 += 2;
			if (v4 >= (unsigned char)v8) {
				return nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(a2 + 20), 4u);
			}
		}
	}
	return nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(a2 + 20), 4u);
}
// 52A440: using guessed type char var_100[256];

//----- (0052AAB0) --------------------------------------------------------
int nox_xxx_XFer_ReadMonsterBuffs_52AAB0(uint32_t* a1) {
	int v1;   // ebp
	char* v2; // ebx
	int v3;   // eax
	int v5;   // ebx
	int v6;   // edi
	int v7;   // ecx
	int v8;   // eax
	int v9;   // eax
	int v10;  // [esp-4h] [ebp-138h]
	int v11;  // [esp+10h] [ebp-124h]
	int v12;  // [esp+14h] [ebp-120h]
	int v13;  // [esp+18h] [ebp-11Ch]
	int v14;  // [esp+1Ch] [ebp-118h]
	int v15;  // [esp+20h] [ebp-114h]
	int v16;  // [esp+24h] [ebp-110h]
	int v17[3];
	char v20[256]; // [esp+34h] [ebp-100h]

	v13 = 2;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v13, 2u);
	if ((short)v13 > 2 || (short)v13 <= 0) {
		return 0;
	}
	LOBYTE(v16) = sub_424CB0((int)a1);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v16, 1u);
	if (nox_crypt_IsReadOnly()) {
		v5 = 0;
		if (!(uint8_t)v16) {
			return 1;
		}
		while (1) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v11, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v20, (unsigned char)v11);
			v20[(unsigned char)v11] = 0;
			v6 = nox_xxx_enchantByName_424880(v20);
			if (v6 == -1) {
				break;
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v15, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v14, 4u);
			v7 = a1[15];
			v17[1] = a1[14];
			v10 = (unsigned char)v15;
			v17[0] = a1;
			v17[2] = v7;
			v8 = nox_xxx_getEnchantSpell_424920(v6);
			nox_xxx_spellAccept_4FD400(v8, (int)a1, a1, (int)a1, v17, v10);
			*((uint16_t*)a1 + v6 + 172) = v14;
			if (v6 == 26 && (short)v13 >= 2) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v12, 4u);
				v9 = sub_4FF2D0(51, (int)a1);
				if (v9) {
					*(uint32_t*)(v9 + 72) = v12;
				}
			}
			if (++v5 >= (unsigned char)v16) {
				return 1;
			}
		}
		return 0;
	}
	v1 = sub_424D00();
	if (v1 == -1) {
		return 1;
	}
	do {
		if (nox_xxx_testUnitBuffs_4FF350((int)a1, v1)) {
			v2 = nox_xxx_getEnchantName_4248F0(v1);
			LOBYTE(v11) = strlen(v2);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v11, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2, (unsigned char)v11);
			LOBYTE(v15) = nox_xxx_buffGetPower_4FF570((int)a1, v1);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v15, 1u);
			v14 = nox_xxx_unitGetBuffTimer_4FF550((int)a1, v1);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v14, 4u);
			if (v1 == 26) {
				v3 = sub_4FF2D0(51, (int)a1);
				if (v3) {
					v12 = *(uint32_t*)(v3 + 72);
				} else {
					v12 = 100;
				}
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v12, 4u);
			}
		}
		v1 = sub_424D20(v1);
	} while (v1 != -1);
	return 1;
}
// 52AAB0: using guessed type char var_100[256];

//----- (0052AD10) --------------------------------------------------------
size_t nox_xxx_readNPCVoiceSet_52AD10(int a1) {
	const char** v1; // eax
	const char** v2; // esi
	size_t result;   // eax
	const char** v4; // eax
	int v5;          // [esp+4h] [ebp-104h]
	char v6[256];    // [esp+8h] [ebp-100h]

	if (nox_crypt_IsReadOnly()) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(v6, (unsigned char)v5);
		v6[(unsigned char)v5] = 0;
		v4 = nox_xxx_getDefaultSoundSet_424350(v6);
		result = nox_xxx_setNPCVoiceSet_424320(a1, (int)v4);
	} else {
		v1 = (const char**)nox_xxx_monsterGetSoundSet_424300(a1);
		v2 = v1;
		if (v1) {
			LOBYTE(v5) = strlen(*v1);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 1u);
			result = nox_xxx_fileReadWrite_426AC0_file3_fread(*v2, (unsigned char)v5);
		} else {
			LOBYTE(v5) = 0;
			result = nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 1u);
		}
	}
	return result;
}
// 52AD10: using guessed type char var_100[256];

//----- (0052ADE0) --------------------------------------------------------
int nox_xxx_XFerNPC_52ADE0(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;            // esi
	char* v2;          // edi
	int result;        // eax
	uint32_t* v4;      // eax
	char* v5;          // eax
	char* v6;          // eax
	char* v7;          // eax
	char* v8;          // eax
	char* v9;          // eax
	char* v10;         // eax
	char* v11;         // eax
	char* v12;         // eax
	char* v13;         // eax
	char* v14;         // edi
	int v15;           // eax
	short v16;         // ax
	unsigned char v17; // bl
	uint8_t* v18;      // edi
	int v19;           // ebx
	unsigned char i;   // bl
	int v21;           // edx
	uint16_t* v22;     // eax
	uint16_t* v23;     // eax
	int v24;           // edx
	uint8_t* v25;      // ebx
	uint32_t* v26;     // eax
	int v27;           // ecx
	char* v28;         // ebp
	bool v29;          // zf
	unsigned char j;   // bl
	int v31;           // eax
	int* v32;          // ebp
	char* v33;         // ebx
	char* v34;         // ebx
	int v35;           // edi
	int v36;           // eax
	int v37;           // eax
	int v38;           // eax
	int v39;           // eax
	uint32_t* k;       // esi
	int v41;           // eax
	int v42;           // eax
	int v43;           // [esp+10h] [ebp-238h]
	int v44;           // [esp+14h] [ebp-234h]
	int v45;           // [esp+18h] [ebp-230h]
	unsigned char v46; // [esp+1Fh] [ebp-229h]
	int v47;           // [esp+20h] [ebp-228h]
	int v48;           // [esp+24h] [ebp-224h]
	int v49;           // [esp+28h] [ebp-220h]
	int v50;           // [esp+2Ch] [ebp-21Ch]
	int v51;           // [esp+30h] [ebp-218h]
	int v52;           // [esp+34h] [ebp-214h]
	char v53[3];       // [esp+38h] [ebp-210h]
	int v54;           // [esp+3Ch] [ebp-20Ch]
	uint32_t v55[2];   // [esp+40h] [ebp-208h]
	char v56[256];     // [esp+48h] [ebp-200h]
	char v57[256];     // [esp+148h] [ebp-100h]

	v1 = a1p->data_update;
	v2 = *(char**)(a1 + 756);
	v54 = *(uint32_t*)(a1 + 136);
	v44 = 62;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 2u);
	if ((short)v44 > 62) {
		return 0;
	}
	if (nox_crypt_IsReadOnly() == 1) {
		v4 = nox_xxx_monsterDefByTT_517560(*(unsigned short*)(a1 + 4));
		*(uint32_t*)(v1 + 484) = v4;
		if (v4) {
			*(uint32_t*)(v1 + 1440) = v4[23];
		}
	}
	result = nox_xxx_mapReadWriteObjData_4F4530((int*)a1, (short)v44);
	if (!result) {
		return 0;
	}
	if (!nox_crypt_IsReadOnly()) {
		nox_xxx_xferIndexedDirection_509E20(*(short*)(a1 + 124), (int2*)v55);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(v55, 8u);
	if (v2) {
		v5 = v2 + 640;
	} else {
		v5 = 0;
	}
	nox_xxx_xferReadScriptHandler_4F5580(v1 + 1232, v5);
	if (v2) {
		v6 = v2 + 896;
	} else {
		v6 = 0;
	}
	nox_xxx_xferReadScriptHandler_4F5580(v1 + 1264, v6);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1220), 2u);
	if (v2) {
		v7 = v2 + 768;
	} else {
		v7 = 0;
	}
	nox_xxx_xferReadScriptHandler_4F5580(v1 + 1224, v7);
	if ((short)v44 >= 32) {
		v8 = v2 ? v2 + 1024 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1240, v8);
		v9 = v2 ? v2 + 1152 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1248, v9);
		v10 = v2 ? v2 + 1280 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1256, v10);
		v11 = v2 ? v2 + 1408 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1272, v11);
		v12 = v2 ? v2 + 1536 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1280, v12);
		v13 = v2 ? v2 + 1664 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1288, v13);
		if ((short)v44 >= 50) {
			if (v2) {
				v14 = v2 + 1792;
			} else {
				v14 = 0;
			}
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1296, v14);
		}
	}
	v47 = 0; // FIXME: set to direction? was uninitialized
	v15 = nox_crypt_IsReadOnly();
	if (nox_crypt_IsReadOnly() != 1 ||
		(v16 = nox_xxx_xferDirectionToAngle_509E00(v55), *(uint16_t*)(a1 + 126) = v16, *(uint16_t*)(a1 + 124) = v16,
		 v15 = nox_crypt_IsReadOnly(), nox_crypt_IsReadOnly() != 1)) {
		if (!v15) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 4u);
		}
	} else {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 4u);
	}
	if (nox_crypt_IsReadOnly() == 1) {
		v17 = 0;
		LOBYTE(v48) = 0;
		do {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v53, 3u);
			nox_xxx_setNPCColor_4E4A90(a1, v48, v53);
			LOBYTE(v48) = ++v17;
		} while (v17 < 6u);
	} else {
		v18 = (uint8_t*)(v1 + 2076);
		v19 = 6;
		do {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v18, 3u);
			v18 += 3;
			--v19;
		} while (v19);
	}
	if (nox_crypt_IsReadOnly() == 1 && (uint16_t)v44 == 31) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 1u);
		for (i = 0; i < v46; ++i) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v50, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v56, (unsigned char)v50);
			v56[(unsigned char)v50] = 0;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
		}
	}
	if ((short)v44 >= 32) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1332), 1u);
		if ((short)v44 < 49) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 2u);
			*(uint32_t*)(v1 + 1440) = (unsigned short)v47;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1440), 4u);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1352), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1336), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1344), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1312), 4u);
		v43 = 0;
		v22 = *(uint16_t**)(a1 + 556);
		if (v22) {
			LOWORD(v21) = *v22;
			v43 = v21;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 2u);
		v23 = *(uint16_t**)(a1 + 556);
		if (v23) {
			*v23 = v43;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1304), 4u);
		*(uint32_t*)(v1 + 1308) = *(uint32_t*)(v1 + 1304);
		if ((short)v44 < 35) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1360), 4u);
		}
		LOBYTE(v49) = strlen((const char*)(v1 + 1364));
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1364), (unsigned char)v49);
		*(uint8_t*)((unsigned char)v49 + v1 + 1364) = 0;
		if ((short)v44 >= 34) {
			if (nox_crypt_IsReadOnly()) {
				memset((void*)(v1 + 1488), 0, 0x224u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 4u);
				for (j = 0; j < v43; LOBYTE(v48) = j) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 1u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v56, (unsigned char)v49);
					v56[(unsigned char)v49] = 0;
					v31 = nox_xxx_spellNameToN_4243F0(v56);
					nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 4 * v31 + 1488), 4u);
					++j;
				}
			} else {
				v24 = 0;
				v25 = (uint8_t*)(v1 + 1488);
				v43 = 0;
				v26 = (uint32_t*)(v1 + 1488);
				v27 = 137;
				do {
					if (*v26) {
						++v24;
					}
					++v26;
					--v27;
				} while (v27);
				v43 = v24;
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 4u);
				v48 = 0;
				v47 = 137;
				do {
					if (*(uint32_t*)v25) {
						v28 = nox_xxx_spellNameByN_424870(v48);
						LOBYTE(v52) = strlen(v28);
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v52, 1u);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v28, (unsigned char)v52);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v25, 4u);
					}
					v25 += 4;
					v29 = v47 == 1;
					++v48;
					--v47;
				} while (!v29);
			}
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1488), 0x224u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1448) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1450) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1448), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1450), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1456) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1458) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1456), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1458), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1464) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1466) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1464), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1466), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1472) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1474) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1472), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1474), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1480) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1482) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1480), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1482), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 >= 33) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1316), 4u);
		}
		if ((short)v44 >= 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2040), 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1324), 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1328), 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1320), 4u);
			if ((short)v44 < 42) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 2u);
				if (!(uint16_t)v43) {
					*(uint8_t*)(v1 + 1445) = 1;
				}
			}
			v32 = (int*)(v1 + 2044);
			v47 = 3;
			do {
				if (nox_crypt_IsReadOnly() == 1) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v56, (unsigned char)v43);
					v56[(unsigned char)v43] = 0;
					*v32 = nox_xxx_spellNameToN_4243F0(v56);
				} else {
					v33 = nox_xxx_spellNameByN_424870(*v32);
					LOBYTE(v43) = strlen(v33);
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v33, (unsigned char)v43);
				}
				++v32;
				--v47;
			} while (v47);
		}
		if ((short)v44 >= 35) {
			if (nox_crypt_IsReadOnly()) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(v57, (unsigned char)v43);
				v57[(unsigned char)v43] = 0;
				*(uint32_t*)(v1 + 1360) = nox_xxx_actionNByNameMB_5345F0(v57);
			} else {
				v34 = sub_5345B0(*(uint32_t*)(v1 + 1360));
				LOBYTE(v43) = strlen(v34);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(v34, (unsigned char)v43);
			}
		}
	}
	if ((short)v44 < 41) {
		v35 = a1;
	} else {
		v35 = a1;
		result = nox_xxx_XFer_ActionData_529CE0(a1);
		if (!result) {
			return result;
		}
	}
	if ((short)v44 >= 42) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1445), 1u);
	}
	if ((short)v44 >= 44) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v1, 4u);
	}
	if ((short)v44 >= 45) {
		v36 = *(uint32_t*)(v35 + 556);
		v45 = 0;
		if (v36) {
			v45 = *(unsigned short*)(v36 + 4);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		v37 = *(uint32_t*)(v35 + 556);
		if (v37) {
			*(uint16_t*)(v37 + 4) = v45;
		}
	}
	if ((short)v44 >= 46) {
		v51 = *(uint32_t*)(v35 + 12) & 0x180;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v51, 4u);
		*(uint32_t*)(v35 + 12) |= v51;
	}
	if ((short)v44 >= 48) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(*(uint8_t**)(v35 + 556), 2u);
	}
	if ((short)v44 >= 51) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v35 + 28), 4u);
	}
	if ((short)v44 >= 52) {
		nox_xxx_readNPCVoiceSet_52AD10(v35);
	}
	if (!((short)v44 < 61 || (result = nox_xxx_XFer_ReadMonsterBuffs_52AAB0((uint32_t*)v35)) != 0)) {
		return result;
	}
	if ((short)v44 < 62) {
		if (nox_crypt_IsReadOnly() == 1) {
			goto LABEL_149;
		}
		goto LABEL_156;
	}
	LOBYTE(v45) = *(uint8_t*)(v35 + 540);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 1u);
	if (nox_crypt_IsReadOnly() != 1) {
		goto LABEL_156;
	}
	if ((uint8_t)v45) {
		nox_xxx_setSomePoisonData_4EEA90(v35, (unsigned char)v45);
		if (nox_crypt_IsReadOnly() == 1) {
			goto LABEL_149;
		}
		goto LABEL_156;
	}
LABEL_149:
	if (!*(uint8_t*)(v1 + 1445)) {
		goto LABEL_170;
	}
	if (nox_common_gameFlags_check_40A5C0(1)) {
		v38 = *(uint32_t*)(v35 + 556);
		if (v38) {
			*(uint16_t*)(v38 + 4) = 0;
			**(uint16_t**)(v35 + 556) = 0;
		}
	}
	if (nox_crypt_IsReadOnly() != 1) {
		goto LABEL_156;
	}
LABEL_170:
	v39 = *(uint32_t*)(v35 + 16);
	if ((v39 & 0x8000) != 0) {
		LOBYTE(v39) = v39 | 0x40;
		*(uint32_t*)(v35 + 16) = v39;
	}
LABEL_156:
	if (!*(uint32_t*)(v35 + 136) || nox_crypt_IsReadOnly() != 1 ||
		(result = nox_xxx_xfer_4F3E30(v44, v35, *(uint32_t*)(v35 + 136))) != 0) {
		sub_52BA70(v35);
		if (nox_crypt_IsReadOnly() == 1) {
			if (nox_common_gameFlags_check_40A5C0(1)) {
				for (k = *(uint32_t**)(v35 + 504); k; k = (uint32_t*)k[124]) {
					v41 = k[4];
					if (v41 & 0x100) {
						v42 = k[2];
						k[4] &= 0xFFFFFEFF;
						if (v42 & 0x1001000) {
							nox_xxx_NPCEquipWeapon_53A2C0(v35, (int)k);
						} else {
							nox_xxx_NPCEquipArmor_53E520(v35, k);
						}
					}
				}
			}
		}
		*(uint32_t*)(v35 + 136) = v54;
		return 1;
	}
	return result;
}
// 52B1CD: variable 'v21' is possibly undefined
// 52ADE0: using guessed type char var_200[256];
// 52ADE0: using guessed type char var_100[256];

//----- (0052BA70) --------------------------------------------------------
int sub_52BA70(int a1) {
	int result; // eax
	int v2;     // edi
	int v3;     // esi
	int v4;     // ecx
	int v5;     // edx

	result = a1;
	v2 = 0;
	v3 = 0;
	if (a1) {
		for (result = *(uint32_t*)(a1 + 504); result; result = *(uint32_t*)(result + 496)) {
			v4 = *(uint32_t*)(result + 16);
			if (v4 & 0x100) {
				v5 = *(uint32_t*)(result + 8);
				if (v5 & 0x1001000 && *(uint32_t*)(result + 12) & 0x7FFE40C) {
					if (v3 == 1) {
						BYTE1(v4) &= 0xFEu;
						*(uint32_t*)(result + 16) = v4;
					} else {
						v2 = 1;
					}
				} else if (v5 & 0x2000000 && *(uint8_t*)(result + 12) & 2) {
					if (v2 == 1) {
						BYTE1(v4) &= 0xFEu;
						*(uint32_t*)(result + 16) = v4;
					} else {
						v3 = 1;
					}
				}
			}
		}
	}
	return result;
}

//----- (0052BAF0) --------------------------------------------------------
int sub_52BAF0(int a1) {
	int v1;       // eax
	int result;   // eax
	uint32_t* v3; // ebx
	int v4;       // edi
	int v5;       // ebp
	int* v6;      // esi
	int v7;       // [esp+0h] [ebp-4h]
	int v8;       // [esp+8h] [ebp+4h]

	v1 = a1;
	v8 = 0;
	result = *(uint32_t*)(v1 + 748);
	v7 = result;
	if (*(uint8_t*)(result + 544) & 0x80) {
		return result;
	}
	v3 = (uint32_t*)(result + 552);
	do {
		v4 = 0;
		v5 = *getMemU32Ptr(0x587000, 255604 + 16 * *v3);
		if (v5 <= 0) {
			goto LABEL_14;
		}
		v6 = v3 + 1;
		do {
			if (*getMemU32Ptr(0x587000, 255608 + 4 * (v4 + 4 * *v3)) == 1) {
				if (*v6) {
					*v6 = sub_4ECF10(*v6);
					goto LABEL_12;
				}
			} else {
				if (*getMemU32Ptr(0x587000, 255608 + 4 * (v4 + 4 * *v3)) != 2) {
					goto LABEL_12;
				}
				if (*v6) {
					*v6 = nox_server_getWaypointById_579C40(*v6);
					goto LABEL_12;
				}
			}
			*v6 = 0;
		LABEL_12:
			++v4;
			v6 += 2;
		} while (v4 < v5);
		result = v7;
	LABEL_14:
		v3 += 6;
		++v8;
	} while (v8 <= *(char*)(result + 544));
	return result;
}
