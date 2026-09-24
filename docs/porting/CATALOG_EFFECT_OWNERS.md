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

Qualification after conversion is pending. Evidence:
[captured baseline](catalog-effect-owners-baseline.json),
[batch commands](catalog-effect-owners-batch.json),
[focused selection](catalog-effect-owners-tests.txt).
