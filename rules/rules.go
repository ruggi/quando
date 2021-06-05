package rules

import (
	"regexp"
	"time"
)

type Rule struct {
	Name string
	Re   *regexp.Regexp

	TimeFn     func(t time.Time, matches map[string]string) (time.Time, error)
	DurationFn func(t time.Duration, matches map[string]string) (time.Duration, error)

	Disabled bool
}
