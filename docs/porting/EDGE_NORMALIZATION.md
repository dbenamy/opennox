# Edge normalization — 2026-09-11

Scope: 411490, the inverse edge-category mapping shared by Go mergeBorderEdge
and C subtile lookup/merge owners. Keep its ABI entry until the C callers move.
It reads width/height bytes from a valid physical EDGE row (64×60 at 85B3FC+28644),
without consulting active count, touching RNG, or writing any state.

The baseline covers every width/height byte pair and deduplicated values at each
ordered branch boundary, adjacent parity cases, signed extrema, zero and negatives.
Rows rotate across all 64 physical entries. Exact 3×3 returns raw input unchanged;
other dimensions preserve zero precedence, first width threshold, middle parity,
last-width threshold and default modular offset. C uses -fno-strict-overflow;
the independent expected default is computed using uint32 modular addition.

Reuse the edge fixture's per-call table/map/guard checks and exact restoration.
Both RNG streams must remain at their initial index after each normalization.
The shared Go RNG implementation remains untouched. No invalid physical row or
unbounded pointer is passed to C. Original-C qualification and Go port pending.

Original-C qualification passes **1,299,167 normalization calls**, plus the
1,145,856 existing edge-mapping calls after extending the shared fixture.
Production C baseline: **140,613 physical lines**, 153 files, zero reference C.
Artifacts: build/port-edge-normalization.
