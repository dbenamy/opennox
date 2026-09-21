#ifndef NOX_PORT_GAME2
#define NOX_PORT_GAME2

#include "defs.h"

char* nox_get_thing_name(int i);
wchar2_t* nox_get_thing_pretty_name(int i);
wchar2_t* nox_get_thing_desc(int i);
int nox_get_thing_pretty_image(int i);
int nox_drawable_link_thing(nox_drawable* a1, int i);
nox_thing* sub_44D330(char* a1);
void sub_44D3A0();
void sub_44D5C0(int a1, int a2);
void sub_44D640();
int sub_44D660(char* a1);
int sub_44D7E0(int a1);
void sub_44D8F0();
int nox_xxx_playDialogFile_44D900(unsigned char* a1, int a2);
int sub_44D930();
void nox_video_stopAllFades_44E040();
int nox_video_inFadeTransition_44E0D0();
void sub_44E320();
int nox_client_getIntroScreenDuration_44E3B0();
void sub_450580();
unsigned char sub_450750();
char sub_450760(char a1);
int sub_4526D0(int a1);
int sub_4526F0(int a1);
int* sub_455800();
int sub_455A00(int a1);
int sub_455F10(int a1);
int* sub_456D60(int a1);
int sub_457460(int a1);
int nox_xxx_guiServerOptionsTryHide_4574D0();
char* sub_459AA0(void* a1);
int sub_459D50(int a1);
int sub_459D80(int a1);
nox_drawable* nox_xxx_cliFirstMinimapObj_459EB0();
nox_drawable* nox_xxx_cliGetSpritePlayer_45A000();
nox_drawable* sub_45A010(nox_drawable* dr);
nox_drawable* sub_45A060();
nox_drawable* nox_drawable_next_45A070(nox_drawable* a1);
void nox_client_toggleSpellbook_45AC70();
int nox_xxx_bookHideMB_45ACA0(int a1);
int nox_xxx_bookClickSpell_45B1F0();
int nox_xxx_bookClickCreature_45B200();
int sub_45CFC0();
int* nox_xxx_bookSetForward_45D200(int* a1, int a2, int2* a3);
void nox_xxx_abilityReward_45D290(int a1, char* a2, int a3);
int sub_45D500(int a1);
void nox_xxx_bookFillAll_45D570(int a1, int a2);
int sub_45D9B0();
void* nox_xxx_book_45DBE0(void* a1, int a2, int a3);
int nox_xxx_clientUpdateButtonRow_45E110(int a1);
int nox_xxx_buttonsGetSelectedRow_45E180();
void* sub_4602F0();
int nox_client_trapSetSelect_4604B0(int a1);
int sub_4604E0();
int sub_460660();
int nox_xxx_quickBarClose_4606B0();
int sub_460940(void* this);

void nox_video_fadeInScreen_44DAB0(int a1, int a2, void* a3);
void nox_video_fadeOutScreen_44DB30(int a1, int a2, void* a3);

#endif // NOX_PORT_GAME2
