package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	// Flags
	_filePath1   *string
	_filePath2   *string
	_outFileName *string
	_print       *bool

	// Writer
	_outWriter *bufio.Writer
)

func main() {
	log.Println("Starting compareBytes")

	// Set and parse flags
	_filePath1 = flag.String("f1", "", "Path of first file")
	_filePath2 = flag.String("f2", "", "Path of second file")
	_outFileName = flag.String("o", "", "If present out file with compare results.")
	_print = flag.Bool("p", false, "Print mismatches, default is true.")
	flag.Parse()

	// Ensure user passed filenames
	if *_filePath1 == "" || *_filePath2 == "" {
		log.Println("f1 and f2 flags must be set")
		flag.PrintDefaults()
		return
	}

	// Creater writer if flag was passed
	writeOut := false
	if *_outFileName != "" {
		writeOut = true
		outFile, err := os.Create(*_outFileName)
		if err != nil {
			log.Printf("Error creating out file: %s\n", err.Error())
			return
		}

		_outWriter = bufio.NewWriter(outFile)
		defer _outWriter.Flush()
	}

	// Open files
	file1, err := os.Open(*_filePath1)
	if err != nil {
		log.Printf("Error opening %s: %s\n", *_filePath1, err.Error())
		return
	}

	file2, err := os.Open(*_filePath2)
	if err != nil {
		log.Printf("Error opening %s: %s\n", *_filePath2, err.Error())
		return
	}

	// Get file sizes
	file1Stat, err := file1.Stat()
	if err != nil {
		log.Printf("Error getting stats for %s: %s\n", *_filePath1, err.Error())
		return
	}

	file2Stat, err := file2.Stat()
	if err != nil {
		log.Printf("Error getting stats for %s: %s\n", *_filePath2, err.Error())
		return
	}

	file1Size := file1Stat.Size()
	file2Size := file2Stat.Size()

	// Create buf reader
	reader1 := bufio.NewReader(file1)
	reader2 := bufio.NewReader(file2)

	count1 := 0
	count2 := 0
	mismatches := 0

	for {
		byte1, err := reader1.ReadByte()
		if err != nil {
			// If EOF, print how much was read of each one
			if errors.Is(err, io.EOF) {
				log.Printf("Read all of %s\n", *_filePath1)
				log.Printf("Read %d bytes of %s\n", count1, *_filePath1)
				log.Printf("Read %d bytes of %s\n", count2, *_filePath2)
				// Print if file2 has unread bytes
				if count2 != int(file2Size) {
					log.Printf("%d bytes still remaining to be read for %s\n", int(file2Size)-count2, *_filePath2)
				}
				log.Printf("Total number of mismatches: %d\n", mismatches)
				return
			}

			log.Printf("Error reading byte for %s: %s", *_filePath1, err.Error())
			return
		}
		count1++

		byte2, err := reader2.ReadByte()
		if err != nil {
			// If EOF, print how much was read of each one
			if errors.Is(err, io.EOF) {
				log.Printf("Read all of %s\n", *_filePath2)
				log.Printf("Read %d bytes of %s\n", count1, *_filePath1)
				log.Printf("Read %d bytes of %s\n", count2, *_filePath2)
				// Print if file2 has unread bytes
				if count1 != int(file1Size) {
					log.Printf("%d bytes still remaining to be read for %s\n", int(file1Size)-count1, *_filePath1)
				}
				log.Printf("Total number of mismatches: %d\n", mismatches)
				return
			}

			log.Printf("Error reading byte for %s: %s", *_filePath2, err.Error())
			return
		}
		count2++

		// Compare byte
		if byte1 != byte2 {
			mismatches++

			logString := fmt.Sprintf("Pos: %2d, f1: %02x, f2: %02x \n", count1, byte1, byte2)

			// Log it if selected
			if *_print {
				log.Printf(logString)
			}

			// Write it out if selected
			if writeOut {
				_outWriter.WriteString(logString)
			}
		}
	}

}
