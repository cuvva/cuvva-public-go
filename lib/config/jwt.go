package config

// JWT configures public (and optionally private) keys and issuer for
// JSON Web Tokens. It is intended to be used in composition rather than a key.
type JWT struct {
	Issuer  string `json:"issuer"`
	Public  string `json:"public"`
	Private string `json:"private,omitempty"`
}
