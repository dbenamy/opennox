# Unused generated dispatch signatures

Qualified on Linux 386/SSE2. Baseline commit: `119b666a`; original source is
`62873a54`, the qualified remaining-fixture conversion.

The generator emitted 76 generic C callback signatures, but only 19 have Go
callers in the repository (one is fixture-only). An AST import/selector audit and
literal search through tracked Go, assembly, C and header source found the other
57 names only in their definitions. Selected dependency metadata lists only
project packages as importers. This removes unused source; it does not translate
57 live callback routes.

The existing generator now emits those 19 signatures. Every retained C definition
is byte-identical, every retained Go wrapper is AST-identical, and regeneration
reproduces the installed dispatcher exactly. No caller, owner, identity, fallback,
ABI layout, scalar conversion or external binding changed.

| Metric | Before | After |
| --- | ---: | ---: |
| Generic embedded C dispatch bodies | 76 | 19 |
| Total embedded production C bodies | 77 | 20 |
| Client / server selected cgo files | 80 / 81 | 80 / 81 |
| Selected legacy C exports | 10 | 10 |
| Tracked headers / physical lines | 157 / 2,741 | 157 / 2,741 |
| Standalone production / test-reference C lines | 0 / 0 | 0 / 0 |

Original root evidence reuses the preceding full-default and focused qualification
with exact source identity. `TestLegacyCallbackAdapters` additionally passed twice
per original profile, checking callback choice, pointer arguments, signed return
boundaries and object callback invocation. After conversion it passed in all three
profiles. The full default corpus passed 2,478 roots with one established
`TestMapPopulationPrerequisiteProbe` skip; focused server/high-resolution suites
passed 559/561 roots with exact names and no skips. Captures/assertions are unchanged.

Safe/static and all production/ABI checks passed. The full asset suite exactly
matches known failure and package results. Fresh headless character creation,
save/load and resume passed on the final production binary; all 1,654 original
asset hashes remain unchanged. Accepted gates share source fingerprints.
A separate GUI preview was unnecessary for deleting unused signatures with
identical retained bodies; the final gameplay gate remained.

Primary completed this batch while Luna was quota-limited; no substitute model.
See [baseline](unused-dispatch-signatures-baseline.json),
[qualification](unused-dispatch-signatures-qualification.json),
[inventory](unused-dispatch-signatures-inventory-after.json),
[manifest](unused-dispatch-signatures-batch.json) and
[selection](unused-dispatch-signatures-tests.txt).
Local artifacts: `build/port-unused-dispatch-signatures/`.

Disk cleanup losslessly archived 36 historical capture groups (372,830,208
allocated bytes), removed four superseded book production/safe binaries
(187,338,752 bytes; rebuild `b2597f97`), and removed 20 obsolete project cache
archives predating `119b666a` (917,471,232 bytes). Hash, source/replacement and host
use checks were recorded; no compiler was active during cache deletion. Restore
capture groups by decompressing their recorded archive to the first path,
verifying SHA-256, then hard-linking the remaining recorded aliases. Final
scenario asset deduplication reclaimed 559,968,256 bytes with a restoration
manifest. Original assets, current binaries, saves and qualification evidence remain.
All completed cleanup and finalization scripts are consumed.
