package structs

import (
	"time"
)

func CurrentTime() Time {
	t := time.Now().UTC()
	y , m , d := t.Date()

	return  Time  {
		Year 	: uint(y),
		Month 	:  uint(m),
		Day 	: uint(d),
		Hour 	: uint(t.Hour()),
		Minute	: uint(t.Minute()),
		Second 	: uint(t.Second()),
		}
}
func GetGoTime(t Time) time.Time {
	return time.Date(int(t.Year) , time.Month(t.Month) , int(t.Day) , int(t.Hour) , int(t.Minute) , int(t.Second) , 0 , time.UTC)
}


// CompareTime : if a < b = -1 , a == b = 0 , a > b = 1
func CompareTime(a time.Time , b time.Time) int {
	if a.Before(b) == true {
		return 1
	}
	if a.After(b) == true {
		return -1
	}
	return 0
}