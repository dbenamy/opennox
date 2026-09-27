# Animation callback identities

## Scope and original behavior

This batch replaces 25 engine C callback addresses used by GUI animation
transitions with stable native identities and typed Go dispatch. The connected
owners are main menu, character selection/creation, options, bindings and server
browser; game-state transitions also install these callbacks.

Preserve the 68-byte animation record and callback offsets, signed integer
results, dynamic lookup of replaceable hooks, and the existing nil checks.
Animation completion updates/clamps position and state before invoking its
callback. List ticking saves the next node before a callback may free itself.
Several completion owners save a callback pointer, free their animation and
windows, then invoke it; identity storage must be independent of that allocation.
Unknown foreign callbacks retain the existing raw C fallback. The incoming
completion slot has no non-null production assignment in the complete field
inventory; its existing void dispatch remains unchanged.

The primary independently checked the 25-symbol/11-getter reference closure,
including saved aliases and direct Go calls. Shared window-event callbacks and
C record definitions are outside this scope. Native wrappers may remain for
live Go callers; removing an export does not require removing its implementation.

## Baseline and qualification plan

Two new porttest-only Go files add five original-path contracts: signed return
boundaries and hook replacement; foreign callback fallback; completion position/
state before self-free; safe linked-list iteration across self-free callbacks;
and the real main-menu completion owner invoking a copied callback after both
animations and their windows are cleaned up. The animation-list fixture owns an
isolated list and restores prior state. These contracts passed a focused original
preflight before the broad baseline.

The affected selection includes 68 roots across animation, character creation,
options, binding editor, server browser and main-menu resource drawing. Each
selection passed twice independently in default, server and highres on original
source, with exact expected root names and no failures/skips.
Existing state/pixel captures remain unchanged. The preceding GUI-adapter
production qualification can be reused for the original baseline because the
only source changes are the two porttest-only files; no affected root-test result
is reused across that test-source change.

After conversion, require the full affected selection in all three profiles,
a fresh default headless preview, the complete accumulated default port corpus,
safe/static checks, three production builds/ABI checks, the established full-suite
result and final headless save/load/resume, plus original asset hashes. The full
corpus is warranted here because the animation dispatcher is shared GUI
infrastructure. Its independently derived expected set has 2,475 roots, including
the single established prerequisite-probe skip.

## Delegation and reversible decisions

One GPT-6 Luna helper audited the field/getter closure and reviewed the five new
contracts. Primary review accepted the test ownership and cleanup, corrected
ambiguous audit wording about direct Go callers, and owns API/baseline acceptance.
Luna prepared the bounded implementation in an isolated ignored draft;
tracked source stays frozen during tests. Final acceptance requires primary
mapping/caller review and actual qualification.

Use a typed registry following the existing tooltip callback pattern. Stable
keys are process-lifetime byte-array addresses; registration closures look up
mutable hooks/owners at call time. Preserve raw fallback, including existing
unguarded nil behavior. Do not modernize allocation/layout in this batch.

Before original profile builds, nine obsolete, hash/stat-verified project cache
archives were removed after host process/open-file checks, reclaiming 420,257,792
allocated bytes. Source/assets/current binaries and evidence were retained.
Journal: `build/port-animation-identities/cache-before-baseline/`.

Local artifacts: `build/port-animation-identities/`. Source inventories and draft
reviews: `build/port-after-gui-scout/`. The original baseline is recoverable at `c6ffcb75`; converted results follow.

Committed evidence: [original baseline](animation-identities-baseline.json),
[batch manifest](animation-identities-batch.json), and
[affected selection](animation-identities-tests.txt).

## Qualified conversion

All 68 affected roots pass in default/server/highres with exact original names.
The full default corpus passes 2,474 roots with its one established prerequisite
skip (2,475 total). Safe/static, three production builds/ABI checks, exact known
suite comparison (304 failure events; 17 passing, two failing, 32 skipped packages),
preview/final headless save/load/resume and all 1,654 original asset hashes qualify
unchanged reviewed source. Frozen captures and baseline tests are unchanged.
See [qualification](animation-identities-qualification.json).

Exports fall 169→144; selected production cgo files fall 119/120→113/114.
Headers remain 157 files / 2,902 physical lines. Embedded production C bodies remain 77;
standalone production/test-reference C remain 0. Process-lifetime native keys route
25 owners; generic foreign fallback stays until its other users are migrated.

Luna's first draft preserved routing but unnecessarily retained 23 orphan wrappers,
including C-typed wrappers. The requested refinement hit Luna's usage limit before
being produced. Primary completed it locally: removed 23 orphan wrappers, retained
the real drawGeneral body and a native-int browser wrapper with direct Go callers,
and eliminated 6 unneeded cgo imports (two empty files deleted). Primary verified all
25 mappings, exact export/prototype delta, every address assignment/getter and six
saved-pointer calls. No replacement helper model was used. This was useful bounded
delegation with a material primary scope correction; no measured cost claim.
