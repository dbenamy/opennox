# Ordered client message queue

The batch replaces six live queue helpers in GAME2_3.c and removes its inert
console walk. Callers are client dispatch, local/remote receive loops, session
initialization/cleanup and the console command. The 32-byte foreign record layout
and existing list owners remain stable; the private implementation moves to Go.

## Original C baseline

[Qualification](client-sequence-c-qualification.json) records nine roots in each
profile and two independent default captures. Four captures contain 1,699 cases:
624 existing insertion cases plus 1,075 delivery/length/clock observations. The
existing insertion hash is unchanged. The preceding production qualification is
reused because only one test file was added; every preceding source hash matches.

Contracts cover ordering, duplicate handling, gaps, strict timeout boundaries,
sequence wrap, clock narrowing and unsigned elapsed time, every payload length,
callback index/order and ignored callback results. Real queue/list/allocation
owners are used; only the packet-delivery callback is observed and restored.

## Decisions to review

Preserve two surprising C behaviors: polling a pending gap before its timeout
clears its retry cursor, and future-only insertion never creates that cursor.
Clock subtraction promotes the truncated clock to uint64, so backward movement
or wrap can immediately expire a gap. These are compatibility decisions, not
new timeout behavior. Initialization does not free an existing queue; cleanup
frees nodes but leaves cursors until reset. Preserve the caller lifecycle.

The console walk only traverses the list and returns null; its caller discards
that value. Removing the walk preserves the command's successful return and lack
of output. No original C is retained solely for tests.

Native implementation and qualification are pending.
