package main

import (
	"log"
	"os"
)

// This device driver implements a sparse file block device

type SparseFile struct {
	fp *os.File
}

func (d *SparseFile) ReadAt(p []byte, off uint) error {
	n, err := d.fp.ReadAt(p, int64(off))
	log.Printf("[SparseFile] READ offset:%d len:%d\n", off, n)
	return err
}

func (d *SparseFile) WriteAt(p []byte, off uint) error {
	n, err := d.fp.WriteAt(p, int64(off))
	log.Printf("[SparseFile] WRITE offset:%d len:%d\n", off, n)
	return err
}

func (d *SparseFile) Disconnect() {
	log.Println("[SparseFile] DISCONNECT")
	d.fp.Sync()
}

func (d *SparseFile) Flush() error {
	log.Println("[SparseFile] FLUSH")
	err := d.fp.Sync()
	return err
}

func (d *SparseFile) Trim(off, length uint) error {
	log.Printf("[SparseFile] TRIM offset:%d len:%d\n", off, length)
	return nil
}
