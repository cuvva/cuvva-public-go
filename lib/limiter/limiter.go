package limiter

// Limiter is an interface implemented by all rate limiting schemes.
type Limiter interface {
	// Allow reports whether an event with the given token can happen
	// within the configured maximum rate. Returned is the current
	// event count and whether the event can happen.
	Allow(token string) (count int64, ok bool)
}

// Tiered is a limiter where multiple Limits can be applied over
// a range of configurations for a single token.
type Tiered []Limiter

// Allow reports when an event can occur with regards to all configured
// rate limiters. The count returned will be the highest of all counts reported.
func (t Tiered) Allow(token string) (count int64, ok bool) {
	for _, l := range t {
		n, allowed := l.Allow(token)
		if !allowed {
			count = n
			return
		}

		if n > count {
			count = n
		}
	}

	ok = true
	return
}
