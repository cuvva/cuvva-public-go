package icache

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	urllib "net/url"
	"strings"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/soap"
	"github.com/cuvva/cuvva-public-go/lib/soap/wss"
)

var (
	ErrUnauthorized = errors.New("experian icache: unauthorized")
	ErrForbidden    = errors.New("experian icache: forbidden")
)

var path = mustParseURL("DelphiForQuotations/InteractiveWS.asmx")

func mustParseURL(url string) *urllib.URL {
	parsed, err := urllib.Parse(url)
	if err != nil {
		panic(err)
	}

	return parsed
}

type Client struct {
	httpClient *http.Client

	baseURL *urllib.URL
}

func NewClient(baseURL string) (*Client, error) {
	url, err := urllib.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL: url,

		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}, nil
}

func (c *Client) GetConsumerData(ctx context.Context, token string, input *Input) (*Output, error) {
	soapy := soap.Envelope{
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
					Input: *input,
				},
			},
		},
	}

	url := c.baseURL.ResolveReference(path).String()

	data, err := xml.Marshal(soapy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/xml")
	req.Header.Set("Content-Type", "text/xml")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result soapEnvelope
	err = handleResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	if result.Body.Fault != nil {
		return nil, cher.New("experian_soap_fault", cher.M{
			"full_soap_envelope": result,

			"code":   result.Body.Fault.Code,
			"string": result.Body.Fault.String,
		})
	}

	if result.Body.Content == nil {
		return nil, cher.New("experian_missing_response", cher.M{
			"full_soap_envelope": result,
		})
	}

	output := result.Body.Content.Root.Output

	if output.Error != nil {
		code := fmt.Sprintf("experian_icache_%s", strings.ToLower(output.Error.ErrorCode))

		return nil, cher.New(code, cher.M{
			"full_soap_envelope": result,

			"code":     output.Error.ErrorCode,
			"message":  output.Error.Message,
			"severity": output.Error.Severity,
		})
	}

	if output.OneShotFailure != nil {
		return nil, cher.New("experian_icache_failure", cher.M{
			"full_soap_envelope": result,

			"reason": output.OneShotFailure.Reason,
		})
	}

	if output.Control == nil {
		return nil, cher.New("experian_missing_response", cher.M{
			"full_soap_envelope": result,
		})
	}

	return &output, nil
}

func handleResponse(resp *http.Response, result interface{}) error {
	switch resp.StatusCode {
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("request failed (%d) - could not read body: %w", resp.StatusCode, err)
		}

		return fmt.Errorf("request failed (%d) - error body:\n\n%s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unexpected response (%d) - could not read body: %w", resp.StatusCode, err)
	}

	err = xml.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("unexpected response (%d) - could not parse xml: %w", resp.StatusCode, err)
	}

	return nil
}
