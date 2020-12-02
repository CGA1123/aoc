package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func valid(r *bufio.Reader, min, max int, pat byte) (bool, error) {
	matches := 0
	var b byte
	var err error

	for {
		b, err = r.ReadByte()
		if err != nil {
			return false, err
		}

		if b == '\n' {
			break
		}

		if b == pat {
			matches = matches + 1
		}
	}

	return (min <= matches) && (matches <= max), nil
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
	count := 0

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

		b, err := valid(r, min, max, pat)
		if err != nil {
			break
		}

		if b {
			count = count + 1
		}
	}
	if err != nil && err != io.EOF {
		log.Printf("error parsing: %v", err)
		return
	}

	log.Printf("count: %v", count)
}
