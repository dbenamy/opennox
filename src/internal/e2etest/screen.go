package e2etest

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"
)

// CheckScreen compares decoded pixels with an existing golden. Updating goldens
// requires update=true; normal checks never create or replace expected images.
// Mismatches retain the actual image and, for equal dimensions, a visual diff.
func CheckScreen(path string, got image.Image, update bool) error {
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, got); err != nil {
		return err
	}
	if update {
		return os.WriteFile(path, encoded.Bytes(), 0644)
	}
	stem := strings.TrimSuffix(path, ".png")
	for _, suffix := range []string{"_got.png", "_diff.png"} {
		if err := os.Remove(stem + suffix); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read screen golden %q (use NOX_E2E_OVERRIDE=true to create it): %w", path, err)
	}
	want, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode screen golden %q: %w", path, err)
	}
	a, b := rgba(got), rgba(want)
	if a.Rect == b.Rect && bytes.Equal(a.Pix, b.Pix) {
		return nil
	}
	if err := os.WriteFile(stem+"_got.png", encoded.Bytes(), 0644); err != nil {
		return err
	}
	if a.Rect != b.Rect {
		return fmt.Errorf("screen %q size mismatch: got %v, want %v", path, a.Rect.Size(), b.Rect.Size())
	}
	// Highlight RGB differences while preserving differences in alpha as well.
	diff := image.NewNRGBA(a.Rect)
	for i, x := range a.Pix {
		d := int(x) - int(b.Pix[i])
		if d < 0 {
			d = -d
		}
		d *= 10
		if d > 255 {
			d = 255
		}
		if i%4 == 3 {
			d = 255 - d
		}
		diff.Pix[i] = byte(d)
	}
	encoded.Reset()
	if err := png.Encode(&encoded, diff); err != nil {
		return err
	}
	if err := os.WriteFile(stem+"_diff.png", encoded.Bytes(), 0644); err != nil {
		return err
	}
	return fmt.Errorf("screen %q pixel mismatch", path)
}

func rgba(src image.Image) *image.NRGBA {
	dst := image.NewNRGBA(image.Rectangle{Max: src.Bounds().Size()})
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	return dst
}
