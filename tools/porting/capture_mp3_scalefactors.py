#!/usr/bin/env python3
"""Capture actual scalar-SSE2 minimp3 scalefactor bytes and float bit patterns."""
import argparse,gzip,hashlib,json,struct,subprocess
from pathlib import Path

SHIM=r'''
#include <stdio.h>
#include <stdlib.h>
#define MINIMP3_ONLY_MP3
#define MINIMP3_NO_SIMD
#define MINIMP3_IMPLEMENTATION
#include "minimp3.h"
_Static_assert(sizeof(int)==4 && sizeof(float)==4,"requires32-bit int/float");
static void rd(void*p,size_t n){if(fread(p,1,n,stdin)!=n)exit(2);}
static uint32_t ru(void){unsigned char b[4];rd(b,4);return(uint32_t)b[0]|(uint32_t)b[1]<<8|(uint32_t)b[2]<<16|(uint32_t)b[3]<<24;}
static void wr(const void*p,size_t n){if(fwrite(p,1,n,stdout)!=n)exit(3);}
static void wu(uint32_t v){unsigned char b[4]={v,v>>8,v>>16,v>>24};wr(b,4);}
static void pattern(uint8_t*p,size_t n,unsigned id){for(size_t i=0;i<n;i++)p[i]=id==0?0:id==1?255:id==3?i:((i*73+19)^(i>>2));}
static int intact(const uint8_t*p,size_t n){for(size_t i=0;i<n;i++)if(p[i]!=0xcc)return 0;return 1;}
int main(void){
 int op;uint8_t data[128];
 while((op=getchar())!=EOF){
  if(op==1){
   uint8_t sizes[4],counts[4];rd(sizes,4);rd(counts,4);int32_t scfsi=(int32_t)ru();uint32_t len=ru(),start=ru(),pat=ru();
   unsigned sum=0;for(int i=0;i<4;i++){if(sizes[i]>5)return 4;sum+=counts[i];}
   if(sum>36||len>64||start>32||(scfsi<0&&scfsi!=-16)||scfsi>15||pat>4)return 5;
   struct {uint8_t pre[8],v[40],post[8];} scf,ist;
   memset(&scf,0xcc,sizeof(scf));memset(&ist,0xcc,sizeof(ist));
   for(int i=0;i<40;i++){scf.v[i]=0xa5^(i*13);ist.v[i]=0x80+i*7;}
   pattern(data,sizeof(data),pat);bs_t bs;bs_init(&bs,data,len);bs.pos=start;
   L3_read_scalefactors(scf.v,ist.v,sizes,counts,&bs,scfsi);
   wu(bs.pos);wu(bs.limit);wr(scf.v,40);wr(ist.v,40);
   wu(intact(scf.pre,8)&&intact(scf.post,8)&&intact(ist.pre,8)&&intact(ist.post,8));
  }else if(op==2){
   uint32_t bits=ru(),exp=ru();if((bits&0x7f800000)==0x7f800000||exp>1024)return 6;
   float y;memcpy(&y,&bits,4);y=L3_ldexp_q2(y,exp);memcpy(&bits,&y,4);wu(bits);
  }else if(op==3){
   uint8_t hdr[4];rd(hdr,4);
   uint32_t len=ru(),start=ru(),pat=ru(),compress=ru(),gain=ru(),layout=ru(),sub=ru(),flags=ru(),ch=ru();
   if(!hdr_valid(hdr)||HDR_GET_LAYER(hdr)!=1||len>64||start>32||pat>4||compress>(HDR_TEST_MPEG1(hdr)?15:511)||gain>255||layout>2||sub>511||flags>0xf03||ch>1)return 7;
   L3_gr_info_t gr;memset(&gr,0,sizeof(gr));gr.scalefac_compress=compress;gr.global_gain=gain;
   gr.n_long_sfb=layout==0?22:layout==1?(HDR_TEST_MPEG1(hdr)?8:6):0;gr.n_short_sfb=layout==0?0:layout==1?30:39;
   gr.subblock_gain[0]=sub&7;gr.subblock_gain[1]=(sub>>3)&7;gr.subblock_gain[2]=(sub>>6)&7;
   gr.preflag=flags&1;gr.scalefac_scale=(flags>>1)&1;gr.scfsi=(flags>>8)&15;
   L3_gr_info_t before=gr;
   struct {uint8_t pre[8];float v[40];uint8_t post[8];} scf;
   struct {uint8_t pre[8],v[40],post[8];} ist;
   memset(&scf,0xcc,sizeof(scf));memset(&ist,0xcc,sizeof(ist));
   for(int i=0;i<40;i++){uint32_t bits=0x4b123456+i;memcpy(&scf.v[i],&bits,4);ist.v[i]=0x80+i*7;}
   pattern(data,sizeof(data),pat);bs_t bs;bs_init(&bs,data,len);bs.pos=start;
   L3_decode_scalefactors(hdr,ist.v,&bs,&gr,scf.v,ch);
   wu(bs.pos);wu(bs.limit);for(int i=0;i<40;i++){uint32_t bits;memcpy(&bits,&scf.v[i],4);wu(bits);}wr(ist.v,40);
   wu(intact(scf.pre,8)&&intact(scf.post,8)&&intact(ist.pre,8)&&intact(ist.post,8));wu(memcmp(&gr,&before,sizeof(gr))==0);
  }else return 8;
 }
 return ferror(stdin)||fflush(stdout);
}
'''
def u(*v):return struct.pack('<'+'I'*len(v),*(x&0xffffffff for x in v))
def requests():
    counts=[(0,0,0,0),(6,5,5,5),(8,9,6,12),(9,9,6,12),(6,6,6,3),(12,12,12,0),(1,0,4,5)]
    sizes=[(x,)*4 for x in range(6)]+[(0,1,2,5),(5,0,3,1)]
    for count in counts:
        for size in sizes:
            for scfsi in (-16,0,1,2,4,8,15):
                for length in (0,1,8,64):
                    for start in (0,3,16):
                        for pat in (0,1,4):yield 1,bytes(size+count)+u(scfsi,length,start,pat),92
    for bits in (0,0x80000000,1,0x007fffff,0x00800000,0x3f800000,0x40000000,0x45000000,0x7f7fffff,0xbf800000):
        for exp in range(1025):yield 2,u(bits,exp),4
    # Contexts cover normal/intensity/MS/mono modes, retained bytes, truncation,
    # byte wrap, both scalefactor scales and all compressor values.
    contexts=[(0,0,0,0,0,0,0,64,0,0),(0x70,1,1,1,15,511,1,64,7,255),
              (0x60,0,0,1,5,1+3*8+7*64,4,0,0,210),(0xc0,0,1,0,10,511,3,17,16,213)]
    for version in (0xe3,0xf3,0xfb):
        for compress in range(16 if version==0xfb else 512):
            for layout in range(3):
                for mode,ch,pre,scale,scfsi,sub,pat,length,start,gain in contexts:
                    h=bytes((255,version,0x90 if version==0xfb else 0x80,mode))
                    if version!=0xfb:pre=int(compress>=500)
                    if layout!=0:scfsi=0
                    flags=pre|(scale<<1)|(scfsi<<8)
                    yield 3,h+u(length,start,pat,compress,gain,layout,sub,flags,ch),216

def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def main():
    ap=argparse.ArgumentParser(description=__doc__);ap.add_argument('output',type=Path);args=ap.parse_args()
    root=Path(__file__).resolve().parents[2];out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
    header=root/'src/legacy/client/audio/mp3/minimp3.h';shim=out/'capture.c';shim.write_text(SHIM)
    flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse','-I',str(header.parent),str(shim),'-lm']
    subprocess.run(flags+['-o',str(out/'capture')],check=True);counts=[0]*3
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
    with (out/'results-0.bin').open('rb') as result,(out/'scalefactors_c.bin.gz').open('wb') as dest:
        with gzip.GzipFile(filename='',mode='wb',fileobj=dest,mtime=0) as gz:
            def write(b):h.update(b);gz.write(b)
            write(b'NMP3SCF1')
            for op,payload,n in requests():
                expected=result.read(n);assert len(expected)==n;write(bytes((op,))+payload+expected)
            assert not result.read(1)
    meta=dict(counts=counts,result_sha256=hashes,ubsan_result_sha256=sha(out/'results-ubsan.bin'),uncompressed_sha256=h.hexdigest(),fixture_sha256=sha(out/'scalefactors_c.bin.gz'),source_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),header_sha256=sha(header),tool_sha256=sha(__file__),shim_sha256=sha(shim),compiler=subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],flags=flags)
    (out/'capture.json').write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta,indent=2))
if __name__=='__main__':main()
