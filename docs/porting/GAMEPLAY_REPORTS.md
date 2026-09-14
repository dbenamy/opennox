# Next candidate: gameplay status and item reports

Candidate after orchestration qualification: **73 connected routines /1,337 C
section lines** in `legacy/GAME3_2.c`, from4D7BE0 up to but excluding4D9EB0.
This includes player status and resources, inventory/equipment notifications,
team and game-result reports, creature/journal updates, fades, reporting caches,
score/elimination bookkeeping and the aggregate player-report routine. Six nearby
helpers are included because the aggregate calls several of them. Variadic text
formatting after this range is a separate scope.

**No fixture or production changes for this candidate are applied.** The local
source/signature/callee inventories and an unapplied thin dispatcher draft are in
`build/port-gameplay-reports`. Re-audit callers after virtual removal to establish
retained versus retired C ABIs before constructing the baseline. Initial virtual
removal finds45 routines with outside C references and28 without; inspect Go
getter/callback uses before finalizing that split. Of the28 candidates, nine have
no outside C/Go references and19 have Go-only references.

## Reuse and scope

The existing guarded shop/player/object fixture already captures queued messages,
routing metadata, player/object mutations and normal message-list contents.
`legacy/shop_pools_porttest.go` records queues plus NetList channel1 for player
indices1,7,31. `player_controls_porttest.go` provides actor/target/update records.
A scoped optional reports extension can reuse these owners and action range1800+;
old fixtures must remain byte-identical when the extension is absent.

The direct cache globals are5D4594:1556320/24/28 (TeamBase/SilverKey/GoldKey type IDs).
Use actual type lookup and record cache initialization/reuse. Elimination paths
also use two fixed text resources and game-rule/countdown services. Inventory,
team and player enumeration should use real fixture owners; recorded service
boundaries may cover external lifecycle transitions. No live recipients are needed.

## Baseline requirements

Cover every routine and require positive messages and state changes, including
multi-player fanout and aggregate reporting. Vary recipients, IDs, player/monster/
item flags, signed deltas, resource limits, float rounding, equipment/modifiers,
linked inventories, teams and visibility. Capture routing/order/channel metadata,
all defined message bytes, state/cache/counter mutations, returns, guards and
floating-point control. Suppression-only cases are insufficient.

Repeat original C captures before locking. Audit any concrete undefined bytes or
inconsistent fields and record corrections separately; do not blindly refresh
hashes. Commit/push the qualified C baseline before converting the connected batch.
Then match captures, retire obsolete bridges, run appropriate accumulated and
production/gameplay checks, record actual C LOC, commit/push and continue.
