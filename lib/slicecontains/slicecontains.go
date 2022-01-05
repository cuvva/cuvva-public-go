package slicecontains

// Elem checks if an element is contained within a slice
func Elem[E comparable](slice []E, val E) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

// SameElems checks that all the elements exist in both a and b.
func SameElems[E comparable](a, b []E) bool {
	if len(a) != len(b) {
		return false
	}

	keys := map[E]struct{}{}
	for _, v := range a {
		keys[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := keys[v]; !ok {
			return false
		}
	}

	return true
}

// Deprecated: Use Elem instead
func String(slice []string, val string) bool {
	return Elem(slice, val)
}

// Deprecated: Use SameElems instead
// SameStrings checks that all of the strings exist in a and b.
func SameStrings(a, b []string) bool {
	return SameElems(a,b)
}
