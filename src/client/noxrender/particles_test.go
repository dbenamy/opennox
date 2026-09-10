package noxrender

import (
	"image"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Each raw-pixel reference was verified against the original PNG golden using
// the historical color expansion (maximum 5-bit channel -> 248, now 255).
// Hashing framebuffer words avoids depending on that presentation conversion.
var particleCases = []struct {
	name string
	opt  particleOpt
	exp  string
}{
	{
		name: "white",
		opt:  particleOpt{rad: 10, blur: 0, intens: 0xff, color: RGB{0xff, 0xff, 0xff}},
		exp:  "caf0fc32de5d0bde3e66c9a50e72e20b93d0ead6da35dced8adc971b22cbca11",
	},
	{
		name: "green",
		opt:  particleOpt{rad: 10, blur: 0, intens: 0xff, color: RGB{140, 220, 80}},
		exp:  "008838a9f38d9ab29138eab6160f9af8508051dec7eb5fc72c78f865b96039ab",
	},
	{
		name: "white3",
		opt:  particleOpt{rad: 10, blur: 3, intens: 0xff, color: RGB{0xff, 0xff, 0xff}},
		exp:  "90c857d2064a533ae876627558ae810699e5cc9c17d8981a1837cd01e070134a",
	},
	{
		name: "green3",
		opt:  particleOpt{rad: 10, blur: 3, intens: 0xff, color: RGB{140, 220, 80}},
		exp:  "80e175e68b0386ece3466f49f7213869715161ef7dee257e7ddd12104596d5cb",
	},
	{
		name: "white32",
		opt:  particleOpt{rad: 10, blur: 32, intens: 0x80, color: RGB{0xff, 0xff, 0xff}},
		exp:  "38c023c57b5f38c158c6c55fc0ce4aecbcca0c22c9f06bfa738094c092c282fe",
	},
	{
		name: "green32",
		opt:  particleOpt{rad: 10, blur: 32, intens: 0x80, color: RGB{140, 220, 80}},
		exp:  "8a337a925f21eef4c0757818fa0a903dc0310ec0258b34f79bbcc92090f61484",
	},
}

func TestDrawParticle(t *testing.T) {
	debug := os.Getenv("NOX_RENDER_DEBUG") == "true"
	const outDir = ".testOut"
	if debug {
		err := os.MkdirAll(outDir, 0755)
		require.NoError(t, err)
	}
	for _, c := range particleCases {
		t.Run(c.name, func(t *testing.T) {
			img := genParticle(c.opt)
			off, sz, _ := img.Meta()
			csz := image.Point{
				X: (sz.X + off.X) * 2,
				Y: (sz.Y + off.Y) * 2,
			}
			pos := image.Pt(sz.X/2, sz.Y/2)

			pix := newBlack16(csz.X, csz.Y)
			d := newRenderData(csz.X, csz.Y)

			r := NewRender(slog.Default(), nil)
			r.SetPixBuffer(pix)
			r.SetData(d)

			r.DrawImage16(img, pos)

			if debug {
				fname := filepath.Join(outDir, "part_"+c.name+".png")
				out, err := os.Create(fname)
				require.NoError(t, err)
				err = png.Encode(out, pix)
				closeErr := out.Close()
				require.NoError(t, err)
				require.NoError(t, closeErr)
			}
			require.Equal(t, c.exp, pixelHash16(pix))
		})
	}
}
