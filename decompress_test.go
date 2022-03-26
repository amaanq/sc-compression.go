package ScCompression

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test(t *testing.T) {
	fd, err := os.Open("coc/assets/logic/achievements.csv")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	// To anyone reading this for help/as a tutorial, pay attention here:
	sc := NewDecompressor(fd) // fd is of type *os.File, which implements io.ReadWriteSeeker
	data, err := sc.Decompress()
	if err != nil {
		panic(err)
	}

	// write data to file
	data_bytes, err := ioutil.ReadAll(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data_bytes))

	decompFile, err := os.Create("coc/assets/logic/achievements_decompressed.csv")
	if err != nil {
		panic(err)
	}
	defer decompFile.Close()
	decompFile.Write(data_bytes)
}
