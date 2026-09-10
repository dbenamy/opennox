package blobs

import (
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatGoOffsetPreservesTerms(t *testing.T) {
	cases := []string{
		"uintptr(x)*13+71276+uintptr(y)",
		"100+x+y", "x+100+y", "x+(100+y)", "(100+x)+y",
		"100-x+y", "100+x-y", "100+x+(y*4)", "100+(x-y)+y",
		"0+x+y", "100+x+20+y", "100+x-20+y", "100+x+(20-y)",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			input := "memmap.Uint8(0x587000, " + expr + ")"
			parse := func(call string) Access {
				list, err := BlobAccessesGo(token.NewFileSet(), "fixture.go", []byte("package p; var v = "+call))
				require.NoError(t, err)
				require.Len(t, list, 1)
				return list[0]
			}
			got := parse(input).String()
			_, err := parser.ParseExpr(got)
			require.NoError(t, err)
			require.Equal(t, got, parse(got).String(), "formatting must be idempotent")
			// Evaluate both expressions with several concrete inputs. This checks the
			// arithmetic independently of the formatter's static/dynamic decomposition.
			offset := strings.TrimSuffix(strings.SplitN(got, ",", 2)[1], ")")
			for _, values := range [][2]string{{"0", "0"}, {"3", "7"}, {"11", "5"}} {
				subst := strings.NewReplacer("x", values[0], "y", values[1])
				original, err := types.Eval(token.NewFileSet(), nil, token.NoPos, subst.Replace(expr))
				require.NoError(t, err)
				formatted, err := types.Eval(token.NewFileSet(), nil, token.NoPos, subst.Replace(offset))
				require.NoError(t, err)
				require.Equal(t, original.Value.ExactString(), formatted.Value.ExactString())
			}
		})
	}
}
