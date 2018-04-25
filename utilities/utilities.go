package utilities

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	// fmt.Print("bResult", b)

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

// NumberInSlice find number in a slice
func NumberInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// InArray find number in a slice
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

// Pagination help
func Pagination(ctx *gin.Context, defaultLimit int) (limit int, offset int) {
	limit = intOr(ctx.Query("limit"), defaultLimit)
	offset, _ = strconv.Atoi(ctx.Query("offset"))
	return
}

func intOr(str string, defaultValue int) int {
	v, _ := strconv.Atoi(str)
	if v == 0 {
		return defaultValue
	}
	return v
}

// ParseTime parse post timestamp
func ParseTime(timestring string) time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)
	if len(timestring) == 0 {
		return now
	}
	datetypes := []string{"Today", "Yesterday"}
	dateString := strings.Split(timestring, ", ")[0]
	recent, _ := InArray(dateString, datetypes)
	timestamp := time.Now().In(loc)

	timeText := strings.Split(timestring, ", ")[1]
	hour, _ := strconv.Atoi(strings.Split(timeText, ":")[0])
	minute, _ := strconv.Atoi(strings.Split(timeText, ":")[1])
	if recent {
		if dateString == "Today" {
			timestamp = time.Date(now.Year(), now.Month(), now.Day(),
				hour, minute, 0, 0, loc)
		} else {
			timestamp = time.Date(now.Year(), now.Month(), now.Day()-1,
				hour, minute, 0, 0, loc)
		}
	} else {
		date, _ := strconv.Atoi(strings.Split(dateString, "-")[0])
		month, _ := strconv.Atoi(strings.Split(dateString, "-")[1])
		year, _ := strconv.Atoi(strings.Split(dateString, "-")[2])
		timestamp = time.Date(year, time.Month(month), date,
			hour, minute, 0, 0, loc)
	}

	return timestamp
}
