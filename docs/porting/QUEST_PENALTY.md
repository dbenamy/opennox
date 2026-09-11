# Quest death-penalty policy

GAME5 54CBD0–54D080 is seven functions / 385 physical C lines. Registered
PlayerDie calls the root when quest lives run out. The complete policy halves
gold, removes gems and equipment, revokes spell/beast-scroll/warrior knowledge,
and sends the corresponding notifications. The six private helpers have no
other caller or registration. PlayerDie and its score handlers remain shared;
this policy's root requires its C ABI for that caller.

## Original-C baseline

The shared callback/lifecycle fixture supplies actual C-owned player/owner/update
memory, guarded inventory objects and initializer buffers, real gem types and
armor-bit lookup, controlled records in the real eligibility tables, real
shop pricing, gold, inventory unlink/dequip, delayed deletion and packet engines.
The existing item-class eligibility seam records each candidate and returns its
configured eligibility. Automatic shield equip is disabled in the fixture.

Every case captures full player, owner, update, item and initializer words;
normalized links; all six gem/price caches; packets and deletion order; both
PRNG indices; and memory guards. Player memory, deletion list, globals and
services are restored. Twenty cases install a real gold protection record,
capture encoded/decoded state and private RNG, and assert decoded gold equals
the player balance. Existing callback/initialization hashes must stay unchanged.

There are 21 policy smoke cases, four mixed-inventory smoke cases, 384 knowledge
contracts, 420 equipment contracts, 420 gem contracts, 1,024 generated cases,
and 20 gold-protection cases (2,293 total). Cases distinguish:

- Root order and the second class-zero armor removal after real unlink/dequip.
- Every protected spell ID, disabled eligibility, levels zero/one/multiple,
  class gates, 136/40 boundaries, and empty candidate RNG ranges.
- Equipped-item exceptions, all four modifier slots, replacement order,
  per-candidate eligibility calls and armor exclusion bits 1/4/0x400.
- Interleaved gem counts zero through five, first odd-gem credit, halving and
  unsigned gold wrap, real float32 shop-price conversion, and warm/cold caches.
- Destroyed items whose repeated delayed-delete calls do not unlink them.

Hashes are locked in src/quest_penalty_porttest_test.go:

- penalty-smoke: `935d5c69aae17f96ba002918c3ab251c162e3befb5b8993ea7b3b5f0e6cd97cb`.
- penalty-inventory-smoke: `284b6f32f7891acca86436e97b60e082df8832fa5fb775ee484d88afbba5149b`.
- penalty-knowledge: `43b81f95b13274ec0a001f340bea3fce8252fb4255214f1d0621a833653ed91f`.
- penalty-equipment: `7599536683bf26456102e04348984b21a5ab53839f7f1091cd899f57b72dc678`.
- penalty-gems: `9678247c376349f62e21c332ecde30986f6b3d1859104d69e373be78988e2c73`.
- penalty-corpus: `9cb2754bde579826b7486dedb7b222c6ce60bb172a8159db5d13818275942c15`.
- penalty-protection: `682c54b2ea10f9335abf9d55cde4669ed874cb5cac1140ebf66774773e9cfd15`.

C bodies remain intact: 134,954 production lines, 153 files, zero reference C.
Ignored logs and captures are under build/port-quest-penalty. Next: native
conversion of the complete family, followed by one qualification cycle.
