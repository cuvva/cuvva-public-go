package icache

type Control struct {
	XMLNS string `xml:"xmlns,attr"`

	ExperianReference   *string
	ClientAccountNumber *string
	ClientBranchNumber  *string
	UserIdentity        *string
}
