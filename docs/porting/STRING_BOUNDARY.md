# Internal string boundaries

Status: original baseline qualified at `6ff98033`; conversion not installed yet.

Remove six local C imports by using Go pointer types at internal boundaries for
configuration, browser hosting, server options, map painting and object transfer.
Keep the live CString/GoString/CWString implementations and their ownership
unchanged. Tile/border C entrypoints retain their signatures. The historical
object-script adapter becomes a direct call with the same raw handle/output
values; its discarded return conversion has no observable effect.

Repository-wide searches found only definitions and internal header declarations
for the input text-buffer get/free pair and `nox_fs_normalize`, and no users of
`CStringArray`, the legacy `CBytes` wrapper or `CWStrSlice`. Retire these unused
interfaces rather than preserving their allocation bugs with new fixtures. No
supported external shared-library API or dynamic symbol lookup was found in the
build/source surface. This is an internal API removal, not a claim that unknown
external consumers are impossible. The direct test-fixture C.CBytes call remains.

## Baseline and coverage

The original source is byte-identical to the qualified allocation chunk, so its
fresh production evidence is reused for this baseline only. Newly run owner
contracts pass 177 roots per production profile, plus the separately audited
`TestConsoleCommandsNames` password owner: 178 distinct roots per profile. The
converted combined selection must match those actual names, not just the count.
Seven safe roots (CString ownership and border selection) pass without skips.

The selected historical object-record contract reaches Trigger versions below 3
with a read-only cryptfile, exercising `objectXferLegacyScript` twice. It checks
the frozen historical hash, stream position and object fields. The prefab runtime
selection also checks the underlying script scanner. Tile/border contracts cover
the retained C entrypoints; painting contracts cover their direct Go callers.

Two broader **original-source** safe selections fail in existing fixture setup:
`TestServerConfigScalarStorage` rejects mapped access at 0x715090, and
`TestTileSelectionABI` rejects its state observer's mapped access at 0x97cb68.
They pass in every normal production profile. Preserve these diagnostic logs;
do not claim the whole owner selection passes with safe runtime checks. Full
safe-build/static qualification remains required after conversion.

The selector audit found 129 existing owner roots missing from the accumulated
pattern. Add the exact newly qualified names, recorded in the baseline JSON.
Their absence from that pattern does not imply their earlier focused batch
qualifications never ran.

## Review and delegation

Luna drafted the bounded migration and retirement set without editing active
source. Primary reviewed every diff, confirmed the historical trigger branch,
and caught a direct border-lookup test bridge needing the same pointer migration.
Keep the root assertions and frozen expectations unchanged. Primary owns original
qualification, integration, source/ABI review and final acceptance. No measured
delegation savings are claimed.

Artifacts: `build/port-string-boundary/`; original qualification:
[string-boundary-c-qualification.json](string-boundary-c-qualification.json).
Draft/audit paths are ignored and are not acceptance evidence by themselves.
