package deribit

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"cs-projects-eth-collar/internal/types"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    string
	httpClient *http.Client
}

type APIResponse struct {
	Result interface{} `json:"result"`
	Error  *APIError   `json:"error"`
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

func (c *Client) GetAccountSummary(currency string) (*types.AccountSummary, error) {
	endpoint := "/private/get_account_summary"
	params := map[string]interface{}{
		"currency": currency,
	}

	var response struct {
		Result types.AccountSummary `json:"result"`
		Error  *APIError            `json:"error"`
	}

	if err := c.makeRequest("GET", endpoint, params, &response); err != nil {
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

	if err := c.makeRequest("GET", endpoint, params, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("API error: %s (code: %d)", response.Error.Message, response.Error.Code)
	}

	return response.Result, nil
}

// GetIndexPrice 获取指数价格 (现货价格)
func (c *Client) GetIndexPrice(currency string) (float64, error) {
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

	if err := c.makeRequest("GET", endpoint, params, &response); err != nil {
		return 0, err
	}

	if response.Error != nil {
		return 0, fmt.Errorf("API error: %s (code: %d)", response.Error.Message, response.Error.Code)
	}

	return response.Result.IndexPrice, nil
}

func (c *Client) makeRequest(method, endpoint string, params map[string]interface{}, result interface{}) error {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce := timestamp

	jsonData, err := json.Marshal(params)
	if err != nil {
		return err
	}

	body := string(jsonData)
	if method == "GET" && len(params) > 0 {
		body = ""
		url := c.baseURL + endpoint + "?"
		for k, v := range params {
			url += fmt.Sprintf("%s=%v&", k, v)
		}
		url = url[:len(url)-1]
	}

	requestData := method + "\n" + endpoint + "\n" + body + "\n"
	signature := c.generateSignature(timestamp, nonce, requestData)

	var req *http.Request
	if method == "GET" {
		url := c.baseURL + endpoint
		if len(params) > 0 {
			url += "?"
			for k, v := range params {
				url += fmt.Sprintf("%s=%v&", k, v)
			}
			url = url[:len(url)-1]
		}
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, c.baseURL+endpoint, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("deri-hmac-sha256 id=%s,ts=%s,nonce=%s,sig=%s",
		c.apiKey, timestamp, nonce, signature))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d, body: %s", resp.StatusCode, string(responseBody))
	}

	return json.Unmarshal(responseBody, result)
}

func (c *Client) generateSignature(timestamp, nonce, requestData string) string {
	message := timestamp + "\n" + nonce + "\n" + requestData
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
