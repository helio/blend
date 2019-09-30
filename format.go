package blend

// FileHeader is at the start of each blender file and gives general decoding information.
type FileHeader struct {
	// File identifier, always BLENDER
	Identifier [7]byte
	// Size of a pointer, '-' indicates 64 bits, '_' indicates 32 bits
	PointerSize byte
	// TypeIdx of byte ordering used, 'v' means little endian, 'V' big endian
	Endianness byte
	// Version of Blender file the file was created in, 280 means version 2.80
	Version [3]byte
}

// FileBlock64 represents a file-block if the file is encoded with 64 bits.
type FileBlock64 struct {
	header *FileBlockHeader64
	data   []byte
}

// FileBlock32 represents a file-block if the file is encoded with 32 bits.
type FileBlock32 struct {
	header *FileBlockHeader32
	data   []byte
}

// FileBlockHeader64 represents a file-block header if the file is encoded with 64 bits.
type FileBlockHeader64 struct {
	// File-block identifier
	Code [4]byte
	// Total length of the data after the file-block header
	Size uint32
	// Memory address the structure was located when written to disk
	OldMemoryAddress uint64
	// Index of the SDNA structure
	SDNAIndex uint32
	// Number of structures located in this file-block
	Count uint32
}

// FileBlockHeader32 represents a file-block header if the file is encoded with 32 bits.
type FileBlockHeader32 struct {
	// File-block identifier
	Code [4]byte
	// Total length of the data after the file-block header
	Size uint32
	// Memory address the structure was located when written to disk
	OldMemoryAddress uint32
	// Index of the SDNA structure
	SDNAIndex uint32
	// Number of structures located in this file-block
	Count uint32
}

// StructureDNA is stored in a file-block with code 'DNA1', just before the 'ENDB' file-block.
// This block contains all the internal structures of the Blender release the file was created in.
type StructureDNA struct {
	Identifier [4]byte
	NameID     [4]byte
	NumNames   uint32
	Names      []string
	TypeID     [4]byte
	NumTypes   uint32
	Types      []string
	LenID      [4]byte
	Lengths    []uint16
	StructID   [4]byte
	NumStructs uint32
	Structs    []struct {
		TypeIdx   uint16
		NumFields uint16
		Fields    []struct {
			TypeIdx uint16
			NameIdx uint16
		}
	}
}
