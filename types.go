package ScCompression

import "io"

type Signature int

const (
	NONE Signature = iota
	LZMA           // starts with 5D 00 00 04
	SC             // starts with SC
	SCLZ           // starts with SC and contains SCLZ
	SIG            // starts with Sig:
)

type Compressor struct {
	writer     io.ReadWriteSeeker
	lzmawriter io.ReadWriter
}

type Decompressor struct {
	reader     io.ReadWriteSeeker
	lzmareader io.ReadWriter
}
