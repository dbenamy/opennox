#!/usr/bin/env python3
"""Run a selected port suite and require every discovered test to execute.

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
    parser.add_argument("--timeout-seconds", type=int, default=600,
                        help="per-package Go test timeout (default: 600)")
    parser.add_argument("--package", action="append", dest="packages",
                        help="affected package (repeatable; default: .)")
    parser.add_argument("--require-no-skips", action="store_true")
    args = parser.parse_args()
    packages = args.packages or ["."]
    if args.timeout_seconds <= 0:
        parser.error("--timeout-seconds must be positive")
    pattern = args.pattern_file.read_text().strip()
    result = {"tags": args.tags, "success": False, "timeout_seconds": args.timeout_seconds, "packages": packages}
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
    # A long 386 process can exhaust its 4 GiB address space before the default
    # heap-growth target triggers collection. Leave room for C and mapped blobs.
    env = dict(os.environ, GOMAXPROCS=os.environ.get("GOMAXPROCS", "2"),
               GOMEMLIMIT=os.environ.get("GOMEMLIMIT", "768MiB"))
    result["runtime_env"] = {k: env.get(k) for k in ("GOMAXPROCS", "GOMEMLIMIT", "GOGC")}
    discovery = args.log.with_suffix(args.log.suffix + ".discovery")
    with discovery.open("w") as log:
        proc = subprocess.run(common + ["-json", "-list", pattern] + packages, cwd=ROOT / "src",
                              env=env, stdout=log, stderr=subprocess.STDOUT)
    result["discovery_exit"] = proc.returncode
    if proc.returncode:
        return finish("test discovery failed")
    expected = set()
    for line in discovery.read_text().splitlines():
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            continue
        name = event.get("Output", "").strip()
        if name.startswith("Test") and not any(c.isspace() for c in name):
            expected.add((event["Package"], name))
    result["selected_tests"] = len(expected)
    result["selected_root_tests"] = sum(pkg == PACKAGE for pkg, _ in expected)
    if not expected:
        return finish("pattern selected no tests")

    ran, completed, skipped = set(), set(), set()
    with args.log.open("w") as log:
        proc = subprocess.Popen(common + ["-count=1", "-timeout", f"{args.timeout_seconds}s", "-json",
                                           "-run", pattern] + packages,
                                cwd=ROOT / "src", env=env, stdout=subprocess.PIPE,
                                stderr=subprocess.STDOUT, text=True, errors="replace")
        for line in proc.stdout:
            log.write(line)
            try:
                event = json.loads(line)
            except json.JSONDecodeError:
                continue
            name = (event.get("Package"), event.get("Test", ""))
            if name not in expected:
                continue
            if event.get("Action") == "run":
                ran.add(name)
            elif event.get("Action") in ("pass", "skip", "fail"):
                completed.add(name)
                if event.get("Action") == "skip":
                    skipped.add(name)
        code = proc.wait()
    result.update(exit=code, started_tests=len(ran), completed_tests=len(completed),
                  started_root_tests=sum(pkg == PACKAGE for pkg, _ in ran),
                  completed_root_tests=sum(pkg == PACKAGE for pkg, _ in completed),
                  skipped_tests=sorted(skipped),
                  missing_started=sorted(expected - ran), missing_completed=sorted(expected - completed))
    if code:
        return finish("selected test suite failed")
    if expected - ran or expected - completed:
        return finish("discovered tests did not all execute and finish")
    if args.require_no_skips and skipped:
        return finish("selected tests skipped despite required prerequisites")
    result["success"] = True
    return finish("all selected tests executed and finished", 0)


if __name__ == "__main__":
    sys.exit(main())
