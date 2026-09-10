# Spell-class eligibility — 2026-09-10

Original-C baseline covers 38,016 ABI calls: all 256 byte classes plus eight
invalid raw 32-bit classes, all combinations of the three class flags with
unrelated bit noise, valid and invalid definitions, missing/nonpositive spell
IDs, and AllowAll both enabled and disabled. An independent arithmetic oracle
checks exact 0/9 results. Synthetic definitions use the real server Flags lookup.

Class 1 accepts any-class or conjurer flags; class 2 accepts any-class or wizard
flags. Other classes always return 9, including warrior and raw 257/258. Invalid
definitions still supply flags; absent/nonpositive definitions do not, even with
AllowAll. Preserve full-width C int arguments and the live C entry point.

All cases pass against the original C on 386 before production changes.
Production C before conversion: **141,215 physical lines**, 153 files, zero
reference C. Local artifacts: `build/port-spell-class/`.
