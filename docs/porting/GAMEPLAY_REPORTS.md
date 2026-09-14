# Next candidate: gameplay status and item reports

Candidate after orchestration qualification: **73 connected routines /1,337 C
section lines** in `legacy/GAME3_2.c`, from4D7BE0 up to but excluding4D9EB0.
This includes player status and resources, inventory/equipment notifications,
team and game-result reports, creature/journal updates, fades, reporting caches,
score/elimination bookkeeping and the aggregate player-report routine. Six nearby
helpers are included because the aggregate calls several of them. Variadic text
formatting after this range is a separate scope.

**The first C baseline checkpoint is qualified; production reporting remains C.** The local
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


### Current — first gameplay-reporting baseline checkpoint qualified

Native map orchestration is complete and pushed in **3fa8a141**: five routines,
−207 C lines; all73,659 accumulated cases/982groups and contracts pass all three
variants, production builds and symbol audits pass, the full suite matches its
1,553 known failures, and fresh unchanged headless gameplay passes37.120s.

The next connected scope is **73 gameplay-reporting routines /1,337 C section
lines** (GAME3_2.c,4D7BE0 through4D9E70). All still run in C. The first baseline
checkpoint covers41 routines with **2,064 captured cases/six locked hashes**, plus
three journal regressions and a creature-order probe. Original-C captures repeat
exactly. Reporting plus all existing player-control tests pass three default
runs46.352s, server16.373s and highres15.493s; old control hashes are unchanged.
These are scoped baseline checks; full application qualification follows the
completed baseline/conversion. See docs/porting/GAMEPLAY_REPORTS.md.

A separate reversible prerequisite initializes the three journal message buffers.
The original regression reproduces nonzero unused bytes in add/remove/update.
The recursive audit permits only unused message-byte changes: nine snapshots,
254bytes; all defined fields and all other state are identical. Corrected tests
pass. No reference C algorithms were added. Physical C remains **101,128 lines /
148files /zero reference C**, confirmed by tools/porting/c_loc.py.

Next: expand the remaining32 routines and deeper branch/integration coverage,
then qualify/commit/push the complete C baseline before converting all73. Unapplied
root test stages results.go.stage and notifications.go.stage under
build/port-gameplay-reports cover winner broadcasts, score changes, creature
resource cascades and additional notifications. Audit and apply after this
checkpoint is committed. Important: direct-list messages are in ResourceMessages;
that existing snapshot drains NetList before later snapshots. Journal record
lists explicitly reset per case. The optional reporting fixture is isolated from
old cases; its C adapter contains calls only. Initial ABI inventory finds45 outside
C callers and28 retirement candidates; recheck Go getter/callback references.

Local evidence: build/port-gameplay-reports/initial-capture-audit.json,
journal-padding-audit.json, c-initial-{default,server,highres}.log and
initial-qualification.json. Historical captures/binaries have SHA-verified
compressed archives; completed gameplay asset copies have restoration manifests.
The latest map-orchestration-port asset copy is now deduplicated (556,358,986bytes
recovered); originals and changed run files remain. Preserve the 7z asset archive.

Standing authorization: continue connected chunks, thoroughly test, document C
LOC, commit/push, summarize and continue until a substantive question or rate
limit. Resolve confident reversible choices and record them for later review.
No pending question; no new agents. Use build/baseline/env.sh, GOMAXPROCS=2 and
go -p 2. Full-suite logs stay local; report metadata only. SSH push to
dbenamy/opennox dev is authorized. Known separate reviews: orchestration failure
paths preserve pending flags/saved objects; theme cleanup remains shallow.
