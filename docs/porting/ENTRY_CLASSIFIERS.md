# Entry character classifiers

Status: direct-libc conversion qualified against committed actual-C baseline
`cbf62bc1`. Standalone C remains 45 lines in four files, zero reference C.

The former forwarding bodies `entryDigit` and `entryAlnum` accepted unsigned
16-bit text units and called libc `iswdigit` and `iswalnum`. The UI only observes zero or
nonzero. The predicates now call those same libc functions directly from Go,
using explicit `C.wint_t` widening. Locale behavior is retained; this is not an
ASCII-only or Unicode-category rewrite and does not remove libc/CGO.

The baseline extracts two small Go predicates around the actual production C
wrappers. The caller still gives each the same uint16 value and rejects false.
A tagged test visits every value 0..65535 for both predicates and packs 131,072
booleans into digit/alphanumeric bitsets. Three fresh-process captures must agree
before freezing their hash. Independent ASCII acceptance/rejection assertions
make a trivial all-zero result fail. No duplicate C reference wrappers are added.

All eleven existing entry-widget roots remain selected, preserving their frozen
captures. A new real-input contract checks both-filter precedence, the event-state
gate, digit acceptance, restoring alphabetic acceptance when only the digit filter
is removed, and rejecting spaces. The keyboard matrix also covers IME/language,
modifiers, event states and character-boundary inputs; composition, ownership,
rendering and lifecycle tests supply related integration checks.

The baseline and direct-libc phases each require the four legacy contracts
(classifiers, numeric storage, raw storage, audio GC), twelve widget roots across
default/server/highres, static checks, safe build, production/ABI, exact known-suite
comparison and fresh headless gameplay/save-load. Frozen default-locale captures
do not establish other locales' character sets; retaining the same libc calls
preserves that dependency. Safe runtime remains outside this build gate.

Luna drafted a read-only plan and mechanically derived the manifests. Primary
review changed its proposed duplicate-reference fixture to capture the actual
production wrappers, reviewed the manifest diffs and wrote the exhaustive and
filter-precedence contracts. Builds and acceptance remain primary-owned.

## Frozen actual-C capture

Three fresh processes agree on all 131,072 classifications; capture SHA256 is
`f5c39db5e866885891bbb1fe8d0b2244d9fbcea7d232cdd4b7e6f90dd03ff6e0`.
The 16,384 bytes contain an 8,192-byte digit bitset followed by an alphanumeric
bitset, ascending code-unit index, low bit first. The new widget precedence
contract passes. Both baseline and direct-libc qualifications pass.

## Baseline qualification

All four legacy contracts and twelve widget roots pass in default/server/highres,
with no skips and unchanged storage/UI captures. Static checks, safe build,
production/ABI, exact known-suite comparison, headless gameplay and explicit
save/load pass. Source fingerprints agree, preflight matches production, and all
ten callback identities remain distinct. See [C qualification](entry-classifiers-c-qualification.json).
Safe runtime was not tested; known full-suite failures remain unchanged.

## Direct-libc qualification

The predicates now call `C.iswdigit(C.wint_t(v))` and
`C.iswalnum(C.wint_t(v))` directly. Both custom C forwarding bodies are removed.
The complete 131,072-case capture is unchanged in default/server/highres, alongside
all four contracts and twelve widget roots with no skips. Fresh static, safe build,
production/ABI, exact known-suite, gameplay and explicit save/load gates pass;
source fingerprints agree and all ten callback addresses remain distinct.
See [native qualification](entry-classifiers-native-qualification.json).

Standalone C remains **45 lines in four files** (zero reference C): this removes
two cgo-preamble C bodies, which that physical `.c` metric never counted. libc,
its locale semantics, generated cgo wrappers and other inline/third-party C remain.
No production correction or golden update was needed after conversion.
