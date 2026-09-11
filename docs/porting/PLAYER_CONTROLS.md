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
Local source/scope audit: build/port-player-controls. Fixtures have not started.

## Locked-door notification: pending compatibility decision

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
It does not interpret the padding. This is source/disassembly evidence; no
runtime packet-instability experiment has been performed or claimed.

Proposed conversion: use a zero-initialized `[52]byte`, preserving the header,
text, selector, recipient, length and transport flags. Compare all defined C
fields and message order against repeated original-C captures; exclude only
undefined padding from the differential expectation and add independent Go
assertions that it is zero. Exercise both real keys, valid lengths 1/47/48,
empty/nil/rejected lengths 49/50, non-player/nil recipients, and selectors 0–4.
Do not change string acceptance or other message formats under this decision.

The user has been asked whether to accept this documented byte-level fix or
defer this one routine. No production controls conversion or baseline lock has
started; resolve that choice before converting this routine. The rest of the
connected baseline can proceed independently.

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
