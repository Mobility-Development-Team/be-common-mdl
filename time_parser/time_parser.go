package time_parser

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"regexp"
	"time"
)

func ValidateDateStringFormat(date string) bool {
	// format yyyy-mm-dd
	re := regexp.MustCompile("^([0-9]{4}|[0-9]{2})[-]([0]?[1-9]|[1][0-2])[-]([0]?[1-9]|[1|2][0-9]|[3][0|1])$")
	logger.Error("re.MatchString(date):", re.MatchString(date), ",", date)
	return re.MatchString(date)
}

func ParseDateStringAsTime(date string, isDayEnd bool) (dateTime time.Time, err error) {
	t := "00:00"
	if isDayEnd {
		t = "23:59"
	}
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		return
	}
	if len(date) < 10 {
		err = errors.New("invalid date")
		return
	}
	dateTime, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", date[:10], t), loc)
	if err != nil {
		return
	}
	return
}

func ParseDateStartEndTimeFull(date, startTime, endTime string) (startDateTime, endDateTime time.Time, err error) {
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		return
	}
	if len(date) < 10 {
		err = errors.New("invalid date")
		return
	}
	startDateTime, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", date[:10], startTime), loc)
	if err != nil {
		return
	}
	endDateTime, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", date[:10], endTime), loc)
	if err != nil {
		return
	}
	return
}

func ParseDateStartEndTime(date, startTime, endTime string) (datePtr, startTimePtr, endTimePtr *time.Time, err error) {
	// Store date
	dateComponent := time.Time{}
	startTimeComponent := time.Time{}
	endTimeComponent := time.Time{}
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		logger.Error("[ParseDateStartEndTime] Failed getting timezone information, ignoring...")
	}
	if date != "" {
		dateComponent, err = time.ParseInLocation("2006-01-02", date[:10], loc)
		if err != nil {
			return
		}
		datePtr = &dateComponent
	}
	if startTime != "" {
		startTimeComponent, err = time.ParseInLocation("15:04", startTime, loc)
		if err != nil {
			return
		}
		startTimeComponent = startTimeComponent.AddDate(1970, 0, 0)
		startTimePtr = &startTimeComponent
	}
	if endTime != "" {
		endTimeComponent, err = time.ParseInLocation("15:04", endTime, loc)
		if err != nil {
			return
		}
		endTimeComponent = endTimeComponent.AddDate(1970, 0, 0)
		endTimePtr = &endTimeComponent
	}
	return
}

func GetTimeDisplay(startTime, endTime string) string {
	switch {
	case startTime != "" && endTime != "":
		return fmt.Sprintf("%s - %s", startTime, endTime)
	case startTime != "":
		return fmt.Sprintf("Starting at %s", startTime)
	case endTime != "":
		return fmt.Sprintf("Ending at %s", endTime)
	default:
		return "Full Day"
	}
}
