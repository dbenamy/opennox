#include "client__gui__guiinv.h"
#include "client__gui__window.h"

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME5_2.h"
#include "client__drawable__drawable.h"
#include "client__gui__guimsg.h"
#include "client__gui__tooltip.h"
#include "common__magic__speltree.h"
#include "common__object__modifier.h"
#include "common__strman.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_8531A0_2576;
extern uint32_t dword_5d4594_1049844;
extern uint32_t dword_5d4594_1049808;
extern uint32_t dword_5d4594_1062484;
extern uint32_t dword_5d4594_1062556;
extern uint32_t dword_5d4594_1062564;
extern uint32_t dword_5d4594_1062560;
extern uint32_t dword_5d4594_1049804;
extern uint32_t dword_5d4594_1062496;
extern uint32_t dword_5d4594_1063120;
extern uint32_t dword_5d4594_1062488;
extern uint32_t dword_5d4594_1062516;
extern uint32_t dword_5d4594_1062476;
extern uint32_t dword_5d4594_1049856;
extern uint32_t dword_5d4594_1049800_inventory_click_row_index;
extern uint32_t dword_5d4594_1062456;
extern uint32_t dword_5d4594_1063636;
extern uint32_t dword_5d4594_1049796_inventory_click_column_index;
extern uint32_t dword_5d4594_1062512;
extern uint32_t dword_5d4594_1049864;
extern uint32_t dword_5d4594_1063116;
extern uint32_t dword_5d4594_1062480;
extern uint32_t array_5D4594_1049872[9];
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_red_2589776;
extern uint32_t nox_color_blue_2650684;
extern uint32_t nox_color_cyan_2649820;
extern uint32_t nox_color_yellow_2589772;
extern uint32_t nox_color_violet_2598268;
extern uint32_t nox_color_black_2650656;
extern uint32_t nox_color_orange_2614256;

extern nox_inventory_cell_t nox_client_inventory_grid_1050020[NOX_INVENTORY_CELLS_MAX];

//----- (00462740) --------------------------------------------------------
int sub_462740() {
	wchar2_t* v0;  // eax
	uint32_t* v1; // eax

	if (wndIsShown_nox_xxx_wndIsShown_46ACC0(*(int*)&dword_5d4594_1062476)) {
		return 0;
	}
	nox_window_set_hidden(*(int*)&dword_5d4594_1062476, 1);
	dword_5d4594_1063116 = 0;
	dword_5d4594_1063120 = 0;
	v0 = nox_strman_loadString_40F1D0("thing.db:IdentifyDescription", 0, "C:\\NoxPost\\src\\Client\\Gui\\guiinv.c",
									  2361);
	nox_wcscpy((wchar2_t*)getMemAt(0x5D4594, 1063124), v0);
	v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1062476, 9156);
	nox_window_call_field_94((int)v1, 16399, 0, 0);
	nox_xxx_wndClearCaptureMain_46ADE0(*(int*)&dword_5d4594_1062456);
	dword_5d4594_1049864 = 0;
	nox_client_setCursorType_477610(0);
	return 1;
}

//----- (00464BD0) --------------------------------------------------------
int sub_464BD0(int a1, int a2, unsigned int a3) {
	int v4;          // eax
	int v5;          // eax
	int v6;          // eax
	int v7;          // eax
	int v8;          // eax
	int v9;          // esi
	int v10;         // edi
	int v14;         // eax
	int v15;         // eax
	int v16;         // esi
	int v17;         // edi
	int v19;         // ecx
	int v20;         // eax
	int v21;         // eax
	int v26;         // eax
	uint32_t* v28;   // esi
	wchar2_t* v29;    // eax
	int v30;         // eax
	int v31;         // eax
	int v32;         // esi
	int v33;         // edi
	int v34;         // eax
	int v36;         // eax
	int v37;         // eax
	int v38;         // eax
	uint32_t* v40;   // ecx
	int v41;         // edx
	int v42;         // edx
	int v45;         // esi
	int v47;         // eax
	int v48;         // esi
	const void* v49; // edi
	wchar2_t* v50;    // eax
	int2 v51;        // [esp-24h] [ebp-7Ch]
	int v52;         // [esp-1Ch] [ebp-74h]
	int v53;         // [esp-18h] [ebp-70h]
	int v54;         // [esp-8h] [ebp-60h]
	int v55;         // [esp-4h] [ebp-5Ch]
	int2 v56;        // [esp+8h] [ebp-50h]
	int2 v57;        // [esp+10h] [ebp-48h]
	int2 v58;        // [esp+18h] [ebp-40h]
	int2 v59;        // [esp+20h] [ebp-38h]

	v59.field_4 = a3 >> 16;
	v59.field_0 = (unsigned short)a3;
	sub_463370(*(uint32_t**)&dword_5d4594_1062456, &v59, &v56);
	if (sub_45D9B0() || getMemByte(0x5D4594, 1049868) != 2) {
		return 1;
	}
	switch (a2) {
	case 5:
		if (nox_xxx_playerAnimCheck_4372B0()) {
			return 1;
		}
		if (dword_5d4594_1049864 == 5) {
			v8 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136352));
			if (v8) {
				v9 = (v56.field_0 - 314) / 50;
				v10 = (dword_5d4594_1062512 + v56.field_4 - 13) / 50;
				if (!sub_464B40(v9, v10)) {
					return 1;
				}
				int v11 = v10 + NOX_INVENTORY_ROW_COUNT * v9;
				if (nox_client_inventory_grid_1050020[v11].field_140) {
					nox_drawable* dr = nox_client_inventory_grid_1050020[v11].field_0;
					dword_5d4594_1063116 = dr;
					dr->field_32 = nox_client_inventory_grid_1050020[v11].field_4;
				} else {
					dword_5d4594_1063116 = 0;
				}
				return 1;
			}
			if (sub_478030()) {
				if (sub_479870()) {
					LOBYTE(v14) = sub_479880(&v56);
					if (v14) {
						dword_5d4594_1063116 = sub_4798A0(&v56);
						return 1;
					}
				}
			}
		} else if (dword_5d4594_1049864 == 6) {
			v15 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136352));
			if (v15) {
				v16 = (v56.field_0 - 314) / 50;
				v17 = (dword_5d4594_1062512 + v56.field_4 - 13) / 50;
				if (sub_464B40(v16, v17)) {
					int v18 = v17 + NOX_INVENTORY_ROW_COUNT * v16;
					if (nox_client_inventory_grid_1050020[v18].field_140) {
						v19 = nox_client_inventory_grid_1050020[v18].field_0;
						v20 = nox_client_inventory_grid_1050020[v18].field_4;
						*(uint32_t*)(v19 + 128) = v20;
						if (v19) {
							nox_xxx_trade_4657B0(v20);
							return 1;
						}
					}
				}
			}
		} else if (getMemByte(0x5D4594, 1049870) != 1 ||
				   (v21 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136336)), v21 != 1)) {
			if (getMemByte(0x5D4594, 1049869) || nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136384)) ||
				nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136400))) {
				if (nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136384)) == 1) {
					return 0;
				}
				if (nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136400)) == 1) {
					return 0;
				}
			} else {
				nox_xxx_wndSetCaptureMain_46ADC0(*(int*)&dword_5d4594_1062456);
				if (sub_479590() == 3) {
					nox_xxx_clientTradeMB_4657E0(&v56);
				} else {
					sub_4658A0(*(int*)&dword_5d4594_1062456, &v56);
				}
				if (*getMemU32Ptr(0x5D4594, 1049848)) {
					nox_xxx_cursorSetDraggedItem_477690(*getMemIntPtr(0x5D4594, 1049848));
					nox_xxx_setKeybTimeout_4160D0(0);
					*(int2*)getMemAt(0x5D4594, 1062572) = v56;
					nox_xxx_clientPlaySoundSpecial_452D80(791, 100);
					return 1;
				}
			}
		}
		return 1;
	case 7:
		if (nox_xxx_playerAnimCheck_4372B0() || dword_5d4594_1049864 == 6) {
			return 1;
		}
		if (!getMemByte(0x5D4594, 1049869)) {
			v26 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136368));
			if (v26) {
				if ((v56.field_4 - 13) / 50 == 1) {
					if (dword_5d4594_1049864 != 5) {
						sub_465CA0();
						return 1;
					}
					sub_462740();
					return 1;
				}
			}
		}
		// fallthrough
	case 6:
		if (nox_xxx_playerAnimCheck_4372B0() || dword_5d4594_1049864 == 6) {
			return 1;
		}
		int v43 = 0;
		if (dword_5d4594_1049864 == 5) {
			if (nox_xxx_cursorGetTypePrev_477630() == 7) {
				sub_462740();
				return 1;
			}
		} else {
			nox_xxx_wndClearCaptureMain_46ADE0(*(int*)&dword_5d4594_1062456);
		}
		if (dword_5d4594_1049864 == 4) {
			v58 = v59;
			sub_473970(&v58, &v58);
			v28 = nox_drawable_find_49ABF0(&v58, 20);
			if (v28) {
				v57.field_0 = nox_win_width / 2;
				v57.field_4 = nox_win_height / 2;
				sub_473970(&v57, &v57);
				if ((v57.field_0 - v28[3]) * (v57.field_0 - v28[3]) + (v57.field_4 - v28[4]) * (v57.field_4 - v28[4]) <=
					5625) {
					dword_5d4594_1049864 = 0;
					return 1;
				}
				v29 = nox_strman_loadString_40F1D0("ObjectTooFar", 0, "C:\\NoxPost\\src\\Client\\Gui\\guiinv.c", 3858);
			} else {
				v29 = nox_strman_loadString_40F1D0("NoObject", 0, "C:\\NoxPost\\src\\Client\\Gui\\guiinv.c", 3869);
			}
			nox_xxx_printCentered_445490(v29);
			dword_5d4594_1049864 = 0;
			return 1;
		}
		if (!*getMemU32Ptr(0x5D4594, 1049848)) {
			return 1;
		}
		if (!nox_xxx_wndPointInWnd_46AAB0(*(uint32_t**)&dword_5d4594_1062456, v59.field_0, v59.field_4) ||
			(v30 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136384)), v30) ||
			(v31 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136400)), v31)) {
			v58 = v59;
			sub_473970(&v58, &v57);
			if (dword_5d4594_1049856 == 1) {
				if (!sub_4C12C0()) {
					nox_xxx_clientDrop_465BE0(&v57);
				}
			} else {
				v47 = dword_5d4594_1049800_inventory_click_row_index +
					  14 * dword_5d4594_1049796_inventory_click_column_index +
					  7 * dword_5d4594_1049796_inventory_click_column_index;
				v48 = nox_client_inventory_grid_1050020[v47].field_140;
				if (nox_client_inventory_grid_1050020[v47].field_140) {
					v49 = 0;
					nox_xxx_wndClearCaptureMain_46ADE0(*(int*)&dword_5d4594_1062456);
					if (*(uint32_t*)(*getMemU32Ptr(0x5D4594, 1049848) + 112) & 0x13001000) {
						v49 = (const void*)(*getMemU32Ptr(0x5D4594, 1049848) + 432);
					}
					sub_4C05F0(0, 0);
					v53 = *(uint32_t*)(*getMemU32Ptr(0x5D4594, 1049848) + 108);
					v52 = *(uint32_t*)(*getMemU32Ptr(0x5D4594, 1049848) + 128);
					v51 = v58;
					v50 = nox_strman_loadString_40F1D0("DropLabel", 0, "C:\\NoxPost\\src\\Client\\Gui\\guiinv.c", 4148);
					nox_gui_itemAmountDialog_4C0430((int)v50, v51.field_0, v51.field_4, v52, v53, v49, v48 + 1, 0,
													sub_465CD0, 0);
				} else if (!sub_4C12C0()) {
					nox_xxx_clientDrop_465BE0(&v57);
				}
			}
			if (dword_5d4594_1049856) {
				goto LABEL_121;
			}
			v55 = dword_5d4594_1049800_inventory_click_row_index;
			v54 = dword_5d4594_1049796_inventory_click_column_index;
			sub_4649B0(*getMemIntPtr(0x5D4594, 1049848), v54, v55);
			goto LABEL_121;
		}
		v32 = *getMemU32Ptr(0x5D4594, 1062572) - v56.field_0;
		v33 = *getMemU32Ptr(0x5D4594, 1062576) - v56.field_4;
		if (!nox_xxx_checkKeybTimeout_4160F0(0, gameFPS() / 3u) && v32 * v32 + v33 * v33 < 100) {
			v34 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136352));
			if (!v34) {
				goto LABEL_121;
			}
			if (!sub_4C12C0()) {
				if (*(uint32_t*)(*getMemU32Ptr(0x5D4594, 1049848) + 112) & 0x3001000) {
					int v35 = dword_5d4594_1049800_inventory_click_row_index +
							  NOX_INVENTORY_ROW_COUNT * dword_5d4594_1049796_inventory_click_column_index;
					if (nox_client_inventory_grid_1050020[v35].field_136) {
						nox_xxx_clientSetAltWeapon_461550(0);
						nox_client_inventory_grid_1050020[v35].field_136 = 0;
					} else if (nox_client_inventory_grid_1050020[v35].field_132) {
						nox_xxx_clientDequip_464B70(*getMemIntPtr(0x5D4594, 1049848));
					} else {
						nox_xxx_clientKeyEquip_465C30(*(int*)&dword_5d4594_1049796_inventory_click_column_index,
													  *(int*)&dword_5d4594_1049800_inventory_click_row_index);
					}
				} else {
					nox_xxx_clientUse_465C70(*getMemIntPtr(0x5D4594, 1049848));
				}
			}
			sub_4649B0(*getMemIntPtr(0x5D4594, 1049848), *(int*)&dword_5d4594_1049796_inventory_click_column_index,
					   *(int*)&dword_5d4594_1049800_inventory_click_row_index);
			goto LABEL_121;
		}
		v36 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136336));
		if (v36 && !getMemByte(0x5D4594, 1049870)) {
			if (!dword_5d4594_1049856) {
				nox_xxx_clientEquip_4623B0(*getMemIntPtr(0x5D4594, 1049848));
				sub_4649B0(*getMemIntPtr(0x5D4594, 1049848), *(int*)&dword_5d4594_1049796_inventory_click_column_index,
						   *(int*)&dword_5d4594_1049800_inventory_click_row_index);
			}
			goto LABEL_121;
		}
		v37 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136352));
		if (!v37) {
			v55 = dword_5d4594_1049800_inventory_click_row_index;
			v54 = dword_5d4594_1049796_inventory_click_column_index;
			sub_4649B0(*getMemIntPtr(0x5D4594, 1049848), v54, v55);
			goto LABEL_121;
		}
		v38 = *(uint32_t*)(*getMemU32Ptr(0x5D4594, 1049848) + 108);
		if (v38 == dword_5d4594_1062560 || v38 == *getMemU32Ptr(0x5D4594, 1049728) ||
			v38 == *getMemU32Ptr(0x5D4594, 1049724) || v38 == dword_5d4594_1062556 || v38 == dword_5d4594_1062564) {
			sub_4649B0(*getMemIntPtr(0x5D4594, 1049848), *(int*)&dword_5d4594_1049796_inventory_click_column_index,
					   *(int*)&dword_5d4594_1049800_inventory_click_row_index);
			goto LABEL_121;
		}
		dword_5d4594_1049804 = (v56.field_0 - 314) / 50;
		dword_5d4594_1049808 = (dword_5d4594_1062512 + v56.field_4 - 13) / 50;
		if (!sub_464B40((v56.field_0 - 314) / 50, (dword_5d4594_1062512 + v56.field_4 - 13) / 50)) {
			goto LABEL_121;
		}
		if (dword_5d4594_1049856) {
			int v39 = dword_5d4594_1049808 + NOX_INVENTORY_ROW_COUNT * dword_5d4594_1049804;
			if (nox_client_inventory_grid_1050020[v39].field_140 && (v40 = (uint32_t*)nox_client_inventory_grid_1050020[v39].field_0) != 0 &&
				((v41 = v40[28], v41 & 0x2000000) && *(uint32_t*)(*getMemU32Ptr(0x5D4594, 1049848) + 112) & 0x2000000 &&
					 v40[29] == *(uint32_t*)(*getMemU32Ptr(0x5D4594, 1049848) + 116) ||
				 v41 & 0x1001000 && *(uint32_t*)(*getMemU32Ptr(0x5D4594, 1049848) + 112) & 0x1001000)) {
				v42 = nox_client_inventory_grid_1050020[v39].field_4;
				*getMemU32Ptr(0x5D4594, 1049860) = 1;
				v40[32] = v42;
				nox_xxx_clientEquip_4623B0((int)v40);
			} else {
				*getMemU32Ptr(0x5D4594, 1049860) = 1;
				nox_xxx_clientDequip_464B70(*getMemIntPtr(0x5D4594, 1049848));
			}
			goto LABEL_121;
		}
		if (nox_client_inventory_grid_1050020[dword_5d4594_1049800_inventory_click_row_index +
								NOX_INVENTORY_ROW_COUNT * dword_5d4594_1049796_inventory_click_column_index]
				.field_140) {
			v55 = dword_5d4594_1049800_inventory_click_row_index;
			v54 = dword_5d4594_1049796_inventory_click_column_index;
			sub_4649B0(*getMemIntPtr(0x5D4594, 1049848), v54, v55);
			goto LABEL_121;
		}
		if (!sub_4649B0(*getMemIntPtr(0x5D4594, 1049848), *(int*)&dword_5d4594_1049804, *(int*)&dword_5d4594_1049808)) {
			sub_4649B0(*getMemIntPtr(0x5D4594, 1049848), *(int*)&dword_5d4594_1049796_inventory_click_column_index,
					   *(int*)&dword_5d4594_1049800_inventory_click_row_index);
			goto LABEL_121;
		}
		nox_xxx_clientPlaySoundSpecial_452D80(792, 100);
		v43 = dword_5d4594_1049800_inventory_click_row_index +
				  NOX_INVENTORY_ROW_COUNT * dword_5d4594_1049796_inventory_click_column_index;
		v45 = nox_client_inventory_grid_1050020[v43].field_136;
		if (v45) {
			int v46 = dword_5d4594_1049808 + NOX_INVENTORY_ROW_COUNT * dword_5d4594_1049804;
			nox_client_inventory_grid_1050020[v46].field_136 = v45;
			nox_client_inventory_grid_1050020[v43].field_136 = 0;
			dword_5d4594_1062480 = &nox_client_inventory_grid_1050020[v46];
		}
		sub_461B50();
	LABEL_121:
		nox_xxx_cursorResetDraggedItem_4776A0();
		if (!dword_5d4594_1049856) {
			nox_xxx_spriteDelete_45A4B0(*(uint64_t**)getMemAt(0x5D4594, 1049848));
		}
		*getMemU32Ptr(0x5D4594, 1049848) = 0;
		dword_5d4594_1049856 = 0;
		return 1;
	case 8:
		return 1;
	case 9:
		if (dword_5d4594_1049864 == 5) {
			sub_462740();
			return 1;
		}
		return 0;
	case 19:
		if (nox_xxx_playerAnimCheck_4372B0()) {
			return 1;
		}
		v6 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136384));
		if (v6) {
			if (dword_5d4594_1049864 == 5) {
				return 1;
			}
			return 0;
		}
		v7 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136400));
		if (v7) {
			if (dword_5d4594_1049864 == 5) {
				return 1;
			}
			return 0;
		}
		nox_window_call_field_94(*(int*)&dword_5d4594_1062456, 16391, *getMemIntPtr(0x5D4594, 1062500), 0);
		return 1;
	case 20:
		if (nox_xxx_playerAnimCheck_4372B0()) {
			return 1;
		}
		v4 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136384));
		if (v4) {
			if (dword_5d4594_1049864 == 5) {
				return 1;
			}
			return 0;
		}
		v5 = nox_xxx_pointInRect_4281F0(&v56, (int4*)getMemAt(0x587000, 136400));
		if (v5) {
			if (dword_5d4594_1049864 == 5) {
				return 1;
			}
			return 0;
		}
		nox_window_call_field_94(*(int*)&dword_5d4594_1062456, 16391, *getMemIntPtr(0x5D4594, 1062504), 0);
		return 1;
	default:
		if (dword_5d4594_1049864 == 5) {
			return 1;
		}
		return 0;
	}
}
//----- (00466160) --------------------------------------------------------
int sub_466160() {
	wchar2_t* v0; // eax

	if (getMemByte(0x5D4594, 1049868) == 2) {
		v0 = nox_strman_loadString_40F1D0("CloseInventoryTT", 0, "C:\\NoxPost\\src\\Client\\Gui\\guiinv.c", 410);
	} else {
		v0 = nox_strman_loadString_40F1D0("OpenInventoryTT", 0, "C:\\NoxPost\\src\\Client\\Gui\\guiinv.c", 414);
	}
	nox_xxx_cursorSetTooltip_4776B0(v0);
	return 1;
}

//----- (004661D0) --------------------------------------------------------
int sub_4661D0() {
	wchar2_t* v0; // eax
	wchar2_t* v2; // eax

	if (dword_5d4594_1062480) {
		v0 = nox_xxx_clientAskInfoMb_4BF050(**(wchar2_t***)&dword_5d4594_1062480);
		nox_xxx_cursorSetTooltip_4776B0(v0);
	} else {
		v2 = nox_strman_loadString_40F1D0("ToolTipWeapon2Area", 0, "C:\\NoxPost\\src\\Client\\Gui\\guiinv.c", 3331);
		nox_xxx_cursorSetTooltip_4776B0(v2);
	}
	return 1;
}

//----- (004667E0) --------------------------------------------------------
int nox_xxx_inventroryOnHovewerSub_4667E0(int a1, int a2, unsigned int a3) {
	int v3;       // edx
	int v4;       // esi
	int v5;       // ecx
	int v6;       // eax
	wchar2_t* v7;  // eax
	int v9;       // ecx
	wchar2_t* v10; // eax
	int v11;      // eax
	wchar2_t* v12; // eax
	int v13;      // eax
	wchar2_t* v14; // eax
	int2 v15;     // [esp+4h] [ebp-8h]

	v3 = 40;
	v15.field_0 = (unsigned short)a3;
	v15.field_4 = a3 >> 16;
	v4 = 0;
	while (v3 <= (unsigned short)a3) {
		v3 += 35;
		++v4;
	}
	v5 = 0;
	do {
		if ((1 << v5) & *getMemU32Ptr(0x5D4594, 1062540) && v5 != 31) {
			--v4;
		}
		if (v4 < 0) {
			break;
		}
		++v5;
	} while (v5 < 32);
	if (v5 != 32) {
		v6 = nox_xxx_getEnchantSpell_424920(v5);
		v7 = (wchar2_t*)nox_xxx_spellTitle_424930(v6);
		nox_xxx_cursorSetTooltip_4776B0(v7);
		return 1;
	}
	v9 = 0;
	do {
		if ((1 << v9) & getMemByte(0x5D4594, 1062536)) {
			--v4;
		}
		if (v4 < 0) {
			break;
		}
		++v9;
	} while (v9 < 6);
	if (v9 != 6) {
		v10 = sub_413480(1 << v9);
		nox_xxx_cursorSetTooltip_4776B0(v10);
		return 1;
	}
	if (!nox_common_gameFlags_check_40A5C0(4096)) {
		nox_xxx_cursorSetTooltip_4776B0(0);
		return 1;
	}
	v11 = nox_xxx_pointInRect_4281F0(&v15, (int4*)getMemAt(0x5D4594, 1049812));
	if (v11 == 1) {
		v12 = nox_strman_loadString_40F1D0("thing.db:AnkhGUI", 0, "C:\\NoxPost\\src\\Client\\Gui\\guiinv.c", 4385);
		nox_xxx_cursorSetTooltip_4776B0(v12);
		return 1;
	}
	v13 = nox_xxx_pointInRect_4281F0(&v15, (int4*)getMemAt(0x5D4594, 1049828));
	if (v13 == 1 && sub_4BFD30() == 1) {
		v14 = nox_strman_loadString_40F1D0("GeneralPrint:TooltipKeyIcon", 0, "C:\\NoxPost\\src\\Client\\Gui\\guiinv.c",
										   4388);
		nox_xxx_cursorSetTooltip_4776B0(v14);
		return 1;
	} else {
		nox_xxx_cursorSetTooltip_4776B0(0);
		return 1;
	}
}
