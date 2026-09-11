# Subtile point predicate and list lookup — 2026-09-11

Scope: 4113A0 and 411350. The latter returns the value of the last matching
node in a C-owned linked list, or the caller's fallback for a null list or no
match. Nodes are five 32-bit words: value, auxiliary, border row, edge, next.
It normalizes each edge with 411490 before applying the 12-category point test.
Neither function writes any input or table state.

Original-C baseline passes 39,200 predicate calls and 2,180 list lookups.
Predicate coverage includes coordinates -1..47, all categories and nearby
invalid values, plus signed extreme coordinates/categories. The independent
oracle uses two diagonal half-plane comparisons and explicit uint32 wrapping
for 46-y. Equality belongs to both adjoining regions.

List checks cover all 64 physical border rows, normalized categories, ordered
mixed edges, every last-match position in lengths 0..31, nonmatching tails,
raw signed fallback values, and null lists with null/non-null points. Node
values are distinct from row/index values. Point and node storage is C-owned
with guards; every word/link and table byte is checked after each call. Blob
contents and guards are restored and verified. Tail links are explicitly null.
No cyclic/invalid lists, invalid physical rows or non-null-list/null-point calls
are passed to C. No copied C algorithm or assets are retained.

Production C baseline: **140,578 physical lines**, 153 files, zero reference C.
Artifacts: build/port-subtile-lookup. Original-C qualification passed; port and
accumulated validation remain. Predicate's only C caller is lookup, so retire
its C bridge after porting both; keep lookup's ABI for grid owner 411160.
