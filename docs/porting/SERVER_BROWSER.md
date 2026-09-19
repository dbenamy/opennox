# Server browser — original C baseline qualified

Qualified parent **d8133587** is committed/pushed. Production source is unchanged:
**16,851 physical C lines /59 files /zero reference C**. Initial scope expanded to
65 connected routines /2,067 body lines across noxworld and GAME1_2/GAME1_3/
GAME2_3/GAME3. This includes UI lifecycle, server-list storage/sorting, selected
server details, map hit testing, popup selection, state and formatting helpers.

See [selection](server-browser-selection.json), [callers](server-browser-callers.json)
and [global candidates](server-browser-global-callers.json). The latter still
includes function-name matches and shared globals; complete the ownership/interface
plan before conversion. Read-only draft notes are under build/port-server-browser.

## Baseline coverage

All18 browser roots pass on original C;18 captures /72,402 records match an
independent190-root affected run and are frozen. Coverage includes all 65,536 mode words,
4,732 radius/quadrant cases, 210 popup clamp cases, real linked-list counts, state
words, endpoint formatting and libc address forms, all ten collection sort modes,
LAN rows/text widths, map marker geometry/images/tooltips, detailed restrictions,
re-sort selected-record lifetime, connection/refresh timers and numeric resource
construction notifications. Controlled string-manager entries, real GUI/list owners
and owned constants avoid accidentally testing zero-filled state.

The C server record is packed to169 bytes, while its existing Go view is172 bytes;
fixtures read only the defined169 bytes. Temporary C strings are explicitly freed.
The re-sort fixture releases original C's detached nodes only after checking that
selected-server data remains valid. Native code must own that lifetime explicitly.

Host-description and shipped-resource proximity popup contracts also pass,
including the100-entry selection table guard and repeated cleanup. The final
interface plan anticipates7 retained selected exports,58 retired interfaces and
34 retired C owners; verify these after source removal. Default/server/highres
baseline qualification passed:190 roots each, no skips;142 captures /100,275
records match on all targets and all2,685 source fingerprints agree. Static
mapped-memory checks pass. Production is identical to qualified parentd8133587;
its three binary hashes were reverified, with new browser capture/repeat added. See [plan](server-browser-plan.md) and
[interface plan](server-browser-interface-plan.json).

The new headless scenario passes twice using the qualified parent binary, with six
exact screen comparisons: open, sort, close, reopen, host class selection. The
optional `OPENNOX_ISOLATE_NETWORK=1` runner setting creates a child network namespace
with only loopback and directs lobby discovery to an unavailable local port. Host
networking is unchanged; these tests need no external service. Both runs exited0.
Original data copies were hash-verified and deduplicated, reclaiming1,112,777,500
bytes; captures, logs and restoration manifests remain under the run directories.

## Findings to preserve or review

- Map hit testing subtracts unsigned coordinates before double conversion. Negative
  differences wrap, making selection asymmetric. Initial contracts confirm it;
  preserve during this port and leave any gameplay correction for separate review.
- Popup Y lower-bound comparison is unsigned; negative values can remain large
  unsigned words. Signed overflow boundaries are included in the C contracts.
- Address duplicate detection takes a signed16-bit port and compares it with an
  unsigned16-bit stored value. High ports need explicit compatibility cases.
- Re-sort copies records, and the original list nodes appear unreclaimed. Audit
  selected-server aliases and GUI ownership before changing lifetimes; do not
  simply free old records while callers may still hold them.
- The region lookup tests four entries; unmatched coordinates fall back to region0.
  The marker tests explicitly preserve unsigned coordinate shifts and this fallback.
- The host-description C code always advertises the highres protocol version,
  because legacy/video_highres.go unconditionally enables that C define. This
  differs from Go default/server constants; preserve the existing C output and
  retain the mismatch as a review item. No version word is normalized in captures.
- Twelve old configuration setters are registered in a blob table that appears
  unconsumed. Follow code reads and relocated pointers before deciding reachability.
  The coordinate getter does have four live browser callers.

Retire the now-unused character initializer C adapter sub_4A5E90_A and empty
selclass header with native source changes; retain its active Go hook. No C
algorithms will remain solely for tests.

## Recovery and next step

Baseline qualification is recorded in [the report](server-browser-c-qualification.json)
and [production identity](server-browser-production-identity.json). Frozen browser
expectations are in the fixtures and [capture index](server-browser-captures.json).
The C baseline is ready to commit; native conversion has not begun.

Final local evidence:build/port-server-browser/c-qualified-{default,server,highres},
static-c-final.log and the two server-browser-c headless runs. All sessions joined.
Freeze and qualification scripts are consumed; do not regenerate expectations.
Follow the conversion plan, compare against these frozen contracts, then run fresh
production/ABI/full-suite/browser/gameplay/save-load qualification before committing
the conversion. Keep the deferred configuration-table audit outside this batch.
