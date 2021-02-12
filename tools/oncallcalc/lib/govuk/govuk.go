package govuk

import (
	"context"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type Client struct {
	*jsonclient.Client
}

func New() *Client {
	return &Client{jsonclient.NewClient("https://www.gov.uk", nil)}
}

func (g Client) GetBankHolidays(ctx context.Context) (res *BankHolidays, err error) {
	return res, g.Do(ctx, http.MethodGet, "/bank-holidays.json", nil, nil, &res)
}

type BankHolidays struct {
	EnglandAndWales Country `json:"england-and-wales"`
	Scotland        Country `json:"scotland"`
	NorthernIreland Country `json:"northern-ireland"`
}

type Country struct {
	Events []Event `json:"events"`
}

type Event struct {
	Title string `json:"title"`
	Date  string `json:"date"`
}
