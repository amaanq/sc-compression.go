// Credit: https://github.com/jeanbmar/sc-compression/blob/master/src/sc-compression.js

package ScCompression

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"strings"

	"github.com/ulikunitz/xz/lzma"
)

// Default method of instantiating this package, we need a seeker to be able to seek to different indices of the reader. This is the only way to do this.
func NewDecompressor(reader io.ReadWriteSeeker) *Decompressor {
	return &Decompressor{
		reader:     reader,
		lzmareader: bytes.NewBuffer(nil),
	}
}

// This will return an io.Reader so you can copy this to a file or readall (not recommended) into a buffer
func (d *Decompressor) Decompress() (io.Reader, error) {
	signature := d.readSignature()
	switch signature {
	case NONE:
		return d.reader, nil
	case LZMA:
		return d.decompressLZMA(0)
	case SC:
		return d.decompressSC()
	case SCLZ:
		return d.decompressSCLZ()
	case SIG:
		return d.decompressSIG()
	default:
		panic("invalid sig")
	}
}

func (d *Decompressor) decompressLZMA(offset int) (io.Reader, error) {
	d.reader.Seek(int64(offset), io.SeekStart) // signature-dependent offset

	uncompressedSizeBuf := make([]byte, 4)
	_, err := d.reader.Seek(5, io.SeekCurrent) // seek 5 to read an int here
	if err != nil {
		panic(err)
	}

	n, err := d.reader.Read(uncompressedSizeBuf)
	if err != nil {
		panic(err)
	}
	_, err = d.reader.Seek(int64(0-n), io.SeekCurrent) // "unread" what was read
	if err != nil {
		panic(err)
	}

	_, err = d.reader.Seek(-5, io.SeekCurrent) // seek back 5 ("beginning")
	if err != nil {
		panic(err)
	}

	uncompressedSize := int32(binary.LittleEndian.Uint32(uncompressedSizeBuf))
	insert := make([]byte, 0)
	if uncompressedSize == -1 {
		insert = []byte{0xFF, 0xFF, 0xFF, 0xFF}
	} else {
		insert = []byte{0x00, 0x00, 0x00, 0x00}
	}

	//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
	// Copy 0-9, insert, then 9-end to lzma reader
	//_______________________________________________________________________

	nine := make([]byte, 9)
	_, err = d.reader.Read(nine) // read 9 bytes of reader
	if err != nil {
		return nil, err
	}
	_, err = d.reader.Seek(0-9, io.SeekCurrent) // "unread" what was read
	if err != nil {
		panic(err)
	}
	_, err = d.lzmareader.Write(nine) // write 9 bytes to lzmareader
	if err != nil {
		return nil, err
	}

	_, err = d.lzmareader.Write(insert) // write 4 bytes "insert" to lzmareader
	if err != nil {
		return nil, err
	}

	_, err = d.reader.Seek(9, io.SeekCurrent) // seek to 9
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(d.lzmareader, d.reader) // copy 9 to end bytes of reader to lzmareader
	if err != nil {
		return nil, err
	}
	d.reader.Seek(-9, io.SeekCurrent) // seek back to "beginning"

	if reader, err := lzma.NewReader(d.lzmareader); err != nil {
		return nil, err
	} else {
		return reader, nil
	}
}

func (d *Decompressor) decompressSC() (io.Reader, error) {
	return d.decompressLZMA(26)
}

func (d *Decompressor) decompressSCLZ() (io.Reader, error) {
	return nil, nil // TODO
}

func (d *Decompressor) decompressSIG() (io.Reader, error) {
	return d.decompressLZMA(68)
}

func (d *Decompressor) readSignature() Signature {
	bufferof31 := make([]byte, 31)
	_, err := d.reader.Read(bufferof31)
	if err != nil {
		panic(err) // user: just make sure the reader is readable
	}

	if hex.EncodeToString(bufferof31[:3]) == "5d0000" {
		return LZMA
	} else if strings.ToLower(string(bufferof31[:2])) == "sc" {
		if len(bufferof31) > 30 && strings.ToLower(string(bufferof31[26:30])) == "sclz" {
			return SCLZ
		}
		return SC
	} else if strings.ToLower(string(bufferof31[:4])) == "sig:" {
		return SIG
	}
	return NONE
}
