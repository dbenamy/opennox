#!/usr/bin/env python3
"""Freeze original scalar-SSE2 synthesis DCT with caller strides and retained tails."""
import argparse,gzip,hashlib,json,struct,subprocess
from pathlib import Path
import capture_mp3_imdct as common
SHIM=common.SHIM.split('int main(void)')[0]+r'''
int main(void){int op;while((op=getchar())!=EOF){
 if(op!=1)return 4;
 unsigned pat=ru(),seed=ru(),mark=ru(),off=ru(),n=ru();if(pat>5||mark>=1216||off>576||n>18)return 5;
 struct{uint8_t pre[8];float v[1216];uint8_t post[8];} data;
 memset(&data,0xcc,sizeof(data));pattern(data.v,1216,pat,seed,mark);
 mp3d_DCT_II(data.v+off,n);
 floats(data.v,1216);wu(intact(data.pre,8)&&intact(data.post,8));
 }return ferror(stdin)||fflush(stdout);}
'''
def requests():
 for pat in range(6):
  for seed in range(1 if pat<2 else 16):
   for off in (0,576):
    for n in (0,1,2,6,12,17,18):yield common.u(pat,seed,off+seed,off,n)
 # Every coordinate in either channel under the actual Layer III n=18 call.
 for off in (0,576):
  for mark in range(576):yield common.u(3,0,off+mark,off,18)
 # First/last active columns at the one-column and MPEG Layer I/II extents.
 for n in (1,12):
  for band in range(32):
   for col in sorted({0,n-1}):yield common.u(3,7,band*18+col,0,n)
def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def main():
 ap=argparse.ArgumentParser(description=__doc__);ap.add_argument('output',type=Path);args=ap.parse_args()
 root=Path(__file__).resolve().parents[2];out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
 header=root/'src/legacy/client/audio/mp3/minimp3.h';shim=out/'capture.c';shim.write_text(SHIM)
 flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse','-I',str(header.parent),str(shim),'-lm']
 subprocess.run(flags+['-o',str(out/'capture')],check=True)
 rows=list(requests());(out/'requests.bin').write_bytes(b''.join(b'\1'+row for row in rows))
 hashes=[]
 for i in range(3):
  with (out/'requests.bin').open('rb') as fi,(out/f'results-{i}.bin').open('wb') as fo:subprocess.run([str(out/'capture')],stdin=fi,stdout=fo,check=True)
  hashes.append(sha(out/f'results-{i}.bin'))
 assert len(set(hashes))==1
 subprocess.run(flags+['-fsanitize=undefined','-fno-sanitize-recover=all','-o',str(out/'capture-ubsan')],check=True)
 with (out/'requests.bin').open('rb') as fi,(out/'results-ubsan.bin').open('wb') as fo:subprocess.run([str(out/'capture-ubsan')],stdin=fi,stdout=fo,check=True)
 assert sha(out/'results-ubsan.bin')==hashes[0]
 h=hashlib.sha256()
 with (out/'results-0.bin').open('rb') as fi,(out/'dct_c.bin.gz').open('wb') as fo:
  with gzip.GzipFile(filename='',mode='wb',fileobj=fo,mtime=0) as gz:
   def write(b):h.update(b);gz.write(b)
   write(b'NMP3DCT1')
   for row in rows:
    expected=fi.read(4868);assert len(expected)==4868;write(b'\1'+row+expected)
   assert not fi.read(1)
 meta=dict(cases=len(rows),result_sha256=hashes,ubsan_result_sha256=sha(out/'results-ubsan.bin'),uncompressed_sha256=h.hexdigest(),fixture_sha256=sha(out/'dct_c.bin.gz'),source_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),header_sha256=sha(header),tool_sha256=sha(__file__),pattern_tool_sha256=sha(common.__file__),shim_sha256=sha(shim),compiler=subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],flags=flags)
 (out/'capture.json').write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta,indent=2))
if __name__=='__main__':main()
