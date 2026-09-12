# Player controls, respawning and observers

Next connected batch: 56 functions / 1,971 removable physical C lines across
GAME3_3.c and GAME4.c. The address blocks contain 1,973 lines; two standalone
forward declarations remain in C. Include player equipment/stamina helpers,
start selection and respawning, action admission/mapping, scheduled spells,
monster-bot transitions, initial player values/default equipment and observer
selection/owned-unit cleanup. All candidate functions remain original C.

Reuse the guarded actual-player, inventory, ownership, ability and AI fixtures.
Supply coherent player IDs and complete player/monster update records. Record
returns, all affected object/player fields, health/mana protection, created
items, ownership changes, messages, clock and RNG state. Model external transport
as dependency recording while retaining its real state owner. Keep child-list
and inventory-list lifetimes distinct during teardown.

Cover movement/attack state admission, status masks and stamina bounds, action
mapping, start/observer selection with excluded and successful candidates,
respawn equipment and resource restoration, bot transitions and scheduled-spell
ordering. Require successful transitions and allocations alongside denied/empty
paths. Preserve signed/narrow returns and arithmetic stores; establish repeated
original-C captures, lock hashes and commit/push before production conversion.

Convert the connected family and qualify accumulated tests, relevant variants,
production builds, full-suite failure identities and unchanged headless gameplay
at the batch boundary. Update C_LOC/recovery docs, commit/push and continue.
Local source/scope audit: build/port-player-controls. A 56-entry thin dispatcher
draft is staged there; it has not been applied or compiled. Full guarded controls
fixtures and locked baselines remain to be implemented.

## Locked-door notification: approved padding correction (2026-09-12)

Source and compiled-code review found undefined padding in `sub_4FADD0`
(GAME4.c, address 004FADD0). The routine sends a fixed 52-byte message:
bytes 0–1 are F0/21, the NUL-terminated localization key starts at byte 2,
and byte 51 holds the key selector. For a valid string of length L (1–48),
bytes L+3 through 50 are never initialized. The two production callers use
`objcoll.c:GateLockedKey` and `objcoll.c:DoorLockedKey`.

The qualified reward-generation binary confirms the missing stores: relative
to ESP after the 0x54-byte frame allocation, the message begins at +0x18;
selector is stored at +0x4b, the header at +0x18, and checked memcpy copies only
strlen+1 bytes to +0x1a. It then sends all 0x34 bytes. There is no initialization
of the intervening padding. Reproduce this read-only inspection with:

```
objdump -d --disassemble=sub_4FADD0 build/port-reward-generation/bin/opennox
```

The client decoder in client__network__cdecode.c, MSG_GAUNTLET subtype 0x21,
reads the NUL-terminated text at +2 and selector at +51, then consumes 52 bytes.
It does not interpret the padding. The isolated compiled-routine experiment
below additionally confirms the difference between original and corrected output.

The user approved zeroing the padding. The current C routine now initializes
its array with `{0}`; the eventual Go port will use `[52]byte`, preserving the header,
text, selector, recipient, length and transport flags. Compare all defined C
fields and message order against repeated original-C captures. The fixed C
now provides deterministic padding for complete future C/Go comparisons, with
independent assertions that it is zero. Exercise both real keys, valid lengths 1/47/48,
empty/nil/rejected lengths 49/50, non-player/nil recipients, and selectors 0–4.
Do not change string acceptance or other message formats under this decision.

The correction is intentionally separate from the 56-function conversion. No
controls C bodies have been removed. C remains 117,956 lines / 149 files, with
zero reference C. The local scope/source audit was refreshed to include the fix.

### Reproduction and focused validation

The tracked [probe](probes/locked_door_probe.py) extracts the actual routine from
GAME4.c into a temporary compilation unit with a recording transport. It does
not retain a duplicate C algorithm. Run from the repository root:

```
python3 docs/porting/probes/locked_door_probe.py --source-ref f701d9a6
python3 docs/porting/probes/locked_door_probe.py --expect-zero
```

Both versions pass 1,440 cases / 400 messages, repeated in separate processes.
Cases include both real strings, empty/nil and boundary lengths, nil/non-player
recipients, high class bits, selectors 0–4, and recipients including 255. Full
object/update/player inputs are checked unchanged. The defined-field digest is
`f628c861746ff1484f294b803a56c51fc38badbd792037b1d092dd4b3f87308f` in all four runs.
Original output had nonzero padding in 245 messages in this experiment; that
count is incidental stack history and is not a compatibility expectation.
Corrected output has zero padding in every message. Probe compilation uses
32-bit GCC, optimization, fortification and stack protection. This isolates
message construction; the accumulated engine regression is recorded below.
Local evidence: build/port-player-controls/padding.

### Review later

This byte-level change is approved; review it with the eventual Go message
implementation. Keep the fixed 52-byte format and zero-padding assertion.
The user also authorized future decisions when the answer is reasonably clear
and reversal would not require substantial effort: implement, document the
reason/evidence and mark for later review. Ask only when uncertainty or reversal
cost warrants user input. See [DECISIONS.md](DECISIONS.md).

## Additional fixture audit

Reuse the player fixture's real protection records and player registry rather
than approximating stat/observer dependencies. Observer selection needs coherent
unit NetCode and PlayerByID membership. Gold-generation fixtures already showed
why fractional and large-plus-small inputs matter: player-value interpolation
and bolt-damage captures must likewise include non-exact arithmetic and store
boundaries. Keep the signed level byte within valid XP-table caller bounds.

Bot fixtures need full guarded 0x898-byte update records, with player data +292
pointing to the bot data and bot +2180 back to the player update data. Capture
both records across morph/restore and normalize their identities. Exercise
all three classes, deterministic RNG choices, all action mappings, and actual
creation plus respawn after/before the two-second threshold. Restore update
pointers before teardown. Waypoint fixtures need all three pointer slots and
both index bytes, with positive near-delete and far-move paths.

Scheduled-spell fixtures must distinguish front removal from stack removal:
front removal shifts entries and zeros the vacated slot, while queue removal
only decrements the count. Record admitted/rejected cast dependency arguments,
text/audio requests, full queue bytes, and integer-coordinate-to-float bits.

Observer baseline implementation notes: distinguish global object traversal
from owned-child traversal. `sub_4E5F40` queries the next object after requesting
deletion; `sub_4E5FC0` saves child/inventory next pointers before deletion.
The observer helpers wrap from the current target to the beginning and skip
dead/status-filtered candidates. Include successful wraparound, GameBall mode
preference, and no eligible candidates. The slave-next helper additionally
requires a non-nil owner before searching siblings. Transfer uses the actor's
owner as the new owner; removal clears every child owner/next pointer and the
actor's head. These are distinct effects and should have complete list snapshots.

### Padding-fix engine regression

The full accumulated port-test selection passes in default/server/highres:
180.175s / 167.954s / 179.293s, with existing expectations unchanged. This checks
all 41,354 previously focused cases along with the other accumulated groups.
The separate message probe covers 1,440 additional cases; these are not yet
part of the Go test fixture count. This one-line correction was not a new
production-build/full-suite/gameplay milestone; those last passed at the reward
conversion. Evidence: build/port-player-controls/padding/qualification.json.
