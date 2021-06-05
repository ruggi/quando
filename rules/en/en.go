package en

import (
	"regexp"
	"strconv"
	"time"

	"github.com/ruggi/quando/rules"
	"github.com/ruggi/quando/timeutil"
)

var monthNumbers = map[string]int{
	"jan": 1,
	"feb": 2,
	"mar": 3,
	"apr": 4,
	"may": 5,
	"jun": 6,
	"jul": 7,
	"aug": 8,
	"sep": 9,
	"oct": 10,
	"nov": 11,
	"dec": 12,
}

var Rules = []rules.Rule{
	{
		Name: "yesterday",
		Re:   regexp.MustCompile(`yesterday`),
		TimeFn: func(t time.Time, m map[string]string) (time.Time, error) {
			return timeutil.Today().Add(-time.Hour * 24), nil
		},
	},
	{
		Name: "tomorrow",
		Re:   regexp.MustCompile(`tomorrow`),
		TimeFn: func(t time.Time, m map[string]string) (time.Time, error) {
			return timeutil.Today().Add(time.Hour * 24), nil
		},
	},
	{
		Name: "today",
		Re:   regexp.MustCompile(`today`),
		TimeFn: func(t time.Time, m map[string]string) (time.Time, error) {
			return timeutil.Today(), nil
		},
	},
	{
		Name: "at time",
		Re:   regexp.MustCompile(`at ((?P<hours>[0-9]{1,2})(([.:,])?(?P<minutes>[0-9]{2}))? ?(?P<ampm>am|pm)?)`),
		TimeFn: func(t time.Time, m map[string]string) (time.Time, error) {
			hour, err := strconv.Atoi(m["hours"])
			if err != nil {
				return t, err
			}

			min := t.Minute()
			if mm := m["minutes"]; mm != "" {
				min, err = strconv.Atoi(mm)
				if err != nil {
					return t, err
				}
			}

			if m["ampm"] == "pm" {
				hour += 12
			}

			return time.Date(t.Year(), t.Month(), t.Day(), hour, min, t.Second(), t.Nanosecond(), t.Location()), nil
		},
	},
	{
		Name: "duration",
		Re:   regexp.MustCompile(`for (?P<duration>[0-9.]+) ?((?P<hours>h((ou)?r?s?)?)|(?P<minutes>m(in(ute)?s?)))`),
		DurationFn: func(t time.Duration, m map[string]string) (time.Duration, error) {
			d, err := strconv.ParseFloat(m["duration"], 64)
			if err != nil {
				return t, err
			}
			mul := time.Duration(0)
			if m["hours"] != "" {
				mul = time.Hour
			}
			if m["minutes"] != "" {
				mul = time.Minute
			}
			return t + time.Duration(float64(mul)*d), nil
		},
	},
	{
		Name: "date",
		Re:   regexp.MustCompile(`(on ?)?(?P<date_pre>[0-9]{1,2})? ?((?P<jan>jan(uary)?)|(?P<feb>feb(ruary)?)|(?P<mar>mar(ch)?)|(?P<apr>apr(il)?)|(?P<may>may()?)|(?P<jun>jun(e)?)|(?P<jul>jul(y)?)|(?P<aug>aug(ust)?)|(?P<sep>sep(tember)?)|(?P<oct>oct(ober)?)|(?P<nov>nov(ember)?)|(?P<dec>dec(ember)?)) (((?P<date_post>[0-9]{1,2})[^0-9]))?([, ]+(?P<year>[0-9]{4}))?`),
		TimeFn: func(t time.Time, m map[string]string) (time.Time, error) {
			var err error

			day := t.Day()
			d := ""
			if date := m["date_pre"]; date != "" {
				d = date
			}
			if date := m["date_post"]; date != "" {
				d = date
			}
			if d != "" {
				day, err = strconv.Atoi(d)
				if err != nil {
					return t, err
				}
			}
			// TODO validate 1 <= day <= 31

			month := t.Month()
			for k, v := range monthNumbers {
				if m[k] != "" {
					month = time.Month(v)
				}
			}

			year := t.Year()
			if y := m["year"]; y != "" {
				year, err = strconv.Atoi(y)
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
		Re:   regexp.MustCompile(`in +(?P<mul>[0-9.]+) +((?P<seconds>seconds?)|(?P<minutes>minutes?)|(?P<hours>hours?)|(?P<days>days?)|(?P<weeks>weeks?)|(?P<months>months?)|(?P<years>years?))`),
		TimeFn: func(t time.Time, m map[string]string) (time.Time, error) {
			mul, err := strconv.ParseFloat(m["mul"], 64)
			if err != nil {
				return t, err
			}
			units := map[string]func(time.Time) time.Time{
				"seconds": func(time.Time) time.Time { return timeutil.Now().Add(time.Duration(mul * float64(time.Second))) },
				"minutes": func(time.Time) time.Time { return timeutil.Now().Add(time.Duration(mul * float64(time.Minute))) },
				"hours":   func(time.Time) time.Time { return timeutil.Now().Add(time.Duration(mul * float64(time.Hour))) },
				"days":    func(t time.Time) time.Time { return t.AddDate(0, 0, int(mul)) },
				"weeks":   func(t time.Time) time.Time { return t.AddDate(0, 0, int(7*mul)) },
				"months":  func(t time.Time) time.Time { return t.AddDate(0, int(mul), 0) },
				"years":   func(t time.Time) time.Time { return t.AddDate(int(mul), 0, 0) },
			}
			for k, v := range units {
				if m[k] != "" {
					return v(t), nil
				}
			}
			return t, nil
		},
	},
	{
		Name: "upcoming",
		Re:   regexp.MustCompile(`((?P<next>next)|(?P<prev>prev(ious)?)|(?P<last>last)) ((?P<day>day)|(?P<week>week)|(?P<month>month)|(?P<year>year)|(?P<sun>sun(day)?)|(?P<mon>mon(day)?)|(?P<tue>tue(sday)?)|(?P<wed>wed(nesday)?)|(?P<thu>thu(rsday)?)|(?P<fri>fri(day)?)|(?P<sat>sat(urday)?))`),
		TimeFn: func(t time.Time, m map[string]string) (time.Time, error) {
			mul := 1
			if m["prev"] != "" || m["last"] != "" {
				mul = -1
			}
			findDay := func(t time.Time, target string, mul int) time.Time {
				for t.Weekday().String() != target {
					t = t.AddDate(0, 0, mul)
				}
				return t
			}
			units := map[string]func(t time.Time, mul int) time.Time{
				"day":   func(t time.Time, mul int) time.Time { return t.AddDate(0, 0, mul) },
				"week":  func(t time.Time, mul int) time.Time { return t.AddDate(0, 0, 7*mul) },
				"month": func(t time.Time, mul int) time.Time { return t.AddDate(0, mul, 0) },
				"year":  func(t time.Time, mul int) time.Time { return t.AddDate(mul, 0, 0) },
				"sun":   func(t time.Time, mul int) time.Time { return findDay(t, "Sunday", mul) },
				"mon":   func(t time.Time, mul int) time.Time { return findDay(t, "Monday", mul) },
				"tue":   func(t time.Time, mul int) time.Time { return findDay(t, "Tuesday", mul) },
				"wed":   func(t time.Time, mul int) time.Time { return findDay(t, "Wednesday", mul) },
				"thu":   func(t time.Time, mul int) time.Time { return findDay(t, "Thursday", mul) },
				"fri":   func(t time.Time, mul int) time.Time { return findDay(t, "Friday", mul) },
				"sat":   func(t time.Time, mul int) time.Time { return findDay(t, "Saturday", mul) },
			}
			for k, v := range units {
				if m[k] != "" {
					return v(t, mul), nil
				}
			}
			return t, nil
		},
	},
}
