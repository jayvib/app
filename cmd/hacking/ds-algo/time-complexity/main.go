package main

// Fun1 demonstrate the time complexity of O(n).
func Fun1(n int) (incremented int) {
	for i := 0; i < n; i++ {
		incremented++
	}
	return
}

// Fun2 demonstrate the time complexity of O(n).
func Fun2(n int) (incremented int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			incremented++
		}
	}
	return
}

// Fun4 demonstrates the time complexity of O(n^3)
func Fun4(n int) (incremented int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				incremented++
			}
		}
	}
	return
}

