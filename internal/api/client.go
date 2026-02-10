package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// APIClient описывает минимальный контракт для общения с backend API VoltaVPN.
type APIClient interface {
	Activate(ctx context.Context, token string) (*ActivateResponse, error)
}

// ActivateRequest — тело запроса на активацию opaque-токена.
type ActivateRequest struct {
	Token string `json:"token"`
}

// ActivateResponse — ответ сервера с сессионным токеном и VPN-профилем.
type ActivateResponse struct {
	SessionToken string `json:"session_token"`
	VPNProfile   string `json:"vpn_profile,omitempty"`
	ProfileURL   string `json:"profile_url,omitempty"`
}

// HTTPClient — реальный клиент на базе net/http.
type HTTPClient struct {
	baseURL *url.URL
	client  *http.Client
}

const (
	defaultTimeout       = 10 * time.Second
	maxResponseBodyBytes = 1 << 20 // 1 MiB
)

func NewHTTPClient(baseURL string) (*HTTPClient, error) {
	if baseURL == "" {
		return nil, errors.New("empty base URL")
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme != "https" {
		return nil, errors.New("API base URL must use HTTPS")
	}

	parsed.RawQuery = ""
	parsed.Fragment = ""

	httpClient := &http.Client{
		Timeout: defaultTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return errors.New("too many redirects")
			}

			if !sameHostHTTPS(req.URL, parsed) {
				return errors.New("redirect to unexpected host")
			}

			return nil
		},
	}

	return &HTTPClient{
		baseURL: parsed,
		client:  httpClient,
	}, nil
}

func sameHostHTTPS(u *url.URL, base *url.URL) bool {
	if u == nil || base == nil {
		return false
	}
	return u.Scheme == "https" && strings.EqualFold(u.Hostname(), base.Hostname())
}

func (c *HTTPClient) Activate(ctx context.Context, token string) (*ActivateResponse, error) {
	if c == nil || c.client == nil || c.baseURL == nil {
		return nil, errors.New("uninitialized HTTP client")
	}

	if token == "" {
		return nil, errors.New("empty token")
	}

	u := *c.baseURL
	u.Path = strings.TrimRight(u.Path, "/") + "/v1/activate"

	body, err := json.Marshal(ActivateRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("unexpected status code from API")
	}

	limited := &io.LimitedReader{
		R: resp.Body,
		N: maxResponseBodyBytes,
	}

	var out ActivateResponse
	dec := json.NewDecoder(limited)
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}

	if out.SessionToken == "" || (out.VPNProfile == "" && out.ProfileURL == "") {
		return nil, errors.New("invalid response payload")
	}

	return &out, nil
}

type MockClient struct{}

func NewClientFromEnv() (APIClient, error) {
	baseURL := strings.TrimSpace(os.Getenv("VOLTA_API_BASE_URL"))
	if baseURL == "" {
		return &MockClient{}, nil
	}

	return NewHTTPClient(baseURL)
}

func (m *MockClient) Activate(ctx context.Context, token string) (*ActivateResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("empty token")
	}

	return &ActivateResponse{
		SessionToken: "mock-session-token",
		VPNProfile:   "mock-vpn-profile",
	}, nil
}
