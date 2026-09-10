#!/usr/bin/env python3
"""Summarize repo-owned C translation units selected by go list -deps -json.

Run once per target. This reports build selection, not runtime reachability.
C included by headers/preambles and external-library C are not counted as units.
"""
import argparse
import json
import subprocess
from pathlib import Path

parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument('graph', type=Path)
parser.add_argument('--binary', type=Path, required=True)
parser.add_argument('--symbol', action='append', default=[])
args = parser.parse_args()
text = args.graph.read_text()
decoder = json.JSONDecoder()
packages = []
while text.strip():
    package, end = decoder.raw_decode(text.lstrip())
    text = text.lstrip()[end:]
    name = package.get('ImportPath', '')
    if name.startswith('github.com/opennox/opennox/v1/') and package.get('CFiles'):
        packages.append({'package': name, 'c_files': sorted(package['CFiles'])})
symbols = subprocess.check_output(['nm', '--defined-only', str(args.binary)], text=True)
defined = {line.split()[-1] for line in symbols.splitlines() if line.split()}
print(json.dumps({
    'translation_units': sum(len(p['c_files']) for p in packages),
    'packages': sorted(packages, key=lambda p: p['package']),
    'requested_symbols_retained': {s: s in defined for s in args.symbol},
}, indent=2))
