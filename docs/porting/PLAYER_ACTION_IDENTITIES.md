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

Luna owns a bounded uninstalled overlay. Primary owns original baseline acceptance,
review against wrappers and full C preambles, integration and qualification.

## Planned gates

Run all 565 affected roots on all profiles, safe/static checks, three production
builds/ABI checks, exact known-suite comparison, fresh preflight/final headless
save/load, original asset hashes and the measured C inventory. Conversion is pending.
