package checkmot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetRecordByVRM(t *testing.T) {
	type testCase struct {
		name           string
		response       []byte
		expectedResult *Vehicle
		expectedErr    error
	}

	tests := []testCase{
		{
			name:           "no results",
			response:       []byte(`[]`),
			expectedResult: nil,
			expectedErr:    ErrNoResults,
		},
		{
			name: "two results same vehicle",
			response: []byte(`
[
  {
    "registration": "CUV 001",
    "make": "BMW",
    "model": "X3",
    "firstUsedDate": "2007.10.02",
    "fuelType": "Diesel",
    "primaryColour": "Grey",
    "motTests": [
      {
        "completedDate": "2017.08.15 12:47:50",
        "testResult": "PASSED",
        "expiryDate": "2018.08.23",
        "odometerValue": "114722",
        "odometerUnit": "mi",
        "motTestNumber": "269281492742",
        "odometerResultType": "READ",
        "rfrAndComments": [
          {
            "text": "Nearside Rear Tyre worn close to the legal limit (4.1.E.1)",
            "type": "ADVISORY",
            "dangerous": false
          }
        ]
      }
    ]
  },
  {
    "registration": "CUV 001",
    "make": "BMW",
    "model": "X3",
    "firstUsedDate": "2007.10.02",
    "fuelType": "Diesel",
    "primaryColour": "Grey",
    "motTests": [
      {
        "completedDate": "2021.09.24 12:07:40",
        "testResult": "PASSED",
        "expiryDate": "2022.09.23",
        "odometerValue": "162236",
        "odometerUnit": "mi",
        "motTestNumber": "752024886501",
        "odometerResultType": "READ",
        "rfrAndComments": [
          {
            "text": "Offside Registration plate lamp inoperative in the case of multiple lamps or light sources (4.7.1 (b) (i))",
            "type": "MINOR",
            "dangerous": false
          }
        ]
      }
    ]
  }
]
`),
			expectedResult: &Vehicle{
				Registration:    "CUV 001",
				Make:            "BMW",
				Model:           "X3",
				FuelType:        "Diesel",
				PrimaryColour:   "Grey",
				ManufactureYear: nil,
				DVLAID:          nil,
				MOTTestDueDate:  nil,
				FirstUsedDate:   datePtr("2007-10-02"),
				MOTTests: []*MOTTest{
					{

						CompletedDate:      Time{time.Date(2017, 8, 15, 12, 47, 50, 0, loc)},
						TestResult:         "PASSED",
						ExpiryDate:         datePtr("2018-08-23"),
						OdometerValue:      "114722",
						OdometerUnit:       "mi",
						MOTTestNumber:      "269281492742",
						OdometerResultType: "READ",
						RfRAndComments: []*RfROrComment{
							{
								Text:      "Nearside Rear Tyre worn close to the legal limit (4.1.E.1)",
								Type:      "ADVISORY",
								Dangerous: false,
							},
						},
					},
					{

						CompletedDate:      Time{time.Date(2021, 9, 24, 12, 7, 40, 0, loc)},
						TestResult:         "PASSED",
						ExpiryDate:         datePtr("2022-09-23"),
						OdometerValue:      "162236",
						OdometerUnit:       "mi",
						MOTTestNumber:      "752024886501",
						OdometerResultType: "READ",
						RfRAndComments: []*RfROrComment{
							{
								Text:      "Offside Registration plate lamp inoperative in the case of multiple lamps or light sources (4.7.1 (b) (i))",
								Type:      "MINOR",
								Dangerous: false,
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "multiple results different vehicles",
			response: []byte(`
[
  {
    "registration": "CUV 001",
    "make": "BMW",
    "model": "X3",
    "firstUsedDate": "2007.10.02",
    "fuelType": "Diesel",
    "primaryColour": "Grey",
    "motTests": [
      {
        "completedDate": "2017.08.15 12:47:50",
        "testResult": "PASSED",
        "expiryDate": "2018.08.23",
        "odometerValue": "114722",
        "odometerUnit": "mi",
        "motTestNumber": "269281492742",
        "odometerResultType": "READ",
        "rfrAndComments": [
          {
            "text": "Nearside Rear Tyre worn close to the legal limit (4.1.E.1)",
            "type": "ADVISORY",
            "dangerous": false
          }
        ]
      }
    ]
  },
  {
    "registration": "CUV 001",
    "make": "BMW",
    "model": "X3",
    "firstUsedDate": "2011.10.02",
    "fuelType": "Diesel",
    "primaryColour": "Grey",
    "motTests": [
      {
        "completedDate": "2021.09.24 12:07:40",
        "testResult": "PASSED",
        "expiryDate": "2022.09.23",
        "odometerValue": "162236",
        "odometerUnit": "mi",
        "motTestNumber": "752024886501",
        "odometerResultType": "READ",
        "rfrAndComments": [
          {
            "text": "Offside Registration plate lamp inoperative in the case of multiple lamps or light sources (4.7.1 (b) (i))",
            "type": "MINOR",
            "dangerous": false
          }
        ]
      }
    ]
  }
]`),
			expectedResult: nil,
			expectedErr:    ErrMultipleVehicles,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.response)
			}))

			client := NewClient(srv.URL, "client-key")
			res, err := client.GetRecordByVRM(context.Background(), "CUV 001")
			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func datePtr(str string) *Date {
	d := Date(str)
	return &d
}
