# MP3 synthesis DCT

## Original-C baseline

The actual scalar-SSE2 C header produces identical results in three native386
processes and UBSan for 2,172 records. Cases cover all576 basis coordinates at
both channel offsets, finite random patterns, signed zero, subnormals, empty and
boundary column counts, and supplemental12-column transforms. The active
MINIMP3_ONLY_MP3 build calls18 columns for every MPEG version;12 is the inactive
Layer I/II caller extent, not a lower-rate Layer III path.

Each result records all1,216 float words, including the other channel and retained
tail, plus external guards. Primary compared all24 table constants with original
C and reviewed scalar expression association and in-place staging. The capture
imports only common data-pattern/I/O helpers from the earlier capture tool, with
both tool hashes in [provenance](mp3-dct-c-capture.json). No C algorithm copies.
The NMP3DCT1 fixture uses byte opcode1, five little-endian request words and4,868
result bytes. Expectations are frozen before testing the Go implementation.

The package remains unimported. Production still uses C; standalone C is six
lines/one file, reference C zero, plus the active header and81 preamble bodies.
