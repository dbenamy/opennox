# Direct Go book drag/drop callbacks

Ten call sites in the icon, page and visibility handlers now call the existing
Go callback variables directly, eliminating Go→C→exported-Go round trips.
The three callback assignments and exported entry points stay intact.

Seven clear calls have no arguments or return conversion. Both save calls retain
signed 32-bit ID normalization (`int(int32(id))`); the getter comparison normalizes
its result to int32. The configured ABI is 386, where Go int and C int agree.
This does not establish a 64-bit port. There are no pointer ownership transfers,
errno/multiple-return operations, macro substitutions or callback-address changes
in this batch. Existing quickbar paths already use these Go callback variables.

Luna inventoried the candidate and independently reviewed test gaps. The primary
verified the export wrappers, root initialization, C headers and all ten sites.
Luna's supplemental test exercises actual client drag state through the real
book hide handler. It verifies type1 clears and type0/type2/-1/int32 extremes
remain, plus blocked and already-hidden exits. Existing frozen captures are not
regenerated. High-bit drag IDs are not injected through unrelated lookup tables;
32-bit normalization is explicit in production.

All 25 spellbook roots and unchanged 24 frozen captures pass in
default/server/highres; the new independent state contract also passed before
conversion. Safe/static, fresh production/ABI and headless creation/save-load
pass. The full suite matches exactly: 304 existing failure events, 17 passing /
2 failing / 32 skipped packages. All four binaries retain the three C exports
and omit their redundant generated C-call bridges. There are no new skips.

Baseline: `2afcb05a` ([record](book-direct-calls-c-qualification.json)).
[Conversion qualification](book-direct-calls-qualification.json) records source
identity, commands/results, exact suite comparison, scenarios and binary hashes.
[Batch manifest](book-direct-calls-batch.json) defines the reproducible checks. The initial baseline uses the
historical broad affected selection (138 default/highres roots, 137 server roots); post-change uses the bounded book
selection because unrelated consumer tests already have qualified evidence.

Standalone production and test-reference C remain zero files/lines; production
C preamble bodies remain79. This removes generated call bridges, not handwritten
C algorithm bodies, and makes no measured whole-game performance claim.


## Storage

Completed immutable binaries with identical SHA-256, mode and ownership now
share storage through hard links. Host PID1/process/open-file checks preceded
replacement, and every original path/content remains verified. Ten duplicate
paths recovered 493,230,372 bytes. Original metadata and per-file progress are in
`build/port-book-direct/duplicate-binaries-shared.json` and
`duplicate-binaries-extra-shared.json` alongside it. To deliberately mutate a
retained artifact in the future, first copy it to a temporary sibling and replace
the path, separating its inode; never overwrite a shared inode in place. Prefer
this verified deduplication to evicting useful caches when possible: the previous
cache cleanup required a server test rebuild during this batch.


Two historical raw captures were losslessly gzip-archived after independent
hash/stat review and host open-file checks: 469,204,794 raw bytes became
10,662,396 compressed bytes. Full round trips verified before original retirement.
The canonical files have 8 and 9 historical symlink aliases; restore canonical
paths before reusing those aliases or old qualification scripts. Record, hashes,
metadata and alias lists: `build/port-book-direct/capture-archive-record.json`.
The pause probe was successful; the host capture's own test passed, but its
broader original run later failed an unrelated pause test. Archival preserves
that provenance and does not reclassify a failed run as successful. No current
qualification depends on these raw captures.

Seven older votes/item-respawn test logs were also archived with the same checks:
278,937,102 raw bytes became 17,386,546 gzip bytes. Full log contents remain
recoverable from their sibling `.gz` files. Paths, hashes and restore metadata
are in `build/port-book-direct/log-archive-record.json`. These are completed
historical outputs, not the active production suite log.
