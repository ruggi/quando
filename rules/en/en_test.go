package en_test

import (
	"testing"
	"time"

	"github.com/ruggi/quando"
	"github.com/ruggi/quando/rules/en"
	"github.com/ruggi/quando/timeutil"
	"github.com/stretchr/testify/require"
)

func TestRules(t *testing.T) {
	timeutil.Now = func() time.Time {
		return time.Unix(0, 0).UTC()
	}

	q := quando.NewParser(quando.WithRules(en.Rules...))

	tests := []struct {
		in       string
		wantTime string
		wantDur  string
		wantText string
	}{
		{"do stuff at 6 pm and be happy", "1970-01-01T18:00:00Z", "0s", "do stuff and be happy"},
		{"a meeting", "1970-01-01T00:00:00Z", "0s", "a meeting"},
		{"a meeting today", "1970-01-01T00:00:00Z", "0s", "a meeting"},
		{"a meeting yesterday", "1969-12-31T00:00:00Z", "0s", "a meeting"},
		{"a meeting tomorrow", "1970-01-02T00:00:00Z", "0s", "a meeting"},
		{"a meeting tomorrow at 3pm", "1970-01-02T15:00:00Z", "0s", "a meeting"},
		{"at 3pm schedule stuff", "1970-01-01T15:00:00Z", "0s", "schedule stuff"},
		{"buy flowers tomorrow at 6", "1970-01-02T06:00:00Z", "0s", "buy flowers"},
		{"buy flowers tomorrow at 6 am", "1970-01-02T06:00:00Z", "0s", "buy flowers"},
		{"block time tomorrow for 2 hours at 3pm", "1970-01-02T15:00:00Z", "2h0m0s", "block time"},
		{"go back in time on oct 26, 1985", "1985-10-26T00:00:00Z", "0s", "go back in time"},
		{"go back in time on 26 oct 1985", "1985-10-26T00:00:00Z", "0s", "go back in time"},
		{"buy flowers in 2 days", "1970-01-03T00:00:00Z", "0s", "buy flowers"},
		{"buy flowers in 2.5 minutes", "1970-01-01T00:02:30Z", "0s", "buy flowers"},
		{"buy flowers in 2 months", "1970-03-01T00:00:00Z", "0s", "buy flowers"},
		{"buy flowers in 2 years", "1972-01-01T00:00:00Z", "0s", "buy flowers"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			r, err := q.Parse(tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.wantTime, r.Time.Format(time.RFC3339))
			require.Equal(t, tt.wantDur, r.Duration.String())
			require.Equal(t, tt.wantText, r.Text)
		})
	}
}
