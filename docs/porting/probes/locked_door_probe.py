from pathlib import Path
import subprocess,json,hashlib,tempfile,argparse
parser=argparse.ArgumentParser(description="Compile the current C locked-door routine with a recording transport, without retaining a C algorithm copy.")
parser.add_argument("--expect-zero", action="store_true")
parser.add_argument("--source-ref", help="Read GAME4.c from a git revision instead of the working tree")
args=parser.parse_args()
work=tempfile.TemporaryDirectory(prefix="opennox-door-probe-")
p=Path(work.name)
source=(subprocess.check_output(['git','show',args.source_ref+':src/legacy/GAME4.c'],text=True) if args.source_ref else Path('src/legacy/GAME4.c').read_text())
a=source.index('void sub_4FADD0(');b=source.index('//-----',a)
body=source[a:b]
harness=r'''
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
static unsigned char obj[780],upd[556],player[4828],before[6164];
static const char *key;
static int expected,selector,recipient,sends,cases,padding_nonzero;
static void check(int ok){if(!ok) abort();}
static int nox_xxx_netSendPacket0_4E5420(int who,const void* raw,int n,int x,int y){
 const unsigned char *p=raw; ++sends;
 check(expected && who==recipient && n==52 && x==0 && y==1);
 check(p[0]==240 && p[1]==33 && p[51]==selector);
 check(memcmp(p+2,key,strlen(key)+1)==0);
 for(size_t i=strlen(key)+3;i<51;i++) if(p[i]) {padding_nonzero++;break;}
 printf("%d:%d:%d:%d:%d:%d:%s\n",cases,who,n,x,y,selector,key);
 return 17;
}
BODY
int main(){
 const char *texts[]={NULL,"","a","objcoll.c:GateLockedKey","objcoll.c:DoorLockedKey",NULL,NULL,NULL,NULL};
 char strings[4][51];int lengths[]={47,48,49,50};
 for(int i=0;i<4;i++){memset(strings[i],'A'+i,lengths[i]);strings[i][lengths[i]]=0;texts[5+i]=strings[i];}
 for(int repeat=0;repeat<8;repeat++)for(int mode=0;mode<4;mode++)for(int k=0;k<9;k++)for(int sel=0;sel<5;sel++){
  memset(obj,0xa5,sizeof(obj));memset(upd,0x5a,sizeof(upd));memset(player,0x69,sizeof(player));
  *(uint32_t*)(obj+8)=mode==2?2:(mode==3?0x80000004u:4);*(uint32_t*)(obj+748)=(uintptr_t)upd;
  *(uint32_t*)(upd+276)=(uintptr_t)player;recipient=repeat==7?255:repeat;player[2064]=recipient;
  memcpy(before,obj,sizeof(obj));memcpy(before+sizeof(obj),upd,sizeof(upd));memcpy(before+sizeof(obj)+sizeof(upd),player,sizeof(player));
  key=texts[k];selector=sel;expected=mode!=0 && mode!=2 && key && strlen(key)>0 && strlen(key)<=48;
  int old=sends;sub_4FADD0(mode==0?0:(int)(uintptr_t)obj,key,sel);check(sends-old==expected);
  check(!memcmp(before,obj,sizeof(obj)) && !memcmp(before+sizeof(obj),upd,sizeof(upd)) && !memcmp(before+sizeof(obj)+sizeof(upd),player,sizeof(player)));
  cases++;
 }
 fprintf(stderr,"{\"cases\":%d,\"messages\":%d,\"messages_with_nonzero_padding\":%d}\n",cases,sends,padding_nonzero);
}
'''
(p/'probe.c').write_text(harness.replace('BODY',body))
subprocess.run(['gcc','-m32','-O2','-fno-strict-aliasing','-fno-strict-overflow','-D_FORTIFY_SOURCE=2','-fstack-protector-strong',str(p/'probe.c'),'-o',str(p/'probe')],check=True)
r=subprocess.run([str(p/'probe')],capture_output=True,check=True)
q=json.loads(r.stderr);q['defined_fields_sha256']=hashlib.sha256(r.stdout).hexdigest()
assert q['cases']==1440 and q['messages']==400
assert q['defined_fields_sha256']=='f628c861746ff1484f294b803a56c51fc38badbd792037b1d092dd4b3f87308f'
if args.expect_zero: assert q['messages_with_nonzero_padding']==0
print(json.dumps(q,indent=2))
work.cleanup()
