package computations

func Average(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0.0
	for _, num := range nums {
		sum += num
	}
	return sum / float64(len(nums))
}

func WeightedAverage(nums []float64, weights []float64) float64 {
	if len(nums) == 0 || len(weights) == 0 || len(nums) != len(weights) {
		return 0
	}
	sum := 0.0
	weightSum := 0.0
	for i, num := range nums {
		sum += num * weights[i]
		weightSum += weights[i]
	}
	if weightSum == 0 {
		return 0
	}
	return sum / weightSum
}

func AverageGradePerYear(grades []float64, years []int) map[int]float64 {
	yearGrades := make(map[int][]float64)
	for i, grade := range grades {
		year := years[i]
		yearGrades[year] = append(yearGrades[year], grade)
	}
	avgPerYear := make(map[int]float64)
	for year, grades := range yearGrades {
		avgPerYear[year] = Average(grades)
	}
	return avgPerYear
}
