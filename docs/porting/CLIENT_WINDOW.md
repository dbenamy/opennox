# Window geometry, state and draw-data helpers

Status: native conversion fully qualified against corrected C baseline
`20970e2a`, following listbox conversion `5cba7b22`.

Scope: all196 original lines in `client__gui__window.c`, plus164 lines in
GAME2_1.c from0046A8A0 through before0046DB80:26 routines /360 original C lines.
Keep the23-line `gui_cgo.c` function-pointer adapters; they are ABI adapters,
not remaining algorithms to remove for a size tally.

Fixtures use actual GUI window/tree owners, callbacks and renderer references.
The font helper registers a face and an opaque handle through the same renderer
maps used by loaded fonts, with restoration. Capture complete window/draw-data
words, geometry outputs and callback-time state, tree ordering and returned
references; normalize only identified owned pointers.

## Prerequisite range termination

Both inclusive ID-range loops increment a signed int after visiting the last ID.
With first=last=INT_MAX, that wraps under the qualified compiler options and the
loop continues. Break after processing the inclusive final ID. Two local C lines
preserve all ordinary ranges and terminate the maximum-ID single-child case.
The original nonterminating case is established by source audit, not executed
as a golden. A contract checks hide, disable and show of an actual child whose
ID isINT_MAX. This is a reversible correction to review later.

Other independent contracts check distinct nil-output behavior, updated geometry
before resize notification, inclusive point bounds and strict ancestry. Matrices
cover signed geometry/bit masks, missing/nil windows, hidden and overlapping
children, callback returns, label-style precedence, image/font handles, exact
332-byte draw copies and inclusive ID ranges.

The first C fixture build required explicit image-handle-to-pointer conversions;
that was corrected before freezing. Capture/compiler evidence is under
`build/port-client-window`. Actual source supersedes historical ignored drafts.

Caller audit retired five interfaces: unused reverse-tab no-op, font getter,
draw-data copier, deepest child lookup and ancestor-visibility query. Their Go
callers now use Go directly. Twenty-one interfaces have remaining C callers.
Do not replace these with higher-level GUI methods without checking semantics:
raw flags/dead-window access, nil outputs and signed/unsigned inclusive bounds
are deliberately distinct in places.

Corrected-baseline C count is88,574 /100 files /zero reference C (+2 prerequisite
lines since listbox). The native conversion removes all362 corrected lines.


The corrected C capture contains7,664 results in two groups: geometry/state5,040
and tree/labels2,624. It includes overflowing parent coordinates and disabled
children during hit testing. The repeat matched both frozen SHA-256 digests;
all six focused tests passed in22.887s. Broader C qualification is recorded below.
The earlier capture passes (`c-b`/`c-c`) failed only because expectations were
intentionally unset; all four contracts passed in both executing runs.


## C baseline qualification

All six focused tests and both capture groups repeated unchanged. Affected client
123 /59.524s, server122 /173.602s and
highres123 /74.692s passed, with all selected tests run and finished.
The client build passed in58.463s and fresh warrior gameplay in
35.401s with null audio and reference override=false. The full asset
suite preserved the exact1,553 failure entries and15pass/3fail/32skip package
outcomes. All1488 Go/C/header source fingerprints remained unchanged.
Artifacts are under `build/port-client-window/c-final-*`, `c-qualification.json`,
`c-gameplay-qualification.json`, `c-full-suite-*`, and `c-source-verification.json`.

## Native conversion and qualification

All26 routines are native. The first executing native run passed all seven
focused tests in167.229s. Both frozen captures (7,664 results /two groups)
matched the corrected C captures byte-for-byte, without golden changes or
behavioral corrections after comparison. Five independent contracts pass.

Coordinate ABI getters preserve the original read/store ordering, including
aliased output pointers. The additional source-derived alias contract expects19
when both outputs share storage; this contract is derived from the original C,
not presented as an extra frozen C execution. Raw owned-window access, nil-output
differences, callback timing, strict ancestry and inclusive signed/unsigned hit
bounds remain explicit. Twenty-one interfaces retain actual C callers; five
unused/private interfaces are retired. Listbox, input and legal-screen Go callers
use native helpers directly. The23-line gui_cgo.c function-pointer adapters remain.

Accumulated default757 /418.355s, affected server123 /174.243s and highres124 /73.836s passed; every selected root test started and finished.

All three production binaries passed build and ABI checks (ELF32/i386, SSE2,
CGO,21 retained interfaces, five retired symbols absent, no test helpers). The full
asset suite preserved exactly1,553 known failure entries and15pass/3fail/32skip
package outcomes; it is not a green legacy suite. Fresh warrior gameplay passed
in37.506s and the existing entry-typing scenario in12.149s, with null audio and
reference override=false. All1489 source fingerprints remained unchanged.

Production C is**88,212 physical lines /99 files /zero reference C**, down362
from the corrected baseline (net360 after the two-line prerequisite). Accumulated
frozen coverage is424,490 results /1,115 groups, plus independent contracts.
Evidence: qualification.json, native-*-result.json, native-capture-comparison.json,
native-source-verification.json and the two client-window gameplay run folders.
Next connected candidate: localized item-hover text and cursor tooltip storage.
