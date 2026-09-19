# Server-browser conversion plan

Selected scope:65 original C bodies /2,067 body lines. Keep UI lifecycle, linked
server records, ordering, map markers, proximity popup, details and host-description
construction in one batch. The marker-image helper sub_437320 has only one private
caller and is included rather than leaving an unnecessary C interface.

## Reachability and boundaries

The selected sub_439CC0 has no reachable caller: its only invocation is inside the
constant-false game-host display block in noxworld.c. Remove that nine-line body
and declaration with the disabled block; do not translate or preserve it for tests.
There are64 live selected routines /2,058 body lines.

Initial outside-body audit is retained under build/port-server-browser/
outside-selected.json. It distinguishes actual outside C calls from Go wrappers,
headers and forward declarations. Retain adapters for sub_43AF30 (GAME5_2),
sub_43B6E0/sub_43B750 (existing client notifications), and the constructor used by
the main-menu animation slot. Animation callbacks sub_438330, sub_438370 and sub_43B490, still stored in the
shared C-layout animation record, need thin adapters. Seven selected adapters
remain in total;58 selected interfaces can retire. The global-owner audit identifies
34 C owners to move and the two shared flags to retain. Go wrappers should call native code
directly. Per-frame update registration can use the existing Go SetUpdateFunc2;
there is no need for a C-only callback just to call Go once per frame.

Use a typed Go browser owner for window references, selection, sort/order and
connection UI state. Preserve existing Go accessors while changing their backing
storage. Audit all remaining literal C symbol references after removal. Shared
connection flags528252/528256 still have external C users and remain shared.
Mapped server settings, quest state and resource tables retain their existing owners.

The12 deferred configuration callbacks remain a separate reachability audit; their
registration table alone does not establish a live consumer. They are not included
in the65-body baseline scope. Retain the live coordinate getter4A7EF0.

## Behavioral constraints

- Preserve signed widths, uint32-before-double radius arithmetic, four-region
  search/fallback, unsigned coordinate shifts and popup clamping.
- Preserve all ten ordering modes, new-before-old ties, high-port lookup behavior,
  endpoint parser forms and name truncation. Do not replace libc address parsing
  with net.ParseIP without matching captured lexical cases.
- Own the selected record across re-sort and clearing. Original C copies nodes on
  re-sort and abandons old allocations; the fixture validates selection before
  explicitly releasing those original allocations. Native ownership should prevent
  leaks without invalidating observable selection data.
- C records are169 bytes; Go's current view is172. Do not over-read C allocations.
- Check event kinds before interpreting numeric construction IDs as pointers.
- Host descriptions currently advertise0x000f039a even on default/server Go builds:
  legacy/video_highres.go has an unconditional C compiler define. Preserve this
  existing output and document the C/Go version-constant disagreement for review.

## Qualification

Freeze only independent, repeated original-C captures. Run the directly affected
browser, class-selection, server-config/options/panels, GUI listbox/entry/window,
quest-runtime and resource-definition contracts on default/server/highres. Reuse
the verified parent production source identity for the C baseline. The isolated
headless browser scenario has already passed capture and exact comparison.

After conversion, require fresh production/ABI builds, exact known full-suite
comparison, static mapped-memory checks, isolated browser screen comparison and
qualified gameplay/save-load. Update C LOC, review notes and checkpoint, commit and
push, then continue. Retire the previously deferred character initializer C adapter
and empty selclass header with this native source batch.
