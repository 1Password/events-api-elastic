package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.1password.io/eventsapibeat/utils"
	"go.1password.io/eventsapibeat/version"
)

const (
	DefaultTimeout = 30 * time.Second
)

var DefaultUserAgent = "1Password Events API Beats / " + version.Version

type Client struct {
	httpClient *http.Client
}

type SignInAttemptResponse struct {
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
	Items   []struct {
		UUID        string    `json:"uuid"`
		SessionUUID string    `json:"session_uuid"`
		Timestamp   time.Time `json:"timestamp"`
		Country     string    `json:"country"`
		Category    string    `json:"category"`
		Type        string    `json:"type"`
		Details     *struct {
			Value string `json:"value"`
		} `json:"details"`
		SignInAttemptTargetUser struct {
			UUID  string `json:"uuid"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"target_user"`
		SignInAttemptClient struct {
			AppName         string `json:"app_name"`
			AppVersion      string `json:"app_version"`
			PlatformName    string `json:"platform_name"`
			PlatformVersion string `json:"platform_version"`
			OSName          string `json:"os_name"`
			OSVersion       string `json:"os_version"`
			IPAddress       string `json:"ip_address"`
		} `json:"client"`
	} `json:"items"`
}

type ItemUsageResponse struct {
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
	Items   []struct {
		UUID          string    `json:"uuid"`
		Timestamp     time.Time `json:"timestamp"`
		UsedVersion   uint32    `json:"used_version"`
		VaultUUID     string    `json:"vault_uuid"`
		ItemUUID      string    `json:"item_uuid"`
		ItemUsageUser struct {
			UUID  string `json:"uuid"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"user"`
		ItemUsageClient struct {
			AppName         string `json:"app_name"`
			AppVersion      string `json:"app_version"`
			PlatformName    string `json:"platform_name"`
			PlatformVersion string `json:"platform_version"`
			OSName          string `json:"os_name"`
			OSVersion       string `json:"os_version"`
			IPAddress       string `json:"ip_address"`
		} `json:"client"`
	} `json:"items"`
}

type IntrospectResponse struct {
	UUID     string    `json:"UUID"`
	IssuedAt time.Time `json:"IssuedAt"`
	Features []string  `json:"Features"`
}

func NewClient(transport http.RoundTripper) (*Client, error) {
	return &Client{
		httpClient: &http.Client{
			Timeout:   DefaultTimeout,
			Transport: transport,
		},
	}, nil
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
