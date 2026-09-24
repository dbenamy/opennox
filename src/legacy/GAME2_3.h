#ifndef NOX_PORT_GAME2_3
#define NOX_PORT_GAME2_3

#include "defs.h"

unsigned int sub_48C6B0(int a1, int a2);
void nox_xxx_drawPointMB_499B70(int xLeft, int yTop, int a3);
void nox_client_drawBorderLines_49CC70(int xLeft, int yTop, int a3, int a4);
void sub_49CD30(int xLeft, int yTop, int a3, int a4, int a5, int a6);
void nox_client_drawRectFilledOpaque_49CE30(int xLeft, int yTop, int a3, int a4);
void nox_client_drawRectFilledAlpha_49CF10(int xLeft, int yTop, int a3, int a4);
int sub_49D1C0(void* a1, int a2, int a3);
int nox_client_drawLineFromPoints_49E4B0();
void nox_client_drawPixel_49EFA0(int a1, int a2);
void nox_client_drawAddPoint_49F500(int a1, int a2);
void nox_xxx_rasterPointRel_49F570(int a1, int a2);
int nox_client_copyRect_49F6F0(int xLeft, int yTop, int a3, int a4);
void nox_xxx_wndDraw_49F7F0();
int sub_49F860();
nox_window* nox_new_window_from_file(char* cname, void* fnc);


#endif // NOX_PORT_GAME2_3
