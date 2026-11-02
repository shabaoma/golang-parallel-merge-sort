package main

import (
	"math/rand"
)

func main() {
}

func GenerateSlice(size int) []int {
	r := rand.New(rand.NewSource(1))
	s := make([]int, size)
	for i := 0; i < size; i++ {
		s[i] = r.Int()
	}
	return s
}
