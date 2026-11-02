# 概要
最近在网上看到一篇[有趣的文章](https://medium.com/@jayaganesh1997/let-us-sort-25e41a4ba854)。
说的是在归并排序中，可以将每个子数组的排序工作交给 n 个协程，再用 n 个协程将这些子数组归并来降低时间复杂度。
据说“可以将运行时间显著降为 O(N)。”

这种说法是否可行呢？

带着疑问我做了下面的实验。

# 实验
## 电脑配置
下面关于当前实现的耗时和内存，是在 Apple M2 上测的，benchmark 名称后缀为 `-8`。

[Benchmark Data](https://github.com/shabaoma/Golang-Parallel-Merge-Sort/blob/master/result/%E5%AE%9E%E9%AA%8C%E7%BB%93%E6%9E%9C.txt) 和两张图是上一版实现的记录，机器为：

* MacBook Pro (13-inch, 2016, Four Thunderbolt 3 Ports)
* 2.9 GHz Intel Core i5
* 8 GB 2133 MHz LPDDR3

上一版的 MergeSortS、MergeSortSSC 每次 merge 都 `make` 一块新切片。MergeSortSSC 把数组固定切成 4 段，使用 6 个协程。

## 算法比较
1. 快速排序（QuickSort）：标准库 `sort.Sort`，比较经 `sort.Interface` 调用
2. 归并排序（MergeSortS）：顺序归并。两块长度为 n 的缓冲轮流写入，区间长度不超过 32 时用插入排序
3. 并行归并排序（MergeSortSSC）：区间长度大于 10000，且递归深度不超过 floor(log2(GOMAXPROCS)) 时，左右两半并行
4. 并行归并排序（MergeSortSC）：每个分治节点起协程，区间长度不超过 32 时停止分裂
5. 并行归并排序（MergeSortJ）：元素经 channel 归并，分裂到长度不超过 1，协程数随 n 增长

## Benchmark
上一版的原始输出：

[Benchmark Data](https://github.com/shabaoma/Golang-Parallel-Merge-Sort/blob/master/result/%E5%AE%9E%E9%AA%8C%E7%BB%93%E6%9E%9C.txt)  
![CPU](https://raw.githubusercontent.com/shabaoma/Golang-Parallel-Merge-Sort/master/result/cpu.png)  
![MEM](https://raw.githubusercontent.com/shabaoma/Golang-Parallel-Merge-Sort/master/result/mem.png)

## 计算时间
长度不超过 10000 时，MergeSortSSC 与 MergeSortS 同级，二者快于 QuickSort，MergeSortJ 最慢。MergeSortSC 和 QuickSort 处在中间，先后随长度变化。  
长度在 10^5 及以上时，MergeSortJ > QuickSort > MergeSortS > MergeSortSC > MergeSortSSC。

## 占用内存
MergeSortJ > MergeSortSC > MergeSortSSC ≒ MergeSortS > QuickSort

# 分析
文章中的方法没有达到 O(n)。归并的总工作量仍是 O(n log n)，顶层 merge 仍是串行的 O(n)。MergeSortJ 每次分裂都创建 channel 和协程，所以它始终最慢。

下面的 pprof 取自旧版 MergeSortJ，当时长度等于 1 才返回。时间主要花在创建 channel 和发送元素：
```
      50ms       50ms      7:func Sort(arr []int, ch chan int) {
      60ms       90ms      8:	defer close(ch)
         .          .      9:	if len(arr) == 1 {
         .      1.92s     10:		ch <- arr[0]
         .      100ms     11:		return
         .          .     12:	}
         .          .     13:
         .          .     14:	mid := len(arr) >> 1
      20ms      3.41s     15:	s1 := make(chan int, mid)
         .      3.84s     16:	s2 := make(chan int, len(arr)-mid)
         .          .     17:
         .          .     18:	// Concurrency established
      10ms      800ms     19:	go Sort(arr[:mid], s1)
         .      790ms     20:	go Sort(arr[mid:], s2)
         .          .     21:	// The sorting of arr[mid:] & arr[:mid] occurs Concurrently now.
         .          .     22:
         .          .     23:	// Merging happens simultaneously and is not blocked on individual sorting.
         .        18s     24:	merge(s1, s2, ch)
         .      170ms     25:}

```

旧版 MergeSortS 和 MergeSortSSC 的内存几乎一样，分配都来自每一层 merge 的新切片，SSC 只多了 6 个协程。那一版在数组足够长时，SSC 的耗时大约是顺序归并和 `sort.Sort` 的 2/3。对应的 merge 如下：
```
      10ms       10ms     11:func merge(s1, s2 []int) []int {
         .      610ms     12:	s := make([]int, len(s1)+len(s2))
         .          .     13:	i, j, k := 0, 0, 0
     110ms      110ms     14:	for i < len(s1) && j < len(s2) {
      10ms       10ms     15:		if s1[i] < s2[j] {
     210ms      210ms     16:			s[k] = s1[i]
      80ms       80ms     17:			i++
         .          .     18:		} else {
     250ms      250ms     19:			s[k] = s2[j]
      70ms       70ms     20:			j++
         .          .     21:		}
         .          .     22:		k++
         .          .     23:	}
      10ms       10ms     24:	for i < len(s1) {
      10ms       10ms     25:		s[k] = s1[i]
         .          .     26:		i++
         .          .     27:		k++
         .          .     28:	}
         .          .     29:	for j < len(s2) {
         .          .     30:		s[k] = s2[j]
         .          .     31:		j++
         .          .     32:		k++
         .          .     33:	}
      10ms       10ms     34:	return s
         .          .     35:}
```

当前 merge 写入已有的两块缓冲，不再每次 `make`。10^7 个 int 时，MergeSortS 和 MergeSortSSC 都是大约 160 MB；分配次数分别是 2 次和十余次。长度不超过 10000 时 SSC 不开协程。长度足够大时，SSC 的耗时大约是顺序归并的 1/4，大约是 `sort.Sort` 的 1/6。

QuickSort 没有单独的快排实现。下面这段是旧工具链的 profile，时间主要在 `doPivot`。Go 1.19 起标准库改为 pdqsort，当前工具链里没有这个函数。`sort.Sort` 原地排序，这次 benchmark 把输入准备放在计时之外，计时区间里几乎没有额外分配。
```
         .       30ms    183:func quickSort(data Interface, a, b, maxDepth int) {
         .          .    184:	for b-a > 12 { // Use ShellSort for slices <= 12 elements
         .          .    185:		if maxDepth == 0 {
         .          .    186:			heapSort(data, a, b)
         .          .    187:			return
         .          .    188:		}
         .          .    189:		maxDepth--
         .      1.96s    190:		mlo, mhi := doPivot(data, a, b)
         .          .    191:		// Avoiding recursion on the larger subproblem guarantees
         .          .    192:		// a stack depth of at most lg(b-a).
         .          .    193:		if mlo-a < b-mhi {
         .      1.44s    194:			quickSort(data, a, mlo, maxDepth)
         .          .    195:			a = mhi // i.e., quickSort(data, mhi, b)
         .          .    196:		} else {
         .      1.62s    197:			quickSort(data, mhi, b, maxDepth)
         .          .    198:			b = mlo // i.e., quickSort(data, a, mlo)
         .          .    199:		}
         .          .    200:	}
         .          .    201:	if b-a > 1 {
         .          .    202:		// Do ShellSort pass with gap 6
         .          .    203:		// It could be written in this simplified form cause b-a <= 12
         .          .    204:		for i := a + 6; i < b; i++ {
         .          .    205:			if data.Less(i, i-6) {
         .          .    206:				data.Swap(i, i-6)
         .          .    207:			}
         .          .    208:		}
         .       10ms    209:		insertionSort(data, a, b)
         .          .    210:	}
         .          .    211:}
```

# 结论
文章所说的 O(n) 不成立。按元素创建协程和 channel，开销会随 n 增长。  
数组长度不超过大约 10000 时，并行分治不会更快。  
数组足够大时，只在有限深度上并行排序子区间，可以缩短时间。归并的额外内存是两块 O(n) 缓冲。  
QuickSort 的额外内存最少。在这台 Apple M2 上，当前归并快于 `sort.Sort`。

# 相关链接
[A Comparison of Parallel Sorting Algorithms on Different Architectures](https://parasol.tamu.edu/publications/download.php?file_id=191)
