package mixpanel

import (
	"bufio"
	"encoding/json"
	"io"
)

type Event struct {
	Event      string
	Properties EventProperties
}

type EventProperties struct {
	DistinctID              *string      `json:"distinct_id,omitempty"`
	Time                    *int         `json:"time,omitempty"`
	OS                      *string      `json:"$os,omitempty"`
	Browser                 *string      `json:"$browser,omitempty"`
	ProjectID               *string      `json:"Project ID,omitempty"`
	Code                    *string      `json:"code,omitempty"`
	ScreenClass             *string      `json:"screen_class,omitempty"`
	ProductCategory         *string      `json:"product_category,omitempty"`
	MPDeviceModel           *string      `json:"mp_device_model,omitempty"`
	Manufacturer            *string      `json:"$manufacturer,omitempty"`
	ActiveDeviceConfigFlags *[]string    `json:"active_device_config_flags,omitempty"`
	ScreenWidth             *int         `json:"$screen_width,omitempty"`
	City                    *string      `json:"city,omitempty"`
	Label                   interface{}  `json:"label,omitempty"` // string array and strings are being send
	UnderwritersIDs         *[]string    `json:"underwriters_ids,omitempty"`
	UserID                  *string      `json:"user_id,omitempty"`
	SameAsResiAddress       *bool        `json:"same_as_resi_address,omitempty"`
	Service                 *string      `json:"service,omitempty"`
	Radio                   *string      `json:"$radio,omitempty"`
	ScreenName              *string      `json:"screen_name,omitempty"`
	AppRelease              interface{}  ` json:"$app_release,omitempty"` // Both strings and integers are being send.
	ScreenHeight            *int         `json:"$screen_height,omitempty"`
	DeviceID                *string      `json:"$device_id,omitempty"`
	MPCountryCode           *string      `json:"mp_country_code,omitempty"`
	Type                    *string      `json:"type,omitempty"`
	Product                 *string      `json:"product,omitempty"`
	OSVersion               *string      `json:"$os_version,omitempty"`
	WiFi                    *bool        `json:"$wifi,omitempty"`
	LibVersion              *string      `json:"$lib_version,omitempty"`
	Error                   *string      `json:"error,omitempty"`
	IOSIFA                  *string      `json:"$ios_ifa,omitempty"`
	ApplePay                *bool        `json:"apple-pay,omitempty"`
	Action                  *string      `json:"action,omitempty"`
	UnderwriterID           *string      `json:"underwriter-id,omitempty"`
	ProductCategories       *string      `json:"product_categories,omitempty"`
	AppBuildNumber          *interface{} `json:"$app_build_number,omitempty"` // Both strings and integers are being send.
	ActiveConfigFlags       *[]string    `json:"active_config_flags,omitempty"`
	AppVersionString        *string      `json:"$app_version_string,omitempty"`
	GUID                    *string      `json:"guid,omitempty"`
	Model                   *string      `json:"$model,omitempty"`
	Region                  *string      `json:"$region,omitempty"`
	MPLib                   *string      `json:"mp_lib,omitempty"`
	Value                   *interface{} `json:"value,omitempty"` // Both strings and integers are being send.
	AppVersion              *string      `json:"$app_version,omitempty"`
	MPProcessingTimeMs      *int         `json:"mp_processing_time_ms,omitempty"`
	Identifier              *string      `json:"identifier,omitempty"`
	ErrorCode               *string      `json:"error_code,omitempty"`
	InsertID                *string      `json:"$insert_id,omitempty"`
}

type ParseExportResultFunc func(data []byte) (event *Event, err error)

type ExportResultScanner struct {
	scanner bufio.Scanner
	parse   ParseExportResultFunc
	err     error
}

// setErr records the first error encountered.
func (ers *ExportResultScanner) setErr(err error) {
	if ers.err != nil || err == nil || err == io.EOF {
		return
	}

	ers.err = err
}

// Err returns the first non-EOF error that was encountered by the ExportResultScanner.
func (ers *ExportResultScanner) Err() error {
	return ers.err
}

func (ers *ExportResultScanner) Event() *Event {
	re, err := ers.parse(ers.scanner.Bytes())
	ers.setErr(err)

	return re
}

func (ers *ExportResultScanner) RawEventBytes() []byte {
	event := ers.Event()
	out, err := json.Marshal(event)
	ers.setErr(err)

	return out
}

func (ers *ExportResultScanner) Bytes() []byte {
	return ers.scanner.Bytes()
}

func (ers *ExportResultScanner) Scan() bool {
	return ers.scanner.Scan()
}

func ParseEvent(data []byte) (rawData *Event, err error) {
	err = json.Unmarshal(data, &rawData)
	if err != nil {
		return
	}

	return
}

func NewExportResultScanner(scanner bufio.Scanner) *ExportResultScanner {
	return &ExportResultScanner{
		scanner: scanner,
		parse:   ParseEvent,
	}
}
