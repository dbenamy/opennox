# MP3 inverse transforms and overlap state

## Original-C baseline

The capture tool includes the actual production minimp3 header with scalar SSE2
386 arithmetic. Three independent runs and UBSan agree on 2,380 records: 98 DCT9,
50 IDCT3, 558 IMDCT12, 720 long-block, 240 short-block, 34 sign changes and 680
wrapper sequences containing 1,258 granule steps. No copied C algorithms are kept.
[Capture provenance](mp3-imdct-c-capture.json) records hashes and compiler flags.

Cases include signed zero, subnormals, dyadic inputs, finite random bit patterns,
basis impulses, empty and boundary band counts, three controlled windows and all
block types. Stateful sequences refresh spectral input while retaining overlap
through long/start/short/stop transitions. The fixture records exact float bits,
retained tails and external guards. Primitive DCT/IDCT captures observe a prefix
and a following tail, rather than the entire allocated workspace; full transforms
observe all 640 spectral and 320 overlap words. Luna independently reviewed source
order, capture bounds and protocol. Primary independently compared all 60 float
constants with C and reviewed state mutation and aliasing in the short transforms.

Magic is NMP3IMD1; byte opcodes1–7 precede little-endian requests and results.
For multi-step records all request pairs precede all expected outputs. The capture
tool defines the complete protocol. Expected results are frozen before Go testing.

The package is still private and unimported. Production continues using C; this
baseline does not retire any implementation. Standalone C remains six lines in one
file, reference C zero, with the active decoder header and 81 production preamble
bodies outside the standalone source count.
