# Shared intrusive lists and player-group membership

Status: the C baseline is qualified. All five roots pass in default, repeat,
server and highres, with unchanged source and no skips (`build/port-lists/c-qualified`,
154.223s; initial independent checks 105.266s). Frozen expectations cover **2,319
records in five groups**, listed in `lists-captures.json`. No production
implementation has changed and no C prerequisite correction was needed.
Starting C: **68,897 lines / 88 files / zero reference C**, after qualified catalog
conversion **b6df7591**. Native integration is next.

The connected owner occupies 300 C lines in GAME1_1.c, from 00425760 through
00425BE0: list initialization, forward/reverse/index traversal, sorted insertion,
unlink, and player-group membership allocation/lookup/removal/free. The root
`lists.go` separately implements several operations; consolidate its raw-pointer
wrappers on the new Go owner too. Keep team assignment at 00418840 outside this
batch; it remains a real C user of the registry. This is a complete shared owner;
do not pad its scope with unrelated code to reach a line-count target.

Contracts cover complete normalized links/payloads and independent order models:
0/1/2/7/32/128 items, signed sort boundaries and duplicates, mixed insertion and
removal across two lists, detached removal/reinsertion, nullable traversal,
addressable index cases, and the existing partial append to an uninitialized
list. Normalize pointers only to fixture node identities. Registry checks cover
repeated initialization/free, duplicate IDs, signed member values, missing/empty
membership, Unicode names within the actual ten-word buffer, and allocation
ownership. Compare return-address identity without dereferencing freed records.

Preserve the three-word list layout and third-word sentinel/rank dual use. The
52-byte group record contains links, ten UTF16 words, ID, team pointer and a child
head; each member is sixteen bytes. Initialization guard remains set after free.
Removing a missing member also frees an already-empty group. Do not retain C
algorithms solely for comparison; the committed C baseline is the oracle.

Review finding: root Go traversal currently dereferences a nil successor where the
C owner returns nil. Consolidating on the C contract will correct this alongside
the conversion, with a native root-wrapper contract. Preserve raw UTF16 units;
do not round-trip names through Unicode decoding. Names beyond the group buffer
have no valid C oracle; the native destination copy must stay bounded.

The caller audit currently retains twelve list interfaces and four registry
interfaces for C users, retiring seven private registry interfaces. Recheck at
integration, including Go wrappers in maps, rules and player/server setup.
Use affected map/rule/player checks during implementation, then all three
production builds, exact known-suite comparison and both forced map-expansion
replays. Run the complete accumulated corpus in every configuration because the
list owner is widely shared. The tracked manifest includes separate milestone phases for default/server/highres.
Run those concurrently in isolated output directories with GOMAXPROCS=1, a 768 MiB
Go memory limit and a 1,800-second test timeout, as qualified in the codec milestone.
Record actual counts and any skips; only the opt-in map-population prerequisite
diagnostic is an expected skip.

Ignored drafts under `build/port-lists` are work in progress. The prepare-baseline
script has been applied; do not rerun it or overwrite installed fixtures with
stale drafts. The native draft is not installed or qualified.
