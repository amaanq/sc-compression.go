// Credit: https://github.com/jeanbmar/sc-compression/blob/master/src/sc-compression.js

package ScCompression

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"strings"

	"github.com/kjk/lzma"
)

type Signature int

const (
	NONE Signature = iota
	LZMA           // starts with 5D 00 00 04
	SC             // starts with SC
	SCLZ           // starts with SC and contains SCLZ
	SIG            // starts with Sig:
)

type ScCompression struct {
	buffer []byte
}

func New(fp string) (*ScCompression, error) {
	bytes, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	return &ScCompression{
		buffer: bytes,
	}, nil
}

func NewFromBuffer(buffer []byte) (*ScCompression, error) {
	return &ScCompression{
		buffer: buffer,
	}, nil
}

func (s *ScCompression) Decompress() []byte {
	signature := s.readSignature()
	switch signature {
	case NONE:
		return s.buffer
	case LZMA:
		return s.decompressLZMA()
	case SC:
		return s.decompressSC()
	case SCLZ:
		return s.decompressSCLZ()
	case SIG:
		return s.decompressSIG()
	default:
		panic("invalid sig")
	}
}

func (s *ScCompression) decompressLZMA() []byte {
	uncompressedSize := int32(binary.LittleEndian.Uint32(s.buffer[5:9]))
	if uncompressedSize == -1 {
		buf := make([]byte, 9)
		copy(buf, s.buffer[:9])
		buf = append(append(buf, []byte{0xFF, 0xFF, 0xFF, 0xFF}...), s.buffer[9:]...)
		s.buffer = buf
	} else {
		buf := make([]byte, 9)
		copy(buf, s.buffer[:9])
		buf = append(append(buf, []byte{0x00, 0x00, 0x00, 0x00}...), s.buffer[9:]...)
		s.buffer = buf
	}
	uncompressedBuf := make([]byte, uncompressedSize)
	n, err := lzma.NewReader(bytes.NewReader(s.buffer)).Read(uncompressedBuf)
	if err != nil {
		panic(err)
	}
	if n != int(uncompressedSize) {
		panic("bad value read")
	}
	return uncompressedBuf
}

func (s *ScCompression) decompressSC() []byte {
	s.buffer = s.buffer[26:]
	return s.decompressLZMA()
}

func (s *ScCompression) decompressSCLZ() []byte {
	return nil
}

func (s *ScCompression) decompressSIG() []byte {
	s.buffer = s.buffer[68:]
	return s.decompressLZMA()
}

func (s *ScCompression) readSignature() Signature {
	if hex.EncodeToString(s.buffer[:3]) == "5d0000" {
		return LZMA
	} else if strings.ToLower(string(s.buffer[:2])) == "sc" {
		if len(s.buffer) > 30 && strings.ToLower(string(s.buffer[26:30])) == "sclz" {
			return SCLZ
		}
		return SC
	} else if strings.ToLower(string(s.buffer[:4])) == "sig:" {
		return SIG
	}
	return NONE
}
