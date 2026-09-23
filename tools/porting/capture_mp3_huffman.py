#!/usr/bin/env python3
"""Freeze actual scalar Layer III Huffman/dequantization state and power helper."""
import argparse,gzip,hashlib,json,struct,subprocess,re
from pathlib import Path
from mp3_source import original_header, original_header_text
import capture_mp3_imdct as common
import capture_mp3_spectrum as spectrum
SHIM=common.SHIM.split('int main(void)')[0]+r'''
static void bytepattern(uint8_t*p,unsigned pat,uint32_t seed){uint32_t s=seed?seed:1;for(int i=0;i<2815;i++){s^=s<<13;s^=s>>17;s^=s<<5;p[i]=pat==0?0:pat==1?255:pat==2?0xaa:s>>24;}}
int main(void){int op;while((op=getchar())!=EOF){
 if(op==1){int x=(int32_t)ru();if(x< -16||x>8206)return 4;float v=L3_pow_43(x);floats(&v,1);continue;}
 if(op!=2&&op!=3)return 5;
 uint8_t hdr[4],side[64];rd(hdr,4);rd(side,64);L3_gr_info_t gr[4];memset(gr,0,sizeof(gr));bs_t bs;bs_init(&bs,side,64);if(!hdr_valid(hdr)||L3_read_side_info(&bs,gr,hdr)<0)return 6;
 unsigned tab=ru(),big=ru(),count1=ru(),pos=ru(),length=ru(),pat=ru(),seed=ru(),scpat=ru();
 if(tab>31||big>288||count1>1||pos>22480||length>4095||pat>3||scpat>2||pos+length>22520)return 7;
 gr[0].big_values=big;gr[0].count1_table=count1;
 gr[0].table_select[0]=tab;gr[0].table_select[1]=(tab+7)%32;gr[0].table_select[2]=(tab+13)%32;
 gr[0].region_count[0]=0;gr[0].region_count[1]=0;gr[0].region_count[2]=255;
 L3_gr_info_t before=gr[0];
 struct{uint8_t pre[8],v[2815],post[8];} input;memset(&input,0xcc,sizeof(input));bytepattern(input.v,pat,seed);if(op==3){unsigned bytes=ru();if(bytes>64)return 8;rd(input.v,bytes);}uint8_t oldinput[2815];memcpy(oldinput,input.v,2815);
 struct{uint8_t pre[8];float v[640];uint8_t post[8];} dst;memset(&dst,0xcc,sizeof(dst));for(int i=0;i<640;i++){uint32_t b=0x4b123456+i;memcpy(dst.v+i,&b,4);}
 float scf[40],oldscf[40];for(int i=0;i<40;i++){uint32_t b=scpat==0?0x3f800000:scpat==1?(uint32_t)((110+i%24)<<23):((i&1)?0x80000000:1);memcpy(scf+i,&b,4);}memcpy(oldscf,scf,sizeof(scf));
 // Keep the actual physical scratch capacity visible beyond the logical bit limit.
 bs_init(&bs,input.v,2815);bs.pos=pos;bs.limit=pos+length;
 L3_huffman(dst.v,&bs,&gr[0],scf,pos+length);
 floats(dst.v,640);wu(bs.pos);wu(bs.limit);wu(memcmp(oldinput,input.v,2815)==0);wu(memcmp(oldscf,scf,sizeof(scf))==0);wu(memcmp(&before,&gr[0],sizeof(before))==0);wu(intact(input.pre,8)&&intact(input.post,8)&&intact(dst.pre,8)&&intact(dst.post,8));
 }return ferror(stdin)||fflush(stdout);}
'''
def requests():
 for x in range(-16,8207):yield 1,common.u(x),4
 base=spectrum.setup(0xfb,0,0)
 for tab in range(32):
  for pos in range(8):
   for big in (0,1,2,7,287,288):
    for pat in (0,1,3):
     for count1 in (0,1):
      yield 2,base+common.u(tab,big,count1,pos,4095,pat,tab*17+pos+1,1),2584
 # Every actual scalefactor row/layout and scale pattern, including count1 cutoff.
 for version in (0xe3,0xf3,0xfb):
  for rate in range(3):
   for layout in range(3):
    ctx=spectrum.setup(version,rate,layout)
    for big in (0,1,143,288):
     for scpat in range(3):
      for length in (0,1,7,31,255,4095):
       for count1 in (0,1):yield 2,ctx+common.u(31,big,count1,3,length,3,rate+layout*19+version,scpat),2584
 # Count1 lookahead near physical end with zero/short logical ranges.
 for pos in (8,15,16383,22400,22479):
  for length in (0,1,7):
   for pat in range(4):
    for count1 in (0,1):yield 2,base+common.u(0,0,count1,pos,length,pat,3,0),2584

 # Enumerate codebook leaves to construct inputs, never expected decoded values.
 source=original_header_text(Path(__file__).resolve().parents[2])
 def table(name):
  text=source[source.index(name+'['):];text=text[text.index('{')+1:text.index('}')]
  return [int(v) for v in re.findall(r'-?\d+',text)]
 tabs,index,linbits=table('tabs'),table('tabindex'),table('g_linbits')
 for tab in range(32):
  book=tabs[index[tab]:];leaves=set()
  def walk(base,width,prefix):
   for i in range(1<<width):
    leaf=book[base+i];bits=f'{i:0{width}b}'
    if leaf<0:walk(-(leaf>>3),leaf&7,prefix+bits)
    else:leaves.add((prefix+bits[:leaf>>8],leaf&255))
  walk(0,5,'')
  for code,leaf in sorted(leaves):
   escaped=linbits[tab] and (leaf&15==15 or leaf>>4==15)
   for extmax in ((0,1) if escaped else (0,)):
    for signs in ((0,1,2,3) if escaped else (0,3)):
     bits=code
     for j in range(2):
      value=(leaf>>(4*j))&15
      if value==15 and linbits[tab]:bits+=('1' if extmax else '0')*linbits[tab]
      if value:bits+=str((signs>>j)&1)
     length=len(bits);encoded=int(bits.ljust((length+7)//8*8,'0') or '0',2).to_bytes((length+7)//8,'big')
     yield 3,base+common.u(tab,1,0,0,length,0,1,0)+common.u(len(encoded))+encoded,2584

def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def main():
 ap=argparse.ArgumentParser(description=__doc__);ap.add_argument('output',type=Path);args=ap.parse_args()
 root=Path(__file__).resolve().parents[2];out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
 header=original_header(root,out);shim=out/'capture.c';shim.write_text(SHIM)
 flags=['gcc','-m32','-std=c11','-O2','-g','-msse2','-mfpmath=sse','-I',str(header.parent),str(shim),'-lm']
 subprocess.run(flags+['-o',str(out/'capture')],check=True);counts=[0,0,0]
 with (out/'requests.bin').open('wb') as fo:
  for op,row,n in requests():fo.write(bytes((op,))+row);counts[op-1]+=1
 hashes=[]
 for i in range(3):
  with (out/'requests.bin').open('rb') as fi,(out/f'results-{i}.bin').open('wb') as fo:subprocess.run([str(out/'capture')],stdin=fi,stdout=fo,check=True)
  hashes.append(sha(out/f'results-{i}.bin'))
 assert len(set(hashes))==1
 subprocess.run(flags+['-fsanitize=undefined','-fno-sanitize-recover=all','-o',str(out/'capture-ubsan')],check=True)
 with (out/'requests.bin').open('rb') as fi,(out/'results-ubsan.bin').open('wb') as fo:subprocess.run([str(out/'capture-ubsan')],stdin=fi,stdout=fo,check=True)
 assert sha(out/'results-ubsan.bin')==hashes[0]
 h=hashlib.sha256()
 with (out/'results-0.bin').open('rb') as fi,(out/'huffman_c.bin.gz').open('wb') as fo:
  with gzip.GzipFile(filename='',mode='wb',fileobj=fo,mtime=0) as gz:
   def write(b):h.update(b);gz.write(b)
   write(b'NMP3HUF1')
   for op,row,n in requests():
    expected=fi.read(n);assert len(expected)==n;write(bytes((op,))+row+expected)
   assert not fi.read(1)
 meta=dict(counts=counts,result_sha256=hashes,ubsan_result_sha256=sha(out/'results-ubsan.bin'),uncompressed_sha256=h.hexdigest(),fixture_sha256=sha(out/'huffman_c.bin.gz'),source_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),header_sha256=sha(header),tool_sha256=sha(__file__),pattern_tool_sha256=sha(common.__file__),side_setup_tool_sha256=sha(spectrum.__file__),side_payload_tool_sha256=sha(Path(spectrum.__file__).with_name('capture_mp3_sideinfo.py')),shim_sha256=sha(shim),compiler=subprocess.check_output(['gcc','--version'],text=True).splitlines()[0],flags=flags)
 (out/'capture.json').write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta,indent=2))
if __name__=='__main__':main()
