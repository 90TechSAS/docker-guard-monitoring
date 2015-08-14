package utils

import (
	"strconv"
)

/*
	Convert int => string
*/
func I2S(i int) string {
	return strconv.Itoa(i)
}

/*
	Convert string => int
*/
func S2I(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

/*
	Convert float64 => string
*/
func F2S(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

/*
	Convert string => float64
*/
func S2F(s string) float64 {
	var f float64
	f, _ = strconv.ParseFloat(s, 64)
	return f
}
