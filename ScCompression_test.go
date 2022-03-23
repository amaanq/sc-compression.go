package ScCompression

import (
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	sc, err := New("animations.csv")
	if err != nil {
		panic(err)
	}
	data := sc.Decompress()
	if false {fmt.Println(string(data))}
	f, _ := os.Create("dc.csv")
	f.Write(data)
	f.Close()
}
