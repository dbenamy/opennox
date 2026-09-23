# Matching-signature direct Go calls

This batch selects 56 production calls across 34 files whose 33 existing Go
export functions have exactly the parameter and result types used by the
corresponding generated cgo call wrappers. The planned change removes `C.` from
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
The planned retirement removes these five exports and their header declarations.
It preserves the public Go book function variables and the balance getter bodies.

The independent balance fixture retains all 104 mode/key/index cases and exact
expected float64 bits. The original C/direct comparison is preserved in Git and
its qualified checkpoint; retaining that test route does not justify retaining
dead production exports. Fresh safe/production binaries must omit all five
retired symbols. No supported external C/plugin interface was found for these
internal executable symbols; this is an in-tree compatibility decision to review
if an out-of-tree user is later identified.

## Baseline and qualification plan

The pre-conversion baseline passes 379 affected roots plus the supplemental
spell-class contract in each of default/server/highres, without skips. Source is
identical to qualified production checkpoint `1bb4b78a`; its four retained binary
hashes were verified for production-evidence reuse. See the
[baseline record](scalar-direct-calls-c-qualification.json).

Freeze the existing affected gameplay/UI contracts before changing production
source. Run default/server/highres with every discovered selected root required
to complete, without skips. A supplemental baseline includes
`TestSpellClassEligibilityABI`, which the initial owner-family pattern omitted.
After conversion include it in the main selection and preserve the numeric/book
contracts. Existing frozen hashes and expected values must remain unchanged.

The selection includes AI combat/state, attack, balance, meters, damage,
equipment, gameplay reports, inventory, death/object state/objectives, controls,
projectiles, browser/options, spawn policy, spellbook/effects/lifecycle and world
motion. These are affected owner contracts, not exhaustive call-site or branch
coverage. Review of exact forwarding signatures complements those tests.

Then qualify safe/static, three fresh production builds, retained/retired ABI,
exact known full-suite results, headless character creation and save/load. The
known full-suite baseline remains 304 failure events (17 passing, two failing,
32 skipped packages); do not reinterpret it as an entirely green suite.

## Delegation and recovery

Luna drafted the exact-list call patch and a bounded owner/test inventory. Primary
checked generated signatures, call expressions, ownership and the retirement
boundary. The owner inventory identified the omitted spell-class test. Its
family-name mappings are explicitly not proofs of per-call execution.

Artifacts: `build/port-scalar-direct`; the source inventory and initial draft are
under `build/port-balance-direct`. Original assets and retained binaries remain.
Two completed cleanup runs archived 167 older successful logs losslessly:
1,823,895,814 raw bytes became 121,160,892 gzip bytes, verified before raw removal.
Restore instructions/hashes are `old-log-archive-record.json` and
`old-log-round2-record.json` under the new artifact directory.

Standalone production/reference C remains **zero files / zero lines**; the
remaining **79 C preamble bodies** are unchanged by this planned batch. Generated
bridges and external libraries are outside that metric. No measured game-speed
claim is made.
