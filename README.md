# logslice

Stream and filter structured JSON logs from files or stdin with a query DSL.

---

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git && cd logslice && go build -o logslice .
```

---

## Usage

```bash
# Filter logs from a file
logslice -f app.log 'level == "error"'

# Pipe from stdin
cat app.log | logslice 'status >= 500 && service == "api"'

# Select specific fields
logslice -f app.log --fields time,level,msg 'level != "debug"'

# Follow a live log file (like tail -f)
logslice -f app.log --follow 'latency_ms > 200'
```

### Query DSL

The query DSL supports basic comparisons (`==`, `!=`, `>`, `<`, `>=`, `<=`), logical operators (`&&`, `||`), and field existence checks (`has(field)`).

```bash
# Check field existence
logslice -f app.log 'has(trace_id) && level == "warn"'
```

---

## Output

By default, logslice pretty-prints matching JSON log lines. Use `--raw` to output compact JSON suitable for piping to other tools.

```bash
logslice -f app.log --raw 'level == "error"' | jq '.msg'
```

---

## Installation via Homebrew

```bash
brew install yourname/tap/logslice
```

---

## License

MIT © yourname