# Native scalar owner and caller boundaries

Status: qualified. Production baseline `a7dc3a36`; scalar baseline `097febf7`.

Seventeen production files used C solely for character, floating-point, boolean
and integer types. Replace them with the exact native representations emitted by
the qualified Linux 386 cgo/GCC toolchain: `char → int8`, `uchar → uint8`,
`float → float32`, `bool → bool`, `longlong → int64`, `ushort → uint16`,
`short → int16`, and explicit 32-bit signed/unsigned integers. The isolated
compiler proof and flag review are recorded in the qualification evidence.

The conversion covers 24 private C-typed signatures and their callers across
27 files: 48 changed and 304 unchanged functions, with none added or removed.
Package AST inspection found no non-call references to these functions.
Six string-reading owners and one return-value caller use the existing
`GoStringP` bridge, retaining C.GoString behavior and ownership. Message and
damage calls use their existing Go owners with equivalent arguments, narrowing,
evaluation order and ignored returns. Arithmetic and float rounding boundaries
remain unchanged. One empty fixture C import and two unused prefab fixture
prototypes are removed. Whole-source header/preamble searches found no live C
callers. Actual libc operations and external native libraries remain unchanged.

The first compile caught three missing `unsafe.Pointer` conversions from
interned `*int8` strings to the message owner's `*byte` arguments. The corrected
calls match the conversions performed by the old wrapper. No tests started in
that failed build. Corrected preflight passes all 25 exact roots and static checks.
Root assertions and frozen expectations remain byte-identical.

Existing independent contracts cover signed boundaries, flag predicates, clocks,
text-buffer identity and guards, byte-string lookup, prefab precision/ownership,
damage, poison, saved-player framing and serialized messages. The selection
combines reverse package-function references, transitive fixture/root helper
calls, entire affected assertion files and conservative owner-family inclusion.
An initial graph that confused common method/field names with package functions
was rejected before execution. Extra related owner tests are intentionally
included; selection is not a claim that every root is a direct caller.

Original qualification passes 968 default/high-resolution and 966 server roots
twice, using exact-source/binary/environment reuse for 166/164/166 first
observations, 802 additional roots per profile, and complete repeats in separate
processes. Server omits exactly two `!server` hover tests. The established opt-in
`TestMapPopulationPrerequisiteProbe` diagnostic is explicitly excluded.
See [baseline evidence](scalar-boundaries-baseline.json).

Converted contracts pass the same 968/966/968 exact root-name sets without failures
or skips. Safe/static checks, all three production/ABI builds, the exact known
failure multiset and package outcomes, fresh headless character creation/save/load/
resume, and all 1,654 original asset hashes pass. Accepted gates share the reviewed
source fingerprints. No width-mirroring tests, full converted corpus or separate
GUI preview are needed for this representation-only scope; new failures or
uncertain dependencies still require broader qualification.

| Metric | Before | After |
| --- | ---: | ---: |
| Selected project cgo files, client/server | 52 / 53 | 35 / 36 |
| Selected legacy C exports | 0 | 0 |
| Embedded production C bodies | 20 | 20 |
| Headers / physical lines | 157 / 2,731 | 157 / 2,731 |
| Standalone production / test-reference C lines | 0 / 0 | 0 / 0 |

See [qualification](scalar-boundaries-qualification.json) and
[dependency inventory](scalar-boundaries-inventory-after.json). Primary handled
this batch because Luna remains quota-limited; no substitute model was used.

Before testing, 198 inactive Linux 386 Go cache archives were removed after exact
path/stat/hash and host-use checks, reclaiming 4,976,803,840 allocated bytes.
Source, assets, module downloads and qualified binaries remain; cache misses
rebuild normally. Proposed old-executable archival plans were not executed.
After qualification, 1,654 verified scenario asset copies were deduplicated
(559,841,280 allocated bytes). Original assets, saves and results remain, with
restoration manifest in `build/baseline/runs/scalar-boundaries-save/`.
Journals and local evidence are under `build/port-scalar-boundaries/`.
All completed cleanup, installation and acceptance scripts are consumed.
