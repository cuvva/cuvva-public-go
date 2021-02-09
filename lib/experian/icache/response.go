package icache

import (
	"encoding/xml"
)

type InteractiveResponse struct {
	XMLName xml.Name `xml:"http://www.uk.experian.com/experian/wbsv/peinteractive/v100 InteractiveResponse"`

	Root OutputRoot
}

type OutputRoot struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/BizTalk/2003/Any OutputRoot"`

	Output Output
}

type Output struct {
	XMLName xml.Name `xml:"http://schema.uk.experian.com/experian/cems/msgs/v1.1/ConsumerData Output"`

	Control *Control

	OneShotFailure *OneShotFailure
	Error          *ServiceError

	BureauMatchKey *BureauMatchKey

	Residencies []Residency `xml:"Residency"`

	ConsumerSummary *ConsumerSummary
}

type OneShotFailure struct {
	FailedLocation int
	Reason         string
}

type ServiceError struct {
	ErrorCode string
	Message   string
	Severity  int
}

type BureauMatchKey struct {
	MatchCategory string
}

type ConsumerSummary struct {
	PremiumValueData ConsumerSummaryPremiumValueData
	Summary          ConsumerSummarySummary
}

type ConsumerSummaryPremiumValueData struct {
	Scoring  Map
	AddrLink Map
	AgeDOB   Map `xml:"AgeDoB"`
}

type ConsumerSummarySummary struct {
	ElectoralRoll Map
	PublicInfo    Map
}
