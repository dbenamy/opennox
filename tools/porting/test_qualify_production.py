"""Production evidence reuse must fail when its qualifying conditions change."""
import contextlib
import hashlib
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch
import qualify_production


class EvidenceReuse(unittest.TestCase):
    def attempt(self, change=None, source=None):
        with tempfile.TemporaryDirectory() as tmp:
            root=Path(tmp)
            previous=root/'old'/'production'
            (previous/'bin').mkdir(parents=True)
            spec=dict(retained=[],retained_c=['compress'],retired=['decompress'],known_suite='known.log',scenarios=[])
            (previous.parent/'manifest.json').write_text(json.dumps(dict(production=spec)))
            (previous.parent/'source.json').write_text(json.dumps({'source':'same'}))
            builds={}
            for name in ['opennox','opennox-hd','opennox-server']:
                data=b'qualified artifact'
                (previous/'bin'/name).write_bytes(data)
                builds[name]=dict(exit=0,abi_verified=True,sha256=hashlib.sha256(data).hexdigest())
            saved=dict(builds=builds,full_suite=dict(exit=1,exact_failure_multiset_match=True,exact_package_result_match=True))
            if change:
                change(spec,saved)
            (previous/'production.json').write_text(json.dumps(saved))
            manifest=root/'batch.json'
            manifest.write_text(json.dumps(dict(production=spec)))
            with patch('sys.argv',['qualify_production.py',str(manifest),'--out',str(root/'new'),
                                   '--reuse-production-from',str(previous)]), \
                 patch.object(qualify_production,'fingerprints',return_value=source or {'source':'same'}), \
                 contextlib.redirect_stderr(io.StringIO()):
                code=qualify_production.main()
            return code,json.loads((root/'new/production.json').read_text())

    def test_exact_evidence_is_reused(self):
        code,result=self.attempt()
        self.assertEqual(code,0)
        self.assertTrue(result['success'])
        self.assertIn('reused_production_from',result)

    def test_changed_source_rejected(self):
        self.assertEqual(self.attempt(source={'source':'changed'})[0],1)

    def test_changed_abi_gate_rejected(self):
        self.assertEqual(self.attempt(change=lambda spec,saved:spec['retired'].append('another'))[0],1)

    def test_changed_binary_rejected(self):
        self.assertEqual(self.attempt(change=lambda spec,saved:saved['builds']['opennox'].update(sha256='wrong'))[0],1)

    def test_unqualified_suite_rejected(self):
        self.assertEqual(self.attempt(change=lambda spec,saved:saved['full_suite'].update(exact_failure_multiset_match=False))[0],1)


if __name__=='__main__':
    unittest.main()
