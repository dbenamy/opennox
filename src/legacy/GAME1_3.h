#ifndef NOX_PORT_GAME1_3
#define NOX_PORT_GAME1_3

#include "defs.h"

void nox_game_addStateCode_43BDD0(int a1);
int nox_game_getStateCode_43BE10();
int nox_game_switchStates_43C0A0();
nox_gui_animation* nox_gui_makeAnimation_43C5B0(nox_window* win, int x1, int y1, int x2, int y2, int in_dx, int in_dy,
												int out_dx, int out_dy);
void nox_game_SetCliDrawFunc(void* a1);
void* nox_xxx_guiFontPtrByName_43F360(char* a1);

void* nox_xxx_dialogMsgBoxCreate_449A10(nox_window* win, wchar2_t* a2, wchar2_t* text, int a4, void* a5, void* a6);

#endif // NOX_PORT_GAME1_3
