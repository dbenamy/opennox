"""Mocked orchestration tests for the prebuilt profile controller."""
import contextlib
import io
import json
import os
from pathlib import Path
import sys
import tempfile
import threading
import time
import unittest
from unittest.mock import patch

import run_profiles


class ProfileController(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        (self.root / "src").mkdir()
        (self.root / "tools/porting").mkdir(parents=True)
        self.pattern = self.root / "patterns.txt"
        self.pattern.write_text("^TestExample$")
        self.source = {"src/a.go": "source-hash"}

    def invoke(self, *, jobs=1, runner=None, source_values=None):
        out = self.root / "out"
        builds, runners = [], []
        source_values = source_values or [self.source]
        source_iter = iter(source_values)

        def fingerprint():
            try:
                return next(source_iter)
            except StopIteration:
                return source_values[-1]

        def default_runner(profile, command):
            returncode, success = 0, True
            result_path = Path(command[command.index("--result") + 1])
            result_path.write_text(json.dumps({"success": success}))
            return type("Result", (), {"returncode": returncode})()

        runner = runner or default_runner

        def fake_run(command, **kwargs):
            if command[0] == "go":
                profile = Path(command[command.index("-o") + 1]).stem.removesuffix(".test")
                builds.append(profile)
                Path(command[command.index("-o") + 1]).write_bytes((profile + " binary").encode())
                return type("Result", (), {"returncode": 0})()
            profile = next(k for k, tags in run_profiles.PROFILES.items() if tags == command[command.index("--tags") + 1])
            runners.append(profile)
            if len(builds) != len(run_profiles.PROFILES):
                raise AssertionError("profile execution began before all builds completed")
            return runner(profile, command)

        argv = ["run_profiles.py", "--out", str(out), "--pattern-file", str(self.pattern),
                "--jobs", str(jobs)]
        with patch.object(run_profiles, "ROOT", self.root), \
             patch.object(run_profiles, "fingerprints", side_effect=fingerprint), \
             patch.object(run_profiles, "supplemental_fingerprints", return_value={"src/a.s": "input-hash"}), \
             patch.object(run_profiles.subprocess, "run", side_effect=fake_run), \
             patch("sys.argv", argv), contextlib.redirect_stdout(io.StringIO()):
            code = run_profiles.main()
        report = json.loads((out / "result.json").read_text())
        return code, report, builds, runners

    def test_all_profiles_build_sequentially_before_any_run(self):
        code, result, builds, runners = self.invoke(jobs=1)
        self.assertEqual(code, 0)
        self.assertEqual(builds, ["server", "highres", "default"])
        self.assertEqual(set(runners), set(run_profiles.PROFILES))

    def test_shared_diagnostic_output_env_is_rejected(self):
        out = self.root / "out"
        argv = ["run_profiles.py", "--out", str(out), "--pattern-file", str(self.pattern)]
        with patch("sys.argv", argv), \
             patch.dict(os.environ, {"OPENNOX_CLIENT_EFFECTS_DIAGNOSTICS": "/shared/captures"}), \
             contextlib.redirect_stderr(io.StringIO()):
            with self.assertRaises(SystemExit) as ex:
                run_profiles.main()
        self.assertEqual(ex.exception.code, 2)
        self.assertFalse(out.exists())

    def test_execution_concurrency_is_capped_at_two(self):
        lock = threading.Lock()
        state = {"active": 0, "maximum": 0}

        def runner(profile, command):
            result_path = Path(command[command.index("--result") + 1])
            with lock:
                state["active"] += 1
                state["maximum"] = max(state["maximum"], state["active"])
            time.sleep(0.04)
            result_path.write_text(json.dumps({"success": True}))
            with lock:
                state["active"] -= 1
            return type("Result", (), {"returncode": 0})()

        code, _, _, _ = self.invoke(jobs=2, runner=runner)
        self.assertEqual(code, 0)
        self.assertEqual(state["maximum"], 2)

    def test_source_change_after_build_refuses_to_start_profiles(self):
        changed = {"src/a.go": "changed"}
        code, result, builds, runners = self.invoke(source_values=[self.source, changed, changed])
        self.assertEqual(code, 1)
        self.assertEqual(builds, ["server"])
        self.assertEqual(runners, [])
        self.assertIn("source changed", result["error"])

    def test_runner_failure_with_zero_exit_marks_batch_failed(self):
        def runner(profile, command):
            result_path = Path(command[command.index("--result") + 1])
            result_path.write_text(json.dumps({"success": profile != "server"}))
            return type("Result", (), {"returncode": 0})()

        code, result, _, runners = self.invoke(jobs=2, runner=runner)
        self.assertEqual(code, 1)
        self.assertEqual(set(runners), set(run_profiles.PROFILES))
        self.assertFalse(next(row for row in result["profiles"] if row["profile"] == "server")["success"])

    def test_controller_joins_other_runs_after_runner_exception(self):
        completed = set()
        lock = threading.Lock()

        def runner(profile, command):
            if profile == "server":
                raise RuntimeError("mocked runner failure")
            time.sleep(0.03)
            result_path = Path(command[command.index("--result") + 1])
            result_path.write_text(json.dumps({"success": True}))
            with lock:
                completed.add(profile)
            return type("Result", (), {"returncode": 0})()

        code, result, _, runners = self.invoke(jobs=2, runner=runner)
        self.assertEqual(code, 1)
        self.assertEqual(set(runners), set(run_profiles.PROFILES))
        self.assertEqual(completed, {"highres", "default"})
        self.assertEqual(len(result["profiles"]), 3)


if __name__ == "__main__":
    unittest.main()
