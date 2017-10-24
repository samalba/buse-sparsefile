package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/samalba/buse-go/buse"
)

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s /dev/nbd0 <filename> <size in MB>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func initFile(path string, size uint64) (*os.File, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Println("File exists, ignoring size argument.")
		fp, err := os.OpenFile(path, os.O_RDWR, 0600)
		if err != nil {
			return nil, fmt.Errorf("Cannot open file %s: %s", path, err)
		}
		return fp, nil
	}
	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("Cannot create file %s: %s", path, err)
	}
	if err := fp.Truncate(int64(size)); err != nil {
		return nil, fmt.Errorf("Cannot set size to %d MB: %s", size/1024/1024, err)
	}
	log.Printf("Initialized new file %s with size %d MB", path, size/1024/1024)
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
	// Convert size to Bytes
	size = size * 1024 * 1024
	fp, err := initFile(args[1], size)
	if err != nil {
		fatal("Cannot initialize file", err)
	}
	drv := &SparseFile{fp}
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
