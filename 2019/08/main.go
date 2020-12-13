package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	high = 6
	wide = 25
	long = high * wide
)

func PartOne(img []int) int {
	numberOfImages := len(img) / long
	var counts []map[int]int
	var idx int

	for i := 0; i < numberOfImages; i++ {
		m := map[int]int{}
		o := i * long

		for j := 0; j < long; j++ {
			m[img[o+j]] += 1
		}

		counts = append(counts, m)

		if counts[i][0] < counts[idx][0] {
			idx = i
		}
	}

	return counts[idx][1] * counts[idx][2]
}

func PartTwo(img []int) []int {
	numberOfImages := len(img) / long
	decoded := make([]int, long)

	for i := 0; i < long; i++ {
		decoded[i] = 2 // default transparent

		for j := 0; j < numberOfImages; j++ {
			o := j * long
			pixel := img[o+i]

			if pixel != 2 {
				decoded[i] = pixel
				break
			}
		}
	}

	return decoded
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	var img []int

	s := bufio.NewScanner(bufio.NewReader(f))
	s.Scan()

	for _, b := range s.Bytes() {
		img = append(img, int(b-'0'))
	}

	log.Printf("pt(1): %v", PartOne(img))

	pt2 := PartTwo(img)

	for i := 0; i < high; i++ {
		for j := 0; j < wide; j++ {
			pixel := pt2[((i * wide) + j)]

			if pixel == 0 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("\u25A0")
			}
		}

		fmt.Print("\n")
	}
}
