package convert

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

const jakartaTZ = "Asia/Jakarta"

var jakartaLocation, _ = time.LoadLocation(jakartaTZ)

// ToFloat64 Convert any value to float64
func ToFloat64(v interface{}) float64 {
	result := float64(0)
	switch v.(type) {
	case string:
		result, _ = strconv.ParseFloat(v.(string), 64)
	case int:
		result = float64(v.(int))
	case int64:
		result = float64(v.(int64))
	case float64:
		result = float64(v.(float64))
	case uint8:
		result, _ = strconv.ParseFloat(string(v.(uint8)), 64)
	default:
		result = float64(0)
	}

	return result
}

// ToByteArr Convert any value to byte array
func ToByteArr(v interface{}) []byte {
	result := []byte("")
	if v == nil {
		return result
	}

	switch v.(type) {
	case string:
		result = []byte(v.(string))
	case int:
		result = []byte(strconv.Itoa(v.(int)))
	case int64:
		result = []byte(strconv.FormatInt(v.(int64), 10))
	case bool:
		result = []byte(strconv.FormatBool(v.(bool)))
	case float64:
		result = []byte(strconv.FormatFloat(v.(float64), 'E', -1, 64))
	case []uint8:
		result = v.([]uint8)
	default:
		resultJSON, err := json.Marshal(v)
		if err == nil {
			result = resultJSON
		} else {
			log.Println("func ToByteArr", err)
		}
	}

	return result
}

// ToInt Convert any value to int
func ToInt(v interface{}) int {
	result := int(0)
	switch v.(type) {
	case string:
		result, _ = strconv.Atoi(v.(string))
	case int:
		result = int(v.(int))
	case int64:
		result = int(v.(int64))
	case float64:
		result = int(v.(float64))
	case uint8:
		result, _ = strconv.Atoi(string(v.(uint8)))
	default:
		result = int(0)
	}

	return result
}

// ToInt64 Convert any value to int64
func ToInt64(v interface{}) int64 {
	result := int64(0)
	switch v.(type) {
	case string:
		result, _ = strconv.ParseInt(v.(string), 10, 64)
	case int:
		result = int64(v.(int))
	case int64:
		result = int64(v.(int64))
	case float64:
		result = int64(v.(float64))
	case uint8:
		result, _ = strconv.ParseInt(string(v.(uint8)), 10, 64)
	default:
		result = int64(0)
	}

	return result
}

// ToString Convert any value to string
func ToString(v interface{}) string {
	result := ""
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		result = v.(string)
	case int:
		result = strconv.Itoa(v.(int))
	case int64:
		result = strconv.FormatInt(v.(int64), 10)
	case bool:
		result = strconv.FormatBool(v.(bool))
	case float64:
		result = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case uint8:
		result = string(v.(uint8))
	case []uint8:
		result = string(v.([]uint8))
	default:
		resultJSON, err := json.Marshal(v)
		if err == nil {
			result = string(resultJSON)
		} else {
			log.Println("func ToString", err)
		}
	}

	return result
}

// FixMySQLTime Convert time to Jakarta/WIB (UTC+7) with offset
// If input time is UTC+0 and using offset 25200 (UTC+7) then output is time in UTC+7 zone without value change
// Else if input time is UTC+0 and using offset for example 28800 (UTC+8)
// then output is time in UTC+7 zone with original input time minus one hour
// E.g.
// input: t = 2018-10-30 00:00:00 +0000 UTC, offset = 25200(JKT). output: 2018-10-30 00:00:00 +0700 WIB
// input: t = 2018-10-30 00:00:00 +0000 UTC, offset = 28800(SG). output: 2018-10-29 23:00:00 +0700 WIB
// input: t = 2018-10-30 00:00:00 +0000 UTC, offset = -14400(EDT DST). output: 2018-10-30 11:00:00 +0700 WIB
func FixMySQLTime(t time.Time, offset int) time.Time {
	return t.In(jakartaLocation).Add(time.Duration(-offset) * time.Second)
}
