# Quantity dialog and player-to-player trade UI

Status: **Go conversion fully qualified**. C baseline
b77d5102 is committed and pushed. Prior inventory
window conversion 24f3f67e and cancellation cleanup e712ca56 are pushed.

Scope: 36 routines, 1,264 address-block C lines plus the trade translation unit's
17-line preamble: about 1,281 removable C lines. The adjacent sub_4C1CA0 belongs to
summon controls and is excluded. Before prerequisites, production was 82,632 C lines / 96 files /
zero reference C.

The initial C prerequisite corrects sub_4C0910/sub_4C11E0 from char to a complete
cell-pointer return, replaces partial-pointer assignments in callers, and makes
outside-both-grid hover avoid dereferencing packed coordinates. These are
reversible decompilation corrections under standing authorization. Independent
contracts cover exact cell addresses at 42,436 grid/edge/outside positions across
two window offsets, plus empty/outside hover. C LOC is unchanged.

The fixture reuses inventoryWindowOwner, the actual parser/widgets, renderer,
player/string manager and drawable factory/pool. It owns the seven additional C
words, eight real 140-byte trade cells, text/image/viewport state and callback
observations. Existing fixture strings are preserved through a JSON round trip
before appending trade-specific entries; ReadJSON alone replaces the manager's
entries. A newly authored Trade.wnd exercises real construction.

Historical initial three-test run: build/port-client-trade-ui/develop-a. The construction
capture has an empty development expectation; no baseline is frozen. Additional
preflight work covers allocation failures, missing-item removal, fixed 32-code
capacity and repeated quantity-dialog lifetime. Expand across every scoped
routine and meaningful integration before freezing the C baseline. Keep original
failure evidence and review justified prerequisites, rather than using invalid
pointer accesses or adjacent-memory writes as an oracle.

## C development and prerequisites (historical)

Development-a passed all three initial tests in 178.944s. Development-b completed
eight root tests in 23.846s and demonstrated two bugs with independent contracts:
quantity reopening closed the dialog and discarded its new drawable, and both
trade grids/preferred-cell lookup accepted full 32-item stacks. Its separate
callback failure was a fixture error: direct RawEvent bypassed the typed static
text event. StaticTextSetText plus independent readback corrects that fixture;
production numeric parsing is unchanged.

Further reversible C prerequisites are now under test:

- Allocate the replacement quantity drawable first; on failure leave the existing
  dialog intact. On success close/release its old item before publishing the new
  item and opening the dialog.
- Skip full stacks in both searches/preferred-cell checks and guard the final
  32-code append. Check the real trade drawable allocation before dereferencing it.
- Return after reporting a missing trade item instead of dereferencing null.
- Declare tooltip mouse coordinates as an unsigned scalar across C/Go, preserving
  the actual 32-bit callback word without treating it as a Go pointer.

These add 17 physical C lines, bringing the working tree to **82,649 / 96 files /
zero reference C**. No complete C baseline or conversion is qualified yet.
Development-c selects 12 tests, adding quantity prices/callback text/placement,
real allocation failures, missing resources and missing-item reporting. The
already owned centered-message storage is now included in trade captures. All
capture expectations remain empty until coverage and repeatability are complete.

Development-c passed all 12 root tests in 179.541s. The next run, develop-d,
selects 16 tests, adding all acceptance-bit values, signed money formatting,
player/partner start gates and reset deletion order/idempotence. Acceptance state
now fully assigns the source byte before shifting, removing the partial-word
read without changing C LOC. The constructor now independently checks its loaded
trade tooltip. Existing development captures are retained; nothing is frozen yet.

Development-d passed 16 tests in 180.778s; all eight c/d common groups were
byte-identical. Development-e completed 19 tests in 37.087s with two failures.
The input contract assumed 30 ticks/sec while the reused fixture runs at 60;
it now explicitly checks both rates. The actual reset-during-drag defect retained
cursor/capture. A focused unchanged-C check in 20.673s additionally proved a
one-item drag leaked, while larger stacks freed their drawable but retained the
cursor/capture.

The reset prerequisite now clears the trade cursor/capture, checks both four-cell
grids for ownership, and deletes a detached drag before normal grid teardown.
Checking actual grid ownership also handles a new item report arriving while the
old item is held. New contracts exercise that sequence. This adds 17 more C
lines: **34 prerequisite lines total; 82,666 working C lines / 96 files / zero
reference C**. Development-f selects 22 tests, adding capacity, actual modifier-ID
lookup and removal value/ordering/lifetime coverage. No expectations are frozen.


Development-f passed all 22 roots in 179.831s. Development-g passed 26 in
63.838s, adding all quantity mouse event classes, pressed-button combinations,
populated trade rendering/hover and direct 32-slot code helpers. Fifteen of the
17 common f/g groups matched; the other two exposed dangling reset label text.
The static-text widget retains its supplied C string, but reset supplied a stack
buffer. Full field differences are retained in f-g-diff-*; do not normalize away
these text changes.

Development-h completed 28 roots in 63.904s. The shipped Trade.wnd/MultMove.wnd
integration passed (14 captures), using real parser/widgets/renderer and authored
trade reports, not a two-client multiplayer scenario. The initial destroy/recreate
contract incorrectly expected show to return 1 (it returns the final text-setter
result); it now checks actual active state. The corrected lifetime-before run
completed both roots in 21.175s and reproduced the retained quantity active flag
and expired reset labels.

Two further local prerequisites: quantity destruction clears its active flag;
trade reset's text buffer has static storage for the existing borrowing widgets.
The latter deliberately avoids changing shared static-text event ownership in
this batch. It preserves the three reset labels until the next reset, when all
three receive the same new text. The native port must give them equivalent stable
storage. These add one physical C line: **35 prerequisite lines total; 82,667 C
lines / 96 files / zero reference C**. Independent contracts check the reset text
and 100 destroy/recreate cycles. Development-i also adds preferred-cell routing,
inactive reports, unstackable equipment and all five quantity modifier words.
Expectations remain unfrozen. The 30-root run and production builds are in flight.
Completed development captures are losslessly compressed with hash manifests;
original assets are untouched.


Development-i was stopped after preserving its stack trace (270.548s): the new
100-cycle fixture omitted the game frame's FreeDestroyed pass and exhausted the
GUI window pool. NewWindowRaw then passed nil to setExt, whose panic left its
mutex locked; test cleanup blocked trying to acquire it. The fixture now performs
deferred GUI reclamation each cycle. General GUI allocation-failure handling is
recorded for a separate review; it is not changed or claimed fixed here.

Development-j passed all **30 roots / 3,216 captured results / 23 groups** in
64.643s, with 42,436 independent grid contracts and repeated ownership/lifecycle
checks. All three C production binaries passed ELF32/i386/SSE2/CGO and 36-symbol
checks (66.034s / 9.121s / 77.183s), with no tagged test helpers. The full asset
suite matches the exact 1,553 known failures and 15/3/32 package outcomes.
The nine-screen unchanged inventory gameplay comparison passed in 49.358s.
The broader default C suite passed all 346 roots in 183.491s; all 23 raw groups
match development-j exactly. Server/highres qualification is still in flight.

Additional quantity gameplay is in development: use the existing local developer
console to create an apple stack and then exercise the inventory drop dialog with
real mouse input. The first setup reached the console but text-input events did
not enter the command; development2 uses ordinary key presses. Do not count these
setup screenshots as quantity-dialog coverage until the dialog is visibly reached.


## Qualified C baseline

Final affected tests: **346 default / 344 server / 346 highres**, all pass in
183.491s / 280.225s / 202.863s. Every one of the 23 groups is byte-identical to
development-j in all three independent processes. Freeze **3,216 results / 23
hashes** in the tagged tests. The final locked run passed all 30 roots in
50.881s. Its first attempt exited 143 after successful discovery but before any
suite output or result file; the cause is not established. Preserve
c-locked-interruption.json and both logs. No source changed for the retry.

All production architecture/interface checks, exact full-suite comparison and
nine-screen headless inventory comparison passed as above. Source fingerprints
were unchanged during qualification; only the 23 expectation literals changed
for the final locked check. Evidence: build/port-client-trade-ui/c-qualification.json.
There are **82,667 physical production C lines / 96 files / zero reference C**.
The 36 routines now occupy 1,299 address-block lines plus 17 TU preamble lines:
**1,316 removable C lines**. The cumulative captured corpus is 458,695 results /
1,208 groups; independent property/lifetime assertions are additional.

### Integration limits

The shipped-window test passes in every target with OPENNOX_TRADE_UI_ASSETS set
to the original data directory. It exercises quantity up/accept callbacks,
pricing layouts and a start/add/draw/accept/finish trade flow with authored
reports. This is genuine widget/resource integration, not two-client gameplay.
The fresh nine-screen game comparison covers the inventory connection and
teardown, but does not itself open a quantity or P2P dialog.

Five extra quantity-gameplay setup runs are retained locally. English console
entry uses ordinary key presses (text-input events take the IME path), and the
existing local developer spawn command creates apples. The attempted fixed
pickup/teleport sequence did not establish an inventory stack; no screenshot
from those runs is a quantity-dialog oracle. Defer further setup refinement
rather than make the port wait on an open-ended gameplay script. The baseline
is qualified by the independent contracts, repeated captures, shipped-window
integration and stated broader checks; retain this integration gap for native
qualification and later multiplayer testing. This is a reversible testing-scope
decision under standing authorization, not a claim that the missing scenario passed.

## Native conversion

All 36 routines now use Go quantity-dialog and trade-window owners. Removed
1,316 physical C lines, including client__gui__guitrade.c: **81,351 C lines / 95
files / zero reference C**. Typed trade cells assert their 140-byte size and
136-byte value offset. Go callers and the tagged fixture dispatcher call Go
directly. Eleven interfaces remain for real C decoder/shop callers and the raw
GUI tooltip callback; 25 private interfaces are retired.

Quantity acceptance passes a C-owned point allocation and four scalar words to
the callback, then releases the point. Stable interned reset labels respect the
existing widget's borrowing lifetime. The existing shared production C number
parser remains in use; no C algorithm is retained solely for testing.

The add-report return declaration is now uint32_t instead of char*, matching its
numeric price. Actual decoder callers ignore the result; the 32-bit return word
is preserved without manufacturing a Go pointer from a price. The quantity
modifier argument is void* in the C declaration to match cgo's generated bridge;
the implementation only copies the modifier bytes. These reversible ABI typing
decisions preserve actual usage and are recorded for review.

Native-a failed discovery on the const/void* declaration conflict in 66.964s.
The job was joined and no compiler remained before correcting that declaration.
Native-b then passed all 30 focused roots in 184.917s: **every one of the 23
capture groups matched the frozen C bytes on the first behavioral Go run**.
No behavioral expectation changed.

Two additional independent ABI tests exercise 100 quantity callback cycles with
forced GC inside the actual callback, and both trade grids across seven price
words (including high-bit values), entering through real retained C interfaces.
The latter also invokes the widget's actual raw C tooltip callback. Native-c
passed all 32 roots in 184.355s with all 3,216 captured results unchanged.

## Native qualification

All **909 accumulated default / 346 affected server / 348 affected highres**
root tests passed in 511.865s / 277.512s / 201.154s. Every selected test started
and finished. All 23 full capture groups match the frozen C hashes in all three
targets; the shipped-window integration ran with original assets enabled.

Three production builds succeeded in 58.802s / 8.400s / 55.419s. Each is
ELF32/i386/SSE2/CGO, with all eleven retained interfaces supplied by Go bridges,
all twenty-five retired symbols absent, and no tagged test helpers. The full
asset suite in 49.523s has exactly the same 1,553 known failure entries and
15 pass / 3 fail / 32 skip package outcomes. The fresh unchanged nine-screen
inventory game comparison passed in **51.832s**, with reference updates disabled.

All 1,602 Go/C/header fingerprints stayed unchanged throughout qualification.
Evidence: build/port-client-trade-ui/native-qualification.json, native-b-parity.json,
native-qualified-capture-compression.json and native-binary-verification.json.
Completed captures are losslessly compressed with SHA-256 manifests. All jobs
are joined. The corpus remains **458,695 results / 1,208 groups**; additional
independent grid/lifetime/ABI assertions are outside that count.

The integration limits above remain: authored reports on real shipped widgets
are covered; a two-client trade game and the extra quantity-pickup gameplay setup
are not qualified. No claim of coverage is based on those unfinished scenarios.
Next is the connected shop UI: 42 scoped routines and about 1,088 removable C
lines before any independently justified baseline prerequisites.
