# MP3 Huffman decoding and dequantization

## Original-C baseline

Three native386 scalar-SSE2 processes and UBSan agree on8,223 power-helper cases,
13,224 patterned Huffman cases and12,906 directed codewords. The power domain is
all0..8206 caller values plus the source table's supported negative prefix -16..-1.
Patterned cases span all32 book indices, all start-bit alignments, zero/partial/full
big-value pairs, both count1 books, every version/rate/band layout, three scale
patterns, logical bit cutoffs and lookahead near the physical scratch end.

Luna's independent review identified that seeded patterns alone do not establish
escape-extreme/sign coverage. Before Go testing, primary added directed inputs
for every leaf in each book. Escaped leaves include zero/max extension bits and
all four sign combinations; ordinary leaves include both positive and negative
signs. Only input construction uses table traversal; expected values come from
the actual C helper. Every byte of the initial result prefix remains unchanged.
All26,130 Huffman records have true external guards and input/gr/scalefactor
immutability checks, explicitly verified before freezing.

[Provenance](mp3-huffman-c-capture.json) includes source/compiler/tool hashes and
prefix preservation. All2,417 numeric table entries independently match C:
145 float power values and2,272 integer Huffman entries. No copied C algorithms.
The primary owns algorithm translation and the runner; Luna supplied bounded table
extraction/power draft and independently reviewed the algorithm and fixture.

## Storage and arithmetic contracts

The decoder's logical bit limit differs from physical scratch capacity. Original
C reads ahead into the actual2,815-byte scratch, including retained bytes beyond
the current logical end. Go uses the backing slice capacity for this lookahead;
it must not zero-fill or truncate those bytes silently at integration. The tests
intentionally shorten the visible length while preserving physical capacity.

Preserve partial-band scale advancement, count1 zero coefficients left untouched,
cutoff checks before sign consumption, and final bit position set to the granule
limit. The caller clears spectral storage before use. Existing side-info validation
and physically sized scratch remain required. NaN/invalid-side-info hardening is
not silently substituted for original behavior in this private helper port.

The NMP3HUF1 fixture uses opcodes1(power),2(patterned),3(directed). All words are
little endian; each Huffman result records640 float words, bit position/limit and
four immutability/guard flags. Directed requests append a byte count and explicit
encoded prefix. The capture tool defines the full protocol.

Production still uses C. Standalone C is six lines/one file, reference zero; the
active decoder header and81 preamble bodies remain outside that source-size metric.
