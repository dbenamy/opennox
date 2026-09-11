# Network alias table and exhaustion fix — 2026-09-11

Scope: reset 57B920, slot selection 57B9A0 and record write 57BA10.
Original-C helper baseline: 19d02832. Native helpers and the approved caller
fix are complete and validated.
The table has 255 eight-byte records: two uint16 keys and a uint32 frame.
Storage belongs to the client blob or embedded Player.NetData16, not these helpers.

Selection starts at the low byte of the full-width first key, mapping 0/255 to 1.
It scans slots 1 through 254, accepting a matching key pair OR a record with
Frame4 strictly less than the current frame. Equality remains occupied for
mismatching keys. It compares zero-extended stored keys against full-width signed
C int inputs, while the writer truncates keys to uint16. A full cycle returns
char -1, whose raw byte is 255. Reset clears exactly 2,040 bytes and returns 0;
write returns its original pointer bits. Preserve those live C ABI contracts.

The existing root Go sub_57B930 is an outbound lookup for matching unexpired
aliases, not the same operation as allocating an incoming alias slot. Keep its
semantics separate; do not use it as an oracle or consolidate it blindly.

## Confirmed exhaustion bug; fix approved 2026-09-11

Both GAME2_3.c callers (494A60 and 494C30) store the selector result in unsigned
char v24, then test v24 != -1. Integer promotion makes that condition always
true. On exhaustion, both write an eight-byte record at table + 8*255, immediately
past its 2,040-byte extent, and send MSG_NEW_ALIAS with reserved alias 255.

Disassembly of the current 386 production binary confirms that neither caller
branches after the selector: it zero-extends AL, computes table + 8*index and
calls 57BA10 unconditionally. Local inspection artifacts are under
build/port-ping-aggregate/alias-caller-494A60.asm and alias-caller-494C30.asm.
This is a logical table overrun; no claim is made about the ownership of the
following bytes in the client blob without further investigation.

The user approved fixing exhaustion during the port on 2026-09-11.
The proposed fix compares against the byte sentinel 255 and skips the table
write and alias announcement when full, while continuing packet processing.
Proceed with the fix and actual-caller regression coverage.

Helper coverage: original-C selection across every start byte, full-width and
negative keys, matching versus expired records, exact frame equality, wrap and
full-table failure; guarded C-owned storage for exact reset/write extents, key
truncation and return bits. Actual caller coverage includes
full tables and successful aliases including 128..254, checking no write or
announcement on failure and unchanged normal packet processing.

## Original-C helper baseline — 2026-09-11

24,576 selection cases cover every starting byte, normal/out-of-range/negative
key combinations, four frame boundaries, matching entries, wraparound, expiry,
equality and slot-zero exclusion. An independent circular-distance oracle ranks
all eligible slots; exact signed-char results and table immutability are checked.
2,040 writes cover every slot and eight key/frame edge patterns, each followed by
a reset (alternating direct ABI and Go wrapper). Checks include key truncation,
returned pointer bits, all table bytes, and 16-byte guards on both sides.
All **28,656 helper operations** pass against original C before replacement.
Production C remains **141,126 physical lines**, 153 files, zero reference C.
Local artifacts: `build/port-network-alias/`.

## Actual-caller regression and approved fix

1,530 cases call the real C packet handlers (both streams × three frame values ×
full-table plus each of 254 usable slots). A minimal server owns a real NetList;
a C-owned drawable proxy records sprite creation and subsequent frame/animation
updates. Check exact consumed lengths, decoded positions, camera calls, sprite
updates, full table bytes, the adjacent eight-byte guard and exact queued messages.
Fixtures restore globals and blob bytes and release netlist/drawable allocations.

Before changing the callers, all 1,524 successful-slot cases passed. Exactly the
six exhausted-table cases failed: the guard became the attempted record and
MSG_NEW_ALIAS advertised 255. The normal packet-processing assertions still
passed. Example stream 1 guard: 01340900ffffffff; message: a5ff01340900ffffffff.
Stream 2 also reproduced frame+60 rollover at frame 0xffffffef. Evidence is in
build/port-network-alias/caller-before-fix.log (local, not committed).

The two checks now compare the unsigned byte to 0xFFu. Exhaustion skips only
the table write and alias announcement; packet processing continues. Successful
slots 128..254 stay accepted. The original helper ABI remains live for C callers,
while the existing Go reset wrapper calls the native reset directly. No reference
C implementation is retained. The outbound lookup remains unchanged.

Production C: **141,082 physical lines (−44)** in 153 files; reference C: **0**.
All accumulated protection/network/waypoint/rules/spell-class/ping tests pass on
386 default/server/highres. All three production binaries build. Fresh
network-alias-port gameplay passes both preserved screenshots with overrides
disabled. The full suite exactly matches the command-rule baseline: 15 passing,
3 known failing and 32 skipped/no-test packages; the same multiset of 1,553
failure entries, with zero added or removed. Local comparison metadata:
build/port-network-alias/full-suite-comparison.json.
