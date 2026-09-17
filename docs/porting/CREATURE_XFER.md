# Monster and NPC serialization

## Scope and baseline plan

Continue after qualified item/reward serialization **a7814f6e**. The next connected
block is GAME4_2.c from 00528DB0 through EOF: **1,570 C lines / ten functions**.
It covers monster/NPC callbacks, spellcaster defaults, action serialization and
its timestamp/argument helpers, buff and voice records, equipment normalization,
and post-load action-reference resolution. Current qualified C remains **58,539
lines / 82 production files / zero reference C**.

Audit existing registrations and callers before choosing retained exports versus
private Go helpers. Keep unrelated declarations where remaining C callers need
them. Reuse the qualified common/item stream and real object/type ownership
fixtures; extend them with actual monster definitions, action metadata, waypoints,
spell/buff owners, sound sets and inventory relationships as required.

Start with current registered-callback read/write contracts and version rejection.
Then cover historical field boundaries, names and masks, scripts, direction and
float conversion, health/default policy, NPC colors/voice, shop payloads, inventory,
buff side effects and timers. Action contracts must cover frame adjustments,
argument kinds, path/waypoint references, stack boundaries and post-load resolution.
Check full state, exact fields/stream positions, checksums, pointer identity and
ownership. Repeat C captures before freezing; establish independent contracts for
any ambiguous behavior before a reversible correction.

No creature callback has been translated. Preparatory caller references are under
build/port-creature-xfer/caller-audit.json. The initial manifest runs creature plus
qualified item/common serialization checks on all targets and a separate repeat.
Final affected and production qualification will be selected after the caller/owner
audit, preserving the exact known suite and gameplay/save-load/flat references.

## Initial caller audit

Retain the Monster/NPC registered callbacks and sub_52BAF0, which the server's
post-load phase calls. The other seven helpers have no callers outside this block
and can become private Go functions after moving their callers. Preserve any
unrelated waypoint declaration needed elsewhere.

Initial fixtures are installed for current records and future-version rejection,
using real registered callbacks and MonsterUpdateData-sized buffers. Timestamp/helper bridges are installed; ignored drafts are consumed.
No expectation is frozen and no creature implementation has changed.

The first current-record writer stopped in its action-name lookup because the small
fixture had not initialized the runtime AI-name pointer table. A debugger located
the missing table at the writer's name-length operation. The fixture now owns and
restores all 72 production AI names while keeping the actual C resolver. This is
fixture setup, not evidence of an engine defect. Next, own the shipped numeric
argument-kind table and allowed default actions as well for extended action tests.

Local disk space was recovered by removing only hash-verified duplicate original
assets from seven older completed runs. Their screenshots/logs and restoration
manifests remain; original assets and the active batch are untouched.

## Initial C evidence

With the required runtime tables owned/restored, **106 leaf cases pass** in
c-initial3.log: current Monster/NPC round trips, future-version rejections and
100 timestamp cases. The current records are 273 and 303 bytes. The timestamp
helper returns the raw wrapped sum while storing at least 1 using a signed
comparison. No expectation is frozen yet.

The action stream/gate fixture is installed. It independently builds an empty
stack/path record, checks frame shifts and serialized fields, and captures both
writer and reader state. Some timestamp adjustments run on writing as well as
reading; preserve their minimum-1 side effect rather than assuming saves are pure.

Current action stream/gate checks pass 28 cases; historical action reads pass
140 cases including signed nonpositive versions, which this helper accepts using
its oldest layout. The full 72-action argument fixture initially assumed a byte
count return; the existing stream adapter actually returns a success flag. Its
contract was corrected against that adapter before freezing. Owner isolation now
also restores the object lookup cache, monster-definition list and voice-set list.

## Expanded helper contracts (not frozen)

All 72 shipped action kinds pass writer/reader argument contracts (144 cases),
including real object/waypoint identities and normalized captures. Post-load
reference conversion passes 220 cases: absent/live/destroyed objects, waypoint
lookup, full 24-entry stacks and disabled signed-stack values. Inventory conflict
normalization passes all 781 sequences of up to four items across five categories;
full object state is checked, allowing only expected equipped-flag changes.

Twelve action edge cases pass: empty/unknown/maximum-length names, absent final
word, supported but unshipped argument kind 6, and unsupported kinds. The C helper
returns the final stream adapter's boolean success; unsupported kinds return the
kind value and leave the argument/tail unread. These are compatibility contracts,
not reasons to regenerate expectations during the translation.

Logs: c-action-arguments-development2.log, c-postload-development.log,
c-equipment-development.log and c-action-edges-development.log under
build/port-creature-xfer. Definition-default contracts are in development. No
creature production source has changed, and the baseline still needs complete
callback/historical/buff coverage plus repeat and integration qualification.

Definition defaults now pass 108 cases against actual linked definitions, covering
health preservation/truncation, selective status bits and retreat/resume overrides.
Seven voice cases cover empty/255-byte names, actual linked lookup, clearing unknown
names and absent sets. Buff writes pass 72 cases, preserving the shipped order and
shield duration-record/default-100 behavior; fourteen version/name rejection cases
check exact stopping positions and unchanged object state.

Seventy buff-application cases use the real spell acceptance and infravision effect:
levels above five clamp, the applied buff power reflects that clamp, and the saved
32-bit duration overwrites the resulting timer through its 16-bit field. Versions
one and two agree. Thirty-two path cases reach all 32 coordinate pairs and sixteen
waypoints, checking actual pending-waypoint lookup and missing-reference results.
The first nonempty waypoint contract exposed a reserved-word detail: C serializes
the completed waypoint-loop count later in the record. Preserve this existing wire
behavior; the independent builder was corrected before freezing.

Additional logs: c-defaults-development.log, c-voice-development.log,
c-buffs-development.log and c-paths-buff-application.log. The first paths compile
used Next instead of the actual WpNext field; this was only a fixture correction.
Separate-process helper captures are in progress. Main callback historical records,
merchant payloads, scripts/spell masks, nonempty action references/stack records and
auto-spell defaults still need coverage before calling this a frozen C baseline.

## Recovery checkpoint: helper baseline

Two separate default-target processes pass **16 roots / 1,734 leaf cases**, zero
failures/skips. Their **1,707 records / thirteen capture groups** match byte for
byte (`c-helper-check-{a,b}.jsonl`, `helper-check-{a,b}-*.json`). These are development
captures, not frozen expectations: callback expansion and all-target/integration
qualification remain outstanding. A recovery commit saves this tested fixture work
without claiming a completed conversion. C remains **58,539 / 82 files**.
