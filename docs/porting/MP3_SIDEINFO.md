# Layer III side-information parser

This second Go decoder preparation translates `L3_read_side_info`. Production
continues using the complete C decoder until the Go path is assembled and qualified.
Standalone C remains six lines/one file, reference C zero; the third-party header
and 81 production preamble bodies remain outside that standalone count.

## Original-C baseline

`tools/porting/capture_mp3_sideinfo.py` calls the actual production header, with no
copied decoder algorithm. It generates structured field mutations, directed
truncations and deterministic arbitrary input. Three separate native 386 processes
agree on all 36,864 records (23,306 successful returns and 13,558 error returns).
A fourth UBSan build also matches with no diagnostics. Source/compiler/fixture
hashes are in [mp3-sideinfo-c-capture.json](mp3-sideinfo-c-capture.json).

Coverage includes all three supported MPEG versions and rates, four channel modes,
one/two/four output slots, four start offsets, all block types, mixed/short/long
selection, big-values 288/289 boundary, MPEG-2 compress 499/500 preflag boundary,
reservoir offsets, part23 budget checks and lengths 0–40 plus larger payloads.
This is directed coverage, not an exhaustive Cartesian product. Random input is
fixed-seed; no shipped asset data is present in this fixture.

Each record starts with a four-byte header, logical byte length, start bit,
initial fill byte (the last three uint32 little-endian), and 64 input bytes. The
rest of a 4,160-byte backing buffer is zero. Output is signed return, final bit
position and limit, then all four granules: 21 scalar fields in declaration order,
a flag indicating that the original table pointer was retained, and 40 normalized
table bytes. Selected static tables are copied through their first zero and padded;
untouched table pointers reference a 40-byte sentinel. Struct padding and raw
addresses are excluded. The magic is `NMP3SID1`; the capture shim defines exact order.

Nonzero initial fields expose writes that successful paths intentionally omit and
partial writes before errors. The physical input backing buffer covers all pointer
construction even for truncated logical input. No negative bit positions, arithmetic
overflow or invalid header table indices enter the original parser.

Primary caught a representation defect in Luna's first draft: recording the table
family alone discarded the selected sample-rate row, and could reinterpret retained
state after errors using a later header. The draft is being corrected to borrow the
specific static table slice, matching the original pointer's role. No faulty draft
has been installed in production or accepted as qualified.
