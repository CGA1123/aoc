package main

import (
	"io/ioutil"
	"log"
)

func main() {
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("err reading: %v", err)
	}

	// trim newline character
	b = b[:len(b)-1]

	log.Printf("%v\n", startOfPacket(b)+1)
	log.Printf("%v\n", startOfMesage(b)+1)
}

func startOf(data []byte, size int) int {
	for i := size - 1; i < len(data); i++ {
		if uniq(data[i-(size-1) : i+1]) {
			return i
		}

	}

	return -1
}

func startOfPacket(data []byte) int {
	return startOf(data, 4)

}

func startOfMesage(data []byte) int {
	return startOf(data, 14)
}

func uniq(data []byte) bool {
	m := make(map[byte]struct{}, len(data))

	for _, b := range data {
		m[b] = struct{}{}
	}

	return len(m) == len(data)
}
