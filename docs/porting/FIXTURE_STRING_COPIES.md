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

## Status

Original baseline accepted: six owner roots twice per profile, six safe roots,
and two private contracts in each of default/server/highres/safe passed. Draft
not installed. Expected fixture imports
14→11; production remains 4 client/highres and 5 server. Standalone C, production
embedded bodies and legacy export bridges remain zero. Headers remain 157/2,731.
Primary handles this batch while Luna quota is unavailable.
Artifacts: `build/port-fixture-string-copies/`.
