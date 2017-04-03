package utils

import (
	"strconv"
	"log"
)

func StringToFloat(numero string) float64 {
	i, err := strconv.ParseFloat(numero, 64)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func FloatToString(input_num float64, precision int) string {
	return strconv.FormatFloat(input_num, 'f', precision, 64)
}



