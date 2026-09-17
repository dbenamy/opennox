# Monster and NPC serialization

## Current status

The C baseline is frozen and qualified; translation is next.
No creature production implementation has changed. Current C remains **58,539
physical lines / 82 production files / zero reference C**. Recovery checkpoint
**48c0ba80** is committed and pushed. Item/reward conversion **a7814f6e** remains
the latest qualified conversion.

Before freezing, default/server/highres and a separate default process each passed
**61 roots / 5,108 leaf cases**, without skips. All **4,908 records / 44 groups**
matched byte for byte: **1,977 creature records / twenty groups**, plus the qualified
2,931 item/common records / 24 groups. Frozen checks also pass on all targets and the separate repeat. All three
production builds pass the C-symbol audit. The full asset suite matches its exact
known 1,553 failure entries (15 packages pass / three fail / 32 skip). Gameplay
matches 41 frames, actual save/reload seven, and flat rendering fourteen with exact
map regeneration. All phase source fingerprints stayed unchanged.

## Scope and retained interfaces

GAME4_2.c from 00528DB0 through EOF: **1,570 C lines / ten functions**. This covers
Monster/NPC callbacks, spellcaster defaults, action serialization and timestamp/
argument helpers, buff and voice records, equipment normalization, and post-load
reference conversion.

Retain Monster/NPC registered callbacks and sub_52BAF0, called by the server's
post-load phase. The other seven helpers have no callers outside this block and
can become private Go functions after moving their callers. Preserve unrelated
waypoint declarations still needed by remaining C. Caller audit:
build/port-creature-xfer/caller-audit.json.

## Contracts and owners

Fixtures use actual registered callbacks, owned monster/type buffers, production
AI names and shipped numeric buff/action/direction tables. Actual linked monster
definitions and voice sets, pending waypoints, live object lists, lookup-cache
owners, inventory chains, duration records and spell definitions preserve runtime
lookup/application behavior. Pointer identity is checked before normalizing only
identified reference fields. Full state, exact fields, return values, stream ends,
wire bytes and checksums are checked as appropriate.

- Current Monster/NPC round trips and future-version rejection.
- Every positive supported historical version: Monster 1–64, NPC 1–62. Independent
  plaintext builders check script layouts, field widths, mirrored movement values,
  health/max-health rules and colors. Old Monster scripts use the actual file-handle
  registry with matching close/release ownership.
- Spell tables: raw/named layouts, empty/sparse/alternating/full 137 slots, clearing
  absent entries, and the current writer's exact count/name/value section.
- Action frame shifts and signed/wrapping timestamp boundaries; all 72 shipped
  argument layouts with present/absent object and waypoint identities; missing tail,
  empty/unknown/255-byte names, unshipped supported kind 6 and unsupported kinds.
- Every valid action-stack length through 24. Paths reach all 32 coordinate pairs
  and sixteen pending-waypoint references. Object-reference arrays reach their
  16/8-element limits and cover absent/live/destroyed lookups. Post-load conversion
  independently checks identities, all action types, maximum stack and signed disable.
- All 781 inventory sequences of up to four items across five categories, checking
  complete item state except independently expected equipped-flag changes.
- Definition defaults: saved health, selective status bits, retreat/resume override
  flags and automatic-spell dispatch. The latter exercises the actual Wizard policy
  with defined/missing definitions, matching/nonmatching types and four flag values.
- Voice records: real linked lookup, unknown-name clearing, empty/255-byte names,
  and absent voice sets. Buff writers cover the shipped iteration order and actual
  shield duration lookup/default. Version/name rejection checks exact stopping points.
  Actual infravision spell application covers versions 1/2, level clamping and timer
  truncation across seventy cases; other spell effects reuse their qualified owners.
- Merchant reads: nine historical boundaries, 0/1/2/60 stock entries, real type lookup,
  versioned item parameters and both vendor tail fields. Shop helpers are already Go.

Captured groups (not every assertion produces a capture):

| Group | Records |
| --- | ---: |
| action-arguments | 144 |
| action-current | 20 |
| action-edges | 12 |
| action-historical | 140 |
| action-references | 27 |
| action-stack | 25 |
| auto-spells | 16 |
| buff-application | 70 |
| buff-writer | 72 |
| current | 2 |
| definition-defaults | 108 |
| equipment-order | 781 |
| merchant | 36 |
| monster-historical | 64 |
| npc-historical | 62 |
| paths | 32 |
| postload | 220 |
| spell-words | 40 |
| timestamps | 100 |
| voices | 6 |

## Compatibility decisions and review notes

1. Timestamp arithmetic wraps at 32 bits. The helper returns the raw sum while
   storing at least 1 using a signed comparison. Several adjustments also run
   while saving; saves have these existing state changes.
2. Action versions with a nonpositive signed value use the oldest layout. Future
   positive versions reject. Unknown action names resolve to action zero.
3. The shared stream adapter returns boolean success, not a byte count. Argument
   serialization returns its final I/O status; an absent final word fails. Unsupported
   argument kinds return the kind value before consuming argument/tail bytes.
4. A reserved action word carries the completed waypoint-loop count. Preserve this
   observed wire behavior, including checksums, rather than writing an assumed zero.
5. Spell application clamps levels above five; the stored 32-bit duration then
   overwrites the resulting timer through its 16-bit field. Existing ordering remains.
6. Definition flags compare exactly with one where C does; other nonzero values
   do not necessarily enable defaults. Preserve health/status masks and ordered
   inventory conflict behavior.

No production defect correction was needed during this baseline. Early failures
were fixture setup or expectation errors: missing AI-name/numeric runtime tables,
wrong stream return assumption, uninitialized fixture file-handle registry, wrong
waypoint field name, and an assumed no-argument idle action (shipped idle uses a
timestamp). Each was corrected before freezing, with existing engine code unchanged.

## Evidence and remaining gates

Pre-freeze all-target/repeat artifacts:
build/port-creature-xfer/c-{default,server,highres,repeat}-development.
Frozen qualification uses c-{default,server,highres,repeat} and c-production-final.
The batch manifest records commands, hashes and whole-source fingerprints.
Development logs remain under build/port-creature-xfer; 48c0ba80's helper recovery
checks passed 1,734 cases and repeated 1,707 records before callback expansion.

Next: commit/push this baseline, then translate the
connected block, review widths/ownership/exports and qualify against these unchanged
expectations plus the complete affected suite and production gates. Matching samples
never replaces the direct C/Go review.

Disk recovery removed only hash-verified duplicate original assets from completed
runs, preserving logs/screenshots and restoration manifests. Original assets and
active runs remain intact. Ignored fixture drafts already installed in source are
consumed; do not reinstall them over reviewed files.

Final baseline timings: default 33.4 s, server 31.9 s, highres 41.9 s, repeat 7.3 s;
production 225.4 s including gameplay/save-load/flat. Client SHA-256:
86f015e11c91ca7ffb93b5657f0d4aa19745fcb4fe3837c167132f2ba7dd9683.
The first production manifest mistakenly requested Go-backed exports for the still-C
callbacks. The existing retained_c option corrected that configuration; no source
or expectation changed. The final production run passes all gates.

Private Go helper drafts are under build/port-creature-xfer/native-*.go. They are
uninstalled and unqualified; review/complete them before use. No C count reduction
is claimed by this baseline commit.
