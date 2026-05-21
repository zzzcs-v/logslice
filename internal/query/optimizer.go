package query

// Optimize takes a parsed filter and applies simplification rules.
// Currently it flattens redundant AND/OR wrappers and removes nil branches.
func Optimize(f *Filter) *Filter {
	if f == nil {
		return nil
	}
	return optimizeFilter(f)
}

func optimizeFilter(f *Filter) *Filter {
	if f == nil {
		return nil
	}

	switch f.Op {
	case OpAnd:
		left := optimizeFilter(f.Left)
		right := optimizeFilter(f.Right)
		// If either side is nil, return the other
		if left == nil {
			return right
		}
		if right == nil {
			return left
		}
		return &Filter{Op: OpAnd, Left: left, Right: right}

	case OpOr:
		left := optimizeFilter(f.Left)
		right := optimizeFilter(f.Right)
		if left == nil {
			return right
		}
		if right == nil {
			return left
		}
		return &Filter{Op: OpOr, Left: left, Right: right}

	default:
		// Leaf node — return as-is
		return f
	}
}
