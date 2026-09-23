#!/usr/bin/env python3
"""Capture complete original scalar MP3 frame calls, state and PCM capacity writes."""
import argparse,gzip,hashlib,json,struct,subprocess,os
from pathlib import Path
import capture_mp3_imdct as common
import capture_mp3_sideinfo as sideinfo
SHIM=common.SHIM.split('int main(void)')[0]+r'''
static void state(const mp3dec_t*d){floats(d->mdct_overlap[0],288);floats(d->mdct_overlap[1],288);for(int i=0;i<960;i++){if(i>=898&&i<956&&(i%4==2||i%4==3))wu(0);else floats(d->qmf_state+i,1);}wu(d->reserv);wu(d->free_format_bytes);wr(d->header,4);wr(d->reserv_buf,511);}
static void shorts(const int16_t*p,int n){for(int i=0;i<n;i++){uint16_t v=p[i];uint8_t b[2]={v,v>>8};wr(b,2);}}
int main(void){int op;while((op=getchar())!=EOF){if(op!=1)return 4;unsigned steps=ru();if(!steps||steps>8)return 5;mp3dec_t dec;memset(&dec,0,sizeof(dec));
 for(unsigned step=0;step<steps;step++){
  unsigned flags=ru(),n=ru();if(flags>3||n>49152)return 6;
  uint8_t packet[49152],before[49152];rd(packet,n);memcpy(before,packet,n);
  struct{uint8_t pre[8];int16_t v[2368];uint8_t post[8];} pcm;memset(&pcm,0xcc,sizeof(pcm));
  mp3dec_frame_info_t info;memset(&info,0x5a,sizeof(info));if(flags&1)mp3dec_init(&dec);
  if(getenv("OPENNOX_QMF_POISON")){for(int i=0;i<15;i++){dec.qmf_state[896+4*i+2]=12345.f+i;dec.qmf_state[896+4*i+3]=-12345.f-i;}}
  int samples=mp3dec_decode_frame(&dec,packet,n,(flags&2)?NULL:pcm.v,&info);
  wu(samples);wu(info.frame_bytes);wu(info.channels);wu(info.hz);wu(info.layer);wu(info.bitrate_kbps);state(&dec);shorts(pcm.v,2368);wu(memcmp(packet,before,n)==0);wu(intact(pcm.pre,8)&&intact(pcm.post,8));
 }}return ferror(stdin)||fflush(stdout);}
'''
def setbits(data,start,n,v):
 for i in range(n):
  at=start+i;mask=1<<(7-at%8)
  data[at//8]=(data[at//8]&~mask)|(((v>>(n-1-i))&1)*mask)
def frame(h,change=None,free=None):
 mpeg1=bool(h[1]&8);rate=[44100,48000,32000][(h[2]>>2)&3]
 if not mpeg1:rate>>=1
 if not h[1]&16:rate>>=1
 bitrate=112 if mpeg1 else 64
 size=(144000*bitrate//(rate*(1 if mpeg1 else 2)))+((h[2]>>1)&1) if free is None else free+((h[2]>>1)&1)
 packet=bytearray(size);packet[:4]=h
 change0=(0,'big',1) if change in ('coded','escape') else change
 side,_=sideinfo.payload(h,0,change0)
 side=bytearray(side)
 if mpeg1:
  n=(1 if h[3]>>6==3 else 2)*2
  # Keep actual SCFSI reuse bits, clear private bits. Original C leaks nonzero
  # private bits into first-granule reuse of uninitialized ist_pos storage.
  setbits(side,9,7+n,(1<<(2*n))-1)
 if change in ('coded','escape'):
  n=(1 if h[3]>>6==3 else 2)*(2 if mpeg1 else 1)
  start=16+n if mpeg1 else 8+n
  setbits(side,start,12,32 if change=='coded' else 128)
  setbits(side,start+12+9+8+(4 if mpeg1 else 9)+1,15,1057 if change=='coded' else 32767)
 off=4+(0 if h[1]&1 else 2);assert size>=off+64;packet[off:off+64]=side
 return bytes(packet)
def requests():
 for version in (0xe3,0xf3,0xfb):
  for rate in range(3):
   for mode in (0,0x50,0x70,0xc0):
    for crc in (0,1):
     for pad in (0,1):
      h=bytes((255,version-crc,0x80|rate*4|pad*2,mode));f=frame(h)
      yield [(0,f*3),(0,f*2),(0,f),(1,f),(2,f)]
    h=bytes((255,version,0x80|rate*4,mode))
    for change in ((0,'block',1),(0,'block',2),(0,'block',3),(0,'mixed',1),'coded','escape'):
     f=frame(h,change);plain=frame(h)
     yield [(0,f*3),(0,f*2),(0,plain),(1,f),(2,f)]
    for change in ((0,'block',0),(0,'big',289),(-1,'main',1),(-1,'main',511 if version==0xfb else 255)):
     f=frame(h,change);plain=frame(h)
     yield [(0,f*3),(0,f*2),(0,plain),(0,f),(1,f)]
 # Prefix/truncation scans, unsupported Layer I/II, reset and metadata-only calls.
 f=frame(bytes.fromhex('fffb8000'),'coded')
 for n in (0,1,3,4,5,31,63,len(f)-1,len(f),len(f)+1,len(f)*2):
  for prefix in (0,1,17):
   packet=(bytes(prefix)+f*3)[:n]
   yield [(0,packet),(2,packet),(0,f*2),(1,packet)]
 for h,size in ((bytes.fromhex('fffd9000'),522),(bytes.fromhex('ffff9000'),312)):
  packet=h+bytes(size-4)
  yield [(2,packet*3),(0,packet*2),(1,packet),(2,packet)]
 for version in (0xe3,0xf3,0xfb):
  h=bytes((255,version,0,0));f=frame(h,free=417)
  yield [(0,f*3),(0,f*2),(0,f),(1,f*3),(2,f)]
 # Changing mode/rate forces validation/reset; long input covers the real48KiB caller.
 a=frame(bytes.fromhex('fffb8000'),'coded');b=frame(bytes.fromhex('fff384c0'),'escape')
 yield [(0,a*3),(0,b*3),(0,a*2),(1,b*3)]
 packet=(a*((49152+len(a)-1)//len(a)))[:49152]
 yield [(0,packet),(2,packet),(1,packet)]
def wire(seq):return b'\1'+common.u(len(seq))+b''.join(common.u(flags,len(p))+p for flags,p in seq)
def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def main():
 ap=argparse.ArgumentParser(description=__doc__);ap.add_argument('output',type=Path);args=ap.parse_args()
 root=Path(__file__).resolve().parents[2];out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
 header=root/'src/legacy/client/audio/mp3/minimp3.h';shim=out/'capture.c';shim.write_text(SHIM)
 flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse','-I',str(header.parent),str(shim),'-lm']
 subprocess.run(flags+['-o',str(out/'capture')],check=True);rows=list(requests());steps=sum(map(len,rows));(out/'requests.bin').write_bytes(b''.join(map(wire,rows)))
 hashes=[]
 for i in range(3):
  with (out/'requests.bin').open('rb') as fi,(out/f'results-{i}.bin').open('wb') as fo:subprocess.run([str(out/'capture')],stdin=fi,stdout=fo,check=True)
  hashes.append(sha(out/f'results-{i}.bin'))
 assert len(set(hashes))==1,hashes
 subprocess.run(flags+['-fsanitize=undefined,float-cast-overflow','-fno-sanitize-recover=all','-o',str(out/'capture-ubsan')],check=True)
 with (out/'requests.bin').open('rb') as fi,(out/'results-ubsan.bin').open('wb') as fo:subprocess.run([str(out/'capture-ubsan')],stdin=fi,stdout=fo,check=True)
 assert sha(out/'results-ubsan.bin')==hashes[0]
 with (out/'requests.bin').open('rb') as fi,(out/'results-poison.bin').open('wb') as fo:subprocess.run([str(out/'capture')],stdin=fi,stdout=fo,check=True,env=dict(os.environ,OPENNOX_QMF_POISON='1'))
 assert sha(out/'results-poison.bin')==hashes[0]
 h=hashlib.sha256()
 with (out/'results-0.bin').open('rb') as fi,(out/'frame_c.bin.gz').open('wb') as fo:
  with gzip.GzipFile(filename='',mode='wb',fileobj=fo,mtime=0) as gz:
   def write(b):h.update(b);gz.write(b)
   write(b'NMP3FRM1')
   for seq in rows:
    expected=fi.read(11435*len(seq));assert len(expected)==11435*len(seq)
    for j in range(len(seq)):assert expected[(j+1)*11435-8:(j+1)*11435]==common.u(1,1)
    write(wire(seq)+expected)
   assert not fi.read(1)
 meta=dict(normalized_qmf_indices=[896+4*i+j for i in range(15) for j in (2,3)],poison_result_sha256=sha(out/'results-poison.bin'),sequences=len(rows),frame_calls=steps,result_sha256=hashes,ubsan_result_sha256=sha(out/'results-ubsan.bin'),uncompressed_sha256=h.hexdigest(),fixture_sha256=sha(out/'frame_c.bin.gz'),source_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),header_sha256=sha(header),tool_sha256=sha(__file__),pattern_tool_sha256=sha(common.__file__),side_payload_tool_sha256=sha(sideinfo.__file__),shim_sha256=sha(shim),compiler=subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],flags=flags)
 (out/'capture.json').write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta,indent=2))
if __name__=='__main__':main()
