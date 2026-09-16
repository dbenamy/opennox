#include <math.h>

#include "compat.h"

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
#include "common__log.h"
#include "common__net_list.h"
#include "common__random.h"
#include "common__system__team.h"
#include "server__mapgen__generate__populate.h"
#include "server__system__server.h"

#include "client__gui__chathelp.h"
#include "client__gui__conntype.h"
#include "client__gui__servopts__guiserv.h"
#include "client__gui__window.h"

#include "client__drawable__update__cloud.h"
#include "client__drawable__update__dball.h"

#include "MixPatch.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__crypt.h"
#include "common__magic__speltree.h"
#include "defs.h"
#include "operators.h"
#include "server__script__builtin.h"

extern uint32_t dword_5d4594_3835368;
extern uint32_t dword_5d4594_3835360;
extern uint32_t dword_5d4594_3835364;
extern uint32_t dword_5d4594_1556128;
extern uint32_t dword_5d4594_1556136;
extern uint32_t dword_5d4594_3835372;
extern uint32_t dword_5d4594_1556144;
extern uint32_t dword_5d4594_1523040;
extern uint32_t dword_5d4594_1563276;
extern uint32_t dword_5d4594_3835392;
extern uint32_t nox_server_sendMotd_108752;
extern uint32_t dword_5d4594_1556856;
extern uint32_t dword_5d4594_1548480;
extern uint32_t dword_5d4594_1523048;
extern uint32_t dword_5d4594_1523044;
extern uint32_t dword_5d4594_1523032;
extern uint32_t nox_server_sanctuaryHelp_54276;
extern uint32_t dword_5d4594_3835312;
extern uint32_t dword_5d4594_3835388;
extern uint32_t dword_5d4594_3835348;
extern uint32_t dword_5d4594_3835352;
extern uint32_t dword_5d4594_1523036;
extern uint32_t dword_5d4594_1548700;
extern uint32_t dword_5d4594_3835356;

extern uint32_t nox_server_connectionType_3596;
extern uint32_t dword_5d4594_1550916;
extern uint32_t dword_5d4594_2649712;
extern uint32_t dword_5d4594_3835396;
extern uint32_t dword_5d4594_1523024;
extern uint32_t dword_5d4594_1523028;
extern uint32_t dword_5d4594_1548476;

extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern int ptr_5D4594_2650668_cap;

nox_list_item_t nox_common_maplist = {0};

//----- (004CEBA0) --------------------------------------------------------
int sub_4CEBA0(int a1, char* a2) {
	uint32_t* v2; // eax
	uint32_t* v3; // edi
	char* v4;     // ebx
	uint32_t* v5; // esi
	uint32_t* v6; // ebp
	uint32_t* v8; // [esp+10h] [ebp-4h]
	char* v9;     // [esp+18h] [ebp+4h]

	dword_5d4594_1523024 = nox_new_window_from_file("rulelist.wnd", sub_4CF060);
	dword_5d4594_1523028 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10170);
	dword_5d4594_1523032 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10171);
	dword_5d4594_1523036 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10172);
	dword_5d4594_1523040 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10173);
	dword_5d4594_1523044 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10174);
	dword_5d4594_1523048 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10175);
	v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10176);
	nox_xxx_wndSetDrawFn_46B340((int)v2, sub_4CEED0);
	sub_46B120(*(uint32_t**)&dword_5d4594_1523024, a1);
	v3 = *(uint32_t**)(dword_5d4594_1523028 + 32);
	v9 = nox_xxx_gLoadImg_42F970("UISlider");
	v4 = nox_xxx_gLoadImg_42F970("UISliderLit");
	v5 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10179);
	v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10177);
	v8 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, 10178);
	*(uint32_t*)(v5[100] + 8) = 16;
	*(uint32_t*)(v5[100] + 12) = 10;
	sub_4B5700((int)v5, 0, 0, (int)v9, (int)v4, (int)v4);
	nox_xxx_wnd_46B280((int)v5, *(int*)&dword_5d4594_1523028);
	nox_xxx_wnd_46B280((int)v6, *(int*)&dword_5d4594_1523028);
	nox_xxx_wnd_46B280((int)v8, *(int*)&dword_5d4594_1523028);
	v3[9] = v5;
	v3[7] = v6;
	v3[8] = v8;
	sub_4CED40(a2);
	return dword_5d4594_1523024;
}

//----- (004CED40) --------------------------------------------------------
void* sub_4CED40(char* a1) {
	HANDLE result;                         // eax
	HANDLE v2;                             // ebp
	struct _WIN32_FIND_DATAA FindFileData; // [esp+8h] [ebp-248h]
	char FileName[64];                     // [esp+148h] [ebp-108h]
	wchar2_t v5[100];                       // [esp+188h] [ebp-C8h]

	nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16399, 0, 0);
	nox_sprintf(FileName, "maps\\%s\\*.rul", a1);
	result = FindFirstFileA(FileName, &FindFileData);
	v2 = result;
	if (result != (HANDLE)-1) {
		FindFileData.cFileName[strlen(FindFileData.cAlternateFileName) + 256] = 0;
		if (nox_strcmpi(a1, FindFileData.cAlternateFileName) && nox_strcmpi("user", FindFileData.cAlternateFileName)) {
			nox_swprintf(v5, L"%S", FindFileData.cAlternateFileName);
			nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16397, (int)v5, -1);
		}
		while (FindNextFileA(v2, &FindFileData)) {
			FindFileData.cFileName[strlen(FindFileData.cAlternateFileName) + 256] = 0;
			if (nox_strcmpi(a1, FindFileData.cAlternateFileName)) {
				if (nox_strcmpi("user", FindFileData.cAlternateFileName)) {
					nox_swprintf(v5, L"%S", FindFileData.cAlternateFileName);
					nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16397, (int)v5, -1);
				}
			}
		}
		result = (HANDLE)FindClose(v2);
	}
	return result;
}

//----- (004CEED0) --------------------------------------------------------
int sub_4CEED0(int a1, int a2) {
	int v2;       // eax
	char v3;      // cl
	uint16_t* v4; // esi
	uint32_t* v5; // eax
	int xLeft;    // [esp+8h] [ebp-8h]
	int yTop;     // [esp+Ch] [ebp-4h]

	nox_client_wndGetPosition_46AA60((uint32_t*)a1, &xLeft, &yTop);
	if ((signed char)*(uint8_t*)(a1 + 4) >= 0) {
		if (*(uint32_t*)(a2 + 20) != 0x80000000) {
			nox_client_drawRectFilledOpaque_49CE30(xLeft, yTop, *(uint32_t*)(a1 + 8), *(uint32_t*)(a1 + 12));
		}
	} else {
		nox_client_drawImageAt_47D2C0(*(uint32_t*)(a2 + 24), xLeft, yTop);
	}
	v2 = nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16404, 0, 0);
	v3 = *(uint8_t*)(dword_5d4594_1523040 + 4);
	if (v2 < 0) {
		if (v3 & 8) {
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523040, 0);
		}
		if (*(uint8_t*)(dword_5d4594_1523044 + 4) & 8) {
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523044, 0);
		}
		if (*(uint8_t*)(dword_5d4594_1523048 + 4) & 8) {
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523048, 0);
		}
	} else {
		if (!(v3 & 8)) {
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523040, 1);
		}
		if (!(*(uint8_t*)(dword_5d4594_1523044 + 4) & 8) && !nox_common_gameFlags_check_40A5C0(49152)) {
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523044, 1);
		}
		if (!(*(uint8_t*)(dword_5d4594_1523048 + 4) & 8)) {
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523048, 1);
		}
	}
	v4 = (uint16_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1523032, 16413, 0, 0);
	v5 = (uint32_t*)nox_xxx_wndGetFocus_46B4F0();
	if (v5 && *v5 == 10171) {
		if (v4 && *v4) {
			if (!(*(uint8_t*)(dword_5d4594_1523036 + 4) & 8)) {
				nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523036, 1);
				return 1;
			}
		} else if (*(uint8_t*)(dword_5d4594_1523036 + 4) & 8) {
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523036, 0);
		}
	}
	return 1;
}

//----- (004CF060) --------------------------------------------------------
int sub_4CF060(int a1, unsigned int a2, int* a3, int a4) {
	uint32_t* v4;       // esi
	const char* v6;     // eax
	const char* v7;     // esi
	char* v8;           // eax
	int v9;             // esi
	int v10;            // eax
	int v11;            // eax
	char* v12;          // eax
	char* v13;          // eax
	char* v14;          // esi
	int v15;            // eax
	int v16;            // eax
	char* v17;          // eax
	char* v18;          // edi
	int v19;            // eax
	int v20;            // esi
	int v21;            // eax
	const wchar2_t* v22; // edi
	char* v23;          // eax
	char* v24;          // eax
	int v25;            // esi
	int v26;            // ebx
	const wchar2_t* v27; // eax
	char v28[16];       // [esp+Ch] [ebp-10h]

	if (a2 > 0x4007) {
		if (a2 == 16400) {
			nox_xxx_wndGetID_46B0A0(a3);
		}
		return 1;
	}
	if (a2 != 16391) {
		if (a2 != 23 && a2 == 16387) {
			v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1523024, a4);
			if (!v4) {
				return 0;
			}
			if ((unsigned short)a3 == 1) {
				return 0;
			}
			v6 = (const char*)nox_window_call_field_94((int)v4, 16413, 0, 0);
			v7 = v6;
			if (v6) {
				if (*v6) {
					atoi(v6);
					if (a4 == 10171) {
						v8 = sub_4165B0();
						if (!nox_strcmpi(v7, v8) || !nox_strcmpi(v7, "user")) {
							nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523036, 0);
							return 1;
						}
					}
				}
			}
		}
		return 1;
	}
	v9 = nox_xxx_wndGetID_46B0A0(a3);
	nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
	switch (v9) {
	case 10172:
		sub_416580();
		v22 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1523032, 16413, 0, 0);
		nox_sprintf(v28, "%S%s", v22, getMemAt(0x587000, 191640));
		v23 = sub_4165B0();
		sub_459AA0((int)v23);
		v24 = sub_4165B0();
		sub_57AAA0(v28, v24, 0);
		v25 = 0;
		v26 = *(uint32_t*)(dword_5d4594_1523028 + 32);
		if (*(short*)(v26 + 44) <= 0) {
			nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16397, (int)v22, -1);
			nox_window_call_field_94(*(int*)&dword_5d4594_1523032, 16414, (int)getMemAt(0x5D4594, 1523056), 0);
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523036, 0);
			return 1;
		}
		break;
	case 10173:
		sub_416580();
		v10 = nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16404, 0, 0);
		v11 = nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16406, v10, 0);
		nox_sprintf(v28, "%S%s", v11, getMemAt(0x587000, 191592));
		v12 = sub_4165B0();
		sub_459AA0((int)v12);
		v13 = sub_4165B0();
		sub_57AAA0(v28, v13, 0);
		nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16403, -1, 0);
		return 1;
	case 10174:
		sub_416580();
		v14 = sub_4165B0();
		v15 = nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16404, 0, 0);
		v16 = nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16406, v15, 0);
		nox_sprintf(v28, "%S%s", v16, getMemAt(0x587000, 191608));
		sub_57A1E0((int*)v14, v28, 0, 7, *((uint16_t*)v14 + 26));
		sub_453F70(v14 + 24);
		sub_4535E0((int*)v14 + 11);
		sub_4535F0(*((uint32_t*)v14 + 12));
		v17 = sub_4165B0();
		sub_459880((int)v17);
		nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16403, -1, 0);
		sub_459D50(1);
		return 1;
	case 10175:
		v18 = sub_4165B0();
		v19 = nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16404, 0, 0);
		v20 = v19;
		v21 = nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16406, v19, 0);
		nox_sprintf(v28, "%S%s", v21, getMemAt(0x587000, 191624));
		sub_57A9F0(v18, v28);
		nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16398, v20, 0);
		return 1;
	default:
		return 1;
	}
	while (1) {
		v27 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16406, v25, 0);
		if (!_nox_wcsicmp(v22, v27)) {
			break;
		}
		if (++v25 >= *(short*)(v26 + 44)) {
			nox_window_call_field_94(*(int*)&dword_5d4594_1523028, 16397, (int)v22, -1);
			nox_window_call_field_94(*(int*)&dword_5d4594_1523032, 16414, (int)getMemAt(0x5D4594, 1523056), 0);
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1523036, 0);
			return 1;
		}
	}
	nox_window_call_field_94(*(int*)&dword_5d4594_1523032, 16414, (int)getMemAt(0x5D4594, 1523052), 0);
	return 1;
}

//----- (004CF470) --------------------------------------------------------
int nox_xxx_mapValidateMB_4CF470(char* a1, int a2) {
	int v2;              // ebx
	int v4;              // [esp+8h] [ebp-408h]
	int v5;              // [esp+Ch] [ebp-404h]
	char FileName[1024]; // [esp+10h] [ebp-400h]

	v2 = 0;
	if (!a2) {
		return 6;
	}
	if (a1) {
		if (strchr(a1, '\\')) {
			strcpy(FileName, a1);
		} else {
			strcpy(FileName, "maps\\");
			strncat(FileName, a1, 1024-6);
			FileName[strlen(FileName)-4] = 0;
			*(uint16_t*)&FileName[strlen(FileName)] = *getMemU16Ptr(0x587000, 191672);
			strcat(FileName, a1);
		}
		if (nox_fs_access(FileName, 0) != -1) {
			v4 = 0;
			if (nox_fs_access(FileName, 2) == -1) {
				v2 = 1;
			}
			if (nox_xxx_cryptOpen_426910(FileName, 1, 19)) {
				v2 |= 2u;
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 4u);
				if (v5 != -86065425 && v5 == -86050098) {
					nox_xxx_fileCryptReadCrcMB_426C20(&v4, 4u);
					if (v4 == a2) {
						v2 |= 4u;
					}
				}
				nox_xxx_cryptClose_4269F0();
			}
		}
	}
	return v2;
}

//----- (004CFC30) --------------------------------------------------------
void nox_xxx_mapFindCrown_4CFC30() {
	int v0; // esi
	int v1; // edi

	if (!*getMemU32Ptr(0x5D4594, 1523076)) {
		*getMemU32Ptr(0x5D4594, 1523076) = nox_xxx_getNameId_4E3AA0("Crown");
	}
	v0 = nox_server_getFirstObject_4DA790();
	if (v0) {
		do {
			v1 = nox_server_getNextObject_4DA7A0(v0);
			if (*(unsigned short*)(v0 + 4) == *getMemU32Ptr(0x5D4594, 1523076)) {
				nox_xxx_delayedDeleteObject_4E5CC0(v0);
				sub_4EC6A0(v0);
			}
			v0 = v1;
		} while (v1);
	}
}

//----- (004CFDF0) --------------------------------------------------------
int sub_4CFDF0(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1523072) = a1;
	return result;
}

//----- (004CFE00) --------------------------------------------------------
int sub_4CFE00() { return *getMemU32Ptr(0x5D4594, 1523072); }

//----- (004CFFA0) --------------------------------------------------------
int nox_xxx_mapGetTypeMB_4CFFA0(void* a1) { return nox_mapToGameFlags_4CFF50(*(uint32_t*)((int)a1 + 1392)); }

//----- (004CFFC0) --------------------------------------------------------
int sub_4CFFC0(int a1) { return nox_mapToGameFlags_4CFF50(*(uint32_t*)(a1 + 28)); }

//----- (004D0010) --------------------------------------------------------
void* nox_objectTypeGetXfer(char* id);
int nox_xxx_interesting_xfer_4D0010(uint32_t* a1, int a2) {
	int i;          // eax
	uint32_t* v3;   // edi
	char* v4;       // eax
	int (*v5)(int); // eax
	int v6;         // esi
	int v7;         // eax
	int v8;         // ecx
	int v9;         // eax
	int v10;        // ecx
	int v11;        // esi
	int v12;        // eax
	int v13;        // ecx
	int v14;        // esi
	char* v15;      // eax
	int v16;        // ecx
	int v17;        // esi
	char* v18;      // eax
	int v19;        // esi
	uint32_t* v20;  // eax
	int v21;        // eax
	int v22;        // esi
	char* v23;      // eax

	for (i = nox_server_getFirstObjectUninited_4DA870(); i; i = nox_server_getNextObjectUninited_4DA880(i)) {
		*(uint32_t*)(i + 44) = *(uint32_t*)(i + 40);
		*(uint32_t*)(i + 40) = a2++;
	}
	v3 = (uint32_t*)nox_server_getFirstObjectUninited_4DA870();
	if (!v3) {
		return a2;
	}
	while (1) {
		v4 = (char*)nox_xxx_getUnitName_4E39D0((int)v3);
		v5 = nox_objectTypeGetXfer(v4);
		if (v5 == nox_xxx_XFerElevator_4F53D0) {
			v6 = v3[187];
			v7 = sub_4CFFE0(*(uint32_t*)(v6 + 8));
			if (!v7) {
				*(uint32_t*)(v6 + 8) = 0;
				*(uint32_t*)(v6 + 4) = 0;
			} else {
				v8 = *(uint32_t*)(v7 + 40);
				*(uint32_t*)(v6 + 4) = v7;
				*(uint32_t*)(v6 + 8) = v8;
			}
		} else if (v5 == nox_xxx_XFerElevatorShaft_4F54A0) {
			v6 = v3[187];
			v9 = sub_4CFFE0(*(uint32_t*)(v6 + 8));
			if (v9) {
				v10 = *(uint32_t*)(v9 + 40);
				*(uint32_t*)(v6 + 4) = v9;
				*(uint32_t*)(v6 + 8) = v10;
			} else {
				*(uint32_t*)(v6 + 8) = 0;
				*(uint32_t*)(v6 + 4) = 0;
			}
		} else if (v5 == nox_xxx_XFerTransporter_4F5300) {
			v11 = v3[187];
			v12 = sub_4CFFE0(*(uint32_t*)(v11 + 16));
			if (v12) {
				v13 = *(uint32_t*)(v12 + 40);
				*(uint32_t*)(v11 + 12) = v12;
				*(uint32_t*)(v11 + 16) = v13;
			} else {
				*(uint32_t*)(v11 + 16) = 0;
				*(uint32_t*)(v11 + 12) = 0;
			}
		} else if (v5 == nox_xxx_XFerHole_4F51D0) {
			v14 = v3[175];
			v15 = nox_xxx_mapGetWallSize_426A70();
			v16 = *(uint32_t*)(v14 + 12);
			*(uint32_t*)(v14 + 8) += *a1 - 23 * *(uint32_t*)v15;
			*(uint32_t*)(v14 + 12) = a1[1] - 23 * *((uint32_t*)v15 + 1) + v16;
		} else if (v5 == nox_xxx_XFerExit_4F4B90) {
			v17 = v3[175];
			v18 = nox_xxx_mapGetWallSize_426A70();
			*(float*)(v17 + 80) = (double)(int)(*a1 - 23 * *(uint32_t*)v18) + *(float*)(v17 + 80);
			*(float*)(v17 + 84) = (double)(int)(a1[1] - 23 * *((uint32_t*)v18 + 1)) + *(float*)(v17 + 84);
		} else if (v5 == nox_xxx_XFerMover_4F5730) {
			v19 = v3[187];
			v20 = (uint32_t*)sub_579C60(*(uint32_t*)(v19 + 8));
			if (v20) {
				*(uint32_t*)(v19 + 8) = *v20;
			} else {
				*(uint32_t*)(v19 + 8) = 0;
			}
			v21 = sub_4CFFE0(*(uint32_t*)(v19 + 32));
			if (v21) {
				*(uint32_t*)(v19 + 32) = *(uint32_t*)(v21 + 40);
			} else {
				*(uint32_t*)(v19 + 32) = 0;
			}
		} else if (v5 == nox_xxx_XFerGlyph_4F5890) {
			v22 = v3[173];
			v23 = nox_xxx_mapGetWallSize_426A70();
			*(float*)(v22 + 28) = (double)(int)(*a1 - 23 * *(uint32_t*)v23) + *(float*)(v22 + 28);
			*(float*)(v22 + 32) = (double)(int)(a1[1] - 23 * *((uint32_t*)v23 + 1)) + *(float*)(v22 + 32);
		}
		v3 = (uint32_t*)nox_server_getNextObjectUninited_4DA880((int)v3);
		if (!v3) {
			return a2;
		}
	}
}

//----- (004D11A0) --------------------------------------------------------
void sub_4D11A0() {
	if (!*getMemU32Ptr(0x5D4594, 1548504)) {
		nox_common_list_clear_425760(getMemAt(0x5D4594, 1548492));
		*getMemU32Ptr(0x5D4594, 1548504) = 1;
	}
}

//----- (004D11D0) --------------------------------------------------------
void sub_4D11D0() {
	int* result; // eax
	int* v1;     // esi
	int* v2;     // edi

	result = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1548492));
	v1 = result;
	if (result) {
		do {
			v2 = nox_common_list_getNextSafe_4258A0(v1);
			nox_common_list_remove_425920((uint32_t**)v1);
			free(v1);
			v1 = v2;
		} while (v2);
	}
}

//----- (004D1210) --------------------------------------------------------
void sub_4D1210(int a1) {
	void* result; // eax
	void* v2;     // esi
	uint32_t* v3; // eax

	result = (void*)sub_4D12A0(a1);
	if (!result) {
		result = nox_common_playerInfoFromNum_417090(a1);
		v2 = result;
		if (result) {
			v3 = calloc(1, 0x10u);
			v3[3] = v2;
			nox_common_list_append_4258E0((int)getMemAt(0x5D4594, 1548492), v3);
		}
	}
}

//----- (004D1250) --------------------------------------------------------
int* sub_4D1250(int a1) {
	int* result; // eax
	int* v2;     // esi

	result = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1548492));
	v2 = result;
	if (result) {
		while (*(unsigned char*)(v2[3] + 2064) != a1) {
			result = nox_common_list_getNextSafe_4258A0(v2);
			v2 = result;
			if (!result) {
				return result;
			}
		}
		nox_common_list_remove_425920((uint32_t**)v2);
		free(v2);
	}
	return result;
}

//----- (004D12A0) --------------------------------------------------------
int sub_4D12A0(int a1) {
	int* v1; // eax

	v1 = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1548492));
	if (!v1) {
		return 0;
	}
	while (*(unsigned char*)(v1[3] + 2064) != a1) {
		v1 = nox_common_list_getNextSafe_4258A0(v1);
		if (!v1) {
			return 0;
		}
	}
	return 1;
}

void nox_xxx_mapSwitchLevel_4D12E0_tileFree() {
	for (int j = 0; j < ptr_5D4594_2650668_cap * 44; j += 44) {
		for (int k = 0; k < ptr_5D4594_2650668_cap; k++) {
			*(uint8_t*)((uint32_t)(ptr_5D4594_2650668[k]) + j) = 0;
			*(uint32_t*)((uint32_t)(ptr_5D4594_2650668[k]) + j + 4) = 255;
			*(uint32_t*)((uint32_t)(ptr_5D4594_2650668[k]) + j + 24) = 255;
			nox_xxx_tileFreeTile_422200((uint32_t)(ptr_5D4594_2650668[k]) + j + 4);
			nox_xxx_tileFreeTile_422200((uint32_t)(ptr_5D4594_2650668[k]) + j + 24);
		}
	}
}

//----- (004D15C0) --------------------------------------------------------
void sub_4D15C0() { *getMemU32Ptr(0x5D4594, 1548508) = 0; }

//----- (004D1600) --------------------------------------------------------
int nox_xxx_scavengerTreasureMax_4D1600() { return *getMemU32Ptr(0x5D4594, 1548528); }

//----- (004D1610) --------------------------------------------------------
void sub_4D1610() { *getMemU32Ptr(0x5D4594, 1548528) = 0; }

//----- (004D23C0) --------------------------------------------------------
void sub_51A100();
int nox_xxx_servResetPlayers_4D23C0() {
	char* i; // esi
	int v2;  // [esp-Ch] [ebp-14h]

	for (i = nox_common_playerInfoGetFirst_416EA0(); i; i = nox_common_playerInfoGetNext_416EE0((int)i)) {
		if (*((uint32_t*)i + 514)) {
			dword_5d4594_2649712 &= ~(1 << i[2064]);
			v2 = *((uint32_t*)i + 514);
			i[3676] = 2;
			nox_xxx_playerMakeDefItems_4EF7D0(v2, 1, 0);
			*((uint32_t*)i + 535) = 0;
			*((uint32_t*)i + 534) = 0;
		}
	}
	sub_51A100();
	nox_common_gameFlags_unset_40A540(0x20000);
	nox_xxx_netGameSettings_4DEF00();
	nox_server_gameUnsetMapLoad_40A690();
	return 1;
}

//----- (004D3050) --------------------------------------------------------
char* nox_xxx_netReportAllLatency_4D3050() {
	char* result; // eax
	bool v1;      // zf
	int i;        // esi
	char v3[5];   // [esp+0h] [ebp-8h]

	v3[0] = -41;
	if (!dword_5d4594_1548700 || (result = nox_common_playerInfoGetNext_416EE0(*(int*)&dword_5d4594_1548700),
								  (dword_5d4594_1548700 = result) == 0)) {
		result = nox_common_playerInfoGetFirst_416EA0();
		dword_5d4594_1548700 = result;
	}
	if (result) {
		for (int k = 0; result[2064] != 31 && k < 32; k++) {
			v1 = sub_554240((unsigned char)result[2064]) == 0;
			result = *(char**)&dword_5d4594_1548700;
			if (!v1) {
				break;
			}
			result = nox_common_playerInfoGetNext_416EE0(*(int*)&dword_5d4594_1548700);
			dword_5d4594_1548700 = result;
			if (!result) {
				result = nox_common_playerInfoGetFirst_416EA0();
				dword_5d4594_1548700 = result;
			}
		}
		if (result) {
			*(uint16_t*)&v3[1] = *((uint16_t*)result + 1030);
			*(uint16_t*)&v3[3] = sub_554240((unsigned char)result[2064]);
			result = nox_common_playerInfoGetFirst_416EA0();
			for (i = (int)result; result; i = (int)result) {
				nox_netlist_addToMsgListCli_40EBC0(*(unsigned char*)(i + 2064), 1, v3, 5);
				result = nox_common_playerInfoGetNext_416EE0(i);
			}
		}
	}
	return result;
}

//----- (004D39F0) --------------------------------------------------------
int sub_4D39F0(const char* a3) {
	unsigned int v1;    // ecx
	char v2;            // dl
	unsigned char* v3;  // edi
	const char* v4;     // esi
	int v5;             // edx
	int v6;             // eax
	unsigned char* v7;  // edi
	unsigned int v8;    // ecx
	unsigned char* v9;  // edi
	const char* v10;    // esi
	unsigned char* v11; // edi
	int v12;            // ecx
	int v13;            // edx
	int v14;            // eax
	char* v15;          // edi
	unsigned char v16;  // cl
	int result;         // eax
	char v18[2048];     // [esp+10h] [ebp-800h]

	*getMemU64Ptr(0x5D4594, 1549772) = nox_platform_get_ticks();
	memset(getMemAt(0x973F18, 35912), 0, 0x48u);
	*getMemU32Ptr(0x973F18, 35912) = 0;
	*getMemU32Ptr(0x973F18, 35916) = 0;
	dword_5d4594_3835348 = 0;
	dword_5d4594_3835356 = 255;
	dword_5d4594_3835352 = 0;
	dword_5d4594_3835360 = 0;
	dword_5d4594_3835364 = 1;
	dword_5d4594_3835368 = 1;
	dword_5d4594_3835372 = 1;
	*getMemU32Ptr(0x973F18, 35948) = 0;
	*getMemU32Ptr(0x973F18, 35952) = 0;
	*getMemU32Ptr(0x973F18, 35956) = 0;
	dword_5d4594_3835388 = 0;
	dword_5d4594_3835392 = 1;
	dword_5d4594_3835396 = -1;
	*getMemU8Ptr(0x973F18, 35972) = 2;
	*getMemU32Ptr(0x973F18, 35976) = 0;
	*getMemU32Ptr(0x973F18, 35980) = 0;
	sub_51D0E0();
	if (a3) {
		v1 = strlen(a3) + 1;
		v2 = v1;
		v1 >>= 2;
		memcpy(getMemAt(0x973F18, 42152), a3, 4 * v1);
		v4 = &a3[4 * v1];
		v3 = getMemAt(0x973F18, 42152 + 4 * v1);
		LOBYTE(v1) = v2;
		v5 = *getMemU32Ptr(0x587000, 197560);
		memcpy(v3, v4, v1 & 3);
		strcpy((char*)getMemAt(0x973F18, 36008), a3);
		v6 = *getMemU32Ptr(0x587000, 197564);
		v7 = getMemAt(0x973F18, 36008 + strlen((const char*)getMemAt(0x973F18, 36008)));
		*(uint32_t*)v7 = *getMemU32Ptr(0x587000, 197556);
		*((uint32_t*)v7 + 1) = v5;
		*((uint32_t*)v7 + 2) = v6;
		v8 = strlen(a3) + 1;
		LOBYTE(v5) = v8;
		v8 >>= 2;
		memcpy(getMemAt(0x973F18, 38056), a3, 4 * v8);
		v10 = &a3[4 * v8];
		v9 = getMemAt(0x973F18, 38056 + 4 * v8);
		LOBYTE(v8) = v5;
		LOWORD(v5) = *getMemU16Ptr(0x587000, 197576);
		memcpy(v9, v10, v8 & 3);
		v11 = getMemAt(0x973F18, 38057 + strlen((const char*)getMemAt(0x973F18, 38056)));
		v12 = *getMemU32Ptr(0x587000, 197572);
		*(uint32_t*)--v11 = *getMemU32Ptr(0x587000, 197568);
		LOBYTE(v6) = getMemByte(0x587000, 197578);
		*((uint32_t*)v11 + 1) = v12;
		*((uint16_t*)v11 + 4) = v5;
		v11[10] = v6;
		nox_fs_remove((const char*)getMemAt(0x973F18, 36008));
		nox_fs_remove((const char*)getMemAt(0x973F18, 38056));
	} else {
		*getMemU8Ptr(0x973F18, 42152) = getMemByte(0x5D4594, 1549780);
		*getMemU8Ptr(0x973F18, 40104) = getMemByte(0x5D4594, 1549784);
		*getMemU8Ptr(0x973F18, 36008) = getMemByte(0x5D4594, 1549788);
		*getMemU8Ptr(0x973F18, 38056) = getMemByte(0x5D4594, 1549792);
	}
	nox_xxx_mapReset_5028E0();
	v13 = *getMemU32Ptr(0x587000, 197584);
	strcpy(v18, a3);
	v14 = *getMemU32Ptr(0x587000, 197588);
	v15 = &v18[strlen(v18)];
	*(uint32_t*)v15 = *getMemU32Ptr(0x587000, 197580);
	v16 = getMemByte(0x587000, 197592);
	*((uint32_t*)v15 + 1) = v13;
	*((uint32_t*)v15 + 2) = v14;
	v15[12] = v16;
	sub_502A50(v18);
	sub_502AB0(v18);
	result = sub_502B10();
	dword_5d4594_3835312 = 0;
	*getMemU32Ptr(0x973F18, 35880) = 0;
	*getMemU32Ptr(0x5D4594, 1599580) = 0;
	return result;
}
// 4D39F0: using guessed type char var_800[2048];

//----- (004D3C50) --------------------------------------------------------
void nox_xxx_tileInitdataClear_4D3C50(const void* a1) { memcpy(getMemAt(0x973F18, 35912), a1, 0x48u); }

//----- (004D3C70) --------------------------------------------------------
unsigned char* sub_4D3C70() { return getMemAt(0x973F18, 35912); }

//----- (004D3C80) --------------------------------------------------------
uint32_t* sub_4D3C80(uint32_t* a1) {
	uint32_t* result; // eax
	int v2;           // ebp
	int v3;           // ecx
	int v4;           // edi
	int v5;           // esi
	int v6;           // edx
	int v7;           // ebx
	int v8;           // ecx
	int v9;           // esi
	int v10;          // [esp+10h] [ebp-10h]
	int v11;          // [esp+14h] [ebp-Ch]
	int v12;          // [esp+1Ch] [ebp-4h]

	result = a1;
	v2 = a1[3];
	v10 = *a1;
	v3 = a1[1];
	v11 = a1[1];
	if (v2 < v3) {
		v3 = a1[3];
		v10 = a1[2];
		v11 = a1[3];
	}
	if (a1[5] < v3) {
		v10 = a1[4];
		v3 = a1[5];
		v11 = a1[5];
	}
	if (a1[7] < v3) {
		v10 = a1[6];
		v11 = a1[7];
	}
	v4 = a1[2];
	v12 = a1[3];
	if (*a1 < v4) {
		v4 = *a1;
		v12 = a1[1];
	}
	if (a1[4] < v4) {
		v4 = a1[4];
		v12 = a1[5];
	}
	v5 = a1[6];
	if (v5 < v4) {
		v4 = a1[6];
		v12 = a1[7];
	}
	v6 = a1[4];
	v7 = a1[5];
	if (*a1 > v6) {
		v6 = *a1;
		v7 = a1[1];
	}
	if (a1[2] > v6) {
		v7 = a1[3];
		v6 = a1[2];
	}
	if (v5 > v6) {
		v6 = a1[6];
		v7 = a1[7];
	}
	v8 = a1[7];
	v9 = a1[6];
	if (a1[1] > v8) {
		v9 = *a1;
		v8 = a1[1];
	}
	if (v2 > v8) {
		v9 = a1[2];
		v8 = a1[3];
	}
	if (a1[5] > v8) {
		v9 = a1[4];
		v8 = a1[5];
	}
	a1[6] = v9;
	*a1 = v10;
	a1[2] = v4;
	a1[7] = v8;
	a1[1] = v11;
	a1[4] = v6;
	a1[5] = v7;
	a1[3] = v12;
	return result;
}

//----- (004D3E30) --------------------------------------------------------
int sub_4D3E30(float2* a1, float2* a2) {
	int result; // eax

	if (!a1 || !a2) {
		return 0;
	}
	if (a1->field_0 <= 80.5) {
		a1->field_0 = 82.5;
	}
	if (a1->field_4 <= 80.5) {
		a1->field_4 = 81.5;
	}
	if (a1->field_0 >= 5853.5) {
		a1->field_0 = 5851.5;
	}
	if (a1->field_4 >= 5853.5) {
		a1->field_4 = 5852.5;
	}
	result = 1;
	a2->field_0 = (a1->field_0 - 1.0 - a1->field_4) * 0.70710677;
	a2->field_4 = (a1->field_4 + a1->field_0 - 5912.0) * 0.70710677;
	return result;
}

//----- (004D6000) --------------------------------------------------------
int sub_4D6000(nox_object_t* a1p) {
	int a1 = a1p;
	int result; // eax
	int v2;     // esi

	result = 0;
	if (a1) {
		v2 = *(uint32_t*)(a1 + 748);
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4652) = 0;
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4656) = 0;
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4660) = 0;
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4664) = 0;
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4668) = 0;
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4672) = 0;
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4676) = 0;
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4680) = 0;
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4684) = 0;
		*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4688) = nox_game_getQuestStage_4E3CC0();
		result = *(uint32_t*)(v2 + 276);
		*(uint32_t*)(result + 4692) = 63;
	}
	return result;
}

//----- (004D60B0) --------------------------------------------------------
int sub_4D60B0() {
	int result; // eax
	int i;      // esi

	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		sub_4D6000(i);
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (004D60E0) --------------------------------------------------------
uint32_t* sub_4D60E0(int a1) {
	uint32_t* result; // eax
	int v2;           // ecx

	result = (uint32_t*)a1;
	if ((*(uint8_t*)(a1 + 16) & 0x20) != 32) {
		v2 = *(uint32_t*)(a1 + 748);
		result = *(uint32_t**)(v2 + 276);
		if (result[1198] == 1) {
			++result[1163];
			result = *(uint32_t**)(v2 + 276);
			result[1173] |= 1u;
		}
	}
	return result;
}

//----- (004D6130) --------------------------------------------------------
int sub_4D6130(int a1) {
	int result; // eax
	int v2;     // eax

	result = a1;
	if (a1) {
		if ((*(uint8_t*)(a1 + 16) & 0x20) != 32) {
			v2 = *(uint32_t*)(a1 + 748);
			++*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4660);
			result = *(uint32_t*)(v2 + 276);
			*(uint32_t*)(result + 4692) |= 2u;
		}
	}
	return result;
}

//----- (004D6170) --------------------------------------------------------
int sub_4D6170(int a1) {
	int result; // eax
	int v2;     // eax

	result = a1;
	if (a1) {
		if ((*(uint8_t*)(a1 + 16) & 0x20) != 32) {
			v2 = *(uint32_t*)(a1 + 748);
			++*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4664);
			result = *(uint32_t*)(v2 + 276);
			*(uint32_t*)(result + 4692) |= 4u;
		}
	}
	return result;
}

//----- (004D61B0) --------------------------------------------------------
void sub_4D61B0(int a1) {
	int v2; // eax
	if (a1) {
		if ((*(uint8_t*)(a1 + 16) & 0x20) != 32) {
			v2 = *(uint32_t*)(a1 + 748);
			++*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4668);
			int result = *(uint32_t*)(v2 + 276);
			*(uint32_t*)(result + 4692) |= 8u;
		}
	}
}

//----- (004D61F0) --------------------------------------------------------
int sub_4D61F0(int a1) {
	int result; // eax
	int v2;     // eax

	result = a1;
	if (a1) {
		if ((*(uint8_t*)(a1 + 16) & 0x20) != 32) {
			v2 = *(uint32_t*)(a1 + 748);
			++*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4672);
			++*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4676);
			result = *(uint32_t*)(v2 + 276);
			*(uint32_t*)(result + 4692) |= 0x10u;
		}
	}
	return result;
}

//----- (004D6540) --------------------------------------------------------
unsigned int sub_4D6540(int a1) {
	int v1;              // edi
	unsigned int* v2;    // eax
	unsigned int* v3;    // ebx
	unsigned int result; // eax
	unsigned int v5;     // ebp
	int i;               // edi
	unsigned int* v7;    // esi
	unsigned int v8;     // eax
	int v9;              // eax
	float v10;           // [esp+0h] [ebp-28h]
	unsigned int v11;    // [esp+14h] [ebp-14h]
	unsigned int v12;    // [esp+18h] [ebp-10h]
	float v13;           // [esp+20h] [ebp-8h]
	unsigned int v14;    // [esp+2Ch] [ebp+4h]

	v1 = nox_xxx_player_4E3CE0();
	v2 = (unsigned int*)nox_common_playerInfoFromNum_417090(a1);
	v3 = v2;
	if (!v2 || !v2[1198]) {
		return 0;
	}
	if (v1 == 1) {
		result = sub_4D66E0(v2[1167], v2[1168], v2[1166], v2[1172]);
	} else {
		v5 = 0;
		v11 = 0;
		v14 = 0;
		v12 = 1;
		v13 = 1.0;
		for (i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
			v7 = *(unsigned int**)(*(uint32_t*)(i + 748) + 276);
			if (v3[1198]) {
				v8 = sub_4D66E0(v7[1167], v7[1168], 0, v7[1172]);
				if (v8 > v5) {
					v5 = v8;
				}
				v11 += v7[1167];
				v12 = v7[1172];
				v14 += v7[1168];
			}
		}
		v9 = sub_4D66E0(v11, v14, 0, v12);
		if (v5 > 0) {
			v13 = (double)(unsigned int)v9 / (double)v5;
		}
		v10 = (double)(unsigned int)sub_4D66E0(v3[1167], v3[1168], v3[1166], v3[1172]) * v13;
		result = nox_float2int(v10);
	}
	if (result > 0x3B9AC9FF) {
		result = 999999999;
	}
	return result;
}

//----- (004D66E0) --------------------------------------------------------
int sub_4D66E0(unsigned int a1, unsigned int a2, unsigned int a3, unsigned int a4) {
	float v5; // [esp+4h] [ebp-10h]

	v5 = nox_double2float(pow((double)a4, *(long double*)getMemAt(0x581450, 10088))) *
		 ((double)a1 * 10.0 + (double)a2 * 35.0 + (double)a3 * 0.1);
	return nox_float2int(v5);
}

//----- (004D6770) --------------------------------------------------------
int sub_4D6770(int a1) {
	char* v1;    // ebp
	int v2;      // ebx
	int v3;      // edi
	char* v4;    // esi
	int v5;      // eax
	char v7[90]; // [esp+Ch] [ebp-5Ch]

	v1 = nox_common_playerInfoFromNum_417090(a1);
	memset(v7, 0, 0x58u);
	*(uint16_t*)&v7[88] = 0;
	v7[0] = -16;
	v7[1] = 12;
	*(uint16_t*)&v7[2] = sub_4D7300();
	v2 = 0;
	v3 = nox_xxx_getFirstPlayerUnit_4DA7C0();
	if (v3) {
		v4 = &v7[8];
		do {
			v5 = *(uint32_t*)(v3 + 748);
			if (*(uint32_t*)(*(uint32_t*)(v5 + 276) + 4792) == 1 && v2 < 6) {
				*((uint16_t*)v4 - 1) = *(uint16_t*)(v3 + 36);
				*(uint16_t*)v4 = *(uint16_t*)(*(uint32_t*)(v5 + 276) + 4668);
				*((uint16_t*)v4 + 3) = *(uint16_t*)(*(uint32_t*)(v5 + 276) + 4664);
				*((uint16_t*)v4 + 1) = *(uint16_t*)(*(uint32_t*)(v5 + 276) + 4672);
				*((uint16_t*)v4 + 2) = *(uint16_t*)(*(uint32_t*)(v5 + 276) + 4680);
				*((uint32_t*)v4 + 2) = sub_4D6540(*(unsigned char*)(*(uint32_t*)(v5 + 276) + 2064));
				++v2;
				v4 += 14;
				*(uint16_t*)&v7[4] = *((uint16_t*)v1 + 2344);
			}
			v3 = nox_xxx_getNextPlayerUnit_4DA7F0(v3);
		} while (v3);
	}
	return nox_xxx_netSendPacket1_4E5390(a1, (int)v7, 90, 0, 1);
}

//----- (004D6880) --------------------------------------------------------
int sub_4D6880(int a1, int a2) {
	char v3[69]; // [esp+8h] [ebp-48h]

	memset(v3, 0, 0x44u);
	v3[68] = 0;
	v3[0] = -16;
	v3[1] = 13;
	if (a2) {
		v3[4] |= 1u;
	}
	if (sub_51A950()) {
		v3[4] |= 2u;
	}
	*(uint16_t*)&v3[2] = nox_game_getQuestStage_4E3CC0();
	nox_server_currentMapGetFilename_409B30();
	strcpy(&v3[5], sub_4D6940());
	nox_server_currentMapGetFilename_409B30();
	strcpy(&v3[37], sub_4D6950());
	return nox_xxx_netSendPacket1_4E5390(a1, (int)v3, 69, 0, 1);
}

//----- (004D6940) --------------------------------------------------------
char* sub_4D6940() { return (char*)getMemAt(0x973F18, 3838); }

//----- (004D6950) --------------------------------------------------------
char* sub_4D6950() { return (char*)getMemAt(0x973F18, 3806); }

//----- (004D6960) --------------------------------------------------------
int nox_game_sendQuestStage_4D6960(int a1) {
	char v2[69] = {0};
	v2[0] = -16;
	v2[1] = 14;
	v2[4] = 0;
	if (sub_51A950()) {
		v2[4] |= 2u;
	}
	*(uint16_t*)&v2[2] = nox_game_getQuestStage_4E3CC0();
	nox_server_currentMapGetFilename_409B30();
	strcpy(&v2[5], sub_4D6940());
	nox_server_currentMapGetFilename_409B30();
	strcpy(&v2[37], sub_4D6950());
	return nox_xxx_netSendPacket1_4E5390(a1, (int)v2, 69, 0, 1);
}

//----- (004D6A20) --------------------------------------------------------
int sub_4D6A20(int a1, int a2) {
	short v3; // cx
	int v5;   // [esp+0h] [ebp-4h]

	v3 = *(uint16_t*)(a2 + 40);
	LOWORD(v5) = 4080;
	HIWORD(v5) = v3;
	return nox_xxx_netSendPacket0_4E5420(a1, &v5, 4, 0, 1);
}

//----- (004D6F50) --------------------------------------------------------
int nox_xxx_isQuest_4D6F50() { return *getMemU32Ptr(0x5D4594, 1556160); }

//----- (004D6F60) --------------------------------------------------------
int nox_xxx_setQuest_4D6F60(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1556160) = a1;
	return result;
}

//----- (004D6F70) --------------------------------------------------------
int sub_4D6F70() { return *getMemU32Ptr(0x5D4594, 1556164); }

//----- (004D6F80) --------------------------------------------------------
int sub_4D6F80(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1556164) = a1;
	return result;
}

//----- (004D6FA0) --------------------------------------------------------
int sub_4D6FA0() { return *getMemU32Ptr(0x5D4594, 1556104); }

//----- (004D70B0) --------------------------------------------------------
char* sub_4D70B0() {
	char* result; // eax

	result = *(char**)getMemAt(0x5D4594, 1556152);
	if (!*getMemU32Ptr(0x5D4594, 1556152)) {
		result = sub_4169F0();
	}
	return result;
}

//----- (004D70C0) --------------------------------------------------------
int nox_xxx_bookCreatureTest_4D70C0(int a1) {
	return nox_common_gameFlags_check_40A5C0(4096) || a1 != 37 && a1 != 38 && a1 != 40 && a1 != 39;
}

//----- (004D7100) --------------------------------------------------------
int sub_4D7100(int a1) {
	return nox_common_gameFlags_check_40A5C0(4096) || a1 != 111 && a1 != 112 && a1 != 114 && a1 != 113;
}

//----- (004D7150) --------------------------------------------------------
int sub_4D7150() {
	int result; // eax
	int i;      // edi
	int* v2;    // esi
	int v3;     // eax

	result = dword_5d4594_1556144;
	if (dword_5d4594_1556144) {
		if (gameFrame() > *(int*)&dword_5d4594_1556144) {
			result = nox_xxx_getFirstPlayerUnit_4DA7C0();
			for (i = result; result; i = result) {
				v2 = *(int**)(i + 748);
				v3 = v2[69];
				if ((*(uint8_t*)(v3 + 3680) & 1) == 1 && *(uint32_t*)(v3 + 4792) == 1 && !v2[78] && !v2[79] &&
					(*(uint32_t*)(v3 + 3680) & 0x10) == 16) {
					sub_4DF3C0(v2[69]);
					nox_xxx_playerLeaveObserver_0_4E6AA0(v2[69]);
					nox_xxx_playerCameraUnlock_4E6040(i);
				}
				result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
			}
		}
	}
	return result;
}

//----- (004D71E0) --------------------------------------------------------
int sub_4D71E0(int a1) {
	int result; // eax

	result = a1;
	dword_5d4594_1556136 = a1;
	return result;
}

//----- (004D71F0) --------------------------------------------------------
void sub_4D72B0(int a1);
unsigned int sub_4D71F0() {
	unsigned int result; // eax
	int v1;              // esi
	char v2;             // al

	result = dword_5d4594_1556136;
	if (dword_5d4594_1556136) {
		if ((unsigned int)(gameFrame() - dword_5d4594_1556136) >= 0x2328) {
			v1 = 0;
			result = nox_xxx_getFirstPlayerUnit_4DA7C0();
			if (result) {
				do {
					if (*(uint32_t*)(*(uint32_t*)(result + 748) + 308)) {
						v1 = 1;
					}
					result = nox_xxx_getNextPlayerUnit_4DA7F0(result);
				} while (result);
				if (v1) {
					result = nox_xxx_player_4E3CE0();
					if (result > 1) {
						sub_4D71E0(0);
						result = sub_4D72C0();
						if (!result) {
							sub_4D72B0(1);
							v2 = sub_4D72C0();
							result = sub_4D7280(255, v2);
						}
					}
				}
			}
		}
	}
	return result;
}

//----- (004D7280) --------------------------------------------------------
int sub_4D7280(int a1, char a2) {
	char v4[3]; // [esp+0h] [ebp-4h]
	v4[0] = -16;
	v4[1] = 24;
	v4[2] = a2;
	return nox_xxx_netSendPacket1_4E5390(a1, (int)v4, 3, 0, 1);
}

//----- (004D72D0) --------------------------------------------------------
int sub_4D72D0(int a1) {
	int result; // eax

	result = dword_5d4594_1556128;
	*getMemU32Ptr(0x5D4594, 1556132) = dword_5d4594_1556128;
	dword_5d4594_1556128 = a1;
	return result;
}

//----- (004D7300) --------------------------------------------------------
int sub_4D7300() { return *getMemU32Ptr(0x5D4594, 1556132); }

//----- (004D7430) --------------------------------------------------------
int sub_4D7430() { return *getMemU32Ptr(0x5D4594, 1556116); }

//----- (004D7440) --------------------------------------------------------
int sub_4D7440(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1556116) = a1;
	return result;
}

//----- (004D7450) --------------------------------------------------------
int sub_4D7450(int a1, short a2) {
	char v3[4]; // [esp+0h] [ebp-4h]

	v3[0] = -16;
	v3[1] = 29;
	*(uint16_t*)&v3[2] = a2;
	return nox_xxx_netSendPacket0_4E5420(a1, v3, 4, 0, 1);
}

//----- (004D7480) --------------------------------------------------------
void sub_4D7480(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;     // edi
	int v2;     // eax
	float2* v3; // ebx

	if (a1 && *(uint8_t*)(a1 + 8) & 4) {
		v1 = *(uint32_t*)(a1 + 748);
		v2 = *(uint32_t*)(v1 + 316);
		if (v2) {
			v3 = *(float2**)(v2 + 700);
			nox_xxx_playerLeaveObserver_0_4E6AA0(*(uint32_t*)(v1 + 276));
			nox_xxx_playerCameraUnlock_4E6040(a1);
			*(uint32_t*)(v1 + 316) = 0;
			nox_xxx_unitMove_4E7010(a1, v3 + 10);
			nox_xxx_aud_501960(312, a1, 2, *(uint32_t*)(a1 + 36));
			nox_xxx_netSendPointFx_522FF0(129, v3 + 10);
		}
	}
}

//----- (004D7520) --------------------------------------------------------
char sub_4D7520(int a1) {
	int v1; // eax
	int i;  // esi
	int v3; // eax
	int v4; // esi
	int v5; // edi

	LOBYTE(v1) = getMemByte(0x5D4594, 1556120);
	if (*getMemU32Ptr(0x5D4594, 1556120) != 1) {
		*getMemU32Ptr(0x5D4594, 1556120) = a1;
		return v1;
	}
	if (a1) {
		*getMemU32Ptr(0x5D4594, 1556120) = a1;
		return v1;
	}
	for (i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
		v3 = *(uint32_t*)(i + 748);
		if (*(uint32_t*)(*(uint32_t*)(v3 + 276) + 4792) && *(uint32_t*)(v3 + 316)) {
			sub_4D7480(i);
		}
	}
	v1 = nox_server_getFirstObject_4DA790();
	v4 = v1;
	if (!v1) {
		*getMemU32Ptr(0x5D4594, 1556120) = a1;
		return v1;
	}
	do {
		v5 = nox_server_getNextObject_4DA7A0(v4);
		LOBYTE(v1) = *(uint8_t*)(v4 + 8);
		if (v1 & 0x20 && *(uint8_t*)(v4 + 12) & 2) {
			LOBYTE(v1) = nox_xxx_objectSetOff_4E7600(v4);
		}
		v4 = v5;
	} while (v5);
	*getMemU32Ptr(0x5D4594, 1556120) = 0;
	return v1;
}

//----- (004D75E0) --------------------------------------------------------
int sub_4D75E0() { return *getMemU32Ptr(0x5D4594, 1556120); }

//----- (004D75F0) --------------------------------------------------------
int sub_4D75F0(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1556108) = a1;
	return result;
}

//----- (004D7600) --------------------------------------------------------
void nox_server_checkWarpGate_4D7600() {
	int exp = nox_xxx_player_4E3CE0();
	if (!exp) {
		return;
	}
	if ((unsigned int)(gameFrame() - *getMemU32Ptr(0x5D4594, 1556108)) < 30) {
		return;
	}
	int inGate = 0;
	for (void* unit = nox_xxx_getFirstPlayerUnit_4DA7C0(); unit; unit = nox_xxx_getNextPlayerUnit_4DA7F0(unit)) {
		int v3 = *(uint32_t*)((int)unit + 748);
		if (*(uint32_t*)(*(uint32_t*)(v3 + 276) + 4792) && *(uint32_t*)(v3 + 316)) {
			++inGate;
		}
	}
	if (exp != inGate) {
		// not all players are in the gate
		return;
	}
	if (!nox_server_questMaybeWarp_4E8F60()) {
		// warp failed
		for (void* unit = nox_xxx_getFirstPlayerUnit_4DA7C0(); unit; unit = nox_xxx_getNextPlayerUnit_4DA7F0(unit)) {
			int v5 = *(uint32_t*)((int)unit + 748);
			if (*(uint32_t*)(*(uint32_t*)(v5 + 276) + 4792) && *(uint32_t*)(v5 + 316)) {
				sub_4D7480(unit);
				if (exp <= 1) {
					nox_xxx_netPriMsgToPlayer_4DA2C0(unit, "Gauntlet.c:WarpRestrictedSolo", 0);
				} else {
					nox_xxx_netPriMsgToPlayer_4DA2C0(unit, "Gauntlet.c:WarpRestrictedMulti", 0);
				}
			}
		}
	}
}

//----- (004D76E0) --------------------------------------------------------
int sub_4D76E0(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1556124) = a1;
	return result;
}

//----- (004D76F0) --------------------------------------------------------
int sub_4D76F0() { return *getMemU32Ptr(0x5D4594, 1556124); }

//----- (004D79A0) --------------------------------------------------------
int sub_4D79A0(char a1) {
	int result; // eax

	result = ~(1 << a1);
	*getMemU32Ptr(0x5D4594, 1556300) &= result;
	return result;
}

//----- (004D79C0) --------------------------------------------------------
int sub_4D79C0(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;     // esi
	int result; // eax
	int v3;     // ecx

	v1 = *(uint32_t*)(a1 + 748);
	sub_4D9D20(255, a1);
	sub_4D6000(a1);
	for (result = nox_xxx_getFirstPlayerUnit_4DA7C0(); result; result = nox_xxx_getNextPlayerUnit_4DA7F0(result)) {
		v3 = *(uint32_t*)(result + 748);
		*(uint8_t*)(*(unsigned char*)(*(uint32_t*)(v1 + 276) + 2064) + v3 + 452) = 0;
		*(uint32_t*)(v3 + 4 * *(unsigned char*)(*(uint32_t*)(v1 + 276) + 2064) + 324) = 0;
		*(uint8_t*)(*(unsigned char*)(*(uint32_t*)(v1 + 276) + 2064) + v3 + 484) = 0;
		*(uint8_t*)(*(unsigned char*)(*(uint32_t*)(v1 + 276) + 2064) + v3 + 516) = 0;
	}
	return result;
}

//----- (004D7A60) --------------------------------------------------------
int sub_4D7A60(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1556172 + 4 * a1) = gameFrame();
	return result;
}

//----- (004D7A80) --------------------------------------------------------
int sub_4D7A80() {
	unsigned char* v0; // ebp
	int v1;            // edi
	int v2;            // esi
	char* v3;          // eax
	int v4;            // eax
	int v5;            // edi
	int v6;            // ecx
	int result;        // eax

	v0 = getMemAt(0x5D4594, 1556172);
	v1 = 324 - (uint32_t)getMemAt(0x5D4594, 1556172);
	v2 = 484;
	do {
		v3 = nox_common_playerInfoFromNum_417090(v2 - 484);
		if (v3 && *((uint32_t*)v3 + 523) && *((uint32_t*)v3 + 514) && *((uint32_t*)v3 + 1198) == 1) {
			*(uint32_t*)v0 = 0;
		} else if (*(uint32_t*)v0 && gameFrame() - *(uint32_t*)v0 > (unsigned int)(30 * gameFPS())) {
			v4 = nox_xxx_getFirstPlayerUnit_4DA7C0();
			if (v4) {
				v5 = (int)&v0[v1];
				do {
					v6 = *(uint32_t*)(v4 + 748);
					*(uint8_t*)(v2 + v6 - 32) = 0;
					*(uint32_t*)(v5 + v6) = 0;
					*(uint8_t*)(v2 + v6) = 0;
					*(uint8_t*)(v2 + v6 + 32) = 0;
					v4 = nox_xxx_getNextPlayerUnit_4DA7F0(v4);
				} while (v4);
				v1 = 324 - (uint32_t)getMemAt(0x5D4594, 1556172);
			}
			*(uint32_t*)v0 = 0;
		}
		++v2;
		v0 += 4;
		result = v2 - 484;
	} while (v2 - 484 < 32);
	return result;
}

//----- (004D7B40) --------------------------------------------------------
int sub_4D7B40() {
	int result; // eax

	result = 0;
	memset(getMemAt(0x5D4594, 1556172), 0, 0x80u);
	return result;
}

// C varargs are formatted here; Go owns message serialization and fanout.
extern int nox_gameplayTextLine(int unit, wchar2_t* text);
extern int nox_gameplayTextAll(char flags, wchar2_t* text);
int nox_xxx_netSendLineMessage_4D9EB0(int unit, wchar2_t* format, ...) {
 if (!unit || !(*(uint8_t*)(unit + 8) & 4)) { return unit; }
 wchar2_t text[256];
 va_list args;
 va_start(args, format);
 nox_vswprintf(text, format, args);
 va_end(args);
 return nox_gameplayTextLine(unit, text);
}
int nox_xxx_printToAll_4D9FD0(char flags, wchar2_t* format, ...) {
 wchar2_t text[256];
 va_list args;
 va_start(args, format);
 nox_vswprintf(text, format, args);
 va_end(args);
 return nox_gameplayTextAll(flags, text);
}

//----- (004DA9A0) --------------------------------------------------------
uint32_t* nox_xxx_unitNewAddShadow_4DA9A0(nox_object_t* a1p) {
	uint32_t* a1 = a1p;
	uint32_t* result; // eax
	int v2;           // ecx

	result = a1;
	v2 = a1[4];
	if (!(v2 & 0x410000)) {
		a1[118] = 0;
		a1[4] = v2 | 0x10000;
		a1[117] = dword_5d4594_1556856;
		if (dword_5d4594_1556856) {
			*(uint32_t*)(dword_5d4594_1556856 + 472) = a1;
		}
		dword_5d4594_1556856 = a1;
	}
	return result;
}

//----- (004DA9F0) --------------------------------------------------------
uint32_t* nox_xxx_action_4DA9F0(nox_object_t* a1p) {
	uint32_t* a1 = a1p;
	uint32_t* result; // eax
	int v2;           // ecx
	int v3;           // ecx
	int v4;           // ecx

	result = a1;
	v2 = a1[4];
	if (v2 & 0x10000) {
		a1[4] = v2 & 0xFFFEFFFF;
		v3 = a1[118];
		if (v3) {
			*(uint32_t*)(v3 + 468) = a1[117];
		} else {
			dword_5d4594_1556856 = a1[117];
		}
		v4 = a1[117];
		if (v4) {
			*(uint32_t*)(v4 + 472) = a1[118];
		}
	}
	return result;
}

//----- (004DC550) --------------------------------------------------------
int nox_client_countSaveFiles_4DC550() {
	int v0;              // ebx
	char* v1;            // edi
	char* v5;            // eax
	int v6;              // esi
	char* v7;            // eax
	char PathName[1024]; // [esp+Ch] [ebp-400h]

	v0 = 0;
	v1 = nox_fs_root();
	strcpy(PathName, v1);
	strcat(PathName, "\\Save\\");
	nox_fs_mkdir(PathName);
	v5 = nox_fs_root();
	nox_sprintf(PathName, "%s\\Save\\AUTOSAVE\\Player.plr", v5);
	if (nox_fs_access(PathName, 0) != -1) {
		v0 = 1;
	}
	v6 = 13;
	do {
		v7 = nox_fs_root();
		nox_sprintf(PathName, "%s\\Save\\SAVE%04d\\Player.plr", v7, v0);
		if (nox_fs_access(PathName, 0) != -1) {
			++v0;
		}
		--v6;
	} while (v6);
	return v0;
}
// 4DC550: using guessed type char PathName[1024];

//----- (004DC630) --------------------------------------------------------
int nox_client_countPlayerFiles02_4DC630() {
	int v0;                                // ebp
	char* v1;                              // edi
	HANDLE v5;                             // esi
	char* v6;                              // eax
	struct _WIN32_FIND_DATAA FindFileData; // [esp+10h] [ebp-E40h]
	char PathName[1024];                   // [esp+150h] [ebp-D00h]
	char v10[1280];                        // [esp+550h] [ebp-900h]
	char v11[1024];                        // [esp+A50h] [ebp-400h]

	v0 = 0;
	v1 = nox_fs_root();
	strcpy(PathName, v1);
	strcat(PathName, "\\Save\\");
	strcpy(v11, PathName);
	nox_fs_mkdir(PathName);
	nox_fs_set_workdir(PathName);
	v5 = FindFirstFileA("*.plr", &FindFileData);
	if (v5 != (HANDLE)-1) {
		if (!(FindFileData.dwFileAttributes & 0x10)) {
			nox_sprintf(PathName, "%s%s", v11, FindFileData.cFileName);
			sub_41A000(PathName, v10);
			if (v10[0] & 2) {
				v0 = 1;
			}
		}
		while (FindNextFileA(v5, &FindFileData)) {
			if (!(FindFileData.dwFileAttributes & 0x10)) {
				nox_sprintf(PathName, "%s%s", v11, FindFileData.cFileName);
				sub_41A000(PathName, v10);
				if (v10[0] & 2) {
					++v0;
				}
			}
		}
		FindClose(v5);
	}
	v6 = nox_fs_root();
	nox_fs_set_workdir(v6);
	return v0;
}
// 4DC630: using guessed type char PathName[1024];

//----- (004DCC70) --------------------------------------------------------
int nox_xxx_mapLoadOrSaveMB_4DCC70(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1563072) = a1;
	return result;
}

//----- (004DCCB0) --------------------------------------------------------
int nox_xxx_game_4DCCB0() {
	int result; // eax
	char* v1;   // eax
	int v2;     // esi
	int v3;     // eax

	if (!nox_common_gameFlags_check_40A5C0(2048)) {
		return 1;
	}
	v1 = nox_common_playerInfoFromNum_417090(31);
	if (!v1 || (v2 = *((uint32_t*)v1 + 514)) == 0 || sub_4DCC90() || sub_4139B0() ||
		(unsigned int)(gameFrame() - *getMemU32Ptr(0x5D4594, 1563068)) < 0x1E || nox_xxx_guiCursor_477600() ||
		(v3 = *(uint32_t*)(v2 + 16), BYTE1(v3) & 0x40)) {
		result = 0;
	} else {
		result = sub_4DCC10(v2) != 0;
	}
	return result;
}

//----- (004DD180) --------------------------------------------------------
int nox_xxx_wall_4DF1E0(int a1);
char* nox_xxx_gameServerReadyMB_4DD180(int a1) {
	char* result; // eax
	char* v2;     // edi
	int v3;       // eax
	int i;        // esi
	int v5;       // eax

	result = nox_common_playerInfoFromNum_417090(a1);
	v2 = result;
	if (result) {
		nox_xxx_netNeedTimestampStatus_4174F0((int)result, 16);
		if (nox_common_gameFlags_check_40A5C0(0x2000) && !nox_common_gameFlags_check_40A5C0(128)) {
			v3 = *((uint32_t*)v2 + 514);
			if (v3) {
				nox_xxx_spellBuffOff_4FF5B0(v3, 23);
				nox_xxx_buffApplyTo_4FF380(*((uint32_t*)v2 + 514), 23, 5 * (uint16_t)gameFPS(), 5);
			}
			for (i = nox_server_getFirstObject_4DA790(); i; i = nox_server_getNextObject_4DA7A0(i)) {
				if (*(uint32_t*)(i + 8) & 0x10000000) {
					nox_xxx_netMarkMinimapObject_417190(a1, i, 1);
				}
			}
		}
		if ((v2[3680] & 1) == 1) {
			v5 = *((uint32_t*)v2 + 514);
			*((uint32_t*)v2 + 908) = *(uint32_t*)(v5 + 56);
			*((uint32_t*)v2 + 909) = *(uint32_t*)(v5 + 60);
			if (nox_common_gameFlags_check_40A5C0(512)) {
				nox_xxx_playerLeaveObserver_0_4E6AA0((int)v2);
			}
		}
		nox_xxx_wall_4DF1E0(a1);
		if (nox_common_gameFlags_check_40A5C0(4096) && sub_40A300() == 1) {
			nox_xxx_netGauntlet_4D9E70(a1);
		}
		if (a1 == 31 && nox_common_gameFlags_check_40A5C0(128)) {
			if (!nox_server_connectionType_3596 && 0) {
				sub_49C820();
				return (char*)nox_xxx_netStatsMultiplier_4D9C20(*((uint32_t*)v2 + 514));
			}
			if (nox_server_sanctuaryHelp_54276 == 1) {
				nox_xxx_cliShowHelpGui_49C560();
				return (char*)nox_xxx_netStatsMultiplier_4D9C20(*((uint32_t*)v2 + 514));
			}
			nox_xxx_guiServerOptsLoad_457500();
		}
		result = (char*)nox_xxx_netStatsMultiplier_4D9C20(*((uint32_t*)v2 + 514));
	}
	return result;
}

//----- (004DD9B0) --------------------------------------------------------
int nox_xxx_netGuiGameSettings_4DD9B0(char a1, const void* a2, int a3) {
	char v4[60]; // [esp+8h] [ebp-3Ch]

	v4[0] = -79;
	v4[1] = a1;
	memcpy(&v4[2], a2, 0x3Au);
	return nox_xxx_netSendPacket1_4E5390(a3, (int)v4, 60, 0, 0);
}

//----- (004DDA90) --------------------------------------------------------
void nox_xxx_netNewPlayerMakePacket_4DDA90(unsigned char* buf, nox_playerInfo* pl) {
	buf[0] = 45; // MSG_NEW_PLAYER
	*(uint16_t*)(&buf[1]) = pl->netCode;
	*(uint16_t*)(&buf[100]) = pl->lessons;
	*(uint16_t*)(&buf[102]) = pl->field_2140;
	*(uint32_t*)(&buf[104]) = pl->field_0;
	*(uint32_t*)(&buf[108]) = pl->field_4;
	buf[116] = pl->field_2152;
	buf[117] = pl->field_2156;
	buf[118] = pl->field_3676 == 3;
	*(uint32_t*)(&buf[112]) = pl->field_3680 & 0x423;
	memcpy(&buf[119], pl->field_2096, strlen(pl->field_2096) + 1);
	memcpy(&buf[3], &pl->info, 97);
}

//----- (004DDE10) --------------------------------------------------------
void sub_4DDE10(int a1, nox_playerInfo* a2p) {
	int a2 = a2p;
	uint32_t* result; // eax
	int i;            // esi
	int v4;           // eax

	result = *(uint32_t**)(a2 + 2056);
	if (result) {
		if (!dword_5d4594_1563276) {
			dword_5d4594_1563276 = nox_xxx_getNameId_4E3AA0("Flag");
		}
		result = *(uint32_t**)(a2 + 2056);
		for (i = result[126]; i; i = *(uint32_t*)(i + 496)) {
			v4 = *(uint32_t*)(i + 16);
			if (!(v4 & 0x100)) {
				result = *(uint32_t**)&dword_5d4594_1563276;
				if (*(unsigned short*)(i + 4) != dword_5d4594_1563276) {
					continue;
				}
			}
			sub_4D82F0(a1, (uint32_t*)i);
		}
	}
}

//----- (004DDF60) --------------------------------------------------------
void nox_xxx_playerSendMOTD_4DD140(int a1);
int nox_xxx_netPlayerIncomingServ_4DDF60(int a1) {
	int v1;             // ebx
	char* v2;           // esi
	int v3;             // eax
	int v4;             // edi
	int v5;             // ebp
	char* i;            // edi
	int v7;             // eax
	char* v8;           // eax
	char* j;            // edi
	unsigned char* v10; // eax
	int k;              // esi
	int v13;            // [esp+Ch] [ebp+4h]

	v1 = a1;
	v2 = nox_common_playerInfoFromNum_417090(a1);
	if (!v2) {
		abort();
	}
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		if (a1 != 31) {
			if (v2) {
				v3 = *((uint32_t*)v2 + 514);
				if (v3) {
					*(uint32_t*)(*(uint32_t*)(v3 + 748) + 552) = 1;
				}
			}
		}
		sub_4D9CF0(a1);
		if (v2 && *((uint32_t*)v2 + 514)) {
			sub_4D6000(*((uint32_t*)v2 + 514));
		}
	}
	sub_57B920(v2 + 16);
	v13 = *((uint32_t*)v2 + 514);
	dword_5d4594_2649712 |= 1 << v1;
	v4 = *(uint32_t*)(v13 + 56);
	v5 = *(uint32_t*)(v13 + 60);
	nox_xxx_newPlayerSendAllPlayers_4DE300(v1);
	*((uint32_t*)v2 + 1175) = 0;
	(*(void (**)(int, uint32_t))(v13 + 688))(v13, 0);
	v2[3676] = 3;
	if (!nox_common_gameFlags_check_40A5C0(512)) {
		*((uint32_t*)v2 + 908) = v4;
		*((uint32_t*)v2 + 909) = v5;
	}
	if (nox_server_sendMotd_108752 && nox_common_gameFlags_check_40A5C0(0x2000) &&
		!nox_common_gameFlags_check_40A5C0(4096)) {
		nox_xxx_playerSendMOTD_4DD140(v1);
	}
	for (i = nox_common_playerInfoGetFirst_416EA0(); i; i = nox_common_playerInfoGetNext_416EE0((int)i)) {
		v7 = *((uint32_t*)i + 514);
		if (v7) {
			if (i != v2) {
				nox_xxx_netMarkMinimapObject_417190(v1, v7, 1);
				nox_xxx_netMarkMinimapObject_417190((unsigned char)i[2064], *((uint32_t*)v2 + 514), 1);
				nox_xxx_netSendSimpleObject2_4DF360((unsigned char)i[2064], *((uint32_t*)v2 + 514));
				if (nox_common_gameFlags_check_40A5C0(4096)) {
					nox_xxx_netSendTeam_4D8670((unsigned char)i[2064], *((uint32_t**)v2 + 514));
					nox_xxx_netSendTeam_4D8670(v1, *((uint32_t**)i + 514));
				}
			}
		}
	}
	nox_xxx_servMinimapRevealFlag_4DE380(v1);
	sub_4DF2E0(v1);
	if (nox_common_gameFlags_check_40A5C0(1024) && !sub_40AA70((int)v2)) {
		nox_xxx_netNeedTimestampStatus_4174F0((int)v2, 256);
	}
	if (0) {
		sub_4161E0();
		sub_416690();
	}
	sub_4E8110(v1);
	if (nox_common_gameFlags_check_40A5C0(64)) {
		v8 = sub_4E8310();
		nox_xxx_netSendBallStatus_4D95F0(v1, *v8, *((uint16_t*)v8 + 1));
	} else if (nox_common_gameFlags_check_40A5C0(32)) {
		for (j = nox_server_teamFirst_418B10(); j; j = nox_server_teamNext_418B60((int)j)) {
			v10 = sub_4E8320(j[57]);
			nox_xxx_netSendFlagStatus_4D95A0(v1, *v10, v10[2], v10[1], *((uint16_t*)v10 + 2));
		}
	}
	nox_xxx_sendAllClientStatus_4175C0((int)v2);
	if (sub_409F40(0x2000)) {
		nox_xxx_sendAllPlayerIDs_4DE270((int)v2);
	}
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		for (k = nox_xxx_getFirstPlayerUnit_4DA7C0(); k; k = nox_xxx_getNextPlayerUnit_4DA7F0(k)) {
			if (*(uint32_t*)(*(uint32_t*)(*(uint32_t*)(k + 748) + 276) + 4792) == 1) {
				sub_4D9D20(v1, k);
			}
		}
	}
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		sub_4D7A60(v1);
	}
	return *(uint32_t*)(v13 + 36);
}

//----- (004DE270) --------------------------------------------------------
int nox_xxx_sendAllPlayerIDs_4DE270(int a1) {
	int result; // eax
	int i;      // esi
	int v3;     // [esp-1Ch] [ebp-28h]
	char v4[7]; // [esp+4h] [ebp-8h]

	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		if (*(uint32_t*)(i + 36) != *(uint32_t*)(a1 + 2060)) {
			if (*(uint32_t*)(*(uint32_t*)(i + 748) + 260)) {
				v4[0] = -46;
				*(uint16_t*)&v4[1] = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)i);
				v3 = *(unsigned char*)(a1 + 2064);
				*(uint16_t*)&v4[3] = *(uint16_t*)(i + 4);
				v4[5] = 1;
				v4[6] = 2;
				nox_xxx_netSendPacket0_4E5420(v3, v4, 7, 0, 1);
			}
		}
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (004DE300) --------------------------------------------------------
char* nox_xxx_newPlayerSendAllPlayers_4DE300(int a1) {
	char* result; // eax
	int i;        // esi
	char v3[132]; // [esp+4h] [ebp-84h]

	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		if (*(unsigned char*)(i + 2064) != a1 &&
			(*(uint8_t*)(i + 2064) != 31 || !nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING))) {
			nox_xxx_netNewPlayerMakePacket_4DDA90((int)v3, i);
			nox_xxx_netSendPacket1_4E5390(a1, (int)v3, 129, 0, 0);
			sub_4DDE10(a1, i);
		}
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
	return result;
}

//----- (004DE380) --------------------------------------------------------
int nox_xxx_servMinimapRevealFlag_4DE380(int a1) {
	int result; // eax
	int i;      // esi
	int v3;     // eax

	if (!*getMemU32Ptr(0x5D4594, 1563284)) {
		*getMemU32Ptr(0x5D4594, 1563284) = nox_xxx_getNameId_4E3AA0("GameBall");
	}
	if (!*getMemU32Ptr(0x5D4594, 1563272)) {
		*getMemU32Ptr(0x5D4594, 1563272) = nox_xxx_getNameId_4E3AA0("Crown");
	}
	result = nox_server_getFirstObject_4DA790();
	for (i = result; result; i = result) {
		if (*(uint8_t*)(i + 16) & 4) {
			v3 = *(unsigned short*)(i + 4);
			if ((unsigned short)v3 == *getMemU32Ptr(0x5D4594, 1563284) || v3 == *getMemU32Ptr(0x5D4594, 1563272)) {
				nox_xxx_netMarkMinimapObject_417190(a1, i, 1);
			}
		}
		result = nox_server_getNextObject_4DA7A0(i);
	}
	return result;
}

//----- (004DE410) --------------------------------------------------------
void sub_4DE410(int a1) {
	int v1;           // ebp
	int v2;           // edi
	uint32_t* result; // eax
	uint32_t* i;      // esi
	int v5;           // ecx
	int v6;           // ebx
	unsigned int v7;  // edi
	int v8;           // [esp+10h] [ebp+4h]

	v1 = a1;
	v2 = 1 << a1;
	v8 = 1 << a1;
	result = (uint32_t*)nox_server_getFirstObject_4DA790();
	for (i = result; result; i = result) {
		v5 = i[2];
		i[38] |= v2;
		if (!(v5 & 0x20400000)) {
			i[37] &= ~v2;
		}
		i[v1 + 140] &= 0xFFFu;
		if (i[2] & 0x20400000) {
			v6 = 0x10000;
			v7 = 1;
			do {
				if (sub_4E4C90((int)i, v7)) {
					i[v1 + 140] |= v6;
				}
				v7 *= 2;
				v6 *= 2;
			} while (v7 < 0x10000);
			v2 = v8;
		}
		result = (uint32_t*)nox_server_getNextObject_4DA7A0((int)i);
	}
}

//----- (004DE4D0) --------------------------------------------------------
int sub_4DE4D0(char a1) {
	int v1;     // esi
	int result; // eax
	char v3;    // cl

	v1 = 1 << a1;
	for (result = nox_server_getFirstObject_4DA790(); result; result = nox_server_getNextObject_4DA7A0(result)) {
		v3 = *(uint8_t*)(result + 16);
		*(uint32_t*)(result + 152) |= v1;
		if (!(v3 & 0x20) && !(*(uint32_t*)(result + 8) & 0x20400006)) {
			*(uint32_t*)(result + 148) &= ~v1;
		}
	}
	return result;
}

//----- (004DE7C0) --------------------------------------------------------
void nox_script_event_playerLeave(nox_playerInfo* pl);
int sub_4FF990(unsigned int a1);
void nox_xxx_playerForceDisconnect_4DE7C0(int ind) {
	nox_playerInfo* plr = nox_common_playerInfoFromNum_417090(ind);
	nox_script_event_playerLeave(plr);
	if (sub_4D12A0(ind)) {
		sub_4D1250(ind);
	}
	if (plr->field_2068) {
		int* v3 = sub_425A70(plr->field_2068);
		if (v3) {
			sub_425B60(v3, ind);
		}
	}
	int v4 = *(uint32_t*)((uint32_t)(plr->playerUnit) + 748);
	if (*(uint32_t*)(v4 + 280)) {
		nox_xxx_shopCancelSession_510DC0(*(uint32_t**)(v4 + 280));
	}
	*(uint32_t*)(v4 + 280) = 0;
	sub_510E20(plr->playerInd);
	sub_4FF990(1 << plr->playerInd);

#ifndef NOX_SERVER
	if (!nox_common_gameFlags_check_40A5C0(2))
#endif // NOX_SERVER
	{
		plr->active = 0;
	}

	char* pl = plr;
	sub_56F4F0((int*)pl + 1146);
	sub_56F4F0((int*)pl + 1148);
	sub_56F4F0((int*)pl + 1149);
	sub_56F4F0((int*)pl + 1150);
	sub_56F4F0((int*)pl + 1151);
	sub_56F4F0((int*)pl + 1152);
	sub_56F4F0((int*)pl + 1153);
	sub_56F4F0((int*)pl + 1154);
	sub_56F4F0((int*)pl + 1155);
	sub_56F4F0((int*)pl + 1156);
	sub_56F4F0((int*)pl + 1157);
	sub_56F4F0((int*)pl + 1158);
	sub_56F4F0((int*)pl + 1159);
	sub_56F4F0((int*)pl + 1147);
	sub_56F4F0((int*)pl + 1160);
	sub_56F4F0((int*)pl + 1161);

	char buf[3];
	buf[0] = 46;
	*(uint16_t*)(&buf[1]) = nox_xxx_netGetUnitCodeServ_578AC0(plr->playerUnit);
	nox_xxx_netSendPacket0_4E5420(ind | 0x80, buf, 3, 0, 0);
	nox_xxx_delayedDeleteObject_4E5CC0(plr->playerUnit);
	plr->playerUnit = 0;
	for (int i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
		int v7 = *(uint32_t*)(i + 748);
		*(uint8_t*)(ind + v7 + 452) = 0;
		*(uint32_t*)(v7 + 4 * ind + 324) = 0;
		*(uint8_t*)(ind + v7 + 484) = 0;
		*(uint8_t*)(ind + v7 + 516) = 0;
	}
	if (nox_xxx_gamePlayIsAnyPlayers_40A8A0()) {
		if (nox_common_gameFlags_check_40A5C0(1024) && !nox_xxx_serverIsClosing_446180() && sub_40A770() == 1) {
			sub_5095E0();
		}
	} else {
		sub_40A1F0(0);
		nox_xxx_playerForceSendLessons_416E50(1);
		nox_server_teamsResetYyy_417D00();
		sub_40A970();
	}
	sub_4E55F0(ind);
	nox_xxx_playerResetImportantCtr_4E4F40(ind);
	sub_4E4F30(ind);
	if (0) {
		sub_4161E0();
		sub_416690();
	}
	if (!nox_common_gameFlags_check_40A5C0(4096)) {
		return;
	}
	if (sub_4E9010() == 1) {
		nox_xxx_mapLoad_4D2450(sub_4E8E50());
	} else if (nox_server_questMaybeWarp_4E8F60()) {
		sub_4D60B0();
		sub_4D76E0(1);
		int v12 = nox_game_getQuestStage_4E3CC0();
		int v13 = nox_server_questNextStageThreshold_4D74F0(v12);
		nox_game_setQuestStage_4E3CD0(v13 - 1);
		nox_xxx_mapLoad_4D2450(sub_4E8E50());
	} else {
		char* result = nox_xxx_getFirstPlayerUnit_4DA7C0();
		if (result) {
			while (!*(uint32_t*)(*((uint32_t*)result + 187) + 312)) {
				result = (char*)nox_xxx_getNextPlayerUnit_4DA7F0((int)result);
				if (!result) {
					return;
				}
			}
			result = (char*)sub_4E8E60();
		}
	}
}

//----- (004DEF00) --------------------------------------------------------
int nox_xxx_netGameSettings_4DEF00() {
	char* v0;    // ebx
	char v2[20]; // [esp+Ch] [ebp-48h]
	char v3[49]; // [esp+20h] [ebp-34h]

	v0 = nox_xxx_cliGamedataGet_416590(0);
	v2[0] = -81;
	*(uint32_t*)&v2[1] = gameFrame();
	*(uint32_t*)&v2[9] = nox_common_gameFlags_getVal_40A5B0() & 0x7FFF0;
	*(uint32_t*)&v2[13] = nox_xxx_getServerSubFlags_409E60();
	*(uint32_t*)&v2[5] = NOX_CLIENT_VERS_CODE;
	v2[17] = nox_xxx_servGetPlrLimit_409FA0();
	v2[18] = nox_xxx_servGamedataGet_40A020(*((uint16_t*)v0 + 26));
	v2[19] = sub_40A180(*((uint16_t*)v0 + 26));
	v3[0] = -80;
	strcpy(&v3[1], nox_xxx_serverOptionsGetServername_40A4C0());
	memcpy(&v3[17], v0 + 24, 0x1Cu);
	if (sub_40A220() && (sub_40A300() || sub_40A180(*((uint16_t*)v0 + 26)))) {
		*(uint32_t*)&v3[45] = sub_40A230();
	} else {
		*(uint32_t*)&v3[45] = 0;
	}
	nox_xxx_netSendPacket1_4E5390(159, (int)v2, 20, 0, 0);
	return nox_xxx_netSendPacket1_4E5390(159, (int)v3, 49, 0, 0);
}

//----- (004DF020) --------------------------------------------------------
char* sub_4DF020() {
	char* result;      // eax
	int v1;            // ecx
	unsigned char* v2; // edi
	char* v3;          // esi
	bool v4;           // zf
	char* i;           // esi
	char v6[60];       // [esp+8h] [ebp-3Ch]

	result = sub_459AA0((int)v6);
	v1 = 29;
	v2 = getMemAt(0x5D4594, 1563214);
	v3 = v6;
	v4 = 1;
	do {
		if (!v1) {
			break;
		}
		v4 = *(uint16_t*)v3 == *(uint16_t*)v2;
		v3 += 2;
		v2 += 2;
		--v1;
	} while (v4);
	if (!v4) {
		result = nox_common_playerInfoGetFirst_416EA0();
		for (i = result; result; i = result) {
			if (i[2064] != 31) {
				nox_xxx_netGuiGameSettings_4DD9B0(1, v6, (unsigned char)i[2064]);
			}
			result = nox_common_playerInfoGetNext_416EE0((int)i);
		}
		memcpy(getMemAt(0x5D4594, 1563214), v6, 0x38u);
		*getMemU16Ptr(0x5D4594, 1563270) = *(uint16_t*)&v6[56];
	}
	return result;
}

//----- (004DF0A0) --------------------------------------------------------
int nox_xxx_wallSendDestroyed_4DF0A0(int a1, int a2) {
	int result; // eax
	int i;      // esi
	char v4[3]; // [esp+4h] [ebp-4h]

	v4[0] = 58;
	*(uint16_t*)&v4[1] = *(uint16_t*)(a1 + 10);
	if (a2 != 32) {
		return nox_xxx_netSendPacket0_4E5420(a2, v4, 3, 0, 1);
	}
	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064), v4, 3, 0, 1);
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (004DF120) --------------------------------------------------------
int sub_4DF120(void* a1) {
	int result; // eax
	int i;      // esi
	char v3[3]; // [esp+4h] [ebp-4h]

	v3[0] = 59;
	*(uint16_t*)&v3[1] = *(uint16_t*)((uint32_t)a1 + 10);
	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064), v3, 3, 0, 1);
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (004DF180) --------------------------------------------------------
int sub_4DF180(void* a1) {
	int result; // eax
	int i;      // esi
	char v3[3]; // [esp+4h] [ebp-4h]

	v3[0] = 60;
	*(uint16_t*)&v3[1] = *(uint16_t*)((uint32_t)a1 + 10);
	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064), v3, 3, 0, 1);
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (004DF2E0) --------------------------------------------------------
void sub_4DF2E0(int a1) {
	char* i; // ebx
	int j;   // esi

	if (a1 != 31) {
		for (i = nox_server_teamFirst_418B10(); i; i = nox_server_teamNext_418B60((int)i)) {
			sub_4197C0((wchar2_t*)i, a1);
			for (j = nox_xxx_getFirstPlayerUnit_4DA7C0(); j; j = nox_xxx_getNextPlayerUnit_4DA7F0(j)) {
				if (nox_xxx_teamCompare2_419180(j + 48, i[57])) {
					sub_4198A0(j + 48, a1, *(uint32_t*)(j + 36));
				}
			}
		}
	}
}

//----- (004DF360) --------------------------------------------------------
int nox_xxx_netSendSimpleObject2_4DF360(int a1, nox_object_t* a2p) {
	int a2 = a2p;
	short v2;   // ax
	float v3;   // ecx
	short v4;   // ax
	float v5;   // edx
	char v7[9]; // [esp+4h] [ebp-Ch]

	v7[0] = 47;
	*(uint16_t*)&v7[3] = *(uint16_t*)(a2 + 4);
	v2 = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)a2);
	v3 = *(float*)(a2 + 56);
	*(uint16_t*)&v7[1] = v2;
	v4 = nox_float2int(v3);
	v5 = *(float*)(a2 + 60);
	*(uint16_t*)&v7[5] = v4;
	*(uint16_t*)&v7[7] = nox_float2int(v5);
	return nox_xxx_netSendPacket1_4E5390(a1, (int)v7, 9, 0, 1);
}

//----- (004DF3C0) --------------------------------------------------------
void sub_4DF3C0(nox_playerInfo* pl) {
	int a1 = pl;
	int v1;   // edi
	char* v2; // eax
	char* v3; // ebx
	char* v4; // ebx
	char* v5; // esi
	int v6;   // ebx
	int v7;   // eax
	int v8;   // esi
	int v9;   // ebx

	v1 = *(uint32_t*)(a1 + 2056);
	v2 = sub_416640();
	v3 = v2;
	if (!v1) {
		return;
	}
	if (!sub_40A740() && !nox_common_gameFlags_check_40A5C0(0x8000)) {
		uint8_t v2b = nox_xxx_getTeamCounter_417DD0();
		if (v2b) {
			v2 = sub_4189D0();
			v4 = v2;
			if (v2) {
				v2 = (char*)nox_xxx_servObjectHasTeam_419130(v1 + 48);
				if (!v2) {
					nox_xxx_createAtImpl_4191D0(v4[57], v1 + 48, 1, *(uint32_t*)(v1 + 36), 1);
				}
			}
		}
		return;
	}
	v2 = *(char**)(a1 + 2068);
	if (!v2) {
		return;
	}
	v5 = nox_server_teamByXxx_418AE0(*(uint32_t*)(a1 + 2068));
	if (v5) {
		v6 = v1 + 48;
		v7 = nox_xxx_servObjectHasTeam_419130(v1 + 48);
	} else {
		v8 = (unsigned char)v3[52];
		if ((nox_common_gameFlags_check_40A5C0(96) ||
			 nox_common_gameFlags_check_40A5C0(16) && nox_xxx_CheckGameplayFlags_417DA0(4)) &&
			v8 > 2) {
			v8 = 2;
		}
		v2 = (char*)(unsigned char)sub_417DE0();
		if ((int)v2 >= v8) {
			return;
		}
		if (!nox_common_gameFlags_check_40A5C0(96) ||
			(v9 = (unsigned char)sub_417DE0(), v2 = (char*)sub_417DC0(), v9 < (int)v2)) {
			v2 = sub_418A10();
			v5 = v2;
			if (!v2) {
				return;
			}
			sub_418800((wchar2_t*)v2, (wchar2_t*)(a1 + 2072), 0);
			sub_418830((int)v5, *(uint32_t*)(a1 + 2068));
			sub_4184D0((wchar2_t*)v5);
			v6 = v1 + 48;
			v7 = nox_xxx_servObjectHasTeam_419130(v1 + 48);
		} else {
			return;
		}
	}
	if (v7) {
		sub_4196D0(v6, (int)v5, *(uint32_t*)(v1 + 36), 0);
	} else {
		nox_xxx_createAtImpl_4191D0(v5[57], v6, 1, *(uint32_t*)(v1 + 36), 0);
	}
	return;
}
