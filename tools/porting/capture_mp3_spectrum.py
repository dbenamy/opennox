#!/usr/bin/env python3
"""Capture actual scalar minimp3 stereo, reorder and antialias state."""
import argparse,gzip,hashlib,json,struct,subprocess
from pathlib import Path
from mp3_source import original_header
from capture_mp3_sideinfo import payload as side_payload

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
static void pattern(float*p,int n,unsigned id,unsigned mark){
 for(int i=0;i<n;i++){
  uint32_t b=0;
  if(id==1||id==4||id==5)b=0x3f000000+(i%257)*0x800;
  if(id==2)b=(0x3f800000+(i%257)*0x1000)|((i&1)?0x80000000:0);
  if(id==3){static const uint32_t small[4]={1,0x007fffff,0x00800000,0x80000001};b=small[i%4];}
  if((id==4||id==5)&&i>=576&&i<1152)b=(id==5&&i==576+(int)mark)?0x3f800000:0;
  memcpy(p+i,&b,4);
 }
}
static int intact(const uint8_t*p,int n){for(int i=0;i<n;i++)if(p[i]!=0xcc)return 0;return 1;}
static void setup(uint8_t h[4],L3_gr_info_t gr[4]){
 uint8_t side[64];rd(h,4);rd(side,64);if(!hdr_valid(h)||HDR_GET_LAYER(h)!=1)exit(4);
 memset(gr,0,sizeof(L3_gr_info_t)*4);bs_t bs;bs_init(&bs,side,64);
 if(L3_read_side_info(&bs,gr,h)<0)exit(5);
}
static uint8_t ipos(unsigned id,int i){switch(id){case 0:return 0;case 1:return i%7;case 2:return 6;case 3:return 7;case 4:return 63;case 5:return 64;case 6:return 255;default:return i*7+3;}}
int main(void){
 int op;while((op=getchar())!=EOF){
  struct{uint8_t pre[8];float v[1192];uint8_t post[8];} buf;
  struct{uint8_t pre[8];float v[576];uint8_t post[8];} scratch;
  struct{uint8_t pre[8],v[39],post[8];} ist;
  memset(&buf,0xcc,sizeof(buf));memset(&scratch,0xcc,sizeof(scratch));memset(&ist,0xcc,sizeof(ist));
  for(int i=0;i<576;i++){uint32_t b=0x4b123456+i;memcpy(scratch.v+i,&b,4);}
  if(op==1||op==2){
   unsigned n=ru(),off=ru(),pat=ru();if(n>576||off>576-n||pat>5)return 6;
   pattern(buf.v,1192,pat,0);
   if(op==1)L3_midside_stereo(buf.v+off,n);
   else{uint32_t bl=ru(),br=ru();float kl,kr;memcpy(&kl,&bl,4);memcpy(&kr,&br,4);L3_intensity_stereo_band(buf.v+off,n,kl,kr);}
   floats(buf.v,1192);wu(intact(buf.pre,8)&&intact(buf.post,8));
  }else if(op==3||op==4||op==5){
   uint8_t h[4];L3_gr_info_t gr[4];setup(h,gr);unsigned pat=ru(),mark=ru();if(pat>5||mark>575)return 7;
   pattern(buf.v,1192,pat,mark);
   if(op==3){
    float before[1192];memcpy(before,buf.v,sizeof(before));int max[3]={99,99,99};
    L3_stereo_top_band(buf.v+576,gr[0].sfbtab,gr[0].n_long_sfb+gr[0].n_short_sfb,max);
    for(int i=0;i<3;i++)wu(max[i]);wu(memcmp(before,buf.v,sizeof(before))==0);wu(intact(buf.pre,8)&&intact(buf.post,8));
   }else if(op==4){
    unsigned shift=ru(),ip=ru();if(shift>1||ip>7)return 8;gr[1].scalefac_compress=shift;
    L3_gr_info_t before[4];memcpy(before,gr,sizeof(gr));for(int i=0;i<39;i++)ist.v[i]=ipos(ip,i);
    L3_intensity_stereo(buf.v,ist.v,gr,h);
    floats(buf.v,1192);wr(ist.v,39);wu(intact(buf.pre,8)&&intact(buf.post,8)&&intact(ist.pre,8)&&intact(ist.post,8));wu(memcmp(before,gr,sizeof(gr))==0);
   }else{
    unsigned ch=ru(),off=ru();if(ch>1||off>72)return 9;
    // The buffer deliberately includes the40float scalefactor work area following
    // the two576float channels. This keeps the legacy8kHz mixed extent defined.
    L3_reorder(buf.v+ch*576+off,scratch.v,gr[0].sfbtab+gr[0].n_long_sfb);
    floats(buf.v,1192);floats(scratch.v,576);wu(intact(buf.pre,8)&&intact(buf.post,8)&&intact(scratch.pre,8)&&intact(scratch.post,8));
   }
  }else if(op==6){
   int32_t bands=(int32_t)ru();unsigned ch=ru(),pat=ru();if(bands< -1||bands>31||ch>1||pat>5)return 10;
   pattern(buf.v,1192,pat,0);L3_antialias(buf.v+576*ch,bands);floats(buf.v,1192);wu(intact(buf.pre,8)&&intact(buf.post,8));
  }else return 11;
 }
 return ferror(stdin)||fflush(stdout);
}
'''
def u(*v):return struct.pack('<'+'I'*len(v),*(x&0xffffffff for x in v))
def setup(version,rate,layout,mode=0x50):
    h=bytes((255,version,0x80|rate*4,mode))
    change=None if layout==0 else (0,'block',2) if layout==2 else (0,'mixed',1)
    side,_=side_payload(h,0,change)
    return h+side

def requests():
    for pattern in range(4):
        for n in (0,1,7,8,17,18,575,576):
            for off in sorted({0,576-n}):yield 1,u(n,off,pattern),4772
    gains=[(0,0x3f800000),(0x3f800000,0),(0x3f000000,0x3f000000),(0x3fb504f3,0x3fb504f3),(0xbf800000,0x3f800000),(0x80000000,0x00800000)]
    for pattern in (1,2,3):
        for n in (0,1,8,576):
            for off in sorted({0,576-n}):
                for kl,kr in gains:yield 2,u(n,off,pattern,kl,kr),4772
    for version in (0xe3,0xf3,0xfb):
        for rate in range(3):
            for layout in range(3):
                s=setup(version,rate,layout)
                for pattern,mark in [(0,0),(2,0),(4,0)]+[(5,m) for m in (0,1,17,35,47,48,71,72,497,498,575)]:
                    yield 3,s+u(pattern,mark),20
                for mode in (0x50,0x70):
                    s=setup(version,rate,layout,mode)
                    for shift in ((0,) if version==0xfb else (0,1)):
                        for pattern,mark in ((2,0),(3,0),(4,0),(5,0),(5,575)):
                            for ip in range(8):yield 4,s+u(pattern,mark,shift,ip),4815
                if layout!=0:
                    off=0 if layout==2 else (72 if version==0xe3 and rate==2 else 36)
                    for ch in (0,1):
                        for pattern in range(4):yield 5,setup(version,rate,layout)+u(pattern,0,ch,off),7076
    for pattern in range(4):
        for bands in (-1,0,1,2,15,30,31):
            for ch in (0,1):yield 6,u(bands,ch,pattern),4772

def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def main():
    ap=argparse.ArgumentParser(description=__doc__);ap.add_argument('output',type=Path);args=ap.parse_args()
    root=Path(__file__).resolve().parents[2];out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
    header=original_header(root,out);shim=out/'capture.c';shim.write_text(SHIM)
    flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse','-I',str(header.parent),str(shim),'-lm']
    subprocess.run(flags+['-o',str(out/'capture')],check=True);counts=[0]*6
    with (out/'requests.bin').open('wb') as f:
        for op,payload,n in requests():f.write(bytes((op,))+payload);counts[op-1]+=1
    hashes=[]
    for i in range(3):
        with (out/'requests.bin').open('rb') as fi,(out/f'results-{i}.bin').open('wb') as fo:subprocess.run([str(out/'capture')],stdin=fi,stdout=fo,check=True)
        hashes.append(sha(out/f'results-{i}.bin'))
    assert len(set(hashes))==1
    subprocess.run(flags+['-fsanitize=undefined','-fno-sanitize-recover=all','-o',str(out/'capture-ubsan')],check=True)
    with (out/'requests.bin').open('rb') as fi,(out/'results-ubsan.bin').open('wb') as fo:subprocess.run([str(out/'capture-ubsan')],stdin=fi,stdout=fo,check=True)
    assert sha(out/'results-ubsan.bin')==hashes[0]
    h=hashlib.sha256()
    with (out/'results-0.bin').open('rb') as result,(out/'spectrum_c.bin.gz').open('wb') as dest:
        with gzip.GzipFile(filename='',mode='wb',fileobj=dest,mtime=0) as gz:
            def write(b):h.update(b);gz.write(b)
            write(b'NMP3SPC1')
            for op,payload,n in requests():
                expected=result.read(n);assert len(expected)==n;write(bytes((op,))+payload+expected)
            assert not result.read(1)
    meta=dict(counts=counts,result_sha256=hashes,ubsan_result_sha256=sha(out/'results-ubsan.bin'),uncompressed_sha256=h.hexdigest(),fixture_sha256=sha(out/'spectrum_c.bin.gz'),source_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),header_sha256=sha(header),tool_sha256=sha(__file__),side_input_tool_sha256=sha(root/'tools/porting/capture_mp3_sideinfo.py'),shim_sha256=sha(shim),compiler=subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],flags=flags)
    (out/'capture.json').write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta,indent=2))
if __name__=='__main__':main()
