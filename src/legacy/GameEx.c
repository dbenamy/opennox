#define DG_DYNARR_IMPLEMENTATION
#include "GameEx.h"
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_3.h"
#include "GAME3_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME5_2.h"
#include "client__gui__guimsg.h"
#include "client__gui__window.h"
#include "common/fs/nox_fs.h"
#include "memmap.h"

#include "client__shell__noxworld.h"
#include "operators.h"


//-------------------------------------------------------------------------
// Data declarations

unsigned int gameex_flags = 0x1E;

enum { GAMEEX_PLAYER_NAME_UNITS = sizeof(((nox_playerInfo*)0)->name_final) / sizeof(wchar2_t) };

// Player names have a fixed-width UTF-16 field. Preserve the old narrow-byte
// conversion, with a terminator even when that entire field is occupied.
int nox_CharToOemW(const wchar2_t* pSrc, char* pDst) {
	int n = 0;
	while (n < GAMEEX_PLAYER_NAME_UNITS && pSrc[n]) {
		pDst[n] = (char)pSrc[n];
		++n;
	}
	pDst[n] = 0;
	return n;
}

//----- (10001C20) --------------------------------------------------------

//----- (10001D40) --------------------------------------------------------
char playerInfoStructParser_0(void* a1p) {
	char* a1 = a1p;
	char* v1;  // esi
	char pDst[GAMEEX_PLAYER_NAME_UNITS + 1];

	if (!a1 || a1 == (char*)-2)
		return 0;
	v1 = nox_common_playerInfoGetFirst_416EA0();
	if (!v1)
		return 0;
	while (1) {
		nox_CharToOemW((const wchar2_t*)v1 + 2352, pDst);
		if (!strcmp(pDst, a1 + 2))
			break;
		v1 = nox_common_playerInfoGetNext_416EE0((int)v1);
		if (!v1)
			return 0;
	}
	void* team = nox_xxx_objGetTeamByNetCode_418C80(*((uint32_t*)v1 + 515));
	if (!team)
		return 0;
	a1[1] = *((uint8_t*)team + 4);
	*a1 = v1[2251];
	return 1;
}

//----- (10001E10) --------------------------------------------------------
char playerInfoStructParser_1(void* a1p, int* a3) {
	int a1 = a1p;
	char* v3;     // eax
	char* v4;     // eax
	uint32_t* v6; // eax
	char pDst[GAMEEX_PLAYER_NAME_UNITS + 1];

	if (!a1 || a1 == -2 || !a3)
		return 0;
	v3 = nox_common_playerInfoGetFirst_416EA0();
	int a2 = v3;
	if (!v3)
		return 0;
	while (1) {
		nox_CharToOemW((const wchar2_t*)(a2 + 4704), pDst);
		if (!strcmp(pDst, (const char*)(a1 + 2)))
			break;
		v4 = nox_common_playerInfoGetNext_416EE0(a2);
		a2 = v4;
		if (!v4)
			return 0;
	}
	v6 = nox_xxx_objGetTeamByNetCode_418C80(*(uint32_t*)(a2 + 2060));
	if (!v6)
		return 0;
	*a3 = (int)v6;
	*(uint8_t*)(a1 + 1) = *((uint8_t*)v6 + 4);
	*(uint8_t*)a1 = *(uint8_t*)(a2 + 2251);
	return 1;
}

//----- (10001EE0) --------------------------------------------------------
char mix_MouseKeyboardWeaponRoll(nox_object_t* playerObjP, char a2) {
	int playerObj = playerObjP;
	int v2;        // esi
	signed int v4; // edi
	int i;         // esi MAPDST
	int v6;        // eax
	int v8;        // eax
	int weapFlags; // eax
	int v11;       // eax
	char v16;      // [esp+Fh] [ebp-1h]

	v16 = 0;
	v2 = *(uint32_t*)(playerObj + 748);
	if (!(*(uint8_t*)(*(uint32_t*)(v2 + 276) + 3680) & 3) && *(uint8_t*)(v2 + 88) != 1) {
		int cur = *(uint32_t*)(v2 + 104);
		v4 = 496;
		if (!a2)
			v4 = 500;
		if (cur) {
			int next = cur;
			while (1) {
				next = *(uint32_t*)(next + v4);
				if (!next)
					break;
				weapFlags = nox_xxx_weaponInventoryEquipFlags_415820(next); // weaponEquipFlags
				if (weapFlags && weapFlags != 2) {
					if (nox_xxx_playerClassCanUseItem_57B3D0(
							next, *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(playerObj + 748) + 276) + 2251))) {
						v11 = nox_xxx_playerCheckStrength_4F3180(playerObj, next);
						if (v11) {
							if (nox_xxx_playerTryDequip_4F2FB0((uint32_t*)playerObj, (const nox_object_t*)cur) &&
								nox_xxx_playerTryEquip_4F2F70((uint32_t*)playerObj, next)) {
								v16 = 1;
							}
							return v16;
						}
					}
				}
			}
		} else {
			for (i = *(uint32_t*)(playerObj + 504); i; i = *(uint32_t*)(i + 496)) {
				v6 = nox_xxx_weaponInventoryEquipFlags_415820(i); // weaponEquipFlags
				if (v6 && v6 != 2) {
					if (nox_xxx_playerClassCanUseItem_57B3D0(
							i, *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(playerObj + 0x2EC) + 0x114) + 0x8CB))) {
						v8 = nox_xxx_playerCheckStrength_4F3180(playerObj, i);
						if (v8) {
							if (nox_xxx_playerTryEquip_4F2F70((uint32_t*)playerObj, i)) {
								v16 = 1;
							}
							return v16;
						}
					}
				}
			}
		}
	}
	return v16;
}
/* Orphan comments:
playerCanUseItem
playerCheckStrength
*/

//----- (10002030) --------------------------------------------------------
char playerDropATrap(int playerObj) {
	int v2; // eax
	int i;  // esi
	// int v5; // [esp+10h] [ebp-14h]
	char v7; // [esp+18h] [ebp-Ch]
	char v8; // [esp+1Fh] [ebp-5h]
	float pos[2] = {0};

	v7 = 17;
	if (!playerObj)
		return 0;
	v8 = 0;
	v2 = *(uint32_t*)(*(uint32_t*)(playerObj + 0x2EC) + 0x114);
	pos[0] = *(float*)(v2 + 0xE30);
	pos[1] = *(float*)(v2 + 0xE34);
	if (!(*(uint8_t*)(*(uint32_t*)(*(uint32_t*)(playerObj + 0x2EC) + 0x114) + 0xE60) &
		  3) // check playerGameStatus/isObs
		&& *(uint8_t*)(*(uint32_t*)(playerObj + 0x2EC) + 0x58) != 1) {
		for (i = *(uint32_t*)(playerObj + 0x1F8); i; i = *(uint32_t*)(i + 0x1F0)) {
			if (*(uint8_t*)(i + 0xA) == v7) // check if something from *(byte*)(unit+0xA)=17
			{
				nox_xxx_drop_4ED810(playerObj, i, pos); // drop this item
				return 1;
			}
		}
	}
	return v8;
}


//----- (10004330) --------------------------------------------------------
int getFlagValueFromFlagIndex(signed int a1) {
	signed int v1;   // edx
	unsigned int v2; // eax
	signed int v3;   // ecx
	int result;      // eax

	v1 = 2;
	v2 = a1;
	if (a1 < 0)
		v2 = -a1;
	v3 = 1;
	while (1) {
		if (v2 & 1)
			v3 *= v1;
		v2 >>= 1;
		if (!v2)
			break;
		v1 *= v1;
	}
	if (a1 >= 0)
		result = v3;
	else
		result = 1 / v3;
	return result;
}
