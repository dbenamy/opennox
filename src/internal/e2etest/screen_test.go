package e2etest

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckScreen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "screen.png")
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.SetNRGBA(0, 0, color.NRGBA{R: 120, G: 30, B: 10, A: 255})
	// A missing reference must fail, not silently bless the current output.
	require.ErrorIs(t, CheckScreen(path, img, false), os.ErrNotExist)
	_, err := os.Stat(path)
	require.True(t, os.IsNotExist(err))
	require.NoError(t, CheckScreen(path, img, true))
	original, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NoError(t, CheckScreen(path, img, false))
	// A changed pixel must fail and produce diagnostics without mutating input
	// or replacing the expected image.
	img.SetNRGBA(0, 0, color.NRGBA{G: 255, A: 255})
	before := append([]byte(nil), img.Pix...)
	require.ErrorContains(t, CheckScreen(path, img, false), "pixel mismatch")
	require.Equal(t, before, img.Pix)
	require.FileExists(t, filepath.Join(filepath.Dir(path), "screen_got.png"))
	require.FileExists(t, filepath.Join(filepath.Dir(path), "screen_diff.png"))
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, original, data)
	// Different dimensions must fail safely, including when the golden is smaller.
	require.ErrorContains(t, CheckScreen(path, image.NewNRGBA(image.Rect(0, 0, 3, 3)), false), "size mismatch")
	require.NoFileExists(t, filepath.Join(filepath.Dir(path), "screen_diff.png"))
	require.NoError(t, CheckScreen(path, img, true))
	require.NoError(t, CheckScreen(path, img, false))
	require.NoFileExists(t, filepath.Join(filepath.Dir(path), "screen_got.png"))
}

func TestCheckScreenFormats(t *testing.T) {
	path := filepath.Join(t.TempDir(), "screen.png")
	pal := image.NewPaletted(image.Rect(0, 0, 2, 2), color.Palette{color.Black, color.White})
	pal.SetColorIndex(1, 1, 1)
	require.NoError(t, CheckScreen(path, pal, true))
	// Compare decoded pixels, allowing different image implementations and origins.
	img := image.NewGray(image.Rect(5, 5, 7, 7))
	img.SetGray(6, 6, color.Gray{Y: 255})
	require.NoError(t, CheckScreen(path, img, false))
}

func TestCheckScreenFileErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "screen.png")
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	require.NoError(t, os.WriteFile(path, []byte("not a PNG"), 0644))
	require.ErrorContains(t, CheckScreen(path, img, false), "decode screen golden")
	// An output directory cannot be overwritten by a golden update.
	require.Error(t, CheckScreen(filepath.Dir(path), img, true))
}
