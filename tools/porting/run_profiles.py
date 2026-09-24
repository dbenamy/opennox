#!/usr/bin/env python3
"""Build root port-test profiles sequentially, then run at most two at once.

Source the port environment first. Keep source and dependencies frozen throughout.
Only root-package tests are supported; production scenarios remain sequential.
"""
import argparse
from concurrent.futures import ThreadPoolExecutor
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import time

from run_batch import fingerprints
from run_tests import PACKAGE, supplemental_fingerprints

ROOT = Path(__file__).resolve().parents[2]
PROFILES = {"server": "porttest,server", "highres": "porttest,highres", "default": "porttest"}


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--out", type=Path, required=True)
    ap.add_argument("--pattern-file", type=Path, required=True)
    ap.add_argument("--jobs", type=int, choices=(1, 2), default=1)
    ap.add_argument("--timeout-seconds", type=int, default=3600)
    args = ap.parse_args()
    if args.timeout_seconds <= 0:
        ap.error("timeout must be positive")
    capture_outputs = sorted(k for k, v in os.environ.items()
                             if v and k.startswith("OPENNOX_")
                             and ("CAPTURE" in k or "DIAGNOSTIC" in k))
    if capture_outputs:
        ap.error("unset capture/diagnostic output variables before this sweep: " + ", ".join(capture_outputs))
    out = args.out.resolve()
    out.mkdir(parents=True, exist_ok=False)
    pattern = args.pattern_file.resolve()
    env = dict(os.environ, GOMAXPROCS=os.environ.get("GOMAXPROCS", "2"),
               GOMEMLIMIT=os.environ.get("GOMEMLIMIT", "768MiB"))
    build_env = dict(env, GOGC="100", GOMEMLIMIT="1536MiB")
    source = fingerprints()
    supplemental = supplemental_fingerprints()
    report = dict(success=False, jobs=args.jobs, builds=[], profiles=[], source=source,
                  supplemental_source=supplemental)
    start = time.monotonic()

    def unchanged():
        return fingerprints() == source and supplemental_fingerprints() == supplemental

    def execute(profile):
        binary = out / (profile + ".test")
        command = [sys.executable, str(ROOT / "tools/porting/run_tests.py"),
                   "--pattern-file", str(pattern), "--tags", PROFILES[profile],
                   "--test-binary", str(binary), "--binary-record", str(out / (profile + "-binary.json")),
                   "--log", str(out / (profile + ".jsonl")),
                   "--result", str(out / (profile + "-result.json")),
                   "--timeout-seconds", str(args.timeout_seconds)]
        now = time.monotonic()
        with (out / (profile + "-runner.log")).open("w") as log:
            proc = subprocess.run(command, cwd=ROOT, env=env, stdout=log, stderr=subprocess.STDOUT)
        result = json.loads((out / (profile + "-result.json")).read_text())
        return dict(profile=profile, exit=proc.returncode, success=proc.returncode == 0 and result["success"],
                    wall_seconds=time.monotonic() - now, command=command)

    try:
        for profile, tags in PROFILES.items():
            binary = out / (profile + ".test")
            command = ["go", "test", "-p", "2", "-c", "-tags", tags, "-o", str(binary), "."]
            now = time.monotonic()
            with (out / (profile + "-build.log")).open("w") as log:
                proc = subprocess.run(command, cwd=ROOT / "src", env=build_env,
                                      stdout=log, stderr=subprocess.STDOUT, timeout=1800)
            report["builds"].append(dict(profile=profile, exit=proc.returncode, command=command,
                                         wall_seconds=time.monotonic() - now))
            if proc.returncode or not unchanged():
                raise RuntimeError(f"{profile}: build failed or source changed")
            record = dict(binary_path=str(binary), sha256=hashlib.sha256(binary.read_bytes()).hexdigest(),
                          source=source, supplemental_source=supplemental, tags=tags, package=PACKAGE)
            (out / (profile + "-binary.json")).write_text(json.dumps(record, indent=2) + "\n")
        # The executor joins every submitted process, including after a failure.
        with ThreadPoolExecutor(max_workers=args.jobs) as pool:
            futures = [(profile, pool.submit(execute, profile)) for profile in PROFILES]
            for profile, future in futures:
                try:
                    row = future.result()
                except Exception as exc:
                    row = dict(profile=profile, success=False, error=str(exc))
                report["profiles"].append(row)
                print(json.dumps(row), flush=True)
        report["success"] = all(row["success"] for row in report["profiles"])
    except Exception as exc:
        report["error"] = str(exc)
    finally:
        report["source_unchanged"] = unchanged()
        report["success"] &= report["source_unchanged"]
        report["wall_seconds"] = time.monotonic() - start
        (out / "result.json").write_text(json.dumps(report, indent=2) + "\n")
    return 0 if report["success"] else 1


if __name__ == "__main__":
    sys.exit(main())
