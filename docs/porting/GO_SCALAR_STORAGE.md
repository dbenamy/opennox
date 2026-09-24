# Go integer types for existing scalar storage

Status: **qualified** against frozen baseline `a5061480` (original source
`12bc387d`). See [qualification](go-scalar-storage-qualification.json),
[baseline evidence](go-scalar-storage-baseline.json), [selected dependency inventory](go-scalar-storage-inventory-after.json)
and [batch manifest](go-scalar-storage-batch.json). Original production/root
baseline evidence was reused by exact source identity from
[the preceding export batch](go-only-exports-qualification.json).

## Scope

Replace C integer types on 395 existing Go globals: 385 `C.uint32_t` and one
`C.uint` become `uint32`; three `C.int` become `int32`; six `C.uint64_t` become
`uint64`. Keep every owner, initializer bit, signed interpretation, pointer alias
and lifetime. These globals already live in Go storage; this changes their type
identity without moving storage or making integer address words into pointers.
The supported target remains Linux 386/SSE2.

Migrate only necessary scalar conversions and preserve casts at remaining
C-typed function/field boundaries. Keep allocation domains, callbacks, export
signatures and external native bindings unchanged. The qualified conversion removes 28 production and eight test C imports/preambles.
Selected production cgo files fall from 382 to 354 in every profile. The three
direct project cgo packages, 1,179 legacy exports, 79 callback bodies and 157
headers / 3,902 lines remain unchanged. Standalone C stays zero.

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
package globals. The original draft parses. Subsequent compiler-boundary fixes were reviewed
separately, including an independent Luna review. Exact reconstruction verifies
all 125 installed source files against the reviewed draft plus those fixes.

The 35 initial C-import removals have no remaining selectors, exports or build
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

## Compiler boundary review

The first compile found 18 diagnostics in six files: a four-word snapshot still
used C element types, protection fixtures retained C numeric casts in assignments
and comparisons, a rules fixture had a mixed multi-assignment, and player-reset
mask clearing retained a C cast. Convert those values to the same-width Go types;
preserve operands, assertions and actual C call/field conversions. The reset file
then has no C selectors, exports or build directives. Its preamble included only
`client__gui__chathelp.h`, an empty guarded header including the already-reviewed
`defs.h` closure, so remove that import too. Final removals are 28 production
and eight test C imports. The failed compile log is retained; the second compile passes. Post-conversion storage contracts pass in default/server/highres and scalar-only
safe, with exact frozen hashes and no skips. All three complete consumer sweeps, safe/static and fresh production
qualification pass. Exact reconstruction
checks all 125 installed source files against the reviewed draft plus compiler
fixes. All 1,654 original asset hashes remain unchanged.

## Local artifact cleanup

After the consumer controller completed, host process/descriptor/mapping checks
and fresh source/copy hashes allowed removal of 3,242 identical asset copies from
the completed 378-export and 265-export scenarios: 1,112,511,348 logical bytes
(1,119,281,152 allocated bytes). Original assets, configurations, saves and records
remain. Recovery manifests were written before deletion. No hard links were made;
the helper's hard-link suggestion was not used. See each scenario's
`deduplicated-assets.json` and local `cleanup-completed.json`.

## Qualified outcome

Default/highres each pass 2,425 root tests; server passes 2,414. Each has only the
expected `TestMapPopulationPrerequisiteProbe` skip, with exact root-name sets and
no failure events. Safe build/static and three fresh production/ABI checks pass.
The full-suite comparison exactly matches the existing 304 failure events and
17-pass/two-fail/32-skip package baseline. Headless character creation, save, load
and resume pass. All phase source fingerprints match. Original assets were
rechecked after production qualification and all 1,654 hashes are unchanged.

Phase progress is now 109/463 selected cgo files eliminated on net (354 remain);
711/1,890 export bridges retired (1,179 remain). No new export was retired here.
The C-type changes preserve storage ownership and x86/32-bit assumptions.

Luna also supplied a bounded read-only follow-on inventory and disk-copy audit.
These helped preparation; the primary verified all accepted changes, qualification
and deletion conditions. No measured subscription or cost-saving claim is made.
