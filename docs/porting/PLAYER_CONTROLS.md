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
