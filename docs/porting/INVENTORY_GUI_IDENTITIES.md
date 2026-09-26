# Inventory GUI bridge conversion

## Scope

Retire 22 remaining inventory display, query, request, state and tooltip exports.
Migrate five connected private C-typed helpers (`sub_463420`, `sub_467420`,
`sub_467430`, `sub_467930`, `sub_467750`) with explicit signed32/signed8/byte/pointer types.
Use existing native `drawableUnitCode` for the two local request calculations;
its C export remains available to unrelated callers. Preserve algorithm bodies,
message bytes, return words, mapped layout and fixture operation IDs.

Five tooltip callbacks use the GUI identity registry established in qualified UI
meter revision `541585fa`. Three draw callbacks use existing typed GUI functions.
Keep their fixture normalization indices unchanged. The selected query and state
functions retain their native owners and existing signed/pointer result behavior.

Leave `sub_465CD0` and its declaration in place: the amount dialog stores that
address and invokes it later. Remove that boundary with its owning dialog in a
subsequent batch, including lifecycle and foreign-fallback coverage.

## Baseline and qualification

An independent original-route contract uses the actual
inventory constructor and all five installed tooltip fields. Retain every prior
state/pixel expectation. The affected selection includes the preceding UI meter
coverage and player-file serialization tests for the inventory query consumers.

All 394 client / 391 server original roots pass without failures/skips, and the new
installed-route contract passes twice per profile in separate processes.
[Accepted baseline](inventory-gui-identities-baseline.json) records exact names,
source and runtime environment. The only added source is this test file;
production matches qualified UI meter `541585fa`. No inventory conversion is
installed or qualified yet. Full converted acceptance
requires affected roots in three profiles, safe/static, three production/ABI builds,
exact known-suite comparison, two headless save/load scenarios and asset hashes.

## Delegation

Luna supplied a bounded caller/fixture audit; primary checked the delayed callback
and selected the scope. Luna is drafting one installed-tooltip contract for primary
review. Source review corrected the draft status-tooltip boundary (x39 selects the
first effect; x40 advances to the next slot). No measured usage savings are claimed.

The initial new-test compile caught a 32-bit constant expression overflow in a
normalization-ID assertion. Primary corrected the test before any execution;
production was unchanged, and failed compile evidence remains in
`build/port-inventory-gui-identities/original-first-compile/`.

The first full original run found only a new-test expectation error: x313/y13
is inside the special-item tray and returns its Gold tooltip. Primary checked the
embedded rectangles and added exact named-item and outside-edge assertions. The
corrected test passed a focused 386 preflight, the complete baseline, and separate
repeat processes in all profiles. Existing production source and expectations were
unchanged throughout. Failed run records and exact test source are retained.
