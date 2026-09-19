# Spell, ability and guide awards — native conversion

All 26 selected C bodies and the private signed enchantment counter are now Go.
Seven C interfaces remain for four direct caller roots and three registered item-use
callbacks; nineteen private interfaces and the C counter are retired. Four C
translation units were removed. Go callers invoke native helpers directly.

Production C: **18,662 physical lines in 61 files, zero reference C**, down **820**.
The reduction includes 629 selected body lines, the counter, and obsolete comments,
address markers and blank lines in touched C files. Header removal is not counted.
Original C baseline: **27554f1f**; qualified parent: **a8015048**.

## Qualification

- All **108 affected roots** pass without skips on default, server and highres.
  All **81 captures / 12,248 records** match C exactly across the three targets.
  The 15 new frozen captures contain 6,763 records. No frozen expectations changed.
- All **2,639 source fingerprints** agree across tests and production qualification.
  Static memory checks pass. The native selection adds existing player-stats and
  book-tooltip contracts for callers moved to Go.
- Three fresh production binaries qualify on ELF32/386/SSE2/CGO, including retained
  and retired symbols and absence of test helpers. Full suite matches the exact
  known 1,553 failure entries and package outcomes: 15 pass, 3 fail, 32 skip.
- Fresh headless gameplay and explicit save/load pass against existing references.
  These scenarios complement the focused contracts; they do not exercise every
  class, award or family branch. Physical display/audio remain manual checks.

Evidence: [native qualification](book-awards-native-qualification.json),
[C qualification](book-awards-c-qualification.json),
[selection](book-awards-selection.json), [literal callers](book-awards-callers.json),
[interface plan](book-awards-interface-plan.json),
[native manifest](book-awards-native-batch.json).
Local runs: `build/port-book-awards/native-final-{default,server,highres,production}`
and `static-native-final.log`. All drafts/installers are consumed.

## Contracts and review decisions

Coverage includes catalog order/bounds/row widths, loader partial failure, localized
strings and image lookup, enchantment count/order, class/ID/level gates, all 17
quest single-level spell IDs, signed overrides, uint32 wrap, family propagation,
bookkeeping, exact report bytes and recipients, both reliable and direct message
queues, notification conditions, deferred audio, item consumption and shop closure.
Fixtures use real player, object, catalog, queue, audio and shop owners. Only external
image-resource loading is observed through its existing boundary.

Two legacy behaviors are preserved for separate review:

- Spell-family propagation applies its quest cap and bookkeeping to the original
  spell ID, rather than consistently using the family member. Independent contracts
  explicitly preserve this behavior.
- An unknown field-guide item is consumed when class/known-guide admission passes,
  even though awarding guide ID zero fails.

The native enchantment iterator checks a nonpositive signed count before subtracting
one, preserving the minimum-int gate. The existing spell-class predicate was
extracted without logic changes; its ABI contract remains selected.

The C fixture review corrected colon-qualified localization IDs, added observations
of the separate direct-message queue, and fixed a fixture observer that assumed
player class in deliberate non-player tests. Owner cleanup now remains valid after
fatal assertions. All were corrected before freezing. Production C was unchanged;
its baseline reused the parent's production evidence after source-identity and
binary-hash checks. The native conversion always received fresh production checks.
