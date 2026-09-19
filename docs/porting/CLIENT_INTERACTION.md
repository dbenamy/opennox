# Client interaction and remaining game dialogs — C baseline in development

Qualified production parent: `b30f96a3`, **13,442 C lines /57 files /zero reference C**.
This checkpoint adds tests only. It is not the completed C baseline or a conversion.
All new source files carry the `porttest` build constraint; existing production
source is unchanged. The parent's production qualification remains the current
production evidence.

## Scope and reachability

The provisional connected selection contains83 C bodies /1,482 body lines across
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

Cover remaining cursor hover, paper-doll drawing, pickup/secondary-weapon actions,
centered messages/chat lengths, HUD visibility and ordered Escape handling. The
textual coverage inventory under ignored build/ is only a review aid; transitive
references do not prove branch execution. Audit shared owners beyond dword names,
raw callback slots, original tables, return conventions and pointer lifetimes.

Text review found that Go callers can supply full strings to fixed centered-message
rows and the chat-format temporary. Do not capture writes beyond their storage
as compatibility expectations. Establish supported boundaries and make any
justified reversible prerequisite correction before freezing; none is installed
in this checkpoint. Then repeat original captures, qualify the affected three
targets and production evidence, commit the final baseline, and translate.

Completed final session-dialog scenario copies were hash-deduplicated after their
qualification, reclaiming another1,669,136,451 bytes. Per-run restoration manifests,
reference screenshots/saves and reports remain; original assets/archive are intact.
