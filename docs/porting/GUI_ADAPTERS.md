# GUI adapter retirement

## Scope

Retire five unused C wrappers: static-text creation and drawing in `gui_widgets.go`,
font lookup in `gui_fonts.go`, animation creation in `gui_anim.go`, and dialog
creation in `gui_dialog.go`. Independent whole-source review found 20 references,
all definitions/directives, prototypes, native-owner comments or translator metadata.
Keep the native Go owners, hook variables and translator mappings.

Remove the unused `asGUIAnim` converter and the `asWindowData`/`asWindowDataP`
chain used only by the retiring static-text wrappers. Preserve shared C records
and existing Go layout assertions. The private legacy RGB helper can return Go
`uint32` directly; its only caller already consumes that word. Keep byte narrowing,
color conversion, the separate native color setter and the unrelated root-package
function with the same legacy name.

This retires cgo imports from five GUI files. Separately remove leftover
includes-only preambles/imports from `client_session.go` and
`server_orchestration.go`: neither has a C selector, export, function body or
`#cgo` directive. Their Go code stays unchanged. Retain compiler-setting files
and live animation callback-address exports; no callback protocol or native backend
replacement is included.

The unused dialog wrapper contains raw function-pointer adapters. With no in-tree
caller or address registration, it follows the established unused internal-export
retirement policy. Native dialog hooks and their live callbacks remain unchanged;
no supported external plugin ABI has been identified for this engine entrypoint.

## Original evidence and qualification

Use the exact qualified rendering/image source at `46f07aba`. Reuse its 26 meter
roots, which exercise the only live color-helper consumer. Add ten existing roots
covering animation cleanup, dialog preview/connections/notices, rename events,
session-quit callbacks, briefing construction/drawing and scoreboard headings/pixels.
Run those ten twice per profile in separate processes on the same qualified binaries.
Require all 36 roots per profile after conversion, with exact names and no skips.

No fixture or frozen expectation changes are planned. Native owner tests establish
regression coverage; the wrappers themselves are unreachable. Also require safe/
static checks, three production/ABI builds, exact known-suite outcomes, two headless
save/load scenarios and unchanged original asset hashes.

## Delegation and review

One GPT-6 Luna helper supplied a bounded reference inventory and ignored source
draft. Primary independently checked the five-export reference sets and converter/
color call graph. Primary corrected two scout descriptions: `asGUIAnim` was already
unused, and the separate color setter has no C signature. The two includes-only
imports were identified separately from the GUI behavior change. Source review
and qualification remain primary-owned; no usage savings are inferred.

The conversion below is qualified.

## Qualified conversion

All 36 roots pass in each profile with exact original names and no failures/skips.
Safe/static, three production/ABI builds, exact known-suite comparison, two headless
save/load scenarios and all 1,654 original asset hashes qualify unchanged reviewed
source. Root tests and frozen captures remain unchanged. See
[qualification](gui-adapters-qualification.json) and [baseline](gui-adapters-baseline.json).

Exports fall 174→169; selected production cgo files fall 126/127→119/120.
Headers remain 157 files / 2,924 physical lines: five prototypes occupied six lines.
Embedded production C bodies remain 77; standalone production and test-reference C
remain zero. Primary accepted the source draft without code corrections and
independently reconstructed both headers and the two includes-only source edits.
