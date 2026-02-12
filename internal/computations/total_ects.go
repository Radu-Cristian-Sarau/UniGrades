package computations

func TotalECTS(ects []float64) float64 {
	totalECTS := 0.0
	for _, ectsValue := range ects {
		totalECTS += ectsValue
	}
	return totalECTS
}
