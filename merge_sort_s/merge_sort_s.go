package merge_sort_s

const insertionCutoff = 32

func Sort(s []int) []int {
	n := len(s)
	if n <= 1 {
		return s
	}
	dst := make([]int, n)
	src := make([]int, n)
	copy(dst, s)
	copy(src, s)
	sequential(src, dst, 0, n)
	return dst
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
