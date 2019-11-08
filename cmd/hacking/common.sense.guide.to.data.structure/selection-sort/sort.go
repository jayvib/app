package selection_sort

func Sort(input []int) {
	for i := 0; i < len(input); i++ {
		lowestNumberIndex := i
		for j := i+1; j < len(input); j++ {
			if input[j] < input[lowestNumberIndex] {
				lowestNumberIndex = j
			}
		}

		if lowestNumberIndex != i {
			temp := input[i]
			input[i] = input[lowestNumberIndex]
			input[lowestNumberIndex] = temp
		}
	}
}
