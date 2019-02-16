package language

var (
	timeDictionary = map[string]map[string]string{}
)

func init() {
	timeDictionary[ID] = map[string]string{
		"DayShort":    "h",
		"HourShort":   "j",
		"MinuteShort": "m",
		"Hour":        "jam",
		"Minute":      "menit",
		"Day":         "hari",
	}

	timeDictionary[EN] = map[string]string{
		"DayShort":    "d",
		"HourShort":   "h",
		"MinuteShort": "m",
		"Hour":        "hour",
		"Minute":      "minute",
		"Day":         "day",
	}

	timeDictionary[MY] = map[string]string{
		"DayShort":    "h",
		"HourShort":   "j",
		"MinuteShort": "m",
		"Hour":        "jam",
		"Minute":      "menit",
		"Day":         "hari",
	}

	timeDictionary[CN] = map[string]string{
		"DayShort":    "h",
		"HourShort":   "j",
		"MinuteShort": "m",
		"Hour":        "jam",
		"Minute":      "menit",
		"Day":         "hari",
	}
}

func getLongDay(lang, longDay string) string {
	switch lang {
	case ID:
		return longDayNamesIdID[longDay]
	case MY:
		return longDayNamesMyMY[longDay]
	case CN:
		return longDayNamesZhCN[longDay]
	}
	return longDay
}

func getShortDay(lang, shortDay string) string {
	switch lang {
	case ID:
		return shortDayNamesIdID[shortDay]
	case MY:
		return shortDayNamesMyMY[shortDay]
	case CN:
		return shortDayNamesZhCN[shortDay]
	}
	return shortDay
}

func getLongMonth(lang, longMonth string) string {
	switch lang {
	case ID:
		return longMonthNamesIdID[longMonth]
	case MY:
		return longMonthNamesMyMY[longMonth]
	case CN:
		return longMonthNamesZhCN[longMonth]
	}
	return longMonth
}

func getShortMonth(lang, shortMonth string) string {
	switch lang {
	case ID:
		return shortMonthNamesIdID[shortMonth]
	case MY:
		return shortMonthNamesMyMY[shortMonth]
	case CN:
		return shortMonthNamesZhCN[shortMonth]
	}
	return shortMonth
}

func timeTranslate(lang, key string) string {
	if _, ok := timeDictionary[lang]; ok {
		return timeDictionary[lang][key]
	}

	return ""
}
