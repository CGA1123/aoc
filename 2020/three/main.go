package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	f, e := os.Open("input.txt")
	if e != nil {
		log.Printf("opening file: %v", e)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)

	tree := byte('#')
	var line []byte
	var i int
	var count int
	var err error

	for {
		line, err = r.ReadBytes('\n')
		if err != nil {
			break
		}

		if line[i] == tree {
			count = count + 1
		}

		i = (i + 3) % (len(line) - 1)
	}

	if err != nil && err != io.EOF {
		log.Printf("error: %v", err)
		return
	}

	log.Printf("count (pt. 1): %v", count)
}
