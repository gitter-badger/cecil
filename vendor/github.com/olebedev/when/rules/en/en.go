package en

import "github.com/olebedev/when/rules"

var All = []rules.Rule{
	Weekday(rules.Override),
	CasualDate(rules.Override),
	CasualTime(rules.Override),
	Hour(rules.Override),
	HourMinute(rules.Override),
	Deadline(rules.Override),
}

var WEEKDAY_OFFSET = map[string]int{
	"sunday":    0,
	"sun":       0,
	"monday":    1,
	"mon":       1,
	"tuesday":   2,
	"tue":       2,
	"wednesday": 3,
	"wed":       3,
	"thursday":  4,
	"thur":      4,
	"thu":       4,
	"friday":    5,
	"fri":       5,
	"saturday":  6,
	"sat":       6,
}

var WEEKDAY_OFFSET_PATTERN = "(?:sunday|sun|monday|mon|tuesday|tue|wednesday|wed|thursday|thur|thu|friday|fri|saturday|sat)"

// var MONTH_OFFSET = map[string]int{
// 	"january":   1,
// 	"jan":       1,
// 	"jan.":      1,
// 	"february":  2,
// 	"feb":       2,
// 	"feb.":      2,
// 	"march":     3,
// 	"mar":       3,
// 	"mar.":      3,
// 	"april":     4,
// 	"apr":       4,
// 	"apr.":      4,
// 	"may":       5,
// 	"june":      6,
// 	"jun":       6,
// 	"jun.":      6,
// 	"july":      7,
// 	"jul":       7,
// 	"jul.":      7,
// 	"august":    8,
// 	"aug":       8,
// 	"aug.":      8,
// 	"september": 9,
// 	"sep":       9,
// 	"sep.":      9,
// 	"sept":      9,
// 	"sept.":     9,
// 	"october":   10,
// 	"oct":       10,
// 	"oct.":      10,
// 	"november":  11,
// 	"nov":       11,
// 	"nov.":      11,
// 	"december":  12,
// 	"dec":       12,
// 	"dec.":      12,
// }

var INTEGER_WORDS = map[string]int{
	"one":    1,
	"two":    2,
	"three":  3,
	"four":   4,
	"five":   5,
	"six":    6,
	"seven":  7,
	"eight":  8,
	"nine":   9,
	"ten":    10,
	"eleven": 11,
	"twelve": 12,
}

var INTEGER_WORDS_PATTERN = `(?:one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve)`

// var ORDINAL_WORDS = map[string]int{
// 	"first":          1,
// 	"second":         2,
// 	"third":          3,
// 	"fourth":         4,
// 	"fifth":          5,
// 	"sixth":          6,
// 	"seventh":        7,
// 	"eighth":         8,
// 	"ninth":          9,
// 	"tenth":          10,
// 	"eleventh":       11,
// 	"twelfth":        12,
// 	"thirteenth":     13,
// 	"fourteenth":     14,
// 	"fifteenth":      15,
// 	"sixteenth":      16,
// 	"seventeenth":    17,
// 	"eighteenth":     18,
// 	"nineteenth":     19,
// 	"twentieth":      20,
// 	"twenty first":   21,
// 	"twenty second":  22,
// 	"twenty third":   23,
// 	"twenty fourth":  24,
// 	"twenty fifth":   25,
// 	"twenty sixth":   26,
// 	"twenty seventh": 27,
// 	"twenty eighth":  28,
// 	"twenty ninth":   29,
// 	"thirtieth":      30,
// 	"thirty first":   31,
// }
//
// var ORDINAL_WORDS_PATTERN = `(?:first|second|third|fourth|fifth|sixth|seventh|eighth|ninth|tenth|eleventh|twelfth|thirteenth|fourteenth|fifteenth|sixteenth|seventeenth|eighteenth|nineteenth|twentieth|twenty[ -]first|twenty[ -]second|twenty[ -]third|twenty[ -]fourth|twenty[ -]fifth|twenty[ -]sixth|twenty[ -]seventh|twenty[ -]eighth|twenty[ -]ninth|thirtieth|thirty[ -]first)`
