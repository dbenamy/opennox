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

extern uint32_t dword_5d4594_527988;
extern uint32_t dword_5d4594_529332;
extern uint32_t dword_5d4594_528264;
extern uint32_t dword_5d4594_528260;
extern uint32_t dword_5d4594_531652;
extern uint32_t dword_5d4594_529336;
extern uint32_t dword_5d4594_531656;
extern uint32_t dword_5d4594_588084;
extern uint32_t dword_5d4594_528252;
extern uint32_t dword_587000_60044;
extern uint32_t dword_5d4594_531648;
extern uint32_t dword_5d4594_528256;
extern uint32_t dword_5d4594_534808;
extern uint32_t dword_5d4594_529340;
extern uint32_t dword_5d4594_2660652;
extern uint32_t dword_5d4594_529316;
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

int sub_41D1A0(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 527720) = a1;
	return result;
}

int sub_41D1B0() { return *getMemU32Ptr(0x5D4594, 527720); }

int sub_41D650() {
	char* v0;   // eax
	int result; // eax

	v0 = sub_41FA40();
	result = sub_41F800(v0);
	if (result) {
		result = 0;
	}
	*getMemU32Ptr(0x5D4594, 371700) = 1;
	return result;
}

int sub_41D670(char* a1) {
	int v1;      // eax
	short v3;    // [esp+2h] [ebp-4Ah]
	char v4[72]; // [esp+4h] [ebp-48h]

	if (!sub_420230(v4, &v3)) {
		return 0;
	}
	v1 = sub_4200E0();
	abort();
	return 0;
}

int sub_41D6C0() {
	int v0;         // ebx
	char* v1;       // ebp
	int v2;         // eax
	int v3;         // ebp
	char* v4;       // ebx
	int v5;         // eax
	int result;     // eax
	char* v7;       // ebp
	int v8;         // eax
	short v9;       // [esp+12h] [ebp-45Ah]
	char* v10;      // [esp+14h] [ebp-458h]
	int v11;        // [esp+18h] [ebp-454h]
	int v12;        // [esp+1Ch] [ebp-450h]
	int v13;        // [esp+20h] [ebp-44Ch]
	char v14[72];   // [esp+24h] [ebp-448h]
	char v15[1024]; // [esp+6Ch] [ebp-400h]

	v0 = 1;
	memset(v15, 0, sizeof(v15));
	v11 = 1;
	v12 = 1;
	if (!sub_420230(v14, &v9)) {
		return 0;
	}
	if ((unsigned int)nox_common_playerInfoCount_416F40() > 25) {
		v1 = nox_common_playerInfoGetFirst_416EA0();
		v13 = 25;
		while (1) {
			if (v1[2096]) {
				if (v0) {
					strcat(v15, v1 + 2096);
					v0 = 0;
				} else {
					*(uint16_t*)&v15[strlen(v15)] = *getMemU16Ptr(0x587000, 58112);
					strcat(v15, v1 + 2096);
				}
			}
			v10 = nox_common_playerInfoGetNext_416EE0((int)v1);
			--v13;
			if (!v13) {
				break;
			}
			v1 = v10;
		}
		if (!v0) {
			v2 = sub_4200E0();
			abort();
		}
		v3 = 1;
		if (!v10) {
			return v11 && v12 == 1;
		}
		v15[0] = 0;
		v4 = nox_common_playerInfoGetNext_416EE0((int)v10);
		if (!v4) {
			return v11 && v12 == 1;
		}
		do {
			if (v4[2096]) {
				if (v3) {
					strcat(v15, v4 + 2096);
					v3 = 0;
				} else {
					*(uint16_t*)&v15[strlen(v15)] = *getMemU16Ptr(0x587000, 58116);
					strcat(v15, v4 + 2096);
				}
			}
			v4 = nox_common_playerInfoGetNext_416EE0((int)v4);
		} while (v4);
		if (!v3) {
			v5 = sub_4200E0();
			abort();
		} else {
			result = v12;
		}
		return v11 && result == 1;
	}
	v7 = nox_common_playerInfoGetFirst_416EA0();
	if (v7) {
		do {
			if (v7[2096]) {
				if (v0) {
					strcat(v15, v7 + 2096);
					v0 = 0;
				} else {
					*(uint16_t*)&v15[strlen(v15)] = *getMemU16Ptr(0x587000, 58120);
					strcat(v15, v7 + 2096);
				}
			}
			v7 = nox_common_playerInfoGetNext_416EE0((int)v7);
		} while (v7);
		if (!v0) {
			v8 = sub_4200E0();
			abort();
		}
	}
	return v11;
}

int sub_41DA10(int a1) {
	uint16_t* v1;   // edi
	int result = 0; // eax

	v1 = *(uint16_t**)getMemAt(0x587000, 58132 + 16 * a1);
	if (v1) {
		result = 0;
		memset(v1, 0, 0x2Cu);
		v1[22] = 0;
		*getMemU32Ptr(0x587000, 58136 + 16 * a1) = 0;
	}
	return result;
}

int sub_41DA70(int a1, short a2) {
	int result; // eax
	int v3;     // ecx
	int v4;     // edx

	result = 16 * a1;
	v3 = *getMemU32Ptr(0x587000, 58132 + 16 * a1);
	if (v3 && (v4 = *getMemU32Ptr(0x587000, 58136 + 16 * a1), v4 < 23)) {
		*(uint16_t*)(v3 + 2 * v4) = a2;
		++*getMemU32Ptr(0x587000, 58136 + 16 * a1);
	} else {
		result = sub_41E300(11);
		if (result) {
			dword_5d4594_2660652 = 0;
		}
	}
	return result;
}

int sub_41E2F0() { return dword_5d4594_527988; }

int sub_41E370() {
	int result; // eax

	result = 0;
	dword_5d4594_528252 = 0;
	dword_5d4594_528256 = 0;
	dword_5d4594_528260 = 0;
	dword_5d4594_528264 = 0;
	return result;
}

int nox_xxx_reconAttempt_41E390() {
	int result; // eax

	if (gameFrame() - dword_5d4594_528264 <= (unsigned int)(3600 * gameFPS())) {
		result = dword_5d4594_528252;
		if (dword_5d4594_528252) {
			result = dword_5d4594_528256;
			if (!dword_5d4594_528256) {
				nox_xxx_networkLog_printf_413D30("RECON: Attempting to re-login");
				sub_40E090();
				result = nox_xxx_officialStringCmp_41FDE0();
				if (result == 1) {
					dword_5d4594_528256 = 1;
				} else {
					result = sub_41E470();
				}
			}
		}
	} else {
		sub_41E370();
		result = sub_41E4B0(1);
	}
	return result;
}

void nox_xxx_reconStart_41E400() {
	if (dword_5d4594_528252 != 1 && dword_5d4594_528256 != 1) {
		if (!dword_5d4594_528260) {
			if (!dword_5d4594_528264) {
				nox_xxx_networkLog_printf_413D30("RECON: Starting reconnection process frame (%d)",
												 gameFrame());
				dword_5d4594_528252 = 1;
				dword_5d4594_528256 = 0;
				dword_5d4594_528264 = gameFrame();
				dword_5d4594_528260 = gameFrame() + 120 * gameFPS();
			}
		}
	}
}

int sub_41E470() {
	int result; // eax

	nox_xxx_networkLog_printf_413D30("RECON: TryReconnectAgain called on frame (%d)", gameFrame());
	dword_5d4594_528256 = 0;
	result = gameFrame() + 120 * gameFPS();
	dword_5d4594_528260 = gameFrame() + 120 * gameFPS();
	return result;
}

int sub_41E4B0(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 528268) = a1;
	return result;
}

int sub_41EC30() {
	uint32_t* v0; // ebx
	int v1;       // esi
	int v2;       // edi
	int result;   // eax
	int v4;       // eax

	v0 = *(uint32_t**)&dword_5d4594_529316;
	dword_5d4594_529332 = 0;
	dword_5d4594_529336 = 0;
	if (dword_5d4594_529316) {
		do {
			v1 = v0[7];
			if (v1) {
				do {
					v2 = *(uint32_t*)(v1 + 40);
					if (*(uint32_t*)v1) {
						free(*(void**)v1);
					}
					free((void*)v1);
					v1 = v2;
				} while (v2);
			}
			v0[7] = 0;
			v0[8] = 0;
			v0[6] = 0;
			v0 = (uint32_t*)v0[9];
		} while (v0);
	}
	result = sub_41E2F0();
	if (result == 7) {
		v4 = sub_41E2F0();
		result = sub_41DA70(v4, 10);
	}
	dword_5d4594_529340 = 0;
	return result;
}

int sub_41F4B0() {
	int v0;     // esi
	int v1;     // edi
	int result; // eax
	int v3;     // eax

	v0 = dword_5d4594_531648;
	if (dword_5d4594_531648) {
		do {
			v1 = *(uint32_t*)(v0 + 20);
			if (*(uint32_t*)v0) {
				free(*(void**)v0);
			}
			free((void*)v0);
			v0 = v1;
		} while (v1);
	}
	dword_5d4594_531648 = 0;
	dword_5d4594_531652 = 0;
	dword_5d4594_531656 = 0;
	result = sub_41E2F0();
	if (result == 7) {
		v3 = sub_41E2F0();
		result = sub_41DA70(v3, 11);
	}
	return result;
}

uint32_t* sub_41F790(const char* a1) {
	uint32_t* v1; // edi

	v1 = *(uint32_t**)&dword_5d4594_531648;
	if (!a1 || !dword_5d4594_531648) {
		return 0;
	}
	while (strcmp((const char*)(*v1 + 36), a1)) {
		v1 = (uint32_t*)v1[5];
		if (!v1) {
			return 0;
		}
	}
	return v1;
}

int sub_41F800(const char* a1) {
	int* v1;    // eax
	int result; // eax

	v1 = sub_41F790(a1);
	if (v1) {
		result = *v1;
	} else {
		result = 0;
	}
	return result;
}

char* sub_41FA40() { return (char*)getMemAt(0x5D4594, 534756); }

void sub_41FA50(const char* a1) {
	if (a1) {
		strcpy((char*)getMemAt(0x5D4594, 534756), a1);
	}
}

int sub_41FBE0(uint32_t* a1, uint32_t* a2) {
	int result; // eax

	if (*(int*)&dword_587000_60044 == -1) {
		return 0;
	}
	*a1 = getMemAt(0x5D4594, 531660 + 24 * dword_587000_60044);
	result = 1;
	*a2 = getMemAt(0x5D4594, 531670 + 24 * dword_587000_60044);
	return result;
}

int nox_xxx_officialStringCmp_41FDE0() {
	int v0;         // ebx
	size_t v1;      // eax
	const char* v3; // [esp+8h] [ebp-8h]
	const char* v4; // [esp+Ch] [ebp-4h]

	memset(getMemAt(0x85B3FC, 10308), 0, 0x6Cu);
	if (sub_41FBE0(&v3, &v4) != 1) {
		return 0;
	}
	v0 = sub_4207E0();
	if (v0) {
		strcpy((char*)(v0 + 228), v3);
		strcpy((char*)(v0 + 238), v4);
		sub_41FA50(v3);
		v1 = strcspn((const char*)(v0 + 24), ":");
		if (!strncmp("Official", (const char*)(v1 + v0 + 25), 8u)) {
			nox_xxx_setGameFlags_40A4D0(0x1000000);
		} else {
			nox_common_gameFlags_unset_40A540(0x1000000);
		}
	}
	return 1;
}

int sub_4200E0() { return *getMemU32Ptr(0x587000, 60072); }

int sub_420100() { return *getMemU32Ptr(0x587000, 60072) >> 8; }

int sub_420230(char* a1, uint16_t* a2) {
	int v2;       // ebx
	char* v3;     // eax
	char* v4;     // eax
	char v6[128]; // [esp+10h] [ebp-80h]

	if (!a1) {
		return 0;
	}
	if (!a2) {
		return 0;
	}
	v2 = dword_5d4594_534808;
	if (!dword_5d4594_534808) {
		return 0;
	}
	while (1) {
		if (!nox_strcmpi((const char*)(v2 + 95), "LAD")) {
			strcpy(v6, (const char*)(v2 + 100));
			*a1 = 0;
			*a2 = 0;
			strtok(v6, ";");
			v3 = strtok(0, ";");
			if (v3) {
				strcpy(a1, v3);
			}
			v4 = strtok(0, ";");
			if (v4) {
				*a2 = atoi(v4);
			}
			if (*a1 && *a2) {
				break;
			}
		}
		v2 = *(uint32_t*)(v2 + 20);
		if (!v2) {
			return 0;
		}
	}
	return 1;
}

int sub_420360(char* a1, uint16_t* a2) {
	int v2;       // ebx
	char* v3;     // eax
	char* v4;     // eax
	char* v5;     // eax
	char v7[128]; // [esp+10h] [ebp-80h]

	if (!a1) {
		return 0;
	}
	if (!a2) {
		return 0;
	}
	*a1 = 0;
	*a2 = 0;
	v2 = dword_5d4594_534808;
	if (!dword_5d4594_534808) {
		return 0;
	}
	while (1) {
		if (!nox_strcmpi((const char*)(v2 + 95), "GAM")) {
			if (nox_common_gameFlags_check_40A5C0(4096)) {
				if (!nox_strcmpi((const char*)(v2 + 24), "Quest gameres server")) {
					strcpy(v7, (const char*)(v2 + 100));
					*a1 = 0;
					*a2 = 0;
					strtok(v7, ";");
					v3 = strtok(0, ";");
					if (v3) {
						strcpy(a1, v3);
					}
					v4 = strtok(0, ";");
					if (v4) {
						*a2 = atoi(v4);
					}
				}
			} else {
				if (!nox_common_gameFlags_check_40A5C0(0x2000)) {
					return 0;
				}
				if (!nox_strcmpi((const char*)(v2 + 24), "Gameres server")) {
					strcpy(v7, (const char*)(v2 + 100));
					*a1 = 0;
					*a2 = 0;
					strtok(v7, ";");
					v5 = strtok(0, ";");
					if (v5) {
						strcpy(a1, v5);
					}
					v4 = strtok(0, ";");
					if (v4) {
						*a2 = atoi(v4);
					}
				}
			}
			if (*a1 && *a2) {
				return 1;
			}
		}
		v2 = *(uint32_t*)(v2 + 20);
		if (!v2) {
			return 0;
		}
	}
}

int sub_4207E0() { return *getMemU32Ptr(0x5D4594, 534812); }

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
