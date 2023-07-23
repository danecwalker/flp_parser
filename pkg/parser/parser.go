package parser

import (
	"fmt"

	"github.com/danecwalker/flp-parser/pkg/defs"
)

// ReadDWord reads a dword from the given byte slice
func ReadDWord(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

// ReadWord reads a word from the given byte slice
func ReadWord(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

// ReadByte reads a byte from the given byte slice
func ReadByte(b []byte) uint8 {
	return uint8(b[0])
}

// Byte Iterator
type ByteIterator struct {
	b []byte
	i int
}

func (bi *ByteIterator) Next() uint8 {
	b := bi.b[bi.i]
	bi.i++
	return b
}

func (bi *ByteIterator) NextWord() uint16 {
	w := ReadWord(bi.b[bi.i:])
	bi.i += 2
	return w
}

func (bi *ByteIterator) NextDWord() uint32 {
	dw := ReadDWord(bi.b[bi.i:])
	bi.i += 4
	return dw
}

func dwordToString(dw uint32) string {
	return string([]byte{byte(dw), byte(dw >> 8), byte(dw >> 16), byte(dw >> 24)})
}

func (bi *ByteIterator) ReadEvent() defs.Event {
	kind := bi.Next()
	if kind <= 0x3F {
		c := bi.Next()
		return defs.NewBYTEEvent(kind, c)
	} else if kind >= 0x40 && kind <= 0x7F {
		c := bi.NextWord()
		return defs.NewWORDEvent(kind, c)
	} else if kind >= 0x80 && kind <= 0xBF {
		c := bi.NextDWord()
		return defs.NewDWORDEvent(kind, c)
	} else if kind >= 0xC0 && kind <= 0xFF {
		var vLeng [4]uint8
		var vLen int = 0
		i := 0
		for {
			// Get a byte
			b := bi.Next()
			// Add the first 7 bits of b to vLen
			vLeng[i] = b & 0x7F
			// If the 8th bit is 0, we're done
			if b&0x80 == 0 {
				// Combine the bytes into a single int
				for j := 0; j <= i; j++ {
					vLen |= int(vLeng[j]) << uint(7*j)
				}
				break
			}

			i++
		}

		ev := defs.NewTextEvent(kind, string(bi.b[bi.i:bi.i+int(vLen)]))
		bi.i += int(vLen)
		return ev
	} else {
		panic("invalid event kind")
	}
}

// checkHeader checks the header of the given byte slice
func parseHeader(b *ByteIterator) (string, error) {
	// check first 4 bytes for "FLhd"
	magic := dwordToString(b.NextDWord())
	if magic != "FLhd" {
		return "", fmt.Errorf("invalid file header")
	}

	lenNext := b.NextDWord()
	if lenNext != 6 {
		return "", fmt.Errorf("invalid file header")
	}

	b.NextWord()
	b.NextWord()
	b.NextWord()

	dataChunk := dwordToString(b.NextDWord())
	if dataChunk != "FLdt" {
		return "", fmt.Errorf("invalid file header")
	}

	b.NextDWord()

	versionEvent := b.ReadEvent()
	if versionEvent.Kind() != defs.EventKindFLP_Version {
		return "", fmt.Errorf("invalid file header")
	}

	version := versionEvent.Content().(string)

	return version, nil
}

// Parse parses the given byte slice into a slice of events
func Parse(b []byte) ([]defs.Event, error) {
	iter := ByteIterator{b: b, i: 0}
	version, err := parseHeader(&iter)

	if err != nil {
		return nil, err
	}

	fmt.Println(version)

	for iter.i < len(b) {
		event := iter.ReadEvent()
		if event.Kind() == defs.EventKindFilePath {
			fmt.Println(event)
		}
	}

	return nil, nil
}
