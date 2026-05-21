package query

// Pipeline combines Parse, Validate, and Optimize into a single call.
// It returns the optimized filter or an error if parsing or validation fails.
func Pipeline(input string) (*Filter, error) {
	if input == "" {
		return nil, nil
	}

	f, err := Parse(input)
	if err != nil {
		return nil, err
	}

	if err := Validate(f); err != nil {
		return nil, err
	}

	return Optimize(f), nil
}
