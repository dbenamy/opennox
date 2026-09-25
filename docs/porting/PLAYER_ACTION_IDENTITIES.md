# Player-action C bridge removal

## Scope

Remove 32 exports and three private C-typed adapters from player controls and
attacks. Their gameplay algorithms already live in Go. Move actual legacy callers
to those owners and remove the two export-only files and selected header prototypes.
Keep external libraries, the 386 layout and unrelated interfaces unchanged.

Two attack callbacks are transient radial-scan visitors. Use Go closures capturing
the original attack record/player; preserve geometry, iteration order, strict
positive-overlap dispatch and the returned center-address word. The two radial
fixtures retain independent C observers behind Go closures. No persistent callback
identity registry is needed for these routes.

## Baseline and coverage

Original source is the qualified monster conversion `456bdf2d`. All 565 affected
roots pass on default/server/highres: 558 results are reused from the immediately
preceding qualification with identical complete source/runtime fingerprints; seven
additional roots ran on the same verified binaries. The additional roots cover
radial observers, lifecycle initialization, session admission and transfer sound/
damage owners. Production evidence is reused from that identical-source checkpoint.
See [baseline](player-action-identities-baseline.json) and
[test selection](player-action-identities-tests.txt).

Selection follows all 35 wrappers, full fixture preambles, actual same-package
callers and private fixture routes. Resolve generic method names by receiver:
`portTestShopPools.run/snapshot` must not pull in unrelated helpers with those names.
The whole-source literal audit also covers the two function-address uses; a
call-only scan would miss them. Root-package functions with the same names are
independent native owners, not calls to the retiring legacy wrappers.

Existing independent radial contracts cover shapes, strict contact boundaries,
ordered callback arguments, nil centers, context bits and center return values.
Keep every assertion and captured expectation unchanged. The combined owner
selection includes player controls/attacks and connected shop/lifecycle fixtures.

## Conversion review requirements

Preserve C char/short sign extension before widening to capture words, double
result bits, pointer-shaped results, sparse operation IDs and byte narrowing.
In particular, warcry returns a truncated pointer through a signed short. Preserve
its existing normalization and the native bot-update override at operation 47.
Keep the independent attack-effect and player-init C observers and their ownership.

Retired callback addresses were also fixture normalization keys. Reserve their
32 entries in the existing capture-ID counter (12 attacks, 20 controls), rather
than registering nil pointers. Nil registration would change zero normalization.
Original export addresses are pairwise distinct on all three verified binaries;
unchanged frozen captures must additionally confirm dynamic ID preservation.

Original baseline committed/pushed as `44a9a065`. Primary verified the frozen
19-path overlay (including two deletions), formatted a reviewed copy, checked
functions through the Go AST and compared all owner/operation mappings and exact
32 prototype removals. No selected legacy symbol remains in the merged source.

Luna removed the complete controls identity getter/loop while its review claimed
native IDs 20/47 were retained. Primary restored both before compilation, removed
a duplicate standard-header include and added explicit candidate lifetime protection
around the two independent C observer calls. The other mappings matched review.
This was a useful bounded draft with a substantive correction; its prose was not
accepted as proof. Root assertions and frozen captures remain unchanged.

The map-start public adapter uses a zero-initialized Go point. Its existing owner
only writes that output and does not retain or expose its address; the previous
C-layout allocation is unnecessary. This reversible cleanup preserves nil-input
zero output and gameplay selection/RNG behavior.

## Qualification scope

Qualification includes all 565 affected roots on all profiles, safe/static checks, three production
builds/ABI checks, exact known-suite comparison, fresh preflight/final headless
save/load, original asset hashes and the measured C inventory. The conversion is fully qualified.


## Qualified result

All 565 affected roots pass on default/server/highres against exact original names,
without failures/skips. Safe/static, three production/ABI checks, exact known
asset-suite comparison, both fresh headless save/load scenarios and all 1,654
original asset hashes pass. See [qualification](player-action-identities-qualification.json).

Legacy exports fall 333→301; selected production cgo files fall to 150 client and
151 server. The 157 tracked headers contain 3,063 physical lines.
Embedded production C bodies remain 77; standalone production/test C remain zero.
External native bindings are unchanged. The last full default corpus was the
preceding shared drawable-update batch `024b2632`, explicitly on that earlier source.
