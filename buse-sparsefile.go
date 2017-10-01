package main

import (
	"bytes"
	"flag"
	"fmt"
	"strconv"
	"log"
	"os"
	"encoding/binary"

	_"github.com/samalba/buse-go/buse"
)

type BlockHeader struct {
	size uint64
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s /dev/nbd0 <filename> <size in MB>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func initFile(path string, size uint64) error {
	header := &BlockHeader{size}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Println("File exists, ignoring size argument.")
		//TODO: read config
		return nil
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, header); err != nil {
		return err
	}
	fmt.Printf("% x", buf.Bytes())
	return nil
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) < 3 {
		usage()
	}
	size, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		fmt.Println("Cannot parse size:", args[2])
		os.Exit(1)
	}
	initFile(args[1], size)
}
