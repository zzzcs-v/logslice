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

## Empty Input

Passing an empty string returns `(nil, nil)`. A `nil` filter passed to
`filter.Match` matches every log entry, so callers can safely use the result
without a nil-check in the common case.
