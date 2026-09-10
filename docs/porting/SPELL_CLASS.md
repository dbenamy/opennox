# Spell-class eligibility — 2026-09-10

Original-C baseline covers 38,016 ABI calls: all 256 byte classes plus eight
invalid raw 32-bit classes, all combinations of the three class flags with
unrelated bit noise, valid and invalid definitions, missing/nonpositive spell
IDs, and AllowAll both enabled and disabled. An independent arithmetic oracle
checks exact 0/9 results. Synthetic definitions use the real server Flags lookup.

Class 1 accepts any-class or wizard flags; class 2 accepts any-class or conjurer
flags. Other classes always return 9, including warrior and raw 257/258. Invalid
definitions still supply flags; absent/nonpositive definitions do not, even with
AllowAll. Preserve full-width C int arguments and the live C entry point.

All cases pass against the original C on 386 before production changes.
Production C before conversion: **141,215 physical lines**, 153 files, zero
reference C. Local artifacts: `build/port-spell-class/`.

## Go conversion

Original-C baseline: `2970e5e9`. The live 57AEA0 C ABI now calls Go and uses
server.Spells.Flags directly, preserving lookup before class rejection and full
C int width. The class masks use the existing named spell constants.

Removed the unused C chat predicate 57A160 and its declaration after confirming
there are no callers or address references. Its live root Go equivalent already
serves network_client.go. The unrelated C list cleanup 57ADF0 remains unchanged.

Production C: **141,180 physical lines (−35)** in 153 files; reference C: **0**.
Accumulated protection/network/waypoint/rules/spell-class tests pass on 386
with default, server and highres tags. All three production binaries build.
The fresh spell-class-port headless scenario passes both preserved screenshot
checks with overrides disabled (exit 0). The full suite was not repeated for
this small leaf immediately after the command-rule milestone; that milestone
still records the unchanged 1,553 known failure entries.
