#!/usr/bin/env python3
"""Freeze integer helper results from the production minimp3 header (no copied algorithms)."""
import argparse
import gzip
import hashlib
import json
from pathlib import Path
from mp3_source import original_header
import struct
import subprocess

SHIM = r'''
#include <stdio.h>
#include <stdlib.h>
#define MINIMP3_ONLY_MP3
#define MINIMP3_NO_SIMD
#define MINIMP3_IMPLEMENTATION
#include "minimp3.h"
_Static_assert(sizeof(int) == 4, "requires int32");
static void rd(void *p, size_t n) { if (fread(p,1,n,stdin)!=n) exit(2); }
static uint32_t ru(void) { unsigned char b[4]; rd(b,4); return (uint32_t)b[0]|(uint32_t)b[1]<<8|(uint32_t)b[2]<<16|(uint32_t)b[3]<<24; }
static void wu(uint32_t v) { unsigned char b[4]={v,v>>8,v>>16,v>>24}; if(fwrite(b,1,4,stdout)!=4) exit(3); }
int main(void) {
 int op; uint8_t h[4],h2[4],buf[2815+32];
 while ((op=getchar())!=EOF) {
  if(op==1) { rd(h,4); wu(hdr_valid(h)); }
  else if(op==2) { rd(h,4);rd(h2,4);wu(hdr_compare(h,h2)); }
  else if(op==3) { rd(h,4); int32_t free_bytes=(int32_t)ru();
   if(!hdr_valid(h)) return 4;
   wu(hdr_bitrate_kbps(h));wu(hdr_sample_rate_hz(h));wu(hdr_frame_samples(h));wu(hdr_frame_bytes(h,free_bytes));wu(hdr_padding(h));
  } else if(op==4) {
   uint32_t pattern=ru(),len=ru(),start=ru(),width=ru();
   if(pattern>4 || len>2815 || start>len*8+32 || width>32) return 5;
   for(unsigned i=0;i<sizeof(buf);i++) buf[i]=pattern==0?0:pattern==1?255:pattern==2?(i&1?170:85):pattern==3?i:((i*73+19)^(i>>2));
   bs_t bs; bs.pos=17;bs.limit=8;bs.buf=buf;bs_init(&bs,buf,len);
   int initial=bs.pos;bs.pos=start;uint32_t value=get_bits(&bs,width);
   wu(value);wu(bs.pos);wu(bs.limit);wu(initial);
  } else return 6;
 }
 return ferror(stdin) || fflush(stdout);
}
'''

def u(*values):
    return struct.pack('<'+'I'*len(values), *(v & 0xffffffff for v in values))

def requests():
    for prefix in (0xff, 0xfe, 0):
        for b1 in range(256):
            for b2 in range(256):
                yield 1, bytes((prefix,b1,b2,0)), 4
    for prefix in range(256):
        yield 1, bytes((prefix,0xfb,0x90,0xff)), 4
    for base in (bytes.fromhex('fffb9000'),bytes.fromhex('fff300c0'),bytes.fromhex('ffe280ff'),bytes.fromhex('00ff9000')):
        for b1 in range(256):
            for b2 in range(256):
                other=bytes((255,b1,b2,255))
                yield 2, base+other, 4
                yield 2, other+base, 4
    # Defined table domains, restricted to structurally valid sync/version/layer bits.
    for b1 in (0xe2,0xe3,*range(0xf2,0xf8),*range(0xfa,0x100)):
        for bitrate in range(15):
            for rate in range(3):
                for low in range(4):
                    for last in (0,0xff):
                        for free in (-2147483648,-1,0,1,2304,2147483647):
                            yield 3, bytes((255,b1,bitrate*16+rate*4+low,last))+u(free),20
    for pattern in range(5):
        for length in (0,1,2,3,4,7,31,64,511,2304,2815):
            end=length*8
            starts=set(range(8))|{max(0,end-i) for i in range(40)}|{end,end+1,end+7,end+32}
            for start in sorted(starts):
                if start>end+32:
                    continue
                for width in range(33):
                    yield 4,u(pattern,length,start,width),16

def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

def main():
    ap=argparse.ArgumentParser(description=__doc__)
    ap.add_argument('output',type=Path,help='new capture directory (must not exist)')
    args=ap.parse_args()
    root=Path(__file__).resolve().parents[2]
    out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
    header=original_header(root,out)
    shim=out/'capture.c';shim.write_text(SHIM)
    exe=out/'capture'
    flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse','-I',str(header.parent),str(shim),'-lm','-o',str(exe)]
    subprocess.run(flags,check=True)
    req=out/'requests.bin'
    counts=[0]*4
    with req.open('wb') as f:
        for op,payload,_ in requests():
            f.write(bytes((op,))+payload);counts[op-1]+=1
    hashes=[]
    for trial in range(3):
        result=out/f'results-{trial}.bin'
        with req.open('rb') as fin,result.open('wb') as fout:
            subprocess.run([str(exe)],stdin=fin,stdout=fout,check=True)
        hashes.append(digest(result))
    assert len(set(hashes))==1, hashes
    fixture=out/'integer_c.bin.gz'; sha=hashlib.sha256();size=0
    with (out/'results-0.bin').open('rb') as results, fixture.open('wb') as dest:
        with gzip.GzipFile(filename='',fileobj=dest,mode='wb',mtime=0) as gz:
            def write(b):
                nonlocal size
                sha.update(b);size+=len(b);gz.write(b)
            write(b'NMP3INT1')
            for op,payload,n in requests():
                expected=results.read(n);assert len(expected)==n
                write(bytes((op,))+payload+expected)
            assert results.read(1)==b''
    metadata=dict(header_sha256=digest(header),tool_sha256=digest(Path(__file__)),shim_sha256=digest(shim),
        compiler=subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],flags=flags,
        opcode_counts=counts,uncompressed_bytes=size,uncompressed_sha256=sha.hexdigest(),
        fixture_sha256=digest(fixture),repeated_result_sha256=hashes,
        source_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip())
    (out/'capture.json').write_text(json.dumps(metadata,indent=2)+'\n')
    print(json.dumps(metadata,indent=2))

if __name__=='__main__':
    main()
