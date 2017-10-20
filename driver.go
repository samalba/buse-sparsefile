package main

import (
	"log"
)

// This device driver implements a sparse file block device

type SparseFile struct {
	dataset []byte
}

func (d *SparseFile) ReadAt(p []byte, off uint) error {
	log.Printf("[SparseFile] READ offset:%d len:%d\n", off, len(p))
	return nil
}

func (d *SparseFile) WriteAt(p []byte, off uint) error {
	log.Printf("[SparseFile] WRITE offset:%d len:%d\n", off, len(p))
	return nil
}

func (d *SparseFile) Disconnect() {
	log.Println("[SparseFile] DISCONNECT")
}

func (d *SparseFile) Flush() error {
	log.Println("[SparseFile] FLUSH")
	return nil
}

func (d *SparseFile) Trim(off, length uint) error {
	log.Printf("[SparseFile] TRIM offset:%d len:%d\n", off, length)
	return nil
}
