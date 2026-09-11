# Network alias table: next-chunk investigation — 2026-09-10

Proposed scope: reset 57B920, slot selection 57B9A0 and record write 57BA10.
Original-C helper fixtures are installed and pass on 386; no production alias
changes yet.
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

Planned tests: original-C selection across every start byte, full-width and
negative keys, matching versus expired records, exact frame equality, wrap and
full-table failure; guarded C-owned storage for exact reset/write extents, key
truncation and return bits. If fixing exhaustion, add actual caller coverage for
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
