#include <math.h>

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
#include "common__net_list.h"
#include "common__random.h"
#include "common__system__team.h"
#include "operators.h"
#include "server__ability__ability.h"
#include "server__magic__plyrspel.h"
#include "server__object__health.h"
#include "server__system__server.h"

#include "common__gamemech__pausefx.h"

#include "client__gui__window.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__crypt.h"
#include "common__magic__speltree.h"
#include "server__script__script.h"
#include "server__script__activator.h"

extern uint32_t nox_xxx_respawnAllow_587000_205200;
extern uint32_t dword_5d4594_1567960;
extern uint32_t dword_5d4594_1568280;
extern uint32_t dword_5d4594_1568288;
extern uint32_t dword_5d4594_1563320;
extern uint32_t dword_5d4594_1567988;
extern uint32_t dword_5d4594_1565628;
extern uint32_t dword_5d4594_1565632;
extern uint32_t dword_5d4594_1565520;
extern uint32_t dword_5d4594_1565516;
extern uint32_t dword_5d4594_1567928;
extern void* nox_alloc_respawn_1568020;
extern uint32_t dword_5d4594_1565616;

extern uint32_t dword_5d4594_2649712;
extern uint64_t qword_581450_10176;
extern uint64_t qword_581450_10256;
extern uint64_t qword_5d4594_1567940;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_1568024;
extern uint32_t dword_5d4594_1565512;
extern uint32_t dword_5d4594_2650652;
extern unsigned int gameex_flags;

//----- (004E41B0) --------------------------------------------------------
int sub_4E41B0(char* a1) {
	FILE* v1; // edi

	v1 = nox_fs_open_text(a1);
	if (v1) {
		do {
			while (1) {
				do {
					if (!nox_fs_fgets(v1, (char*)getMemAt(0x5D4594, 1563936), 1024)) {
						nox_fs_close(v1);
						return 1;
					}
				} while (!strtok((char*)getMemAt(0x5D4594, 1563936), "\r\t\n"));
				if (strcmp((const char*)getMemAt(0x5D4594, 1563936), "[Banned Users]")) {
					break;
				}
				if (!sub_4E42C0(v1)) {
					goto LABEL_11;
				}
			}
		} while (strcmp((const char*)getMemAt(0x5D4594, 1563936), "[Allowed Users]") || sub_4E4390(v1));
	LABEL_11:
		nox_fs_close(v1);
	}
	return 0;
}

//----- (004E42C0) --------------------------------------------------------
int sub_4E42C0(FILE* a1) {
	char* v1;       // eax
	char* v2;       // eax
	wchar2_t v4[26]; // [esp+Ch] [ebp-34h]

	while (1) {
		if (!nox_fs_fgets(a1, (char*)getMemAt(0x5D4594, 1563936), 1024)) {
			return 1;
		}
		v1 = strtok((char*)getMemAt(0x5D4594, 1563936), "\r\t\n");
		if (!v1) {
			return 1;
		}
		nox_swprintf(v4, L"%S", v1);
		if (!nox_fs_fgets(a1, (char*)getMemAt(0x5D4594, 1563936), 1024)) {
			return 1;
		}
		v2 = strtok((char*)getMemAt(0x5D4594, 1563936), "\r\t\n");
		if (!v2) {
			break;
		}
		if (!strcmp(v2, "0")) {
			sub_416770(0, v4, 0);
		} else {
			sub_416770(0, v4, v2);
		}
	}
	return 0;
}

//----- (004E4390) --------------------------------------------------------
int sub_4E4390(FILE* a1) {
	char* v1;       // eax
	wchar2_t v3[26]; // [esp+4h] [ebp-34h]

	while (nox_fs_fgets(a1, (char*)getMemAt(0x5D4594, 1563936), 1024)) {
		v1 = strtok((char*)getMemAt(0x5D4594, 1563936), "\r\t\n");
		if (!v1) {
			break;
		}
		nox_swprintf(v3, L"%S", v1);
		sub_4168A0(v3);
	}
	return 1;
}

//----- (004E43F0) --------------------------------------------------------
FILE* sub_4E43F0(char* a1) {
	FILE* result; // eax
	FILE* v2;     // edi
	int* i;       // esi
	int* j;       // esi
	char name[26];

	result = nox_fs_create_text(a1);
	v2 = result;
	if (result) {
		nox_fs_fprintf(result, "%s\n", getMemAt(0x587000, 202212));
		for (i = sub_416900(); i; i = sub_416910(i)) {
			if (!*((uint64_t*)i + 8)) {
				nox_sprintf(name, "%S", i + 3);
				nox_fs_fprintf(v2, "%s\n", name);
				if (*((uint8_t*)i + 72)) {
					nox_fs_fprintf(v2, "%s\n", i + 18);
				} else {
					nox_fs_fprintf(v2, "0\n");
				}
			}
		}
		nox_fs_fprintf(v2, "\n%s\n", getMemAt(0x587000, 202228));
		for (j = sub_4168E0(); j; j = sub_4168F0(j)) {
			nox_sprintf(name, "%S", j + 3);
		nox_fs_fprintf(v2, "%s\n", name);
		}
		nox_fs_close(v2);
		result = (FILE*)1;
	}
	return result;
}

int sub_50B510();
//----- (004EC520) --------------------------------------------------------
int nox_xxx_unitsHaveSameTeam_4EC520(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	int v2; // edi
	int v3; // esi

	v2 = a1;
	if (a1 && a2) {
		while (2) {
			v3 = a2;
			do {
				if (nox_xxx_servCompareTeams_419150(v2 + 48, v3 + 48) || v2 == v3) {
					return 1;
				}
				v3 = *(uint32_t*)(v3 + 508);
			} while (v3);
			v2 = *(uint32_t*)(v2 + 508);
			if (v2) {
				continue;
			}
			break;
		}
	}
	return 0;
}

//----- (004EC5B0) --------------------------------------------------------
void sub_4EC5B0() {
	dword_5d4594_1568024 = 0;
	nox_alloc_class_free_all(*(uint32_t**)&nox_alloc_respawn_1568020);
	nox_xxx_respawnAllow_587000_205200 = 1;
}

//----- (004EC5E0) --------------------------------------------------------
uint32_t* nox_xxx_respawnAdd_4EC5E0(nox_object_t* a1p) {
	int a1 = a1p;
	uint32_t* result;  // eax
	uint32_t* v2;      // ebx
	unsigned short v3; // cx
	uint8_t* v4;       // ebp

	result = *(uint32_t**)&nox_xxx_respawnAllow_587000_205200;
	if (nox_xxx_respawnAllow_587000_205200) {
		result = nox_alloc_class_new_obj_zero(*(uint32_t**)&nox_alloc_respawn_1568020);
		v2 = result;
		if (result) {
			v3 = *(uint16_t*)(a1 + 4);
			result[1] = a1;
			*result = v3;
			result[2] = *(uint32_t*)(a1 + 56);
			result[3] = *(uint32_t*)(a1 + 60);
			*((uint16_t*)result + 8) = *(uint16_t*)(a1 + 124);
			if (*(uint32_t*)(a1 + 8) & 0x13001000) {
				memcpy(result + 7, *(const void**)(a1 + 692), 0x14u);
			}
			if (*(uint32_t*)(a1 + 8) & 0x1000000 && nox_xxx_weaponInventoryEquipFlags_415820(a1) & 0x82) {
				v4 = *(uint8_t**)(a1 + 736);
				*((uint8_t*)v2 + 48) = v4[1];
				*((uint8_t*)v2 + 49) = *v4;
			}
			v2[14] = 0;
			v2[13] = dword_5d4594_1568024;
			result = *(uint32_t**)&dword_5d4594_1568024;
			if (dword_5d4594_1568024) {
				*(uint32_t*)(dword_5d4594_1568024 + 56) = v2;
			}
			dword_5d4594_1568024 = v2;
		}
	}
	return result;
}

//----- (004EC6A0) --------------------------------------------------------
void sub_4EC6A0(int a1) {
	int v1;       // eax
	uint64_t* v2; // ecx
	int v3;       // eax
	int v4;       // ecx
	int v5;       // ecx

	v1 = dword_5d4594_1568024;
	if (*(uint32_t*)(dword_5d4594_1568024 + 4) == a1) {
		v2 = *(uint64_t**)&dword_5d4594_1568024;
		dword_5d4594_1568024 = *(uint32_t*)(dword_5d4594_1568024 + 52);
		v3 = *(uint32_t*)(v1 + 52);
		if (v3) {
			*(uint32_t*)(v3 + 56) = 0;
		}
		nox_alloc_class_free_obj_first(*(unsigned int**)&nox_alloc_respawn_1568020, v2);
	} else if (dword_5d4594_1568024) {
		while (*(uint32_t*)(v1 + 4) != a1) {
			v1 = *(uint32_t*)(v1 + 52);
			if (!v1) {
				return;
			}
		}
		v4 = *(uint32_t*)(v1 + 56);
		if (v4) {
			*(uint32_t*)(v4 + 52) = *(uint32_t*)(v1 + 52);
		}
		v5 = *(uint32_t*)(v1 + 52);
		if (v5) {
			*(uint32_t*)(v5 + 56) = *(uint32_t*)(v1 + 56);
		}
		nox_alloc_class_free_obj_first(*(unsigned int**)&nox_alloc_respawn_1568020, (uint64_t*)v1);
	}
}

//----- (004ECA60) --------------------------------------------------------
int nox_xxx_allocItemRespawnArray_4ECA60() {
	nox_alloc_respawn_1568020 = nox_new_alloc_class("Respawn", 60, 384);
	return nox_alloc_respawn_1568020 != 0;
}

//----- (004ECA90) --------------------------------------------------------
void sub_4ECA90() { nox_free_alloc_class(*(void**)&nox_alloc_respawn_1568020); }

//----- (004ED050) --------------------------------------------------------
void sub_4ED050(int a1, int a2) {
	int v2; // eax
	int i;  // esi
	int v4; // edi

	LOWORD(v2) = *getMemU16Ptr(0x5D4594, 1568248);
	if (!*getMemU32Ptr(0x5D4594, 1568248)) {
		v2 = nox_xxx_getNameId_4E3AA0("Crown");
		*getMemU32Ptr(0x5D4594, 1568248) = v2;
	}
	for (i = *(uint32_t*)(a1 + 516); i; i = *(uint32_t*)(i + 512)) {
		LOWORD(v2) = *(uint16_t*)(i + 4);
		if ((unsigned short)v2 == *getMemU32Ptr(0x5D4594, 1568248)) {
			v4 = *(uint32_t*)(i + 748);
			LOWORD(v2) = nox_xxx_dropCrown_4ED5E0(a1, i, (int*)(a1 + 56));
			*(uint32_t*)(v4 + 4) = a2;
		}
	}
}

