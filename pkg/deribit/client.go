package deribit

import (
	"bytes"
	"cs-projects-eth-collar/internal/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	apiKey         string
	apiSecret      string
	baseURL        string
	httpClient     *http.Client
	accessToken    string
	tokenExpiresAt time.Time
	authMutex      sync.RWMutex
}

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewClient(config types.DeribitConfig) *Client {
	var baseURL string

	// 测试网
	if config.TestNet {
		baseURL = "https://test.deribit.com/api/v2"
	} else {
		baseURL = "https://www.deribit.com/api/v2"
	}

	return &Client{
		apiKey:    config.APIKey,
		apiSecret: config.APISecret,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Authenticate() error {
	c.authMutex.Lock()
	defer c.authMutex.Unlock()

	// Deribit认证API使用JSON-RPC格式
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "public/auth",
		"params": map[string]interface{}{
			"grant_type":    "client_credentials",
			"client_id":     c.apiKey,
			"client_secret": c.apiSecret,
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal auth request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth HTTP error %d: %s", resp.StatusCode, string(responseBody))
	}

	var authResponse struct {
		Result struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int64  `json:"expires_in"`
		} `json:"result"`
		Error *APIError `json:"error"`
	}

	if err := json.Unmarshal(responseBody, &authResponse); err != nil {
		return fmt.Errorf("failed to unmarshal auth response: %w", err)
	}

	if authResponse.Error != nil {
		return fmt.Errorf("authentication API error: %s (code: %d)", authResponse.Error.Message, authResponse.Error.Code)
	}

	c.accessToken = authResponse.Result.AccessToken
	// 设置过期时间，提前1分钟过期以避免边界情况
	expiresIn := authResponse.Result.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600 // 默认1小时
	}
	c.tokenExpiresAt = time.Now().Add(time.Duration(expiresIn-60) * time.Second)

	return nil
}

func (c *Client) isTokenValid() bool {
	// 检查当前token是否有效
	c.authMutex.RLock()
	defer c.authMutex.RUnlock()

	return c.accessToken != "" && time.Now().Before(c.tokenExpiresAt)
}

func (c *Client) ensureAuthenticated() error {
	// 如果未认证或token过期则重新认证
	if c.isTokenValid() {
		return nil
	}

	return c.Authenticate()
}

func (c *Client) GetAccountSummary(currency string) (*types.AccountSummary, error) {
	endpoint := "/private/get_account_summary"
	params := map[string]interface{}{
		"currency": currency,
	}

	var response struct {
		Result types.AccountSummary `json:"result"`
		Error  *APIError            `json:"error"`
	}

	if err := c.makePrivateRequest("GET", endpoint, params, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("API error: %s (code: %d)", response.Error.Message, response.Error.Code)
	}

	return &response.Result, nil
}

func (c *Client) GetPositions(currency string) ([]interface{}, error) {
	endpoint := "/private/get_positions"
	params := map[string]interface{}{
		"currency": currency,
		"kind":     "future",
	}

	var response struct {
		Result []interface{} `json:"result"`
		Error  *APIError     `json:"error"`
	}

	if err := c.makePrivateRequest("GET", endpoint, params, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("API error: %s (code: %d)", response.Error.Message, response.Error.Code)
	}

	return response.Result, nil
}

func (c *Client) GetIndexPrice(currency string) (float64, error) {
	// 获取指数价格 (现货价格)
	endpoint := "/public/get_index_price"
	params := map[string]interface{}{
		"index_name": currency + "_usd", // 例如: eth_usd
	}

	var response struct {
		Result struct {
			IndexPrice float64 `json:"index_price"`
		} `json:"result"`
		Error *APIError `json:"error"`
	}

	if err := c.makePublicRequest("GET", endpoint, params, &response); err != nil {
		return 0, err
	}

	if response.Error != nil {
		return 0, fmt.Errorf("API error: %s (code: %d)", response.Error.Message, response.Error.Code)
	}

	return response.Result.IndexPrice, nil
}

func (c *Client) makePublicRequest(method, endpoint string, params map[string]interface{}, result interface{}) error {
	return c.makeRequest(method, endpoint, params, result, false)
}

func (c *Client) makePrivateRequest(method, endpoint string, params map[string]interface{}, result interface{}) error {
	// 确保认证
	if err := c.ensureAuthenticated(); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	return c.makeRequest(method, endpoint, params, result, true)
}

func (c *Client) makeRequest(method, endpoint string, params map[string]interface{}, result interface{}, isPrivate bool) error {
	// 构建完整的 URL
	var fullURL string
	var req *http.Request
	var err error

	if method == "GET" && len(params) > 0 {
		fullURL = c.baseURL + endpoint + "?"
		for k, v := range params {
			fullURL += fmt.Sprintf("%s=%v&", k, v)
		}
		fullURL = fullURL[:len(fullURL)-1]
		req, err = http.NewRequest(method, fullURL, nil)
	} else if method == "GET" {
		fullURL = c.baseURL + endpoint
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
		fullURL = c.baseURL + endpoint
		jsonData, jsonErr := json.Marshal(params)
		if jsonErr != nil {
			return fmt.Errorf("failed to marshal params: %w", jsonErr)
		}
		req, err = http.NewRequest(method, fullURL, bytes.NewBuffer(jsonData))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	if err != nil {
		return fmt.Errorf("failed to create request for %s: %w", fullURL, err)
	}

	// 设置认证头
	if isPrivate {
		c.authMutex.RLock()
		token := c.accessToken
		c.authMutex.RUnlock()

		if token == "" {
			return fmt.Errorf("access token is required for private requests")
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed to %s: %w", fullURL, err)
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d for %s: %s", resp.StatusCode, fullURL, string(responseBody))
	}

	// 解析 JSON 响应
	if err := json.Unmarshal(responseBody, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(responseBody))
	}

	return nil
}
