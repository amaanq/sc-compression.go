package ScCompression

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test(t *testing.T) {
	fd, err := os.OpenFile("coc/assets/logic/achievements.csv", os.O_RDWR, 0666)
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
	by, err := ioutil.ReadAll(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(by))
}
