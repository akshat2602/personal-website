package internal

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// GenerateOGImage creates a per-post OG image by overlaying title text onto the base image.
// It saves the result to outputDir/og/<slug>.png and returns the relative path.
func GenerateOGImage(post *Post, baseImagePath, outputDir string) error {
	ogDir := filepath.Join(outputDir, "og")
	if err := os.MkdirAll(ogDir, 0755); err != nil {
		return err
	}

	outPath := filepath.Join(ogDir, post.Slug+".png")

	// Load base image
	f, err := os.Open(baseImagePath)
	if err != nil {
		// If no base image, skip silently
		return nil
	}
	defer f.Close()

	baseImg, _, err := image.Decode(f)
	if err != nil {
		return nil
	}

	// Create a new RGBA image from the base
	bounds := baseImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, baseImg, bounds.Min, draw.Src)

	// Draw a semi-transparent dark overlay for text contrast
	overlay := image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 160})
	overlayRect := image.Rect(
		bounds.Min.X, bounds.Max.Y-120,
		bounds.Max.X, bounds.Max.Y,
	)
	draw.DrawMask(img, overlayRect, overlay, image.Point{}, nil, image.Point{}, draw.Over)

	// Draw title text
	title := post.Title
	if len(title) > 60 {
		title = title[:57] + "..."
	}

	drawWrappedText(img, title, bounds.Max.X, bounds.Max.Y)

	// Write output
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, img)
}

// drawWrappedText draws title text onto the image using the basic font.
// For a production site, you'd want to use a proper TTF with golang.org/x/image/font/opentype.
func drawWrappedText(img *image.RGBA, title string, imgWidth, imgHeight int) {
	face := basicfont.Face7x13

	textColor := color.RGBA{R: 248, G: 248, B: 242, A: 255} // Dracula foreground

	margin := 30
	lineHeight := 20
	maxWidth := imgWidth - margin*2

	words := strings.FieldsFunc(title, func(r rune) bool { return unicode.IsSpace(r) })
	var lines []string
	var current string
	for _, w := range words {
		candidate := current
		if candidate != "" {
			candidate += " "
		}
		candidate += w
		if measureText(face, candidate) > maxWidth && current != "" {
			lines = append(lines, current)
			current = w
		} else {
			current = candidate
		}
	}
	if current != "" {
		lines = append(lines, current)
	}

	// Start Y from near the bottom
	startY := imgHeight - 90 + (len(lines)-1)*lineHeight/2

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
	}

	for i, line := range lines {
		y := startY + i*lineHeight
		d.Dot = fixed.Point26_6{
			X: fixed.I(margin),
			Y: fixed.I(y),
		}
		d.DrawString(line)
	}
}

func measureText(face font.Face, s string) int {
	d := &font.Drawer{Face: face}
	return int(d.MeasureString(s) >> 6)
}
