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
