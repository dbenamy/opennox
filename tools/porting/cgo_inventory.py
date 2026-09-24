#!/usr/bin/env python3
"""Inventory production cgo dependencies without compiling or changing modules.

Source build/baseline/env.sh first. Raw go-list output stays in --out; the compact
summary is suitable for a reviewed, committed snapshot. go list -e is metadata,
not a successful build, even when Error/DepsErrors are empty.
"""
import argparse
import collections
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[2]
MODULE = "github.com/opennox/opennox/v1"
PROFILES = {"default": "", "highres": "highres", "server": "server"}


def decode_stream(raw):
    decoder = json.JSONDecoder()
    pos = 0
    while pos < len(raw):
        while pos < len(raw) and raw[pos].isspace():
            pos += 1
        if pos == len(raw):
            break
        value, pos = decoder.raw_decode(raw, pos)
        yield value


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--out", type=Path, required=True)
    args = parser.parse_args()
    out = args.out.resolve()
    out.mkdir(parents=True, exist_ok=True)
    env = dict(os.environ, GOOS="linux", GOARCH="386", GO386="sse2",
               GOPROXY="off", GOTOOLCHAIN="local", GOMAXPROCS="2")
    result = {
        "revision": subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip(),
        "toolchain": subprocess.check_output(["go", "version"], env=env, text=True).strip(),
        "target": "linux/386/SSE2", "entrypoint": "./cmd/opennox",
        "scope": "Production dependency metadata; no porttest/safe; not compilation or qualification",
        "selector_metric": "Lexical C.selector mentions (including comments), not resolved calls",
        "profiles": {},
    }
    for cgo in ("1", "0"):
        for profile, tag in PROFILES.items():
            cmd = ["go", "list", "-mod=readonly", "-e", "-deps", "-json"]
            if tag:
                cmd += ["-tags", tag]
            cmd += ["./cmd/opennox"]
            run = subprocess.run(cmd, cwd=ROOT / "src", env=dict(env, CGO_ENABLED=cgo),
                                 capture_output=True, text=True, timeout=45)
            stem = out / f"{profile}-cgo{cgo}"
            stem.with_suffix(".jsonstream").write_text(run.stdout)
            stem.with_suffix(".stderr").write_text(run.stderr)
            if run.returncode:
                raise RuntimeError(f"{cmd} exited {run.returncode}; inspect {stem}.stderr")
            packages = list(decode_stream(run.stdout))
            if not any(p.get("ImportPath") == MODULE + "/cmd/opennox" for p in packages):
                raise RuntimeError("missing production entrypoint")
            by_name = {p["ImportPath"]: p for p in packages}
            direct = []
            for p in packages:
                if not p.get("CgoFiles"):
                    continue
                item = {"package": p["ImportPath"], "standard": p.get("Standard", False),
                        "files": sorted(p["CgoFiles"]),
                        "imported_by": sorted(q["ImportPath"] for q in packages
                                              if p["ImportPath"] in q.get("Imports", []))}
                if p["ImportPath"].startswith(MODULE + "/"):
                    selectors = collections.Counter()
                    exports, empty, hashes = [], [], {}
                    for name in item["files"]:
                        path = Path(p["Dir"]) / name
                        source = path.read_text()
                        hashes[str(path.relative_to(ROOT))] = hashlib.sha256(path.read_bytes()).hexdigest()
                        names = re.findall(r"\bC\.([A-Za-z_]\w*)", source)
                        exported = re.findall(r"^//export (\w+)", source, re.M)
                        selectors.update(names)
                        exports.extend(exported)
                        if not names and not exported and "#cgo" not in source:
                            empty.append(name)
                    item.update(source_sha256=hashes, export_directives=len(exports),
                                c_selector_mentions=dict(sorted(selectors.items())),
                                apparently_unused_import_candidates=empty)
                direct.append(item)
            errors = []
            for p in packages:
                for error in ([p["Error"]] if p.get("Error") else []) + p.get("DepsErrors", []):
                    entry = {"package": p["ImportPath"], "error": error["Err"],
                             "import_stack": error.get("ImportStack", [])}
                    entry["error"] = entry["error"].replace(str(ROOT), "<repo>")
                    if entry not in errors:
                        errors.append(entry)
            # Transitive project packages that need any selected cgo package.
            memo = {}
            def needs_cgo(name, seen):
                if name in memo:
                    return memo[name]
                if name in seen:
                    raise RuntimeError(f"dependency cycle: {name}")
                p = by_name.get(name, {})
                value = bool(p.get("CgoFiles")) or any(
                    needs_cgo(dep, seen | {name}) for dep in p.get("Imports", []) if dep != "C")
                memo[name] = value
                return value
            transitively = sorted(p["ImportPath"] for p in packages
                                  if p["ImportPath"].startswith(MODULE + "/")
                                  and needs_cgo(p["ImportPath"], set()))
            key = f"{profile}-cgo{cgo}"
            result["profiles"][key] = {
                "command": cmd, "go_list_exit": run.returncode,
                "package_count": len(packages), "direct_cgo": direct,
                "project_packages_reaching_cgo": transitively, "metadata_errors": errors,
            }
            print(key, "packages", len(packages), "direct cgo", len(direct),
                  "metadata errors", len(errors), flush=True)
    (out / "inventory.json").write_text(json.dumps(result, indent=2) + "\n")


if __name__ == "__main__":
    main()
