# Query Pipeline

The `Pipeline` function provides a single entry point for turning a raw query
string into an optimized `*Filter` ready for use by the filter engine.

## Stages

1. **Parse** (`parser.go`)  
   Tokenises the input via the lexer and builds an AST (`*Filter`).

2. **Validate** (`validator.go`)  
   Ensures the AST is semantically correct (non-empty field names, known
   operators, etc.).

3. **Optimize** (`optimizer.go`)  
   Simplifies the AST by collapsing redundant `AND`/`OR` nodes whose children
   are `nil`, reducing unnecessary branching at match time.

## Usage

```go
f, err := query.Pipeline(`level = "error" AND msg exists`)
if err != nil {
    log.Fatal(err)
}
// pass f to filter.Match or filter.StreamFilter
```

## Error Handling

Errors returned by `Pipeline` are typed and can be inspected for more detail:

- `*query.ParseError` – the input could not be tokenised or parsed. The `Pos`
  field indicates the byte offset in the input string where parsing failed.
- `*query.ValidationError` – the AST was built successfully but failed semantic
  checks. The `Field` field names the offending field or operator.

Callers that want to surface a user-friendly message can use `errors.As` to
unwrap these types and extract context before formatting the error.

## Empty Input

Passing an empty string returns `(nil, nil)`. A `nil` filter passed to
`filter.Match` matches every log entry, so callers can safely use the result
without a nil-check in the common case.
