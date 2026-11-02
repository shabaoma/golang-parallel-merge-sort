package merge_sort_ssc

import (
	"math/bits"
	"runtime"
	"sync"
)

const (
	insertionCutoff   = 32
	parallelThreshold = 10000
)

func Sort(s []int) []int {
	n := len(s)
	if n <= 1 {
		return s
	}
	dst := make([]int, n)
	src := make([]int, n)
	copy(dst, s)
	copy(src, s)
	parallel(src, dst, 0, n, parallelDepth())
	return dst
}

func parallelDepth() int {
	p := runtime.GOMAXPROCS(0)
	if p <= 1 {
		return 0
	}
	return bits.Len(uint(p)) - 1
}

func parallel(src, dst []int, lo, hi, depth int) {
	if depth <= 0 || hi-lo <= parallelThreshold {
		sequential(src, dst, lo, hi)
		return
	}
	mid := lo + (hi-lo)>>1
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		parallel(dst, src, lo, mid, depth-1)
	}()
	parallel(dst, src, mid, hi, depth-1)
	wg.Wait()
	merge(src, dst, lo, mid, hi)
}

func sequential(src, dst []int, lo, hi int) {
	if hi-lo <= insertionCutoff {
		insertion(dst, lo, hi)
		return
	}
	mid := lo + (hi-lo)>>1
	sequential(dst, src, lo, mid)
	sequential(dst, src, mid, hi)
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
