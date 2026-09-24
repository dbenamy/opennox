#ifndef NOX_PORT_GAME2_1
#define NOX_PORT_GAME2_1

#include "common__savegame.h"
#include "defs.h"

int sub_460EA0(int a1);
int nox_xxx_clientSetAltWeapon_461550(int a1);
int sub_4615C0();
int sub_461600(int a1);
int sub_461930();
unsigned char* sub_461B50();
char* sub_461EF0(int a1);
int nox_xxx_clientEquip_4623B0(int a1);
int sub_4625D0(uint32_t* a1);
int sub_463370(uint32_t* a1, nox_point* pos, uint32_t* a3);
int nox_xxx_guiDrawInventoryTray_4643B0(int a1, int a2);
int sub_464B40(int a1, int a2);
int nox_xxx_clientDequip_464B70(int a1);
int nox_xxx_trade_4657B0(short a1);
int nox_xxx_clientTrade_465870(short a1);
int nox_xxx_clientDrop_465BE0(int2* a1);
int nox_xxx_clientKeyEquip_465C30(int a1, int a2);
void nox_xxx_clientUse_465C70(int a1);
void sub_465CD0(uint32_t* a1, int a2, int a3, int a4);
int sub_465D50_draw(int a1);
int nox_xxx_inventoryDrawProc_466580(uint32_t* a1);
int sub_466620(int a1, int a2, unsigned int a3);
int sub_466F50(uint32_t* a1, int* a2);
void nox_client_invAlterWeapon_4672C0();
int sub_4673F0(int a1, int a2);
int sub_4674A0();
void nox_window_set_visible_unk5(int visible);
int sub_467590();
int sub_467650();
void sub_467680();
nox_window* nox_xxx_wndGetHandle_4676A0();
int sub_467700(int a1);
int sub_467810(int a1, int a2);
char* sub_467870(int a1, int a2);
int sub_4678B0();
int sub_4678C0();
int sub_467BB0();
int sub_467C10();
int sub_467C80();
int nox_xxx_wndRetNULL_46A8A0();
int sub_46AE10(int a1, int a2);
int nox_xxx_wndSetOffsetMB_46AE40(int a1, int a2, int a3);
int nox_xxx_wndSetIcon_46AE60(int a1, int a2);
int nox_xxx_wndSetIconLit_46AEA0(int a1, int a2);
int sub_46AEC0(int a1, int a2);
int sub_46AEE0(int a1, int a2);
wchar2_t* sub_46AF00(void* a1);
int sub_470CC0();
int sub_470CD0();
void sub_470D70();
int sub_470D90(int a1, int a2);
int nox_xxx_cliGetMana_470DD0();
int sub_470E90(int a1, int a2);
void nox_win_init_cur_weapon(nox_window* a1, int a2, int a3, int w, int h);
int sub_470F40_draw(nox_window* win);
int sub_471250(uint32_t* a1);
int sub_471450(uint32_t* a1);
int nox_xxx_guiBottleSlotDrawFn_471A80(uint32_t* a1);
int nox_xxx_guiBottleSlotProc_471B90(int a1, int a2);
int nox_xxx_drawHealthManaBar_471C00(int a1);
int nox_xxx_guiHealthManaTubeProc_472100(int a1, int a2);
int sub_4721A0(int a1);
wchar2_t* sub_472280();
unsigned char* sub_472310();

#endif // NOX_PORT_GAME2_1
