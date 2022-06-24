package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.1password.io/eventsapibeat/utils"
	"go.1password.io/eventsapibeat/version"
)

var DefaultUserAgent = "1Password Events API Beats / " + version.Version

type Client struct {
	httpClient *http.Client
}

type SignInAttemptResponse struct {
	Cursor  string          `json:"cursor"`
	HasMore bool            `json:"has_more"`
	Items   []SignInAttempt `json:"items"`
}

type SignInAttempt struct {
	UUID                    string                  `json:"uuid"`
	SessionUUID             string                  `json:"session_uuid"`
	Timestamp               time.Time               `json:"timestamp"`
	Country                 string                  `json:"country"`
	Category                string                  `json:"category"`
	Type                    string                  `json:"type"`
	Details                 *SignInAttemptDetails   `json:"details"`
	SignInAttemptTargetUser SignInAttemptTargetUser `json:"target_user"`
	SignInAttemptClient     SignInAttemptClient     `json:"client"`
	SignInAttemptLocation   *SignInAttemptLocation  `json:"location"`
}

type SignInAttemptDetails struct {
	Value string `json:"value"`
}

type SignInAttemptTargetUser struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SignInAttemptClient struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	PlatformName    string `json:"platform_name"`
	PlatformVersion string `json:"platform_version"`
	OSName          string `json:"os_name"`
	OSVersion       string `json:"os_version"`
	IPAddress       string `json:"ip_address"`
}

type SignInAttemptLocation struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ItemUsageResponse struct {
	Cursor  string      `json:"cursor"`
	HasMore bool        `json:"has_more"`
	Items   []ItemUsage `json:"items"`
}

type ItemUsage struct {
	UUID              string             `json:"uuid"`
	Timestamp         time.Time          `json:"timestamp"`
	UsedVersion       uint32             `json:"used_version"`
	VaultUUID         string             `json:"vault_uuid"`
	ItemUUID          string             `json:"item_uuid"`
	Action            string             `json:"action"`
	ItemUsageUser     ItemUsageUser      `json:"user"`
	ItemUsageClient   ItemUsageClient    `json:"client"`
	ItemUsageLocation *ItemUsageLocation `json:"location"`
}

type ItemUsageUser struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ItemUsageClient struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	PlatformName    string `json:"platform_name"`
	PlatformVersion string `json:"platform_version"`
	OSName          string `json:"os_name"`
	OSVersion       string `json:"os_version"`
	IPAddress       string `json:"ip_address"`
}

type ItemUsageLocation struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type IntrospectResponse struct {
	UUID     string    `json:"UUID"`
	IssuedAt time.Time `json:"IssuedAt"`
	Features []string  `json:"Features"`
}

// RetryPolicyWithContextErrors is similar to DefaultRetryPolicy, except that
// we want to retry on context.DeadlineExceeded.
func RetryPolicyWithContextErrors(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return true, nil
		} else {
			return false, ctx.Err()
		}
	}

	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

func NewClient(logger retryablehttp.LeveledLogger, insecureSkipVerify bool) (*Client, error) {
	retryHTTPClient := retryablehttp.NewClient()
	retryHTTPClient.RetryMax = 1000000
	retryHTTPClient.CheckRetry = RetryPolicyWithContextErrors
	retryHTTPClient.Logger = logger
	if httpTransport, ok := retryHTTPClient.HTTPClient.Transport.(*http.Transport); insecureSkipVerify && ok {
		httpTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		}
	}

	client := &Client{
		httpClient: retryHTTPClient.StandardClient(),
	}

	return client, nil
}

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *Client) Introspect(ctx context.Context, bearerToken string) (*IntrospectResponse, error) {
	request, err := c.newAPIRequest(ctx, http.MethodGet, bearerToken, "/api/auth/introspect", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new API request. %w", err)
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	_ = response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %s", response.Status)
	}

	var introspectResponse IntrospectResponse
	err = json.NewDecoder(response.Body).Decode(&introspectResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response. %w", err)
	}

	return &introspectResponse, nil
}

func (c *Client) SignInAttempts(ctx context.Context, bearerToken string, cursor string) (*SignInAttemptResponse, error) {
	request, err := c.newAPIRequest(ctx, http.MethodPost, bearerToken, "/api/v1/signinattempts", strings.NewReader(cursor))
	if err != nil {
		return nil, fmt.Errorf("failed to create new API request. %w", err)
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %s", response.Status)
	}

	var signInAttemptResponse SignInAttemptResponse
	err = json.NewDecoder(response.Body).Decode(&signInAttemptResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response. %w", err)
	}

	return &signInAttemptResponse, nil
}

func (c *Client) ItemUsages(ctx context.Context, bearerToken string, cursor string) (*ItemUsageResponse, error) {
	request, err := c.newAPIRequest(ctx, http.MethodPost, bearerToken, "/api/v1/itemusages", strings.NewReader(cursor))
	if err != nil {
		return nil, fmt.Errorf("failed to create new API request. %w", err)
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %s", response.Status)
	}

	var itemUsageResponse ItemUsageResponse
	err = json.NewDecoder(response.Body).Decode(&itemUsageResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response. %w", err)
	}

	return &itemUsageResponse, nil
}

func (c *Client) newAPIRequest(ctx context.Context, method string, bearerToken string, path string, body io.Reader) (*http.Request, error) {
	jwt, err := utils.ParseJWTClaims(bearerToken)
	if err != nil {
		return nil, err
	}

	url, err := jwt.GetEventsURL()
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", url, path), body)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	request.Header.Add("User-Agent", DefaultUserAgent)
	return request, nil
}
