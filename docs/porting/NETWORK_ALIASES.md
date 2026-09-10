# Network alias table: next-chunk investigation — 2026-09-10

Proposed scope: reset 57B920, slot selection 57B9A0 and record write 57BA10.
No production alias changes or original-C test fixtures have been installed.
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

## Confirmed exhaustion bug; user decision pending

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

Asked the user whether to fix exhaustion during the port (recommended) or defer.
The proposed fix compares against the byte sentinel 255 and skips the table
write and alias announcement when full, while continuing packet processing.
Do not change production alias behavior until that decision arrives.

Planned tests: original-C selection across every start byte, full-width and
negative keys, matching versus expired records, exact frame equality, wrap and
full-table failure; guarded C-owned storage for exact reset/write extents, key
truncation and return bits. If fixing exhaustion, add actual caller coverage for
full tables and successful aliases including 128..254, checking no write or
announcement on failure and unchanged normal packet processing.
