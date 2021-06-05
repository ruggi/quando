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
		indexes := r.Re.FindStringIndex(s)
		if indexes == nil {
			continue
		}
		matches := r.Re.FindStringSubmatch(s)
		names := make(map[string]string, r.Re.NumSubexp())
		for _, n := range r.Re.SubexpNames() {
			i := r.Re.SubexpIndex(n)
			if i < 0 {
				continue
			}
			names[n] = matches[i]
		}
		var err error
		if r.DurationFn != nil {
			res.Duration, err = r.DurationFn(res.Duration, names)
		} else {
			res.Time, err = r.TimeFn(res.Time, names)
		}
		if err != nil {
			return Result{}, err
		}
		res.Boundaries = append(res.Boundaries, Boundary{
			From: indexes[0],
			To:   indexes[1],
		})
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
