package filter

import (
	"bufio"
	"io"

	"github.com/yourorg/logslice/internal/query"
)

// StreamFilter reads newline-delimited JSON from r and writes matching lines to w.
func StreamFilter(r io.Reader, w io.Writer, f *query.Filter) (matched int, total int, err error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		total++

		ok, matchErr := Match(line, f)
		if matchErr != nil {
			// skip malformed lines silently
			continue
		}
		if ok {
			matched++
			w.Write(line)
			w.Write([]byte("\n"))
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		err = scanErr
	}
	return
}
