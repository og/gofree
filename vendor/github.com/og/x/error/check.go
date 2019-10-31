package ge

import (
	"time"
)

// Check (err error) equal `if err != nil { panic(err) }`
func Check (err error) {
	if err != nil { panic(err) }
}

func GetInt(i int, err error) int {
	Check(err)
	return i
}
func GetIntList(i []int, err error) []int {
	Check(err)
	return i
}
func GetInt32List(i []int32, err error) []int32 {
	Check(err)
	return i
}

func GetFloat64(i float64, err error) float64 {
	Check(err)
	return i
}
func GetFloat64List(i []float64, err error) []float64 {
	Check(err)
	return i
}
func GetFloat32(i float32, err error) float32 {
	Check(err)
	return i
}
func GetFloat32List(i []float32, err error) []float32 {
	Check(err)
	return i
}

func GetString(s string, err error) string {
	Check(err)
	return s
}
func GetStringList(s []string, err error) []string {
	Check(err)
	return s
}

func GetBool(b bool, err error) bool {
	Check(err)
	return b
}
func GetBoolList(b []bool, err error) []bool {
	Check(err)
	return b
}

func GetAny(v interface{}, err error) interface{} {
	Check(err)
	return v
}
func GetAnyList(v []interface{}, err error) []interface{} {
	Check(err)
	return v
}
func GetTime(v time.Time, err error) time.Time {
	Check(err)
	return v
}