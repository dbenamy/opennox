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
