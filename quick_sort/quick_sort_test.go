package quick_sort

import (
	"slices"
	"sort"
	"testing"
)

func TestSortInterface(t *testing.T) {
	cases := [][]int{
		nil,
		{},
		{1},
		{2, 1},
		{1, 1, 1},
		{1, 5, 2, 245, 574, -3, 55, -10, 0},
	}
	for _, in := range cases {
		s := append([]int(nil), in...)
		sort.Sort(Nums(s))
		want := append([]int(nil), in...)
		sort.Ints(want)
		if !slices.Equal(s, want) {
			t.Fatalf("got %v want %v", s, want)
		}
	}
}
