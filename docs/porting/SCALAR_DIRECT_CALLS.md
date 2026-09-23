# Matching-signature direct Go calls

This batch selects 56 production calls across 34 files whose 33 existing Go
export functions have exactly the parameter and result types used by the
corresponding generated cgo call wrappers. The conversion removes `C.` from
those expressions. It preserves argument evaluation, casts, narrowing, return
conversions, side effects and the existing implementations.

## Boundary review

The [signature inventory](scalar-direct-calls-signatures.json) records each
source site, Go declaration and generated C-call signature. The generated type
snapshot came from the qualified balance checkpoint (`1bb4b78a`); its SHA256 is
`6f4510a777e30f7139e1e166dd60d2f8e6dba8f6f4f722b38794761e4cd1d533`.
Grouped Go arguments were expanded before comparing types. C `_Bool` and void
were normalized to the corresponding generated types. Source review found no
macro aliases or local declarations shadowing the selected Go function names.
All 34 caller files still need their C imports.

“Scalar signature” describes the ABI types, not purity: several integers encode
existing object/window/team addresses and several calls mutate gameplay state.
Their conversions and object ownership remain unchanged. The supported target
is 386/SSE2; this batch does not claim new 64-bit support.

`sub_50B510` is excluded: its local C prototype returns `int`, but the existing
Go export returns no value. Its current caller discards the result. Resolve this
existing mismatch separately; it is not part of an exact-signature conversion.

## Unused exports

After the preceding book and balance call migrations, three book C entry points
and two balance C entry points have no production C callers or address references
in this repository. The only remaining balance C callers are baseline tests.
The retirement removes these five exports and their header declarations.
It preserves the public Go book function variables and the balance getter bodies.

The independent balance fixture retains all 104 mode/key/index cases and exact
expected float64 bits. The original C/direct comparison is preserved in Git and
its qualified checkpoint; retaining that test route does not justify retaining
dead production exports. Fresh safe/production binaries must omit all five
retired symbols. No supported external C/plugin interface was found for these
internal executable symbols; this is an in-tree compatibility decision to review
if an out-of-tree user is later identified.

## Qualification

The pre-conversion baseline passes 379 affected roots plus the supplemental
spell-class contract in each of default/server/highres, without skips. Source is
identical to qualified production checkpoint `1bb4b78a`; its four retained binary
hashes were verified for production-evidence reuse. See the
[baseline record](scalar-direct-calls-c-qualification.json).

After conversion all 380 selected roots pass in default/server/highres without
skips. The discovered test-name sets match the original baseline plus the direct
spell-class contract, allowing only the balance contract's documented rename.
All 104 numeric cases and their independent expected values remain. Existing
source-frozen capture hashes are unchanged. See the
[conversion record](scalar-direct-calls-qualification.json).

The selection includes AI combat/state, attack, balance, meters, damage,
equipment, gameplay reports, inventory, death/object state/objectives, controls,
projectiles, browser/options, spawn policy, spellbook/effects/lifecycle and world
motion. These are affected owner contracts, not exhaustive call-site or branch
coverage. Review of exact forwarding signatures complements those tests.

Safe/static checks, three fresh production builds, ABI checks, headless character
creation and save/load pass. All four binaries omit the five retired C exports
and the 33 redundant C-call bridges; required exports remain Go-backed and no
porttest symbols are present. The full suite matches the existing 304 failure
events (17 passing, two failing, 32 skipped packages). This is an exact known
baseline comparison, not an entirely green full suite.

## Delegation and recovery

Luna drafted the exact-list call patch and a bounded owner/test inventory. Primary
checked generated signatures, call expressions, ownership and the retirement
boundary. The owner inventory identified the omitted spell-class test. Its
family-name mappings are explicitly not proofs of per-call execution. Independent
diff review also identified the old balance-specific selector's stale test name;
that selector now names the retained Go numeric contract. The active broader
selection already matched it. Review found no accidental implementation or
expected-value changes.

Artifacts: `build/port-scalar-direct`; the source inventory and initial draft are
under `build/port-balance-direct`. Original assets and retained binaries remain.
Two completed cleanup runs archived 167 older successful logs losslessly:
1,823,895,814 raw bytes became 121,160,892 gzip bytes, verified before raw removal.
Restore instructions/hashes are `old-log-archive-record.json` and
`old-log-round2-record.json` under the new artifact directory.

Standalone production/reference C remains **zero files / zero lines**; the
remaining **79 C preamble bodies** are unchanged by this batch. Generated
bridges and external libraries are outside that metric. No measured game-speed
claim is made.

A follow-up inventory distinguishes Go C-selectors from calls and address uses in
C preambles. Several test drivers still invoke the converted scalar exports from
C blocks, and two production callback addresses remain live. Do not infer an
unused export solely from absence of `C.name` selectors. Ten scalar exports have
only C prototypes remaining and are candidates for the next reviewed cleanup;
the others need their test-driver or callback dependencies handled explicitly.
