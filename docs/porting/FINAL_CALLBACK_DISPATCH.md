# Remaining native callback dispatch

## Scope and producer review

Retire the eight remaining generic C dispatch bodies. Seven production owner files
cover object create/init/update/use/death/drop/pickup/collision/transfer/damage-sound,
drawable draw/update, inventory drop, particle update and spatial-force callbacks.
Keep the pure-Go key registry; remove its unused C dispatcher output and generator.
Seventeen fixture/helper files move shared observers to native registrations or
remove their private C-typed adapters. Native algorithms and existing assertions
remain unchanged.

Production named registrations already install native handlers; the only remaining
non-native named initializer registrations carry nil placeholders. Parser/object
copy paths carry those keys. Audits include named fields, raw callback offsets,
secondary drawable-update slots and helper-returned fixture keys. Particle slot124
has one raw fixture observer; the separate Go particle system uses handles/native
functions. Spatial-force production callers pass nil callbacks; its recorder is
fixture-only. Unknown keys in these audited families will panic explicitly.

Preserve the shared death observer across create/init/update/death and the transfer
observer across transfer/damage-sound/use/drawable-update. Do not add fixture names
to parser tables. Initializers keep native one-argument priority and a separate
argument-aware registration for the controls observer. Inventory drop passes a
live position and exact signed result; server drop keeps copied-position and
boolean behavior. Preserve signed results, typed nils and existing lifetime barriers.

Collision observer keys retain a zero low byte through aligned interior pointers
into global byte arrays. Capture counters retain each original cap, count and reset
policy. Native bounds panics replace invalid C array writes outside the qualified
capture domain. Spatial force preserves float bits, opaque integer words and GC
ordering. Screen-particle allocation uses its existing matching native 52-byte type.
The controls fixture retains C.free for freshBots until allocator ownership work.
Five private typed wrappers have only their migrated temporary-update callers.
The unused sound export/flag can retire; its native audio observer remains.

## Baseline and qualification plan

Use qualified damage/monster source `945c190b`. Reuse its broad evidence only after
exact source, supplemental inputs, environment, binary/log hashes and discovered
name checks. Run 28 focused original roots in each default/server/highres profile,
three safe forwarding contracts, and private `TestLegacyCallbackAdapters` separately
in default/server/highres/safe. Existing assertions and captures stay frozen.

The broad selection remains 2,482/2,471/2,482 roots with eight explicit predecessor
exclusions. This is not the complete corpus; the baseline records the excluded
names without classifying them all as known failures. After conversion, run native
focused/static and private adapter checks, a fresh default-client save/load before
the broad sweep, all broad roots and focused repeats, safe contracts/build/static,
three production builds/ABI, exact known-suite outcomes and final fresh save/load.
Check original asset hashes and build-selected inventory afterward.

Removing the generator also removes one no-test package. The new known-suite
expectation removes only the callgen package's skip row; all other rows, including
every known failure, stay byte-identical. Original production evidence uses the
preceding expectation; native qualification uses the adjusted one.

## Review and delegation

Primary reviewed the 24-file draft and three deletions against original source.
Function comparison reports 35 additions, 63 changes, 20 removals and 355 unchanged;
these are review aids, not coverage claims. Review caught three missed world-log
reads and a missing types import in the ignored draft before any compilation.
Luna remains unavailable due to quota; no substitute model was used.
Artifacts: `build/port-final-callback-dispatch/`. The conversion is qualified against accepted baseline commit `0807f897`.

## Local recovery

Before baseline runs, verified and removed 35 old Linux386 Go cache archives,
recovering 1,343,545,344 allocated bytes. Host compiler/process/file-use checks
passed; source, original assets, current qualified binaries and module sources
remain. Exact hashes and paths are in `cache-headroom-approved.json` and its
journal. These cache entries rebuild automatically.

Original baseline accepted: all reuse checks and fresh contracts passed at
`945c190b`. See [baseline](final-callback-dispatch-baseline.json).

Removed seven superseded GUI binaries after committed-source, replacement hashes
and host-use verification, recovering 386,985,984 bytes. Rebuild `0cf5064c` if
needed; records remain in `gui-cleanup-approved.json` and its journal. Current
damage/monster binaries and all original logs remain.

During native tests, removed seven superseded native-integer binaries after the
same committed-source, replacement and host-use checks (387,166,208 bytes).
Rebuild `a7dc3a36` for those historical outputs; logs/metadata remain. Records:
`integer-cleanup-approved.json` and journal. Active tests use separate new binaries.

Before the final scenario, verified and removed four remaining superseded
unused-export build binaries (193,003,520 bytes). Rebuild `f6f5ee4c` if needed.
The old test binaries were already absent; the initial seven-path validation
stopped before deletion, then a four-path plan passed all source/hash/host-use
checks. Records: `unused-exports-cleanup-approved.json` and journal.

## Qualification and progress

Native broad selections passed **2,482/2,471/2,482** roots without skips. All 28
focused roots per profile, three safe contracts and private legacy adapters in
four profiles passed before/after conversion. Native focused/static checks and
fresh default save/load passed before the broad sweep. Safe build/static, three
production builds/ABI, exact known-suite comparison and final fresh save/load
passed. All 1,654 original asset hashes are unchanged. Assertions and captures
remain unchanged; retired C dispatchers/observers are absent from qualified builds.

Embedded production C bodies: **8→0 (79/79 retired)**. Production cgo files:
**5→4 client/highres; 6→5 server**, across **two** project packages. Fixture C
imports: **83→69**. Headers remain **157 files /2,731 lines**; selected legacy C
exports and standalone production/test C remain **zero**. These are dependency
counts, not an estimate of remaining effort. The internal-glue milestone is not
complete: allocator ownership, reachable abort, flags and fixture/type dependencies
remain. Evidence: [qualification](final-callback-dispatch-qualification.json),
[inventory](final-callback-dispatch-inventory-after.json).

The evidence collector initially used an overly broad `go_call_` substring check,
which also matched Go runtime cgo traceback/symbolizer helpers. It now checks the
eight exact engine dispatcher names and their wrappers; those are absent. Native
source and successful qualification runs needed no correction or repetition.

After all qualification jobs finished, verified and removed 3,308 identical
scenario asset copies (1,120,485,376 allocated bytes). Host-use checks passed and
original assets, outputs, binaries and metadata remain. Restore inputs with:

```sh
python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/final-callback-dispatch-preflight-save/deduplicated-assets.json
python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/final-callback-dispatch-save/deduplicated-assets.json
```
