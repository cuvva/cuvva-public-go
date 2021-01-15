package appsflyer

type Clicks struct {
	RawEventFields
}

type ClicksRetargeting struct {
	RawEventFields
}

type ConversionsRetargeting struct {
	RawEventFields
}

type Impressions struct {
	RawEventFields
}

type ImpressionsRetargeting struct {
	RawEventFields
}

type InApps struct {
	RawEventFields
}

type InAppsRetargeting struct {
	RawEventFields
}

type Installs struct {
	RawEventFields
}

type OrganicUninstalls struct {
	RawEventFields
}

type Sessions struct {
	RawEventFields
}

type Uninstalls struct {
	RawEventFields
}

type WebEvents struct {
	RawEventFields
}

type WebTouchPoints struct {
	RawEventFields
}

type WebToApp struct {
	RawEventFields
}

type RawEventFields struct {
	Ad                              string `json:"ad,omitempty"`
	AdID                            string `json:"ad_id,omitempty"`
	AdRevenueAdUnit                 string `json:"adrevenue_ad_unit,omitempty"`
	AdRevenueImpressions            string `json:"adrevenue_impressions,omitempty"`
	AdRevenueNetwork                string `json:"network,omitempty"`
	AdRevenuePlacement              string `json:"adrevenue_placement,omitempty"`
	AdRevenueSegment                string `json:"adrevenue_segment,omitempty"`
	AdSet                           string `json:"af_adset,omitempty"`
	AdSetID                         string `json:"adset_id,omitempty"`
	AdSetName                       string `json:"adset_name,omitempty"`
	AdType                          string `json:"ad_type,omitempty"`
	AdUnit                          string `json:"ad_unit,omitempty"`
	AdvertisingID                   string `json:"advertising_id,omitempty"`
	AfAd                            string `json:"af_ad,omitempty"`
	AfAdID                          string `json:"af_ad_id,omitempty"`
	AfAdSetID                       string `json:"af_adset_id,omitempty"`
	AfAdType                        string `json:"af_ad_type,omitempty"`
	AfCID                           string `json:"af_c_id,omitempty"`
	AfChannel                       string `json:"af_channel,omitempty"`
	AfCostCurrency                  string `json:"af_cost_currency,omitempty"`
	AfCostModel                     string `json:"af_cost_model,omitempty"`
	AfCostValue                     string `json:"af_cost_value,omitempty"`
	AfDeepLink                      string `json:"af_deeplink,omitempty"`
	AfKeywords                      string `json:"af_keywords,omitempty"`
	AfPartner                       string `json:"af_prt,omitempty"`
	AfReEngagementWindow            string `json:"af_reengagement_window,omitempty"`
	AfSiteID                        string `json:"af_siteid,omitempty"`
	AfSubParam1                     string `json:"af_sub1,omitempty"`
	AfSubParam2                     string `json:"af_sub2,omitempty"`
	AfSubParam3                     string `json:"af_sub3,omitempty"`
	AfSubParam4                     string `json:"af_sub4,omitempty"`
	AfSubParam5                     string `json:"af_sub5,omitempty"`
	AfSubSiteID                     string `json:"af_sub_siteid,omitempty"`
	AmazonAID                       string `json:"amazon_aid,omitempty"`
	AmazonFireID                    string `json:"amazon_fire_id,omitempty"`
	AndroidID                       string `json:"android_id,omitempty"`
	AppID                           string `json:"app_id,omitempty"`
	AppName                         string `json:"app_name,omitempty"`
	AppVersion                      string `json:"app_version,omitempty"`
	AppsFlyerID                     string `json:"appsflyer_id,omitempty"`
	AttributedTouchTime             string `json:"attributed_touch_time,omitempty"`
	AttributedTouchType             string `json:"attributed_touch_type,omitempty"`
	AttributionID                   string `json:"attribution_id,omitempty"`
	AttributionLookBack             string `json:"af_attribution_lookback,omitempty"`
	AttributionType                 string `json:"attribution_type,omitempty"`
	AttributionWindow               string `json:"attribution_window,omitempty"`
	BinnedModel                     string `json:"binned_model,omitempty"`
	BlockedReason                   string `json:"blocked_reason,omitempty"`
	BlockedReasonRule               string `json:"blocked_reason_rule,omitempty"`
	BlockedReasonValue              string `json:"blocked_reason_value,omitempty"`
	BlockedSubReason                string `json:"blocked_sub_reason,omitempty"`
	Brand                           string `json:"brand,omitempty"`
	BundleID                        string `json:"bundle_id,omitempty"`
	BundleName                      string `json:"bundle_name,omitempty"`
	Campaign                        string `json:"campaign,omitempty"`
	CampaignID                      string `json:"campaign_id,omitempty"`
	Carrier                         string `json:"carrier,omitempty"`
	Channel                         string `json:"channel,omitempty"`
	City                            string `json:"city,omitempty"`
	Contributor1AfPartner           string `json:"contributor_1_af_prt,omitempty"`
	Contributor1Campaign            string `json:"contributor_1_campaign,omitempty"`
	Contributor1MatchType           string `json:"contributor_1_match_type,omitempty"`
	Contributor1MediaSource         string `json:"contributor_1_media_source,omitempty"`
	Contributor1Partner             string `json:"contributor_1_partner,omitempty"`
	Contributor1TouchTime           string `json:"contributor_1_touch_time,omitempty"`
	Contributor1TouchType           string `json:"contributor_1_touch_type,omitempty"`
	Contributor1UnmaskedMediaSource string `json:"contributor_1_unmasked_media_source,omitempty"`
	Contributor2AfPartner           string `json:"contributor_2_af_prt,omitempty"`
	Contributor2Campaign            string `json:"contributor_2_campaign,omitempty"`
	Contributor2MatchType           string `json:"contributor_2_match_type,omitempty"`
	Contributor2MediaSource         string `json:"contributor_2_media_source,omitempty"`
	Contributor2Partner             string `json:"contributor_2_partner,omitempty"`
	Contributor2TouchTime           string `json:"contributor_2_touch_time,omitempty"`
	Contributor2TouchType           string `json:"contributor_2_touch_type,omitempty"`
	Contributor2UnmaskedMediaSource string `json:"contributor_2_unmasked_media_source,omitempty"`
	Contributor3AfPartner           string `json:"contributor_3_af_prt,omitempty"`
	Contributor3Campaign            string `json:"contributor_3_campaign,omitempty"`
	Contributor3MatchType           string `json:"contributor_3_match_type,omitempty"`
	Contributor3MediaSource         string `json:"contributor_3_media_source,omitempty"`
	Contributor3Partner             string `json:"contributor_3_partner,omitempty"`
	Contributor3TouchTime           string `json:"contributor_3_touch_time,omitempty"`
	Contributor3TouchType           string `json:"contributor_3_touch_type,omitempty"`
	Contributor3UnmaskedMediaSource string `json:"contributor_3_unmasked_media_source,omitempty"`
	CostCurrency                    string `json:"cost_currency,omitempty"`
	CostModel                       string `json:"cost_model,omitempty"`
	CostValue                       string `json:"cost_value,omitempty"`
	Country                         string `json:"country,omitempty"`
	CountryCode                     string `json:"country_code,omitempty"`
	CustomData                      string `json:"custom_data,omitempty"`
	CustomerUserID                  string `json:"customer_user_id,omitempty"`
	DMA                             string `json:"dma,omitempty"`
	DMAWithUnderscores              string `json:"_d_m_a,omitempty"`
	DeepLinkURL                     string `json:"deeplink_url,omitempty"`
	DeepLinkURLWithUnderscore       string `json:"deep_link_url,omitempty"`
	DeviceCategory                  string `json:"device_category,omitempty"`
	DeviceDownloadTime              string `json:"device_download_time,omitempty"`
	DeviceInstallTime               string `json:"device_install_time,omitempty"`
	DeviceType                      string `json:"device_type,omitempty"`
	EventID                         string `json:"event_id,omitempty"`
	EventName                       string `json:"event_name,omitempty"`
	EventRevenue                    string `json:"event_revenue,omitempty"`
	EventRevenueCurrency            string `json:"event_revenue_currency,omitempty"`
	EventRevenueUSD                 string `json:"event_revenue_usd,omitempty"`
	EventRevenueUSDWithUnderscores  string `json:"event_revenue_u_s_d,omitempty"`
	EventSource                     string `json:"event_source,omitempty"`
	EventTime                       string `json:"event_time,omitempty"`
	EventURL                        string `json:"event_url,omitempty"`
	EventValue                      string `json:"event_value,omitempty"`
	FinalData                       string `json:"final_data,omitempty"`
	GPBroadcastReferrer             string `json:"gp_broadcast_referrer,omitempty"`
	GPClickTime                     string `json:"gp_click_time,omitempty"`
	GPReferrer                      string `json:"gp_referrer,omitempty"`
	GeoRegion                       string `json:"geo_region,omitempty"`
	GeoState                        string `json:"geo_state,omitempty"`
	GooglePlayBroadcastReferrer     string `json:"google_play_broadcast_referrer,omitempty"`
	GooglePlayClickTime             string `json:"google_play_click_time,omitempty"`
	GooglePlayInstallBegin          string `json:"google_play_install_begin,omitempty"`
	GooglePlayInstallBeginTime      string `json:"gp_install_begin,omitempty"`
	GooglePlayReferrer              string `json:"google_play_referrer,omitempty"`
	HTTPReferrer                    string `json:"http_referrer,omitempty"`
	IDFA                            string `json:"idfa,omitempty"`
	IDFV                            string `json:"idfv,omitempty"`
	IMEI                            string `json:"imei,omitempty"`
	IP                              string `json:"ip,omitempty"`
	IPAddress                       string `json:"ip_address,omitempty"`
	Impressions                     string `json:"impressions,omitempty"`
	InstallAppStore                 string `json:"install_app_store,omitempty"`
	InstallID                       string `json:"install_id,omitempty"`
	InstallTime                     string `json:"install_time,omitempty"`
	InstallType                     string `json:"install_type,omitempty"`
	IsAttributed                    string `json:"is_attributed,omitempty"`
	IsOrganic                       string `json:"is_organic,omitempty"`
	IsPrimaryAttribution            string `json:"is_primary_attribution,omitempty"`
	IsPurchaseValidated             string `json:"is_purchase_validated,omitempty"`
	IsReTargeting                   string `json:"is_retargeting,omitempty"`
	IsReceiptValidated              string `json:"is_receipt_validated,omitempty"`
	IsReinstall                     string `json:"is_reinstall,omitempty"`
	JourneyID                       string `json:"journey_id,omitempty"`
	Keywords                        string `json:"keywords,omitempty"`
	KeywordsMatchType               string `json:"keyword_match_type,omitempty"`
	Language                        string `json:"language,omitempty"`
	MatchType                       string `json:"match_type,omitempty"`
	MediaChannel                    string `json:"media_channel,omitempty"`
	MediaSource                     string `json:"media_source,omitempty"`
	MediaType                       string `json:"media_type,omitempty"`
	MediationNetwork                string `json:"mediation_network,omitempty"`
	MobileCampaign                  string `json:"mobile_campaign,omitempty"`
	MobileCountry                   string `json:"mobile_country,omitempty"`
	MobileDeviceCategory            string `json:"mobile_device_category,omitempty"`
	MobileMediaSource               string `json:"mobile_media_source,omitempty"`
	MobilePlatform                  string `json:"mobile_platform,omitempty"`
	Model                           string `json:"model,omitempty"`
	ModelType                       string `json:"model_type,omitempty"`
	MonetizationNetwork             string `json:"monetization_network,omitempty"`
	NetworkAccountID                string `json:"network_account_id,omitempty"`
	OAID                            string `json:"oaid,omitempty"`
	OSVersion                       string `json:"os_version,omitempty"`
	Operator                        string `json:"operator,omitempty"`
	OriginalURL                     string `json:"original_url,omitempty"`
	PID                             string `json:"pid,omitempty"`
	ParentEventID                   string `json:"parent_event_id,omitempty"`
	Partner                         string `json:"partner,omitempty"`
	PathTouchType                   string `json:"path_touch_type,omitempty"`
	Placement                       string `json:"placement,omitempty"`
	Platform                        string `json:"platform,omitempty"`
	PostalCode                      string `json:"postal_code,omitempty"`
	ReEngagementWindow              string `json:"reengagement_window,omitempty"`
	ReceiptIDs                      string `json:"receipt_ids,omitempty"`
	Region                          string `json:"region,omitempty"`
	Resolution                      string `json:"resolution,omitempty"`
	RetargetingConversionType       string `json:"retargeting_conversion_type,omitempty"`
	SDKVersion                      string `json:"sdk_version,omitempty"`
	ScreenSize                      string `json:"screen_size,omitempty"`
	Segment                         string `json:"segment,omitempty"`
	SiteID                          string `json:"site_id,omitempty"`
	Source                          string `json:"source,omitempty"`
	State                           string `json:"state,omitempty"`
	SubParam1                       string `json:"sub_1,omitempty"`
	SubParam2                       string `json:"sub_2,omitempty"`
	SubParam3                       string `json:"sub_3,omitempty"`
	SubParam4                       string `json:"sub_4,omitempty"`
	SubParam5                       string `json:"sub_5,omitempty"`
	SubSiteID                       string `json:"sub_site_id,omitempty"`
	UTMCampaign                     string `json:"utm_campaign,omitempty"`
	UTMContent                      string `json:"utm_content,omitempty"`
	UTMID                           string `json:"utm_id,omitempty"`
	UTMMedium                       string `json:"utm_medium,omitempty"`
	UTMSource                       string `json:"utm_source,omitempty"`
	UTMTerm                         string `json:"utm_term,omitempty"`
	UnmaskedHTTPReferrer            string `json:"unmasked_http_referrer,omitempty"`
	UnmaskedIP                      string `json:"unmasked_ip,omitempty"`
	UnmaskedMediaSource             string `json:"unmasked_media_source,omitempty"`
	UnmaskedURL                     string `json:"unmasked_url,omitempty"`
	UserAgent                       string `json:"user_agent,omitempty"`
	UserDataPermission              string `json:"user_data_permission,omitempty"`
	WIFI                            string `json:"wifi,omitempty"`
	WebCampaign                     string `json:"web_campaign,omitempty"`
	WebCountry                      string `json:"web_country,omitempty"`
	WebDeviceCategory               string `json:"web_device_category,omitempty"`
	WebEventType                    string `json:"web_event_type,omitempty"`
	WebID                           string `json:"af_web_id,omitempty"`
	WebMediaChannel                 string `json:"web_media_channel,omitempty"`
	WebMediaSource                  string `json:"web_media_source,omitempty"`
	WebMediaType                    string `json:"web_media_type,omitempty"`
	WebPID                          string `json:"web_pid,omitempty"`
	WebPlatform                     string `json:"web_platform,omitempty"`
	WebReferrer                     string `json:"web_referrer,omitempty"`
	WebTimestamp                    string `json:"web_timestamp,omitempty"`
	WebUTMCampaign                  string `json:"web_utm_campaign,omitempty"`
	WebUTMContent                   string `json:"web_utm_content,omitempty"`
	WebUTMID                        string `json:"web_utm_id,omitempty"`
	WebUTMMedium                    string `json:"web_utm_medium,omitempty"`
	WebUTMSource                    string `json:"web_utm_source,omitempty"`
	WebUTMTerm                      string `json:"web_utm_term,omitempty"`
}
