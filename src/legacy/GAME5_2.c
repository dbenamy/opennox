#include <float.h>
#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
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
#include "common__binfile.h"
#include "common__random.h"
#include "common__strman.h"

#include "client__gui__guiinv.h"
#include "client__shell__mainmenu.h"
#include "common__system__team.h"

#include "MixPatch.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"
#include "operators.h"

extern uint32_t dword_5d4594_2523804;
extern uint32_t dword_5d4594_2523776;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_2523780;
extern uint32_t dword_5d4594_2650652;
extern uint32_t dword_8531A0_2576;

//----- (00554040) --------------------------------------------------------
unsigned int nox_server_makeServerInfoPacket_554040(const char* inBuf, int inSz, char* out) {
	char buf[72];

	char* v3 = sub_416640();
	char* game = nox_xxx_cliGamedataGet_416590(0);
	if (!sub_43AF30() || sub_4D6F30()) {
		return 0;
	}
	char playerLimit = nox_xxx_servGetPlrLimit_409FA0();
	char playerCount = nox_common_playerInfoCount_416F40();
	if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING)) {
		--playerCount;
		--playerLimit;
	}
	char* srvName = nox_xxx_serverOptionsGetServername_40A4C0();
	buf[0] = 0;
	buf[1] = 0;
	buf[2] = 13;
	buf[3] = playerCount;
	buf[4] = playerLimit;
	buf[5] = v3[101] & 0xF;
	buf[6] = ((unsigned char)v3[101]) >> 4;
	*(uint32_t*)&buf[7] = *((uint32_t*)game + 11);
	strcpy(&buf[10], nox_xxx_mapGetMapName_409B40());
	buf[19] = v3[102] | sub_43BE50_get_video_mode_id();
	buf[20] = v3[100];
	buf[21] = v3[100] & 0x10;
	*(uint32_t*)&buf[24] = *((uint32_t*)v3 + 12);
	unsigned int gameFlags = nox_common_gameFlags_getVal_40A5B0();
	if (nox_xxx_isQuest_4D6F50()) {
		gameFlags = (gameFlags & 0xFFFFFF7Fu) | 0x1000u;
		*(uint16_t*)&buf[68] = nox_game_getQuestStage_4E3CC0();
	}
	*(uint32_t*)&buf[28] = gameFlags;
	*(uint32_t*)&buf[32] = *((uint32_t*)game + 12);
	*(uint16_t*)&buf[36] = *(uint16_t*)(v3 + 105);
	*(uint16_t*)&buf[38] = *(uint16_t*)(v3 + 107);
	*(uint32_t*)&buf[40] = *((uint32_t*)v3 + 11);
	*(uint32_t*)&buf[44] = *(uint32_t*)(&inBuf[8]); // timestamp of the packet
	memcpy(&buf[48], game + 24, 20);
	memcpy(&out[0], buf, 72);
	strcpy(&out[72], srvName);
	return 72 + strlen(srvName) + 1;
}

//----- (0057ADF0) --------------------------------------------------------
int* sub_57ADF0(int* a1) {
	int* result; // eax
	int* v2;     // esi
	int* v3;     // edi

	result = nox_common_list_getFirstSafe_425890(a1);
	v2 = result;
	if (result) {
		do {
			v3 = nox_common_list_getNextSafe_4258A0(v2);
			nox_common_list_remove_425920((uint32_t**)v2);
			free(v2);
			v2 = v3;
		} while (v3);
	}
	return result;
}

//----- (0057AF20) --------------------------------------------------------
int nox_xxx_get_57AF20() { return dword_5d4594_2523804; }

//----- (0057B0A0) --------------------------------------------------------
void sub_57B0A0() {
	int result;   // eax
	uint32_t* v1; // ecx

	result = dword_5d4594_2523804;
	if (!result) {
		return;
	}
	v1 = *(uint32_t**)&dword_5d4594_2523780;
	if (dword_5d4594_2523780 && (!*getMemU32Ptr(0x5D4594, 2523772) || *getMemU32Ptr(0x5D4594, 2523772) == 1)) {
		nox_xxx_netSendPointFx_522FF0(154, (float2*)(dword_5d4594_2523780 + 56));
		v1 = *(uint32_t**)&dword_5d4594_2523780;
	}
	if (dword_5d4594_2523776) {
		nox_xxx_delayedDeleteObject_4E5CC0(*(int*)&dword_5d4594_2523776);
		v1 = *(uint32_t**)&dword_5d4594_2523780;
	}
	dword_5d4594_2523776 = 0;
	if (v1) {
		nox_xxx_playerSetState_4FA020(v1, 13);
	}
	dword_5d4594_2523780 = 0;
	if (!sub_45D9B0()) {
		sub_413A00(0);
	}
	dword_5d4594_2523804 = 0;
}

//----- (0057B180) --------------------------------------------------------
long long nox_xxx___Getcvt_57B180() { return *getMemU64Ptr(0x5D4594, 2523788); }

//----- (0057B3D0) --------------------------------------------------------
int nox_cheat_allowall = 0;

//----- (0057C090) --------------------------------------------------------
int nox_server_getNextMapGroup_57C090(int a1) {
	int result; // eax

	if (a1) {
		result = *(uint32_t*)(a1 + 88);
	} else {
		result = 0;
	}
	return result;
}

//----- (0057CDB0) --------------------------------------------------------
int sub_57F2A0(float2* a1, int a2, int a3);
char sub_57F1D0(float2* a1);
int sub_57CDB0(int2* a1, float* a2, float2* a3) {
	int2* v3;   // esi
	char v4;    // bl
	int result; // eax
	float2* v6; // eax
	float2* v7; // eax
	float2* v8; // eax
	char v9;    // [esp+10h] [ebp+4h]

	v3 = a1;
	v9 = sub_57F2A0((float2*)a2, a1->field_0, a1->field_4);
	v4 = sub_57F1D0((float2*)a2 + 1);
	switch (sub_57B500(v3->field_0, v3->field_4, 64)) {
	case 0:
		if (v9 != 1 && v9) {
			v8 = a3;
			v8->field_0 = 0.70709997;
			v8->field_4 = 0.70709997;
			return 1;
		}
		a3->field_0 = -0.70709997;
		a3->field_4 = -0.70709997;
		return 1;
	case 1:
		if (v9 == 1 || v9 == 2) {
			a3->field_0 = 0.70709997;
			a3->field_4 = -0.70709997;
			return 1;
		}
		a3->field_0 = -0.70709997;
		a3->field_4 = 0.70709997;
		return 1;
	case 2:
		switch (v9) {
		case 0:
			v8 = a3;
			a3->field_0 = -0.70709997;
			if (!(v4 & 2)) {
				a3->field_4 = -0.70709997;
				return 1;
			}
			v8->field_4 = 0.70709997;
			return 1;
		case 1:
			v6 = a3;
			if (!(v4 & 1)) {
				v6->field_0 = 0.70709997;
				v6->field_4 = -0.70709997;
				return 1;
			}
			a3->field_0 = -0.70709997;
			a3->field_4 = -0.70709997;
			result = 1;
			break;
		case 2:
			v8 = a3;
			a3->field_0 = 0.70709997;
			if (v4 & 1) {
				v8->field_4 = 0.70709997;
				return 1;
			}
			a3->field_4 = -0.70709997;
			return 1;
		case 3:
			v8 = a3;
			if (!(v4 & 4)) {
				v8->field_0 = 0.70709997;
				v8->field_4 = 0.70709997;
				return 1;
			}
			a3->field_0 = -0.70709997;
			a3->field_4 = 0.70709997;
			return 1;
		default:
			return 1;
		}
		return result;
	case 3:
		if (!v9) {
			v8 = a3;
			a3->field_0 = -0.70709997;
			if (!(v4 & 2)) {
				a3->field_4 = -0.70709997;
				return 1;
			}
			v8->field_4 = 0.70709997;
			return 1;
		}
		if (v9 != 1) {
			v8 = a3;
			v8->field_0 = 0.70709997;
			v8->field_4 = 0.70709997;
			return 1;
		}
		v7 = a3;
		if (!(v4 & 1)) {
			v7->field_0 = 0.70709997;
			v7->field_4 = -0.70709997;
			return 1;
		}
		v7->field_0 = -0.70709997;
		v7->field_4 = -0.70709997;
		return 1;
	case 4:
		if (v9 == 1) {
			v7 = a3;
			if (!(v4 & 1)) {
				v7->field_0 = 0.70709997;
				v7->field_4 = -0.70709997;
				return 1;
			}
			v7->field_0 = -0.70709997;
			v7->field_4 = -0.70709997;
			return 1;
		}
		if (v9 == 2) {
			v8 = a3;
			a3->field_0 = 0.70709997;
			if (v4 & 1) {
				v8->field_4 = 0.70709997;
				return 1;
			}
			a3->field_4 = -0.70709997;
			return 1;
		}
		a3->field_0 = -0.70709997;
		a3->field_4 = 0.70709997;
		return 1;
	case 5:
		if (v9 == 2) {
			v8 = a3;
			a3->field_0 = 0.70709997;
			if (v4 & 1) {
				v8->field_4 = 0.70709997;
				return 1;
			}
			a3->field_4 = -0.70709997;
			return 1;
		}
		if (v9 == 3) {
			v8 = a3;
			if (!(v4 & 4)) {
				v8->field_0 = 0.70709997;
				v8->field_4 = 0.70709997;
				return 1;
			}
			a3->field_0 = -0.70709997;
			a3->field_4 = 0.70709997;
			return 1;
		}
		a3->field_0 = -0.70709997;
		a3->field_4 = -0.70709997;
		return 1;
	case 6:
		if (!v9) {
			v8 = a3;
			a3->field_0 = -0.70709997;
			if (!(v4 & 2)) {
				a3->field_4 = -0.70709997;
				return 1;
			}
			v8->field_4 = 0.70709997;
			return 1;
		}
		if (v9 == 3) {
			v8 = a3;
			if (!(v4 & 4)) {
				v8->field_0 = 0.70709997;
				v8->field_4 = 0.70709997;
				return 1;
			}
			a3->field_0 = -0.70709997;
			a3->field_4 = 0.70709997;
			return 1;
		}
		v6 = a3;
		v6->field_0 = 0.70709997;
		v6->field_4 = -0.70709997;
		return 1;
	case 7:
		if (v9 == 1) {
			v7 = a3;
			if (!(v4 & 1)) {
				v7->field_0 = 0.70709997;
				v7->field_4 = -0.70709997;
				return 1;
			}
			v7->field_0 = -0.70709997;
			v7->field_4 = -0.70709997;
			return 1;
		}
		v8 = a3;
		if (v4 & 1) {
			v8->field_0 = 0.70709997;
			v8->field_4 = 0.70709997;
			return 1;
		}
		a3->field_0 = -0.70709997;
		a3->field_4 = 0.70709997;
		return 1;
	case 8:
		if (v9 == 2) {
			v8 = a3;
			a3->field_0 = 0.70709997;
			if (v4 & 1) {
				v8->field_4 = 0.70709997;
				return 1;
			}
			a3->field_4 = -0.70709997;
			return 1;
		} else {
			v8 = a3;
			a3->field_0 = -0.70709997;
			if (!(v4 & 1)) {
				v8->field_4 = 0.70709997;
				return 1;
			}
			a3->field_4 = -0.70709997;
			return 1;
		}
	case 9:
		if (v9 == 3) {
			v8 = a3;
			if (!(v4 & 4)) {
				v8->field_0 = 0.70709997;
				v8->field_4 = 0.70709997;
				return 1;
			}
			a3->field_0 = -0.70709997;
			a3->field_4 = 0.70709997;
			return 1;
		} else {
			v7 = a3;
			if (v4 & 4) {
				v7->field_0 = 0.70709997;
				v7->field_4 = -0.70709997;
				return 1;
			} else {
				v7->field_0 = -0.70709997;
				v7->field_4 = -0.70709997;
				return 1;
			}
		}
	case 0xA:
		if (v9) {
			if (v4 & 2) {
				a3->field_0 = 0.70709997;
				a3->field_4 = -0.70709997;
				return 1;
			}
			v8 = a3;
			v8->field_0 = 0.70709997;
			v8->field_4 = 0.70709997;
			return 1;
		} else {
			v8 = a3;
			a3->field_0 = -0.70709997;
			if (!(v4 & 2)) {
				a3->field_4 = -0.70709997;
				return 1;
			}
			v8->field_4 = 0.70709997;
			return 1;
		}
	default:
		return 0;
	}
}

//----- (0057F1D0) --------------------------------------------------------
char sub_57F1D0(float2* a1) {
	char v1;          // bl
	int v2;           // edi
	double v3;        // st7
	double v4;        // st6
	double v5;        // st7
	unsigned char v7; // [esp+Ch] [ebp-4h]

	v1 = 0;
	v2 = nox_float2int(a1->field_0);
	v7 = nox_float2int(a1->field_4) % 23;
	v3 = (double)(unsigned char)(v2 % 23);
	if (v3 >= 11.5) {
		v4 = (double)v7;
		if (v4 >= 11.5) {
			v1 = 4;
		}
		if (v4 <= 11.5) {
			v1 |= 1u;
		}
	}
	if (v3 <= 11.5) {
		v5 = (double)v7;
		if (v5 >= 11.5) {
			v1 |= 8u;
		}
		if (v5 <= 11.5) {
			v1 |= 2u;
		}
	}
	return v1;
}

//----- (0057F2A0) --------------------------------------------------------
int sub_57F2A0(float2* a1, int a2, int a3) {
	int v3;     // esi
	int v4;     // eax
	int result; // eax
	int v6;     // eax
	float v7;   // [esp+0h] [ebp-Ch]
	float v8;   // [esp+0h] [ebp-Ch]

	v7 = a1->field_0 - (double)(23 * a2);
	v3 = nox_float2int(v7);
	v8 = a1->field_4 - (double)(23 * a3);
	v4 = nox_float2int(v8);
	if (v3 <= v4) {
		LOBYTE(v4) = 22 - v3 <= v4;
		v6 = v4 - 1;
		LOBYTE(v6) = v6 & 0xFD;
		result = v6 + 3;
	} else {
		LOBYTE(v4) = 22 - v3 <= v4;
		result = v4 + 1;
	}
	return result;
}

void nullsub_10(uint32_t a1) {}

int sub_42CC50(void** this) { return sub_42C770(this); }

int nox_xxx_j_inventoryNameSignInit_467460(void) { return nox_xxx_inventoryNameSignInit_4671E0(); }

int nullsub_8(int a1, int a2) { return 0; }
void nullsub_28(uint32_t a1) {}
void nullsub_30(uint32_t a1) {}
void nullsub_29(void) {}
void nullsub_35(uint32_t a1, uint32_t a2) {}
void nullsub_24(uint32_t a1) {}
void nullsub_31(uint32_t a1) {}
void nullsub_9(uint32_t a1) {}
void nox_xxx_j_allocHitArray_511840(void) { nox_xxx_allocHitArray_5486D0(); }
