package benchmark

func HeavyJob(a int) int {
	result := 0
	for i := 1; i <= a; i++ {
		result += i
	}

	return result
}
