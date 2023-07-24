package parser

import (
	"bytes"
	"fmt"
	"os"

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

		ev := defs.NewTextEvent(kind, uint8(i), string(bi.b[bi.i:bi.i+int(vLen)]))
		bi.i += int(vLen)
		return ev
	} else {
		panic("invalid event kind")
	}
}

// checkHeader checks the header of the given byte slice
func parseHeader(b *ByteIterator) (*defs.Header, error) {
	var H = defs.Header{}

	// check first 4 bytes for "FLhd"
	m := b.NextDWord()
	magic := dwordToString(m)
	if magic != defs.F_SIG {
		return nil, fmt.Errorf("invalid file header")
	}
	H.FSig = m

	lenNext := b.NextDWord()
	if lenNext != 6 {
		return nil, fmt.Errorf("invalid file header")
	}
	H.ChunkSize = lenNext

	H.Format = b.NextWord()
	H.NChannels = b.NextWord()
	H.BeatDivPerQNote = b.NextWord()

	d := b.NextDWord()
	dataChunk := dwordToString(d)
	if dataChunk != defs.F_DAT {
		return nil, fmt.Errorf("invalid file header")
	}
	H.FDat = d

	H.FDatChunkSize = b.NextDWord()

	return &H, nil
}

// Parse parses the given byte slice into a slice of events
func Parse(b []byte) (*defs.Project, error) {
	iter := ByteIterator{b: b, i: 0}
	proj := defs.Project{}
	header, err := parseHeader(&iter)
	proj.Header = header

	if err != nil {
		return nil, err
	}

	for iter.i < len(b) {
		event := iter.ReadEvent()
		proj.Events = append(proj.Events, event)
	}

	fmt.Printf("%+v\n", proj.Header)
	fmt.Printf("%+v\n", len(proj.Events))

	return &proj, nil
}

func WriteByte(b uint8, buf *bytes.Buffer) {
	buf.WriteByte(b)
}

func WriteWord(w uint16, buf *bytes.Buffer) {
	buf.WriteByte(byte(w))
	buf.WriteByte(byte(w >> 8))
}

func WriteDWord(dw uint32, buf *bytes.Buffer) {
	buf.WriteByte(byte(dw))
	buf.WriteByte(byte(dw >> 8))
	buf.WriteByte(byte(dw >> 16))
	buf.WriteByte(byte(dw >> 24))
}

func WriteEvent(ev defs.Event, buf *bytes.Buffer) {
	kind := ev.Kind()
	WriteByte(kind, buf)

	if kind <= 0x3F {
		WriteByte(ev.Value().(uint8), buf)
	} else if kind >= 0x40 && kind <= 0x7F {
		WriteWord(ev.Value().(uint16), buf)
	} else if kind >= 0x80 && kind <= 0xBF {
		WriteDWord(ev.Value().(uint32), buf)
	} else if kind >= 0xC0 && kind <= 0xFF {
		v := ev.Value().(string)
		vLen := uint32(len(v))
		for {
			b := byte(vLen & 0x7F)
			vLen >>= 7
			if vLen == 0 {
				WriteByte(b, buf)
				break
			} else {
				WriteByte(b|0x80, buf)
			}
		}

		buf.WriteString(v)
	} else {
		panic("invalid event kind")
	}
}

func Write(project *defs.Project, path string) error {

	// Write header
	buf := bytes.Buffer{}
	WriteDWord(project.Header.FSig, &buf)
	WriteDWord(project.Header.ChunkSize, &buf)
	WriteWord(project.Header.Format, &buf)
	WriteWord(project.Header.NChannels, &buf)
	WriteWord(project.Header.BeatDivPerQNote, &buf)
	WriteDWord(project.Header.FDat, &buf)
	WriteDWord(project.Header.FDatChunkSize, &buf)

	// Write events
	for _, ev := range project.Events {
		WriteEvent(ev, &buf)
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsExist(err) {
		fmt.Println("File already exist")
		return err
	}

	// Write to file
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file")
		return err
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	fmt.Println("Successfully wrote to file")

	return nil
}
