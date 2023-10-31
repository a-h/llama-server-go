package llamaservergo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	Client     *http.Client
	Middleware []Middleware
}

type Middleware interface {
	Request(req *http.Request) error
	Repsonse(res *http.Response) error
}

func WithRequestHeader(key, value string) Middleware {
	return &requestHeaderMiddleware{key: key, value: value}
}

type requestHeaderMiddleware struct {
	key   string
	value string
}

func (m *requestHeaderMiddleware) Request(req *http.Request) error {
	req.Header.Set(m.key, m.value)
	return nil
}

func (m *requestHeaderMiddleware) Repsonse(res *http.Response) error {
	return nil
}

func WithAuthorization(authorization string) Middleware {
	return WithRequestHeader("Authorization", authorization)
}

func WithContentType(contentType string) Middleware {
	return WithRequestHeader("Content-Type", contentType)
}

func New(baseURL string, timeout time.Duration, middleware ...Middleware) (Client, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return Client{}, fmt.Errorf("could not parse base URL: %w", err)
	}
	return Client{
		BaseURL:    parsedBaseURL,
		Client:     &http.Client{Timeout: timeout},
		Middleware: middleware,
	}, nil
}

func Post[TReq, TResp any](ctx context.Context, client http.Client, middleware []Middleware, url string, request TReq) (response TResp, err error) {
	buf, err := json.Marshal(request)
	if err != nil {
		return response, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}
	for _, m := range middleware {
		if err := m.Request(req); err != nil {
			return response, fmt.Errorf("middleware failed to modify request: %w", err)
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	for _, m := range middleware {
		if err := m.Repsonse(res); err != nil {
			return response, fmt.Errorf("middleware failed to modify response: %w", err)
		}
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return response, fmt.Errorf("api responded with non-success status %d: message: %s", res.StatusCode, string(body))
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("api responded with 2xx status code, but the response could not be decoded: %w", err)
	}
	return response, nil
}

func Get[TResp any](ctx context.Context, client http.Client, middleware []Middleware, url string) (response TResp, ok bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response, false, fmt.Errorf("failed to create request: %w", err)
	}
	for _, m := range middleware {
		if err := m.Request(req); err != nil {
			return response, false, fmt.Errorf("middleware failed to modify request: %w", err)
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return response, false, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer res.Body.Close()
	for _, m := range middleware {
		if err := m.Repsonse(res); err != nil {
			return response, false, fmt.Errorf("middleware failed to modify response: %w", err)
		}
	}
	if res.StatusCode == http.StatusNotFound {
		return response, false, nil
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return response, false, fmt.Errorf("api responded with non-success status %d: message: %s", res.StatusCode, string(body))
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return response, false, fmt.Errorf("api responded with 2xx status code, but the response could not be decoded: %w", err)
	}
	return response, true, nil
}

func (client *Client) PropsGet(ctx context.Context) (response PropsGetResponse, ok bool, err error) {
	u := *client.BaseURL
	u.Path = path.Join(u.Path, "props")
	return Get[PropsGetResponse](ctx, *client.Client, client.Middleware, u.String())
}

func (client *Client) CompletionPost(ctx context.Context, request CompletionPostRequest) (response CompletionPostResponse, err error) {
	u := *client.BaseURL
	u.Path = path.Join(u.Path, "completion")
	return Post[CompletionPostRequest, CompletionPostResponse](ctx, *client.Client, client.Middleware, u.String(), request)
}

func (client *Client) EmbeddingPost(ctx context.Context, request EmbeddingPostRequest) (response EmbeddingPostResponse, err error) {
	u := *client.BaseURL
	u.Path = path.Join(u.Path, "embedding")
	return Post[EmbeddingPostRequest, EmbeddingPostResponse](ctx, *client.Client, client.Middleware, u.String(), request)
}

func (client *Client) TokenizePost(ctx context.Context, request TokenizePostRequest) (response TokenizePostResponse, err error) {
	u := *client.BaseURL
	u.Path = path.Join(u.Path, "tokenize")
	return Post[TokenizePostRequest, TokenizePostResponse](ctx, *client.Client, client.Middleware, u.String(), request)
}
