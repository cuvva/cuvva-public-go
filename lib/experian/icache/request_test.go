package icache

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/cuvva/cuvva-public-go/lib/soap"
	"github.com/cuvva/cuvva-public-go/lib/soap/wss"
	"github.com/stretchr/testify/assert"
)

var expected = `<?xml version="1.0" encoding="UTF-8"?>
<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
	<Header xmlns="http://schemas.xmlsoap.org/soap/envelope/">
		<Security xmlns="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
			<BinarySecurityToken xmlns="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" ValueType="ExperianWASP">MDAtMUEtRkY=</BinarySecurityToken>
		</Security>
	</Header>
	<Body xmlns="http://schemas.xmlsoap.org/soap/envelope/">
		<Interactive xmlns="http://www.uk.experian.com/experian/wbsv/peinteractive/v100">
			<Root xmlns="http://schemas.microsoft.com/BizTalk/2003/Any">
				<Input xmlns="http://schema.uk.experian.com/experian/cems/msgs/v1.1/ConsumerData">
					<Control xmlns=""></Control>
					<Application xmlns="">
						<ApplicationType>QI</ApplicationType>
					</Application>
					<ThirdPartyData xmlns="">
						<OptOut>Y</OptOut>
						<TransientAssocs>N</TransientAssocs>
						<HHOAllowed>N</HHOAllowed>
					</ThirdPartyData>
					<Applicant xmlns="">
						<ApplicantIdentifier>1</ApplicantIdentifier>
						<Name>
							<Forename>Alan</Forename>
							<Surname>Blagg</Surname>
						</Name>
						<DateOfBirth>
							<CCYY>1970</CCYY>
							<MM>8</MM>
							<DD>19</DD>
						</DateOfBirth>
					</Applicant>
					<LocationDetails xmlns="">
						<LocationIdentifier>1</LocationIdentifier>
						<UKLocation>
							<HouseNumber>4</HouseNumber>
							<Street>Admirals Walk</Street>
							<Postcode>EN11 8AE</Postcode>
							<Country>UK</Country>
						</UKLocation>
					</LocationDetails>
					<Residency xmlns="">
						<ApplicantIdentifier>1</ApplicantIdentifier>
						<LocationIdentifier>1</LocationIdentifier>
						<LocationCode>01</LocationCode>
						<ResidencyDateFrom>
							<CCYY>2021</CCYY>
							<MM>2</MM>
							<DD>8</DD>
						</ResidencyDateFrom>
						<ResidencyDateTo>
							<CCYY>2021</CCYY>
							<MM>2</MM>
							<DD>8</DD>
						</ResidencyDateTo>
					</Residency>
				</Input>
			</Root>
		</Interactive>
	</Body>
</Envelope>`

func TestRequest(t *testing.T) {
	token := "00-1A-FF" // actual tokens are typically 1631 chars long

	now := time.Date(2021, 2, 8, 23, 6, 0, 0, time.UTC)

	v := soap.Envelope{
		Header: soap.Header{
			Content: wss.Security{
				Token: wss.BinarySecurityToken{
					ValueType: "ExperianWASP",
					Token:     base64.StdEncoding.EncodeToString([]byte(token)),
				},
			},
		},

		Body: soap.Body{
			Content: InteractiveRequest{
				Root: InputRoot{
					Input: Input{
						Control: Control{},

						Application: Application{
							ApplicationType: "QI",
						},

						ThirdPartyData: ThirdPartyData{
							OptOut:          true,
							TransientAssocs: false,
							HHOAllowed:      false,
						},

						Applicants: []Applicant{
							{
								ApplicantIdentifier: 1,

								Name: ApplicantName{
									Forename: "Alan",
									Surname:  "Blagg",
								},

								DateOfBirth: Date{1970, 8, 19},
							},
						},

						Locations: []LocationDetails{
							{
								LocationIdentifier: 1,

								UKLocation: LocationDetailsUKLocation{
									HouseNumber: ptr.String("4"),
									Street:      ptr.String("Admirals Walk"),
									Postcode:    ptr.String("EN11 8AE"),
									Country:     ptr.String("UK"),
								},
							},
						},

						Residencies: []Residency{
							{
								ApplicantIdentifier: 1,
								LocationIdentifier:  1,

								LocationCode: "01",

								ResidencyDateFrom: NewDate(now),
								ResidencyDateTo:   NewDate(now),
							},
						},
					},
				},
			},
		},
	}

	data, err := xml.MarshalIndent(v, "", "\t")
	if assert.NoError(t, err) {
		res := fmt.Sprintf("%s%s", xml.Header, data)

		assert.Equal(t, expected, res)
	}
}
