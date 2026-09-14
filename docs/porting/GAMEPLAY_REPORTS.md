# Gameplay status and item reports

Candidate after orchestration qualification: **73 connected routines /1,337 C
section lines** in `legacy/GAME3_2.c`, from4D7BE0 up to but excluding4D9EB0.
This includes player status and resources, inventory/equipment notifications,
team and game-result reports, creature/journal updates, fades, reporting caches,
score/elimination bookkeeping and the aggregate player-report routine. Six nearby
helpers are included because the aggregate calls several of them. Variadic text
formatting after this range is a separate scope.

**The native conversion is fully qualified. Original C is recoverable at fbfac731.** The local
source/signature/callee inventories and thin dispatcher source are in
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


### Native conversion — qualified

Baseline `fbfac731` is committed and pushed. All73 implementations now use Go;
45 entry points retain their required C ABI and28 have been retired. All19
Go-only callers now invoke native helpers, including equipment, resources,
attacks, spells, player controls and inventory paths. The fixture dispatcher
calls native helpers too. A whole-source audit found no remaining references to
the28 retired names. No C algorithm remains solely for testing.

The first native focused run passed in20.984s: all7,809 cases /28 hashes match the
original C baseline, plus journal-padding and creature-order checks. Accumulated
default/server/highres tests passed (77.441s/114.117s/17.081s wall), covering81,468
cases /1,010 groups and contracts in every variant. Allthree production builds
pass ELF386/SSE2 and symbol audits:28 retired symbols absent,45 required bridges
present, no test helpers. The full suite matches all1,553 baseline failure entries
exactly (15pass/3fail/32skip packages). Fresh unchanged headless gameplay passes
35.208s, with overrides disabled and null audio. No qualification failures were
found in this conversion. `qualification.json` gathers the local evidence.
Measured physical C is **99,791 lines /148 files /zero reference C**, a reduction
of1,337 lines. Local evidence: native-first.log, native-retirement-audit.json,
variants.json and native-qualification-progress.log.

Serialization preserves message bytes, recipient/ordering/priority, narrow signed
returns, full-width caches and mutation ordering. Engine objects and queued
records remain C-owned. Temporary Go byte buffers pass to retained queue routines
that synchronously copy them into owned storage; no pointer is retained. Direct
message-list insertion uses the existing Go queue. Player, team and inventory
iteration use the existing native owners, with the original valid-domain behavior.
The original baseline has1296 elimination-rule cases, plus frame wrap, repeated
aggregate caches, NaN/signed zero, equipment modifiers and byte/word boundaries.
All report actions still assert preserved x87 control on a locked OS thread.

Capture archival has SHA-verified gzip round trips; the latest pass archived208
additional files and recovered2,472,958,742bytes. Captures and generated scripts
are local evidence; tracked tests and baseline Git history reproduce the oracle.
