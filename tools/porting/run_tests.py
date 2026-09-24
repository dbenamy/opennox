#!/usr/bin/env python3
"""Run a selected port suite and require every discovered test to execute.

Source build/baseline/env.sh first. Raw Go output stays in the supplied log files.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import time

from run_batch import fingerprints

ROOT = Path(__file__).resolve().parents[2]
PACKAGE = "github.com/opennox/opennox/v1"


def normalize_tags(value):
    return ",".join(sorted(part.strip() for part in value.split(",") if part.strip()))


def parse_build_metadata(text):
    result = {}
    for line in text.splitlines():
        parts = line.split()
        if not parts:
            continue
        if parts[0] == "path" and len(parts) >= 2:
            result["path"] = parts[1]
        elif parts[0] == "build" and len(parts) >= 2 and "=" in parts[1]:
            key, value = parts[1].split("=", 1)
            result["tags" if key == "-tags" else key] = value
    return result


def supplemental_fingerprints(source=None):
    """Hash tracked src inputs omitted by run_batch.fingerprints()."""
    if source is None:
        source = fingerprints()
    proc = subprocess.run(["git", "ls-files", "-z", "src"], cwd=ROOT,
                          capture_output=True)
    if proc.returncode:
        raise RuntimeError("git ls-files failed while checking supplemental inputs")
    out = {}
    for raw in proc.stdout.split(b"\0"):
        if not raw:
            continue
        rel = raw.decode("utf-8")
        if rel in source:
            continue
        path = ROOT / rel
        if path.is_file():
            out[rel] = hashlib.sha256(path.read_bytes()).hexdigest()
    return out


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--pattern-file", type=Path, required=True)
    parser.add_argument("--tags", default="porttest")
    parser.add_argument("--log", type=Path, required=True)
    parser.add_argument("--result", type=Path, required=True)
    parser.add_argument("--timeout-seconds", type=int, default=600,
                        help="per-package Go test timeout (default: 600)")
    parser.add_argument("--build-memory-limit", default="1536MiB",
                        help="GOMEMLIMIT for discovery/build (default: 1536MiB)")
    parser.add_argument("--package", action="append", dest="packages",
                        help="affected package (repeatable; default: .)")
    parser.add_argument("--require-no-skips", action="store_true")
    parser.add_argument("--test-binary", type=Path,
                        help="run a previously built test binary instead of go test")
    parser.add_argument("--binary-record", type=Path,
                        help="provenance record required with --test-binary")
    args = parser.parse_args()
    packages = args.packages or ["."]
    if bool(args.test_binary) != bool(args.binary_record):
        parser.error("--test-binary and --binary-record must be supplied together")
    if args.test_binary and packages != ["."]:
        parser.error("prebuilt test binaries support only --package .")
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
    prebuilt = bool(args.test_binary)
    common = ["go", "test", "-p", "2", "-tags", args.tags]
    # A long 386 process can exhaust its 4 GiB address space before the default
    # heap-growth target triggers collection. Leave room for C and mapped blobs.
    env = dict(os.environ, GOMAXPROCS=os.environ.get("GOMAXPROCS", "2"),
               GOMEMLIMIT=os.environ.get("GOMEMLIMIT", "768MiB"))
    result["runtime_env"] = {k: env.get(k) for k in ("GOMAXPROCS", "GOMEMLIMIT", "GOGC")}
    # Discovery compiles before listing tests. Give the host compiler more room;
    # actual test execution below keeps the 386 runtime budget unchanged.
    discovery_env = dict(env, GOMEMLIMIT=args.build_memory_limit)
    result["discovery_env"] = {k: discovery_env.get(k) for k in ("GOMAXPROCS", "GOMEMLIMIT", "GOGC")}
    discovery = args.log.with_suffix(args.log.suffix + ".discovery")

    if prebuilt:
        binary = args.test_binary.resolve()
        record_path = args.binary_record.resolve()
        result["binary_verification"] = {"path": str(binary), "record_path": str(record_path),
                                         "verified": False}
        try:
            record = json.loads(record_path.read_text())
            actual_sha = hashlib.sha256(binary.read_bytes()).hexdigest()
            current_sources = fingerprints()
            result["binary_verification"].update(actual_sha256=actual_sha,
                                                  recorded_sha256=record.get("sha256"))
            if record.get("binary_path") != str(binary):
                return finish("binary provenance path mismatch")
            if record.get("sha256") != actual_sha:
                return finish("binary hash mismatch")
            result["binary_verification"]["hash_matches"] = True
            if record.get("source") != current_sources:
                return finish("binary source fingerprint mismatch")
            result["binary_verification"]["source_matches"] = True
            if normalize_tags(record.get("tags", "")) != normalize_tags(args.tags):
                return finish("binary provenance tags mismatch")
            if record.get("package") != PACKAGE:
                return finish("binary provenance package mismatch")
            current_inputs = supplemental_fingerprints(current_sources)
            if record.get("supplemental_source") != current_inputs:
                return finish("binary supplemental input fingerprint mismatch")
            result["binary_verification"]["inputs_match"] = True
            meta = subprocess.run(["go", "version", "-m", str(binary)], cwd=ROOT / "src",
                                  env=env, text=True, capture_output=True)
            if meta.returncode:
                return finish("binary build metadata inspection failed")
            metadata = parse_build_metadata(meta.stdout)
            required = {"GOOS": "linux", "GOARCH": "386", "GO386": "sse2", "CGO_ENABLED": "1"}
            if metadata.get("path") != PACKAGE + ".test":
                return finish("binary is not the expected root-package test binary")
            if any(metadata.get(k) != v for k, v in required.items()):
                return finish("binary target metadata mismatch")
            if normalize_tags(metadata.get("tags", "")) != normalize_tags(args.tags):
                return finish("binary build tags mismatch")
            result["binary_verification"] = {
                "path": str(binary), "sha256": actual_sha, "source_matches": True,
                "record_matches": True, "build_metadata": metadata, "verified": True,
            }
        except (OSError, ValueError, KeyError, TypeError, AttributeError, RuntimeError) as exc:
            return finish(f"binary provenance verification failed: {exc}")

        discover_cmd = ["go", "tool", "test2json", "-t", "-p", PACKAGE, str(binary),
                        "-test.list=" + pattern]
        with discovery.open("w") as log:
            proc = subprocess.run(discover_cmd, cwd=ROOT / "src", env=env,
                                  stdout=log, stderr=subprocess.STDOUT)
        result["discovery_command"] = discover_cmd
    else:
        with discovery.open("w") as log:
            proc = subprocess.run(common + ["-json", "-list", pattern] + packages, cwd=ROOT / "src",
                                  env=discovery_env, stdout=log, stderr=subprocess.STDOUT)
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

    ran, completed, skipped, failed, failure_events = set(), set(), set(), set(), set()
    with args.log.open("w") as log:
        command = (["go", "tool", "test2json", "-t", "-p", PACKAGE, str(binary),
                    "-test.v=test2json", "-test.count=1", "-test.timeout=" + f"{args.timeout_seconds}s",
                    "-test.run=" + pattern] if prebuilt else
                   common + ["-count=1", "-timeout", f"{args.timeout_seconds}s", "-json",
                             "-run", pattern] + packages)
        result["execution_command"] = command
        proc = subprocess.Popen(command,
                                cwd=ROOT / "src", env=env, stdout=subprocess.PIPE,
                                stderr=subprocess.STDOUT, text=True, errors="replace")
        for line in proc.stdout:
            log.write(line)
            try:
                event = json.loads(line)
            except json.JSONDecodeError:
                continue
            name = (event.get("Package"), event.get("Test", ""))
            if event.get("Action") == "fail":
                failure_events.add(name)
            if name not in expected:
                continue
            if event.get("Action") == "run":
                ran.add(name)
            elif event.get("Action") in ("pass", "skip", "fail"):
                completed.add(name)
                if event.get("Action") == "skip":
                    skipped.add(name)
                elif event.get("Action") == "fail":
                    failed.add(name)
        code = proc.wait()
    result.update(exit=code, started_tests=len(ran), completed_tests=len(completed),
                  started_root_tests=sum(pkg == PACKAGE for pkg, _ in ran),
                  completed_root_tests=sum(pkg == PACKAGE for pkg, _ in completed),
                  skipped_tests=sorted(skipped),
                  failed_tests=sorted(failed),
                  failure_events=sorted(failure_events),
                  missing_started=sorted(expected - ran), missing_completed=sorted(expected - completed))
    if code or failure_events:
        return finish("selected test suite failed")
    if expected - ran or expected - completed:
        return finish("discovered tests did not all execute and finish")
    if args.require_no_skips and skipped:
        return finish("selected tests skipped despite required prerequisites")
    result["success"] = True
    return finish("all selected tests executed and finished", 0)


if __name__ == "__main__":
    sys.exit(main())
