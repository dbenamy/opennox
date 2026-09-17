# Reliable game-message queue

## Scope and status

C baseline qualified after object-report conversion `eb2e61c4`; native conversion next.
Twenty live functions / 615 physical C lines in GAME3_3.c, 004E4DE0 through
004E5770, ending before the unrelated sub_50B510 forward declaration. Six private
helpers have live selected callers; preliminary whole-source audit found no orphan.
The prerequisite C correction brings production C to 55,584 lines / 82 files,
zero reference C.

## Contract plan

Reuse the actual sparse players, reliable allocation/list owner, NetList and
frame/rate state. Snapshot payloads, recipients, sequence values, acknowledged and
sent masks, retry state, timestamps and linked-list order. Normalize only related
object addresses to fixture identities; assert actual next/previous/tail links.

Cover initialization and resets, insertion/unlink/trim/coalescing, payload capacity,
sequence wrap, broadcast/exclusion masks, related object-known bits, acknowledgement,
ordinary/ordered framing, retry timing/budgets, queue capacity and rate adaptation.
The delivery paths append to the real game-message queues; no sockets are needed.
Own the actual reserved-byte count rather than substituting a delivery callback.

Capture new C contracts on all targets and repeat independently. Reuse the prior
qualified production baseline while production source remains identical. Preserve
existing consumer expectations and run the full affected corpus and fresh
production/integration checks after native conversion.

The existing Go sequence-reset helper in src/player.go duplicates sub_4E4F30;
compare before unifying it. Include C preambles in bridge-retirement audits.
The pressure ownership correction below is the sole intended behavior change.

## Baseline correction for review

An independent pressure-cleanup regression reproduces a SIGSEGV in original C
(`build/port-reliable-reports/pressure-original-c.jsonl`). The scan selects an
oldest message, then removes a slow player's messages; that removal can free the
selected node. Cleanup subsequently reads the freed links and frees it again.

The corrected C implementation searches the surviving list for the selected node
before unlinking it. It still returns success if player removal already released
that node. This preserves recipient acknowledgement and related-object bit updates,
unlike simply unlinking the oldest message before removing the player. Regression
coverage checks same-player selection, surviving oldest messages for other players
and broadcasts, status mutation, list integrity and subsequent pool reuse.

This is an authorized reversible baseline bug fix, to be carried into Go. It adds
7 C lines temporarily (55,584 production lines / 82 files); the corrected selected
block is 622 lines. No reference C is added. Because production C changed, fresh C
production/integration gates replace the originally proposed baseline reuse.

## C contract coverage

Fifteen roots / 6,989 leaf cases, with 15 capture groups / 9,032 records. Every
selected entry is exercised through the real allocation, players, frame/rate state
and NetList. All three targets and the independent repeat pass frozen checks. Captures normalize related-object addresses only; list and pool
ownership is checked directly. No network connection or substitute callback is used.

Coverage includes 0/1/149/150/151-byte enqueue bounds; ordinary, broadcast and
exclusion audiences; payload ownership; sequence wrap and insertion order; reset
preservation; all acknowledgement classes and known-object mutations; real delivery,
retry budgets/countdowns and delivery timestamps; host/client framing; queue byte
capacity and reserved bytes; allocation limits/dynamic growth; pressure selection;
rate thresholds/adaptation; recipient removal; host gating and replay acknowledgement.

Preserved quirks include empty exclusion audiences returning success before length
validation, zero-length enqueue allocation, signed recipient interpretation,
equal-frame pressure selection favoring the newest node, the 999999999 pressure
sentinel, and the host wrapper reporting success even when its queue is full.

C baseline reproduction: source build/baseline/env.sh, then run
`python3 tools/porting/run_batch.py docs/porting/reliable-reports-batch.json --phase c-default --out NEW_DIRECTORY`
(and c-server/c-highres/c-repeat). The committed C production manifest is
`docs/porting/reliable-reports-c-production.json`; use its production phase while
checked out at the C baseline. The first local C production run used an equivalent
ignored manifest before this recoverable copy was added. The latter also explicitly
checks presence of the seven C interfaces retained by the future Go conversion.


## Qualified C baseline evidence

`build/port-reliable-reports/c-{default,server,highres,repeat}`: 15 roots /
6,989 leaf cases per target, zero skips, all 15 groups / 9,032 records identical.
Fresh `c-production` passes all three builds/ABI audits, the exact known full-suite
multiset (1,553 entries; 15 packages pass, 3 fail, 32 skip), 41-frame gameplay,
7-frame actual save/load and 14-frame flat rendering with exact map regeneration.
The seven future retained C interfaces are present in all three baseline binaries.

All 1,988 source fingerprints agree across these gates. The corrected C client
SHA-256 is `990ec5fb1e033ec9de17a5138608807003532a3dcc248f53027217496b5e460a`.
The original crash capture remains separate from the corrected frozen contracts.
No C algorithm is added solely for tests; this committed baseline provides recovery.
