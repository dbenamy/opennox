#!/usr/bin/env python3
"""Run a reference gameplay scenario on fresh assets under headless X.

Optional forced map expansion removes only files in the new run copy, then
requires a loaded map to be regenerated exactly. Original assets are read-only.
"""
from pathlib import Path
import hashlib
import json
import os
import shutil
import subprocess
import sys
import time

ROOT = Path(__file__).resolve().parents[2]
name, mode = sys.argv[1:3]
if not name.replace('-', '').replace('_', '').isalnum() or mode not in ('capture','compare'):
    raise SystemExit('invalid run name or mode')
reference = Path(sys.argv[3]) if len(sys.argv)>3 else None
if mode == 'compare' and reference is None:
    raise SystemExit('comparison requires a reference')
run = ROOT/'build/baseline/runs'/name
run.mkdir(exist_ok=False)
shutil.copytree(ROOT/'build/assets/extracted/drive_c/Nox',run/'data',symlinks=True)
# Some shipped campaign maps have no compressed counterpart. Generate only
# requested run-copy counterparts with the unchanged production compressor.
generated=[]
for relative in json.loads(os.environ.get('OPENNOX_COMPRESS_MAPS','[]')):
    path=run/'data/maps'/relative
    if Path(relative).is_absolute() or '..' in Path(relative).parts or path.is_symlink():
        raise SystemExit('invalid run-copy map path')
    subprocess.run([os.environ['OPENNOX_MAP_COMPRESSOR'],str(path),str(path.with_suffix('.nxz'))],check=True)
    generated.append(relative)
(run/'generated-compressed-maps.json').write_text(json.dumps(generated,indent=2)+'\n')
removed=[]
if os.environ.get('OPENNOX_REQUIRE_MAP_DECOMPRESSION') == '1':
    for compressed in sorted((run/'data/maps').glob('*/*.nxz')):
        expanded = compressed.with_suffix('.map')
        if expanded.is_file() and not expanded.is_symlink():
            removed.append(dict(path=str(expanded.relative_to(run/'data')),
                                sha256=hashlib.sha256(expanded.read_bytes()).hexdigest()))
            expanded.unlink()
    (run/'removed-maps.json').write_text(json.dumps(removed,indent=2)+'\n')
    if not removed:
        raise SystemExit('no copied maps available for forced decompression')
for path in list((run/'data').iterdir()):
    if path.name.lower() == 'save':
        shutil.rmtree(path)
(run/'data/save').mkdir()
if reference:
    shutil.copytree(reference/'testdata',run/'testdata')
else:
    (run/'testdata').mkdir()
shutil.copyfile(Path(os.environ['OPENNOX_UI_SCENARIO']),run/'scenario.yaml')
binary = Path(os.environ['OPENNOX_DISPLAY_BINARY'])
env = dict(os.environ,GODEBUG=os.environ.get('GODEBUG','')+',randautoseed=0',
           NOX_E2E=str(run/'scenario.yaml'),ALSOFT_DRIVERS='null',
           NOX_E2E_OVERRIDE='true' if mode=='capture' else 'false')
cmd = ['xvfb-run','-a','-s','-screen 0 1280x960x24 -nolisten tcp',str(binary),
       '-data',str(run/'data'),'-window','-pprof','127.0.0.1:0']
start = time.monotonic()
with (run/'output.log').open('w') as log:
    try:
        code = subprocess.run(cmd,env=env,stdout=log,stderr=subprocess.STDOUT,timeout=120).returncode
    except subprocess.TimeoutExpired:
        code = 124
report = dict(exit=code,process_exit=code,elapsed=time.monotonic()-start,capture=mode=='capture',
              godebug=env['GODEBUG'],reference=str(reference) if reference else None,
              binary_sha256=hashlib.sha256(binary.read_bytes()).hexdigest(),command=cmd,
              display_implementation=os.environ.get('OPENNOX_DISPLAY_IMPLEMENTATION',''),
              removed_maps=len(removed))
if code == 0 and removed:
    regenerated=[]
    for item in removed:
        path = run/'data'/item['path']
        if not path.exists() and path.parent.is_dir():
            matches=[p for p in path.parent.iterdir() if p.name.casefold()==path.name.casefold()]
            if len(matches)>1:
                report['error']='ambiguous regenerated map name: '+item['path']
                code=1
                break
            if matches:
                path=matches[0]
        if path.exists():
            if hashlib.sha256(path.read_bytes()).hexdigest()!=item['sha256']:
                report['error']='regenerated map mismatch: '+item['path']
                code=1
                break
            regenerated.append(dict(item, actual_path=str(path.relative_to(run/"data"))))
    if not regenerated:
        report['error']='scenario did not regenerate any compressed map'
        code=1
    (run/'regenerated-maps.json').write_text(json.dumps(regenerated,indent=2)+'\n')
    report['regenerated_maps']=len(regenerated)
report['exit']=code
(run/'result.json').write_text(json.dumps(report,indent=2)+'\n')
print(name,'exit',code,flush=True)
# Run optional asset deduplication separately after inspecting this result. Local
# maintenance failures must not change the gameplay validation process status.
raise SystemExit(code)
