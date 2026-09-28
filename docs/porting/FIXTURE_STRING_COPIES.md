# Fixture string and byte copies

## Scope and review

Three fixture files replace C.CString with the existing legacy CString helper,
and C.CBytes with legacyMalloc plus an exact-length copy. Existing legacyFree
calls and the nil map-cycle path remain unchanged. Private wrapper C char/int
casts become int8/int32 on the qualified target. Three fixture C imports retire.
No production code, root test assertions or frozen captures change.

Go 1.26 cgo wrappers and the existing normal/safe allocator adapters were reviewed
side by side. String copies preserve embedded bytes and add one terminal NUL;
byte copies add none. Normal allocation retains cgo's fatal failure behavior and
its conditional malloc(0) retry. Safe allocation retains tracked ownership and
recoverable failure. The package-wide safe malloc/free defines are part of this
review; merely inspecting included headers would miss them.

AST review found five changed functions, 23 unchanged, none added or removed.
Both private wrappers are confined to their own fixture files. The ten-header
removed include closure has no startup hooks. Caller traversal identifies six
root owners: four command-rule tests, rule removal and map-cycle line endings.
Existing cases cover empty names, embedded NULs, line endings, nil input and
unchanged trailing bytes. Empty nonnil unterminated map-cycle buffers are outside
the valid caller contract; this conversion does not expand that domain.

The original internal-glue plan has been aligned with PORT.md: whole-build cgo-off
is not required while external native bindings remain. Its obsolete deferred
monster-callback note now points to the live checkpoint. This is a documentation
correction, not a change to the agreed milestone.

## Qualification plan

Run all six owners twice per profile and under safe before conversion. Run private
legacy TestStringBoundaryBytes and TestStringMallocFailure in default/server/highres
and safe, checking exact root names and allocation disposition. Preserve accepted
original evidence before installing the draft. Converted qualification starts with
safe contracts, followed by private contracts, focused preflight, all affected
owners and repeats, static checks, safe/production builds, ABI/known-suite checks
and fresh save/load. Every target binary compiles the complete root test package.

## Qualified result

Original and converted paths passed all six affected roots twice per profile and
six safe contracts, without skips. Two private string/allocation contracts passed
in each of default, server, highres and safe. Converted safe contracts ran first.
All target binaries compiled every root; static checks, safe build, production
builds/ABI, exact known-suite outcomes and fresh save/load passed. Root assertions,
frozen captures and original asset hashes are unchanged.

Fixture C imports fell 14→11. Production remains 4 client/highres and 5 server;
standalone C, embedded production bodies and legacy exports remain zero. Headers
remain 157 files / 2,731 physical lines. External native backends are unchanged.

Evidence: [baseline](fixture-string-copies-baseline.json),
[qualification](fixture-string-copies-qualification.json),
[inventory](fixture-string-copies-inventory-after.json).
Primary handled review and qualification while Luna quota was unavailable.
Artifacts: `build/port-fixture-string-copies/`.

## Artifact cleanup and recovery

After qualification and all jobs completed, verified scenario-copy deduplication
reclaimed 560,046,080 bytes. Restore with:

```sh
python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/fixture-string-copies-save/deduplicated-assets.json
```

Seventeen superseded binaries reclaimed another 912,945,152 bytes: ten from shared
fixture types (6048add7, including three byte-identical outputs from its interrupted
scenario attempt) and seven from corrected fixture allocation (d2df2eb6). Sources
matched Git; newer qualified replacements, hashes and host process/open-file/maps
checks passed. Rebuild those revisions if needed. Source, current binaries, useful
logs and original assets remain intact. See [cleanup record](fixture-string-copies-cleanup.json).
