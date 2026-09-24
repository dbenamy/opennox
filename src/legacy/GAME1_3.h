#ifndef NOX_PORT_GAME1_3
#define NOX_PORT_GAME1_3

#include "defs.h"

void nox_game_decStateInd_43BDC0();
void nox_game_addStateCode_43BDD0(int a1);
int nox_game_getStateCode_43BE10();
int sub_43BE30();
void sub_43BE40(int a1);
int sub_43BE50_get_video_mode_id();
wchar2_t* get_video_mode_string(int v1);
void nox_xxx_gameGetScreenBoundaries_43BEB0_get_video_mode(int* w, int* h, int* d);
int nox_game_switchStates_43C0A0();
nox_gui_animation* nox_gui_makeAnimation_43C5B0(nox_window* win, int x1, int y1, int x2, int y2, int in_dx, int in_dy,
												int out_dx, int out_dy);
int sub_43C6E0();
void sub_43CF40();
void sub_43CF70();
int sub_43DB20();
int sub_43DB30(int a1);
char* sub_43DB40(int a1);
void sub_43DBE0();
void nox_game_SetCliDrawFunc(void* a1);
int sub_43E940(void* a1);
void sub_43E9F0();
int sub_43EA20(void* a1);
int sub_43EC10();
int sub_43EC30(void* a1);
int sub_43ECB0(void* a1);
int sub_43ED00(uint32_t* a1);
int sub_43EFD0(void* a1);
int sub_43F010(void* a1);
int sub_43F030(int a1);
int sub_43F050();
int sub_43F060(uint32_t* a1);
int sub_43F0D0();
int sub_43F130();
int nox_xxx_guiFontHeightMB_43F320(void* a1);
void* nox_xxx_guiFontPtrByName_43F360(char* a1);
int nox_xxx_drawString_43F6E0(void* a1, wchar2_t* a2, int a3, int a4);
int nox_draw_drawStringHL_43F730(void* a1, wchar2_t* a2, int a3, int a4);
int nox_xxx_drawStringStyle_43F7B0(void* a1, wchar2_t* a2, int a3, int a4);
int nox_xxx_drawGetStringSize_43F840(void* a1, wchar2_t* a2, int* a3, int* a4, int a5);
int nox_xxx_bookGetStringSize_43FA80(void* a1, wchar2_t* a2, int* a3, int* a4, int a5);
int nox_xxx_drawStringWrap_43FAF0(void* a1, wchar2_t* str, int a3, int a4, int a5, int a6);
int nox_xxx_drawStringWrapHL_43FD00(void* a1, wchar2_t* a2, int a3, int a4, int a5, int a6);
int nox_xxx_bookDrawString_43FA80_43FD80(void* a1, wchar2_t* a2, int a3, int a4, int a5, int a6);
void nox_client_clearScreen_440900();
void nox_client_quit_4460C0();
int nox_xxx_serverIsClosing_446180();
void sub_446380();
int sub_449E60(char a1);
void sub_44A360(int a1);
void sub_44A400();
int sub_44A4A0();
void sub_44A4B0();
int sub_44A4E0();
void nox_xxx____setargv_4_44B000();

void* nox_xxx_dialogMsgBoxCreate_449A10(nox_window* win, wchar2_t* a2, wchar2_t* text, int a4, void* a5, void* a6);

#endif // NOX_PORT_GAME1_3
