#!/usr/bin/env python3
"""Capture real minimp3 side-information parsing, including partial output writes."""
import argparse
import gzip
import hashlib
import json
from pathlib import Path
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
static void wu(uint32_t v){unsigned char b[4]={v,v>>8,v>>16,v>>24};if(fwrite(b,1,4,stdout)!=4)exit(3);}
int main(void){
 uint8_t hdr[4],data[4096+64],sentinel[40];int first;
 while((first=getchar())!=EOF){
  hdr[0]=first;rd(hdr+1,3);uint32_t len=ru(),start=ru(),seed=ru();
  if(!hdr_valid(hdr)||HDR_GET_LAYER(hdr)!=1||len>4096||start>len*8+32||start>512||seed>255)return 4;
  memset(data,0,sizeof(data));rd(data,64);memset(sentinel,seed,40);
  L3_gr_info_t gr[4];memset(gr,seed,sizeof(gr));for(int i=0;i<4;i++)gr[i].sfbtab=sentinel;
  bs_t bs;bs_init(&bs,data,len);bs.pos=start;
  int ret=L3_read_side_info(&bs,gr,hdr);wu(ret);wu(bs.pos);wu(bs.limit);
  for(int i=0;i<4;i++){
   L3_gr_info_t*g=&gr[i];
   wu(g->part_23_length);wu(g->big_values);wu(g->scalefac_compress);
   wu(g->global_gain);wu(g->block_type);wu(g->mixed_block_flag);wu(g->n_long_sfb);wu(g->n_short_sfb);
   for(int j=0;j<3;j++)wu(g->table_select[j]);
   for(int j=0;j<3;j++)wu(g->region_count[j]);
   for(int j=0;j<3;j++)wu(g->subblock_gain[j]);
   wu(g->preflag);wu(g->scalefac_scale);wu(g->count1_table);wu(g->scfsi);
   wu(g->sfbtab==sentinel);
   uint8_t table[40]={0};
   if(g->sfbtab==sentinel)memcpy(table,sentinel,40);
   else {int j=0;do {table[j]=g->sfbtab[j];}while(table[j]&&++j<40);}
   if(fwrite(table,1,40,stdout)!=40)return 5;
  }
 }
 return ferror(stdin)||fflush(stdout);
}
'''

def pack(values):
    bits=[]
    for v,n in values:
        assert 0<=v<1<<n
        bits.extend((v>>i)&1 for i in reversed(range(n)))
    assert len(bits)<=512
    out=bytearray(64)
    for i,b in enumerate(bits):out[i//8]|=b<<(7-i%8)
    return bytes(out),len(bits)

def payload(header,start,change=None):
    mpeg1=bool(header[1]&8);mono=header[3]>>6==3;n=(1 if mono else 2)*(2 if mpeg1 else 1)
    values=[(0,start)] if start else []
    main=change[2] if change and change[:2]==(-1,'main') else 0
    values.extend([(main,9),(0x7ff&((1<<(7+n))-1),7+n)] if mpeg1 else [(main,8),(3&((1<<n)-1),n)])
    for i in range(n):
        d=dict(part=0,big=0,gain=210,compress=0,switch=0,block=2,mixed=0,tables=0,sub0=0,sub1=0,sub2=0,region0=0,region1=0,pre=0,scale=0,count=0)
        if change and change[0]==i:d[change[1]]=change[2]
        # block/mixed/subblock mutations require a switched block.
        if change and change[0]==i and change[1] in ('block','mixed','sub0','sub1','sub2'):d['switch']=1
        values.extend([(d['part'],12),(d['big'],9),(d['gain'],8),(d['compress'],4 if mpeg1 else 9),(d['switch'],1)])
        if d['switch']:
            values.extend([(d['block'],2),(d['mixed'],1),(d['tables']&1023,10),(d['sub0'],3),(d['sub1'],3),(d['sub2'],3)])
        else:values.extend([(d['tables'],15),(d['region0'],4),(d['region1'],3)])
        if mpeg1:values.append((d['pre'],1))
        values.extend([(d['scale'],1),(d['count'],1)])
    return pack(values)

def requests():
    rng=random.Random(0x4e4f58)
    for version in (0xe3,0xf3,0xfb):
        for rate in range(3):
            for mode in (0,0x40,0x80,0xc0):
                h=bytes((255,version,0x80|rate*4,mode))
                n=(1 if mode==0xc0 else 2)*(2 if version==0xfb else 1)
                changes=[None,(-1,'main',1),(-1,'main',255)]
                if version==0xfb:changes.append((-1,'main',511))
                domains=dict(part=(1,7,8,4095),big=(1,287,288,289,511),gain=(0,255),compress=((0,15) if version==0xfb else (0,1,499,500,511)),switch=(1,),block=(0,1,2,3),mixed=(1,),tables=(32767,),sub0=(7,),sub1=(7,),sub2=(7,),region0=(15,),region1=(7,),pre=(1,),scale=(1,),count=(1,))
                for i in range(n):
                    for field,values in domains.items():
                        if field=='pre' and version!=0xfb:continue
                        changes.extend((i,field,v) for v in values)
                for start in (0,3,7,32):
                    for change in changes:
                        data,bits=payload(h,start,change)
                        # Exact side-info bound vs ample payload; equal and one-bit-short budget.
                        for length in sorted({(bits+7)//8,512}):
                            yield h,length,start,0xa5,data
                    data,bits=payload(h,start)
                    for length in range(41):
                        if start<=length*8+32:
                            yield h,length,start,0,data
                            yield h,length,start,255,data
                # Arbitrary fields exercise combinations and error prefixes.
                for i in range(128):
                    data=rng.randbytes(64);length=rng.choice((0,1,2,8,16,32,64,512,4096));start=rng.randrange(min(length*8+32,64)+1)
                    yield h,length,start,rng.choice((0,0x5a,0xa5,255)),data

def wire(row):
    h,length,start,seed,data=row
    return h+struct.pack('<III',length,start,seed)+data

def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()

def main():
    ap=argparse.ArgumentParser(description=__doc__);ap.add_argument('output',type=Path);ap.add_argument('--ubsan',action='store_true');args=ap.parse_args()
    root=Path(__file__).resolve().parents[2];out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
    header=root/'src/legacy/client/audio/mp3/minimp3.h';shim=out/'capture.c';shim.write_text(SHIM);exe=out/'capture'
    flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse']
    if args.ubsan:flags+=['-fsanitize=undefined','-fno-sanitize-recover=all']
    flags+=['-I',str(header.parent),str(shim),'-lm','-o',str(exe)];subprocess.run(flags,check=True)
    req=out/'requests.bin';count=0
    with req.open('wb') as f:
        for row in requests():f.write(wire(row));count+=1
    hashes=[]
    for trial in range(3):
        result=out/f'results-{trial}.bin'
        with req.open('rb') as fin,result.open('wb') as fout:subprocess.run([str(exe)],stdin=fin,stdout=fout,check=True)
        hashes.append(sha(result));assert result.stat().st_size==count*524
    assert len(set(hashes))==1
    digest=hashlib.sha256();fixture=out/'sideinfo_c.bin.gz';success=errors=0
    with (out/'results-0.bin').open('rb') as result,fixture.open('wb') as dest:
        with gzip.GzipFile(filename='',mode='wb',fileobj=dest,mtime=0) as gz:
            def write(b):digest.update(b);gz.write(b)
            write(b'NMP3SID1')
            for row in requests():
                expected=result.read(524);assert len(expected)==524
                if struct.unpack('<i',expected[:4])[0]<0:errors+=1
                else:success+=1
                write(wire(row)+expected)
            assert not result.read(1)
    meta=dict(source_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),header_sha256=sha(header),tool_sha256=sha(__file__),shim_sha256=sha(shim),compiler=subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],flags=flags,cases=count,success=success,errors=errors,uncompressed_sha256=digest.hexdigest(),fixture_sha256=sha(fixture),result_sha256=hashes)
    (out/'capture.json').write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta,indent=2))

if __name__=='__main__':main()
