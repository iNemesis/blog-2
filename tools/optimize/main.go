package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := flag.String("dir", "static/images", "Dossier contenant les images à optimiser")
	contentDir := flag.String("content", "content", "Dossier des articles markdown pour la mise à jour des références")
	lossy := flag.Bool("lossy", false, "Active la compression visuelle (perte imperceptible, qualité 80)")
	quality := flag.Int("quality", 0, "Qualité de compression visuelle (1-100, 0 = mode lossless strict)")
	flag.Parse()

	q := *quality
	if *lossy && q == 0 {
		q = 80
	}

	if err := run(*dir, *contentDir, q); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		os.Exit(1)
	}
}

func run(dir, contentDir string, quality int) error {
	var totalBefore, totalAfter int64
	var count int

	if quality > 0 {
		fmt.Printf("Mode visuel activé (qualité %d/100)\n\n", quality)
	} else {
		fmt.Println("Mode lossless strict (0 altération de pixels)")
		fmt.Println()
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			return nil
		}

		before, after, newPath, err := optimizeFile(path, contentDir, quality)
		if err != nil {
			fmt.Printf("⚠️  %s: %v\n", path, err)
			return nil
		}

		count++
		totalBefore += before
		totalAfter += after

		saved := before - after
		if saved > 0 {
			pct := float64(saved) / float64(before) * 100
			if newPath != path {
				fmt.Printf("✓ %s → %s: %s → %s (-%.1f%%)\n", path, filepath.Base(newPath), formatBytes(before), formatBytes(after), pct)
			} else {
				fmt.Printf("✓ %s: %s → %s (-%.1f%%)\n", path, formatBytes(before), formatBytes(after), pct)
			}
		} else {
			fmt.Printf("· %s: déjà optimal (%s)\n", path, formatBytes(before))
		}
		return nil
	})

	if err != nil {
		return err
	}

	if count == 0 {
		fmt.Printf("Aucune image trouvée dans %s\n", dir)
		return nil
	}

	savedTotal := totalBefore - totalAfter
	pctTotal := 0.0
	if totalBefore > 0 {
		pctTotal = float64(savedTotal) / float64(totalBefore) * 100
	}
	fmt.Printf("\nRésumé : %d image(s) traitée(s)\n", count)
	fmt.Printf("Taille totale : %s → %s (gain : %s / -%.1f%%)\n",
		formatBytes(totalBefore), formatBytes(totalAfter), formatBytes(savedTotal), pctTotal)

	return nil
}

func optimizeFile(path, contentDir string, quality int) (before, after int64, newPath string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, path, err
	}
	before = int64(len(data))
	newPath = path

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".jpg", ".jpeg":
		var optimized []byte
		if quality > 0 {
			optimized, err = compressLossyJPEG(data, quality)
		} else {
			optimized, err = optimizeLosslessJPEG(data)
		}
		if err != nil {
			return 0, 0, path, err
		}
		after = int64(len(optimized))
		if after < before {
			if err := os.WriteFile(path, optimized, 0644); err != nil {
				return 0, 0, path, err
			}
			return before, after, path, nil
		}

	case ".png":
		// Mode visuel : si le PNG est opaque, on le convertit en JPEG haute efficacité
		if quality > 0 {
			img, err := png.Decode(bytes.NewReader(data))
			if err == nil && isOpaque(img) {
				var buf bytes.Buffer
				if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err == nil {
					jpegData := buf.Bytes()
					if int64(len(jpegData)) < before {
						newPath = strings.TrimSuffix(path, filepath.Ext(path)) + ".jpg"
						if err := os.WriteFile(newPath, jpegData, 0644); err != nil {
							return 0, 0, path, err
						}
						_ = os.Remove(path)

						// Met à jour les références dans les articles markdown
						oldRef := "/" + filepath.ToSlash(filepath.Clean(path))
						oldRef = strings.TrimPrefix(oldRef, "/static")
						newRef := strings.TrimSuffix(oldRef, filepath.Ext(oldRef)) + ".jpg"
						_ = updateMarkdownRefs(oldRef, newRef, contentDir)

						return before, int64(len(jpegData)), newPath, nil
					}
				}
			}
		}

		// Recompression PNG lossless par défaut
		optimized, err := optimizePNG(data)
		if err != nil {
			return 0, 0, path, err
		}
		after = int64(len(optimized))
		if after < before {
			if err := os.WriteFile(path, optimized, 0644); err != nil {
				return 0, 0, path, err
			}
			return before, after, path, nil
		}
	}

	return before, before, path, nil
}

// optimizePNG compresse les données PNG au niveau maximal sans perte de pixels.
func optimizePNG(data []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// compressLossyJPEG réencode un JPEG avec un facteur de qualité (1-100).
func compressLossyJPEG(data []byte, quality int) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// optimizeLosslessJPEG supprime les métadonnées superflues (Exif, IPTC, commentaires)
// sans jamais altérer les coefficients DCT (0 perte de qualité).
func optimizeLosslessJPEG(data []byte) ([]byte, error) {
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		return nil, errors.New("format JPEG invalide")
	}

	var out bytes.Buffer
	out.Write([]byte{0xFF, 0xD8}) // SOI

	i := 2
	for i < len(data) {
		if data[i] != 0xFF {
			return nil, errors.New("marqueur JPEG corrompu")
		}
		for i < len(data) && data[i] == 0xFF {
			i++
		}
		if i >= len(data) {
			break
		}
		marker := data[i]
		i++

		if marker == 0xD9 { // EOI
			out.Write([]byte{0xFF, 0xD9})
			break
		}
		if marker >= 0xD0 && marker <= 0xD7 { // RST0..RST7
			out.Write([]byte{0xFF, marker})
			continue
		}

		if i+2 > len(data) {
			return nil, errors.New("longueur de segment tronquée")
		}
		length := int(binary.BigEndian.Uint16(data[i : i+2]))
		if length < 2 || i+length > len(data) {
			return nil, errors.New("segment JPEG tronqué")
		}

		segData := data[i : i+length]
		i += length

		// Marqueurs ignorés : 0xE1 (Exif/XMP), 0xED (IPTC), 0xFE (Commentaires)
		if marker == 0xE1 || marker == 0xED || marker == 0xFE {
			continue
		}

		out.WriteByte(0xFF)
		out.WriteByte(marker)
		out.Write(segData)

		if marker == 0xDA { // Start Of Scan
			streamStart := i
			for i < len(data)-1 {
				if data[i] == 0xFF {
					next := data[i+1]
					if next == 0xD9 {
						out.Write(data[streamStart : i+2])
						return out.Bytes(), nil
					}
					if next != 0x00 && !(next >= 0xD0 && next <= 0xD7) && next != 0xFF {
						break
					}
				}
				i++
			}
			out.Write(data[streamStart:])
			return out.Bytes(), nil
		}
	}

	return out.Bytes(), nil
}

func isOpaque(img image.Image) bool {
	type opaquer interface {
		Opaque() bool
	}
	if o, ok := img.(opaquer); ok {
		return o.Opaque()
	}
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a < 0xffff {
				return false
			}
		}
	}
	return true
}

func updateMarkdownRefs(oldRef, newRef, contentDir string) error {
	return filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if bytes.Contains(data, []byte(oldRef)) {
			replaced := bytes.ReplaceAll(data, []byte(oldRef), []byte(newRef))
			if err := os.WriteFile(path, replaced, 0644); err != nil {
				return err
			}
			fmt.Printf("  ↳ référence mise à jour dans %s (%s → %s)\n", path, oldRef, newRef)
		}
		return nil
	})
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
