package cher

import (
	"errors"
)

type timeoutError interface {
	Error() string
	Timeout() bool
}

func CoerceThirdPartyTimeout(err error) error {
	var netErr timeoutError
	if errors.As(err, &netErr) && netErr.Timeout() {
		return New(ThirdPartyTimeout, M{
			"error": netErr.Error(),
		})
	}

	return err
}
