#!/usr/bin/env python3
"""Capture actual scalar-SSE2 minimp3 inverse transforms and overlap sequences."""
import argparse,gzip,hashlib,json,struct,subprocess
from pathlib import Path
from mp3_source import original_header

SHIM=r'''
#include <stdio.h>
#include <stdlib.h>
#define MINIMP3_ONLY_MP3
#define MINIMP3_NO_SIMD
#define MINIMP3_IMPLEMENTATION
#include "minimp3.h"
_Static_assert(sizeof(int)==4&&sizeof(float)==4,"requires32bit int/float");
static void rd(void*p,size_t n){if(fread(p,1,n,stdin)!=n)exit(2);}
static uint32_t ru(void){unsigned char b[4];rd(b,4);return(uint32_t)b[0]|(uint32_t)b[1]<<8|(uint32_t)b[2]<<16|(uint32_t)b[3]<<24;}
static void wr(const void*p,size_t n){if(fwrite(p,1,n,stdout)!=n)exit(3);}
static void wu(uint32_t v){unsigned char b[4]={v,v>>8,v>>16,v>>24};wr(b,4);}
static void floats(const float*p,int n){for(int i=0;i<n;i++){uint32_t b;memcpy(&b,p+i,4);wu(b);}}
static void pattern(float*p,int n,unsigned id,uint32_t seed,unsigned mark){
 uint32_t state=seed?seed:1;
 for(int i=0;i<n;i++){
  uint32_t b=0;
  if(id==1)b=0x80000000;
  if(id==2){float v=(float)(i+1+(seed%17))*(1.f/16);if(i&1)v=-v;memcpy(&b,&v,4);}
  if(id==3)b=i==(int)mark?0x3f800000:0;
  if(id==4){state^=state<<13;state^=state>>17;state^=state<<5;b=(state&0x807fffff)|((100+(state>>24)%40)<<23);}
  if(id==5){static const uint32_t small[4]={1,0x007fffff,0x00800000,0x80000001};b=small[(i+seed)%4];}
  memcpy(p+i,&b,4);
 }
}
static int intact(const uint8_t*p,int n){for(int i=0;i<n;i++)if(p[i]!=0xcc)return 0;return 1;}
int main(void){
 int op;while((op=getchar())!=EOF){
  uint32_t pat=ru(),seed=ru(),mark=ru();if(pat>5||mark>639)return 4;
  struct{uint8_t pre[8];float v[640];uint8_t post[8];} data;
  struct{uint8_t pre[8];float v[320];uint8_t post[8];} overlap;
  memset(&data,0xcc,sizeof(data));memset(&overlap,0xcc,sizeof(overlap));
  pattern(data.v,640,pat,seed,mark);pattern(overlap.v,320,pat<2?pat:4,seed^0x50434d,mark);
  if(op==1){L3_dct3_9(data.v);floats(data.v,18);wu(intact(data.pre,8)&&intact(data.post,8));}
  else if(op==2){
   float before[3];memcpy(before,data.v,sizeof(before));
   L3_idct3(data.v[0],data.v[1],data.v[2],overlap.v);floats(overlap.v,9);wu(memcmp(before,data.v,sizeof(before))==0);wu(intact(overlap.pre,8)&&intact(overlap.post,8));
  }else if(op==3){
   uint32_t off=ru();if(off>2)return 5;float before[20];memcpy(before,data.v,sizeof(before));
   struct{uint8_t pre[8];float v[12];uint8_t post[8];} dst;memset(&dst,0xcc,sizeof(dst));pattern(dst.v,12,4,seed^0x445354,0);
   L3_imdct12(data.v+off,dst.v,overlap.v);floats(dst.v,12);floats(overlap.v,9);wu(memcmp(before,data.v,sizeof(before))==0);wu(intact(dst.pre,8)&&intact(dst.post,8)&&intact(overlap.pre,8)&&intact(overlap.post,8));
  }else if(op==4||op==5){
   uint32_t bands=ru();if(bands>32)return 6;
   if(op==4){uint32_t win=ru();if(win>2)return 7;float window[18];for(int i=0;i<18;i++)window[i]=win==0?1:win==1?(i&1):(float)(i+1)*(1.f/32);L3_imdct36(data.v,overlap.v,window,bands);}
   else L3_imdct_short(data.v,overlap.v,bands);
   floats(data.v,640);floats(overlap.v,320);wu(intact(data.pre,8)&&intact(data.post,8)&&intact(overlap.pre,8)&&intact(overlap.post,8));
  }else if(op==6){L3_change_sign(data.v);floats(data.v,640);wu(intact(data.pre,8)&&intact(data.post,8));}
  else if(op==7){
   uint32_t steps=ru();if(!steps||steps>6)return 8;
   for(unsigned step=0;step<steps;step++){
    unsigned block=ru(),longs=ru();if(block>3||longs>32)return 9;
    pattern(data.v,640,pat,seed+step*37,(mark+step*17)%576);
    L3_imdct_gr(data.v,overlap.v,block,longs);
    floats(data.v,640);floats(overlap.v,320);wu(intact(data.pre,8)&&intact(data.post,8)&&intact(overlap.pre,8)&&intact(overlap.post,8));
   }
  }else return 10;
 }
 return ferror(stdin)||fflush(stdout);
}
'''
def u(*v):return struct.pack('<'+'I'*len(v),*(x&0xffffffff for x in v))
def requests():
    for pat in range(6):
        for seed in range(1 if pat<2 else 8):
            marks=range(9) if pat==3 else (seed,)
            for mark in marks:yield 1,u(pat,seed,mark),76
            marks=range(3) if pat==3 else (seed,)
            for mark in marks:yield 2,u(pat,seed,mark),44
            for off in range(3):
                marks=range(20) if pat==3 else (seed,)
                for mark in marks:yield 3,u(pat,seed,mark,off),92
            for bands in (0,1,2,4,31,32):
                for win in range(3):yield 4,u(pat,seed,seed,bands,win),3844
                yield 5,u(pat,seed,seed,bands),3844
            yield 6,u(pat,seed,seed),2564
            for block in range(4):
                for longs in (0,2,4,32):yield 7,u(pat,seed,seed,1,block,longs),3844
            for longs in (0,2,4):
                seq=[(0,0),(1,0),(2,longs),(2,longs),(3,0),(0,0)]
                yield 7,u(pat,seed,seed,len(seq))+b''.join(u(*row) for row in seq),3844*len(seq)
            seq=[(2,32),(0,32),(3,32)]
            yield 7,u(pat,seed,seed,len(seq))+b''.join(u(*row) for row in seq),3844*len(seq)
    # Basis impulses at all18coordinates of first and last subbands.
    for bands,base in ((1,0),(32,558)):
        for mark in range(base,base+18):
            for win in range(3):yield 4,u(3,0,mark,bands,win),3844
            yield 5,u(3,0,mark,bands),3844

def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def main():
    ap=argparse.ArgumentParser(description=__doc__);ap.add_argument('output',type=Path);args=ap.parse_args()
    root=Path(__file__).resolve().parents[2];out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
    header=original_header(root,out);shim=out/'capture.c';shim.write_text(SHIM)
    flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse','-I',str(header.parent),str(shim),'-lm']
    subprocess.run(flags+['-o',str(out/'capture')],check=True);counts=[0]*7;wrapper_steps=0
    with (out/'requests.bin').open('wb') as f:
        for op,payload,n in requests():
            f.write(bytes((op,))+payload);counts[op-1]+=1
            if op==7:wrapper_steps+=n//3844
    hashes=[]
    for i in range(3):
        with (out/'requests.bin').open('rb') as fi,(out/f'results-{i}.bin').open('wb') as fo:subprocess.run([str(out/'capture')],stdin=fi,stdout=fo,check=True)
        hashes.append(sha(out/f'results-{i}.bin'))
    assert len(set(hashes))==1
    subprocess.run(flags+['-fsanitize=undefined','-fno-sanitize-recover=all','-o',str(out/'capture-ubsan')],check=True)
    with (out/'requests.bin').open('rb') as fi,(out/'results-ubsan.bin').open('wb') as fo:subprocess.run([str(out/'capture-ubsan')],stdin=fi,stdout=fo,check=True)
    assert sha(out/'results-ubsan.bin')==hashes[0]
    h=hashlib.sha256()
    with (out/'results-0.bin').open('rb') as result,(out/'imdct_c.bin.gz').open('wb') as dest:
        with gzip.GzipFile(filename='',mode='wb',fileobj=dest,mtime=0) as gz:
            def write(b):h.update(b);gz.write(b)
            write(b'NMP3IMD1')
            for op,payload,n in requests():
                expected=result.read(n);assert len(expected)==n;write(bytes((op,))+payload+expected)
            assert not result.read(1)
    meta=dict(counts=counts,wrapper_steps=wrapper_steps,result_sha256=hashes,ubsan_result_sha256=sha(out/'results-ubsan.bin'),uncompressed_sha256=h.hexdigest(),fixture_sha256=sha(out/'imdct_c.bin.gz'),source_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),header_sha256=sha(header),tool_sha256=sha(__file__),shim_sha256=sha(shim),compiler=subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],flags=flags)
    (out/'capture.json').write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta,indent=2))
if __name__=='__main__':main()
