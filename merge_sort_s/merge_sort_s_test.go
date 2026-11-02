package merge_sort_s

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
	for _, n := range []int{31, 32, 33, 1000, 25000} {
		assertSort(t, random(n))
	}
}

func TestMerge(t *testing.T) {
	src := []int{1, 3, 5, 7, 9, 33, 56, 78, 2, 4, 6, 8, 14, 34, 45, 66}
	dst := make([]int, len(src))
	merge(src, dst, 0, 8, 16)
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 14, 33, 34, 45, 56, 66, 78}
	if !slices.Equal(dst, want) {
		t.Fatalf("got %v", dst)
	}
}

func assertSort(t *testing.T, in []int) {
	t.Helper()
	original := append([]int(nil), in...)
	got := Sort(in)
	want := append([]int(nil), original...)
	sort.Ints(want)
	if !slices.Equal(got, want) {
		t.Fatalf("n=%d", len(original))
	}
	if !slices.Equal(in, original) {
		t.Fatal("input mutated")
	}
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
