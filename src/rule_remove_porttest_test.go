//go:build porttest

package opennox

import (
	"reflect"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestRulesRemoveABI(t *testing.T) {
	tests := []struct {
		name          string
		mapName, file string
		dirs          []string
		files         map[string]string
		removed       bool
		want          []legacy.PortTestRuleRemoveEntry
	}{
		{
			name: "missing", mapName: "Arena", file: "gone.rul",
			dirs:    []string{"maps/Arena"},
			files:   map[string]string{"maps/Arena/keep.rul": "keep", "maps/Other/keep.rul": "other"},
			removed: false,
			want: []legacy.PortTestRuleRemoveEntry{
				{Path: "maps", Dir: true}, {Path: "maps/Arena", Dir: true}, {Path: "maps/Arena/keep.rul", Data: "keep"},
				{Path: "maps/Other", Dir: true}, {Path: "maps/Other/keep.rul", Data: "other"},
			},
		},
		{
			name: "existing leaves directory", mapName: "Arena", file: "user.rul",
			dirs: []string{"maps/Arena"}, files: map[string]string{"maps/Arena/user.rul": "old"}, removed: true,
			want: []legacy.PortTestRuleRemoveEntry{{Path: "maps", Dir: true}, {Path: "maps/Arena", Dir: true}},
		},
		{
			name: "case insensitive backslashes", mapName: "mIxEd\\nEsT", file: "tArGeT.rUl",
			files: map[string]string{"maps/MiXeD/NeSt/Target.RUL": "remove", "maps/MiXeD/NeSt/Sibling.RUL": "stay"}, removed: true,
			want: []legacy.PortTestRuleRemoveEntry{
				{Path: "maps", Dir: true}, {Path: "maps/MiXeD", Dir: true}, {Path: "maps/MiXeD/NeSt", Dir: true}, {Path: "maps/MiXeD/NeSt/Sibling.RUL", Data: "stay"},
			},
		},
		{
			name: "nested file path", mapName: "Arena", file: "rules\\custom\\user.rul",
			files: map[string]string{"maps/Arena/rules/custom/user.rul": "remove", "maps/Arena/rules/keep.rul": "keep", "maps/Arena/other.rul": "other"}, removed: true,
			want: []legacy.PortTestRuleRemoveEntry{
				{Path: "maps", Dir: true}, {Path: "maps/Arena", Dir: true}, {Path: "maps/Arena/other.rul", Data: "other"}, {Path: "maps/Arena/rules", Dir: true},
				{Path: "maps/Arena/rules/custom", Dir: true}, {Path: "maps/Arena/rules/keep.rul", Data: "keep"},
			},
		},
		{
			name: "remove empty directory", mapName: "Arena", file: "empty",
			dirs: []string{"maps/Arena/empty"}, removed: true,
			want: []legacy.PortTestRuleRemoveEntry{{Path: "maps", Dir: true}, {Path: "maps/Arena", Dir: true}},
		},
		{
			name: "nonempty directory fails", mapName: "Arena", file: "full",
			files: map[string]string{"maps/Arena/full/keep": "keep"},
			want:  []legacy.PortTestRuleRemoveEntry{{Path: "maps", Dir: true}, {Path: "maps/Arena", Dir: true}, {Path: "maps/Arena/full", Dir: true}, {Path: "maps/Arena/full/keep", Data: "keep"}},
		},
		{
			name: "empty filename selects directory", mapName: "Arena", file: "",
			dirs: []string{"maps/Arena"}, removed: true,
			want: []legacy.PortTestRuleRemoveEntry{{Path: "maps", Dir: true}},
		},
		{
			name: "embedded NUL truncates arguments", mapName: "Arena\x00ignored", file: "user.rul\x00ignored",
			files: map[string]string{"maps/Arena/user.rul": "remove", "maps/Arena/keep.rul": "stay"}, removed: true,
			want: []legacy.PortTestRuleRemoveEntry{{Path: "maps", Dir: true}, {Path: "maps/Arena", Dir: true}, {Path: "maps/Arena/keep.rul", Data: "stay"}},
		},
		{
			name: "empty map name", mapName: "", file: "user.rul",
			files: map[string]string{"maps/user.rul": "remove", "maps/keep.rul": "stay"}, removed: true,
			want: []legacy.PortTestRuleRemoveEntry{{Path: "maps", Dir: true}, {Path: "maps/keep.rul", Data: "stay"}},
		},
		{
			name: "bounded long filename", mapName: "ArenaABC", file: strings.Repeat("x", 238),
			files: map[string]string{"maps/ArenaABC/" + strings.Repeat("x", 238): "remove"}, removed: true,
			want: []legacy.PortTestRuleRemoveEntry{{Path: "maps", Dir: true}, {Path: "maps/ArenaABC", Dir: true}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := legacy.PortTestRuleRemove(legacy.PortTestRuleRemoveSpec{Dir: t.TempDir(), Map: tc.mapName, File: tc.file, Dirs: tc.dirs, Files: tc.files})
			if err != nil {
				t.Fatal(err)
			}
			wantResult := 0
			if tc.removed {
				wantResult = 1
			}
			if got.Result != wantResult || !reflect.DeepEqual(got.After, tc.want) {
				t.Fatalf("result=%v after=%+v want result=%v after=%+v", got.Result, got.After, wantResult, tc.want)
			}
			if tc.removed == false && !reflect.DeepEqual(got.Before, got.After) {
				t.Fatalf("missing removal changed tree: before=%+v after=%+v", got.Before, got.After)
			}
		})
	}
}
