package quando

import (
	"strings"
	"time"

	"github.com/ruggi/quando/rules"
	"github.com/ruggi/quando/rules/en"
	"github.com/ruggi/quando/timeutil"
)

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

type Option func(*Parser)

func WithRules(rules ...rules.Rule) Option {
	return func(p *Parser) {
		p.rules = rules
	}
}

var defaultOptions = []Option{
	WithRules(en.Rules...),
}

type boundary struct {
	From int
	To   int
}

type Result struct {
	Time       time.Time
	Duration   time.Duration
	Text       string
	boundaries []boundary
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

			res.boundaries = append(res.boundaries, boundary{
				From: idx[i][0],
				To:   idx[i][0] + len(sub[0]),
			})
		}
	}

	chars := []byte(s)
	for i := range chars {
		for _, b := range res.boundaries {
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
	res.Text = strings.TrimSpace(res.Text)

	return res, nil
}
