package common

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func MaxFloat64(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func MinFloat64(x, y float64) float64 {
	if x > y {
		return y
	}
	return x
}
