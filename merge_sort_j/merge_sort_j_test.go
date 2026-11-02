package merge_sort_j

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
	for _, n := range []int{2, 3, 31, 32, 33, 200} {
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
	ch := make(chan int, len(s))
	Sort(s, ch)
	got := make([]int, 0, len(s))
	for v := range ch {
		got = append(got, v)
	}
	return got
}

func samples() [][]int {
	return [][]int{
		nil,
		{},
		{1},
		{2, 1},
		{3, 1, 2},
		{1, 1, 1},
		{1, 5, 2, 245, 574, -3, 55, -10, 0},
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
