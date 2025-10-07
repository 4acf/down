package utils

func ClampInt(value, min, max int) int {
	ternary := func(condition bool, a, b int) int {
		if condition {
			return a
		}
		return b
	}
	return ternary(value > max, max, ternary(min > value, min, value))
}

func ClampFloat64(value, min, max float64) float64 {
	ternary := func(condition bool, a, b float64) float64 {
		if condition {
			return a
		}
		return b
	}
	return ternary(value > max, max, ternary(min > value, min, value))
}

func ScaleInt(value, maxIn, maxOut int) int {
	return int(float64(value) / float64(maxIn) * float64(maxOut))
}
