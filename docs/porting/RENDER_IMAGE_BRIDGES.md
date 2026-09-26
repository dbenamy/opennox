# Rendering and image bridge retirement

## Scope and behavior

Retire 17 unused exports: 16 rendering wrappers in `legacy/draw.go` and the
image-loading wrapper in `video_bag.go`. Independent whole-source review found
69 references, all definitions/directives, header prototypes, Go owner comments
or two translator keys. Preserve the translator metadata and native renderer
implementations. Remove the exact prototypes alongside their wrappers.

Remove the unused Go viewport conversion chain and image-reference C alias.
Keep the shared C records: viewport observers still use the viewport layout,
and the C bag layout still contains the opaque image pointer type.

The live private image interfaces move to existing `noxrender.ImageHandle`:
`asImage`, `drawImageAt`, the meter image caller and the ability-icon forwarder.
Keep the two ordered `GetClient` calls in image drawing: evaluate the draw receiver
before resolving the image through the second lookup. Preserve ability-hook calls,
32-bit handle words, allocation and lifetimes. Book/quickbar callers already have
compatible explicit pointer conversions. No rendering backend or algorithm changes.

## Baseline and coverage

The original source is qualified window revision `3838462b`. Select 92 client /
89 server roots by actual callers rather than carrying forward the whole 401/398
window selection. Reuse 88 exact-source roots per profile; run four additional
client roots (main menu and three status-overlay roots), or one server root
(main menu), twice in separate processes on the qualified binaries. The first
baseline collector assumed the overlay roots also existed in server; its count
assertion caught their explicit `!server` tags after a successful server run.
Completed results were preserved and validated; only outstanding runs resumed.

Coverage includes complete existing meter, spellbook and quickbar root families,
plus image-drawing inventory, shop/trade, amount-dialog, main-menu and status-overlay
contracts. Frozen pixels/state and independent assertions remain unchanged.
The helper's narrower suggestion omitted some live image consumers; primary added
those owners and retained the broader three related families. Summon/scoreboard
use `bookDrawImage`, whose separate implementation is unchanged; those extra
unaffected families are not required for this interface change.

The meter's conditional main-image branch is not directly exercised by the
existing meter matrix. Its condition and caller are unchanged; the same image
helper is exercised with actual image handles by other selected owners. Preserve
this limitation rather than claiming full branch coverage.

The renderer package is not all green: `TestDrawImage` has established failures;
particle, pixel-hash, circle and line tests pass in the recorded suite. Require
exact known failure/package outcomes, alongside all selected contracts, safe/static,
three production/ABI builds, two headless save/load runs and original asset hashes.
The accepted baseline and qualification records are linked below.

## Delegation

One GPT-6 Luna helper audited references and test scope and drafted the bounded
conversion. Primary independently checks reachability, native pointer interfaces,
evaluation order, test selection and all acceptance evidence. No usage savings
are inferred. Primary accepted the source draft without code corrections and
corrected a stale sentence in its formatting report. The converted results below
are accepted.

## Qualified conversion

All 92 client / 89 server roots pass with exact expected names and no skips or
failures. Safe/static, three production/ABI builds, exact known-suite comparison,
two headless save/load scenarios and unchanged hashes for all 1,654 original
assets qualify the same reviewed source. Root tests and frozen captures are unchanged.
See [qualification](render-image-bridges-qualification.json) and
[baseline](render-image-bridges-baseline.json).

Exports fall 191→174; selected production cgo files fall 128/129→126/127.
Headers remain 157 files / 2,930 physical lines. Embedded production C
bodies remain 77; standalone production and test-reference C remain zero.
