# GroupBy

`GroupBy` partitions a slice of `json.RawMessage` log entries by the distinct
values of a named field, returning an ordered list of `GroupResult` values.

## Signature

```go
func GroupBy(entries []json.RawMessage, field string) []GroupResult
```

## Behaviour

- If `field` is empty or `entries` is nil/empty, all entries are returned in a
  single group with key `""`.
- Entries whose JSON cannot be decoded, or that do not contain the requested
  field, are placed in the `""` (empty-key) group.
- Groups are returned sorted by key in ascending lexicographic order.

## Example

```go
entries := []json.RawMessage{
    []byte(`{"level":"info","msg":"started"}`),
    []byte(`{"level":"error","msg":"oops"}`),
    []byte(`{"level":"info","msg":"done"}`),
}

groups := query.GroupBy(entries, "level")
// groups[0].Key     == "error"  (1 entry)
// groups[1].Key     == "info"   (2 entries)
```

## Integration

`GroupBy` is designed to compose with `SortEntries`, `Aggregate`, and the
`Paginator` — apply filtering and sorting within each group before rendering.
