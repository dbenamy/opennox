# Shop UI

Status: **C baseline development; no shop conversion or prerequisite applied**.
Quantity/trade UI conversion de1bb918 is qualified and pushed. Production remains
**81,351 C lines / 95 files / zero reference C**.

Scope: 42 routines, GAME2_2.c 00478030 through 004798A0 and the whole shop
translation unit. The unmarked sub_479680 is a separate callback. Exclude the
following 00479950 subsystem. There are 1,059 address/function-block lines plus
29 TU preamble lines: **1,088 removable C lines** before prerequisites.

Reuse inventoryWindowOwner and the actual GUI parser/widgets, renderer, player,
string manager and drawable factory/pool. Optional extra fixture type names add
the eight shopkeeper variants without changing existing callers' type lists.
The shop owns 60 cells of 140 bytes, 32 item codes per cell, column-major layout
with row-first scan order, image/text/viewport state and actual shop widgets.
Construct the real dialogue owner explicitly to isolate shop-close playback state.
Shop actions queue/clear dialogue filenames; this fixture does not pump audio
streams. Headless gameplay separately exercises normal audio initialization.

Initial four tests exercise construction, exact lookup across cells/code slots,
inclusive/clamped hit regions at seven scroll offsets, and selection when the
first/all stacks are full. Development-a failed fixture setup in 181.051s: the minimal environment has no
initialized dialogue owner to copy. Only construction started; the driver correctly
reported the other three roots missing. All jobs joined before replacing that
assumption with the actual dialogue constructor. Development-b now runs the same
four contracts. No C changed and expectations are not frozen.
The suspected capacity defect is tested without appending beyond the C array.
Further scope includes start/close/reset lifetime, add/remove/modifiers/pricing,
scroll/buttons/mouse/tooltip/draw, buy/sell/repair/cancel request bytes and actual
quantity callbacks, allocation failures, resource failures and shipped windows.

The read-only caller audit predicts twelve remaining C command/report entry
points, the raw shop tooltip and five quantity callbacks: approximately eighteen
retained and twenty-four retired interfaces. Move native Go callers directly and
verify actual production symbols after translation. Audit evidence and original
C blocks are retained under build/port-client-shop-ui. No source generator is a
substitute for reviewing the final implementation and frozen expectations.

## Capacity prerequisite

Development-b completed all four roots in 21.849s. Construction and code lookup
passed, as did **437,675 independent hit/inside/scroll-coordinate contracts**.
All three full-stack selection assertions failed: the first full stack won over
an empty/partial cell, and an all-full shop still selected it. No overflowing
append was needed to establish the defect.

Under standing authorization, skip stacks with Count>=32 in lookup and guard
the final append. This small reversible bound correction adds **four C lines**:
working production **81,355 lines / 95 files / zero reference C**. Development-c
adds three tests for filling all 1,920 legal stock slots through the real add
owner, preserving cell state on allocation failure, and sixty-drawable reset
ordering/exactly-once deletion/idempotence. No captures are frozen yet.

Development-c passed all **seven roots** in 181.581s, including all 1,920 legal
adds, the rejected extra item, allocation failure and reset ownership/order.
Development-d adds three roots for the actual information message box (with an
authored resource in a temporary working directory), failed sell/repair quantity
allocation and buy-price multiplication overflow. No additional C correction has
been made before observing those contracts.

## Quantity and price prerequisites

Development-d completed all ten roots in 21.931s, with three failures. The real
sell/repair paths retained their pending flags after quantity drawable allocation
failed. Prices 0x08000000 and 0x80000000 with count 32 and zero player gold
wrapped their product to zero and incorrectly opened a quantity dialog. Zero
price correctly offered the full stack; other tested positive prices rejected it.
The separate information-dialog failure was a fixture capture omission: its
actual entry widgets were unsupported by inventoryWindowResult. Add their data
to that shared capture; existing fixture windows/expectations are unchanged.

Compare affordable quantity as gold/unit-price when price is nonzero, preserving
free goods without multiplication overflow. For sell/repair, record the previous
quantity drawable and arm the pending flag only if successful show publishes a
replacement. A failed allocation preserves the previous dialog and does not arm
a new request. These reversible local corrections add seven more C lines:
**eleven prerequisite lines total; 81,362 working production C lines**.
Development-e reruns all ten roots and adds independent successful retry/cancel
checks after the real allocation pool recovers. No expectation has been frozen.

Development-e passed all **ten roots** in 183.319s, including retry and actual
cancellation callbacks after allocation recovers. Development-f completed fourteen
roots in 24.150s. Its one failure was a test interpretation error: sub_478080
returns the active shop drawable, not its price. Correct the assertion/name and
retain zero-code/full-code-slot lookup behavior. The exact 192 request-byte cases,
full-inventory rejection and removal ordering/lifetime cases passed.
Development-g adds mode/scroll/panel/mouse/resource and modifier/selection tests
for 22 selected roots. The entry-widget capture now normalizes its identified
IME child-window pointer; existing non-IME captures keep the same zero word.

Development-g passed all **22 roots** in 29.804s. Development-h adds start/close,
shopkeeper pictures, 100 destroy/recreate cycles, draw/hover and real quantity
button flows for 29 roots. The capture now records GUI-render state, owned
message-box state, dialogue filename changes and identified string pointers.
This is baseline development, so no frozen expectation changes. Existing
inventory/trade test schemas and expectations remain intact.

## Lifetime repeatability and closing during quantity input

Development-h passed **29 roots** in 37.002s. Development-i passed **30 roots**
in 39.714s, including shipped Shop.wnd/MultMove.wnd with authored reports and
actual quantity callbacks: **1,291 captured results / 24 groups**. Twenty-two of
23 common H/I groups matched byte for byte. Reset differed only in the cached
slider/up/down widget addresses after destruction, across two rows (four JSON
paths). Retain those known live widget identities before invocation so captures
can normalize cached addresses after destruction without dereferencing them.
No C behavior is changed by this capture correction; retain h-i-reset-ownership-diff.json.

Development-j completed **31 roots** in 37.549s. Its new closure contract proved
that both normal shop closure and destruction retained owned buy/sell/repair
quantity dialogs and their drawable; some pending flags also survived. Explicit
cases with an unrelated inventory quantity dialog establish the ownership boundary.
Cancel only when the quantity accept callback is one of the three actual shop
callbacks, before teardown; clear both shop pending flags. This adds eleven C
lines: **22 prerequisite lines total; 81,373 production C lines / 95 files**.
Development-k reruns all 31 roots. Original failing captures/source remain local.

All three production builds passed before the close prerequisite (58.417s /
8.535s / 51.363s); retained in pre-close-bin. New production builds are running
for the corrected behavior. A bounded shop-scene setup attempt uses the prior
binary only for development; it does not qualify the new source.

## Qualified C baseline

Development-k passed all 31 roots in 202.260s: 1,299 captured results / 25 groups,
plus independent coordinate, capacity and lifetime contracts. All 24 groups
shared with development-j are unchanged after the targeted close correction.
Affected default/server/highres passed 379 / 377 / 379 roots in 199.343s /
295.693s / 223.351s. All shop hashes repeat exactly in every target, and all
23 frozen quantity/trade groups remain unchanged. Production builds passed in
59.715s / 10.501s / 70.294s; each is ELF32/i386/SSE2/CGO with all 42 shop
interfaces and no test helpers. Full assets retain exactly 1,553 known failure
entries, with 15 passing / 3 failing / 32 skipped packages; no added failures.
Unchanged nine-screen inventory gameplay passes from fresh assets/save.
Source fingerprints stayed unchanged throughout qualification.

The new tracked [shop gameplay scenario](shop-ui.yaml) also matches all ten
full reference screenshots in an independent fresh process. Run with NOX_DEV=true,
seeded RNG, Xvfb 1280x960 and OpenAL null, reference updates disabled. Ordinary
console key events create a generic Shopkeeper, then normal interaction opens
its shop, switches buy/sell/repair modes, opens a priced pants sell-quantity
dialog, cancels it and closes the shop. The generic NPC has empty stock and no
map-authored localized name. The historically named repair_quantity checkpoint
shows repair inspection, not a repair quantity dialog. Stocked purchase and
repair quantity callbacks are instead covered by the actual-owner fixtures and
shipped Shop.wnd/MultMove.wnd tests. This is not complete campaign/shop coverage.
Initial scene setup attempts remain development artifacts, not qualification.

Evidence: build/port-client-shop-ui/c-qualification.json, c-*-result.json,
c-qualified-capture-compression.json, frozen-captures.json; gameplay runs
client-shop-ui-c-baseline and client-shop-scene-c-repeat against
client-shop-scene-quantity. Completed development captures are losslessly
compressed with byte-hash manifests; original assets remain untouched.

Exactly 1,110 C lines are now removable across the 42 routines and the shop
translation-unit prelude (1,088 original + 22 justified prerequisites).
Working production C: 81,373 lines / 95 files / zero test-reference C.
The 25 expectation literals are frozen; the final locked repeat passed all 31
roots in 37.984s. No native shop source has been applied at this baseline.
