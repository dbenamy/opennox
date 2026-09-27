# Native record, handle and list boundaries

Status: original baseline qualified; conversion draft is not installed.
Production baseline: `1b80dbe1` (qualified scalar boundaries).

Replace private C record aliases with the existing Go owners: Object, Drawable,
Player, Window and Team. Remove unused Thing, ScreenParticle and Waypoint aliases;
keep the independent native Waypoint size assertion. This changes type spelling
at owner/caller boundaries, not memory allocation, field layout or ownership.
The original-source AST audit found 148 pointer-type uses, 34 C-qualified type
selectors and six alias declarations, with no record-by-value construction or
size-based allocation through these aliases.

The [original layout capture](native-layout-types-original-layout.json) records
independently compiled 386 C sizes, alignments, constants and field offsets, with
a byte-identical repeat and reflection of existing Go owners. Object and Drawable
have C prefixes of 772 and 512 bytes and existing Go extensions (full sizes 780
and 516). Player has matching size/offsets but C alignment 1 versus Go alignment 4.
These distinctions are preserved; the conversion does not equate full Go sizes
with legacy C prefixes. Independent contracts freeze the native sizes, offsets,
alignment and extension positions.

`FILE` becomes a non-zero-sized opaque type for protected-region handle pointers.
It is never dereferenced or allocated as libc storage. Real file lookup, cursors,
closure and retained handle identity stay unchanged. The new contract exercises
two real files, repeat lookup, independent/shared cursors, nil/idempotent close,
underlying file closure and the existing stale identity after closure. Its first
run omitted handle-region initialization; adding the established Init/Release
fixture setup made both new contracts pass on unchanged production code.

Rules use the existing 12-byte native list node, retaining next/previous links
and the raw third word. Fixture allocation/free observers stay in place. Point
and rectangle adapters use 32-bit native coordinates with the existing field
order and rectangle normalization. The sole trace-wrapper caller uses its native
owner directly, retaining signed-byte flag narrowing. Wall damage keeps a wrapper
to preserve argument evaluation relative to the replaceable server lookup.
String reads use the existing GoStringP bridge, preserving its current behavior.
The inventory background flag uses gui.StatusBelow. The compiled maximum window
size remains 3840×2160 in every profile because video_highres.go already defines
NOX_HIGH_RES without a build tag.

Five orphan private helpers and three unused aliases are removed after whole-source
Go/header/preamble searches; the trace wrapper is removed with its sole caller
migrated. The constant-false diagnostic in AsWindowP is retired. Reachable aborts,
actual libc operations, allocator ownership and foreign callback fallbacks are
outside this batch. External native libraries remain unchanged.

The original baseline starts with 1,230 default/highres roots and 1,228 server
roots. A separate signature/type-use audit adds 82 roots per profile for a final
selection of 1,312/1,310/1,312. Type identifiers require explicit AST handling; the
function-reference graph alone excludes them. Selection follows changed functions
and alias users through package and root-fixture helpers, then includes entire
affected assertion files. The actual C spell callback retains explicit C pointer
casts; its native callers belong to the expanded selection. Server omits
two client-only hover roots; the established opt-in map-population diagnostic is
excluded. All 1,312/1,310/1,312 selected roots pass twice per profile in separate
processes, with exact name-set verification and no skips. See
[baseline evidence](native-layout-types-baseline.json). The first launcher
omitted asset environment settings and was stopped; none of its partial results
are accepted. Its verified source-matching binaries are reused by the corrected
launcher with the full manifest environment.

Planned converted gates: focused preflight; complete default port corpus and
selected server/highres roots; safe/static and three production/ABI builds;
exact known-suite outcomes; fresh headless save/load/resume; original asset hashes.
Full default coverage is appropriate here because shared record aliases span
many owners. Root assertions and frozen captures stay unchanged except for the
two independently authored contracts above.

Primary owns this batch while Luna is quota-limited. No substitute model is used.
Local drafts, audits, probe captures and test evidence:
`build/port-native-layout-types/`. Draft and refinement generators are consumed.
