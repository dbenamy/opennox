"""Verify the warrior save/load scenario exercised both save and reload paths."""
import hashlib
from pathlib import Path


def verify_save_load(run):
    run = Path(run).resolve()
    log = (run / 'output.log').read_text()

    def after(marker, start):
        pos = log.find(marker, start)
        if pos < 0:
            raise ValueError('save/load evidence missing or out of order: ' + marker)
        return pos + len(marker)

    before = after('[E2E] STATE: before_save', 0)
    saved = after('[console] system> Game Saved.', before)
    after_save = after('[E2E] STATE: after_save', before)
    if saved > after_save:
        raise ValueError('explicit save did not finish before after_save')
    moved = after('[E2E] STATE: moved_since_save', after_save)
    confirmed = after('[E2E] STATE: confirm saved game load', moved)
    saved_map = run / 'data/save/WORKING/war01a/war01a.map'
    server_read = after('[game] loading map "' + str(saved_map) + '"', confirmed)
    client_read = after('[map] client reading map: "' + str(saved_map) + '"', server_read)
    restored = after('[E2E] STATE: after_reload', client_read)
    after('[E2E] STATE: resumed_after_reload', restored)

    artifacts = {}
    for relative in ('Player.plr', 'war01a/war01a.map'):
        saved = (run / 'data/save/AUTOSAVE' / relative).read_bytes()
        working = (run / 'data/save/WORKING' / relative).read_bytes()
        if not saved or saved != working:
            raise ValueError('saved and reloaded files differ or are empty: ' + relative)
        artifacts[relative] = dict(bytes=len(saved), sha256=hashlib.sha256(saved).hexdigest())
    return dict(explicit_save=True, saved_map_loaded=True, resumed_after_reload=True,
                artifacts=artifacts)


if __name__ == '__main__':
    import json
    import sys
    print(json.dumps(verify_save_load(sys.argv[1]), indent=2))
