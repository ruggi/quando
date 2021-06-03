package en

import (
	"regexp"
	"strconv"
	"time"

	"github.com/ruggi/quando/rules"
	"github.com/ruggi/quando/timeutil"
)

var Rules = []rules.Rule{
	{
		Name: "yesterday",
		Re:   regexp.MustCompile(`yesterday`),
		TimeFn: func(t time.Time, s []string) (time.Time, error) {
			return timeutil.Today().Add(-time.Hour * 24), nil
		},
	},
	{
		Name: "tomorrow",
		Re:   regexp.MustCompile(`tomorrow`),
		TimeFn: func(t time.Time, s []string) (time.Time, error) {
			return timeutil.Today().Add(time.Hour * 24), nil
		},
	},
	{
		Name: "today",
		Re:   regexp.MustCompile(`today`),
		TimeFn: func(t time.Time, s []string) (time.Time, error) {
			return timeutil.Today(), nil
		},
	},
	{
		Name: "at time",
		Re:   regexp.MustCompile(`at (([0-9]{1,2})(([.:,])?([0-9]{2}))? ?(am|pm)?)`),
		TimeFn: func(t time.Time, s []string) (time.Time, error) {
			amPM := ""
			if last := s[len(s)-1]; last == "am" || last == "pm" {
				amPM = last
			}

			var err error

			hour, err := strconv.Atoi(s[2])
			if err != nil {
				return t, err
			}

			min := t.Minute()
			if s[5] != "" {
				min, err = strconv.Atoi(s[5])
				if err != nil {
					return t, err
				}
			}

			if amPM == "pm" {
				hour += 12
			}

			return time.Date(t.Year(), t.Month(), t.Day(), hour, min, t.Second(), t.Nanosecond(), t.Location()), nil
		},
	},
	{
		Name: "duration",
		Re:   regexp.MustCompile(`for ([0-9.]+) ?((h((ou)?r?s?)?)|(m(in(ute)?s?)?))`),
		DurationFn: func(t time.Duration, s []string) (time.Duration, error) {
			d, err := strconv.ParseFloat(s[1], 64)
			if err != nil {
				return t, err
			}
			mul := time.Duration(0)
			if s[3] != "" {
				mul = time.Hour
			}
			if s[6] != "" {
				mul = time.Minute
			}
			return t + time.Duration(float64(mul)*d), nil
		},
	},
	{
		Name: "date",
		Re:   regexp.MustCompile(`(on ?)?([0-9]{1,2})? ?((jan(uary)?)|(feb(ruary)?)|(mar(ch)?)|(apr(il)?)|(may()?)|(jun(e)?)|(jul(y)?)|(aug(ust)?)|(sep(tember)?)|(oct(ober)?)|(nov(ember)?)|(dec(ember)?))( (([0-9]{1,2})[^0-9]))?([, ]+([0-9]{4}))?`),
		TimeFn: func(t time.Time, s []string) (time.Time, error) {
			var err error

			day := t.Day()
			d := ""
			if s[2] != "" {
				d = s[2]
			}
			if s[28] != "" {
				d = s[30]
			}
			if d != "" {
				day, err = strconv.Atoi(d)
				if err != nil {
					return t, err
				}
			}
			// TODO validate 1 <= day <= 31

			months := map[int]int{
				4:  1,
				6:  2,
				8:  3,
				10: 4,
				12: 5,
				14: 6,
				16: 7,
				18: 8,
				20: 9,
				22: 10,
				24: 11,
				26: 12,
			}
			month := t.Month()
			for k, v := range months {
				if s[k] != "" {
					month = time.Month(v)
				}
			}

			year := t.Year()
			if s[31] != "" {
				year, err = strconv.Atoi(s[32])
				if err != nil {
					return t, err
				}
			}

			return time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()), nil
		},
	},
	// deadlines
	{
		Name: "deadlines",
		Re:   regexp.MustCompile(`in +([0-9.]+) +((seconds?)|(minutes?)|(hours?)|(days?)|(weeks?)|(months?)|(years?))`),
		TimeFn: func(t time.Time, s []string) (time.Time, error) {
			mul, err := strconv.ParseFloat(s[1], 64)
			if err != nil {
				return t, err
			}
			units := map[int]func(time.Time) time.Time{
				// seconds
				3: func(time.Time) time.Time { return timeutil.Now().Add(time.Duration(mul * float64(time.Second))) },
				// minutes
				4: func(time.Time) time.Time {
					return timeutil.Now().Add(time.Duration(mul * float64(time.Minute)))
				},
				// hours
				5: func(time.Time) time.Time { return timeutil.Now().Add(time.Duration(mul * float64(time.Hour))) },
				// days
				6: func(t time.Time) time.Time { return t.AddDate(0, 0, int(mul)) },
				// weeks
				7: func(t time.Time) time.Time { return t.AddDate(0, 0, int(7*mul)) },
				// months
				8: func(t time.Time) time.Time { return t.AddDate(0, int(mul), 0) },
				// years
				9: func(t time.Time) time.Time { return t.AddDate(int(mul), 0, 0) },
			}
			for k, v := range units {
				if s[k] != "" {
					return v(t), nil
				}
			}
			return t, nil
		},
	},
	{
		Name: "upcoming",
		Re:   regexp.MustCompile(`((next)|(prev(ious)?)|(last)) ((day)|(week)|(month)|(year))`),
		TimeFn: func(t time.Time, s []string) (time.Time, error) {
			mul := 1
			if s[3] != "" || s[5] != "" {
				mul = -1
			}
			units := map[int]func(t time.Time, mul int) time.Time{
				7:  func(t time.Time, mul int) time.Time { return t.AddDate(0, 0, mul) },   // day
				8:  func(t time.Time, mul int) time.Time { return t.AddDate(0, 0, 7*mul) }, // week
				9:  func(t time.Time, mul int) time.Time { return t.AddDate(0, mul, 0) },   // month
				10: func(t time.Time, mul int) time.Time { return t.AddDate(mul, 0, 0) },   // year
			}
			for k, v := range units {
				if s[k] != "" {
					return v(t, mul), nil
				}
			}
			return t, nil
		},
	},
}
