# Numeric storage used by fixture accessors

Status: original-C baseline qualified; native draft not yet installed. Production
remains the qualified numeric-storage checkpoint `2a3eaec7`:166 C lines /6 files.
Only the tagged owner contract changes before the baseline is frozen.

## Scope

Move52 numeric owners (35 map generation,10 audio,7 sustained spells) into Go.
Their remaining C body uses are fixture address/get-set adapters. Replace those
adapters while keeping the same owners, index mappings, null slots and expected
values. Audio slot2 remains a C-owned pointer word; slots5/8 remain nil. The
production populationGlobal and sustainedGlobal helpers keep their mappings.
Three pointer-typed scalar declarations and pointer/array owners remain deferred.

One batch covers the shared ownership/accessor mechanism and avoids three separate
production cycles. The focused selection covers183 consumer roots, including map,
tile/border/floor, audio assets/streams/events, options/music/dialog integrations,
and sustained spells. Shipped audio assets are provided. The standalone population
diagnostic child is excluded; its four-case regression parent remains included.

## Baseline and review

Extend the existing owner contract from343 to395 owners:54,905 cases. New widths
and initializer bits come from the original C definitions; existing343 entries
remain unchanged. Check alias visibility, widths/signs, nonoverlap, neighboring
owners, GC and restoration. Require the original343 observations to match the
previous frozen capture before freezing the extended capture. No C algorithm is
introduced for these contracts. Production evidence remains reusable while only
the tagged fixture changes; native conversion requires fresh production gates.

Luna drafted the guarded migration and mapped consumers; primary caught a helper
name collision and an incomplete audio test selection. The corrected draft uses
unique fixture helper names and includes the additional direct consumers. C-only
helpers declare eight symbols as unsigned int, but existing Go-facing selectors
use uint32_t; the new contract follows the Go-facing declarations. Primary owns
baseline acceptance, exact accessor-mapping review and final integration.

The first probe passes all54,905 cases; the original343 observations are identical
to the prior frozen capture. Extended capture SHA-256:
`6de9fccecf723a0f1e0693560cc57abcf5eac137eebbaac634b771dfd7aa6b2b`.
Formal baseline contracts run twice under default and once under server/highres;
183 default consumers and static checks complete baseline acceptance. The recent
full three-profile consumer sweep is supporting evidence for unchanged production.

Baseline qualification passed: two independent default storage processes and
server/highres match the frozen capture;183 default consumer roots and all static
checks pass without skips. Source fingerprints agree; only the tagged fixture
changed. See [fixture-storage-c-qualification.json](fixture-storage-c-qualification.json).
