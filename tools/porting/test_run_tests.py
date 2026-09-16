"""Fail-closed tests for test discovery/accounting (no game compilation)."""
import contextlib
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch
import run_tests


class Accounting(unittest.TestCase):
    def run_driver(self, discovered, events, code=0, extra=()):
        with tempfile.TemporaryDirectory() as tmp:
            p = Path(tmp)
            (p / 'pattern').write_text('^TestExample$')
            def discover(command, **kw):
                self.assertIn('-json', command)
                for pkg, test in discovered:
                    kw['stdout'].write(json.dumps(dict(Package=pkg, Output=test+'\n'))+'\n')
                return type('Result', (), {'returncode': 0})()
            class Process:
                stdout = [json.dumps(e)+'\n' for e in events]
                def wait(self):
                    return code
            with patch('sys.argv', ['run_tests.py', '--pattern-file', str(p/'pattern'),
                       '--log', str(p/'log'), '--result', str(p/'result'), *extra]), \
                 patch.object(run_tests.subprocess, 'run', side_effect=discover), \
                 patch.object(run_tests.subprocess, 'Popen', return_value=Process()) as launch, \
                 contextlib.redirect_stdout(io.StringIO()):
                result = run_tests.main()
            return result, json.loads((p/'result').read_text()), launch

    def test_all_packages_accounted(self):
        pairs = [(run_tests.PACKAGE, 'TestExample'), ('other/package', 'TestExample')]
        events = [dict(Package=p, Test=t, Action=a) for p,t in pairs for a in ('run','pass')]
        code, result, launch = self.run_driver(pairs, events, extra=('--package','.', '--package','./other'))
        self.assertEqual(code, 0)
        self.assertEqual(result['completed_tests'], 2)
        self.assertEqual(launch.call_args.args[0][-2:], ['.', './other'])

    def test_missing_execution_fails_even_with_zero_exit(self):
        self.assertEqual(self.run_driver([('p','TestExample')], [])[0], 1)

    def test_empty_selection_fails(self):
        self.assertEqual(self.run_driver([], [])[0], 1)

    def test_failure_exit_fails(self):
        events = [dict(Package='p', Test='TestExample', Action=a) for a in ('run','fail')]
        self.assertEqual(self.run_driver([('p','TestExample')], events, code=1)[0], 1)

    def test_required_prerequisite_skip_fails(self):
        events = [dict(Package='p', Test='TestExample', Action=a) for a in ('run','skip')]
        self.assertEqual(self.run_driver([('p','TestExample')], events,
                                        extra=('--require-no-skips',))[0], 1)

    def test_root_only_default(self):
        events = [dict(Package='p', Test='TestExample', Action=a) for a in ('run','pass')]
        code, _, launch = self.run_driver([('p','TestExample')], events)
        self.assertEqual(code, 0)
        self.assertEqual(launch.call_args.args[0][-1], '.')


if __name__ == '__main__':
    unittest.main()
