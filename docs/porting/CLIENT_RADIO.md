# Radio-button selection and drawing

Status: C baseline and GUI lifecycle correction fully qualified; radio routines
remain C until the next conversion commit.

Scope: five C callbacks in GAME3.c (004A84E0 through 004A93C0), 296 physical
lines. Input, exclusive group selection, focus/text updates, callback setup and
colored/image rendering move together. The constructor is already Go. No C
callers remain after replacing its two Go wrapper references, so all five C
interfaces can retire.

The actual GUI fixture owns three radio buttons and a non-radio sibling. Tests
observe group state during notification callbacks to prove ordering, including
the original behavior of clearing selection on any sibling with the same group
word, regardless of widget style. Cases include root/child windows, default or
explicit notification owners, selected/highlight states, disabled/no-focus flags,
mouse input, exact keyboard states and programmatic selection. Drawing varies
image/color mode, colors/images, image offsets, smoothing, label alignment and
raw UTF-16 text boundaries, narrow/ordinary widths and odd/even heights. Raw window words, copied data, notifications, focus,
pixels and renderer state supply the oracle; only known pointers are normalized.

An ownership regression checks the constructor's copied allocation through
actual GUI destruction without reading freed memory. The original C callback
lacks a destroy case while the Go constructor allocates RadioButtonData. The
applied correction keeps allocation/free responsibility in the Go constructor's
callback wrapper; the five C behaviors remain the baseline for other events.
The original ownership contract failed, and the corrected allocation contract passes.

Preserve unusual but defined behavior: notification happens before deselecting
siblings; an already-selected mouse action returns one only when highlighted;
keyboard activation toggles the selection bit after notification, whereas mouse
and programmatic selection set it. Group iteration considers all sibling styles.
Programmatic selection notifies only for raw argument one. Text copies exactly
63 UTF-16 code units and forces a terminator, retaining unpaired surrogate bits.
Tab-navigation callbacks are original no-ops. These compatibility decisions and
the ownership correction are reversible and should be reviewed later.

Two C capture groups now repeat exactly: 11,760 input results and 4,114
drawing/text results (15,874 / two groups). The initial six focused tests pass, including
selection ordering, radio allocation release, deferred callback delivery and
slider allocation release. Broad C-baseline qualification passes.


## GUI lifecycle correction

The original ownership contract fails with a live allocation after actual GUI
destruction. A second contract receives zero deferred cleanup callbacks. The
shared GUI marks a window destroyed, clears callbacks, then later tries ordinary
Func94 dispatch—which rejects destroyed windows. Also, the duplicate-destroy
check reads flags through a getter that hides destroyed windows.

The applied correction saves a dedicated cleanup callback while disabling ordinary
callbacks, invokes that callback once during FreeDestroyed, and rejects repeated
Destroy directly. The radio constructor separately frees its owned data. Deferred
timing and ordinary-event rejection are preserved. Broad qualification of the shared fix passes; the additional callback-mutation
contract also passes in all three configurations.


The deferred-cleanup contract also checks duplicate Destroy calls, ordinary-event
rejection after queuing, parent/child queue order, one-time delivery and a second
window queued by a cleanup callback. The slider contract proves its copied data
stays live until FreeDestroyed and is then released in both orientations. Original
failures and corrected captures are in build/port-client-radio. These frozen expectations precede the native radio conversion.

The GUI correction preserves its existing queue order and timing. It stores the
cleanup callback outside ordinary dispatch and clears it before invoking it; live
window layout and C ABI do not change. Full accumulated tests and fresh gameplay
are required here because the correction affects more than radio buttons.


## Broad C-baseline qualification

With the shared GUI correction applied and all five radio routines still in C:
accumulated726 passes (402.806s, one optional skip), affected server92/highres93
pass, the C client builds, fresh gameplay passes in36.342s, and the full asset
suite exactly matches the known1,553 failure entries and package outcomes
(15 pass /3 fail /32 no-test packages). Source fingerprints remained unchanged.

A final independent contract, added after those checks without changing production
source, has the notification owner set selection and change group. It distinguishes
keyboard XOR from mouse/programmatic OR and proves sibling iteration observes
the changed group. Seven focused tests pass under each configuration;
existing C capture hashes remain unchanged.
