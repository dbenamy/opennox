# Player controls, respawning and observers

Connected batch: 56 functions / 1,971 removable physical C lines across
GAME3_3.c and GAME4.c. The address blocks contain 1,973 lines; two standalone
forward declarations remain in C. Include player equipment/stamina helpers,
start selection and respawning, action admission/mapping, scheduled spells,
monster-bot transitions, initial player values/default equipment and observer
selection/owned-unit cleanup. Baseline `1a87b410` was pushed before conversion. All 56 functions are now Go;
qualification is complete.

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
Local source/scope audit: build/port-player-controls. The 56-entry thin dispatcher and guarded fixture are implemented; the locked
baseline below is ready for conversion.

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

## Locked C baseline (2026-09-12)

All 56 entries are covered by 3,205 cases / 58 complete capture groups. Locked
controls tests pass (9.785s), and the shared resource/inventory/equipment/effects/
world/objectives/attack/projectile/damage/state/reward regression passes (109.281s).
Every controls capture matches byte-for-byte across these separate runs. The
baseline includes both documented initialization corrections; all 56 algorithms
still run in C. Commit and push this checkpoint before production conversion.

Coverage includes all 256 equipment availability masks, fractional stat and bolt
arithmetic, signed stamina and narrow returns, all three bot classes, actual bot
allocation/update, paired morphing, corpse-part creation, real default equipment,
observer wraparound, waypoint movement/deletion, and admitted/rejected spell
queues. Independent positive checks require allocations, known health/mana,
glyph counts, inversion callbacks, queue mutations and zero fifth modifier words.

The first locked run caught one missing callback identity: observer exit installs
`nox_xxx_updatePlayer_4F8100` at object offset 744. Relinking moved its address.
The fixture now normalizes this exact function to identity 91601; the four copies
of that field in each snapshot are checked, rather than ignored. All other capture
bytes were unchanged. Stable evidence: c-locked-*.json and c-confirm-*.json, with
58 expected hashes committed in player_controls_porttest_test.go. Preliminary
c-first/c-repeat files predate this final normalization and are not the baseline.

Actual root ability cancellation and monster updating are bound to the fixture's
server owner; this avoids invoking an unrelated global server or substituting
the algorithms. Guarded player/stat protection records now optionally include
strength, speed, character-name checksum, XP, level and ability checksum fields.
Created inventory items are adopted for capture and teardown without inventing
placement or mutating their state. The retained placement callback is recorded.

### Review later: default-equipment fifth modifier word

`nox_xxx_playerMakeDefItems_4EF7D0` filled only offsets 0, 4, 8 and 12 of a
20-byte local modifier record. `stateAttributes` copies all 20 bytes, including
the uninitialized word at +16, into created equipment. The C initializer now
zeros the whole record before populating descriptors. This deliberate behavior
change follows the user's standing instruction for confident, inexpensive-to-
reverse decisions. It is not an exact preservation of indeterminate bytes.
Keep the Go record zero-initialized and assert created default equipment's fifth
word is zero. Review alongside the already-approved locked-door padding fix.
The C baseline for this batch will include both corrections. C LOC is unchanged.

Current evidence: build/port-player-controls/c-locked.log, c-regression.log,
c-locked-*.json, c-confirm-*.json and baseline-hashes.json. Native drafts remain
in the ignored build directory until the baseline checkpoint is pushed.

## Completed native conversion

All 56 routines now live in player_controls_{helpers,stats,movement,respawn,bot,
messages,exports}.go. The 3,205 cases / 58 complete native captures match the
locked C baseline byte-for-byte (10.976s). Direct Go caller connections are included in the passing accumulated
qualification below. Production C
is 115,985 lines / 149 files / zero reference C: 1,971 lines removed, retaining
the two adjacent forward declarations. No C algorithms are retained for tests.

Native movement/input/observer wrappers, attack and projectile damage callers,
and spawned-object cleanup call the new owners directly. Default equipment calls
native modifier/resource owners; protected XP uses updateProtectionFloat. The
initial port mistakenly used IEEE bits for XP protection; exact captures caught
that mismatch, and the existing numeric encoding was restored without changing
expected hashes. Spawn eligibility retains actual team-list membership checks.

C-owned bot records remain compatible with the existing object teardown. Legacy
signed-byte pointer returns and short respawn returns retain their original low
bits. The C declaration of sub_4FADD0 drops its const qualifier to match cgo's
export declaration; the implementation reads the supplied string. Defined valid
keys remain limited to 1–48 bytes. Original overlong inputs that overflowed the
fixed array are outside the preserved contract; the Go helper rejects them.

Both initialization corrections are retained: locked-door padding and the fifth
modifier word are zero. Review-later decisions above remain explicit. The fixture
also verifies a caller-supplied nonzero fifth modifier word is copied normally.

Evidence: native-final-*.json, native-final.log and c-locked-*.json under
build/port-player-controls. Preliminary and confirmation captures were losslessly
compressed as .json.gz to recover VM disk space; the locked C and final native
captures remain raw. No performance claim is based on fixture execution time.

Accumulated tests, including 44,559 focused cases / 422 groups, pass in default,
server and highres: 192.132s / 179.762s / 196.813s. All three production binaries
build and are verified ELF32/i386, SSE2, CGO enabled. The full suite, run with
NOX_DATA set to the extracted assets and count=1, has the exact known failure
multiset: 1,553 entries; 15 packages pass, 3 fail, 32 skip. Fresh unchanged repeat-a
headless gameplay passes in 42.972s. No expectations were overridden.

All 56 production algorithms are native. The final caller audit found 30 C exports
that became unused after connecting native callers. Retire those in the immediate
follow-up, switching their fixture dispatch to direct Go calls while preserving
all 58 hashes. This is ABI cleanup, not retained C reference algorithms. The
next algorithm batch is spell casting and buff lifecycle: 29 blocks / roughly
1,300 C lines; its source/declaration audit is staged in build/port-spell-lifecycle.

Final evidence: qualification.json, binary-verification.json, ports-*.log,
full-suite-comparison.json and build.log under build/port-player-controls;
gameplay result under build/baseline/runs/player-controls-port. Commit/push this
qualified conversion, summarize and continue the 30-export cleanup without asking.

## Completed export cleanup

The qualified conversion was committed/pushed as `0a3f446d`. Its 30 obsolete
exports and header declarations are now retired. Their test operations call the
native owners directly; the thin C dispatcher covers only retained production
ABIs. All 3,205 cases / 58 original hashes still pass in default/server/highres
(9.281s / 9.427s / 14.374s). Function-identity normalization skips absent symbols,
so nil remains nil. No production algorithm or expected capture changed.

All three production targets rebuilt and were verified ELF32/i386/SSE2/CGO;
`nm` confirms all 30 symbols absent from each binary. Fresh unchanged repeat-a
headless gameplay passes in 34.974s. The preceding conversion's accumulated
matrix and exact full-suite failure comparison remain the algorithm milestone;
this ABI-only follow-up reran the affected corpus across variants, all builds,
symbol checks and gameplay. C remains 115,985 lines / 149 files / zero reference C.
Evidence: build/port-controls-bridges and baseline/runs/controls-bridges-port.

Next: spell casting/buff lifecycle, 29 functions / 1,326 removable lines. Audit
unused exports before qualification in future batches, avoiding a second build
cycle for cleanup that can be included in the conversion itself.
