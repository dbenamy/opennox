# Spell scalar boundaries

Status: converted and qualified; baseline `b6049004`.

Remove the unused specialized `Nox_spells_call_intint6_go` Go wrapper and its
embedded C adapter. Convert three live spell predicates from `C.bool` to `bool`
and the flags helper from `C.uint` to `uint32`. Their server calls, argument
conversions and all caller behavior remain unchanged. The qualified target is
Linux 386; no new portability claim is made.

Whole tracked-code search (including function-value references) found only the
adapter definition, wrapper definition and wrapper call in `src/legacy/spells.go`:
`git grep -n -I -e Nox_spells_call_intint6_go -e nox_spells_call_intint6_go -- . ':!docs' ':!PORTING_STATE.md'`.
The scalar helpers remain live in book rewards/rendering, quickbar input/actions/
rendering and objectives. They delegate directly to established Go owners.
AST review finds four changed functions, one removed and 101 unchanged.

Baseline reuses qualified `32df9553` results under exact source, supplemental
source, binary SHA256, runtime-environment and discovered-name checks. No fixtures
or expectations change. Package-aware inverse references plus whole assertion
files conservatively select nearly the entire corpus; book/quickbar/objective
families are explicitly included for indirect GUI dispatch. Selected compiled
roots: 2,482 default/highres and 2,471 server, all present in the prior accepted
pass sets. The opt-in map-population diagnostic is excluded; profile-specific
uncompiled roots are recorded individually in the baseline JSON. The safe-only
memory bridge test is not part of these three profiles.

Qualification passed the matching three-profile selection, 94-root focused
consumer preflight, safe build/static check, three production builds/ABI, exact
known-suite outcomes, fresh headless save/load and unchanged original assets.
Verified the specialized adapter symbol substring is absent from all production
binaries; exact-name checks alone cannot match cgo's generated hash prefix.
No C implementation is retained solely as a test reference.

Qualified counts: 5 client / 6 server production cgo files,
19 embedded C dispatch bodies, zero legacy exports, zero standalone production
or test C lines. Headers remain 157 files / 2,731 physical lines.

Primary performed this bounded batch locally: Luna remains quota-unavailable.
Ignored draft/review/reuse evidence lives in `build/port-spell-scalar-boundaries/`.
The committed baseline and batch manifest provide reproducible scope; previous
production evidence supplies only the original baseline, not converted acceptance.

All exact selected names passed without skips; no expectations changed. The
94-root native consumer preflight supplies the early behavior gate for this
return-type-only live change; the fresh production scenario supplies save/load
integration. No additional preview build/scenario was needed. Original C type
generation also confirms `_Bool` uses Go bool and C unsigned int uses uint32
on this target. External native-library selections remain unchanged.

Removed seven superseded placement binaries after source/replacement/hash and
host-use checks (387,072,000 allocated bytes); source maps and logs remain.
The final scenario deduplicated 1654 verified original-asset copies
(559968256 allocated bytes); saves/results and recovery manifest remain.
Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py
build/baseline/runs/spell-scalar-boundaries-save/deduplicated-assets.json`.
All batch scripts are consumed after completion.

Six older original-baseline executables from `9dcf1b9a` and `7034a7e4` were also
removed after committed-source, replacement/hash and host-use checks
(400,736,256 allocated bytes). Rebuild those revisions with retained binary
records; current original evidence uses the qualified string binaries.
Records: `build/port-spell-scalar-boundaries/old-baselines-{approved.json,deleted.jsonl}`.
