package merge_sort_sc

import "sync"

const insertionCutoff = 32

func Sort(s []int, res chan []int) {
	n := len(s)
	if n <= 1 {
		res <- s
		return
	}
	dst := make([]int, n)
	src := make([]int, n)
	copy(dst, s)
	copy(src, s)
	parallel(src, dst, 0, n)
	res <- dst
}

func parallel(src, dst []int, lo, hi int) {
	if hi-lo <= insertionCutoff {
		insertion(dst, lo, hi)
		return
	}
	mid := lo + (hi-lo)>>1
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		parallel(dst, src, lo, mid)
	}()
	parallel(dst, src, mid, hi)
	wg.Wait()
	merge(src, dst, lo, mid, hi)
}

func insertion(a []int, lo, hi int) {
	for i := lo + 1; i < hi; i++ {
		v := a[i]
		j := i
		for j > lo && a[j-1] > v {
			a[j] = a[j-1]
			j--
		}
		a[j] = v
	}
}

func merge(src, dst []int, lo, mid, hi int) {
	i, j, k := lo, mid, lo
	for i < mid && j < hi {
		if src[i] <= src[j] {
			dst[k] = src[i]
			i++
		} else {
			dst[k] = src[j]
			j++
		}
		k++
	}
	for i < mid {
		dst[k] = src[i]
		i++
		k++
	}
	for j < hi {
		dst[k] = src[j]
		j++
		k++
	}
}
