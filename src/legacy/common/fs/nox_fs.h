#ifndef NOX_COMMON_FS_H
#define NOX_COMMON_FS_H

#include <stdbool.h>
#include <stdio.h>

#define NOX_FILEPATH_MAX 1024


// nox_fs_set_workdir sets current work directory.
bool nox_fs_set_workdir(char* path);



// nox_fs_open opens the file for reading (in binary mode).
FILE* nox_fs_open(char* path);
// nox_fs_open_text opens the file for reading (in text mode).
FILE* nox_fs_open_text(char* path);
// nox_fs_create creates the file for writing (in binary mode).
FILE* nox_fs_create(char* path);
// nox_fs_create_text creates or opens the file for writing (in text mode).
FILE* nox_fs_create_text(char* path);
// nox_fs_open_rw opens the file for reading and writing (in binary mode).
FILE* nox_fs_open_rw(char* path);


int nox_fs_fseek(FILE* f, long off, int mode);


#define nox_fs_fseek_start(f, off) nox_fs_fseek(f, off, SEEK_SET)
#define nox_fs_fseek_cur(f, off) nox_fs_fseek(f, off, SEEK_CUR)

#endif // NOX_COMMON_FS_H
