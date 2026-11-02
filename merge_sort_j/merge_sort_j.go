/*
This script was created by https://gist.github.com/jayaganesh1997
*/

package merge_sort_j

func Sort(arr []int, ch chan int) {
	defer close(ch)
	if len(arr) <= 1 {
		if len(arr) == 1 {
			ch <- arr[0]
		}
		return
	}

	mid := len(arr) >> 1
	s1 := make(chan int, mid)
	s2 := make(chan int, len(arr)-mid)

	// Concurrency established
	go Sort(arr[:mid], s1)
	go Sort(arr[mid:], s2)
	// The sorting of arr[mid:] & arr[:mid] occurs Concurrently now.

	// Merging happens simultaneously and is not blocked on individual sorting.
	merge(s1, s2, ch)
}

func update(s chan int, ch chan int, c *int, ok *bool) {
	ch <- *c
	*c, *ok = <-s
}

func merge(s1, s2, ch chan int) {
	// v, ok = <-s; ok returns false if there's no more element to be received from s.
	v1, ok1 := <-s1
	v2, ok2 := <-s2
	for ok1 && ok2 {
		if v1 <= v2 {
			update(s1, ch, &v1, &ok1)
		} else {
			update(s2, ch, &v2, &ok2)
		}
	}
	for ok1 {
		update(s1, ch, &v1, &ok1)
	}
	for ok2 {
		update(s2, ch, &v2, &ok2)
	}
}
