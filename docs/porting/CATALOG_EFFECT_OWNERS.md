# Catalog and client-effect fixture owners

## Frozen baseline

Production source is `200d459b`. A new independent rain-orb creation contract
passes through the original C entrypoint in all three Linux 386 profiles. It
checks allocation failure, pointer/list membership, guarded input storage,
unsigned-short/signed-byte narrowing, and random consumption. All 57 focused
roots pass without skips in default/highres; server passes 56 because the
client-only occlusion contract has a `!server` build constraint.

The conversion targets eleven production C exports whose only callers are test
fixtures: five map catalog/cycle wrappers, two object helpers, three effect
creation wrappers and the curve-segment wrapper. Their existing Go algorithms
stay unchanged. Keep live rendering/update callback exports and their C fixture
routes. No targeted wrapper has a production address registration.

The curve wrapper's specialized C body calls a Go test observer through a function
pointer. Replace that transport with a local Go observer; retain exact ordered
segment/token assertions and point guards. Retire the unused test callback export
as well. No independent C algorithm is removed from test coverage.

Original evidence:
[captured baseline](catalog-effect-owners-baseline.json),
[batch commands](catalog-effect-owners-batch.json),
[focused selection](catalog-effect-owners-tests.txt).


## Qualified conversion

Eleven production C exports and seven private C-typed wrappers are retired.
Catalog/cycle and effect/object fixture calls use existing Go owners; the curve
fixture uses a local observer. One production cgo file and one test cgo import
are removed. The specialized curve C callback body is gone. All live export
bodies and live effect fixture C routes remain text-identical.

All 57 default/highres and 56 server roots pass without skips. Safe/static checks,
three fresh production builds and retained/retired ABI checks pass. Headless
character creation and save/load/resume pass. The full suite matches the known
baseline exactly (304 failure events; 17 passing, two failing, 32 skipped packages).
All phases have identical source fingerprints; all ten changed/deleted files match
the independently reviewed drafts. Frozen expectations, the new rain contract and
1,654 original asset hashes are unchanged.

Selected production cgo files fall 225→224 (239/463 eliminated); exports 989→978
(912/1,890 retired). Embedded production C bodies fall 78→77: 76 generic dispatchers
and one specialized spell adapter. Headers remain 157 files, now 3,701 physical lines.
Standalone production/test-reference C remains zero. External bindings unchanged.

Luna audited caller routes and drafted five files. Primary reviewed widths,
reconstructed all 33 edits, added catalog/header changes, and verified retained
functions/cases. Primary added the missing rain creation contract and declined
unnecessary migration of 23 live update test routes. Curve transport retirement is
intentional: the removed observer was already Go, and no production caller uses
that C interface. Pointer countdown return handling was checked explicitly.
The first baseline controller invocation had a command-argument error before any
build; corrected invocation and all conversion gates pass. Test-name acceptance
accounts for the pre-existing client-only occlusion build constraint.

[Qualification](catalog-effect-owners-qualification.json),
[inventory](catalog-effect-owners-inventory-after.json).
Local artifacts: `build/port-catalog-effect-owners/`.
