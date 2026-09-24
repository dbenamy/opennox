# Go integer types for existing scalar storage

Status: original storage baseline qualified on source `12bc387d`; conversion
is uninstalled. See [baseline evidence](go-scalar-storage-baseline.json). Production/root baseline evidence is reused by exact source
identity from [the preceding export batch](go-only-exports-qualification.json).
See [batch manifest](go-scalar-storage-batch.json).

## Scope

Replace C integer types on 395 existing Go globals: 385 `C.uint32_t` and one
`C.uint` become `uint32`; three `C.int` become `int32`; six `C.uint64_t` become
`uint64`. Keep every owner, initializer bit, signed interpretation, pointer alias
and lifetime. These globals already live in Go storage; this changes their type
identity without moving storage or making integer address words into pointers.
The supported target remains Linux 386/SSE2.

Migrate only necessary scalar conversions and preserve casts at remaining
C-typed function/field boundaries. Keep allocation domains, callbacks, export
signatures and external native bindings unchanged. The reviewed draft proposes
removing 27 production and eight test C imports/preambles; measure the final
selected counts after compilation and qualification. Standalone C stays zero.

## Original contracts and acceptance

The existing `legacy` package contracts match all 395 owner names exactly.
`TestScalarStorageOwners` checks 139 patterns per owner (54,905 cases): widths,
initial bits, signed reads, typed/raw alias visibility, distinct storage, unchanged
neighbors and GC while patterned numeric data is installed. `TestRawStorageOwners`
adds layout, registry, raw-owner and scalar-owner disjointness checks. These
package-local contracts are separate from the root consumer sweep.

Before installation, run both contracts twice in independent default processes
and once each in server/highres; run the scalar contract separately in safe.
Compare captures against existing frozen
hashes from raw-storage-native-batch.json; do not regenerate expectations. Review
original failures before freezing. After conversion, repeat both contracts in the three normal profiles
and the scalar contract in safe, plus complete default/server/highres root sweeps, safe build/static,
fresh three-profile production/ABI, exact known-suite comparison, and headless
character creation/save/load/resume. Reuse the preceding root/production baseline
only under verified source identity.

## Review and delegation

Primary identified the owner boundary, fixed-width mapping, existing contracts and
required C-interface conversions. Luna drafted 395 declaration edits, 52 matching
initializer casts and 702 simple assignment casts, leaving multi-assignments
unchanged. Primary caught a uniform generation error: the draft changed selector
suffixes while retaining `C.` prefixes. The independently reconstructed reviewed
copy replaces whole selectors. Primary also handled 11 multi-assignments with
23 same-width casts.

An independent Go AST check verifies all 1,172 initial edit contexts across 122
files: package declarations/initializers or direct assignments to unshadowed
package globals. The uninstalled draft parses; actual compiler/boundary fixes
remain. Exact reconstruction preserves all other bytes before import cleanup.
Future boundary edits must be reviewed separately, not inferred from that check.

The 35 proposed C-import removals have no remaining selectors, exports or build
flags. Primary reviewed their full preambles and 24-header closure. Payloads are
includes plus one prototype, without executable definitions. No candidate has an
`init` function. Preserve all other Go initialization. In particular,
`scalarStorageInitial` is an eagerly invoked capture closure, despite the helper's
initial description as a function value; its original initial-value checks must
still pass. The tile callback initializer is a function value. Numeric conversions
and compile-time size assertions keep their original meaning.

Bounded edit manifests and import inventories were useful, but both mechanical
spans and interpretation needed correction. Keep algorithm design, binding review,
baseline acceptance, integration and qualification with the primary. No measured
cost-saving claim is made.

Local work: `build/port-go-scalar-storage/`. Use `reviewed/`, not the superseded
`draft/`. Completed reconstruction/cleanup scripts are consumed. No source change
is authorized while a build/test is reading the checkout.

## Baseline outcome and safe fixture boundary

Default, an independent default repeat, server and highres pass both storage
contracts with exact test-name sets and no skips. Safe passes the scalar contract
separately. Every capture equals the existing frozen hashes; production/test
source is unchanged from `12bc387d`. No golden or engine code was changed.

The first optional safe run attempted both contracts. Raw storage intentionally
probes every registered field offset in the original blob, including slots whose
owners have moved to Go. Safe's runtime address check rejects the first such
intersection (`byte_581450_1416`) at raw_storage_porttest_test.go:67, before the
scalar test starts. Preserve this failed-run evidence; do not disable the guard
or rewrite expectations to make that fixture pass. Its raw-alias coverage stays
in the three normal profiles. Safe scalar coverage was then run and passed in a
fresh process. No complete safe runtime suite or safe raw-fixture pass is claimed.

Also remove the empty import block left in `session_entry_exports.go` by the
preceding wrapper retirement; the file has no remaining declarations or effects.
