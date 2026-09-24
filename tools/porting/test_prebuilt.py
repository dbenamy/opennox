"""Mocked tests for optional prebuilt root-test support."""
import contextlib
import hashlib
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

import run_tests


METADATA = """path github.com/opennox/opennox/v1.test
mod github.com/opennox/opennox/v1 (devel)
build -tags=server,porttest
build CGO_ENABLED=1
build GOARCH=386
build GOOS=linux
build GO386=sse2
"""


class PrebuiltRunner(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        (self.root / "src").mkdir()
        (self.root / "src/a.s").write_bytes(b"supplemental-input")
        self.binary = self.root / "server.test"
        self.binary.write_bytes(b"test-binary")
        self.source = {"src/a.go": "frozen-sha"}

    def invoke(self, events=(), exit_code=0, *, tags="porttest,server", record_update=None,
               metadata=METADATA, extra=()):
        p = self.root
        pattern = p / "pattern"
        pattern.write_text("^TestExample$")
        record = {
            "binary_path": str(self.binary.resolve()),
            "sha256": hashlib.sha256(self.binary.read_bytes()).hexdigest(),
            "source": self.source,
            "supplemental_source": {"src/a.s": hashlib.sha256(b"supplemental-input").hexdigest()},
            "tags": tags,
            "package": run_tests.PACKAGE,
        }
        if record_update:
            record_update(record)
        rec_path = p / "binary.json"
        rec_path.write_text(json.dumps(record))
        discovered = [{"Action": "output", "Package": run_tests.PACKAGE,
                       "Output": "TestExample\n"}]
        calls = []

        def run(command, **kwargs):
            calls.append(command)
            if command[:3] == ["go", "version", "-m"]:
                return type("Result", (), {"returncode": 0, "stdout": metadata})()
            if command[:3] == ["git", "ls-files", "-z"]:
                return type("Result", (), {"returncode": 0, "stdout": b"src/a.go\0src/a.s\0"})()
            for event in discovered:
                kwargs["stdout"].write(json.dumps(event) + "\n")
            return type("Result", (), {"returncode": 0})()

        class Process:
            stdout = [json.dumps(event) + "\n" for event in events]
            def wait(self):
                return exit_code

        argv = ["run_tests.py", "--pattern-file", str(pattern), "--log", str(p / "run.jsonl"),
                "--result", str(p / "result.json"), "--tags", tags,
                "--test-binary", str(self.binary), "--binary-record", str(rec_path), *extra]
        with patch.object(run_tests, "ROOT", self.root), \
             patch.object(run_tests, "fingerprints", return_value=self.source), \
             patch.object(run_tests.subprocess, "run", side_effect=run), \
             patch.object(run_tests.subprocess, "Popen", return_value=Process()) as launch, \
             patch("sys.argv", argv), contextlib.redirect_stdout(io.StringIO()):
            code = run_tests.main()
        result = json.loads((p / "result.json").read_text())
        return code, result, calls, launch

    def test_successful_prebuilt_discovery_and_execution(self):
        events = [
            {"Action": "run", "Package": run_tests.PACKAGE, "Test": "TestExample"},
            {"Action": "pass", "Package": run_tests.PACKAGE, "Test": "TestExample"},
        ]
        code, result, calls, launch = self.invoke(events)
        self.assertEqual(code, 0)
        self.assertTrue(result["binary_verification"]["verified"])
        self.assertEqual(result["selected_root_tests"], 1)
        self.assertIn("-test.list=^TestExample$", calls[2])
        self.assertEqual(launch.call_args.args[0][1:5], ["tool", "test2json", "-t", "-p"])
        self.assertIn("-test.v=test2json", launch.call_args.args[0])

    def test_wrong_hash_is_rejected_before_execution(self):
        code, result, calls, launch = self.invoke(record_update=lambda r: r.update(sha256="wrong"))
        self.assertEqual(code, 1)
        self.assertIn("hash mismatch", result["reason"])
        self.assertEqual(calls, [])
        launch.assert_not_called()

    def test_wrong_record_path_is_rejected(self):
        code, result, calls, launch = self.invoke(
            record_update=lambda r: r.update(binary_path=str(self.root / "other.test")))
        self.assertEqual(code, 1)
        self.assertIn("path mismatch", result["reason"])
        self.assertEqual(calls, [])
        launch.assert_not_called()

    def test_stale_source_is_rejected(self):
        code, result, calls, _ = self.invoke(record_update=lambda r: r.update(source={"src/old.go": "x"}))
        self.assertEqual(code, 1)
        self.assertIn("source fingerprint mismatch", result["reason"])
        self.assertEqual(calls, [])

    def test_stale_supplemental_tracked_input_is_rejected(self):
        code, result, calls, _ = self.invoke(record_update=lambda r: r.update(supplemental_source={"src/a.s": "stale"}))
        self.assertEqual(code, 1)
        self.assertIn("supplemental input fingerprint mismatch", result["reason"])
        self.assertEqual(calls[-1][:3], ["git", "ls-files", "-z"])

    def test_wrong_tags_rejected(self):
        code, result, calls, _ = self.invoke(record_update=lambda r: r.update(tags="porttest"))
        self.assertEqual(code, 1)
        self.assertIn("provenance tags mismatch", result["reason"])
        self.assertEqual(calls, [])

    def test_wrong_build_metadata_tags_rejected(self):
        code, result, _, launch = self.invoke(metadata=METADATA.replace("server,porttest", "porttest"))
        self.assertEqual(code, 1)
        self.assertIn("build tags mismatch", result["reason"])
        launch.assert_not_called()

    def test_wrong_target_metadata_rejected(self):
        code, result, _, launch = self.invoke(metadata=METADATA.replace("GOARCH=386", "GOARCH=amd64"))
        self.assertEqual(code, 1)
        self.assertIn("target metadata mismatch", result["reason"])
        launch.assert_not_called()

    def test_prebuilt_rejects_non_root_package(self):
        with patch("sys.argv", ["run_tests.py", "--pattern-file", "x", "--log", "l", "--result", "r",
                                "--test-binary", "b", "--binary-record", "r", "--package", "./other"]):
            with self.assertRaises(SystemExit) as ex:
                run_tests.main()
        self.assertEqual(ex.exception.code, 2)

    def test_missing_execution_fails_with_zero_exit(self):
        code, result, _, _ = self.invoke(events=())
        self.assertEqual(code, 1)
        self.assertEqual(result["started_tests"], 0)
        self.assertIn("did not all execute", result["reason"])

    def test_fail_event_fails_even_if_process_exits_zero(self):
        events = [
            {"Action": "run", "Package": run_tests.PACKAGE, "Test": "TestExample"},
            {"Action": "fail", "Package": run_tests.PACKAGE, "Test": "TestExample"},
        ]
        code, result, _, _ = self.invoke(events, exit_code=0)
        self.assertEqual(code, 1)
        self.assertEqual(len(result["failed_tests"]), 1)
        self.assertIn("suite failed", result["reason"])

    def test_unselected_subtest_fail_event_fails_zero_exit(self):
        events = [
            {"Action": "run", "Package": run_tests.PACKAGE, "Test": "TestExample"},
            {"Action": "pass", "Package": run_tests.PACKAGE, "Test": "TestExample"},
            {"Action": "fail", "Package": run_tests.PACKAGE, "Test": "TestExample/subcase"},
        ]
        code, result, _, _ = self.invoke(events, exit_code=0)
        self.assertEqual(code, 1)
        self.assertIn((run_tests.PACKAGE, "TestExample/subcase"),
                      [tuple(item) for item in result["failure_events"]])

    def test_package_level_fail_event_fails_zero_exit(self):
        events = [
            {"Action": "run", "Package": run_tests.PACKAGE, "Test": "TestExample"},
            {"Action": "pass", "Package": run_tests.PACKAGE, "Test": "TestExample"},
            {"Action": "fail", "Package": run_tests.PACKAGE},
        ]
        code, result, _, _ = self.invoke(events, exit_code=0)
        self.assertEqual(code, 1)
        self.assertIn((run_tests.PACKAGE, ""), [tuple(item) for item in result["failure_events"]])


if __name__ == "__main__":
    unittest.main()
