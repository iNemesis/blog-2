package main

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestOptimizePNG(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 5), G: uint8(y * 5), B: 100, A: 255})
		}
	}

	var raw bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.NoCompression}
	if err := enc.Encode(&raw, img); err != nil {
		t.Fatalf("encodage initial: %v", err)
	}

	optimized, err := optimizePNG(raw.Bytes())
	if err != nil {
		t.Fatalf("optimizePNG: %v", err)
	}

	if len(optimized) >= raw.Len() {
		t.Errorf("attendu compression < %d, obtenu %d", raw.Len(), len(optimized))
	}

	decoded, err := png.Decode(bytes.NewReader(optimized))
	if err != nil {
		t.Fatalf("décodage image optimisée: %v", err)
	}
	if decoded.Bounds() != img.Bounds() {
		t.Errorf("bounds différentes: attendu %v, obtenu %v", img.Bounds(), decoded.Bounds())
	}
}

func TestOptimizeLosslessJPEG(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		t.Fatalf("jpeg encode: %v", err)
	}

	app1 := []byte{0xFF, 0xE1, 0x00, 0x08, 'E', 'x', 'i', 'f', 0x00, 0x00}
	withMeta := append([]byte{0xFF, 0xD8}, append(app1, rawJPEG(buf.Bytes())...)...)

	optimized, err := optimizeLosslessJPEG(withMeta)
	if err != nil {
		t.Fatalf("optimizeLosslessJPEG: %v", err)
	}

	if len(optimized) >= len(withMeta) {
		t.Errorf("attendu taille < %d, obtenu %d", len(withMeta), len(optimized))
	}

	_, err = jpeg.Decode(bytes.NewReader(optimized))
	if err != nil {
		t.Fatalf("décodage jpeg optimisé: %v", err)
	}
}

func rawJPEG(b []byte) []byte {
	return b[2:]
}

func TestCompressLossyJPEG(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 4), G: uint8(y * 4), B: 200, A: 255})
		}
	}

	var highQuality bytes.Buffer
	if err := jpeg.Encode(&highQuality, img, &jpeg.Options{Quality: 98}); err != nil {
		t.Fatal(err)
	}

	lossy, err := compressLossyJPEG(highQuality.Bytes(), 75)
	if err != nil {
		t.Fatalf("compressLossyJPEG: %v", err)
	}

	if len(lossy) >= highQuality.Len() {
		t.Errorf("attendu gain lossy: %d vs %d", len(lossy), highQuality.Len())
	}
}

func TestIsOpaque(t *testing.T) {
	opaque := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			opaque.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	if !isOpaque(opaque) {
		t.Error("l'image devrait être détectée comme opaque")
	}

	transparent := image.NewRGBA(image.Rect(0, 0, 10, 10))
	transparent.Set(5, 5, color.RGBA{R: 255, G: 0, B: 0, A: 128})
	if isOpaque(transparent) {
		t.Error("l'image avec alpha=128 ne devrait pas être opaque")
	}
}

func TestUpdateMarkdownRefs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-md-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	mdFile := filepath.Join(tmpDir, "post.md")
	content := []byte("image: /images/test.png\n\n![alt](/images/test.png)")
	if err := os.WriteFile(mdFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	if err := updateMarkdownRefs("/images/test.png", "/images/test.jpg", tmpDir); err != nil {
		t.Fatalf("updateMarkdownRefs: %v", err)
	}

	updated, err := os.ReadFile(mdFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "image: /images/test.jpg\n\n![alt](/images/test.jpg)"
	if string(updated) != expected {
		t.Errorf("attendu %q, obtenu %q", expected, string(updated))
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		in   int64
		want string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
	}

	for _, tt := range tests {
		got := formatBytes(tt.in)
		if got != tt.want {
			t.Errorf("formatBytes(%d) = %s, attendu %s", tt.in, got, tt.want)
		}
	}
}
