#!/usr/bin/env python3
"""Execute a reviewed batch manifest, recording commands, timing and source identity.

No shell expansion, implicit retries, source mutation, capture regeneration or
approval decisions. Each invocation gets a new output directory. Source the port
environment before invoking. See docs/porting/PROCESS_TRIAL.md.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import time

ROOT = Path(__file__).resolve().parents[2]


def fingerprints():
    return {str(p.relative_to(ROOT)): hashlib.sha256(p.read_bytes()).hexdigest()
            for p in sorted((ROOT / 'src').rglob('*'))
            if p.is_file() and p.suffix in ('.go', '.c', '.h', '.mod', '.sum')}


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('manifest', type=Path)
    ap.add_argument('--phase', required=True)
    ap.add_argument('--out', type=Path, required=True)
    args = ap.parse_args()
    manifest = json.loads(args.manifest.read_text())
    steps = manifest['phases'][args.phase]
    out = args.out.resolve()
    out.mkdir(parents=True, exist_ok=False)
    substitutions = dict(root=str(ROOT), out=str(out), phase=args.phase)
    expand = lambda value: value.format_map(substitutions)
    env = dict(os.environ, **{k: expand(v) for k, v in manifest.get('env', {}).items()})
    source = fingerprints()
    (out / 'source.json').write_text(json.dumps(source, indent=2) + '\n')
    (out / 'manifest.json').write_text(json.dumps(manifest, indent=2) + '\n')
    report = dict(batch=manifest['batch'], phase=args.phase, success=False, steps=[],
                  revision=subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=ROOT, text=True).strip())
    start = time.monotonic()
    try:
        for step in steps:
            name = step['name']
            if not name.replace('-', '').replace('_', '').isalnum():
                raise ValueError('step names must be simple filenames')
            command = [expand(v) for v in step['command']]
            now = time.monotonic()
            with (out / (name + '.log')).open('w') as log:
                proc = subprocess.run(command, cwd=expand(step.get('cwd', '{root}')),
                                      env=env, stdout=log, stderr=subprocess.STDOUT)
            entry = dict(name=name, command=command, exit=proc.returncode,
                         wall_seconds=time.monotonic() - now)
            report['steps'].append(entry)
            print(json.dumps(entry), flush=True)
            if proc.returncode != 0:
                raise RuntimeError(f'{name} failed; see {out / (name + ".log")}')
            for spec in step.get('hashes', []):
                path = Path(expand(spec['path']))
                if hashlib.sha256(path.read_bytes()).hexdigest() != spec['sha256']:
                    raise RuntimeError(f'frozen capture mismatch: {path}')
            if fingerprints() != source:
                raise RuntimeError('source changed while validation was running')
        report['success'] = True
    except Exception as exc:
        report['error'] = str(exc)
        print(str(exc), file=sys.stderr)
    finally:
        report['source_unchanged'] = fingerprints() == source
        report['success'] = report['success'] and report['source_unchanged']
        report['wall_seconds'] = time.monotonic() - start
        report['c_loc'] = json.loads(subprocess.check_output(
            [sys.executable, str(ROOT / 'tools/porting/c_loc.py')], cwd=ROOT))
        (out / 'result.json').write_text(json.dumps(report, indent=2) + '\n')
    return 0 if report['success'] else 1


if __name__ == '__main__':
    sys.exit(main())
