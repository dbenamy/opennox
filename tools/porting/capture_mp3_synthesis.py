#!/usr/bin/env python3
"""Capture scalar PCM rounding, filterbank and persistent synthesis state."""
import argparse,gzip,hashlib,json,struct,subprocess
from pathlib import Path
from mp3_source import original_header
import capture_mp3_imdct as common
SHIM=common.SHIM.split('int main(void)')[0]+r'''
static void shorts(const int16_t*p,int n){for(int i=0;i<n;i++){uint16_t v=p[i];uint8_t b[2]={v,v>>8};wr(b,2);}}
static void scaled(float*p,int n,unsigned pat,uint32_t seed,unsigned mark,unsigned shift){pattern(p,n,pat,seed,mark);float scale=1.f;for(unsigned i=0;i<shift;i++)scale*=.5f;for(int i=0;i<n;i++)p[i]*=scale;}
#define GUARD(x) (intact(x.pre,8)&&intact(x.post,8))
int main(void){int op;while((op=getchar())!=EOF){
 if(op==1){uint32_t bits=ru();float v;memcpy(&v,&bits,4);if((bits&0x7fffffff)>0x7f800000)return 4;wu((uint32_t)(int32_t)mp3d_scale_pcm(v));continue;}
 unsigned pat=ru(),seed=ru(),mark=ru(),nch=ru(),shift=ru();if(pat>5||mark>=2176||nch<1||nch>2||shift>24)return 5;
 struct{uint8_t pre[8];float v[1216];uint8_t post[8];} data;
 struct{uint8_t pre[8];float v[2176];uint8_t post[8];} lins;
 struct{uint8_t pre[8];float v[1024];uint8_t post[8];} qmf;
 struct{uint8_t pre[8];int16_t v[1216];uint8_t post[8];} pcm;
 memset(&data,0xcc,sizeof(data));memset(&lins,0xcc,sizeof(lins));memset(&qmf,0xcc,sizeof(qmf));memset(&pcm,0xcc,sizeof(pcm));
 scaled(data.v,1216,pat,seed,mark,shift);scaled(lins.v,2176,pat,seed^0x4c494e,mark,shift);scaled(qmf.v,1024,pat,seed^0x514d46,mark,shift);
 if(op==2){unsigned off=ru();if(off>125)return 6;float before[2176];memcpy(before,lins.v,sizeof(before));mp3d_synth_pair(pcm.v,nch,lins.v+off);shorts(pcm.v,96);wu(memcmp(before,lins.v,sizeof(before))==0);wu(GUARD(lins)&&GUARD(pcm));}
 else if(op==3){unsigned off=ru();if(off>16||(off&1))return 7;mp3d_synth(data.v+off,pcm.v,nch,lins.v);floats(data.v,1216);floats(lins.v,2176);shorts(pcm.v,160);wu(GUARD(data)&&GUARD(lins)&&GUARD(pcm));}
 else if(op==4){unsigned steps=ru();if(steps<1||steps>4)return 8;for(unsigned step=0;step<steps;step++){
  unsigned bands=ru(),channels=ru();if(bands>18||(bands&1)||channels<1||channels>2)return 9;
  scaled(data.v,1216,pat,seed+step*37,(mark+step*17)%1152,shift);
  memset(pcm.v,0xcc,sizeof(pcm.v));
  mp3d_synth_granule(qmf.v,data.v,bands,channels,pcm.v,lins.v);
  floats(data.v,1216);floats(lins.v,2176);floats(qmf.v,1024);shorts(pcm.v,1216);wu(GUARD(data)&&GUARD(lins)&&GUARD(qmf)&&GUARD(pcm));
 }}else return 10;
 }return ferror(stdin)||fflush(stdout);}
'''
def requests():
 # Every half-integer rounding boundary in the unclipped int16 domain, plus its
 # adjacent float32 values. Encode bits directly; never pass a NaN to C casts.
 bits={0,0x80000000,1,0x80000001,0x7f7fffff,0xff7fffff,0x7f800000,0xff800000}
 for i in range(-32768,32768):
  b=struct.unpack('<I',struct.pack('<f',i+.5))[0]
  bits.update((b-1,b,b+1))
 for b in sorted(bits):yield 1,common.u(b),4
 for pat in range(6):
  for seed in range(1 if pat<2 else 4):
   for nch in (1,2):
    for shift in (0,16):
     for off in (0,60,61,124,125):yield 2,common.u(pat,seed,seed,nch,shift,off),200
     for off in (0,2,16):yield 3,common.u(pat,seed,seed,nch,shift,off),13892
     for bands in (0,2,12,18):yield 4,common.u(pat,seed,seed,nch,shift,1,bands,nch),20100
     seq=[(18,nch),(18,nch),(18,3-nch),(18,nch)]
     yield 4,common.u(pat,seed,seed,nch,shift,len(seq))+b''.join(common.u(*p) for p in seq),20100*len(seq)
 # Each input row at both channel offsets, avoiding saturated-only coverage.
 for nch in (1,2):
  for ch in range(nch):
   for band in range(32):
    for col in (0,1):
     mark=ch*576+band*18+col
     yield 3,common.u(3,0,mark,nch,8,0),13892
 # Isolate every polyphase history coefficient in synth_pair's input.
 for mark in range(0,15*64,64):
  for delta in (0,2):yield 2,common.u(3,0,mark+delta,2,8,0),200

def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def main():
 ap=argparse.ArgumentParser(description=__doc__);ap.add_argument('output',type=Path);args=ap.parse_args()
 root=Path(__file__).resolve().parents[2];out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
 header=original_header(root,out);shim=out/'capture.c';shim.write_text(SHIM)
 flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse','-I',str(header.parent),str(shim),'-lm']
 subprocess.run(flags+['-o',str(out/'capture')],check=True)
 counts=[0]*4;steps=0
 with (out/'requests.bin').open('wb') as fo:
  for op,row,n in requests():
   fo.write(bytes((op,))+row);counts[op-1]+=1
   if op==4:steps+=n//20100
 hashes=[]
 for i in range(3):
  with (out/'requests.bin').open('rb') as fi,(out/f'results-{i}.bin').open('wb') as fo:subprocess.run([str(out/'capture')],stdin=fi,stdout=fo,check=True)
  hashes.append(sha(out/f'results-{i}.bin'))
 assert len(set(hashes))==1
 subprocess.run(flags+['-fsanitize=undefined,float-cast-overflow','-fno-sanitize-recover=all','-o',str(out/'capture-ubsan')],check=True)
 with (out/'requests.bin').open('rb') as fi,(out/'results-ubsan.bin').open('wb') as fo:subprocess.run([str(out/'capture-ubsan')],stdin=fi,stdout=fo,check=True)
 assert sha(out/'results-ubsan.bin')==hashes[0]
 h=hashlib.sha256()
 with (out/'results-0.bin').open('rb') as fi,(out/'synthesis_c.bin.gz').open('wb') as fo:
  with gzip.GzipFile(filename='',mode='wb',fileobj=fo,mtime=0) as gz:
   def write(b):h.update(b);gz.write(b)
   write(b'NMP3SYN1')
   for op,row,n in requests():
    expected=fi.read(n);assert len(expected)==n;write(bytes((op,))+row+expected)
   assert not fi.read(1)
 meta=dict(counts=counts,granule_steps=steps,result_sha256=hashes,ubsan_result_sha256=sha(out/'results-ubsan.bin'),uncompressed_sha256=h.hexdigest(),fixture_sha256=sha(out/'synthesis_c.bin.gz'),source_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),header_sha256=sha(header),tool_sha256=sha(__file__),pattern_tool_sha256=sha(common.__file__),shim_sha256=sha(shim),compiler=subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],flags=flags)
 (out/'capture.json').write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta,indent=2))
if __name__=='__main__':main()
