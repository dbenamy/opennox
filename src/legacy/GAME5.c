#include <float.h>
#include <math.h>
#include <stdio.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_3.h"
#include "GAME2_2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "common__random.h"
#include "common__crypt.h"
#include "operators.h"
#include "server__dbase__objdb.h"
#include "server__script__script.h"
#include "common__system__team.h"

#include "client__gui__window.h"
#include "client__video__draw_common.h"
extern uint32_t dword_5d4594_2491716;
extern uint32_t dword_5d4594_2491580;
extern uint32_t dword_5d4594_2491676;
extern uint32_t dword_5d4594_2491588;
extern uint32_t dword_5d4594_2491592;
extern uint32_t dword_5d4594_2491704;
extern uint32_t dword_5d4594_2489460;
extern uint32_t dword_5d4594_2650652;
extern uint32_t nox_player_netCode_85319C;

//----- (00548CD0) --------------------------------------------------------
void nox_xxx_script_forcedialog_548CD0(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	int v2; // ecx
	int v3; // eax

	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 4) {
				if (*(uint8_t*)(a2 + 8) & 2) {
					if (!(*(uint32_t*)(a1 + 16) & 0x8020)) {
						v2 = *(uint32_t*)(a2 + 748);
						if (*(uint8_t*)(a2 + 20) & 0x10) {
							v3 = *(uint32_t*)(v2 + 2096);
							if (v3 != -1 && *(int*)(v2 + 2100) != -1) {
								nox_script_callByIndex_507310(v3, a1, a2);
							}
						}
					}
				}
			}
		}
	}
}

//----- (00548D30) --------------------------------------------------------
void nox_xxx_scriptDialog_548D30(nox_object_t* a1p, char a2) {
	int a1 = a1p;
	int v2; // ebx
	int v3; // edi
	int v5; // ebp
	int v6; // esi

	v2 = a1;
	v3 = *(uint32_t*)(a1 + 748);
	nox_xxx_unitUnFreeze_4E7A60(a1, 0);
	v5 = *(uint32_t*)(v3 + 284);
	if (v5) {
		v6 = *(uint32_t*)(v5 + 748);
		if (*(int*)(v6 + 2096) != -1 && *(int*)(v6 + 2100) != -1) {
			LOWORD(a1) = 1232;
			nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(*(uint32_t*)(v3 + 276) + 2064), &a1, 2, 0, 1);
			*(uint32_t*)(v3 + 284) = 0;
			if (*(uint8_t*)(v6 + 2104) == 1) {
				*(uint8_t*)(v6 + 2105) = a2;
			} else {
				*(uint8_t*)(v6 + 2105) = 0;
			}
			nox_script_callByIndex_507310(*(uint32_t*)(v6 + 2100), v2, v5);
		}
	}
}

//----- (0054AF40) --------------------------------------------------------
nox_object_t* nox_xxx_findObjectAtCursor_54AF40(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;     // eax
	double v2;  // st7
	float2 a1a; // [esp+0h] [ebp-8h]

	v1 = *(uint32_t*)(*(uint32_t*)(a1 + 748) + 276);
	a1a.field_0 = (double)*(int*)(v1 + 2284);
	v2 = (double)*(int*)(v1 + 2288);
	dword_5d4594_2491592 = a1;
	a1a.field_4 = v2;
	*getMemU32Ptr(0x5D4594, 2491596) = 0;
	*getMemU32Ptr(0x5D4594, 2491600) = 0;
	nox_xxx_unitsGetInCircle_517F90(&a1a, 100.0, (int)nox_xxx_playerCursorScanFn_54AFB0, (int)&a1a);
	return *getMemU32Ptr(0x5D4594, 2491596);
}

//----- (0054AFB0) --------------------------------------------------------
void nox_xxx_playerCursorScanFn_54AFB0(int a1, float* a2) {
	float* v2;  // esi
	char* v3;   // eax
	int v4;     // eax
	double v5;  // st7
	float v6;   // ecx
	int v7;     // eax
	double v8;  // st7
	double v9;  // st6
	double v10; // st7
	double v11; // st6
	double v12; // st7
	float v13;  // [esp+0h] [ebp-28h]
	float v14;  // [esp+10h] [ebp-18h]
	float2 a3;  // [esp+18h] [ebp-10h]
	float2 a1a; // [esp+20h] [ebp-8h]
	float v17;  // [esp+2Ch] [ebp+4h]

	if (!*getMemU32Ptr(0x5D4594, 2491604)) {
		*getMemU32Ptr(0x5D4594, 2491604) = nox_xxx_getNameId_4E3AA0("Polyp");
	}
	v2 = (float*)a1;
	if (a1 != dword_5d4594_2491592 && !(*(uint32_t*)(a1 + 16) & 0x8020) &&
		(!nox_xxx_testUnitBuffs_4FF350(a1, 0) || nox_xxx_testUnitBuffs_4FF350(*(int*)&dword_5d4594_2491592, 21)) &&
		(*(uint32_t*)(a1 + 8) & 0x80000206 || *(unsigned short*)(a1 + 4) == *getMemU32Ptr(0x5D4594, 2491604))) {
		if (nox_xxx_mapCheck_537110(a1, *(int*)&dword_5d4594_2491592)) {
			if (!(*(uint8_t*)(a1 + 8) & 4) ||
				(*(uint32_t*)(a1 + 36) != nox_player_netCode_85319C ||
				 !nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING)) &&
					(v3 = nox_common_playerInfoGetByID_417040(*(uint32_t*)(a1 + 36))) != 0 && !(v3[3680] & 1)) {
				if ((*(uint32_t*)(a1 + 8) & 0x200) != 512 || (v4 = *(uint32_t*)(a1 + 16), BYTE1(v4) & 0x40)) {
					v5 = *(float*)(a1 + 60) - *(float*)(a1 + 104);
					v14 = v5;
					v17 = v5;
					if (*((uint32_t*)v2 + 43) == 2) {
						v13 = v2[44] * v2[44];
						v7 = nox_float2int(v13);
						a3.field_0 = v2[14] - v2[44];
						a1a.field_0 = v2[44] + v2[14];
						if (a2[1] > (double)v17) {
							v8 = *a2 - v2[14];
							v9 = a2[1] - v17;
							if ((double)v7 <= v9 * v9 + v8 * v8) {
								return;
							}
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
						v10 = *a2;
						if (a2[1] < (double)v14) {
							v11 = a2[1] - v14;
							if ((double)v7 <= v11 * v11 + (v10 - v2[14]) * (v10 - v2[14])) {
								return;
							}
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
						if (v10 > a3.field_0 && *a2 < (double)a1a.field_0) {
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
					} else if (*((uint32_t*)v2 + 43) == 3) {
						v6 = v2[14];
						a1a.field_4 = v5;
						a3.field_0 = *a2;
						a1a.field_0 = v6;
						a3.field_4 = a2[1];
						if (nox_xxx_map_57B850(&a1a, v2 + 43, &a3)) {
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
						a1a.field_4 = v5;
						if (nox_xxx_map_57B850(&a1a, v2 + 43, &a3) || v2[50] + v2[14] < *a2 && *a2 < (double)v2[14] &&
																		  v14 + v2[51] < a2[1] &&
																		  v17 + v2[51] > a2[1]) {
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
						if (*a2 >= (double)v2[14] && v2[52] + v2[14] > *a2 && v14 + v2[53] < a2[1] &&
							v17 + v2[53] > a2[1]) {
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
					}
				}
			}
		}
	}
}

//----- (0054D2B0) --------------------------------------------------------
int nox_xxx_diePlayer_54D2B0(int a1) {
	int v1;                // edi
	int v2;                // ebp
	int v3;                // esi
	uint32_t* v4;          // eax
	char* v5;              // eax
	int v6;                // ebx
	char v7;               // cl
	int v8;                // eax
	int v9;                // edx
	short v10;             // ax
	short v11;             // dx
	int v12;               // eax
	int v13;               // edx
	short v14;             // ax
	int v15;               // edx
	int v16;               // eax
	int v17;               // edx
	int result;            // eax
	int v19;               // eax
	int v20;               // eax
	float v21;             // [esp+0h] [ebp-28h]
	char* v22;             // [esp+14h] [ebp-14h]
	unsigned char v23[14]; // [esp+18h] [ebp-10h]
	int v24;               // [esp+2Ch] [ebp+4h]

	v1 = a1;
	v2 = 0;
	v3 = *(uint32_t*)(a1 + 748);
	if (!*getMemU32Ptr(0x5D4594, 2491688)) {
		*getMemU32Ptr(0x5D4594, 2491688) = nox_xxx_getNameId_4E3AA0("AnkhTradable");
	}
	if (nox_common_gameFlags_check_40A5C0(2048)) {
		sub_4DB170(0, 0, 0);
	}
	v24 = *(uint32_t*)(a1 + 520);
	if (v24) {
		v24 = nox_xxx_findParentChainPlayer_4EC580(v24);
	}
	v4 = *(uint32_t**)(v3 + 276);
	if (v4[900] && gameFrame() - v4[902] < (unsigned int)(10 * gameFPS())) {
		v5 = nox_common_playerInfoFromNum_417090(v4[901]);
		v6 = (int)v5;
		v22 = v5;
		if (v5) {
			if (*((uint32_t*)v5 + 523) && *((uint32_t*)v5 + 514)) {
				v2 = nox_server_getObjectFromNetCode_4ECCB0(*((uint32_t*)v5 + 515));
			} else {
				v22 = 0;
				v2 = 0;
				v6 = 0;
			}
			if (v24 == v2) {
				v2 = 0;
			}
			if (v1 == v2) {
				v2 = 0;
			}
		}
	} else {
		v22 = 0;
		v6 = 0;
	}
	if (!nox_common_gameFlags_check_40A5C0(0x2000)) {
		goto LABEL_38;
	}
	v7 = 0;
	v23[10] = 0;
	*(uint32_t*)v23 = 0;
	*(uint32_t*)&v23[4] = 0;
	*(uint16_t*)&v23[8] = 0;
	if (v24 && *(uint8_t*)(v24 + 8) & 4) {
		*(uint16_t*)&v23[2] = *(uint16_t*)(v24 + 36);
	}
	v8 = *(uint32_t*)(v1 + 520);
	if (v8) {
		v9 = *(uint32_t*)(v8 + 8);
		if (v9 & 2) {
			v23[10] = 1;
			v10 = *(uint16_t*)(v8 + 4);
			*(uint16_t*)&v23[8] = v10;
			goto LABEL_31;
		}
		if (v9 & 4) {
			v7 = *(uint8_t*)(v3 + 304);
			v11 = *(uint16_t*)(v3 + 300);
			v23[10] = *(uint8_t*)(v3 + 304);
			*(uint16_t*)&v23[8] = v11;
		} else {
			v12 = *(uint32_t*)(v8 + 508);
			if (v12) {
				v13 = *(uint32_t*)(v12 + 8);
				if (v13 & 4) {
					v7 = *(uint8_t*)(v3 + 304);
					v14 = *(uint16_t*)(v3 + 300);
					v23[10] = *(uint8_t*)(v3 + 304);
					*(uint16_t*)&v23[8] = v14;
				} else if (v13 & 2) {
					v23[10] = 1;
					*(uint16_t*)&v23[8] = *(uint16_t*)(v12 + 4);
					goto LABEL_31;
				}
			}
		}
	}
	if (!v7) {
		v10 = *(uint16_t*)(v3 + 300);
		v23[10] = *(uint8_t*)(v3 + 304);
		*(uint16_t*)&v23[8] = v10;
	}
LABEL_31:
	if (v2 && *(uint8_t*)(v2 + 8) & 4) {
		*(uint16_t*)&v23[4] = *(uint16_t*)(v2 + 36);
	}
	*(uint16_t*)&v23[6] = *(uint16_t*)(v1 + 36);
	nox_xxx_netInformTextMsg2_4DA180(14, v23);
	if (v23[10] == 2 && *(uint16_t*)&v23[8] == 2) {
		sub_4FC0B0(v24, 1);
	}
	v6 = (int)v22;
	*(uint32_t*)(v3 + 304) = 0;
LABEL_38:
	if (*(uint32_t*)(v1 + 524) == 16) {
		nox_xxx_aud_501960(299, v1, 0, 0);
	} else if (*(uint8_t*)(*(uint32_t*)(v3 + 276) + 2252)) {
		nox_xxx_aud_501960(331, v1, 0, 0);
	} else {
		nox_xxx_aud_501960(321, v1, 0, 0);
	}
	v15 = *(uint32_t*)(v1 + 16);
	BYTE1(v15) |= 0x80u;
	*(uint32_t*)(v1 + 16) = v15;
	nox_xxx_playerSetState_4FA020((uint32_t*)v1, 3);
	*(uint8_t*)(v3 + 188) = 0;
	*(uint32_t*)(v3 + 216) = 0;
	*(uint32_t*)(v3 + 192) = 0;
	*(uint32_t*)(v3 + 196) = 0;
	*(uint32_t*)(v3 + 200) = 0;
	*(uint32_t*)(v3 + 204) = 0;
	*(uint32_t*)(v3 + 208) = 0;
	*(uint8_t*)(v3 + 212) = 0;
	v16 = nox_xxx_gamePlayIsAnyPlayers_40A8A0();
	if (v16) {
		if (nox_common_gameFlags_check_40A5C0(256)) {
			nox_xxx_playerUpdateScore_54D980(v1, v24, v2, v6);
		} else if (nox_common_gameFlags_check_40A5C0(16)) {
			nox_xxx_playerHandleKotrDeath_54DC40(v1, v24);
		} else if (nox_common_gameFlags_check_40A5C0(1024)) {
			nox_xxx_playerHandleElimDeath_54D7A0(v1, v24);
		}
	}
	if (nox_common_gameFlags_check_40A5C0(1024) && nox_xxx_servGamedataGet_40A020(1024) &&
		*(uint32_t*)(*(uint32_t*)(v3 + 276) + 2140) >= (int)(unsigned short)nox_xxx_servGamedataGet_40A020(1024)) {
		nox_xxx_playerRemoveSpawnedStuff_4E5AD0(v1);
	}
	*(uint32_t*)(v1 + 16) |= 0x10u;
	nox_xxx_action_4DA9F0((uint32_t*)v1);
	if (!nox_common_gameFlags_check_40A5C0(4096)) {
		nox_xxx_dropAllItems_4EDA40((uint32_t*)v1);
	}
	nox_xxx_netNotifyPlayerDied_54DF00(v1);
	v17 = *(uint32_t*)(v3 + 276);
	*(uint16_t*)(v3 + 4) = 0;
	nox_xxx_protectMana_56F9E0(*(uint32_t*)(v17 + 4596), 0);
	nox_xxx_setUnitBuffFlags_4E48F0(v1, 0);
	nox_xxx_playerCancelAbils_4FC180(v1);
	*(uint32_t*)(*(uint32_t*)(v3 + 276) + 3600) = 0;
	nox_xxx_playerCancelSpells_4FEAE0(v1);
	nox_xxx_unitClearBuffs_4FF580(v1);
	if (*(uint32_t*)(v3 + 280)) {
		nox_xxx_shopCancelSession_510DC0(*(uint32_t**)(v3 + 280));
	}
	*(uint32_t*)(v3 + 280) = 0;
	result = nox_common_gameFlags_check_40A5C0(4096);
	if (result) {
		v19 = *(uint32_t*)(v3 + 320);
		if (v19) {
			*(uint32_t*)(v3 + 320) = v19 - 1;
			result = sub_4D6130(v1);
		} else {
			v20 = *(uint32_t*)(v3 + 276);
			*(uint32_t*)(v3 + 548) = gameFrame();
			v23[0] = -16;
			v23[1] = 2;
			*(uint16_t*)&v23[8] = *(uint16_t*)(v20 + 4688);
			*(uint16_t*)&v23[2] = *(uint16_t*)(v20 + 4668);
			*(uint16_t*)&v23[6] = *(uint16_t*)(v20 + 4664);
			*(uint16_t*)&v23[4] = *(uint16_t*)(v20 + 4672);
			*(uint32_t*)&v23[10] = 0;
			nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(v20 + 2064), v23, 14, 0, 1);
			sub_4D6000(v1);
			sub_54CBD0(v1);
			v21 = nox_xxx_gamedataGetFloat_419D40("QuestGameStartingExtraLives");
			*(uint32_t*)(v3 + 320) = nox_float2int(v21);
			result = *(uint32_t*)(v3 + 276);
			*(uint8_t*)(*(unsigned char*)(result + 2064) + v3 + 452) = *(uint8_t*)(v3 + 320);
		}
	}
	return result;
}
// 54D56D: variable 'v16' is possibly undefined

//----- (0054D7A0) --------------------------------------------------------
void nox_xxx_playerHandleElimDeath_54D7A0(int a1, int a2) {
	int v2;   // edi
	int v3;   // ebx
	char* v4; // ebp
	int v5;   // eax
	int v6;   // [esp-4h] [ebp-18h]
	int v7;   // [esp-4h] [ebp-18h]
	char* v8; // [esp+10h] [ebp-4h]
	int v9;   // [esp+18h] [ebp+4h]

	v2 = a1;
	v3 = 0;
	v4 = 0;
	v8 = 0;
	v6 = a1 + 48;
	v9 = *(uint32_t*)(a1 + 748);
	if (nox_xxx_servObjectHasTeam_419130(v6)) {
		v8 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(v2 + 52));
	}
	if (a2) {
		v3 = *(uint32_t*)(a2 + 748);
		if (nox_xxx_servObjectHasTeam_419130(a2 + 48)) {
			v4 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a2 + 52));
		}
	}
	if (a2 == v2) {
		nox_xxx_playerSubLessons_4D8EC0(v2, 1);
		nox_xxx_playerIncrementElimDeath_4D8D40(v2);
		nox_xxx_netReportLesson_4D8EF0(v2);
		if (v8) {
			nox_xxx_netChangeTeamID_419090((int)v8, *((uint32_t*)v8 + 13) + 1);
		}
		if (dword_5d4594_2650652) {
			if (v3) {
				sub_425CA0(*(uint32_t*)(v3 + 276), *(uint32_t*)(v3 + 276));
			}
		}
		return;
	}
	if (a2) {
		if (*(uint8_t*)(a2 + 8) & 4) {
			if (v4) {
				if (v4 == v8) {
					nox_xxx_playerSubLessons_4D8EC0(a2, 1);
					nox_xxx_netReportLesson_4D8EF0(a2);
					if (dword_5d4594_2650652 && v3) {
						sub_425CA0(*(uint32_t*)(v3 + 276), *(uint32_t*)(v3 + 276));
					}
					goto LABEL_32;
				}
			} else if (!v8) {
				nox_xxx_changeScore_4D8E90(a2, 1);
				nox_xxx_netReportLesson_4D8EF0(a2);
				if (dword_5d4594_2650652 && v3 && v9) {
					sub_425CA0(*(uint32_t*)(v3 + 276), *(uint32_t*)(v9 + 276));
				}
				goto LABEL_32;
			}
			nox_xxx_changeScore_4D8E90(a2, 1);
			nox_xxx_netReportLesson_4D8EF0(a2);
			if (dword_5d4594_2650652) {
				if (v3 && v9) {
					v5 = *(uint32_t*)(v3 + 276);
					v7 = *(uint32_t*)(v9 + 276);
					sub_425CA0(v5, v7);
					goto LABEL_32;
				}
			}
		}
	} else if (dword_5d4594_2650652 && v9) {
		v5 = *(uint32_t*)(v9 + 276);
		v7 = *(uint32_t*)(v9 + 276);
		sub_425CA0(v5, v7);
	}
LABEL_32:
	nox_xxx_playerIncrementElimDeath_4D8D40(v2);
	nox_xxx_netReportLesson_4D8EF0(v2);
	if (v8) {
		nox_xxx_netChangeTeamID_419090((int)v8, *((uint32_t*)v8 + 13) + 1);
	}
}

//----- (0054D980) --------------------------------------------------------
void nox_xxx_playerUpdateScore_54D980(int a1, int a2, int a3, int a4) {
	int v4;       // ebx
	char* v5;     // edi
	int v6;       // ebp
	int v7;       // eax
	char* result; // eax
	char* v9;     // esi
	int v10;      // edi
	int v11;      // ecx
	int v12;      // [esp-4h] [ebp-20h]
	int v13;      // [esp-4h] [ebp-20h]
	char* v14;    // [esp+10h] [ebp-Ch]
	char* v15;    // [esp+14h] [ebp-8h]
	char* v16;    // [esp+18h] [ebp-4h]
	int v17;      // [esp+20h] [ebp+4h]

	v4 = a1;
	v5 = 0;
	v12 = a1 + 48;
	v14 = 0;
	v15 = 0;
	v17 = *(uint32_t*)(a1 + 748);
	v6 = 0;
	v16 = 0;
	if (nox_xxx_servObjectHasTeam_419130(v12)) {
		v14 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(v4 + 52));
	}
	if (a2) {
		v6 = *(uint32_t*)(a2 + 748);
		if (nox_xxx_servObjectHasTeam_419130(a2 + 48)) {
			v5 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a2 + 52));
		}
	}
	if (a4) {
		if (a3) {
			v16 = *(char**)(a3 + 748);
			if (nox_xxx_servObjectHasTeam_419130(a3 + 48)) {
				v15 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a3 + 52));
			}
		}
	}
	if (a2 == v4) {
		goto LABEL_31;
	}
	if (a2) {
		if (!(*(uint8_t*)(a2 + 8) & 4)) {
			nox_xxx_playerIncrementElimDeath_4D8D40(v4);
			result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
			v9 = v14;
			goto LABEL_36;
		}
		if (v5) {
			if (v5 == v14) {
				nox_xxx_playerSubLessons_4D8EC0(a2, 1);
				nox_xxx_netReportLesson_4D8EF0(a2);
				nox_xxx_netChangeTeamID_419090((int)v5, *((uint32_t*)v5 + 13) - 1);
				if (!dword_5d4594_2650652 || !v6) {
					nox_xxx_playerIncrementElimDeath_4D8D40(v4);
					result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
					v9 = v14;
					goto LABEL_36;
				}
				v7 = *(uint32_t*)(v6 + 276);
				v13 = *(uint32_t*)(v6 + 276);
				sub_425CA0(v7, v13);
				nox_xxx_playerIncrementElimDeath_4D8D40(v4);
				result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
				v9 = v14;
				goto LABEL_36;
			}
		} else if (!v14) {
			nox_xxx_changeScore_4D8E90(a2, 1);
			nox_xxx_netReportLesson_4D8EF0(a2);
			if (!dword_5d4594_2650652 || !v6 || !v17) {
				nox_xxx_playerIncrementElimDeath_4D8D40(v4);
				result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
				v9 = v14;
				goto LABEL_36;
			}
			v7 = *(uint32_t*)(v6 + 276);
			v13 = *(uint32_t*)(v17 + 276);
			sub_425CA0(v7, v13);
			nox_xxx_playerIncrementElimDeath_4D8D40(v4);
			result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
			v9 = v14;
			goto LABEL_36;
		}
		nox_xxx_changeScore_4D8E90(a2, 1);
		nox_xxx_netReportLesson_4D8EF0(a2);
		nox_xxx_netChangeTeamID_419090((int)v5, *((uint32_t*)v5 + 13) + 1);
		if (dword_5d4594_2650652 && v6 && v17) {
			sub_425CA0(*(uint32_t*)(v6 + 276), *(uint32_t*)(v17 + 276));
		}
		nox_xxx_playerIncrementElimDeath_4D8D40(v4);
		result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
		v9 = v14;
		goto LABEL_36;
	}
	if (a3) {
		nox_xxx_playerIncrementElimDeath_4D8D40(v4);
		result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
		v9 = v14;
		goto LABEL_36;
	}
LABEL_31:
	nox_xxx_playerSubLessons_4D8EC0(v4, 1);
	nox_xxx_netReportLesson_4D8EF0(v4);
	v9 = v14;
	if (v14) {
		nox_xxx_netChangeTeamID_419090((int)v14, *((uint32_t*)v14 + 13) - 1);
	}
	result = *(char**)&dword_5d4594_2650652;
	if (dword_5d4594_2650652 && v6) {
		result = sub_425CA0(*(uint32_t*)(v6 + 276), *(uint32_t*)(v6 + 276));
	}
LABEL_36:
	if (!a3) {
		return;
	}
	if (v5) {
		result = v15;
		if (v5 == v15) {
			return;
		}
		v10 = (int)v15;
	} else {
		v10 = (int)v15;
		if (!v15) {
			goto LABEL_44;
		}
	}
	if (v9) {
		if (v9 == (char*)v10) {
			return;
		}
	}
	if (v10) {
		nox_xxx_changeScore_4D8E90(a3, 1);
		nox_xxx_netReportLesson_4D8EF0(a3);
		if (v10) {
			nox_xxx_netChangeTeamID_419090(v10, *(uint32_t*)(v10 + 52) + 1);
		}
		result = *(char**)&dword_5d4594_2650652;
		if (dword_5d4594_2650652) {
			result = v16;
			if (v16) {
				v11 = v17;
				if (v17) {
					sub_425CA0(*((uint32_t*)result + 69), *(uint32_t*)(v11 + 276));
					return;
				}
			}
		}
		return;
	}
LABEL_44:
	nox_xxx_changeScore_4D8E90(a3, 1);
	nox_xxx_netReportLesson_4D8EF0(a3);
	if (dword_5d4594_2650652) {
		result = v16;
		if (v16) {
			v11 = v17;
			if (v17) {
				sub_425CA0(*((uint32_t*)result + 69), *(uint32_t*)(v11 + 276));
				return;
			}
		}
	}
}

//----- (0054DC40) --------------------------------------------------------
void nox_xxx_playerHandleKotrDeath_54DC40(int a1, int a2) {
	char* v2;     // edi
	char* v3;     // ebx
	char* result; // eax
	int v5;       // ebp
	double v6;    // st7
	int v7;       // ebx
	int v8;       // ebx
	int v9;       // eax
	double v10;   // st7
	int v11;      // eax
	float v12;    // [esp+0h] [ebp-18h]
	float v13;    // [esp+0h] [ebp-18h]
	int v14;      // [esp+0h] [ebp-18h]
	float v15;    // [esp+0h] [ebp-18h]
	int v16;      // [esp+14h] [ebp-4h]

	v2 = 0;
	v3 = 0;
	v16 = *(uint32_t*)(a1 + 748);
	result = (char*)nox_xxx_servObjectHasTeam_419130(a1 + 48);
	if (result) {
		result = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a1 + 52));
		v3 = result;
	}
	if (a2) {
		v5 = *(uint32_t*)(a2 + 748);
		result = (char*)nox_xxx_servObjectHasTeam_419130(a2 + 48);
		if (result) {
			result = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a2 + 52));
			v2 = result;
		}
		if (*(uint8_t*)(a2 + 8) & 4) {
			if (a2 == a1 || v3 == v2 && v3) {
				if (!nox_xxx_unitIsCrown_4E7BE0(a2)) {
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				nox_xxx_playerSubLessons_4D8EC0(a2, 1);
				nox_xxx_netReportLesson_4D8EF0(a2);
				if (v2) {
					nox_xxx_netChangeTeamID_419090((int)v2, *((uint32_t*)v2 + 13) - 1);
				}
				if (!dword_5d4594_2650652 || !v5) {
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				v9 = *(uint32_t*)(v5 + 276);
				v14 = *(uint32_t*)(v5 + 276);
			} else {
				if (!v2 || v2 == v3) {
					if (nox_xxx_unitIsCrown_4E7BE0(a2) || nox_xxx_unitIsCrown_4E7BE0(a1)) {
						if (nox_xxx_unitIsCrown_4E7BE0(a2)) {
							v10 = nox_xxx_gamedataGetFloat_419D40("KotRKingKillsPawnPoints");
						} else {
							v10 = nox_xxx_gamedataGetFloat_419D40("KotRPawnKillsKingPoints");
						}
						v15 = v10;
						v11 = nox_float2int(v15);
						nox_xxx_changeScore_4D8E90(a2, v11);
						nox_xxx_netReportLesson_4D8EF0(a2);
						if (dword_5d4594_2650652 && v5 && v16) {
							sub_425CA0(*(uint32_t*)(v5 + 276), *(uint32_t*)(v16 + 276));
						}
						if (!nox_xxx_CheckGameplayFlags_417DA0(4) && nox_xxx_unitIsCrown_4E7BE0(a1)) {
							sub_4ED050(a1, a2);
						}
					}
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				if (nox_xxx_unitIsCrown_4E7BE0(a2)) {
					if (nox_xxx_unitIsCrown_4E7BE0(a1)) {
						v6 = nox_xxx_gamedataGetFloat_419D40("KotRKingKillsKingPoints");
					} else {
						v6 = nox_xxx_gamedataGetFloat_419D40("KotRKingKillsPawnPoints");
					}
					v12 = v6;
					v7 = nox_float2int(v12);
					nox_xxx_changeScore_4D8E90(a2, v7);
					nox_xxx_netChangeTeamID_419090((int)v2, v7 + *((uint32_t*)v2 + 13));
					nox_xxx_netReportLesson_4D8EF0(a2);
					if (dword_5d4594_2650652 && v5) {
						if (v16) {
							sub_425CA0(*(uint32_t*)(v5 + 276), *(uint32_t*)(v16 + 276));
						}
					}
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				if (!nox_xxx_unitIsCrown_4E7BE0(a1) ||
					(v13 = nox_xxx_gamedataGetFloat_419D40("KotRPawnKillsKingPoints"), v8 = nox_float2int(v13),
					 nox_xxx_changeScore_4D8E90(a2, v8),
					 nox_xxx_netChangeTeamID_419090((int)v2, v8 + *((uint32_t*)v2 + 13)),
					 nox_xxx_netReportLesson_4D8EF0(a2), !dword_5d4594_2650652) ||
					!v5 || !v16) {
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				v9 = *(uint32_t*)(v5 + 276);
				v14 = *(uint32_t*)(v16 + 276);
			}
			sub_425CA0(v9, v14);
			nox_xxx_playerIncrementElimDeath_4D8D40(a1);
			nox_xxx_netReportLesson_4D8EF0(a1);
			return;
		}
	}
}

//----- (0054DF00) --------------------------------------------------------
void nox_xxx_netNotifyPlayerDied_54DF00(int a1) {
	int v1 = 0; // ecx
	short v2;   // cx
	int v4;     // [esp+0h] [ebp-4h]

	v4 = v1;
	v2 = *(uint16_t*)(a1 + 36);
	LOBYTE(v4) = -24;
	*(uint16_t*)((char*)&v4 + 1) = v2;
	nox_xxx_netSendPacket1_4E5390(255, (int)&v4, 3, 0, 0);
}

//----- (0054E6F0) --------------------------------------------------------
int sub_54E6F0(int a1, int a2) {
	int result; // eax

	result = sub_54E730(a2, a1);
	if (result) {
		result = !nox_xxx_unitsHaveSameTeam_4EC520(a1, a2) || nox_xxx_GetGameplayFlags_417D90() & 1;
	}
	return result;
}

//----- (0054E730) --------------------------------------------------------
int sub_54E730(int a1, int a2) {
	int v2;     // ecx
	int v3;     // eax
	int result; // eax
	int v5;     // eax

	if (*(uint8_t*)(a2 + 8) & 1) {
		return 0;
	}
	v2 = *(uint32_t*)(a1 + 16);
	if (v2 & 0x20) {
		return 0;
	}
	v3 = *(uint32_t*)(a2 + 16);
	if (v3 & 0x20 || !*(uint32_t*)(a1 + 696) || !*(uint32_t*)(a2 + 696) || v3 & 0x40) {
		return 0;
	}
	if ((v3 & 0x80u) != 0) {
		return 1;
	}
	if (v2 & 0x11 && v3 & 0x4000 || v3 & 0x11 && v2 & 0x4000 ||
		(v2 & 0x400 || v3 & 0x400) && nox_xxx_unitsHaveSameTeam_4EC520(a2, a1) ||
		(v5 = *(uint32_t*)(a1 + 508)) != 0 && *(uint8_t*)(a1 + 8) & 1 && !(*(uint8_t*)(a1 + 12) & 2) &&
			*(uint8_t*)(v5 + 8) & 2 && *(uint8_t*)(a2 + 8) & 2 &&
			(!nox_xxx_unitIsEnemyTo_5330C0(v5, a2) || nox_xxx_unitsHaveSameTeam_4EC520(a2, *(uint32_t*)(a1 + 508)))) {
		return 0;
	} else {
		return 1;
	}
}

//----- (0054E810) --------------------------------------------------------
int sub_54E810(int a1, float2* a2, int a3) {
	int a3a[4]; // [esp+0h] [ebp-10h]

	a3a[0] = a1;
	a3a[1] = 0;
	a3a[2] = a3;
	a3a[3] = (int)a2;
	sub_517B70(a2, sub_54E850, (int)a3a);
	return a3a[1];
}

//----- (0054E850) --------------------------------------------------------
void sub_54E850(int a1, int a2) {
	int v2;       // eax
	float2* v3;   // ecx
	float2* v4;   // edx
	float v5;     // edx
	uint32_t* v6; // edx
	uint32_t* v7; // eax
	float4 a2a;   // [esp+8h] [ebp-20h]
	float4 a1a;   // [esp+18h] [ebp-10h]

	if ((signed char)*(uint8_t*)(a1 + 8) >= 0) {
		if (sub_54E730(*(uint32_t*)a2, a1) && sub_547DB0(a1, *(float2**)(a2 + 12))) {
			*(uint32_t*)(a2 + 4) = a1;
		}
	} else {
		v2 = *(uint32_t*)(a1 + 748);
		if (!(*(uint8_t*)(a1 + 12) & 4)) {
			v3 = *(float2**)(a2 + 8);
			a1a.field_0 = v3->field_0;
			v4 = *(float2**)(a2 + 12);
			a1a.field_4 = v3->field_4;
			a1a.field_8 = v4->field_0;
			v5 = v4->field_4;
			a2a.field_0 = *(float*)(a1 + 56);
			a1a.field_C = v5;
			a2a.field_4 = *(float*)(a1 + 60);
			a2a.field_8 = (double)*getMemIntPtr(0x587000, 196184 + 8 * *(uint32_t*)(v2 + 12)) + a2a.field_0;
			a2a.field_C = (double)*getMemIntPtr(0x587000, 196188 + 8 * *(uint32_t*)(v2 + 12)) + a2a.field_4;
			if (sub_427980(&a1a, &a2a)) {
				v6 = *(uint32_t**)(a2 + 8);
				v7 = *(uint32_t**)(a2 + 12);
				*v7 = *v6;
				v7[1] = v6[1];
				*(uint32_t*)(a2 + 4) = a1;
			}
		}
	}
}
