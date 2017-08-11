package utils

import(
    "strconv"
)

type binFloat float64
var Float *binFloat


func (bin *binFloat) String(i float64) string {
	return strconv.FormatFloat(i, 'f', -1, 64)
}






