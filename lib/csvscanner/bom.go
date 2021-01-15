package csvscanner

import (
	"encoding/hex"
	"strings"
)

const BOMUTF8 string = "efbbbf"
const BOMUTF16BE string = "feff"
const BOMUTF16LE string = "fffe"
const BOMUTF32BE string = "0000feff"
const BOMUTF32LE string = "fffe0000"

func startsWithBOM(s string) bool {
	encodings := []string{
		BOMUTF8,
		BOMUTF16BE,
		BOMUTF16LE,
		BOMUTF32BE,
		BOMUTF32LE,
	}

	h := hex.EncodeToString([]byte(s))

	for _, e := range encodings {
		if strings.HasPrefix(h, e) {
			return true
		}
	}

	return false
}
