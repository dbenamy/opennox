# Object initialization and reward generation

Completed connected batch: 21 address blocks / 1,862 physical C lines in
GAME3_3.c, from 004F0390 through 004F2210. Includes simple object initializers,
generator setup, weighted reward categories and tiers, spell/ability/field-guide
books, armor/weapon modifiers, potions/gold/gems and marker selection/placement.
All candidate production functions are now native Go; baseline history follows.

Use the existing guarded object/owner fixture and real object allocator, type
registry and RNG. Add guarded reward-marker data, definition and modifier-table
inputs, and enough initialized created-object storage for books and modifiers.
Record full created objects, modifier assignments, placement/deletion requests,
input guards, returned identities, globals and RNG position. Use actual shipped
constant tables where relevant; make definition availability and class masks
explicit fixture inputs. Avoid empty-table-only coverage: require successful
creation for every reward family and each modifier slot, plus absent/filtered
entries, explicit book choices and marker probabilities.

Test tier boundaries and weighted-selection endpoints across seeded runs,
including the original odd-tier table overwrite. Check class restrictions,
duplicate-modifier suppression, masks, creation failures, quest stages, player
counts and marker flags. Separate invalid caller inputs from supported edge
cases. Repeat complete original-C captures in separate processes, lock hashes,
commit/push before any production conversion. Then port the connected batch,
run accumulated variants/builds/full-suite/headless qualification once, update
C_LOC and recovery docs, commit/push and continue.

Local audit: build/port-reward-generation/candidate-scope.json and
candidate-source.txt. Guarded fixtures and the initial corpus are implemented; original-C runs are
in progress. Initialization (288 cases) and tier selection (896 cases) passed.
The first factory run exposed missing fixture ownership of unplaced returned
objects; the dispatcher now registers their actual returned allocation for
capture/teardown without synthesizing placement or mutating object state.
No production reward conversion has started.

## Original-C baseline

All 21 production functions remain C. **6,926 cases / 48 complete capture groups**
repeat byte-for-byte in separate processes and their hashes are locked in
`src/reward_generation_porttest_test.go`. The combined regression passes
**41,306 cases / 363 groups** in 122.877s. Baseline-only locked rerun passes (20.166s).
Production C remains **119,818 lines / 149 files / zero reference C**.

The corpus exercises initializers, tier and category boundaries, explicit book
lists (including bytes other than 1), missing/filtered definitions, every reward
factory, modifier eligibility/fallback and duplicate IDs, wand replenishment,
gold/XP scaling and marker placement. Captures include 288 successful armor
creations, 288 weapon creations, all four modifier slots, 144 Ankh placements
and 16 cases requesting potion deletion. A permanent positive test checks an
actual four-modifier weapon and an explicitly selected ability-book payload.

Fixture storage uses guarded 256-byte marker inputs, guarded created init/use
buffers and the shared 64-byte update payload/16-byte guard. Factory-returned
objects are registered for capture and teardown without inventing placement
calls. Definition/type tables and gold constants are restored after each case.
No C algorithms were copied into the fixture. Sparse modifier masks and equal
identifier bytes on different descriptor pointers distinguish mask filtering
and duplicate suppression from superficially similar implementations.

Local evidence: build/port-reward-generation/c-final-*.json, c-confirm-*.json,
baseline-hashes.json, positive-coverage.json and combined-c-regression.log.
Commit and push this baseline before production conversion. Expected batch C
reduction: 1,862 lines, leaving 117,956 (subject to the final physical count).

### Gold arithmetic addendum before cutover

The baseline was pushed as `9b211b92`. Review found that the initial gold XP
cases mostly had exactly representable averages. Added 48 cases for fractional
averages and large-plus-small XP sums, repeated against the still-original C
(0.282s / 0.223s), with a locked hash. Total: **6,974 cases / 49 groups**.
The two draft Go reward helpers are not yet exported or called by production C;
all 21 C bodies remain unchanged until this addendum is committed and pushed.

## Native conversion

Baseline `9b211b92` and gold-rounding addendum `9caa2c85` were pushed before
C cutover. All 21 functions are now native: weighted selection, book/item
factories, modifier filtering, marker placement and object initializers.
Shared helpers preserve table order, original RNG draws, odd-tier overwrite,
modifier fallback and ID-based duplicate suppression. Shop callers use the
native marker helper directly. The first native run matches all **6,974 cases /
49 complete captures** byte-for-byte (21.068s), including gold rounding.

Exactly **1,862 C lines** are removed, leaving **117,956 / 149 files / zero
reference C**. Full qualification passed as recorded below.
Confirmation and preliminary baseline JSON files were losslessly compressed to
`.json.gz` to save space; c-final and native-first captures remain directly
readable. compare-captures.py supports the compressed confirmation files.
Next planned batch: [player controls, respawning and observers](PLAYER_CONTROLS.md).

## Qualification

All 6,974 cases / 49 complete captures match original C byte-for-byte (21.068s).
Accumulated tests, including 41,354 focused cases / 364 groups, pass in
default/server/highres: 193.708s / 169.508s / 178.540s. Three production binaries are verified
ELF32/i386, GO386=sse2 and CGO enabled. Full-suite failure identities and
multiplicities match exactly: 1,553 entries; 15 packages pass, 3 fail, 32 skip.
Fresh unchanged repeat-a headless gameplay passes in 34.990s
using Xvfb and null audio, without updating expectations.

Evidence: build/port-reward-generation/qualification.json and sibling captures,
logs and binaries; build/baseline/runs/reward-generation-port. Next batch:
[player controls, respawning and observers](PLAYER_CONTROLS.md).
