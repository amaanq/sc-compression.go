package ScCompression

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	sc, err := New("helpshift.csv")
	if err != nil {
		panic(err)
	}
	data := sc.Decompress()
	fmt.Println(string(data))
}
