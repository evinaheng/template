package language

import (
	"fmt"
	"strings"
	"time"
)

// GetDefault get default language
func GetDefault(currLang string) string {

	// Default language is Indonesian
	if currLang != ID && currLang != EN {
		return ID
	}

	return currLang
}

//GetLangID get language ID
func GetLangID(currLang string) int {
	switch currLang {
	case ID:
		return 0
	case EN:
		return 1
	case MY:
		return 2
	case CN:
		return 3
	}
	return -1
}

//GetLangTypeByID get language string
func GetLangTypeByID(id int) string {
	switch id {
	case 0:
		return ID
	case 1:
		return EN
	case 2:
		return MY
	case 3:
		return CN
	}
	return ""
}

// Date convertion from English to other language
func Date(currLang, layout string, date time.Time) string {
	if currLang == EN {
		return date.Format(layout)
	}

	formattedDate := date.Format(layout)

	englishTime := strings.Split(date.Format("Monday$Mon$January$Jan"), "$")

	r := strings.NewReplacer(englishTime[0], getLongDay(currLang, englishTime[0]), englishTime[1], getShortDay(currLang, englishTime[1]), englishTime[2], getLongMonth(currLang, englishTime[2]), englishTime[3], getShortMonth(currLang, englishTime[3]))

	formattedDate = r.Replace(formattedDate)

	return formattedDate
}

// Duration takes in the parameter : lang - language preference, value - duration in minutes, isShort - requested for abbreviation or not
func Duration(currLang string, value int, isShort bool) string {

	var minute = value % 60
	var hour = (value / 60) % 24
	var day = (value / 24) / 60

	if isShort {
		var result string
		if day > 0 {
			result += fmt.Sprintf("%d%s", day, timeTranslate(currLang, "DayShort"))
		}
		if hour > 0 {
			result += fmt.Sprintf(" %d%s", hour, timeTranslate(currLang, "HourShort"))
		}
		if minute > 0 {
			result += fmt.Sprintf(" %d%s", minute, timeTranslate(currLang, "MinuteShort"))
		}
		return strings.TrimSpace(result)
	}

	var duration string
	if minute > 0 {
		duration = fmt.Sprintf("%d %s", minute, timeTranslate(currLang, "Minute"))
		if minute > 1 && currLang == EN {
			duration += "s"
		}
	}

	if hour > 0 || day > 0 {
		hourText := fmt.Sprintf("%d %s", hour, timeTranslate(currLang, "Hour"))
		if hour > 1 && currLang == EN {
			hourText += "s"
		}

		duration = hourText + " " + duration
	}

	if day > 0 {
		dayText := fmt.Sprintf("%d %s", day, timeTranslate(currLang, "Day"))
		if day > 1 && currLang == EN {
			dayText += "s"
		}

		duration = dayText + " " + duration
	}

	return strings.TrimSpace(duration)

}
