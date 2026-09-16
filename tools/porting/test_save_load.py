from pathlib import Path
import tempfile
import unittest
from save_load import verify_save_load


class SaveLoadTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.run = Path(self.tmp.name)
        for slot in ('AUTOSAVE', 'WORKING'):
            for relative in ('Player.plr', 'war01a/war01a.map'):
                p = self.run / 'data/save' / slot / relative
                p.parent.mkdir(parents=True, exist_ok=True)
                p.write_bytes(b'saved-' + relative.encode())
        loaded = str(self.run / 'data/save/WORKING/war01a/war01a.map')
        self.log = '\n'.join([
            '[console] system> Game Saved.',  # Campaign's initial automatic save.
            '[E2E] STATE: before_save',
            '[console] system> Game Saved.',
            '[E2E] STATE: after_save',
            '[E2E] STATE: moved_since_save',
            '[E2E] STATE: confirm saved game load',
            '[game] loading map "' + loaded + '"',
            '[map] client reading map: "' + loaded + '"',
            '[E2E] STATE: after_reload',
            '[E2E] STATE: resumed_after_reload'])

    def verify(self, text=None):
        (self.run / 'output.log').write_text(self.log if text is None else text)
        return verify_save_load(self.run)

    def test_complete_save_reload(self):
        result = self.verify()
        self.assertTrue(result['saved_map_loaded'])
        self.assertEqual(len(result['artifacts']), 2)

    def test_initial_autosave_is_not_explicit_save(self):
        text = self.log.replace('[E2E] STATE: before_save\n[console] system> Game Saved.',
                                '[E2E] STATE: before_save')
        with self.assertRaises(ValueError):
            self.verify(text)

    def test_original_map_load_is_not_saved_map_load(self):
        with self.assertRaises(ValueError):
            self.verify(self.log.replace('data/save/WORKING/war01a/war01a.map', 'data/maps/war01a/war01a.map'))

    def test_no_resumed_gameplay(self):
        with self.assertRaises(ValueError):
            self.verify(self.log.replace('[E2E] STATE: resumed_after_reload', ''))

    def test_changed_working_map(self):
        (self.run / 'data/save/WORKING/war01a/war01a.map').write_bytes(b'other')
        with self.assertRaises(ValueError):
            self.verify()

    def test_empty_save(self):
        for slot in ('AUTOSAVE', 'WORKING'):
            (self.run / 'data/save' / slot / 'Player.plr').write_bytes(b'')
        with self.assertRaises(ValueError):
            self.verify()


if __name__ == '__main__':
    unittest.main()
