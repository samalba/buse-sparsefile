package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/samalba/buse-go/buse"
)

type BlockHeader struct {
	size uint64
}

func fatal(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s /dev/nbd0 <filename> <size in MB>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func initFile(path string, size uint64) (*os.File, error) {
	header := &BlockHeader{size}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Println("File exists, ignoring size argument.")
		fp, err := os.OpenFile(path, os.O_RDWR, 0600)
		if err != nil {
			return nil, fmt.Errorf("Cannot open file %s: %s", path, err)
		}
		buf := make([]byte, binary.Size(header))
		if _, err := fp.Read(buf); err != nil {
			return nil, fmt.Errorf("Cannot read header from file %s: %s", path, err)
		}
		bufr := bytes.NewReader(buf)
		if err := binary.Read(bufr, binary.LittleEndian, header); err != nil {
			return nil, fmt.Errorf("Cannot decode header from file %s: %s", path, err)
		}
		log.Printf("Read size %d MB from header", size/1024/1024)
		return fp, nil
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, header); err != nil {
		return nil, fmt.Errorf("Cannot encode header: %s", err)
	}
	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("Cannot create file %s: %s", path, err)
	}
	if _, err := fp.Write(buf.Bytes()); err != nil {
		return nil, fmt.Errorf("Cannot write header to file %s: %s", path, err)
	}
	log.Println("Initialized new file", path)
	return fp, nil
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
		fatal("Cannot parse size:", args[2])
	}
	initFile(args[1], size)
	os.Exit(1)
	drv := &SparseFile{}
	device, err := buse.CreateDevice(args[0], uint(size), drv)
	if err != nil {
		fatal("Cannot create the device:", err)
	}
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		if err := device.Connect(); err != nil {
			log.Printf("Buse device stopped with error: %s", err)
		} else {
			log.Println("Buse device stopped gracefully.")
		}
	}()
	<-sig
	// Received SIGTERM, cleanup
	fmt.Println("SIGINT, disconnecting...")
	device.Disconnect()
}
