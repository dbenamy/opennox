# Client interaction and remaining game dialogs — qualified C baseline

The corrected C baseline is now qualified and frozen; conversion is next.
Production C: **13,456 lines /57 files /zero reference C**, +14 prerequisite lines
from parent `b30f96a3`. Recovery checkpoints `6c851a78`, `ac2c9182` and `c856a315`
preceded full qualification. Historical development sections below describe their
then-pending gates; the final qualification section supersedes those statuses.

## Scope and reachability

The provisional connected selection contains83 C bodies /1,496 body lines (including14 prerequisite lines) across
chat input/messages, cursor hover, talk/use/trade/pickup, HUD/modal/escape gates,
observer/vote/chat indicators, conversation, game-over, help and key/paper-doll UI.
See [selection](client-interaction-selection.json),
[caller inventory](client-interaction-callers.json) and
[reachability](client-interaction-reachability.json). The graph has55 textual
external roots and79 reachable bodies; this inventory still needs manual root,
callback and enclosing-condition review. Keep packet reorder/update decoding
separate from the GUI behavior batch.

Four apparently orphaned bodies are435690 and the connection-dialog graph
49C820→49C910/49CA60. The latter is the only production writer of owner1305684;
its boolean getter49CB40 remains in HUD/MOTD guards. Current literal/raw-address
search finds only metadata outside the named accesses. Complete the registration
review before retirement. The older MOTD fixture artificially sets that owner:
if removed, explicitly retire its unreachable gate and compare every surviving
capture row with the committed original baseline. Do not keep dead production
code solely to satisfy a fixture.

## Initial contracts and evidence

Thirteen original-C roots pass and independently repeat on the default client,
with all capture hashes unchanged, no skips, and static
memory checks. [Checkpoint](client-interaction-initial-checkpoint.json) records
observed capture hashes, not final frozen expectations. Source fingerprints stayed
unchanged during the run. No server/highres or complete affected sweep is claimed.

- Mouse modes, signed arithmetic, all32 valid buff bits, animation/observer gates,
  UTF-16 text state, absolute cursor movement and32-to64-bit tick storage.
- Actual shipped chat-window construction, start/reentry/close/destroy, ownership,
  numeric resource notifications and16-bit length-field mutation.
- Chat/observer/vote indicator geometry, visibility and image offsets, including
  odd viewports and player-present/player-absent drawing.
- 714 talk/use/trade cases through the real message queue: nil owners, blocking
  flags, exact predicates, wire bytes and unit-code boundaries.
- Key-window state transitions/sounds; game-over labels, countdown wraparound,
  host restart messages and quit callback ordering.
- Conversation resource/list links, choice/replay/close actions, byte-sized choice
  mutation and background anchoring. Replay uses the actual dialog queue owner.
- All eight shipped help resource-table entries, font-height fallback, host/client
  key text, quest/co-op automatic close and checkbox persistence.
- Modal/frame stamps, unsigned frame subtraction and console short-circuit order.

Two failed fixture assumptions were corrected before the passing run. Resource
background images use the root image loader, separately from the legacy image
hook; the pixel contract now assigns an owned nonempty image. Registering a font
named default does not change the renderer fallback face; the help contract now
owns both. No production behavior or expected C capture was changed to conceal
a port mismatch. The first test-driver invocation also used a regex where the
runner expects a pattern file; that invocation ran no tests and is superseded.

Rerun from the repository root with a fresh output directory:

```bash
source build/baseline/env.sh
GOMAXPROCS=2 GOMEMLIMIT=768MiB python3 tools/porting/run_batch.py \
  docs/porting/client-interaction-initial-batch.json --phase focused \
  --out build/port-client-interaction/check-initial
```

## Remaining baseline work

Complete hover enumeration, paper-doll drawing, chat draw/submit/input, HUD
visibility and ordered Escape handling. The extended checkpoint below covers
centered/chat/console text, pickup/secondary and direct cursor geometry/gates. The
textual coverage inventory under ignored build/ is only a review aid; transitive
references do not prove branch execution. Audit shared owners beyond dword names,
raw callback slots, original tables, return conventions and pointer lifetimes.

Text review found that Go callers can supply full strings to fixed centered-message
rows and the chat-format temporary. Do not capture writes beyond their storage
as compatibility expectations. Establish supported boundaries and make any
justified reversible prerequisite correction before freezing; the text
prerequisites described below are now installed. Then repeat original captures, qualify the affected three
targets and production evidence, commit the final baseline, and translate.

Completed final session-dialog scenario copies were hash-deduplicated after their
qualification, reclaiming another1,669,136,451 bytes. Per-run restoration manifests,
reference screenshots/saves and reports remain; original assets/archive are intact.


## Text prerequisites in progress (not qualified)

Direct Go callers do not impose the widget's250-unit chat limit or the centered
row's317-unit text limit. Bound centered-row copying to317 UTF-16 units plus NUL,
preserving expiry/flag neighbors and the complete console message. Allocate the
chat-format temporary from full input length while retaining the original byte
wire count, encoding choice, leading-space handling and queue return value.

The full-console path also reached a shared static512-unit format buffer. Use
bounded formatting on a512-unit local fast path, and an input-sized allocation
and variadic retry for longer output. The Go print adapter copies the resulting
string synchronously; no C caller retains the temporary. Return failure on an
unrepresentable size or allocation failure. The other remaining C use is a
localized console notice in client network decoding; review and production
qualification must include that shared helper. This is a reversible prerequisite
under the standing authorization, not a claim of defined compatibility for old
writes beyond these arrays.

The first long-text test run was stopped during discovery to install the console
correction before exercising the unsafe path. Independent tests check neighboring
rows/control words, unchanged input, full console text, ASCII/high-byte/raw UTF-16
messages, lengths through4,096, byte-count wrap, and mixed variadic arguments across
the fast-path boundary. All-target and production qualification are pending.
Working C is13,456 /57 files /zero reference, +14 prerequisite lines; header changes
are not included in this metric. Production-qualified C remains13,442.


## Extended development checkpoint

The [extended manifest](client-interaction-extended-checkpoint.json) records
19/18/19 passing roots on default/server/highres, with matching captures and2,737
source fingerprints. Default repeats independently with identical captures;
static-extended passes. Run [the batch](client-interaction-extended-batch.json)
with `--phase default`, `--phase server` or `--phase highres` and a fresh output
directory. This checkpoint preserves tested prerequisites and contracts; it does
not claim complete affected-corpus or production qualification. Capture hashes
remain observations pending final baseline review/freeze.

New independent contracts cover27 centered-row boundaries,624 chat length/encoding
cases,44 mixed variadic console-format cases,48 message pixel/expiry cases and90
pickup/secondary requests. Pickup uses the actual inventory grid, tests reserved
rows versus usable space, stack limits and currency exceptions, and checks real
queue bytes, sounds and full console rejection text. Cursor tests compare2,250
circle/box grid points with independent capsule/diamond-prism geometry,18 eligibility
gates and equal/greater depth ordering. They use actual client visibility; only
this root is omitted on server because its visibility implementation deliberately
panics as unreachable. Both client variants cover it.

The nonempty-pixel assertion caught transparent fixture colors; the text-spacing
comparison caught use of cell height instead of11px cap height. Cursor comparison
caught an initially incorrect expected top cap; tracing the unchanged C established
the capsule contract. These were fixture corrections, not production changes.
Minor import/discovery failures produced no accepted evidence and are superseded.
All initial thirteen capture hashes remain unchanged after the text prerequisites.

Disk maintenance removed62 reproducible Go cache entries >=64MiB, unmodified and
unaccessed for six hours, only after every build joined. Reclaimed4,364,972,826 bytes;
ignored audit/removal manifests remain. Assets, source, binaries and qualification
evidence were preserved. Removal is consumed; Go regenerates missing cache data.


## Additional controller and rendering contracts

The default development run `c-escape-complete-initial` passes30 roots without
skips. All26 captures from `c-doll-layer-initial` match `c-blink-second`; all29
captures from that run match the new thirty-root run. All-target repetition is
in progress; these are observed development results, not the frozen baseline.

New contracts exercise192 HUD combinations,33 chat layout cases and24 submit
cases;352 ordered Escape combinations run both directly and through the actual
chat key callback. Eight combinations cover quantity, expanded quickbar and
identification closure before conversation dispatch. Another45 cases cover
pending spellbook addition completion and the earlier modal gates, with actual
particle records and RNG consumption captured.

Rendering uses111 shipped paper-doll asset names/table slots,256 independent
material/layer cases,640 full composition cases and108 conversation blink cases.
Full composition checks gender/base/bald choices, armor/cloak/weapon ordering,
palette state, background and nonempty pixels. Hover enumeration uses the actual
spatial index in both insertion orders and two world offsets.

Fixture corrections preserve production behavior: initialize the actual queue
before chat submission; delegate mouse position to the owned input rather than
the reused synthetic getter; explicitly assign allocated point coordinates
because alloc.New does not copy its sample; advance frame sequence with Input.Tick
rather than its separate event counter. The first blink run was terminated and
joined after that counter mistake, and supplies no accepted evidence. All
installed source drafts are consumed. No additional production changes were made.


The [contract checkpoint](client-interaction-contract-checkpoint.json) now records
30/28/30 passing roots and31/29/31 captures on default/server/highres. All shared
captures match, default independently repeats, static-contracts passes, and2,746
source fingerprints agree. Both client-only hover roots are excluded on server.
Rerun [the contract batch](client-interaction-contract-batch.json) with the desired
`--phase default|server|highres` and a fresh output directory. This checkpoint
remains separate from final baseline and production qualification.


## Corrected-C qualification in progress

Recovery checkpoint `c856a315` is pushed. The affected selection passes493/490/493
roots on default/server/highres in152.3/149.6/149.6 seconds, without skips. All
227/226/227 inherited captures and31/29/31 new interaction captures agree with
the preceding observations. Source fingerprints agree across all runs.
The [qualification batch](client-interaction-batch.json) retains the parent's
exact capture gates and adds inventory, quickbar, spellbook, console, visibility
and chat dependencies. Fresh production qualification also compares the existing
nine-screen inventory-window scenario, alongside dialog/gameplay/save-load.

The refreshed literal caller audit still finds four unreachable bodies. Ten
selected symbols have references in C outside the selected bodies, primarily
client decoding and centered notices; review callback slots separately before
retiring interfaces. Forty-nine selected symbols have Go-to-C references to route
directly to native Go. The ignored previews are audit aids, not final retirement
lists. C remains13,456 physical lines in57 files, zero test-reference C.


## Final corrected-C qualification

All493/490/493 affected roots pass on default/server/highres. All227/226/227 parent
captures remain unchanged, and31/29/31 interaction captures agree across targets
and independent repeats. The two hover roots and inherited inventory world
selection are client-only. Static checks and all2,746 source fingerprints agree.

Three fresh binaries verify ELF32/i386/SSE2/CGO, required interfaces, retired
symbols and absence of test helpers. The asset suite matches the exact1,553 known
failure entries and15 pass /3 fail /32 skip package outcomes. Headless dialog,
gameplay, save/load and nine-screen inventory scenarios all match their existing
references with replacement disabled. Evidence is under c-affected-* and
c-production in build/port-client-interaction; all jobs are joined.

[Qualification](client-interaction-c-qualification.json) and
[frozen interaction captures](client-interaction-c-captures.json) record the
accepted baseline. [The batch manifest](client-interaction-batch.json) now enforces
those hashes alongside the inherited captures. No C algorithms are retained
solely for tests. Translation and native qualification are next.


After qualification, verified duplicate files in the four completed C scenario
copies were removed, reclaiming2,225,495,437 bytes. Restoration manifests preserve
paths, hashes and metadata. Original assets/archive and captures/reports remain;
the deduplication script's deletion mode is consumed.
