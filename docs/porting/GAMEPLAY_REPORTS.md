# Next candidate: gameplay status and item reports

Candidate after orchestration qualification: **73 connected routines /1,337 C
section lines** in `legacy/GAME3_2.c`, from4D7BE0 up to but excluding4D9EB0.
This includes player status and resources, inventory/equipment notifications,
team and game-result reports, creature/journal updates, fades, reporting caches,
score/elimination bookkeeping and the aggregate player-report routine. Six nearby
helpers are included because the aggregate calls several of them. Variadic text
formatting after this range is a separate scope.

**The complete C baseline is qualified; production reporting remains C.** The local
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


### Current — complete gameplay-reporting C baseline qualified

The next connected conversion is **73 gameplay-reporting routines /1,337 C section
lines**, GAME3_2.c from4D7BE0 up to4D9EB0. All remain C in this checkpoint. The
baseline is complete: **7,809 cases/28 locked, repeated capture hashes**, plus three
journal-padding regressions and a creature-order probe. Every routine has positive
coverage. Full accumulated default ports pass **81,468 cases/1,010 groups** and
contracts (324.199s); focused reporting/player-control checks pass server112.707s
and highres42.952s. All previous capture hashes are unchanged. First reporting
checkpoint1e20588e and preceding native orchestration3fa8a141 are already pushed.

Physical C remains **101,128 lines /148 files /zero reference C**. The initial
journal prerequisite zeros unused message bytes; its audit permits only254 unused
bytes across nine snapshots, with all defined fields and other state identical.
No reporting algorithm has been converted or retained solely for testing.

Next: commit/push this complete C baseline, then convert the connected batch to
native Go. The ABI audit confirms **45 retained C bridges and28 removable entry
points**, with no production Go address references requiring extra retention.
Replace Go-only calls directly, retire obsolete declarations/getters, switch the
fixture dispatcher to native helpers, compare all28 captures, then run full
accumulated tests on allthree variants, production builds/symbol audits, the known
full-suite comparison and fresh unchanged headless gameplay. Record C LOC,
commit/push, summarize and continue. No pending user question.

Coverage includes byte/word/signed boundaries, all modifier combinations, weapon/
armor precedence, NPC colors and equipment, sparse-player fanout, proximity
falloff, journal name limits, score/frame wrap, zero-max item-health suppression,
full/cooperative reporting caches, NaN/signed-zero comparisons and1,296 elimination
rule cases. Real queues, guarded records, player/team/member iteration and caches
are exercised. Countdown calls are recorded at the existing service boundary;
its two empty fixture text slots are initialized/restored with the ID used by
server.go. Optional recipient masks and rule globals restore per case. Every
report action checks/preserves x87 control on a locked OS thread. Direct-list
messages are in ResourceMessages, which drains the list before later snapshots.

See docs/porting/GAMEPLAY_REPORTS.md. Local evidence and reusable scripts are in
build/port-gameplay-reports: baseline-qualification.json, capture-audit.json,
abi-audit.json, journal-padding-audit.json, qualify-c.py and test-pattern.txt.
Original captures may be raw or SHA-verified gzip (capture-archive.json). Core
fixtures and expected hashes are tracked, so the baseline can be regenerated.

Disk recovery preserves originals and changed files. Ten additional completed
runs have asset restoration manifests (ai-callback, player-attack, inventory,
equipment, resources, object-death, ai-main, rule-command, world-mechanisms,
temporary-updates; each name ends in -port), recovering5,563,589,860bytes. Latest
map-orchestration-port was also deduplicated. Use the existing deduplicate script
with --restore when needed. Preserve nox-iso-from-archive-org.7z.

Standing authorization: continue connected chunks, thoroughly test, document C
LOC, commit/push, summarize and continue until a substantive question or rate
limit. Resolve confident reversible choices and record them for later review.
No new agents. Use build/baseline/env.sh, GOMAXPROCS=2 and go -p 2. Full-suite/crash
logs stay local; report metadata only. SSH push to dbenamy/opennox dev is authorized.
Known separate reviews: orchestration failure paths preserve pending flags/saved
objects; theme cleanup remains shallow. Native orchestration is fully qualified
in3fa8a141, with fresh unchanged headless gameplay passing37.120s.
