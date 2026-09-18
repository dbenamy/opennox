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


//----- (0051A550) --------------------------------------------------------


//----- (0051A5A0) --------------------------------------------------------


//----- (0051A7A0) --------------------------------------------------------


//----- (0051A930) --------------------------------------------------------


//----- (0051A940) --------------------------------------------------------


//----- (0051A950) --------------------------------------------------------





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
