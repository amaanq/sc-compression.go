# **A Go package to decompress Supercell game asset files**

## **Installation**
```go
go get -u github.com/amaanq/sc-compression.go
```
<br />

## **Usage**
### Let's say you have a Supercell game file called "missions.csv" located in the project directory and you'd like to decompress it, here's how to do it:
```go
package main

import (
    "ioutil"
    "os"

    "github.com/amaanq/sc-compression.go"
)

func main() {
    fd, err := os.Open("missions.csv")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

    decompressor := ScCompression.NewDecompressor(fd)
    data, err := decompressor.Decompress()
    if err != nil {
        panic(err)
    }

    data_bytes, err := ioutil.ReadAll(data)
    fmt.Println(string(data_bytes))

    decompFile, err := os.Create("missions_decompressed.csv")
    if err != nil {
        panic(err)
    }
    defer decompFile.Close()
    decompFile.Write(data_bytes)
}
```
<br />

## **Want to Contribute?**  
### *Submit a pull request or reach out to me (check my profile)*
<br />

## **TODO**
### - [x] *Use readers instead of byte arrays to prevent out of memory panics*
### - [ ] *Implement LZHAM for the SCLZ signature*
### - [ ] *Add compressor capability*
<br />
<br />

## **Special Credits to @jeanbmar as his tool in Javascript was most helpful in making this myself: [sc-compression](https://github.com/jeanbmar/sc-compression)**