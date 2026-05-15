// Command stego encodes a flag into a PNG image using LSB steganography
// and optionally appends a false flag after the PNG's IEND chunk.
//
// Usage:
//
//	go run . -in photo.png -out stego.png -flag 'eros{real_flag}' -decoy 'eros{fake_hint}'
//	go run . -decode -in stego.png
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	var (
		inPath  = flag.String("in", "", "input PNG path")
		outPath = flag.String("out", "", "output PNG path (encode mode)")
		flagStr = flag.String("flag", "", "flag to encode via LSB")
		decoy   = flag.String("decoy", "", "false flag to append after IEND")
		decode  = flag.Bool("decode", false, "decode mode: extract LSB flag from input")
	)
	flag.Parse()

	if *inPath == "" {
		fatal("missing -in")
	}

	if *decode {
		runDecode(*inPath)
		return
	}

	if *outPath == "" {
		fatal("missing -out")
	}
	if *flagStr == "" {
		fatal("missing -flag")
	}

	runEncode(*inPath, *outPath, *flagStr, *decoy)
}

func runEncode(inPath, outPath, flagStr, decoy string) {
	img := readPNG(inPath)

	// Encode the flag into the LSB of the image pixels.
	encoded := encodeLSB(img, []byte(flagStr))

	// Write the PNG to a temporary buffer so we can append the decoy.
	f, err := os.Create(outPath)
	if err != nil {
		fatal("creating output: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, encoded); err != nil {
		fatal("encoding PNG: %v", err)
	}

	// Append the decoy false flag after the IEND chunk.
	if decoy != "" {
		if _, err := f.Write([]byte(decoy)); err != nil {
			fatal("appending decoy: %v", err)
		}
	}

	fmt.Printf("encoded %q into %s (%d pixels used of %d available)\n",
		flagStr, outPath, bitsNeeded(len(flagStr)), img.Bounds().Dx()*img.Bounds().Dy()*3)

	if decoy != "" {
		fmt.Printf("appended decoy %q after IEND\n", decoy)
	}
}

func runDecode(inPath string) {
	img := readPNG(inPath)
	msg := decodeLSB(img)
	fmt.Printf("LSB decoded: %s\n", msg)

	// Also check for appended data after IEND.
	appended := readAppended(inPath)
	if len(appended) > 0 {
		fmt.Printf("appended data: %s\n", appended)
	}
}

// encodeLSB encodes data into the least significant bits of an image's RGB channels.
// Format: 32-bit big-endian length prefix, then the message bytes.
// Each bit of the payload is stored in the LSB of one colour channel value.
func encodeLSB(img image.Image, data []byte) *image.NRGBA {
	bounds := img.Bounds()
	out := image.NewNRGBA(bounds)

	// Copy all pixels first.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			out.Set(x, y, img.At(x, y))
		}
	}

	// Build the payload: 4-byte length + message.
	payload := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(payload[:4], uint32(len(data)))
	copy(payload[4:], data)

	capacity := bounds.Dx() * bounds.Dy() * 3 // 3 channels per pixel (RGB)
	needed := len(payload) * 8
	if needed > capacity {
		fatal("image too small: need %d bits, have %d (%d pixels)",
			needed, capacity, bounds.Dx()*bounds.Dy())
	}

	bitIdx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y && bitIdx < needed; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && bitIdx < needed; x++ {
			r, g, b, a := out.At(x, y).RGBA()
			// RGBA() returns 16-bit values; scale back to 8-bit.
			channels := [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}

			for c := range 3 {
				if bitIdx >= needed {
					break
				}
				byteIdx := bitIdx / 8
				bitPos := 7 - (bitIdx % 8) // MSB first
				bit := (payload[byteIdx] >> bitPos) & 1

				// Clear LSB and set it to our bit.
				channels[c] = (channels[c] & 0xFE) | bit
				bitIdx++
			}

			out.SetNRGBA(x, y, color.NRGBA{
				R: channels[0],
				G: channels[1],
				B: channels[2],
				A: uint8(a >> 8),
			})
		}
	}

	return out
}

// decodeLSB extracts LSB-encoded data from an image.
func decodeLSB(img image.Image) string {
	bounds := img.Bounds()

	// First, read 32 bits to get the length.
	length := extractBits(img, 0, 32)
	msgLen := binary.BigEndian.Uint32(length)

	if msgLen == 0 || int(msgLen) > bounds.Dx()*bounds.Dy() {
		fatal("no LSB data found (decoded length: %d)", msgLen)
	}

	// Now read the message.
	data := extractBits(img, 32, int(msgLen)*8)
	return string(data)
}

// extractBits reads n bits starting from startBit from the LSBs of the image's RGB channels.
func extractBits(img image.Image, startBit, n int) []byte {
	bounds := img.Bounds()
	result := make([]byte, (n+7)/8)

	bitIdx := 0
	channelIdx := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			channels := [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}

			for c := range 3 {
				if channelIdx >= startBit && bitIdx < n {
					bit := channels[c] & 1
					byteIdx := bitIdx / 8
					bitPos := 7 - (bitIdx % 8)
					result[byteIdx] |= bit << bitPos
					bitIdx++
				}
				channelIdx++
				if bitIdx >= n {
					return result
				}
			}
		}
	}
	return result
}

// readAppended reads any data appended after the PNG IEND chunk.
func readAppended(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		fatal("reading file: %v", err)
	}

	// PNG IEND chunk is: 4-byte length (0x00000000) + "IEND" + 4-byte CRC.
	// The IEND marker sequence in the file is the bytes for the chunk.
	iend := []byte{0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}
	for i := range len(data) - len(iend) + 1 {
		if bytesEqual(data[i:i+len(iend)], iend) {
			after := data[i+len(iend):]
			if len(after) > 0 {
				return after
			}
			return nil
		}
	}
	return nil
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func bitsNeeded(msgLen int) int {
	return (4 + msgLen) * 8 // 4-byte length prefix + message
}

func readPNG(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		fatal("opening %s: %v", path, err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		fatal("decoding PNG: %v", err)
	}
	return img
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "stego: "+format+"\n", args...)
	os.Exit(1)
}
