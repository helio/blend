package blend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFile_readExampleHeader(t *testing.T) {
	name := "cubus-animated.blend"
	r, err := readExample(name)
	if err != nil {
		t.Fatalf("Unable to read example file '%s': %s", name, err)
	}
	defer r.Close()

	f, err := NewFile(r)
	if err != nil {
		t.Fatalf("Expected nil error, got: %v", err)
	}

	if string(f.header.Identifier[:7]) != "BLENDER" {
		t.Errorf("expected Identifier 'BLENDER', got: '%v'", f.header.Identifier)
	}
	if f.header.PointerSize != '-' {
		t.Errorf("expected PointerSize 8 bytes / 64 bits ('-'), got: '%v'", f.header.PointerSize)
	}
	if f.header.Endianness != 'v' {
		t.Errorf("expected little endian byte order ('v'), got: '%v'", f.header.Endianness)
	}
	if string(f.header.Version[:3]) != "280" {
		t.Errorf("expected version 280, got: '%v'", f.header.Version)
	}
}

func TestNewFile_readHeader(t *testing.T) {
	testTable := []struct {
		name        string
		pointerSize byte
		endianness  byte
		version     string
		order       binary.ByteOrder
	}{
		{
			name:        "test 32Bit Pointer, LittleEndian, Version 280",
			pointerSize: '_',
			endianness:  'v',
			version:     "280",
			order:       binary.LittleEndian,
		},
		{
			name:        "test 64Bit Pointer, LittleEndian, Version 280",
			pointerSize: '-',
			endianness:  'v',
			version:     "280",
			order:       binary.LittleEndian,
		},
		{
			name:        "test 64Bit Pointer, BigEndian, Version 280",
			pointerSize: '-',
			endianness:  'V',
			version:     "280",
			order:       binary.BigEndian,
		},
		{
			name:        "test 64Bit Pointer, BigEndian, Version 100",
			pointerSize: '-',
			endianness:  'V',
			version:     "100",
			order:       binary.BigEndian,
		},
	}

	for i, f := range testTable {
		t.Run(fmt.Sprintf("#%d %s", i, f.name), func(t *testing.T) {
			data := bytes.NewBuffer(header(f.pointerSize, f.endianness, f.version))
			file, err := NewFile(data)
			if err != nil {
				t.Errorf("expected nil error, got '%s'", err)
			}

			if file.header.PointerSize != f.pointerSize {
				t.Errorf("expected '%v', got '%v'", f.pointerSize, file.header.PointerSize)
			}
			if file.header.Endianness != f.endianness {
				t.Errorf("expected '%v', got '%v'", f.pointerSize, file.header.Endianness)
			}
			if file.order != f.order {
				t.Errorf("expected '%v', got '%v'", f.order, file.order)
			}
			if string(file.header.Version[:3]) != f.version {
				t.Errorf("expected '%v', got '%v'", f.version, file.header.Version)
			}
		})
	}
}

func TestNewFile_readPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("read() should panic if called before readHeader()")
		}
	}()
	f := bytes.NewBuffer(header('-', 'v', "280"))
	file := File{
		r: f,
	}
	var data interface{}
	err := file.read(1, &data)
	if err != nil {
		t.Fatalf("expected read() to panic, got error instead: %s", err)
	}
}

func TestNewFile_headerInvalidIdentifier(t *testing.T) {
	f := bytes.NewBuffer(rawHeader("NOBLEND", '-', 'v', "280"))
	_, err := NewFile(f)
	if err == nil {
		t.Error("expected NewFile to error in readHeader() because of invalid identifier")
	}
	expected := "blend: invalid identifier"
	if err.Error() != expected {
		t.Errorf("expected error '%s', got: '%s'", expected, err)
	}
}

func header(pointerSize, endianness byte, version string) []byte {
	return rawHeader("BLENDER", pointerSize, endianness, version)
}

func rawHeader(identifier string, pointerSize byte, endianness byte, version string) []byte {
	var b []byte
	b = append(b, identifier...)
	b = append(b, pointerSize, endianness)
	return append(b, version...)
}

func readExample(name string) (io.ReadCloser, error) {
	return os.Open(filepath.Join("./examples", name))
}
