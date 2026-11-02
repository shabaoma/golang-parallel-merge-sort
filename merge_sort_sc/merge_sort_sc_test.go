package merge_sort_sc

import (
	"math/rand"
	"slices"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	for _, in := range samples() {
		assertSort(t, in)
	}
	for _, n := range []int{2, 3, 31, 32, 33, 1000, 25000} {
		assertSort(t, random(n))
	}
}

func assertSort(t *testing.T, in []int) {
	t.Helper()
	original := append([]int(nil), in...)
	got := run(in)
	want := append([]int(nil), original...)
	sort.Ints(want)
	if !slices.Equal(got, want) {
		t.Fatalf("n=%d", len(original))
	}
	if !slices.Equal(in, original) {
		t.Fatal("input mutated")
	}
}

func run(s []int) []int {
	ch := make(chan []int, 1)
	Sort(s, ch)
	return <-ch
}

func samples() [][]int {
	reverse := make([]int, 64)
	for i := range reverse {
		reverse[i] = len(reverse) - i
	}
	return [][]int{
		nil,
		{},
		{1},
		{2, 1},
		{3, 1, 2},
		{1, 1, 1},
		{1, 5, 2, 245, 574, -3, 55, -10, 0},
		reverse,
	}
}

func random(n int) []int {
	r := rand.New(rand.NewSource(1))
	s := make([]int, n)
	for i := range s {
		s[i] = r.Intn(100) - 50
	}
	return s
}
