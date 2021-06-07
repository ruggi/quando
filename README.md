# Quando

Parse time and date with natural language.

The `parser.Parse` function returns a `Result`:

```go
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
```

## Usage

```go
package main

import (
    "github.com/ruggi/quando"
)

func main() {
    p := quando.NewParser()
    res, err := p.Parse("buy groceries tomorrow at 2pm")
    if err != nil {
        // ...
    }

    fmt.Println(res.Time) // tomorrow at 2pm
    fmt.Println(res.Text) // "buy groceries"
}
```

## Some examples

| Input                                           | Time                | Duration | Text                   |
| ----------------------------------------------- | ------------------- | -------- | ---------------------- |
| `buy groceries`                                 | today (midnight)    | 0        | `buy groceries`        |
| `buy groceries in 5 minutes`                    | now + 5 minutes     | 0        | `buy groceries`        |
| `gym for 3 hours in 2 weeks`                    | 2 weeks from now    | 3h       | `gym`                  |
| `send christmas cards on dec 23, 2050 at 16:30` | 2050/12/23 16:30:00 | 0        | `send christmas cards` |

## Credits

This started inspired by [when](https://github.com/olebedev/when), but with a slightly different take (e.g. support for durations, timezones, boundaries, etc.).
