# Owned tile-grid fixture helpers

## Scope and review

Retire four C helpers in tile_worklist_porttest.go: row allocation, release, cell
writes and byte equality. Native helpers retain a 128-pointer table, 128 separately
allocated rows and 44-byte cells. Fields 1 and 6 stay at offsets 4 and 24, using raw
uint32 words to preserve signed C bit patterns. Independent compiler probes confirm
pointer width, cell size/alignment and offsets in all three target profiles.

Keep allocation order, reverse partial-failure cleanup, forward complete cleanup
and the nil outer-pointer guard. Existing legacyCalloc/legacyFree adapters retain
normal raw and safe tracked ownership. Equality still compares all physical bytes.
Native slice bounds checks affect only indices outside the original fixture's
valid range. Production code and previously frozen assertions/captures are unchanged.
Nine headers in the removed include closure contain no startup hooks.

## Independent original-path contracts

The baseline adds TestTileWorklistAllocationCleanup against the original C helpers.
It fails the outer allocation and each of 128 row allocations, plus two success
cases. It checks allocation sizes, zero-filled cells, valid frees, complete cleanup
and tracker balance. It already passed focused preflight in normal and safe builds.
Safe Calloc currently records a nil marker on failure; the fixture verifies that
behavior and removes only its own marker after observation. This does not change
the allocator's behavior. Three small Go wrappers expose the existing C observer;
its implementation and OS-thread scope are unchanged.

## Caller selection and qualification plan

The conservative helper graph reaches 563 roots. Main-grid setup in PortTestRoam
requires a nonnil Main specification; auditing its typed producers still reaches
545 roots because spell/callback/damage fixtures inherit that base. All 545 are
already in the conservative selection, which is retained. Add the existing foreign
thread observer contract explicitly: 564 owners in each profile. Every target
compiles all root tests. Repeat nine focused owners and six safe contracts.

Fresh original binaries include the new allocation contract; preceding binaries
are used only for test-name discovery. Accept and commit the original baseline
before installing the conversion. Native qualification starts with safe contracts,
then focused preflight, exact owner selections and repeats, static checks, safe
build, production builds/ABI, exact known-suite outcomes and fresh save/load.
Review found four changed Go functions and four newly native private helpers;
four other Go functions are unchanged.

## Build-cache headroom

After all baseline jobs finished, 105 reproducible Linux 386 cache archives of at
least 10 MB, untouched for two hours, were removed. Hash, file identity, archive
format and host compiler/open-file/mapping checks passed. This reclaimed
4,423,565,312 bytes; cache entries rebuild automatically. Current binaries, module
sources and original assets remain intact. See
[cleanup result](fixture-tile-grid-cache-cleanup.json). The shorter cache age is a
reversible disk-headroom choice, with active-use checks still required.

## Qualified result

Original and native owner selections passed 564 roots per profile without skips.
Nine focused repeats per profile, six safe contracts, static checks, safe build,
production builds/ABI, exact known-suite outcomes and fresh save/load passed.
The independent allocation contract checks all 129 failure positions and two
success cases, including safe tracking. Native safe contracts ran before the broad
sweep. All target binaries compiled every root; prior captures/assertions and
original asset hashes are unchanged.

Four fixture C helper bodies and one fixture C import are retired (11→10 files).
Production cgo remains 4 client/highres and 5 server; embedded production bodies,
legacy exports and standalone C remain zero. Headers remain 157/2,731 physical lines.

Evidence: [baseline](fixture-tile-grid-baseline.json),
[qualification](fixture-tile-grid-qualification.json),
[inventory](fixture-tile-grid-inventory-after.json).
Primary handled review and qualification; Luna quota remained unavailable.
Artifacts: `build/port-fixture-tile-grid/`.

## Runtime observation and scenario recovery

Owner execution times were 341.8→328.8 seconds (default), 329.8→324.7 (server),
and 331.0→326.5 (highres). These are single qualification runs, with some isolated
probe work overlapping the native runs; they do not establish a controlled
performance improvement. See [observations](fixture-tile-grid-runtime-observations.json).

Completed scenario-copy deduplication reclaimed 559,837,184 bytes after host-use
and original/copy hash checks. Restore with:

```sh
python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/fixture-tile-grid-save/deduplicated-assets.json
```

Original assets and current binaries remain intact. See
[scenario cleanup](fixture-tile-grid-scenario-cleanup.json).
