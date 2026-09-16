#!/usr/bin/env python3
"""Build/ABI, exact known-suite comparison, and headless scenarios for a batch.

Invoked as a manifest step after focused/affected tests. The manifest's production
section declares retained/retired C symbols, known-suite log, and scenario/reference
pairs. Raw artifacts remain in the supplied output directory.
"""
import argparse
import collections
import hashlib
import json
import os
import shutil
from run_batch import fingerprints
from pathlib import Path
import struct
import subprocess
import sys
import time

ROOT = Path(__file__).resolve().parents[2]


def suite_events(path):
    failures, packages = collections.Counter(), collections.Counter()
    for line in path.read_text().splitlines():
        try:
            e = json.loads(line)
        except json.JSONDecodeError:
            continue
        action = e.get('Action')
        if action == 'fail':
            failures[e.get('Package', ''), e.get('Test', '')] += 1
        if action in ('pass', 'fail', 'skip') and not e.get('Test'):
            packages[action] += 1
    return failures, packages


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('manifest', type=Path)
    ap.add_argument('--out', type=Path, required=True)
    ap.add_argument("--reuse-production-from", type=Path)
    args = ap.parse_args()
    spec = json.loads(args.manifest.read_text())['production']
    out = args.out.resolve()
    out.mkdir(parents=True, exist_ok=True)
    (out/'bin').mkdir(exist_ok=True)
    report = {'success': False, 'builds': {}, 'scenarios': {}}
    def save():
        (out/'production.json').write_text(json.dumps(report, indent=2)+'\n')
    def run(name, command, cwd=ROOT, env=None):
        now = time.monotonic()
        with (out/(name+'.log')).open('w') as log:
            p = subprocess.run(command, cwd=cwd, env=env, stdout=log, stderr=subprocess.STDOUT)
        return {'exit': p.returncode, 'wall_seconds': time.monotonic()-now}
    try:
        if args.reuse_production_from:
            previous=args.reuse_production_from.resolve()
            saved=json.loads((previous/'production.json').read_text())
            if fingerprints()!=json.loads((previous.parent/'source.json').read_text()):
                raise RuntimeError('source differs from reused production evidence')
            previous_spec=json.loads((previous.parent/'manifest.json').read_text())['production']
            for field in ['retained','retired','retained_c','known_suite','assets']:
                if previous_spec.get(field)!=spec.get(field):
                    raise RuntimeError('requested gate differs from reused evidence: '+field)
            full=saved['full_suite']
            if full['exit']!=1 or not full['exact_failure_multiset_match'] or not full['exact_package_result_match']:
                raise RuntimeError('reused full suite is not qualified')
            for name in ['opennox','opennox-hd','opennox-server']:
                row=saved['builds'][name]
                binary=previous/'bin'/name
                if row['exit']!=0 or not row['abi_verified'] or hashlib.sha256(binary.read_bytes()).hexdigest()!=row['sha256']:
                    raise RuntimeError('reused binary evidence mismatch')
                # Qualified binaries are immutable artifacts; share their storage.
                try:
                    os.link(binary,out/'bin'/name)
                except OSError:
                    shutil.copyfile(binary,out/'bin'/name)
                    shutil.copymode(binary,out/'bin'/name)
            report['builds']=saved['builds']
            report['full_suite']=full
            report['reused_production_from']=str(previous)
            save()
        else:
            for name, tags in [('opennox',''), ('opennox-hd','highres'), ('opennox-server','server')]:
                binary = out/'bin'/name
                cmd = ['go','build','-p','2','-o',str(binary)]
                if tags:
                    cmd += ['-tags',tags]
                result = run(name, cmd+['./cmd/opennox'], cwd=ROOT/'src')
                report['builds'][name] = result
                save()
                if result['exit']:
                    raise RuntimeError(name+' build failed')
                data = binary.read_bytes()
                if data[:5] != b'\x7fELF\x01' or struct.unpack_from('<H',data,18)[0] != 3:
                    raise RuntimeError('wrong binary target')
                nm = subprocess.check_output(['go','tool','nm',str(binary)], text=True)
                symbols = {line.split()[-1] for line in nm.splitlines() if line.split()}
                for symbol in spec['retained']:
                    if symbol not in symbols or not any('_cgoexp_' in s and s.endswith('_'+symbol) for s in symbols):
                        raise RuntimeError('missing Go-backed C export: '+symbol)
                for symbol in spec.get('retained_c', []):
                    if symbol not in symbols:
                        raise RuntimeError('missing live C interface: '+symbol)
                for symbol in spec['retired']:
                    if symbol in symbols:
                        raise RuntimeError('retired C symbol remains: '+symbol)
                for marker in ['PortTest', 'portTestTrade', 'portTestShop', 'nox_porttest_client_sound', 'portTestSpellForceCollect']:
                    if marker in nm:
                        raise RuntimeError('test helper in production: '+marker)
                info = subprocess.check_output(['go','version','-m',str(binary)], text=True)
                if 'porttest' in info or any(v not in info for v in ['GOARCH=386','GO386=sse2','CGO_ENABLED=1']):
                    raise RuntimeError('wrong production build settings')
                result.update(sha256=hashlib.sha256(data).hexdigest(), abi_verified=True)
                save()
            env = dict(os.environ, NOX_DATA=str(ROOT/spec['assets']))
            full = run('full-suite', ['go','test','-p','2','-count=1','-json','./...'], cwd=ROOT/'src', env=env)
            old, old_pkgs = suite_events(ROOT/spec['known_suite'])
            new, new_pkgs = suite_events(out/'full-suite.log')
            full.update(failure_entries=sum(new.values()), packages=dict(new_pkgs),
                        exact_failure_multiset_match=new==old, exact_package_result_match=new_pkgs==old_pkgs)
            report['full_suite'] = full
            save()
            if full['exit'] != 1 or not old or new != old or new_pkgs != old_pkgs:
                raise RuntimeError('full suite differs from known baseline')
        for scenario in spec['scenarios']:
            name = scenario['name']
            result = run(name, [sys.executable, str(ROOT/spec.get('scenario_runner','build/port-client-shop-ui/run-scenario.py')),
                               name, 'compare', str(ROOT/scenario['reference'])],
                         env=dict(os.environ, **scenario.get('env', {}), OPENNOX_COMPRESS_MAPS=json.dumps(spec.get('compress_maps',[])), OPENNOX_MAP_COMPRESSOR=str(ROOT/spec.get('map_compressor','build/port-map-decompression/map-compress')), OPENNOX_REQUIRE_MAP_DECOMPRESSION='1' if spec.get('force_map_decompression') else '0', OPENNOX_DISPLAY_BINARY=str(out/'bin/opennox'),
                                  OPENNOX_DISPLAY_IMPLEMENTATION=spec['description'],
                                  OPENNOX_UI_SCENARIO=str(ROOT/scenario['scenario'])))
            report['scenarios'][name] = result
            save()
            if result['exit']:
                raise RuntimeError(name+' scenario failed')
            result['capture'] = json.loads((ROOT/'build/baseline/runs'/name/'result.json').read_text())
        report['success'] = True
    except Exception as exc:
        report['error'] = str(exc)
        print(str(exc), file=sys.stderr)
    finally:
        save()
    return 0 if report['success'] else 1


if __name__ == '__main__':
    sys.exit(main())
