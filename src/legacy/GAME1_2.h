#ifndef NOX_PORT_GAME1_2
#define NOX_PORT_GAME1_2

#include "defs.h"

unsigned char nox_xxx_wall_42A6C0(unsigned char a1, unsigned char a2);
int nox_xxx_objectTOCgetTT_42C2B0(unsigned short a1);
nox_video_bag_image_t* nox_xxx_readImgMB_42FAA0(int known_idx, char a2, char* a3);
void* nox_video_getImagePixdata_42FB30(nox_video_bag_image_t* img);
nox_things_imageRef_t* nox_xxx_gLoadAnim_42FA20(char* a1);
nox_video_bag_image_t* nox_xxx_gLoadImg_42F970(char* name);
int sub_430B50(int a1, int a2, int a3, int a4);
void sub_430C30_set_video_max(int w, int h);
void nox_xxx_screenGetSize_430C50_get_video_max(int* w, int* h);
void nox_client_resetScreenParticles_431510();
nox_screenParticle* nox_client_newScreenParticle_431540(int a1, int a2, int a3, int a4, int a5, int a6, char a7,
														char a8, char a9, char a10);
void sub_434080(int a1);
void nox_xxx_drawSetTextColor_434390(int a1);
void nox_xxx_drawSetColor_4343E0(int a1);
void nox_client_drawSetColor_434460(int a1);
void nox_client_drawEnableAlpha_434560(int a1);
void nox_client_drawSetAlpha_434580(unsigned char a1);
void sub_4345F0(int a1);
void nox_xxx_draw_434600(int a1);
void sub_434990(int a1, int a2, int a3);
void nox_draw_splitColor_435280(short a1, uint8_t* a2, uint8_t* a3, uint8_t* a4);
void nox_xxx_getSomeCoods_435670(int2* a1);
void nox_xxx_spriteLoadError_4356E0();
void nox_client_setServerConnectAddr_435720(char* addr);
int sub_436550();
nox_draw_viewport_t* nox_draw_getViewport_437250();
void sub_437260();
void sub_437290();
int sub_438330();
int sub_438370();
int nox_client_joinGame_438A90();
int sub_43AF30();
int nox_xxx_cliDrawConnectedLoop_43B360();
int sub_43B490();
void nox_xxx_serverHost_43B4D0();

#endif // NOX_PORT_GAME1_2
