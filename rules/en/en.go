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

			hour := t.Hour()
			hour, err = strconv.Atoi(s[2])
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
		Re:   regexp.MustCompile(`for ([0-9.]+) ?((h(ours?)?)|(m(inutes?)?))`),
		DurationFn: func(t time.Duration, s []string) (time.Duration, error) {
			d, err := strconv.ParseFloat(s[1], 64)
			if err != nil {
				return t, err
			}
			mul := time.Duration(0)
			if s[3] != "" {
				mul = time.Hour
			}
			if s[5] != "" {
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
}
