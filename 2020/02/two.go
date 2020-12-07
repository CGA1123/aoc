package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func partOne(bytes []byte, min, max int, pat byte) bool {
	var matches int
	for _, b := range bytes {
		if b == pat {
			matches = matches + 1
		}
	}

	return (min <= matches) && (matches <= max)
}

func partTwo(bytes []byte, i, j int, pat byte) bool {
	size := len(bytes) - 1 // new line character

	ib := i < size && bytes[i] == pat
	ij := j < size && bytes[j] == pat

	return ib != ij
}

func valid(r *bufio.Reader, min, max int, pat byte) (bool, bool, error) {
	bytes, err := r.ReadBytes('\n')
	if err != nil {
		return false, false, err
	}

	return partOne(bytes, min, max, pat), partTwo(bytes, min-1, max-1, pat), nil
}

func readInt(r *bufio.Reader, delim byte) (int, error) {
	str, err := readStr(r, delim)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(str)
}

func readStr(r *bufio.Reader, delim byte) (string, error) {
	str, err := r.ReadString(delim)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(str, string(delim)), nil
}

func main() {
	f, e := os.Open("input.txt")
	if e != nil {
		log.Printf("opening file: %v", e)
		return
	}
	defer f.Close()

	var err error
	r := bufio.NewReader(f)
	var partOne int
	var partTwo int

	for {
		min, err := readInt(r, '-')
		if err != nil {
			break
		}

		max, err := readInt(r, ' ')
		if err != nil {
			break
		}

		pat, err := r.ReadByte()
		if err != nil {
			break
		}

		// eat ': '
		_, err = r.Discard(2)
		if err != nil {
			break
		}

		b, c, err := valid(r, min, max, pat)
		if err != nil {
			break
		}

		if b {
			partOne = partOne + 1
		}

		if c {
			partTwo = partTwo + 1
		}
	}
	if err != nil && err != io.EOF {
		log.Printf("error parsing: %v", err)
		return
	}

	log.Printf("count (pt 1): %v", partOne)
	log.Printf("count (pt 2): %v", partTwo)
}
