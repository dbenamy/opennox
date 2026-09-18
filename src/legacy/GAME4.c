#include <errno.h>
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
#include "GAME5_2.h"
#include "common__binfile.h"
#include "common__crypt.h"
#include "common__net_list.h"
#include "common__random.h"
#include "common__system__team.h"
#include "operators.h"
#include "server__magic__plyrspel.h"
#include "server__magic__spell__execdur.h"
#include "server__script__file.h"
#include "server__script__script.h"

#include "client__gui__window.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"

extern uint32_t dword_5d4594_3835368;
extern uint32_t nox_server_kickQuestPlayerMinVotes_229992;
extern uint32_t nox_server_resetQuestMinVotes_229988;
extern uint32_t dword_5d4594_3835392;
extern uint32_t dword_5d4594_1568868;
extern void* nox_alloc_magicEnt_1569668;
extern void* nox_alloc_vote_1599652;
extern uint32_t dword_5d4594_1569672;
extern uint32_t dword_5d4594_1599656;
extern uint32_t dword_5d4594_2650652;



int nox_cheat_charmall = 0;

int nox_setImaginaryCaster();
int sub_57AEE0(int a1, nox_object_t* a2);
//----- (00500540) --------------------------------------------------------


//----- (005005E0) --------------------------------------------------------


//----- (005006B0) --------------------------------------------------------


//----- (00500750) --------------------------------------------------------


//----- (00500770) --------------------------------------------------------


//----- (00500790) --------------------------------------------------------


//----- (005007E0) --------------------------------------------------------


//----- (005009B0) --------------------------------------------------------


//----- (00500A60) --------------------------------------------------------


//----- (00500B70) --------------------------------------------------------

// 500B70: using guessed type char var_100[256];

//----- (00500C70) --------------------------------------------------------
// Sends information to the player that an unit order happened
int nox_xxx_orderUnitLocal_500C70(int owner, int orderType) {
	*((uint32_t*)nox_common_playerInfoFromNum_417090(owner) + 912) = orderType;
	return nox_xxx_netCreatureCmd_4D7EE0(owner, orderType);
}

nox_object_t* nox_xxx_unitDoSummonAt_5016C0(int a1, float* a2, nox_object_t* a3, unsigned char a4);
//----- (00502670) --------------------------------------------------------


//----- (00502790) --------------------------------------------------------


//----- (005029A0) --------------------------------------------------------


//----- (005029F0) --------------------------------------------------------


//----- (00502A20) --------------------------------------------------------


//----- (00502A50) --------------------------------------------------------


//----- (00502AB0) --------------------------------------------------------


//----- (00502B10) --------------------------------------------------------

// 502B10: using guessed type char var_40[64];

//----- (00502D70) --------------------------------------------------------


//----- (00502DA0) --------------------------------------------------------


//----- (00502DF0) --------------------------------------------------------


//----- (00502E10) --------------------------------------------------------


//----- (00502E70) --------------------------------------------------------


//----- (00502EA0) --------------------------------------------------------


//----- (00503830) --------------------------------------------------------


//----- (00503B30) --------------------------------------------------------
void nox_script_readWriteZzz_541670(char* path, char* path2, char* dst);
void nox_xxx_waypoint_5799C0();
void sub_579D20();
nox_waypoint_t* sub_579890();


//----- (00503EC0) --------------------------------------------------------


//----- (005040A0) --------------------------------------------------------


//----- (00504150) --------------------------------------------------------
extern uint32_t nox_tile_def_cnt;
extern nox_tileDef_t nox_tile_defs_arr[176];


//----- (00504290) --------------------------------------------------------


//----- (005042F0) --------------------------------------------------------


//----- (00504330) --------------------------------------------------------


//----- (005044B0) --------------------------------------------------------


//----- (00504560) --------------------------------------------------------
void sub_579E90(nox_waypoint_t* a1);


//----- (005048A0) --------------------------------------------------------


//----- (00504910) --------------------------------------------------------


//----- (00504980) --------------------------------------------------------


//----- (005049C0) --------------------------------------------------------


//----- (005049D0) --------------------------------------------------------


//----- (005049E0) --------------------------------------------------------


//----- (00504A10) --------------------------------------------------------


//----- (00505060) --------------------------------------------------------


//----- (00505080) --------------------------------------------------------

// 505080: using guessed type char var_400[1024];

//----- (00505C30) --------------------------------------------------------
int sub_57C130(uint32_t* a1, int a2);
int nox_server_mapLoadAddGroup_57C0C0(char* a1, unsigned int a2, unsigned char a3);

// 505C30: using guessed type char var_E4[76];

//----- (00506260) --------------------------------------------------------
nox_waypoint_t* nox_xxx_waypointNewNotMap_579970(int a1, float a2, float a3);

// 506260: using guessed type char var_4C[76];

//----- (005066D0) --------------------------------------------------------
int nox_xxx_allocVoteArray_5066D0() {
	int result; // eax

	result = nox_new_alloc_class("VoteClass", 52, 64);
	nox_alloc_vote_1599652 = result;
	if (result) {
		dword_5d4594_1599656 = 0;
		result = 1;
	}
	return result;
}

//----- (00506720) --------------------------------------------------------
int sub_506720() {
	int result; // eax

	nox_free_alloc_class(*(void**)&nox_alloc_vote_1599652);
	result = 0;
	nox_alloc_vote_1599652 = 0;
	dword_5d4594_1599656 = 0;
	return result;
}

//----- (00506740) --------------------------------------------------------
int sub_506740(nox_object_t* a1p) {
	int a1 = a1p;
	int result; // eax
	int v2;     // esi
	int v3;     // ecx
	int v4;     // edi

	result = a1;
	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 4) {
			result = dword_5d4594_1599656;
			v2 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
			if (dword_5d4594_1599656) {
				do {
					v3 = *(uint32_t*)(result + 8);
					v4 = *(uint32_t*)(result + 44);
					if (v3 & v2) {
						*(uint32_t*)(result + 8) = ~v2 & v3;
						--*(uint8_t*)(result + 4);
					}
					if (!*(uint8_t*)(result + 4)) {
						sub_5067B0(result);
					}
					result = v4;
				} while (v4);
			}
		}
	}
	return result;
}

//----- (005067B0) --------------------------------------------------------
void sub_5067B0(int a1) {
	int v1; // esi

	if (a1) {
		if (*(uint32_t*)a1 == 2) {
			v1 = 0;
			do {
				if ((1 << v1) & *(uint32_t*)(a1 + 8)) {
					nox_xxx_netSendVote_506840(v1);
				}
				++v1;
			} while (v1 < 32);
		}
		sub_506810(a1);
		nox_alloc_class_free_obj_first(*(unsigned int**)&nox_alloc_vote_1599652, (uint64_t*)a1);
		if (!dword_5d4594_1599656) {
			sub_507190(255, 0);
		}
	}
}

//----- (00506810) --------------------------------------------------------
int sub_506810(int a1) {
	int result; // eax
	int v2;     // ecx
	int v3;     // ecx

	result = a1;
	v2 = *(uint32_t*)(a1 + 44);
	if (v2) {
		*(uint32_t*)(v2 + 48) = *(uint32_t*)(a1 + 48);
	}
	v3 = *(uint32_t*)(a1 + 48);
	if (v3) {
		result = *(uint32_t*)(a1 + 44);
		*(uint32_t*)(v3 + 44) = result;
	} else {
		dword_5d4594_1599656 = *(uint32_t*)(a1 + 44);
	}
	return result;
}

//----- (00506840) --------------------------------------------------------
int nox_xxx_netSendVote_506840(int a1) {
	char v2[2]; // [esp+0h] [ebp-2h]

	v2[0] = -18;
	v2[1] = 7;
	return nox_xxx_netSendPacket1_4E5390(a1, (int)v2, 2, 0, 1);
}

//----- (00506870) --------------------------------------------------------
char sub_506870(int a1, int a2, wchar2_t* a3) {
	char result; // al

	result = a2;
	if (a2 && *(uint8_t*)(a2 + 8) & 4) {
		switch (a1) {
		case 0:
			result = sub_5068E0(0, a2, a3);
			break;
		case 1:
			result = sub_5068E0(1, a2, a3);
			break;
		case 2:
			result = (unsigned int)sub_506B00(2, a2);
			break;
		case 3:
			result = (unsigned int)sub_506B80(3, a2, a3);
			break;
		default:
			return result;
		}
	}
	return result;
}

//----- (005068E0) --------------------------------------------------------
char sub_5068E0(int a1, int a2, wchar2_t* a3) {
	int v3; // eax
	int v4; // ebp
	int v5; // esi
	int v6; // edi
	int v7; // esi

	LOBYTE(v3) = getMemByte(0x587000, 229980);
	if (*getMemU32Ptr(0x587000, 229980) > 0x20u) {
		return v3;
	}
	if (*getMemU32Ptr(0x587000, 229980) == 0) {
		return v3;
	}
	if (!a3) {
		return v3;
	}
	if (!a2) {
		return v3;
	}
	v4 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a2 + 748) + 276) + 2064);
	v3 = nox_common_playerInfoGetFirst_416EA0();
	v5 = v3;
	if (!v3) {
		return v3;
	}
	while (1) {
		if (*(uint32_t*)(v5 + 2092) == 1) {
			v3 = nox_wcscmp((const wchar2_t*)(v5 + 4704), a3);
			if (!v3) {
				break;
			}
		}
		v3 = nox_common_playerInfoGetNext_416EE0(v5);
		v5 = v3;
		if (!v3) {
			return v3;
		}
	}
	if (*(uint8_t*)(v5 + 2064) == 31) {
		return v3;
	}
	v6 = *(uint32_t*)(v5 + 2056);
	if (!v6) {
		return v3;
	}
	if (a2 == v6) {
		return v3;
	}
	if (!nox_xxx_CheckGameplayFlags_417DA0(4) || (v3 = nox_xxx_servCompareTeams_419150(a2 + 48, v6 + 48)) != 0) {
		v7 = dword_5d4594_1599656;
		if (dword_5d4594_1599656) {
			while (*(uint32_t*)v7 != a1 || *(uint32_t*)(v7 + 28) != v6) {
				v7 = *(uint32_t*)(v7 + 44);
				if (!v7) {
					break;
				}
			}
		}
		if (!v7) {
			v3 = sub_506A20(a1, a2);
			v7 = v3;
			if (!v3) {
				return v3;
			}
			*(uint32_t*)(v3 + 28) = v6;
			if (nox_xxx_CheckGameplayFlags_417DA0(4)) {
				*(uint32_t*)(v7 + 20) = 1;
			}
		}
		v3 = *(uint32_t*)(v7 + 8);
		if (!(v4 & v3)) {
			LOBYTE(v3) = *(uint8_t*)(v7 + 4) + 1;
			*(uint32_t*)(v7 + 8) |= v4;
			*(uint8_t*)(v7 + 4) = v3;
		}
	}
	return v3;
}

//----- (00506A20) --------------------------------------------------------
uint32_t* sub_506A20(int a1, int a2) {
	int v2;       // ebx
	uint32_t* v3; // esi

	v2 = 0;
	if (!a2 || !(*(uint8_t*)(a2 + 8) & 4)) {
		return 0;
	}
	if (!dword_5d4594_1599656) {
		v2 = 1;
	}
	v3 = nox_alloc_class_new_obj_zero(*(uint32_t**)&nox_alloc_vote_1599652);
	if (!v3) {
		return 0;
	}
	*v3 = a1;
	v3[6] = gameFrame();
	v3[4] = a2 + 48;
	switch (a1) {
	case 0:
	case 1:
		*((uint8_t*)v3 + 12) = getMemByte(0x587000, 229980);
		break;
	case 2:
	case 3:
		*((uint8_t*)v3 + 12) = 6;
		break;
	default:
		*((uint8_t*)v3 + 12) = getMemByte(0x587000, 229984);
		break;
	}
	nox_xxx_voteAddMB_506AD0((int)v3);
	if (v2) {
		sub_507190(255, 1);
	}
	return v3;
}

//----- (00506AD0) --------------------------------------------------------
int nox_xxx_voteAddMB_506AD0(int a1) {
	int result; // eax

	result = a1;
	*(uint32_t*)(a1 + 48) = 0;
	*(uint32_t*)(a1 + 44) = dword_5d4594_1599656;
	if (dword_5d4594_1599656) {
		*(uint32_t*)(dword_5d4594_1599656 + 48) = a1;
	}
	dword_5d4594_1599656 = a1;
	return result;
}

//----- (00506B00) --------------------------------------------------------
uint32_t* sub_506B00(int a1, int a2) {
	uint32_t* result; // eax
	int v3;           // esi
	char v4;          // cl

	result = *(uint32_t**)&nox_server_resetQuestMinVotes_229988;
	if (nox_server_resetQuestMinVotes_229988) {
		if (a2) {
			result = *(uint32_t**)(*(uint32_t*)(a2 + 748) + 276);
			v3 = 1 << *((uint8_t*)result + 2064);
			if (result[1198]) {
				result = *(uint32_t**)&dword_5d4594_1599656;
				if (dword_5d4594_1599656) {
					while (*result != a1) {
						result = (uint32_t*)result[11];
						if (!result) {
							break;
						}
					}
				}
				if (!result) {
					result = sub_506A20(a1, a2);
					if (!result) {
						return result;
					}
					result[5] = 0;
				}
				if (!(result[2] & v3)) {
					v4 = *((uint8_t*)result + 4) + 1;
					result[2] |= v3;
					*((uint8_t*)result + 4) = v4;
				}
			}
		}
	}
	return result;
}

//----- (00506B80) --------------------------------------------------------
uint32_t* sub_506B80(int a1, int a2, wchar2_t* a3) {
	uint32_t* result;  // eax
	int v4;            // edi
	const wchar2_t* v5; // esi
	int v6;            // esi
	char v7;           // cl

	result = *(uint32_t**)&nox_server_kickQuestPlayerMinVotes_229992;
	if (nox_server_kickQuestPlayerMinVotes_229992) {
		if (a3) {
			result = (uint32_t*)a2;
			if (a2) {
				result = *(uint32_t**)(*(uint32_t*)(a2 + 748) + 276);
				v4 = 1 << *((uint8_t*)result + 2064);
				if (result[1198]) {
					result = nox_common_playerInfoGetFirst_416EA0();
					v5 = (const wchar2_t*)result;
					if (result) {
						while (1) {
							if (*((uint32_t*)v5 + 523) == 1) {
								result = (uint32_t*)nox_wcscmp(v5 + 2352, a3);
								if (!result) {
									break;
								}
							}
							result = nox_common_playerInfoGetNext_416EE0((int)v5);
							v5 = (const wchar2_t*)result;
							if (!result) {
								return result;
							}
						}
						if (*((uint8_t*)v5 + 2064) != 31) {
							result = (uint32_t*)*((uint32_t*)v5 + 1198);
							if (result) {
								v6 = *((uint32_t*)v5 + 514);
								if (v6) {
									if (a2 != v6) {
										result = *(uint32_t**)&dword_5d4594_1599656;
										if (dword_5d4594_1599656) {
											while (*result != a1 || result[7] != v6) {
												result = (uint32_t*)result[11];
												if (!result) {
													break;
												}
											}
										}
										if (!result) {
											result = sub_506A20(a1, a2);
											if (!result) {
												return result;
											}
											result[7] = v6;
										}
										if (!(result[2] & v4)) {
											v7 = *((uint8_t*)result + 4) + 1;
											result[2] |= v4;
											*((uint8_t*)result + 4) = v7;
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return result;
}

//----- (00506C90) --------------------------------------------------------
void sub_506C90(int a1, int a2, wchar2_t* a3) {
	if (a2 && *(uint8_t*)(a2 + 8) & 4) {
		switch (a1) {
		case 0:
			sub_506D00(a2, a3);
			break;
		case 1:
			sub_506D00(a2, a3);
			break;
		case 2:
			sub_506DE0(a2);
			break;
		case 3:
			sub_506E50(a2, a3);
			break;
		default:
			return;
		}
	}
}

//----- (00506D00) --------------------------------------------------------
void sub_506D00(int a1, wchar2_t* a2) {
	char* v2; // esi
	int v3;   // esi
	int v4;   // eax
	int v5;   // edx
	int v6;   // esi
	bool v7;  // zf

	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 4) {
				v2 = nox_common_playerInfoGetFirst_416EA0();
				if (v2) {
					while (*((uint32_t*)v2 + 523) != 1 || nox_wcscmp((const wchar2_t*)v2 + 2352, a2)) {
						v2 = nox_common_playerInfoGetNext_416EE0((int)v2);
						if (!v2) {
							return;
						}
					}
					if (v2[2064] != 31) {
						v3 = *((uint32_t*)v2 + 514);
						if (v3) {
							v4 = dword_5d4594_1599656;
							v5 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
							if (dword_5d4594_1599656) {
								while (*(uint32_t*)v4 || *(uint32_t*)(v4 + 28) != v3 || !(v5 & *(uint32_t*)(v4 + 8))) {
									v4 = *(uint32_t*)(v4 + 44);
									if (!v4) {
										return;
									}
								}
								if (v4) {
									v6 = ~v5 & *(uint32_t*)(v4 + 8);
									v7 = (*(uint8_t*)(v4 + 4))-- == 1;
									*(uint32_t*)(v4 + 8) = v6;
									if (v7) {
										sub_5067B0(v4);
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

//----- (00506DE0) --------------------------------------------------------
void sub_506DE0(int a1) {
	int result;
	int v2;  // edx
	char v3; // cl
	int v4;  // esi

	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 4) {
			result = dword_5d4594_1599656;
			v2 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
			if (dword_5d4594_1599656) {
				while (*(uint32_t*)result != 2) {
					result = *(uint32_t*)(result + 44);
					if (!result) {
						return;
					}
				}
				if (result && v2 & *(uint32_t*)(result + 8)) {
					v3 = *(uint8_t*)(result + 4) - 1;
					v4 = ~v2 & *(uint32_t*)(result + 8);
					*(uint8_t*)(result + 4) = v3;
					*(uint32_t*)(result + 8) = v4;
					if (!v3) {
						sub_5067B0(result);
					}
				}
			}
		}
	}
}

//----- (00506E50) --------------------------------------------------------
void sub_506E50(int a1, wchar2_t* a2) {
	char* v2; // esi
	int v3;   // esi
	int v4;   // eax
	int v5;   // edx
	int v6;   // esi
	bool v7;  // zf

	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 4) {
				v2 = nox_common_playerInfoGetFirst_416EA0();
				if (v2) {
					while (*((uint32_t*)v2 + 523) != 1 || nox_wcscmp((const wchar2_t*)v2 + 2352, a2)) {
						v2 = nox_common_playerInfoGetNext_416EE0((int)v2);
						if (!v2) {
							return;
						}
					}
					if (v2[2064] != 31) {
						v3 = *((uint32_t*)v2 + 514);
						if (v3) {
							v4 = dword_5d4594_1599656;
							v5 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
							if (dword_5d4594_1599656) {
								while (*(uint32_t*)v4 != 3 || *(uint32_t*)(v4 + 28) != v3 ||
									   !(v5 & *(uint32_t*)(v4 + 8))) {
									v4 = *(uint32_t*)(v4 + 44);
									if (!v4) {
										return;
									}
								}
								if (v4) {
									v6 = ~v5 & *(uint32_t*)(v4 + 8);
									v7 = (*(uint8_t*)(v4 + 4))-- == 1;
									*(uint32_t*)(v4 + 8) = v6;
									if (v7) {
										sub_5067B0(v4);
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

//----- (00506F80) --------------------------------------------------------
void sub_506F80(int a1) {
	int v1; // esi
	int v3; // esi

	v1 = *(uint32_t*)(a1 + 28);
	if (*(uint8_t*)(v1 + 16) & 0x20) {
		sub_5067B0(a1);
		return;
	}
	*(uint32_t*)(a1 + 16) = v1 + 48;
	if (sub_507000(a1) == 1) {
		v3 = *(uint32_t*)(v1 + 748);
		nox_xxx_playerCallDisconnect_4DEAB0(*(unsigned char*)(*(uint32_t*)(v3 + 276) + 2064), 4);
		sub_416770(15, (wchar2_t*)(*(uint32_t*)(v3 + 276) + 4704), (const char*)(*(uint32_t*)(v3 + 276) + 2112));
		sub_5067B0(a1);
	}
}

//----- (00507000) --------------------------------------------------------
int sub_507000(int a1) {
	int v1; // edi
	int i;  // esi
	int j;  // eax

	v1 = 0;
	if (*(uint8_t*)(a1 + 4) >= *(uint8_t*)(a1 + 12)) {
		return 1;
	}
	if (*(uint32_t*)(a1 + 20) == 1) {
		for (i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
			if (nox_xxx_servCompareTeams_419150(*(uint32_t*)(a1 + 16), i + 48)) {
				++v1;
			}
		}
	} else {
		for (j = nox_xxx_getFirstPlayerUnit_4DA7C0(); j; j = nox_xxx_getNextPlayerUnit_4DA7F0(j)) {
			++v1;
		}
	}
	return *(unsigned char*)(a1 + 4) >= (unsigned int)(v1 - 1) && *(uint8_t*)(a1 + 4) >= 2u;
}

//----- (00507090) --------------------------------------------------------
void sub_507090(int a1) {
	int i;    // esi
	int v3;   // eax
	char* v4; // eax

	nox_xxx_player_4E3CE0();
	if (*(unsigned char*)(a1 + 4) >= nox_xxx_player_4E3CE0()) {
		for (i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
			v3 = *(uint32_t*)(*(uint32_t*)(i + 748) + 276);
			if (*(uint32_t*)(v3 + 4792) == 1) {
				nox_xxx_playerRespawn_4F7EF0(*(uint32_t*)(v3 + 2056));
			}
		}
		nox_game_setQuestStage_4E3CD0(0);
		v4 = nox_xxx_getQuestMapFile_4D0F60();
		nox_xxx_mapLoad_4D2450(v4);
		sub_5067B0(a1);
	}
}

//----- (00507100) --------------------------------------------------------
void sub_507100(int a1) {
	int v1;              // edi
	int v2;              // ebx
	unsigned int v3;     // eax
	unsigned int result; // eax

	v1 = *(uint32_t*)(a1 + 28);
	if (!v1) {
		sub_5067B0(a1);
		return;
	}
	if (*(uint8_t*)(v1 + 16) & 0x20) {
		sub_5067B0(a1);
		return;
	}
	v2 = *(uint32_t*)(v1 + 748);
	if (!*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4792)) {
		sub_5067B0(a1);
		return;
	}
	if (*(uint8_t*)(a1 + 4) >= *(uint8_t*)(a1 + 12)) {
		goto LABEL_8;
	}
	v3 = nox_xxx_player_4E3CE0();
	if (v3 <= 1) {
		sub_5067B0(a1);
		return;
	}
	result = v3 - 1;
	if (!(*(unsigned char*)(a1 + 4) >= result && *(uint8_t*)(a1 + 4) >= 2u)) {
		return;
	}
LABEL_8:
	sub_4DCFB0(v1);
	sub_416770(15, (wchar2_t*)(*(uint32_t*)(v2 + 276) + 4704), (const char*)(*(uint32_t*)(v2 + 276) + 2112));
	sub_5067B0(a1);
	return;
}

//----- (00507190) --------------------------------------------------------
int sub_507190(int a1, char a2) {
	char v4[3]; // [esp+0h] [ebp-4h]
	v4[0] = -18;
	v4[1] = 6;
	v4[2] = a2;
	return nox_xxx_netSendPacket1_4E5390(a1, (int)v4, 3, 0, 1);
}

//----- (005071C0) --------------------------------------------------------
int sub_5071C0() { return dword_5d4594_1599656 != 0; }

//----- (00509120) --------------------------------------------------------
