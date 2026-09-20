#ifndef NOX_PORT_GAME2_3
#define NOX_PORT_GAME2_3

#include "defs.h"

void sub_48C580(pixel8888* a1, int num);
unsigned int sub_48C6B0(int a1, int a2);
int sub_48CAD0();
int sub_48D4A0();
int sub_48D4B0(int a1);
int sub_48D4F0(unsigned short a1, unsigned short a2);
int sub_48D560(unsigned short a1);
uint32_t* sub_48D5A0(int a1);
int sub_48D660();
int sub_48D740();
void sub_48D760();
int* sub_48D7B0();
void sub_48E8E0(int a1);
void sub_48E940();
nox_drawable* nox_xxx_spriteCreate_48E970(int a1, unsigned short a2, int a3, int a4);
char* sub_4947E0(int a1);
int sub_4948B0(int a1);
int nox_xxx_netCliProcUpdateStream_494A60(unsigned char* a1, int a2, uint32_t* a3);
unsigned char* nox_xxx_netCliUpdateStream2_494C30(unsigned char* a1, int a2, int* a3);
int sub_495060(int a1, short a2, short a3);
int sub_4950C0(int a1);
int sub_4950F0(int a1, char a2);
int sub_495120(int a1, short a2, short a3);
int sub_495150(int a1, short a2);
int nox_xxx_unitSpriteCheckAlly_4951F0(int a1);
int sub_495210(int a1);
void sub_4959B0();
uint32_t* nox_xxx_cliAddObjFriend_4959F0(int a1);
void sub_495A20(int a1);
nox_point sub_499290(int a1);
int sub_4992B0(int a1, int a2);
void nox_xxx_fxDrawTurnUndead_499880(short* a1);
void nox_xxx_drawPointMB_499B70(int xLeft, int yTop, int a3);
void* nox_npc_by_id(int id);
void nox_xxx_clientEquip_49A3D0(char a1, int a2, int a3, int a4);
uint16_t* nox_xxx_cliAddHealthChange_49A650(int a1, short a2);
void nox_xxx_sprite_49AA00_drawable(nox_drawable* dr);
void nox_xxx_updateSpritePosition_49AA90(nox_drawable* dr, int a2, int a3);
void nox_xxx_forEachSprite_49AB00(int4*, void*, int);
nox_drawable* nox_drawable_find_49ABF0(nox_point* pt, int r);
int sub_49AEA0();
void nox_xxx_spriteTransparentDecay_49B950(nox_drawable* a1, int a2);
void nox_xxx_spriteToSightDestroyList_49BAB0_drawable(nox_drawable* a1);
void nox_xxx_spriteToList_49BC80_drawable(nox_drawable* a1);
void sub_49BCD0(nox_drawable* dr);
void nox_xxx_clientAddRayEffect_49C160(int a1);
void nox_xxx_clientRemoveRayEffect_49C450(int a1);
void nox_client_drawBorderLines_49CC70(int xLeft, int yTop, int a3, int a4);
void sub_49CD30(int xLeft, int yTop, int a3, int a4, int a5, int a6);
void nox_client_drawRectFilledOpaque_49CE30(int xLeft, int yTop, int a3, int a4);
void nox_client_drawRectFilledAlpha_49CF10(int xLeft, int yTop, int a3, int a4);
int sub_49D1C0(void* a1, int a2, int a3);
int nox_client_drawLineFromPoints_49E4B0();
int sub_49E4F0(int a1);
void nox_client_drawPixel_49EFA0(int a1, int a2);
void nox_client_drawAddPoint_49F500(int a1, int a2);
void nox_xxx_rasterPointRel_49F570(int a1, int a2);
int nox_client_copyRect_49F6F0(int xLeft, int yTop, int a3, int a4);
void sub_49F7C0_def();
void nox_xxx_wndDraw_49F7F0();
int sub_49F860();
nox_window* nox_new_window_from_file(char* cname, void* fnc);


#endif // NOX_PORT_GAME2_3
