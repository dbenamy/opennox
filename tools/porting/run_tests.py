#!/usr/bin/env python3
"""Run a selected port suite and require every discovered root test to execute.

Source build/baseline/env.sh first. Raw Go output stays in the supplied log files.
"""
import argparse
import json
import os
from pathlib import Path
import subprocess
import sys
import time

ROOT = Path(__file__).resolve().parents[2]
PACKAGE = "github.com/opennox/opennox/v1"


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--pattern-file", type=Path, required=True)
    parser.add_argument("--tags", default="porttest")
    parser.add_argument("--log", type=Path, required=True)
    parser.add_argument("--result", type=Path, required=True)
    args = parser.parse_args()
    pattern = args.pattern_file.read_text().strip()
    result = {"tags": args.tags, "success": False}
    started = time.monotonic()

    def finish(reason, code=1):
        result.update(reason=reason, wall_seconds=time.monotonic() - started)
        args.result.parent.mkdir(parents=True, exist_ok=True)
        args.result.write_text(json.dumps(result, indent=2) + "\n")
        print(json.dumps(result), flush=True)
        return code

    if not pattern or "\n" in pattern or "\r" in pattern:
        return finish("test pattern must be one nonempty line")
    args.log.parent.mkdir(parents=True, exist_ok=True)
    common = ["go", "test", "-p", "2", "-tags", args.tags]
    env = dict(os.environ, GOMAXPROCS=os.environ.get("GOMAXPROCS", "2"))
    discovery = args.log.with_suffix(args.log.suffix + ".discovery")
    with discovery.open("w") as log:
        proc = subprocess.run(common + ["-list", pattern, "."], cwd=ROOT / "src",
                              env=env, stdout=log, stderr=subprocess.STDOUT)
    result["discovery_exit"] = proc.returncode
    if proc.returncode:
        return finish("test discovery failed")
    expected = {line for line in discovery.read_text().splitlines()
                if line.startswith("Test") and not any(c.isspace() for c in line)}
    result["selected_root_tests"] = len(expected)
    if not expected:
        return finish("pattern selected no root tests")

    ran, completed = set(), set()
    with args.log.open("w") as log:
        proc = subprocess.Popen(common + ["-count=1", "-timeout", "600s", "-json",
                                           "-run", pattern, "./..."],
                                cwd=ROOT / "src", env=env, stdout=subprocess.PIPE,
                                stderr=subprocess.STDOUT, text=True, errors="replace")
        for line in proc.stdout:
            log.write(line)
            try:
                event = json.loads(line)
            except json.JSONDecodeError:
                continue
            name = event.get("Test", "")
            if event.get("Package") != PACKAGE or name not in expected:
                continue
            if event.get("Action") == "run":
                ran.add(name)
            elif event.get("Action") in ("pass", "skip", "fail"):
                completed.add(name)
        code = proc.wait()
    result.update(exit=code, started_root_tests=len(ran), completed_root_tests=len(completed),
                  missing_started=sorted(expected - ran), missing_completed=sorted(expected - completed))
    if code:
        return finish("selected test suite failed")
    if expected - ran or expected - completed:
        return finish("discovered root tests did not all execute and finish")
    result["success"] = True
    return finish("all selected root tests executed and finished", 0)


if __name__ == "__main__":
    sys.exit(main())
