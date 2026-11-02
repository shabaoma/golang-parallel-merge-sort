package main

import (
	"fmt"
	"sort"
	"sort_test/merge_sort_j"
	"sort_test/merge_sort_s"
	"sort_test/merge_sort_sc"
	"sort_test/merge_sort_ssc"
	"sort_test/quick_sort"
	"testing"
)

var slice_length = 10000000

func TestQuickSort(t *testing.T) {
	raw := GenerateSlice(slice_length)
	s := make([]int, len(raw))
	copy(s, raw)
	sort.Sort(quick_sort.Nums(s))
	mustMatch(t, raw, s)
}

func TestMergeSortS(t *testing.T) {
	raw := GenerateSlice(slice_length)
	got := merge_sort_s.Sort(raw)
	mustMatch(t, raw, got)
}

// func TestMergeSortSC(t *testing.T) {
// 	s := GenerateSlice(slice_length)
// 	ch := make(chan []int)
// 	merge_sort_sc.Sort(s, ch)
// }

func TestMergeSortSSC(t *testing.T) {
	raw := GenerateSlice(slice_length)
	got := merge_sort_ssc.Sort(raw)
	mustMatch(t, raw, got)
}

func TestMergeSortJ(t *testing.T) {
	raw := GenerateSlice(slice_length)
	ch := make(chan int, len(raw))
	merge_sort_j.Sort(raw, ch)
	got := make([]int, 0, len(raw))
	for v := range ch {
		got = append(got, v)
	}
	mustMatch(t, raw, got)
}

func mustMatch(t *testing.T, raw, got []int) {
	t.Helper()
	if len(got) != len(raw) || !sort.IntsAreSorted(got) {
		t.Fatal("not sorted")
	}
	var a, b int
	for i := range raw {
		a += raw[i]
		b += got[i]
	}
	if a != b {
		t.Fatal("checksum")
	}
}

// func BenchmarkQuickSort(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		s := quick_sort.Nums(GenerateSlice(slice_length))
// 		sort.Sort(s)
// 	}
// }

// func BenchmarkMergeSortS(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		s := GenerateSlice(slice_length)
// 		_ = merge_sort_s.Sort(s)
// 	}
// }

// func BenchmarkMergeSortSC(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		s := GenerateSlice(slice_length)
// 		ch := make(chan []int)
// 		merge_sort_sc.Sort(s, ch)
// 	}
// }

// func BenchmarkMergeSortSSC(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		s := GenerateSlice(slice_length)
// 		_ = merge_sort_ssc.Sort(s)
// 	}
// }

// func BenchmarkMergeSortJ(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		s := GenerateSlice(slice_length)
// 		ch := make(chan int, len(s))
// 		merge_sort_j.Sort(s, ch)
// 	}
// }

func BenchmarkQuickSortIncrease(b *testing.B) {
	for n := 10; n <= slice_length; n *= 10 {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			data := GenerateSlice(n)
			buf := make([]int, n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				copy(buf, data)
				b.StartTimer()
				sort.Sort(quick_sort.Nums(buf))
			}
		})
	}
}

func BenchmarkMergeSortSIncrease(b *testing.B) {
	for n := 10; n <= slice_length; n *= 10 {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			data := GenerateSlice(n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = merge_sort_s.Sort(data)
			}
		})
	}
}

func BenchmarkMergeSortSCIncrease(b *testing.B) {
	for n := 10; n <= slice_length; n *= 10 {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			data := GenerateSlice(n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ch := make(chan []int, 1)
				merge_sort_sc.Sort(data, ch)
				_ = <-ch
			}
		})
	}
}

func BenchmarkMergeSortSSCIncrease(b *testing.B) {
	for n := 10; n <= slice_length; n *= 10 {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			data := GenerateSlice(n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = merge_sort_ssc.Sort(data)
			}
		})
	}
}

func BenchmarkMergeSortJIncrease(b *testing.B) {
	for n := 10; n <= slice_length; n *= 10 {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			data := GenerateSlice(n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ch := make(chan int, len(data))
				merge_sort_j.Sort(data, ch)
				for range ch {
				}
			}
		})
	}
}
