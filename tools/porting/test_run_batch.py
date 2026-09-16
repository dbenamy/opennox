import contextlib
import hashlib
import io
import json
from pathlib import Path
import sys
import tempfile
import unittest
from unittest.mock import patch
import run_batch


class ManifestTests(unittest.TestCase):
    def run_manifest(self, command, hashes=(), changed=False):
        with tempfile.TemporaryDirectory() as tmp:
            p = Path(tmp)
            manifest = dict(batch='runner-test', phases={'test':[
                dict(name='probe', command=command, hashes=list(hashes))]})
            (p/'manifest.json').write_text(json.dumps(manifest))
            with patch('sys.argv', ['run_batch.py', str(p/'manifest.json'), '--phase','test','--out',str(p/'out')]), \
                 patch.object(run_batch, 'fingerprints', side_effect=[{'s':'a'}, {'s':'b' if changed else 'a'}, {'s':'b' if changed else 'a'}]), \
                 contextlib.redirect_stdout(io.StringIO()), contextlib.redirect_stderr(io.StringIO()):
                code = run_batch.main()
            return code, json.loads((p/'out/result.json').read_text())

    def test_matching_artifact(self):
        code, result = self.run_manifest([sys.executable,'-c',"from pathlib import Path; Path('{out}/capture').write_bytes(b'abc')"],
            [dict(path='{out}/capture', sha256=hashlib.sha256(b'abc').hexdigest())])
        self.assertEqual(code, 0)
        self.assertTrue(result['source_unchanged'])

    def test_mismatch_fails(self):
        code, result = self.run_manifest([sys.executable,'-c',"from pathlib import Path; Path('{out}/capture').write_bytes(b'bad')"],
            [dict(path='{out}/capture', sha256=hashlib.sha256(b'abc').hexdigest())])
        self.assertEqual(code, 1)
        self.assertIn('mismatch', result['error'])

    def test_source_change_fails(self):
        self.assertEqual(self.run_manifest([sys.executable,'-c','pass'], changed=True)[0], 1)

    def test_command_failure_fails(self):
        self.assertEqual(self.run_manifest([sys.executable,'-c','raise SystemExit(3)'])[0], 1)


if __name__ == '__main__':
    unittest.main()
