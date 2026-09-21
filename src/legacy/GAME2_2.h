#ifndef NOX_PORT_GAME2_2
#define NOX_PORT_GAME2_2

#include "defs.h"

void nox_video_setCutSize_4766A0(int a1);
int nox_video_getCutSize_4766D0();
void nox_draw_setCutSize_476700(int cutPerc, int a2);
int nox_client_setCursorType_477610(int a1);
int nox_xxx_cursorGetTypePrev_477630();
void nox_xxx_bookSaveSpellForDragDrop_477640(int a1, int a2);
void nox_xxx_bookSpellDnDclear_477660();
int nox_xxx_bookGetSpellDnDType_477670();
void nox_xxx_cursorSetDraggedItem_477690(nox_drawable* a1);
void nox_xxx_cursorResetDraggedItem_4776A0();
void nox_xxx_cursorSetTooltip_4776B0(wchar2_t* a1);
int sub_478030();
int sub_478040();
void sub_478850(int a1, short a2, int a3, int a4);
int sub_478E50(int a1, int a2, unsigned int a3);
int sub_479590();
void sub_4795A0(int a1);
int sub_479690(int a1, short a2, short a3, int a4);
void sub_479810();
int sub_479820(int a1, short a2);
int sub_479D00();
void nox_client_drawImageAt_47D2C0(nox_video_bag_image_t* img, int x, int y);
void sub_47D370(int a1);
void sub_47D400(int a1, char a2);
int nox_draw_imageMeta_47D5C0(nox_video_bag_image_t* a1, uint32_t* a2, uint32_t* a3, uint32_t* a4, uint32_t* a5);
unsigned char sub_47DBC0();
int sub_47FCE0(uint32_t* a1, int a2);
uint8_t* sub_480220(uint8_t* a1, uint8_t* a2);
uint16_t* sub_480250(uint8_t* a1, uint16_t* a2);
long long sub_484C00(int a1, int a2);
long long nox_xxx_spriteChangeLightSize_484C30(int a1, int a2);
int sub_484C60(float a1);
int sub_484CE0(int a1, float a2);
int sub_4862E0(void* a3, int a4);
void* sub_486320(void* a1, int a2);
int sub_486350(void* a1, int a2);
int sub_486380(void* a1, uint32_t a2, int32_t a3, uint32_t a4);
int sub_4863B0(void* a2);
void* sub_4864A0(void* a3);
int sub_486520(void* a2);
int sub_486550(void* a1);
void sub_486570(void* a1, void* a2);
void sub_486620(void* a1);
int sub_4873C0(int a3);
int* sub_487810(int a1, int a2);
int nox_xxx_wndEditProc_487D70(nox_window* a1, int a2, int a3, int a4);
int nox_xxx_wndEditDrawNoImage_488160(int a1, int a2);
int nox_xxx_wndStaticDrawNoImage_488D00(nox_window* a1p, nox_window_data* a2p);
nox_window* nox_gui_newStaticText_489300(nox_window* a1, int a2, int a3, int a4, int a5, int a6, nox_window_data* a7p, nox_staticText_data* a8p);
int nox_xxx_setSomeFunc_48A210(int a1);

#endif // NOX_PORT_GAME2_2
