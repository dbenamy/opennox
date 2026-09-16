# Spellbook icon-release baseline supplement

Original-C behavior is qualified for 480 additional icon-release records, bringing
spellbook coverage to **4,694 records in 24 roots**. The original 23 capture hashes
are unchanged. [The capture list](spellbook-captures.json) and
[batch manifest](spellbook-batch.json) include the supplemental frozen hash.

This was measured in an isolated checkout of original-C baseline `337f7c77`, with
only `spellbook_icon_release_porttest_test.go` and `legacy/book_cursor_porttest.go`
added to its source tree. The supplemental baseline commit deliberately leaves all
native production changes unstaged, so checking out that commit still runs C.
No production correction or replacement C algorithm was introduced for these tests.

The cases use actual book, quickbar, cursor, ability, input-timeout and message-list
owners. They cover aimed/immediate abilities, aimed/immediate/deferred spells,
creature-guide entries, existing spell/ability cursors and queued casts, inclusive
window edges, outside/quickbar drops, release events 6/7 and unrelated events 8/12.
Independent contracts check drag cleanup, outside capture/selection cleanup,
immediate ability message bytes and non-release stability. Complete normalized
book/quickbar/cursor state and queued message bytes are frozen.

`c-icon-release2` passes the initial single root in 118.287 seconds. The earlier
launch passed a regex to the pattern-file option, was rejected before discovery,
and remains preserved. The finalized original-C qualification runs all 24 roots,
not just the new one, to check fixture cleanup against the existing corpus:

| Configuration | Roots | Result | Seconds |
| --- | ---: | --- | ---: |
| Default | 24 | pass, no skips | 60.552 |
| Independent repeat | 24 | pass, no skips | 43.822 |
| Server | 24 | pass, no skips | 137.033 |
| Highres | 24 | pass, no skips | 61.513 |

Total manifest time: **303.563 seconds**, source unchanged. Evidence is preserved
under `build/port-book/c-supplement-qualified`, including exact source fingerprints,
commands and the asset-path-adjusted manifest. The isolated worktree is temporary;
Git history preserves the full original-C oracle after it is removed. Native
qualification is separate and does not regenerate these expectations.
