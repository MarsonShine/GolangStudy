package demo3

import "math"

func GetPrimes(max int) []int {
	if max <= 1 {
		return []int{}
	}
	marks := make([]bool, max)
	var count int
	squareRoot := int(math.Sqrt(float64(max)))
	for i := 2; i <= squareRoot; i++ {
		if !marks[i] {
			for j := i * i; j < max; j += i {
				if !marks[j] {
					marks[j] = true
					count++
				}
			}
		}
	}
	primes := make([]int, 0, max-count)
	for i := 2; i < max; i++ {
		if !marks[i] {
			primes = append(primes, i)
		}
	}
	return primes
}
