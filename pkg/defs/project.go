package defs

const F_SIG = "FLhd"
const F_DAT = "FLdt"

// Header is the header of a project file
type Header struct {
	FSig            uint32
	ChunkSize       uint32
	Format          uint16
	NChannels       uint16
	BeatDivPerQNote uint16
	FDat            uint32
	FDatChunkSize   uint32
}

type Project struct {
	Header          *Header
	Events          []Event
	FLStudioFactory string
}
