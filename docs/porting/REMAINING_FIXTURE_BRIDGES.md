# Remaining fixture export bridges

Qualified on Linux 386/SSE2. Baseline commit: `61bca9ef`, with original production
and test source from `b2597f97`. No assertions or captures were regenerated.

Seventy selected engine C exports had no production C consumer. Sixty existing
Go adapters were fixture-only and moved under `porttest`; ten have real Go callers
and keep their bodies without C exports. The ten live production C callback
addresses remain. This removes 15 production cgo imports and five empty files.

| Metric | Before | After |
| --- | ---: | ---: |
| Client / server selected cgo files | 95 / 96 | 80 / 81 |
| Selected legacy C exports | 80 | 10 |
| Tracked headers / physical lines | 157 / 2,809 | 157 / 2,741 |
| Embedded production C bodies | 77 | 77 |
| Standalone production / test-reference C lines | 0 / 0 | 0 / 0 |

All 70 adapter signatures/bodies and 374 remaining production functions in the
touched files are AST-identical to the original. Native owners, layouts and
external native-library bindings stay unchanged. Transitional C scalar types in
fixture adapters remain for subsequent cleanup.

Whole-source review checked Go calls/addresses, import-C preambles, headers and
macro aliases. Protection setters retain signed32, byte and ushort narrowing.
Temporary dispatch preserves argument order, field offsets, sparse capture IDs,
live FlameCleanse callback17 and independent C observers. Five unused callback
addresses become nil; skipping nil identity registration preserves canonical zero.
Two C-int size assertions move with their fixture adapter. The declaration audit
caught thirteen pointer-returning declarations missed by the initial filter,
including a multiline prototype, before installation.

Original and converted full default corpora each passed 2,478 roots with only the
established `TestMapPopulationPrerequisiteProbe` skip. Each of 305 newly selected
roots passed twice per original profile (full default supplies its first run).
Prior same-source book coverage of 256 client/high-resolution and 254 server roots
was reused with source/binary/environment verification. Converted focused suites
passed 561 high-resolution and 559 server roots with exact names and no skips.
Broader indirect fixture sharing is covered by the full default corpus, not claimed
from the focused selection alone.

Safe/static and all three production/ABI checks passed, including absence of all
70 retired exports. Fresh preview and final headless character creation/save/load/
resume passed on the same production binary. The full asset suite exactly matches
the known failure and package-result sets. All 1,654 original asset hashes are
unchanged, and accepted gates share source fingerprints.

The profile controller now accepts an optional default-only pattern, allowing the
full default sweep without a redundant focused-default pass. Existing behavior
without the option is unchanged; all seven orchestration checks pass.
Primary completed this batch locally because Luna remained quota-limited; no
substitute model was used. No qualifying source was edited during active runs.

See [baseline](remaining-fixture-bridges-baseline.json),
[qualification](remaining-fixture-bridges-qualification.json),
[inventory](remaining-fixture-bridges-inventory-after.json),
[manifest](remaining-fixture-bridges-batch.json) and
[focused selection](remaining-fixture-bridges-tests.txt).
Local artifacts: `build/port-remaining-fixture-bridges/`.

Disk recovery removed four superseded binaries (187,367,424 allocated bytes;
rebuild `1238c985`) and three book test binaries (200,716,288 bytes; rebuild
`b2597f97`) after source/hash/replacement/host-use checks. Lossless gzip archival
reclaimed 711,507,968 bytes from 30 historical capture groups, 260,030,464 bytes
from six accepted original logs, and 165,019,648 bytes from three converted logs.
Cleanup/archive journals record exact paths, hashes and restoration instructions.
Restore logs with `gzip -dk`; for capture groups, decompress to the first recorded
path and hard-link the others. Verify the recorded uncompressed SHA-256.
Preview/final each deduplicated 1,654 verified asset copies; restore from their
`deduplicated-assets.json` manifests with the existing restoration tool. Current
binaries, original assets, saves and committed evidence remain available.
