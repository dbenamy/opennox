#ifndef NOX_PORT_GAME3
#define NOX_PORT_GAME3

#include "defs.h"
#include "common__savegame.h"

void sub_4A1A40(int a1);
int sub_4A1BE0(int a1);
int nox_client_guiXxxDestroy_4A24A0();
int nox_xxx_wndListboxProcWithoutData10_4A28E0(uint32_t* a1, int a2, unsigned int a3, int a4);
int nox_xxx_wndListboxProcPre_4A30D0(nox_window* win, unsigned int a2, uint32_t a3, int a4);
int sub_4A4800(int a1);
int nox_game_showSelClass_4A4840();
int sub_4A4970();
int sub_4A49A0();
int sub_4A50A0();
int sub_4A50D0();
int sub_4A6890();
int sub_4A6C90();
int sub_4A7A70(int a1);
nox_window* nox_gui_newButtonOrCheckbox_4A91A0(nox_window* parent, int a2, int a3, int a4, int a5, int a6, nox_window_data* draw);
int nox_xxx_compassGenStrings_4A9C80();
int nox_game_showOptions_4AA6B0();
int sub_4AA9C0();
int sub_4AAA10();
int sub_4AB0C0();
int nox_client_mapSpecialRWObjectData_4AC610();
int sub_4AD9B0(int a1);
int sub_4ADA40();
void sub_4AE6F0(int a1, int a2, int a3, int a4, int a5);
long long sub_4AEE30();
void nox_client_drawPoint_4B0BC0(int a1, int a2, int a3);
int nox_xxx_wndScrollBoxDraw_4B4BA0(int a1, int a2, unsigned int a3, int a4);
nox_window* nox_gui_newSlider_4B4EE0(int a1, int a2, int a3, int a4, int a5, int a6, uint32_t* a7, float* a8);
void sub_4B5700(nox_window* a1, void* a2, void* a3, void* a4, void* a5, void* a6);
void sub_4B6720(int2* a1, int a2, int a3, char a4);
uint32_t* nox_xxx_netHandleSummonPacket_4B7C40(short a1, unsigned short* a2, unsigned short a3, unsigned char a4,
											   short a5);
void sub_4B7EE0(short a1);
int nox_xxx_spriteShieldLoad_4B7F90();
uint32_t* nox_xxx_fxShield_4B8090(unsigned int a1, int a2);
void nox_xxx_spriteScanForShield_4B81E0(int a1, int a2);
uint32_t* sub_4B8E10(uint32_t* a1, char* a2);

#endif // NOX_PORT_GAME3
