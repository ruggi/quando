package quando

import (
	"regexp"
	"strings"
	"time"

	"github.com/ruggi/quando/rules"
	"github.com/ruggi/quando/rules/en"
	"github.com/ruggi/quando/timeutil"
)

var (
	defaultOptions = []Option{
		WithRules(en.Rules...),
	}
	reMultipleSpaces = regexp.MustCompile(`[\s]+`)
)

type Option func(*Parser)

func WithRules(rules ...rules.Rule) Option {
	return func(p *Parser) {
		p.rules = rules
	}
}

type Boundary struct {
	From int
	To   int
}

type Result struct {
	// The parsed time
	Time time.Time
	// The parsed duration
	Duration time.Duration
	// The text resulting from removing the matching time/date tokens from the original string
	Text string
	// The boundaries of the matching tokens
	Boundaries []Boundary
}

type Parser struct {
	rules []rules.Rule
}

func NewParser(options ...Option) *Parser {
	p := &Parser{}
	for _, opt := range defaultOptions {
		opt(p)
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}

func (p *Parser) Parse(s string) (Result, error) {
	res := Result{
		Time: timeutil.Today(),
	}

	for _, r := range p.rules {
		if r.Disabled {
			continue
		}
		idx := r.Re.FindAllStringSubmatchIndex(s, -1)
		if len(idx) <= 0 {
			continue
		}
		submatches := r.Re.FindAllStringSubmatch(strings.ToLower(s), -1)
		for i, sub := range submatches {
			var err error
			if r.DurationFn != nil {
				res.Duration, err = r.DurationFn(res.Duration, sub)
			} else {
				res.Time, err = r.TimeFn(res.Time, sub)
			}
			if err != nil {
				return Result{}, err
			}

			res.Boundaries = append(res.Boundaries, Boundary{
				From: idx[i][0],
				To:   idx[i][0] + len(sub[0]),
			})
		}
	}

	chars := []byte(s)
	for i := range chars {
		for _, b := range res.Boundaries {
			if i >= b.From && i < b.To {
				chars[i] = 0
			}
		}
	}
	for _, c := range chars {
		if c == 0 {
			continue
		}
		res.Text += string(c)
	}
	res.Text = strings.TrimSpace(reMultipleSpaces.ReplaceAllString(res.Text, " "))

	return res, nil
}
