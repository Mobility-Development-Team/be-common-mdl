package timeutil

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	logger "github.com/sirupsen/logrus"
)

// ValidateDateTimeStringFormat check if the given Time value is in 'yyyy-mm-dd hh:mm:ss' format
func ValidateDateTimeStringFormat(time string) bool {
	// format hh:mm
	re := regexp.MustCompile("^([0-9]{4}|[0-9]{2})[-]([0]?[1-9]|[1][0-2])[-]([0]?[1-9]|[1|2][0-9]|[3][0|1]) ([0-1]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$")
	return re.MatchString(time)
}

// ValidateDateStringFormat check if the given Date value is in 'yyyy-mm-dd' format
func ValidateDateStringFormat(date string) bool {
	// format yyyy-mm-dd
	re := regexp.MustCompile("^([0-9]{4}|[0-9]{2})[-]([0]?[1-9]|[1][0-2])[-]([0]?[1-9]|[1|2][0-9]|[3][0|1])$")
	return re.MatchString(date)
}

// ValidateTimeStringFormat check if the given Time value is in 'hh:mm' format
func ValidateTimeStringFormat(time string) bool {
	// format hh:mm
	re := regexp.MustCompile("^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$")
	return re.MatchString(time)
}

func ParseTimeAsFormattedStr(date *time.Time, format string) (dateTime *string) {
	if date == nil {
		return nil
	}
	fd := date.Local().Format(format)
	return &fd
}

// ParseDateStringAsTime parse a single given date string in full format with either Day Start/ Day End
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

// ParseDateStringAsTimeCustomTime parse a single given date string in full format with either Day Start/ Day End
func ParseDateStringAsTimeCustomTime(dateStr, timeStr string) (dateTime time.Time, err error) {
	if !ValidateDateStringFormat(dateStr) || !ValidateTimeStringFormat(timeStr) || len(dateStr) < 10 {
		err = errors.New("invalid date, time format")
		return
	}
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		return
	}
	dateTime, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", dateStr[:10], timeStr), loc)
	if err != nil {
		return
	}
	return
}

// ParseDateStartEndTimeFull parse a single given date string with specified Day Start/ Day End time in Time
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

func ParseDateTimeStrAsDatetime(dateTimeStr string, loc *time.Location) (timePtr *time.Time, err error) {
	if nil == loc {
		loc, err = time.LoadLocation("Asia/Hong_Kong")
		if err != nil {
			logger.Error("[ParseDateStartEndTime] Failed getting timezone information, ignoring...")
		}
	}
	timeComponent := time.Time{}
	if dateTimeStr != "" {
		timeComponent, err = time.ParseInLocation("2006-01-02 15:04:05", dateTimeStr, loc)
		if err != nil {
			return
		}
		timePtr = &timeComponent
	}
	return
}

func ParseTimeStrAsDatetime(timeStr string, loc *time.Location) (timePtr *time.Time, err error) {
	if nil == loc {
		loc, err = time.LoadLocation("Asia/Hong_Kong")
		if err != nil {
			logger.Error("[ParseDateStartEndTime] Failed getting timezone information, ignoring...")
		}
	}
	timeComponent := time.Time{}
	if timeStr != "" {
		timeComponent, err = time.ParseInLocation("15:04", timeStr, loc)
		if err != nil {
			return
		}
		timeComponent = timeComponent.AddDate(1970, 0, 0)
		timePtr = &timeComponent
	}
	return
}

// ParseDateStartEndTime parse a single given date string with specified Day Start/ Day End time in Time
func ParseDateStartEndTime(date, startTime, endTime string) (datePtr, startTimePtr, endTimePtr *time.Time, err error) {
	// Store date
	dateComponent := time.Time{}
	// startTimeComponent := time.Time{}
	// endTimeComponent := time.Time{}
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
		startTimePtr, err = ParseTimeStrAsDatetime(startTime, loc)
		if err != nil {
			return
		}
		// startTimePtr = &startTimeComponent
	}
	if endTime != "" {
		endTimePtr, err = ParseTimeStrAsDatetime(endTime, loc)
		if err != nil {
			return
		}
		// endTimePtr = &endTimeComponent
	}
	return
}

// GetTimeDisplay return a time display based on the given time string
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
