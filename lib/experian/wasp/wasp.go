package wasp

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	urllib "net/url"
	"strings"
	"time"
)

var loginPath = mustParseURL("WASPAuthenticator/tokenService.asmx/LoginWithCertificate")

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

func NewClient(baseURL, certPEM, keyPEM string) (*Client, error) {
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return nil, err
	}

	url, err := urllib.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL: url,

		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			},
		},
	}, nil
}

// Login generates a token to use with Experian services.
//
// The `application` argument is an arbitrary string used to identify the
// the source of the requests for auditing purposes.
//
// The returned token is not encoded. You will likely need to base64 encode it.
func (c *Client) Login(ctx context.Context, application string) (string, error) {
	params := urllib.Values{
		"application": {application},
		"checkIP":     {"true"},
	}

	url := c.baseURL.ResolveReference(loginPath).String()
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "text/xml")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result string
	err = handleResponse(resp, &result)
	if err != nil {
		return "", err
	}

	// unfortunately there is no other way to determine the presence of an error (even with their SOAP interface!)
	if strings.Contains(strings.ToLower(result), "error") {
		return "", fmt.Errorf("error returned (%d):\n\n%s", resp.StatusCode, result)
	}

	return result, nil
}

func handleResponse(resp *http.Response, result interface{}) error {
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
