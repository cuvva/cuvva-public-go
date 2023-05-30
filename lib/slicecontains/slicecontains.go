package slicecontains

func Int(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

func String(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

// SameStrings checks that all of the strings exist in a and b.
func SameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	keys := map[string]struct{}{}
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
