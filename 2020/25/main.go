package main

import (
	"log"
)

const (
	InitialSubject = 7
	KeyA           = 1965712
	KeyB           = 19072108
)

func TransformStep(subject, i int64) int64 {
	return (i * subject) % 20201227
}

func Transform(subject, loop int64) int64 {
	v := int64(1)

	for i := int64(0); i < loop; i++ {
		v = TransformStep(v, subject)
	}

	return v
}

func CrackLoopSize(key int64) int64 {
	v := int64(1)

	i := 0
	for v != key {
		v = TransformStep(InitialSubject, v)
		i++
	}

	return int64(i)
}

func PartOne(a, b int64) int64 {
	return Transform(b, CrackLoopSize(a))
}

func main() {
	log.Printf("pt(1): %v", PartOne(KeyA, KeyB))
}
