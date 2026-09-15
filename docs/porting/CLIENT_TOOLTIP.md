# Item-hover text and cursor tooltip storage

Status: original C baseline fully qualified after window conversion `98013201`.
Production is still C; the native draft is not yet applied or qualified.

Scope: all314 lines of client__gui__tooltip.c and the15-line cursor setter block
in GAME2_2.c: two connected routines /329 C lines. Both interfaces have live C
callers. Go callers should use native helpers when the conversion is applied.
No production correction is needed to establish this baseline.

## Owners and coverage

Reuse the actual drawable allocation/type registry and renderer owner, real
projectile/armor and modifier definitions, actual StringManager language/file
lookups, server spell and ability definitions, and mapped creature-title table.
Every fixture restores owned state. Metadata requests use the existing host-game
NetList path and capture queued bytes; this establishes UI message construction
and request suppression, not socket transport behavior.

Capture complete drawable words, normalized known modifier references, return
identity (borrowed pretty name versus shared scratch), raw UTF16 text and full
scratch/cursor tails, and queued messages. Only known owned pointers normalize.
The static fallback is verified empty; the fixture owns its two-byte cell without
assuming additional storage beyond it. The scratch and cursor capacities are
1,024 and256 UTF16 units respectively.

| Frozen group | Results | Coverage |
| --- | ---: | --- |
| tooltip-equipment | 3,840 | Eight languages, five class patterns, three subclass patterns, all16 modifier combinations, repeated calls |
| tooltip-books | 7,680 | Eight languages, all three subtype bits/precedence, unit-code boundaries/class marking, missing/known/pending metadata, repeated calls |
| tooltip-names | 240 | Nil/empty/borrowed/raw UTF16 names, missing equipment definitions, nil/empty/raw modifier descriptors |
| tooltip-cursor | 30 | Nil/empty/short/long text, lengths254–257 and up to1,024, raw surrogates, source bytes and unchanged tails |

Total:11,790 captured results /four groups. Independent contracts check exact
request bytes and pending sentinels137/41/6; borrowed pointers; language-specific
modifier spacing and book ordering; narrow-name localized formatting; cursor
truncation; and raw names up to the exact1,023-unit scratch limit across languages.
Counts describe executed inputs, including intentional repeated states.

## Baseline development

The first C run compiled and all four initial contracts passed. Captures were
intentionally rejected while expected hashes were unset. Reviewing text outputs
caught a fixture error: StringManager entries must include the ToolTip.c prefix;
unqualified entries had produced missing-string messages. Correct the fixture
before freezing, add explicit book-order and missing-definition formatting
contracts, and rerun. No production source or behavior changed.

The corrected c-b run passed all six then-present contracts and produced the four
frozen captures. Audited examples include “BookOf Spark”, “BookOf Flame Ω” and
“Missing: WhiteOrb”. Placeholder expectations were the only failures. The frozen
repeat additionally includes the exact-capacity raw-name contract. Artifacts and
ignored drafts are under build/port-client-tooltip; current source supersedes
historical fixture drafts. Go/C source stays immutable during checks.

## Translation review

Preserve raw UTF16 rather than decoding/re-encoding modifier and borrowed text.
Languages2 and6 contain unusual spaces that are observable compatibility behavior.
Projectile-versus-armor lookup and modifier suppression use the original class
and subclass masks. Book subtype precedence is spell, lore, then ability;
sentinel states suppress duplicate requests. The ordinary-name path retains a
borrowed pointer, while assembled names remain in the shared mapped buffer.

The existing variadic string formatter still owns the localized missing-equipment
format. A tiny production C ABI adapter may call it; no item-selection/assembly
algorithm should remain in C solely for testing. The shared formatter itself is
outside this batch. Cursor rendering is unchanged; these tests establish the
text/storage supplied to it, with headless gameplay at qualification boundaries.

The frozen c-c repeat passed all11 focused tests in20.467s, including the seventh
contract for exact-capacity names. All four capture files match c-b byte-for-byte.
The supplied asset's actual NoArmsInfo format also uses the narrow `%S` argument,
confirmed by reading nox.csf; decoded relevant entries are recorded in the ignored
asset-tooltip-strings.json. The fixture's explicit localized format tests that
same argument kind without depending on asset installation.

## C baseline qualification

affected: all135 selected root tests passed in58.836s. server: all134 selected root tests passed in175.086s. highres: all135 selected root tests passed in75.278s. Every selected test started and finished.

Client build passed in4.220s; fresh warrior gameplay
passed in37.233s with null audio and reference override=false.
All1494 Go/C/header source fingerprints remained unchanged.
No production behavior changed for this fixture-only baseline. The full asset
suite milestone remains the immediately preceding window conversion: exact1,553
known failure entries,15pass/3fail/32skip packages. Native qualification will repeat
that comparison. C remains88,212 /99 files /zero reference C.
