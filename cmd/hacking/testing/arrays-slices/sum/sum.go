package sum

func Sum(numbers []int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

func SumAll(numbersToSum ...[]int) []int {

	var numTotals []int
	for _, numberGroup := range numbersToSum {
		numTotals = append(numTotals, Sum(numberGroup))
	}

	return numTotals
}

func SumAllTails(numbersToSum ...[]int) []int {
	var res []int
	for _, numberGroup := range numbersToSum {
		if len(numberGroup) == 0 {
			res = append(res, 0)
			continue
		}
		res = append(res, Sum(numberGroup[1:]))
	}
	return res
}
