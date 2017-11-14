package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Round is a custom implementation for rounding values as
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// JSONMarshal is a custom Marshal in order to overcome the default
// behavior of the JSON encoder
func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {

	b, err := json.Marshal(v)
	fmt.Print("bResult", b)

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

// ParseThreadURL is a thread id filter function
func ParseThreadURL(_url string) int {
	url, _ := strconv.Atoi(strings.Split(_url, "t=")[1])
	return url
}
