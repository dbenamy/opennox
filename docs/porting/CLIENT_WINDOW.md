# Window geometry, state and draw-data helpers

Status: corrected C baseline qualified, after listbox conversion
`5cba7b22` was fully qualified and pushed.

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

The first fixture build required explicit image-handle-to-pointer conversions;
that is corrected. The baseline is now frozen; no native qualification is claimed yet. Capture/compiler evidence and ignored drafts: `build/port-client-window`.
Go drafts prepared during listbox qualification must still be reviewed against
this batch's actual baseline before integration.

Caller audit can retire five interfaces: unused reverse-tab no-op, font getter,
draw-data copier, deepest child lookup and ancestor-visibility query. Their Go
callers should use Go directly. Twenty-one interfaces have remaining C callers.
Do not replace these with higher-level GUI methods without checking semantics:
raw flags/dead-window access, nil outputs and signed/unsigned inclusive bounds
are deliberately distinct in places.

Corrected-baseline C count is88,574 /100 files /zero reference C (+2 prerequisite
lines since listbox). The full corrected scope will remove362 lines when ported.


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

The native draft preserves output-store ordering for callers that alias coordinate
output storage; computing a complete point first would differ from C in that case.
An additional source-derived ABI contract will check this during native comparison.
The draft remains unqualified until it passes the frozen captures and broader
checks. Go copies of historical fixture drafts must not replace current source.
