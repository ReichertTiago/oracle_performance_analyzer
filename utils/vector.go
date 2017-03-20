package utils

func Extend(slice []string, element string) []string {
	n := len(slice)
	if n == cap(slice) {
		newSlice := make([]string, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}


/*
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num * output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
 */