#include <math.h>
#include <time.h>

// For htonl
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
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME5_2.h"
#include "common__system__team.h"
#include "server__ability__ability.h"
#include "server__magic__plyrgide.h"
#include "server__magic__plyrspel.h"
#include "server__system__server.h"

#include "client__gui__guiinv.h"
#include "client__gui__guijourn.h"
#include "client__gui__servopts__guiserv.h"
#include "client__gui__servopts__playrlst.h"
#include "common__binfile.h"
#include "common__crypt.h"
#include "common__log.h"
#include "common__random.h"
#include "operators.h"

#include "client__gui__guibook.h"

#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"
#include "server__script__builtin.h"
#include "server__script__script.h"

extern uint32_t dword_5d4594_588084;
extern uint32_t nox_player_netCode_85319C;
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern int ptr_5D4594_2650668_cap;

double sub_419A10(float a1) {
	*getMemFloatPtr(0x5D4594, 527672) = a1;
	**(uint32_t**)getMemAt(0x587000, 55744) &= 0x7FFFFFFFu;
	return *getMemFloatPtr(0x5D4594, 527672);
}

unsigned int sub_419A30(float a1) {
	unsigned int result; // eax

	if (a1 < 0.0) {
		return 0;
	}
	*getMemU32Ptr(0x5D4594, 527668) = getMemAt(0x5D4594, 527676);
	*getMemFloatPtr(0x5D4594, 527676) = a1 + 8388608.0;
	result = 0x7fffff & *getMemU32Ptr(0x5D4594, 527676);
	*getMemU32Ptr(0x5D4594, 527680) = 0x7fffff & *getMemU32Ptr(0x5D4594, 527676);
	return result;
}

int nox_float2int(float a1) { return (int)a1; }

short nox_float2int16(float a1) { return (int)a1; }

float nox_double2float(double a1) { return (float)a1; }

int nox_double2int(double a1) { return (int)a1; }





















































int sub_4254A0(int a1, uint8_t* a2) {
	*(uint32_t*)a1 = a2;
	*(uint8_t*)(a1 + 4) = 0;
	return *a2 & 1;
}

bool sub_4254C0(unsigned char** a1) {
	char v1;           // cl
	unsigned char* v2; // ecx

	v1 = *((uint8_t*)a1 + 4) + 1;
	*((uint8_t*)a1 + 4) = v1;
	if (v1 == 8) {
		v2 = *a1;
		*((uint8_t*)a1 + 4) = 0;
		*a1 = v2 + 1;
	}
	return ((1 << *((uint8_t*)a1 + 4)) & **a1) > 0;
}

uint8_t* sub_425500(int a1, uint8_t* a2, char a3) {
	uint8_t* result; // eax

	result = a2;
	*(uint32_t*)a1 = a2;
	*(uint8_t*)(a1 + 4) = 0;
	*a2 = a3;
	return result;
}

char sub_425520(int a1, char a2) {
	char v2;     // cl
	uint8_t* v3; // ecx
	char result; // al

	v2 = *(uint8_t*)(a1 + 4) + 1;
	*(uint8_t*)(a1 + 4) = v2;
	if (v2 == 8) {
		v3 = *(uint8_t**)a1;
		*(uint8_t*)(a1 + 4) = 0;
		*(uint32_t*)a1 = ++v3;
		*v3 = 0;
	}
	result = a2 << *(uint8_t*)(a1 + 4);
	**(uint8_t**)a1 |= result;
	return result;
}

int sub_425550(uint8_t* a1, uint8_t* a2, int a3) {
	int v3;     // edi
	int v4;     // esi
	char v5;    // al
	int v6;     // ebx
	char v8[8]; // [esp+8h] [ebp-10h]
	char v9[8]; // [esp+10h] [ebp-8h]
	bool v10;   // [esp+1Ch] [ebp+4h]

	v3 = 1;
	v4 = 0;
	v5 = sub_4254A0((int)v9, a1);
	sub_425500((int)v8, a2, v5);
	if (a3 == 1) {
		return 1;
	}
	v6 = a3 - 1;
	do {
		if (!(++v4 % 7u)) {
			sub_425520((int)v8, 1);
			++v3;
		}
		v10 = sub_4254C0((unsigned char**)v9);
		sub_425520((int)v8, v10);
		--v6;
	} while (v6);
	return v3;
}

int nox_xxx_countObserverPlayers_425BF0() {
	int v0;  // esi
	char* i; // eax
	int v2;  // ecx

	v0 = 0;
	if (nox_common_gameFlags_check_40A5C0(0x8000)) {
		for (i = nox_common_playerInfoGetFirst_416EA0(); i; i = nox_common_playerInfoGetNext_416EE0((int)i)) {
			v2 = *((uint32_t*)i + 920);
			if (v2 & 1 && !(v2 & 0x20) && i[2064] != 31) {
				++v0;
			}
		}
	}
	return v0;
}

int nox_xxx_wallGet_426A30() { return *getMemU32Ptr(0x5D4594, 739992); }

char* nox_xxx_mapGetWallSize_426A70() { return (char*)getMemAt(0x5D4594, 739980); }

void nox_xxx_mapWall_426A80(int* a1) {
	*getMemU32Ptr(0x5D4594, 739980) = a1[0];
	*getMemU32Ptr(0x5D4594, 739984) = a1[1];
}
