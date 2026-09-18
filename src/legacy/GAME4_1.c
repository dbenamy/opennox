#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__system__ctrlevnt.h"
#include "common__binfile.h"
#include "common__crypt.h"
#include "common__net_list.h"
#include "common__random.h"
#include "operators.h"
#include "common__system__team.h"
#include "server__gamemech__explevel.h"
#include "server__magic__plyrspel.h"
#include "server__system__trade.h"

#include "client__gui__window.h"
#include "client__video__draw_common.h"
#include "common__magic__speltree.h"
#include "server__script__script.h"

extern uint32_t dword_5d4594_3835364;
extern uint32_t dword_5d4594_1599708;
extern uint32_t dword_587000_234176;
extern uint32_t dword_5d4594_2487244;
extern uint32_t dword_5d4594_2386228;
extern void* nox_alloc_spawn_2386216;
extern uint32_t dword_5d4594_3835348;
extern void* nox_alloc_tradeSession_2386492;
extern void* nox_alloc_monsterList_2386220;
extern uint32_t dword_5d4594_2386500;
extern uint32_t dword_5d4594_2386212;
extern void* nox_alloc_tradeItems_2386496;
extern uint32_t dword_5d4594_2386224;

extern uint32_t nox_tile_def_cnt;
extern nox_tileDef_t nox_tile_defs_arr[176];





















































































//----- (00516570) --------------------------------------------------------
extern uint32_t nox_gameDisableMapDraw_5d4594_2650672;


//----- (00516D00) --------------------------------------------------------


//----- (00516EE0) --------------------------------------------------------


//----- (00516F10) --------------------------------------------------------


//----- (00516F30) --------------------------------------------------------


//----- (00516F90) --------------------------------------------------------


//----- (00516FC0) --------------------------------------------------------


//----- (00517010) --------------------------------------------------------


//----- (00517090) --------------------------------------------------------


//----- (00517140) --------------------------------------------------------


//----- (00517170) --------------------------------------------------------


//----- (005174F0) --------------------------------------------------------


//----- (00517520) --------------------------------------------------------


//----- (00517560) --------------------------------------------------------




//----- (0051A500) --------------------------------------------------------
int sub_51A500(int a1) {
	int v1;           // ecx
	int v2;           // edx
	unsigned char* i; // eax
	int v4;           // esi

	if (!*getMemU32Ptr(0x5D4594, 2388664)) {
		sub_51A550();
	}
	if (!a1) {
		return 0;
	}
	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 249904)) {
		return 0;
	}
	HIWORD(v2) = 0;
	for (i = getMemAt(0x587000, 249904);; i += 16) {
		LOWORD(v2) = *(uint16_t*)(a1 + 4);
		if (*((uint32_t*)i - 1) == v2) {
			break;
		}
		v4 = *((uint32_t*)i + 4);
		++v1;
		if (!v4) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 249908 + 16 * v1);
}

//----- (0051A550) --------------------------------------------------------
char* sub_51A550() {
	char* result;      // eax
	unsigned char* v1; // esi

	result = *(char**)getMemAt(0x587000, 249904);
	if (*getMemU32Ptr(0x587000, 249904)) {
		v1 = getMemAt(0x587000, 249896);
		do {
			*((uint32_t*)v1 + 3) = nox_xxx_getNameId_4E3AA0(result);
			*((uint32_t*)v1 + 1) = nox_xxx_getNameId_4E3AA0(*(char**)v1);
			result = (char*)*((uint32_t*)v1 + 6);
			v1 += 16;
		} while (result);
		*getMemU32Ptr(0x5D4594, 2388664) = 1;
	} else {
		*getMemU32Ptr(0x5D4594, 2388664) = 1;
	}
	return result;
}

//----- (0051A5A0) --------------------------------------------------------
void* nox_xxx_objectTypeByIndHealthData(int a1);
void nox_xxx_spawnHecubahQuest_51A5A0(int* a1) {
	uint32_t* v1;  // edi
	uint32_t* v2;  // esi
	int v3;        // eax
	int v4;        // eax
	double v5;     // st7
	int v6;        // eax
	uint32_t* v7;  // esi
	int v8;        // eax
	uint32_t* v9;  // eax
	int v10;       // eax
	uint32_t* v11; // eax
	int v12;       // eax
	uint32_t* v13; // eax
	int v14;       // eax
	uint32_t* v15; // eax
	float v16;     // [esp+8h] [ebp-8h]
	float v17;     // [esp+Ch] [ebp-4h]

	v1 = nox_xxx_newObjectByTypeID_4E3810("Hecubah");
	v16 = sub_4E40F0();
	if (v1) {
		v2 = (uint32_t*)v1[187];
		v3 = v2[121];
		if (v3) {
			v4 = *(uint32_t*)(v3 + 72);
		} else {
			v4 = *(unsigned short*)((int)nox_xxx_objectTypeByIndHealthData(*((unsigned short*)v1 + 2)) + 4);
		}
		if (v16 < 1.0) {
			v16 = 1.0;
		}
		v5 = (double)v4 * v16;
		v17 = v5;
		nox_xxx_unitSetHP_4E4560((int)v1, (long long)v5);
		*(uint16_t*)(v1[139] + 4) = nox_float2int(v17);
		if (!*(uint16_t*)v1[139]) {
			nox_xxx_unitSetHP_4E4560((int)v1, 1u);
		}
		v6 = v1[139];
		if (!*(uint16_t*)(v6 + 4)) {
			*(uint16_t*)(v6 + 4) = 1;
		}
		v2[411] = 0x10000000;
		v2[423] = 0x10000000;
		v2[340] = 4;
		v2[326] = 1062501089;
		v2[510] = 3;
		v2[410] = 0x8000000;
		v2[444] = 0x20000000;
		v2[388] = 0x40000000;
		v2[415] = 0x40000000;
		nox_xxx_gamedataGetFloat_419D40("HecubahQuestSkill");
		v2[330] = 1062836634;
		nox_xxx_createAt_4DAA50((int)v1, 0, *(float*)a1, *((float*)a1 + 1));
		v7 = nox_xxx_newObjectByTypeID_4E3810("RewardMarker");
		if (v7) {
			v8 = nox_game_getQuestStage_4E3CC0();
			v9 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v8 + 2);
			if (v9) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v9, 0);
			}
			v10 = nox_game_getQuestStage_4E3CC0();
			v11 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v10 + 2);
			if (v11) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v11, 0);
			}
			v12 = nox_game_getQuestStage_4E3CC0();
			v13 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v12 + 2);
			if (v13) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v13, 0);
			}
			v14 = nox_game_getQuestStage_4E3CC0();
			v15 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v14 + 2);
			if (v15) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v15, 0);
			}
			nox_xxx_objectFreeMem_4E38A0((int)v7);
		}
	}
}

//----- (0051A7A0) --------------------------------------------------------
void nox_xxx_spawnNecroQuest_51A7A0(int* a1) {
	uint32_t* v1; // edi
	uint32_t* v2; // esi
	int v3;       // eax
	int v4;       // eax
	double v5;    // st7
	int v6;       // eax
	uint32_t* v7; // esi
	int v8;       // eax
	uint32_t* v9; // eax
	float v10;    // [esp+8h] [ebp-8h]
	float v11;    // [esp+Ch] [ebp-4h]

	v1 = nox_xxx_newObjectByTypeID_4E3810("Necromancer");
	v10 = sub_4E40F0();
	if (v1) {
		v2 = (uint32_t*)v1[187];
		v3 = v2[121];
		if (v3) {
			v4 = *(uint32_t*)(v3 + 72);
		} else {
			v4 = *(unsigned short*)((int)nox_xxx_objectTypeByIndHealthData(*((unsigned short*)v1 + 2)) + 4);
		}
		if (v10 < 1.0) {
			v10 = 1.0;
		}
		v5 = (double)v4 * v10;
		v11 = v5;
		nox_xxx_unitSetHP_4E4560((int)v1, (long long)v5);
		*(uint16_t*)(v1[139] + 4) = nox_float2int(v11);
		if (!*(uint16_t*)v1[139]) {
			nox_xxx_unitSetHP_4E4560((int)v1, 1u);
		}
		v6 = v1[139];
		if (!*(uint16_t*)(v6 + 4)) {
			*(uint16_t*)(v6 + 4) = 1;
		}
		v2[340] = 4;
		v2[411] = 0x10000000;
		v2[423] = 0x10000000;
		v2[326] = 1062501089;
		v2[510] = 1;
		v2[410] = 0x8000000;
		v2[444] = 0x20000000;
		v2[415] = 0x40000000;
		nox_xxx_createAt_4DAA50((int)v1, 0, *(float*)a1, *((float*)a1 + 1));
		v7 = nox_xxx_newObjectByTypeID_4E3810("RewardMarker");
		if (v7) {
			v8 = nox_game_getQuestStage_4E3CC0();
			v9 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v8 + 2);
			if (v9) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v9, 0);
			}
			nox_xxx_objectFreeMem_4E38A0((int)v7);
		}
	}
}

//----- (0051A930) --------------------------------------------------------
int nox_xxx_getQuestStage_51A930() { return *getMemU32Ptr(0x5D4594, 2388660); }

//----- (0051A940) --------------------------------------------------------
int sub_51A940(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 2388656) = a1;
	return result;
}

//----- (0051A950) --------------------------------------------------------
int sub_51A950() { return *getMemU32Ptr(0x5D4594, 2388656); }




//----- (0051D0E0) --------------------------------------------------------
void sub_51D0E0() { dword_5d4594_2487244 = 0; }

//----- (0051D0F0) --------------------------------------------------------
int sub_51D0F0(char a1) {
	*getMemU8Ptr(0x973F18, 35972) = a1;
	return 1;
}

//----- (0051D100) --------------------------------------------------------
int sub_51D100(int a1) {
	if (a1 != 1 && a1) {
		return 0;
	}
	*getMemU32Ptr(0x973F18, 35976) = a1;
	return 1;
}

//----- (0051D120) --------------------------------------------------------
nox_waypoint_t* nox_xxx_waypointNew_5798F0(float a1, float a2);
uint32_t* sub_51D120(float* a1) {
	uint32_t* result; // eax
	uint32_t* v2;     // esi
	float2 v3;        // [esp+4h] [ebp-8h]

	result = (float*)nox_xxx_mapGenFixCoords_4D3D90((float2*)a1, &v3);
	if (result) {
		result = nox_xxx_waypointNew_5798F0(v3.field_0, v3.field_4);
		v2 = result;
		if (result) {
			if (dword_5d4594_2487244) {
				if (*getMemU32Ptr(0x973F18, 35976) == 1) {
					sub_51D300(*(int*)&dword_5d4594_2487244, (int)result, getMemByte(0x973F18, 35972));
					sub_51D300((int)v2, *(int*)&dword_5d4594_2487244, getMemByte(0x973F18, 35972));
				}
			}
			dword_5d4594_2487244 = v2;
			result = v2;
		}
	}
	return result;
}

//----- (0051D1A0) --------------------------------------------------------
nox_waypoint_t* sub_579AD0(float a1, float a2);
float* sub_51D1A0(float2* a1) {
	float* result; // eax
	float2 a2;     // [esp+0h] [ebp-8h]

	result = (float*)nox_xxx_mapGenFixCoords_4D3D90(a1, &a2);
	if (result) {
		result = sub_579AD0(a2.field_0, a2.field_4);
	}
	return result;
}

//----- (0051D3F0) --------------------------------------------------------
float2* sub_51D3F0(float2* a1, float2* a2) {
	float2* result; // eax
	float2* v3;     // esi

	result = a1;
	if (a1) {
		if (a2) {
			result = (float2*)sub_51D1A0(a1);
			v3 = result;
			if (result) {
				result = (float2*)sub_51D1A0(a2);
				if (result) {
					result = (float2*)sub_51D2C0((int)v3, (int)result);
				}
			}
		} else {
			result = 0;
		}
	}
	return result;
}
