"""Recover the retired, hash-pinned MP3 C header for historical baseline tools.

No C implementation is retained in the working source tree. Captures recover the
original from local Git history into their ignored output directory when needed.
"""
from functools import lru_cache
import hashlib
from pathlib import Path
import subprocess

REVISION = '6702100ba617a26cf952b2f8ba629a6f4256426e'
RELATIVE = 'src/legacy/client/audio/mp3/minimp3.h'
SHA256 = 'e03e87c847cbdfd79cd3bf74a3b1fe24bc8c590245e92ce8153e87b23e17507f'

@lru_cache(maxsize=1)
def original_header_bytes(root):
    root = Path(root)
    path = root / RELATIVE
    data = path.read_bytes() if path.exists() else subprocess.check_output(
        ['git', 'show', f'{REVISION}:{RELATIVE}'], cwd=root)
    if hashlib.sha256(data).hexdigest() != SHA256:
        raise RuntimeError('MP3 capture header differs from the qualified original')
    return data

def original_header_text(root):
    return original_header_bytes(root).decode('utf-8')

def original_header(root, output):
    source = Path(root) / RELATIVE
    data = original_header_bytes(root)
    if source.exists():
        return source
    recovered = Path(output) / 'original-c' / 'minimp3.h'
    recovered.parent.mkdir(parents=True, exist_ok=True)
    recovered.write_bytes(data)
    return recovered
