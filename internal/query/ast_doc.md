# Query AST

This package defines the Abstract Syntax Tree (AST) types used by the logslice query DSL.

## Types

### `Filter`

Represents a single filter predicate with an optional chain via `And` or `Or`.

| Field   | Type       | Description                              |
|---------|------------|------------------------------------------|
| `Field` | `string`   | The JSON log field to match against      |
| `Op`    | `Operator` | The comparison or existence operator     |
| `Value` | `string`   | The value to compare (empty for exists)  |
| `And`   | `*Filter`  | Optional chained AND condition           |
| `Or`    | `*Filter`  | Optional chained OR condition            |

### `Operator`

Supported operators:

| Operator | Symbol   | Description              |
|----------|----------|--------------------------|
| `OpEq`   | `=`      | Exact equality           |
| `OpNeq`  | `!=`     | Not equal                |
| `OpContains` | `~`  | Substring match          |
| `OpExists`   | `exists` | Field presence check |
| `OpGt`   | `>`      | Greater than (numeric)   |
| `OpLt`   | `<`      | Less than (numeric)      |
| `OpGte`  | `>=`     | Greater than or equal    |
| `OpLte`  | `<=`     | Less than or equal       |

## Example

```
level = error AND service ~ auth
```

Parses into a `Filter{Field:"level", Op:OpEq, Value:"error", And: &Filter{Field:"service", Op:OpContains, Value:"auth"}}`.
