# Direct Go scalar and pointer forwarding

Planned scope: 69 existing Go→C→Go calls (29 scalar calls and 40 pointer calls),
plus ten C exports with no remaining callers or address references. The existing
Go implementations remain. The rejected `sub_50B510` prototype mismatch remains
outside this batch.

## Type and ownership review

The [signature inventory](typed-direct-calls-signatures.json) records 41 wrapper
signatures and every selected site. Primary compared them against freshly
obtained cgo metadata from qualified checkpoint `a69e1da8`. Pointer types were
normalized through explicit Go/C aliases, never merely similar layouts. No
function-name macro aliases were found. Scalar argument/result conversions retain
C widths on the supported 386/SSE2 target; no new 64-bit compatibility claim.
All nine random-number calls retain their original evaluation positions, bounds
and surrounding arithmetic. The generator and its state are unchanged.

Pointer calls preserve identities and existing conversions. Movement copies a
point by value before calling the server owner. Start-position and damage output
pointers remain typed and synchronous. Respawn reads/copies the attribute array
into object storage; it does not retain that borrowed array. Map loading converts
the interned C string to a Go string before passing it to the server. Existing
object/player owners and returned object pointers remain unchanged. These calls
are not assumed pure just because their signatures match.

The ten retired-interface candidates were checked for Go C-selectors, calls and
address uses inside C preambles/headers, and other tracked tool/external references.
Only their prototypes remain in source C regions. Their same-package Go functions
remain for direct callers. Other exports with actual callback or test-driver
references are excluded; a selector-only search is insufficient evidence of an
unused interface.

## Baseline and gates

The exact preceding 380-root qualification is reused after verifying all
production/test source fingerprints, environment settings and four retained
binary hashes. A fresh additional selection passes 110/109/110 roots in
default/server/highres, without skips. The sole server exclusion is
`TestClientInventoryWindowWorldSelection`, explicitly built with `!server`.
The initial finalizer assumed a flat110 additional roots and correctly stopped;
its count check was corrected after inspecting that build constraint and the
actual discovered-name difference. No tests or expectations were changed.
See [baseline evidence](typed-direct-calls-c-qualification.json).

The complete converted selection must pass 490/489/490 roots. Its discovered names
must exactly equal the union of reused and freshly executed baseline selections.
Retain existing frozen captures and independent expectations. This covers affected
owners, not every possible call-site branch. Safe/static, fresh production/ABI,
exact known full-suite comparison, headless character creation and save/load
remain required after conversion.

## Delegation and recovery

Luna supplied both bounded call drafts; primary owns signature normalization,
integer/ownership review, scope selection and qualification. The baseline reuse
avoids rerunning identical recently qualified tests; new owner contracts and all
converted contracts still execute. No game-speed improvement is claimed.

Artifacts: `build/port-typed-direct`; drafts and candidate inventories are under
`build/port-scalar-direct`. Five completed captures were losslessly archived after
matching every hardlink, verifying completion provenance, checking host open files,
and verifying gzip round trips. 677,499,524 raw bytes became 16,273,443 gzip bytes.
All 44 hardlink paths and five symlink aliases are recorded; symlinks are preserved
and become usable again after restoration. Use the archive script's `--restore`
before historical finalizers. `capture-archive-plan.json` and
`capture-archive-record.json` contain hashes, metadata and paths. The archival
script is CONSUMED; restoration remains supported. Assets and binaries are intact.

Standalone production/reference C remains **zero files / zero lines**, with
**79 production C preamble bodies**. This metric excludes headers, generated
bridges and external libraries.
