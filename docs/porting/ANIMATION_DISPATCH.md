# Native animation dispatch

Status: original baseline accepted; native draft not installed.

Route all four animation callback slots through the existing native registry,
discarding the result for completion slots. Retire CallIntVoid/CallVoidVoid and
their generator entries. Replace the options fixture's Go-to-C-to-Go counter and
a new original-path completion observer with stable registered Go identities.
Keep the 25 production handlers, mutable-hook lookups, signed returns, identity
stability, nil completion guards and the 68-byte Anim representation unchanged.

The producer audit resolves all 33 production field writes to nil or one of the
25 native registrations, including getters and aliases. Existing non-registered
keys occur only in the options fixture, consumed by animation, options and
character lifecycle assertions. The newly added completion observer is likewise
a test-only C export. Whole-source constructor/field, consumer and C-layout checks
find no production foreign animation-key producer. The in-completion field remains
supported through native identities; it is not discarded because current game
constructors initialize it to nil.

Unregistered animation keys now deliberately panic rather than entering raw C.
This ends arbitrary foreign-pointer support at this internal animation boundary;
no engine producer depends on it. Recoverability/diagnostics for invalid keys are
a reversible behavior correction for review, not claimed identical C-crash output.
Other callback families retain their existing behavior pending their own audits.
External library bindings are unaffected.

## Original contract and fixture correction

TestAnimationInCompletionOrdering checks both movement directions, nil/non-nil
completion, callback-driven freeing, final state/position before callback, focus
after callback, list unlink and one-shot completion. Its six cases pass twice in
each of default/server/highres/safe on the original C route. The first draft used
newEntryOwner: default passed twice, but safe rejected an unrelated effects fixture
access to mapped global dword_5d4594_1313532 before reaching animation assertions.
The corrected fixture owns a real gui.GUI and real windows directly; no drawing or
entry-widget environment is needed. Repeated all four profiles after correction.
Initial source/logs and the correction are preserved in the batch artifacts.

## Selection and acceptance

Captured 69 owner roots in 52 assertion files covering animation, character
creation, options, bindings, server browser and main-menu resources. Use the exact
same compiled-name sets after conversion. The registered-key dispatch branch and
all handler mappings are unchanged; the changed fallback has only controlled
fixture consumers, while in-completion behavior has the new independent contract.
This selection does not claim to requalify every renderer/game-loop implementation
called by an unchanged handler.

A package/name inverse scan conservatively selected 2,485 roots. Giving init
functions separate identities reduced it to 1,677, exposing false links through
merged init/local names. Method/field ambiguities still over-select unrelated
subsystems. The selected scope instead follows audited producers, dispatch owners
and the changed branch; the graph is an audit aid rather than proof of test impact.
No root-baseline reuse is allowed after adding the new fixture: recapture the
chosen original selection. Existing frozen expectations remain unchanged.

Require focused native completion/static checks, matching three-profile owner
roots, the new contract twice in four profiles, safe build/static, three production
builds/ABI, exact known-suite outcomes, fresh save/load and asset integrity.
Check both retired dispatcher symbol substrings and both fixture export names are
absent where applicable; retained dispatcher bodies and all 25 mappings must match
baseline source. Metadata discovery is not build evidence.

Expected qualified counts: 5 client / 6 server production cgo files, 15 embedded
C bodies, zero legacy exports and zero standalone C. Fixture C-import source files
are 101 at the parent, 102 with the temporary baseline observer, and 100 after
conversion. Headers remain 157 files / 2,731 lines.

Primary only; Luna remains quota-unavailable. Artifacts:
`build/port-animation-dispatch/`. Installation/acceptance scripts are single-use.

All 69 original roots passed without skips in all three profiles; the new
completion contract passed twice in each of four profiles. Native conversion and
qualification remain pending. Original static memory checks passed.
