package rules

import (
	"regexp"
	"time"
)

type Rule struct {
	Name string
	Re   *regexp.Regexp

	TimeFn     func(t time.Time, s []string) (time.Time, error)
	DurationFn func(t time.Duration, s []string) (time.Duration, error)

	Disabled bool
}
