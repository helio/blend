package blend

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type File struct {
	r            io.Reader
	header       *FileHeader
	order        binary.ByteOrder
	pointerSize  uint8
	fileBlocks64 map[string]FileBlock64
	fileBlocks32 map[string]FileBlock32
}

// NewFile initializes the File struct and reads the header.
// This automatically determines the byte order, after which the rest of the file can be read if needed.
func NewFile(r io.Reader) (*File, error) {
	f := File{
		r: r,
	}
	if err := f.readHeader(); err != nil {
		return nil, err
	}
	if f.pointerSize == 64 {
		f.fileBlocks64 = make(map[string]FileBlock64)
	} else {
		f.fileBlocks32 = make(map[string]FileBlock32)
	}

	return &f, nil
}

// readNextBytes reads number of bytes from file.
// shamelessly stolen from https://www.jonathan-petitcolas.com/2014/09/25/parsing-binary-files-in-go.html
func readNextBytes(r io.Reader, n int) ([]byte, error) {
	bytes := make([]byte, n)

	// FIXME: take care about `n`
	_, err := r.Read(bytes)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}

// readHeader reads the first 12 bytes which represent a blender file header.
// most importantly the byte order is determined upon which the rest of the file can be read successfully.
func (f *File) readHeader() error {
	header := FileHeader{}
	data, err := readNextBytes(f.r, 12)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(data)

	// determine byte order before trying to parse
	// byte order is within the file header at offset 8, c type `char`
	var order binary.ByteOrder = binary.LittleEndian
	if data[8] == 'V' {
		order = binary.BigEndian
	}
	if err = binary.Read(buffer, order, &header); err != nil {
		return err
	}
	identifier := string(header.Identifier[:])
	if identifier != "BLENDER" {
		return errors.New("blend: invalid identifier")
	}

	f.pointerSize = 64
	if header.PointerSize == '_' {
		f.pointerSize = 32
	}

	f.order = order
	f.header = &header
	return nil
}

// readFileBlocks reads all file blocks and builds up the cache structure.
func (f *File) readFileBlocks() error {
	for {
		if f.pointerSize == 64 {
			header, err := f.readFileBlockHeader64()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
			data, err := readNextBytes(f.r, int(header.Size))
			if err != nil {
				return err
			}

			f.fileBlocks64[byteSliceToString(header.Code[:])] = FileBlock64{
				header: header,
				data:   data,
			}
		} else {
			header, err := f.readFileBlockHeader32()
			if err != nil {
				return err
			}
			data, err := readNextBytes(f.r, int(header.Size))
			if err != nil {
				return err
			}

			f.fileBlocks32[byteSliceToString(header.Code[:])] = FileBlock32{
				header: header,
				data:   data,
			}
		}
	}
}

func (f *File) readFileBlockHeader64() (*FileBlockHeader64, error) {
	header := FileBlockHeader64{}
	return &header, f.read(24, &header)
}
func (f *File) readFileBlockHeader32() (*FileBlockHeader32, error) {
	header := FileBlockHeader32{}
	return &header, f.read(20, &header)
}

// read reads the next `n` bytes into the structured `data`.
// This function panics if byte order has not been determined yet, which should be done when initializing File.
func (f *File) read(n int, data interface{}) error {
	if f.order == nil {
		panic("blend: unable to read bytes before reading header")
	}
	binData, err := readNextBytes(f.r, n)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(binData)

	return binary.Read(buffer, f.order, data)
}

func byteSliceToString(s []byte) string {
	n := bytes.IndexByte(s, 0)
	// if byte array doesn't contain any 0 bytes
	if n == -1 {
		return string(s[:])
	}
	return string(s[:n])
}
