package bubble_sort

// Bubble sort has an efficiency of O(N^2)
func Sort(in []int) {
	unsortedUntilIndex := len(in)-1
	isSorted := false
	for !isSorted {
		isSorted = true
		for i := 0; i < unsortedUntilIndex; i++ {
			if in[i] > in[i+1] {
				isSorted = false
				in[i], in[i+1] = in[i+1], in[i]
			}
		}
		unsortedUntilIndex--
	}
}