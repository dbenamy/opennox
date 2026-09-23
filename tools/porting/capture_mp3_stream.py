#!/usr/bin/env python3
"""Capture actual minimp3 frame scanning, initialization and byte-reservoir state."""
from pathlib import Path
from mp3_source import original_header
import argparse
import gzip
import hashlib
import json
import random
import struct
import subprocess

SHIM=r'''
#include <stdio.h>
#include <stdlib.h>
#define MINIMP3_ONLY_MP3
#define MINIMP3_NO_SIMD
#define MINIMP3_IMPLEMENTATION
#include "minimp3.h"
_Static_assert(sizeof(int)==4,"requires int32");
static void rd(void*p,size_t n){if(fread(p,1,n,stdin)!=n)exit(2);}
static uint32_t ru(void){unsigned char b[4];rd(b,4);return(uint32_t)b[0]|(uint32_t)b[1]<<8|(uint32_t)b[2]<<16|(uint32_t)b[3]<<24;}
static void wr(const void*p,size_t n){if(fwrite(p,1,n,stdout)!=n)exit(3);}
static void wu(uint32_t v){unsigned char b[4]={v,v>>8,v>>16,v>>24};wr(b,4);}
static void pattern(unsigned char*p,size_t n,uint32_t seed){for(size_t i=0;i<n;i++)p[i]=(i*73+seed)^(i>>2);}
int main(void){
 int op;unsigned char data[49152+64];mp3dec_t dec;mp3dec_scratch_t scratch;
 while((op=getchar())!=EOF){
  if(op==1 || op==2){
   uint32_t n=ru(),hint=ru();if(n>49152||hint>2304)return 4;
   memset(data,0,sizeof(data));rd(data,n);
   int free_bytes=hint,frame_bytes=0x12345678;
   if(op==1){int offset=mp3d_find_frame(data,n,&free_bytes,&frame_bytes);wu(offset);wu(free_bytes);wu(frame_bytes);}
   else {if(n<4||!hdr_valid(data))return 5;wu(mp3d_match_frame(data,n,hint));}
  }else if(op==3){
   uint32_t seed=ru();memset(&dec,seed,sizeof(dec));
   mp3dec_t before;memcpy(&before,&dec,sizeof(dec));mp3dec_init(&dec);
   wu(dec.header[0]);before.header[0]=0;wu(memcmp(&before,&dec,sizeof(dec))==0);
  }else if(op==4){
   uint32_t have=ru(),begin=ru(),length=ru(),pos=ru(),seed=ru();
   if(have>511||begin>511||length>2304||pos>length*8)return 6;
   memset(&dec,seed,sizeof(dec));dec.reserv=have;pattern(dec.reserv_buf,511,seed);
   memset(&scratch,seed^255,sizeof(scratch));pattern(data,sizeof(data),seed+19);
   bs_t outer;bs_init(&outer,data,length);outer.pos=pos;
   mp3dec_t before;memcpy(&before,&dec,sizeof(dec));bs_t outer_before=outer;
   int ret=L3_restore_reservoir(&dec,&outer,&scratch,begin);
   wu(ret);wu(scratch.bs.pos);wu(scratch.bs.limit);wu(scratch.bs.buf==scratch.maindata);
   wu(memcmp(&before,&dec,sizeof(dec))==0);wu(outer.buf==outer_before.buf&&outer.pos==outer_before.pos&&outer.limit==outer_before.limit);
   wr(scratch.maindata,sizeof(scratch.maindata));
  }else if(op==5){
   uint32_t length=ru(),pos=ru(),seed=ru();if(length>2815||pos>length*8+32)return 7;
   memset(&dec,seed,sizeof(dec));pattern(dec.reserv_buf,511,seed);
   memset(&scratch,0,sizeof(scratch));pattern(scratch.maindata,sizeof(scratch.maindata),seed+19);
   bs_init(&scratch.bs,scratch.maindata,length);scratch.bs.pos=pos;
   mp3dec_scratch_t before;memcpy(&before,&scratch,sizeof(scratch));
   L3_save_reservoir(&dec,&scratch);wu(dec.reserv);wr(dec.reserv_buf,511);wu(memcmp(&before,&scratch,sizeof(scratch))==0);
  }else return 8;
 }
 return ferror(stdin)||fflush(stdout);
}
'''

def u(*v):return struct.pack('<'+'I'*len(v),*v)

def requests():
    # Byte packets exercise both normal rates and header-only Layer I/II handling.
    headers=[(bytes.fromhex(h),size) for h,size in [('fffb9000',417),('fff380c0',208),('ffe38000',417),('fffd9000',522),('ffff9000',312)]]
    for n in range(25):yield 1,u(n,0)+bytes(n),12
    for h,size in headers:
        for count in (1,2,3,10,11):
            packet=bytearray(size*count+4)
            for i in range(count):packet[i*size:i*size+4]=h
            for prefix in (0,1,3,17):
                full=bytes(prefix)+packet
                for trim in (-1,0,1,3,4):
                    n=prefix+size*count+trim
                    if 0<=n<=len(full):yield 1,u(n,0)+full[:n],12
            for n in (4,size-1,size,size+1,len(packet)):
                if n<=len(packet):yield 2,u(n,0)+packet[:n],4
        for padding in (0,1):
            padded=bytearray(h);padded[2]|=padding*2
            padbytes=(4 if h[1]&6==6 else 1)*padding
            sizepad=size+padbytes
            packet=bytearray(sizepad*3)
            for i in range(3):packet[i*sizepad:i*sizepad+4]=padded
            yield 1,u(len(packet),0)+packet,12
    # Free-format discovery requires a consistent third header; changing padding
    # shifts successive boundaries. Probe either side of the2304 search bound.
    for version in (0xe3,0xf3,0xfb,0xfd,0xff):
        for size in (4,5,16,417,2303,2304):
            for padded in (False,True):
                pad=4 if version&6==6 else 1
                h0=bytes((255,version,2 if padded else 0,0));h1=bytes((255,version,0,0))
                first=size+(pad if padded else 0)
                packet=bytearray(first+2*size+8)
                for at,h in ((0,h0),(first,h1),(first+size,h0)):
                    packet[at:at+4]=h
                for n in sorted({4,first,first+size,first+size+4,len(packet)}):
                    for hint in (0,size):yield 1,u(n,hint)+packet[:n],12
                # Defined direct matcher with zero hint intentionally checks the
                # original bounded no-progress loop for unpadded free format.
                yield 2,u(len(packet),0)+packet,4
    rng=random.Random(0x535452)
    for _ in range(512):
        n=rng.randrange(1025);packet=bytearray(rng.randbytes(n))
        if n>=4:
            at=rng.randrange(n-3);packet[at:at+4]=rng.choice(headers)[0]
        yield 1,u(n,rng.choice((0,4,417,2304)))+packet,12
    for seed in (0,1,0x5a,0xa5,255):yield 3,u(seed),8
    for have in (0,1,255,510,511):
        for begin in (0,1,255,510,511):
            for length in (0,1,17,511,2304):
                for pos in sorted({0,min(1,length*8),min(7,length*8),length*8}):
                    yield 4,u(have,begin,length,pos,0xa5),24+2815
    for length in (0,1,2,510,511,512,513,2304,2815):
        for pos in sorted(set(range(9))|{max(0,length*8-i) for i in range(9)}|{length*8+1,length*8+7,length*8+32}):
            if pos<=length*8+32:
                for seed in (0,0xa5,255):yield 5,u(length,pos,seed),4+511+4


    # Actual ail decoder buffer is 16*1024*3 bytes. Keep the preceding baseline
    # records byte-identical and append coverage beyond the initial probe bound.
    for n in (32769,49151,49152):
        yield 1,u(n,0)+bytes(n),12
    for h,size in headers:
        for prefix in (0,49152-3*size):
            packet=bytearray(49152)
            for i in range(3):packet[prefix+i*size:prefix+i*size+4]=h
            yield 1,u(len(packet),0)+packet,12
            if prefix==0:yield 2,u(len(packet),0)+packet,4

def main():
 ap=argparse.ArgumentParser(description=__doc__)
 ap.add_argument('output',type=Path,help='new capture directory')
 args=ap.parse_args()
 root=Path(__file__).resolve().parents[2]
 out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
 header=original_header(root,out)
 shim=out/'capture.c';shim.write_text(SHIM)
 flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse','-I',str(header.parent),str(shim),'-lm']
 subprocess.run(flags+['-o',str(out/'capture')],check=True)
 counts=[0]*5
 with (out/'requests.bin').open('wb') as f:
  for op,payload,n in requests():f.write(bytes((op,))+payload);counts[op-1]+=1
 hashes=[]
 for i in range(3):
  with (out/'requests.bin').open('rb') as fi,(out/f'results-{i}.bin').open('wb') as fo:subprocess.run([str(out/'capture')],stdin=fi,stdout=fo,check=True)
  hashes.append(hashlib.sha256((out/f'results-{i}.bin').read_bytes()).hexdigest())
 assert len(set(hashes))==1
 subprocess.run(flags+['-fsanitize=undefined','-fno-sanitize-recover=all','-o',str(out/'capture-ubsan')],check=True)
 with (out/'requests.bin').open('rb') as fi,(out/'results-ubsan.bin').open('wb') as fo:subprocess.run([str(out/'capture-ubsan')],stdin=fi,stdout=fo,check=True)
 assert hashlib.sha256((out/'results-ubsan.bin').read_bytes()).hexdigest()==hashes[0]
 h=hashlib.sha256()
 with (out/'results-0.bin').open('rb') as result,(out/'stream_c.bin.gz').open('wb') as dest:
  with gzip.GzipFile(filename='',mode='wb',fileobj=dest,mtime=0) as gz:
   def write(b):h.update(b);gz.write(b)
   write(b'NMP3STR1')
   for op,payload,n in requests():
    expected=result.read(n);assert len(expected)==n;write(bytes((op,))+payload+expected)
   assert not result.read(1)

 def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
 meta={'counts':counts,'result_sha256':hashes,'ubsan_result_sha256':sha(out/'results-ubsan.bin'),
       'uncompressed_sha256':h.hexdigest(),'fixture_sha256':sha(out/'stream_c.bin.gz'),
       'source_commit':subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),
       'header_sha256':sha(header),'tool_sha256':sha(__file__),'shim_sha256':sha(shim),
       'compiler':subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],
       'flags':flags,'ubsan_flags':['-fsanitize=undefined','-fno-sanitize-recover=all']}
 (out/'capture.json').write_text(json.dumps(meta,indent=2)+'\n')
 print(json.dumps(meta,indent=2))

if __name__=='__main__':main()
