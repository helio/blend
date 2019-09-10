package blend

type FileHeader struct {
	Identifier  [7]byte
	PointerSize byte
	Endianness  byte
	Version     [3]byte
}

type FileBlock64 struct {
	header *FileBlockHeader64
	data   *FileBlockData
}
type FileBlockHeader64 struct {
	Code             [4]byte
	Size             uint32
	OldMemoryAddress uint64
	SDNAIndex        uint32
	Count            uint32
}
type FileBlockHeader32 struct {
	Code             [4]byte
	Size             uint32
	OldMemoryAddress uint32
	SDNAIndex        uint32
	Count            uint32
}
type FileBlockData struct {
}
type StructureDNA struct {
}
