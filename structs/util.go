package structs

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

// StrToDate should be in format yyyy/mm/dd 
func StrToDate(s string) (*Date , error){
	var (
		y, 
		m,
		d string
	)
	_ , err := fmt.Sscanf(s , "%s-%s-%s\n" , &y , &m , &d)
	if err != nil {
		return nil , err
	}
	
	year , err := strconv.Atoi(y)
	if err != nil {
		return nil , err
	}
	month , err := strconv.Atoi(m)
	if err != nil {
		return nil , err
	}
	day , err := strconv.Atoi(d)
	if err != nil {
		return nil , err
	}
	
	return &Date{
		Year : uint(year),
		Month : uint(month),
		Day : uint(day),
	} , nil
}


func GetAge(dob string) (int , error){
	var (y , m , d string)
	for i := 0;i < len(dob);i++{
		if i < 4 {
			y += string(dob[i])
		} else if i < 7 && i > 4 {
			m += string(dob[i])
		} else if i < 10 && i > 7 {
			d += string(dob[i])
		}
	}
	year , err := strconv.Atoi(y)
	month , err := strconv.Atoi(m)
	day , err := strconv.Atoi(d)
	if err != nil {
		return 0, err
	}
	bDate := Date{
		Year : uint(year),
		Month : uint(month),
		Day : uint(day),
	}
	
	aDate := CurrentTime()
	
	ageYears := int(aDate.Year - bDate.Year)
	ageMonths := int(aDate.Month - bDate.Month)
	ageDays := int(aDate.Day - bDate.Day)
	
	
	age := 0
	
	if ageMonths < 0 && ageDays < 0 {
		age = ageYears - 1
	} else {
		age = ageYears
	}
	
	return age , nil
}



func GenerateID() string {
	return uuid.NewString()
}

func GenerateUserId() string {
	id := GenerateID()
	id = "user-" + id
	return id
}

func GenerateChatId() string {
	id := GenerateID()
	id = "chat-" + id
	return id
}

func GenerateFileId() string {
	id := GenerateID()
	id = "file-" + id
	return id
}
