package blend

type FileHeader struct {
	Identifier  [7]byte
	PointerSize byte
	Endianness  byte
	Version     [3]byte
}

type FileBlock struct {
}
type FileBlockHeader struct {
}
type FileBlockData struct {
}
type StructureDNA struct {
}
