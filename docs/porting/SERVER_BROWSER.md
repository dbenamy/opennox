# Server browser — Go conversion

The 64 live routines in the selected browser batch are Go. One unused helper was
removed, 58 private C interfaces and 34 C global owners retired, and the empty
noxworld C file was deleted. Seven C entry points remain for existing callers and
animation slots. The preceding character initializer's unused C adapter and empty
header also retire. Production C is **14,691 physical lines in 58 files**, zero
reference C: **−2,160** from the qualified parent.

See [selection](server-browser-selection.json), [interface plan](server-browser-interface-plan.json),
[original-C qualification](server-browser-c-qualification.json) and
[native qualification](server-browser-native-qualification.json).

## Evidence

The original-C baseline is committed as **04061130**; its unchanged production
parent is **d8133587**. Native qualification runs 197 affected roots on each of
default/server/highres, with no skips. All 142 artifacts /100,275 records match the
original C captures. The 18 browser-specific captures contain 72,402 records.
Fresh binaries, ABI inventories, the exact known full-suite result and isolated
browser, character-creation/gameplay and save/load scenarios are required gates.
All gates pass on the final source. The known full suite matches exactly: 1,553
failure entries, 15 passing /3 failing /32 skipped packages.

Local evidence: `build/port-server-browser/native-final-{default,server,highres,production}`
and `static-native-final.log`. Frozen expectations are in the fixtures and
[capture index](server-browser-captures.json); do not regenerate them.

Coverage includes all 65,536 mode words; 4,732 radius/quadrant cases; signed/unsigned
popup-position boundaries; actual linked lists and all 10 sort modes; selected
record lifetime; libc IPv4 forms and high ports; LAN text, marker geometry/images,
server restrictions, host descriptions and the shipped proximity-popup resource.
The 100-entry popup boundary, repeated cleanup and numeric parser notifications
have independent checks. Additional original-C contracts cover byte-widened text,
clock truncation/addition/equality, all 11 connection error texts, password/connected
dialog transitions, connected draw routing, column synchronization and UTF-16
password truncation. The latter uses a local send stub. A real static-text widget
regression verifies browser labels remain valid after temporary allocations are reused.

The headless browser scenario opens, sorts, closes, reopens and proceeds to host
class selection, comparing six exact screens. `OPENNOX_ISOLATE_NETWORK=1` confines
that scenario to a child network namespace with loopback and an unavailable local
lobby endpoint. It requires no external service and leaves host networking intact.

The first native browser scenario caught freed static-label text that the focused
listbox tests did not cover. Static labels retain their pointers; listbox rows copy
their text. Use persistent GUI text storage for the four browser labels. The
regression was demonstrated before the fix. A second screenshot comparison isolated
a scrollbar difference: the original constructor offsets the thumb child, whereas
the initial translation offset its parent slider. Corrected the child target; the
unchanged six-screen original-C scenario serves as the regression. All final gates
run on the corrected source. A fresh default-client browser-only preflight then
passed all six screens before the final qualification sweep. For GUI batches, run
this relevant scenario early enough to catch constructor/rendering differences
before rebuilding all targets and repeating the full suite.

## Decisions for review

- Keep the original unsigned coordinate subtraction before radius calculation,
  unsigned popup Y comparison, four-region fallback and signed-short high-port
  lookup behavior. These surprising behaviors are covered rather than silently
  corrected during the port.
- Displayed narrow strings widen bytes individually; they do not decode UTF-8.
  Sort wide strings as UTF-16 code units and narrow names with ASCII case folding.
- Platform ticks truncate to 32 bits at the old C boundary. Connection-test and
  password deadlines add 20,000 before widening; delayed joining widens before
  adding 1,000. Preserve both and their distinct equality comparisons.
- Preserve the unconditional legacy highres protocol version and maximum video
  dimensions, including default/server builds. The C define came from a Go file
  without a highres build constraint; normalizing it would change compatibility.
- Re-sort reuses list nodes. Before replacing or freeing a selected node, copy its
  169-byte record into one owned snapshot. Keep that snapshot and the 12-byte list
  sentinel for browser/session lifetime: gameplay still reads the selected
  endpoint after the UI closes. This replaces the original abandoned-list leaks.
  Invalid sort selectors return without allocating a lost record.
- Initialize the formerly undefined password-request header padding and the color/
  level fields of the browser map-polygon temporary record. The sender writes the
  first three header bytes; payload truncation remains unchanged. These are
  deterministic behavior corrections under the user's reversible-decision policy.
- Browser state now has a Go owner. Its fixed-width ABI types preserve existing
  interop layouts; no selected C global definitions or algorithms remain.

Twelve old configuration setters are still registered in a blob table that appears
unconsumed. Audit actual reads and relocated pointers before removing them. The
coordinate getter had live browser callers and was included in this conversion.

## Recovery and disk

All installed drafts, installers, boundary-retirement and completed cleanup scripts
are consumed. Original selection byte offsets are stale after source removal;
source and committed expectations govern recovery.

Completed C browser asset copies were verified/deduplicated, reclaiming 1,112,777,500
bytes. During native conversion, 24,120 compiler-cache entries not read or written
in 9 hours were removed after all builds joined and inode/size/timestamps matched
an audit manifest, reclaiming 7,527,475,663 bytes. Original assets/archive, captures,
logs and qualified binaries remain. Local manifests are under
`build/port-server-browser` and the scenario run directories.

Completed native browser asset copies were also verified against the original
assets and deduplicated, reclaiming 1,669,166,250 bytes. Restoration manifests,
screenshots and diagnostic logs remain; the deletion helper is consumed.
