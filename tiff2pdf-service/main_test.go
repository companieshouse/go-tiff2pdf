package main

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/companieshouse/go-tiff2pdf/tiff2pdf"
)

func TestStringTo512Cchar(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{name: "empty"},
		{name: "ASCII", input: "document title"},
		{name: "511 bytes", input: strings.Repeat("a", 511)},
		{
			name:    "512 bytes",
			input:   strings.Repeat("a", 512),
			wantErr: "metadata exceeds maximum length",
		},
		{
			name:    "513 bytes",
			input:   strings.Repeat("a", 513),
			wantErr: "metadata exceeds maximum length",
		},
		{
			name:    "large input",
			input:   strings.Repeat("a", 4096),
			wantErr: "metadata exceeds maximum length",
		},
		{
			name:  "UTF-8 preserved",
			input: "caf\u00e9 \u4e16\u754c \U0001f600",
		},
		{
			name:  "UTF-8 at 511-byte boundary",
			input: strings.Repeat("a", 507) + "\U0001f600",
		},
		{
			name:    "512 UTF-8 bytes",
			input:   strings.Repeat("\u00e9", 256),
			wantErr: "metadata exceeds maximum length",
		},
		{
			name:    "wide rune starts at index 511",
			input:   strings.Repeat("a", 511) + "\U0001f600",
			wantErr: "metadata exceeds maximum length",
		},
		{
			name:    "embedded NUL",
			input:   "before\x00after",
			wantErr: "metadata contains a NUL byte",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tiff2pdf.StringTo512Cchar(tt.input)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("got error %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for i := 0; i < len(tt.input); i++ {
				if byte(got[i]) != tt.input[i] {
					t.Fatalf(
						"byte %d: got %#02x, want %#02x",
						i, byte(got[i]), tt.input[i],
					)
				}
			}
			for i := len(tt.input); i < len(got); i++ {
				if got[i] != 0 {
					t.Fatalf("byte %d: expected zero", i)
				}
			}
		})
	}
}

func minimalTIFF() []byte {
	const entries = 9
	const pixelOffset = 8 + 2 + entries*12 + 4

	data := make([]byte, pixelOffset+1)
	copy(data, "II")
	binary.LittleEndian.PutUint16(data[2:], 42)
	binary.LittleEndian.PutUint32(data[4:], 8)
	binary.LittleEndian.PutUint16(data[8:], entries)

	fields := []struct {
		tag   uint16
		typ   uint16
		value uint32
	}{
		{256, 4, 1},           // ImageWidth
		{257, 4, 1},           // ImageLength
		{258, 3, 8},           // BitsPerSample
		{259, 3, 1},           // Compression: none
		{262, 3, 1},           // PhotometricInterpretation: BlackIsZero
		{273, 4, pixelOffset}, // StripOffsets
		{277, 3, 1},           // SamplesPerPixel
		{278, 4, 1},           // RowsPerStrip
		{279, 4, 1},           // StripByteCounts
	}

	for i, field := range fields {
		offset := 10 + i*12
		binary.LittleEndian.PutUint16(data[offset:], field.tag)
		binary.LittleEndian.PutUint16(data[offset+2:], field.typ)
		binary.LittleEndian.PutUint32(data[offset+4:], 1)
		binary.LittleEndian.PutUint32(data[offset+8:], field.value)
	}

	return data
}

func TestConvertTiffToPDFMetadata(t *testing.T) {
	input := minimalTIFF()

	t.Run("accepts 511-byte metadata in every field", func(t *testing.T) {
		config := tiff2pdf.DefaultConfig()
		value := strings.Repeat("a", 511)
		config.Creator = value
		config.Author = value
		config.Subject = value
		config.Title = value

		out, err := tiff2pdf.ConvertTiffToPDF(input, config, "input.tif", "output.pdf")
		if err != nil {
			t.Fatalf("conversion failed: %v", err)
		}
		if out == nil || out.PageCount != 1 ||
			!strings.HasPrefix(string(out.PDF), "%PDF-") {
			t.Fatalf("expected a one-page PDF, got %+v", out)
		}
	})

	fields := []struct {
		name string
		set  func(*tiff2pdf.Config, string)
	}{
		{"Creator", func(c *tiff2pdf.Config, s string) { c.Creator = s }},
		{"Author", func(c *tiff2pdf.Config, s string) { c.Author = s }},
		{"Subject", func(c *tiff2pdf.Config, s string) { c.Subject = s }},
		{"Title", func(c *tiff2pdf.Config, s string) { c.Title = s }},
	}
	invalid := []struct {
		name    string
		value   string
		wantErr string
	}{
		{
			"too long",
			strings.Repeat("a", 512),
			"metadata exceeds maximum length",
		},
		{
			"NUL",
			"before\x00after",
			"metadata contains a NUL byte",
		},
	}

	for _, field := range fields {
		for _, value := range invalid {
			t.Run(field.name+"/"+value.name, func(t *testing.T) {
				config := tiff2pdf.DefaultConfig()
				field.set(config, value.value)

				out, err := tiff2pdf.ConvertTiffToPDF(
					input, config, "input.tif", "output.pdf",
				)
				if err == nil || !strings.Contains(err.Error(), value.wantErr) {
					t.Fatalf("got error %v, want %q", err, value.wantErr)
				}
				if out != nil {
					t.Fatal("expected nil output on validation failure")
				}
			})
		}
	}
}
