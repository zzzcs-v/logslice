# TimeRange

The `TimeRange` type provides time-window filtering for structured log entries.

## Overview

Many log pipelines need to restrict output to a specific time window — for
example, "show me everything between 09:00 and 10:00 on 2024-06-15". The
`TimeRange` type and its helpers make this straightforward.

## Parsing

```go
tr, err := query.ParseTimeRange("2024-06-15T09:00:00Z", "2024-06-15T10:00:00Z")
```

Accepted formats:
- `RFC3339` — `2006-01-02T15:04:05Z07:00`
- `RFC3339Nano` — nanosecond precision variant
- `2006-01-02` — date-only shorthand (time treated as midnight UTC)

Either bound may be empty string to leave it open-ended.

## Checking membership

```go
if tr.Contains(entry.Time) {
    // process entry
}
```

Bounds are **inclusive** on both ends.

## Zero value

A `TimeRange` with both `From` and `To` set to `nil` is considered zero.
`IsZero()` returns `true` and `Contains()` always returns `true`, meaning
no time filtering is applied — all entries pass through.

## Integration with filter layer

`filter.MatchTimeRange(raw, tr)` inspects a raw JSON entry for any of the
well-known timestamp keys (`time`, `timestamp`, `ts`, `@timestamp`) and
delegates to `tr.Contains`. Entries with no recognisable timestamp field
are passed through unchanged.
