#!/usr/bin/env python3
"""Count physical lines in tracked src/**/*.c, with porttest references separate.

Blank lines and comments count; headers, C embedded in Go, and dependencies do
not. This is a source-size checkpoint, not a measure of active code or effort.
Pass a Git revision to count history; omit it to count the current tracked files.
"""
import argparse
import json
import subprocess
from pathlib import Path

parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument('revision', nargs='?')
args = parser.parse_args()
root = Path(subprocess.check_output(['git', 'rev-parse', '--show-toplevel'], text=True).strip())
if args.revision:
    revision = subprocess.check_output(['git', 'rev-parse', '--verify', args.revision + '^{commit}'], cwd=root, text=True).strip()
    paths = subprocess.check_output(['git', 'ls-tree', '-r', '--name-only', '-z', revision, '--', 'src'], cwd=root).split(b'\0')
else:
    revision = 'working tree (tracked files)'
    paths = subprocess.check_output(['git', 'ls-files', '-z', '--', 'src'], cwd=root).split(b'\0')
counts = {name: {'files': 0, 'lines': 0} for name in ('production', 'test_reference')}
for raw in paths:
    path = raw.decode()
    if not path.endswith('.c'):
        continue
    if args.revision:
        data = subprocess.check_output(['git', 'show', revision + ':' + path], cwd=root)
    else:
        file = root / path
        if not file.exists():
            continue
        data = file.read_bytes()
    group = 'test_reference' if data.startswith(b'//go:build porttest\n') else 'production'
    counts[group]['files'] += 1
    counts[group]['lines'] += len(data.splitlines())
print(json.dumps({'revision': revision, **counts}, indent=2))
