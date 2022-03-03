package icache

import (
	"encoding/xml"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/stretchr/testify/assert"
)

var responseXML = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
	<soap:Body>
		<InteractiveResponse xmlns="http://www.uk.experian.com/experian/wbsv/peinteractive/v100">
			<OutputRoot xmlns="http://schemas.microsoft.com/BizTalk/2003/Any">
				<Output xmlns="http://schema.uk.experian.com/experian/cems/msgs/v1.1/ConsumerData">
					<Applicant xmlns="">
						<ApplicantIdentifier>1</ApplicantIdentifier>
						<Name>
							<Forename>Alan</Forename>
							<MiddleName/>
							<Surname>Blagg</Surname>
						</Name>
						<DateOfBirth>
							<CCYY>1970</CCYY>
							<MM>8</MM>
							<DD>19</DD>
						</DateOfBirth>
						<FormattedName>Mr A Blagg</FormattedName>
						<FormattedDOB>19/08/1970</FormattedDOB>
						<Age>
							<Years>50</Years>
							<Months>5</Months>
						</Age>
					</Applicant>
					<ConsumerSummary xmlns="">
						<PremiumValueData>
							<TPD>
								<NDHHO>N</NDHHO>
								<NDOPTOUTVALID>Y</NDOPTOUTVALID>
							</TPD>
							<CII>
								<NDSPCII>-9</NDSPCII>
								<NDSPACII>-99</NDSPACII>
							</CII>
							<Mosaic>
								<EA4M01>7876</EA4M01>
								<EA4M02>1588</EA4M02>
								<EA4M03>18</EA4M03>
								<EA4M04>100</EA4M04>
								<EA4M05>920</EA4M05>
								<EA4M06>155</EA4M06>
								<EA4M07>75</EA4M07>
								<EA4M08>0</EA4M08>
								<EA4M09>0</EA4M09>
								<EA4M10>0</EA4M10>
								<EA4M11>0</EA4M11>
								<EA4M12>0</EA4M12>
								<EA4T01>30</EA4T01>
								<EA5T01>0</EA5T01>
								<EA5T02>21</EA5T02>
								<NDG01>4</NDG01>
								<EA4N01>1190</EA4N01>
								<EA4N02>115</EA4N02>
								<EA4N03>2591</EA4N03>
								<EA4N04>-1535</EA4N04>
								<EA4N05>-2253</EA4N05>
								<NDG02>872</NDG02>
								<NDG03>0</NDG03>
								<NDG04>0</NDG04>
								<NDG05>0</NDG05>
								<NDG06>0</NDG06>
								<NDG07>0</NDG07>
								<NDG08>0</NDG08>
								<NDG09>0</NDG09>
								<NDG10>0</NDG10>
								<NDG11>0</NDG11>
								<NDG12>0</NDG12>
							</Mosaic>
							<Scoring>
								<NDSI22>DFLT6</NDSI22>
								<NDSI23>POQID</NDSI23>
								<E5S01>1</E5S01>
								<E5S02>3</E5S02>
								<E5S041>90</E5S041>
								<E5S051>-999</E5S051>
								<E5S042>0</E5S042>
								<E5S052>699</E5S052>
								<E5S043>0</E5S043>
								<E5S053>0</E5S053>
								<NDHHOSCORE>-999</NDHHOSCORE>
								<NDVALSCORE>0</NDVALSCORE>
							</Scoring>
							<AddrLink>
								<NDLNK01>N</NDLNK01>
							</AddrLink>
							<Director>
								<NDDIRSP>N</NDDIRSP>
							</Director>
							<AgeDoB>
								<EA5S01>50</EA5S01>
								<EA4S01>-1</EA4S01>
								<EA4S03>-1</EA4S03>
								<EA4S05>-1</EA4S05>
								<EA4S07>-1</EA4S07>
							</AgeDoB>
							<AdditDelphiBlocks>
								<Utilisationblock>
									<SPA01>N</SPA01>
									<SPA02>N</SPA02>
									<SPA03>N</SPA03>
									<SPA04>N</SPA04>
									<SPA05>N</SPA05>
									<SPA06>N</SPA06>
									<SPA07>N</SPA07>
									<SPA08>N</SPA08>
									<SPA09>N</SPA09>
									<SPA10>N</SPA10>
									<SPB111>N</SPB111>
									<SPB112>N</SPB112>
									<SPB113>N</SPB113>
									<SPB114>N</SPB114>
									<SPB115>N</SPB115>
									<SPB116>N</SPB116>
									<SPB117>N</SPB117>
									<SPB218>N</SPB218>
									<SPB219>N</SPB219>
									<SPB220>N</SPB220>
									<SPB221>N</SPB221>
									<SPB322>N</SPB322>
									<SPB323>N</SPB323>
									<SPC24>N</SPC24>
									<SPD25>N</SPD25>
									<SPE126>N</SPE126>
									<SPE127>N</SPE127>
									<SPE128>N</SPE128>
									<SPF129>N</SPF129>
									<SPF130>N</SPF130>
									<SPF131>N</SPF131>
									<SPF232>N</SPF232>
									<SPF233>N</SPF233>
									<SPF334>N</SPF334>
									<SPF335>N</SPF335>
									<SPF336>N</SPF336>
									<SPG37>-1</SPG37>
									<SPG38>-1</SPG38>
									<SPH39>N</SPH39>
									<SPH40>N</SPH40>
									<SPH41>N</SPH41>
									<SPCIICHECKDIGIT>N</SPCIICHECKDIGIT>
									<SPAA01>N</SPAA01>
									<SPAA02>N</SPAA02>
									<SPAA03>N</SPAA03>
									<SPAA04>N</SPAA04>
									<SPAA05>N</SPAA05>
									<SPAA06>N</SPAA06>
									<SPAA07>N</SPAA07>
									<SPAA08>N</SPAA08>
									<SPAA09>N</SPAA09>
									<SPAA10>N</SPAA10>
									<SPAB111>N</SPAB111>
									<SPAB112>N</SPAB112>
									<SPAB113>N</SPAB113>
									<SPAB114>N</SPAB114>
									<SPAB115>N</SPAB115>
									<SPAB116>N</SPAB116>
									<SPAB117>N</SPAB117>
									<SPAB218>N</SPAB218>
									<SPAB219>N</SPAB219>
									<SPAB220>N</SPAB220>
									<SPAB221>N</SPAB221>
									<SPAB322>N</SPAB322>
									<SPAB323>N</SPAB323>
									<SPAC24>N</SPAC24>
									<SPAD25>N</SPAD25>
									<SPAE126>N</SPAE126>
									<SPAE127>N</SPAE127>
									<SPAE128>N</SPAE128>
									<SPAF129>N</SPAF129>
									<SPAF130>N</SPAF130>
									<SPAF131>N</SPAF131>
									<SPAF232>N</SPAF232>
									<SPAF233>N</SPAF233>
									<SPAF334>N</SPAF334>
									<SPAF335>N</SPAF335>
									<SPAF336>N</SPAF336>
									<SPAG37>N</SPAG37>
									<SPAG38>N</SPAG38>
									<SPAH39>N</SPAH39>
									<SPAH40>N</SPAH40>
									<SPAH41>N</SPAH41>
									<SPACIICHECKDIGIT>N</SPACIICHECKDIGIT>
								</Utilisationblock>
							</AdditDelphiBlocks>
						</PremiumValueData>
						<Summary>
							<ElectoralRoll>
								<E4Q01>4</E4Q01>
								<E4Q02>4</E4Q02>
								<E4Q03>-1</E4Q03>
								<E4Q04>-1</E4Q04>
								<E4Q05>5</E4Q05>
								<E4Q06>5</E4Q06>
								<E4Q07>-1</E4Q07>
								<E4Q08>-1</E4Q08>
								<E4Q09>5</E4Q09>
								<E4Q10>5</E4Q10>
								<E4Q11>-1</E4Q11>
								<E4Q12>-1</E4Q12>
								<E4Q13>5</E4Q13>
								<E4Q14>5</E4Q14>
								<E4Q15>-1</E4Q15>
								<E4Q16>-1</E4Q16>
								<E4Q17>5</E4Q17>
								<E4R01>0</E4R01>
								<E4R02>EN118AE</E4R02>
								<EA4R01PM>0</EA4R01PM>
								<EA4R01CJ>0</EA4R01CJ>
								<EA4R01PJ>0</EA4R01PJ>
								<NDERL01>-1</NDERL01>
								<NDERL02>-1</NDERL02>
								<EA2Q02>0</EA2Q02>
								<NDERLMACA>Y</NDERLMACA>
								<EA5U01>0</EA5U01>
								<EA5U02>-1</EA5U02>
							</ElectoralRoll>
							<PublicInfo>
								<E1A01>-1</E1A01>
								<E1A02>-1</E1A02>
								<E1A03>-1</E1A03>
								<EA1C01>T</EA1C01>
								<EA1D01>-1</EA1D01>
								<EA1D02>-1</EA1D02>
								<EA1D03>-1</EA1D03>
								<E2G01>0</E2G01>
								<E2G02>0</E2G02>
								<E2G03>0</E2G03>
								<EA2J01>0</EA2J01>
								<EA2J02>0</EA2J02>
								<EA2J03>0</EA2J03>
								<EA4Q06>N</EA4Q06>
								<SPBRPRESENT>N</SPBRPRESENT>
								<SPABRPRESENT>N</SPABRPRESENT>
							</PublicInfo>
							<CAIS>
								<E1A04>0</E1A04>
								<E1A05>0</E1A05>
								<E1A06>0</E1A06>
								<E1A07>0</E1A07>
								<E1A08>0</E1A08>
								<E1A09>0</E1A09>
								<E1A10>0</E1A10>
								<E1A11>0</E1A11>
								<E1B01>-1</E1B01>
								<E1B02>0</E1B02>
								<E1B03>N</E1B03>
								<E1B04>0</E1B04>
								<E1B05>N</E1B05>
								<E1B06>0</E1B06>
								<E1B07>N</E1B07>
								<E1B08>N</E1B08>
								<E1B09>-1</E1B09>
								<E1B10>0</E1B10>
								<E1B11>0</E1B11>
								<E1B12>0</E1B12>
								<E1B13>0</E1B13>
								<NDECC01>0</NDECC01>
								<NDECC02>0</NDECC02>
								<NDECC03>-1</NDECC03>
								<NDECC04>0</NDECC04>
								<NDECC07>-1</NDECC07>
								<NDECC08>0</NDECC08>
								<E1C01>N</E1C01>
								<E1C02>0</E1C02>
								<E1C03>0</E1C03>
								<E1C04>0</E1C04>
								<E1C05>0</E1C05>
								<E1C06>0</E1C06>
								<EA1B02>0</EA1B02>
								<E1D01>0</E1D01>
								<E1D02>0</E1D02>
								<E1D03>0</E1D03>
								<E1D04>0</E1D04>
								<NDHAC01>0</NDHAC01>
								<NDHAC02>0</NDHAC02>
								<NDHAC03>0</NDHAC03>
								<NDHAC04>0</NDHAC04>
								<NDHAC05>0</NDHAC05>
								<NDHAC09>0</NDHAC09>
								<NDINC01>0</NDINC01>
								<EA1F01>N</EA1F01>
								<EA1F02>N</EA1F02>
								<EA1F03>N</EA1F03>
								<E2G04>0</E2G04>
								<E2G05>0</E2G05>
								<E2G06>0</E2G06>
								<E2G07>0</E2G07>
								<E2G08>0</E2G08>
								<E2G09>0</E2G09>
								<E2G10>0</E2G10>
								<E2G11>0</E2G11>
								<E2H01>0</E2H01>
								<E2H02>0</E2H02>
								<E2H03>N</E2H03>
								<E2H04>0</E2H04>
								<E2H05>N</E2H05>
								<E2H06>0</E2H06>
								<E2H07>N</E2H07>
								<E2H08>N</E2H08>
								<E2H09>0</E2H09>
								<E2H10>0</E2H10>
								<E2H11>0</E2H11>
								<E2H12>0</E2H12>
								<E2H13>0</E2H13>
								<NDECC05>0</NDECC05>
								<NDECC09>0</NDECC09>
								<NDECC10>0</NDECC10>
								<E2I01>N</E2I01>
								<E2I02>0</E2I02>
								<E2I03>0</E2I03>
								<E2I04>0</E2I04>
								<E2I05>0</E2I05>
								<E2I06>0</E2I06>
								<EA2H02>0</EA2H02>
								<E2J01>0</E2J01>
								<E2J02>0</E2J02>
								<E2J03>0</E2J03>
								<E2J04>0</E2J04>
								<NDHAC10>0</NDHAC10>
								<NDHAC06>0</NDHAC06>
								<NDHAC07>0</NDHAC07>
								<NDHAC08>0</NDHAC08>
								<NDINC02>0</NDINC02>
								<EA2L02>N</EA2L02>
								<EA2L03>N</EA2L03>
								<NDECC06>0</NDECC06>
								<NDINC03>0</NDINC03>
							</CAIS>
							<CAPS>
								<E1E01>0</E1E01>
								<E1E02>0</E1E02>
								<EA1B01>0</EA1B01>
								<NDPSD01>0</NDPSD01>
								<NDPSD02>0</NDPSD02>
								<NDPSD03>0</NDPSD03>
								<NDPSD04>0</NDPSD04>
								<NDPSD05>0</NDPSD05>
								<NDPSD06>0</NDPSD06>
								<EA1E01>0</EA1E01>
								<EA1E02>0</EA1E02>
								<EA1E03>0</EA1E03>
								<EA1E04>0</EA1E04>
								<E2K01>0</E2K01>
								<E2K02>0</E2K02>
								<EA2H01>0</EA2H01>
								<NDPSD07>0</NDPSD07>
								<NDPSD08>0</NDPSD08>
								<NDPSD09>0</NDPSD09>
								<NDPSD10>0</NDPSD10>
								<EA2K01>0</EA2K01>
								<EA2K02>0</EA2K02>
								<EA2K03>0</EA2K03>
								<EA2K04>0</EA2K04>
								<NDPSD11>0</NDPSD11>
							</CAPS>
							<CIFAS>
								<EA1A01>T</EA1A01>
								<EA4P01>T</EA4P01>
							</CIFAS>
							<CML>
								<EA1C02>T</EA1C02>
							</CML>
							<GAIN>
								<EA1G01>N</EA1G01>
								<EA1G02>N</EA1G02>
							</GAIN>
							<NOC>
								<EA4Q07>N</EA4Q07>
								<EA4Q08>N</EA4Q08>
								<EA4Q09>N</EA4Q09>
								<EA4Q10>N</EA4Q10>
								<EA4Q11>N</EA4Q11>
								<EA4Q01>N</EA4Q01>
								<EA4Q02>N</EA4Q02>
								<EA4Q03>N</EA4Q03>
								<EA4Q04>N</EA4Q04>
								<EA4Q05>N</EA4Q05>
							</NOC>
							<TPD>
								<NDOPTOUT>Y</NDOPTOUT>
							</TPD>
						</Summary>
					</ConsumerSummary>
					<Control xmlns="">
						<ExperianReference>6BXSSNQKN9</ExperianReference>
						<ClientAccountNumber>J6433</ClientAccountNumber>
					</Control>
					<LocationDetails xmlns="">
						<LocationIdentifier>1</LocationIdentifier>
						<UKLocation>
							<HouseNumber>4</HouseNumber>
							<Street>ADMIRALS WALK</Street>
							<PostTown>HODDESDON</PostTown>
							<County>HERTFORDSHIRE</County>
							<Postcode>EN118AE</Postcode>
						</UKLocation>
						<RMC>0130900</RMC>
						<RegionNumber>7</RegionNumber>
						<FormattedLocation>4, Admirals Walk, Hoddesdon, Hertfordshire, EN118AE</FormattedLocation>
					</LocationDetails>
					<Residency xmlns="">
						<ApplicantIdentifier>1</ApplicantIdentifier>
						<LocationIdentifier>1</LocationIdentifier>
						<LocationCode>01</LocationCode>
						<TimeAt>
							<Years>0</Years>
							<Months>8</Months>
						</TimeAt>
						<ResidencyDateFrom>
							<CCYY>2015</CCYY>
							<MM>01</MM>
							<DD>22</DD>
						</ResidencyDateFrom>
						<ResidencyDateTo>
							<CCYY>2015</CCYY>
							<MM>09</MM>
							<DD>30</DD>
						</ResidencyDateTo>
					</Residency>
					<ThirdPartyData xmlns="">
						<OptOut>Y</OptOut>
						<TransientAssocs>N</TransientAssocs>
						<HHOAllowed>N</HHOAllowed>
					</ThirdPartyData>
					<BureauMatchKey xmlns="">
						<MatchCategory>04</MatchCategory>
					</BureauMatchKey>
				</Output>
			</OutputRoot>
		</InteractiveResponse>
	</soap:Body>
</soap:Envelope>`

var responseExpected = soapEnvelope{
	XMLName: xml.Name{"http://schemas.xmlsoap.org/soap/envelope/", "Envelope"},

	Body: soapBody{
		XMLName: xml.Name{"http://schemas.xmlsoap.org/soap/envelope/", "Body"},

		Content: &InteractiveResponse{
			XMLName: xml.Name{"http://www.uk.experian.com/experian/wbsv/peinteractive/v100", "InteractiveResponse"},

			Root: OutputRoot{
				XMLName: xml.Name{"http://schemas.microsoft.com/BizTalk/2003/Any", "OutputRoot"},

				Output: Output{
					XMLName: xml.Name{"http://schema.uk.experian.com/experian/cems/msgs/v1.1/ConsumerData", "Output"},

					Control: &Control{
						ExperianReference:   ptr.String("6BXSSNQKN9"),
						ClientAccountNumber: ptr.String("J6433"),
					},

					BureauMatchKey: &BureauMatchKey{
						MatchCategory: "04",
					},

					Residency: &Residency{
						ApplicantIdentifier: 1,
						LocationIdentifier:  1,
						LocationCode:        "01",
						ResidencyDateFrom:   Date{2015, 1, 22},
						ResidencyDateTo:     Date{2015, 9, 30},
					},

					ConsumerSummary: &ConsumerSummary{
						PremiumValueData: ConsumerSummaryPremiumValueData{
							Scoring: Map{
								"NDSI22":     "DFLT6",
								"NDSI23":     "POQID",
								"E5S01":      "1",
								"E5S02":      "3",
								"E5S041":     "90",
								"E5S051":     "-999",
								"E5S042":     "0",
								"E5S052":     "699",
								"E5S043":     "0",
								"E5S053":     "0",
								"NDHHOSCORE": "-999",
								"NDVALSCORE": "0",
							},

							AddrLink: Map{
								"NDLNK01": "N",
							},

							AgeDOB: Map{
								"EA5S01": "50",
								"EA4S01": "-1",
								"EA4S03": "-1",
								"EA4S05": "-1",
								"EA4S07": "-1",
							},

							AdditDelphiBlocks: AdditDelphiBlocks{
								Utilisationblock: Map{
									"SPA01":            "N",
									"SPA02":            "N",
									"SPA03":            "N",
									"SPA04":            "N",
									"SPA05":            "N",
									"SPA06":            "N",
									"SPA07":            "N",
									"SPA08":            "N",
									"SPA09":            "N",
									"SPA10":            "N",
									"SPB111":           "N",
									"SPB112":           "N",
									"SPB113":           "N",
									"SPB114":           "N",
									"SPB115":           "N",
									"SPB116":           "N",
									"SPB117":           "N",
									"SPB218":           "N",
									"SPB219":           "N",
									"SPB220":           "N",
									"SPB221":           "N",
									"SPB322":           "N",
									"SPB323":           "N",
									"SPC24":            "N",
									"SPD25":            "N",
									"SPE126":           "N",
									"SPE127":           "N",
									"SPE128":           "N",
									"SPF129":           "N",
									"SPF130":           "N",
									"SPF131":           "N",
									"SPF232":           "N",
									"SPF233":           "N",
									"SPF334":           "N",
									"SPF335":           "N",
									"SPF336":           "N",
									"SPG37":            "-1",
									"SPG38":            "-1",
									"SPH39":            "N",
									"SPH40":            "N",
									"SPH41":            "N",
									"SPCIICHECKDIGIT":  "N",
									"SPAA01":           "N",
									"SPAA02":           "N",
									"SPAA03":           "N",
									"SPAA04":           "N",
									"SPAA05":           "N",
									"SPAA06":           "N",
									"SPAA07":           "N",
									"SPAA08":           "N",
									"SPAA09":           "N",
									"SPAA10":           "N",
									"SPAB111":          "N",
									"SPAB112":          "N",
									"SPAB113":          "N",
									"SPAB114":          "N",
									"SPAB115":          "N",
									"SPAB116":          "N",
									"SPAB117":          "N",
									"SPAB218":          "N",
									"SPAB219":          "N",
									"SPAB220":          "N",
									"SPAB221":          "N",
									"SPAB322":          "N",
									"SPAB323":          "N",
									"SPAC24":           "N",
									"SPAD25":           "N",
									"SPAE126":          "N",
									"SPAE127":          "N",
									"SPAE128":          "N",
									"SPAF129":          "N",
									"SPAF130":          "N",
									"SPAF131":          "N",
									"SPAF232":          "N",
									"SPAF233":          "N",
									"SPAF334":          "N",
									"SPAF335":          "N",
									"SPAF336":          "N",
									"SPAG37":           "N",
									"SPAG38":           "N",
									"SPAH39":           "N",
									"SPAH40":           "N",
									"SPAH41":           "N",
									"SPACIICHECKDIGIT": "N",
								},
							},

							CII: Map{
								"NDSPCII":  "-9",
								"NDSPACII": "-99",
							},

							Mosaic: Map{
								"EA4M01": "7876",
								"EA4M02": "1588",
								"EA4M03": "18",
								"EA4M04": "100",
								"EA4M05": "920",
								"EA4M06": "155",
								"EA4M07": "75",
								"EA4M08": "0",
								"EA4M09": "0",
								"EA4M10": "0",
								"EA4M11": "0",
								"EA4M12": "0",
								"EA4T01": "30",
								"EA5T01": "0",
								"EA5T02": "21",
								"NDG01":  "4",
								"EA4N01": "1190",
								"EA4N02": "115",
								"EA4N03": "2591",
								"EA4N04": "-1535",
								"EA4N05": "-2253",
								"NDG02":  "872",
								"NDG03":  "0",
								"NDG04":  "0",
								"NDG05":  "0",
								"NDG06":  "0",
								"NDG07":  "0",
								"NDG08":  "0",
								"NDG09":  "0",
								"NDG10":  "0",
								"NDG11":  "0",
								"NDG12":  "0",
							},
						},

						Summary: ConsumerSummarySummary{
							ElectoralRoll: Map{
								"E4Q01":     "4",
								"E4Q02":     "4",
								"E4Q03":     "-1",
								"E4Q04":     "-1",
								"E4Q05":     "5",
								"E4Q06":     "5",
								"E4Q07":     "-1",
								"E4Q08":     "-1",
								"E4Q09":     "5",
								"E4Q10":     "5",
								"E4Q11":     "-1",
								"E4Q12":     "-1",
								"E4Q13":     "5",
								"E4Q14":     "5",
								"E4Q15":     "-1",
								"E4Q16":     "-1",
								"E4Q17":     "5",
								"E4R01":     "0",
								"E4R02":     "EN118AE",
								"EA4R01PM":  "0",
								"EA4R01CJ":  "0",
								"EA4R01PJ":  "0",
								"NDERL01":   "-1",
								"NDERL02":   "-1",
								"EA2Q02":    "0",
								"NDERLMACA": "Y",
								"EA5U01":    "0",
								"EA5U02":    "-1",
							},

							PublicInfo: Map{
								"E1A01":        "-1",
								"E1A02":        "-1",
								"E1A03":        "-1",
								"EA1C01":       "T",
								"EA1D01":       "-1",
								"EA1D02":       "-1",
								"EA1D03":       "-1",
								"E2G01":        "0",
								"E2G02":        "0",
								"E2G03":        "0",
								"EA2J01":       "0",
								"EA2J02":       "0",
								"EA2J03":       "0",
								"EA4Q06":       "N",
								"SPBRPRESENT":  "N",
								"SPABRPRESENT": "N",
							},

							CAIS: Map{
								"E1A04":   "0",
								"E1A05":   "0",
								"E1A06":   "0",
								"E1A07":   "0",
								"E1A08":   "0",
								"E1A09":   "0",
								"E1A10":   "0",
								"E1A11":   "0",
								"E1B01":   "-1",
								"E1B02":   "0",
								"E1B03":   "N",
								"E1B04":   "0",
								"E1B05":   "N",
								"E1B06":   "0",
								"E1B08":   "N",
								"E1B07":   "N",
								"E1B09":   "-1",
								"E1B10":   "0",
								"E1B11":   "0",
								"E1B12":   "0",
								"E1B13":   "0",
								"NDECC01": "0",
								"NDECC02": "0",
								"NDECC03": "-1",
								"NDECC04": "0",
								"NDECC07": "-1",
								"NDECC08": "0",
								"E1C01":   "N",
								"E1C02":   "0",
								"E1C03":   "0",
								"E1C04":   "0",
								"E1C05":   "0",
								"E1C06":   "0",
								"EA1B02":  "0",
								"E1D01":   "0",
								"E1D02":   "0",
								"E1D03":   "0",
								"E1D04":   "0",
								"NDHAC01": "0",
								"NDHAC02": "0",
								"NDHAC03": "0",
								"NDHAC04": "0",
								"NDHAC05": "0",
								"NDHAC09": "0",
								"NDINC01": "0",
								"EA1F01":  "N",
								"EA1F02":  "N",
								"EA1F03":  "N",
								"E2G04":   "0",
								"E2G05":   "0",
								"E2G06":   "0",
								"E2G07":   "0",
								"E2G08":   "0",
								"E2G09":   "0",
								"E2G10":   "0",
								"E2G11":   "0",
								"E2H01":   "0",
								"E2H02":   "0",
								"E2H03":   "N",
								"E2H04":   "0",
								"E2H05":   "N",
								"E2H06":   "0",
								"E2H07":   "N",
								"E2H08":   "N",
								"E2H09":   "0",
								"E2H10":   "0",
								"E2H11":   "0",
								"E2H12":   "0",
								"E2H13":   "0",
								"NDECC05": "0",
								"NDECC09": "0",
								"NDECC10": "0",
								"E2I01":   "N",
								"E2I02":   "0",
								"E2I03":   "0",
								"E2I04":   "0",
								"E2I05":   "0",
								"E2I06":   "0",
								"EA2H02":  "0",
								"E2J01":   "0",
								"E2J02":   "0",
								"E2J03":   "0",
								"E2J04":   "0",
								"NDHAC10": "0",
								"NDHAC06": "0",
								"NDHAC07": "0",
								"NDHAC08": "0",
								"NDINC02": "0",
								"EA2L02":  "N",
								"EA2L03":  "N",
								"NDECC06": "0",
								"NDINC03": "0",
							},

							CAPS: Map{
								"E1E01":   "0",
								"E1E02":   "0",
								"EA1B01":  "0",
								"NDPSD01": "0",
								"NDPSD02": "0",
								"NDPSD03": "0",
								"NDPSD04": "0",
								"NDPSD05": "0",
								"NDPSD06": "0",
								"EA1E01":  "0",
								"EA1E02":  "0",
								"EA1E03":  "0",
								"EA1E04":  "0",
								"E2K01":   "0",
								"E2K02":   "0",
								"EA2H01":  "0",
								"NDPSD07": "0",
								"NDPSD08": "0",
								"NDPSD09": "0",
								"NDPSD10": "0",
								"EA2K01":  "0",
								"EA2K02":  "0",
								"EA2K03":  "0",
								"EA2K04":  "0",
								"NDPSD11": "0",
							},
						},
					},
				},
			},
		},
	},
}

func TestResponse(t *testing.T) {
	var result soapEnvelope
	err := xml.Unmarshal([]byte(responseXML), &result)
	if assert.NoError(t, err) {
		assert.Equal(t, responseExpected, result)
	}
}
