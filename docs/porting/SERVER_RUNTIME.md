# Remaining runtime helpers

Status: **qualified original-C baseline**; native conversion is next. Production
remains identical to `39d24e36`. See [qualification](server-runtime-c-qualification.json).
C remaining: **2,705 physical lines in 23 production files**, zero reference C.

## Scope

The batch connects the remaining numeric conversions, observer count, meter-wave
initialization, modifier/equipment lookups, rejected-player list disposal,
pause-effect lifecycle and saved-creature ownership transfer. Move their Go
callers with the implementations and remove obsolete C interfaces. No C algorithm
will remain solely for testing.

Defer `GameEx.c` player-name conversion and the server-listing formatter to a
following text/layout batch. Their existing local-buffer and reserved-byte
behavior needs a separate review before freezing expectations.

## Original-C contracts

The baseline invokes actual C through the existing application owners and thin
porttest adapters. It checks independent invariants as well as recording results.
Default, repeated default, server and highres each pass all 13 roots without
skips. All eleven capture hashes /55,986 cases agree, and static checks pass.
Production source matches the preceding qualified conversion; only tagged tests
and documentation changed, so its production evidence is reused.

- Double conversion: 37,059 cases, including raw exponent/mantissa combinations,
  signed boundaries, NaNs/infinities and all 12 x87 precision/rounding modes.
  Existing float conversion contracts cover about 1.4 million inputs.
- Observers: 6,144 cases covering active player patterns, every status byte,
  the enable flag and host-slot exclusion. Player indices remain valid throughout.
- Meter wave: 18 complete 320-entry snapshots with shipped coefficients, varied
  scales and zero-angle controls; output guards and source constants are checked.
- Rejected-player disposal: 21 empty/single/multiple-node list cases, up to 1,024
  nodes. The existing allocation observer checks traversal and exactly-once frees.
- Equipment lookup/loading: 960 lookup and 72 loading cases; first-match order,
  ASCII/non-ASCII queries, UTF-16 labels, missing versus empty strings, all mask
  bits, sentinels, initial ready states and repeated calls.
- Material cache: 288 cases, including missing-then-available definitions,
  preloaded caches, class gates and null modifier slots. Cache lookup precedes
  the class check; a missing material and a null slot can compare equal.
- Armor conductivity: 480 class/type/definition cases with float boundary values.
- Modifier images/labels: 9,216 cases covering every byte, duplicated and high-bit
  table keys, image failures and cached/repeated lookup. The C argument is signed:
  values 128–255 never match an unsigned table key. Image loads occur in row order.
- Pause lifecycle: 960 passing probe cases, including early-return gates,
  supplied/cached actors, missing factories, effect kinds, animation outcomes,
  clock narrowing and repeated starts/stops. Effects use the real object factory;
  animation calls are observed at the existing retained-callee boundary.
- Host ownership: 768 passing probe cases, using real server object and owned
  lists, host records and minimap reporting. Check child traversal while the owner
  mutation reverses the destination list, class/status gates and marker teardown.

## Reachability decisions

- The old bit encoder and its four private bit helpers have no caller.
- The C modifier-next helper and next-map-group getter are unreferenced.
- The old C material-draw helper has no caller. The renderer already uses its
  separate Go implementation; similar spelling does not connect the two.
- The old memory-accounting list has no production population path and its output
  table has no consumer. An older console fixture fabricates the list. Remove the
  orphaned accounting implementation and that artificial fixture while retaining
  the console command contract.
- Four old session words have only reset-to-zero writers. Their reset and
  conditional game-loop callbacks are inert. The separate pause-effect state is
  live and covered explicitly.
- The old audio-path callback clears its own buffer and always returns null.
  Preserve its public empty-string result without translating unreachable suffix
  assembly.
- Nine no-op modifier callbacks are still registered; two are explicitly compared
  by pointer identity. They are not dead merely because their bodies are empty.
  Keep those identities until the callback registrations are migrated.

The inventory and detailed symbol searches are in
`build/port-final-server-helpers`; the tracked qualification manifest records the final source identity, counts and
capture hashes. All jobs joined before committing this baseline.

## Probe and storage notes

Probes 1–3 pass their successive numeric, lookup and image contracts. Probe4 found
that fixture action 1600 selects the spell dispatcher before the controls owner;
use unused controls actions 1498/1499 for the two new lifecycle contracts. The
pause cases also explicitly enable cooperative mode, as its real pause service
requires. Probe5 passed ownership but found that the book-busy flag uses a separate C global,
not the mapped address. The fixture now owns the actual flag and the unpause-hold
global. Probe6 passes all eleven roots /55,986 captured cases. Production source is
unchanged by these fixture corrections.

Cleanup removed 1,112,747,701 bytes of verified-identical completed-run asset
copies from `client-resources-native` and `client-resources-native-save`, preserving
original assets, changed run files and all evidence. The cleanup apply is consumed;
restore with `python3 build/port-final-server-helpers/deduplicate-resource-assets.py
--restore RUN` when rerunning either completed scenario.

Identical completed captures were consolidated with verified relative symlinks,
recovering another 1,663,983,267 bytes while preserving their paths and hashes.
The local audit is `build/port-final-server-helpers/deduplicated-captures.json`.
