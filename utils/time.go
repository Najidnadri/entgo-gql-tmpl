package utils

import "time"

func GetCurrentTimeInMalaysia() (time.Time, error) {
	//init the loc
	loc, err := time.LoadLocation("Asia/Kuala_Lumpur")
	if err != nil {
		return time.Now(), err
	}

	//set timezone,
	now := time.Now().In(loc)

	return now, nil
}

func ConvertToMalaysiaTime(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	return t.In(loc)
}
