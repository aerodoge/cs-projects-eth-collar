package deribit

import (
	"cs-projects-eth-collar/internal/types"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestClient(t *testing.T) *Client {
	cfg := types.DeribitConfig{
		APIKey:    "",
		APISecret: "",
		TestNet:   false,
	}
	return NewClient(cfg)
}

func TestGetIndexPrice(t *testing.T) {
	client := setupTestClient(t)
	price, err := client.GetIndexPrice("eth")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", price)
}

func TestGetAccountSummary(t *testing.T) {
	client := setupTestClient(t)
	summary, err := client.GetAccountSummary("ETH")

	assert.Nil(t, err)
	fmt.Printf("%+v", summary)
}

func TestGetAccountSummaries(t *testing.T) {
	client := setupTestClient(t)
	summaries, err := client.GetAccountSummaries()

	assert.Nil(t, err)
	assert.NotNil(t, summaries)
	assert.Greater(t, len(summaries.Summaries), 0)

	fmt.Printf("Change Margin Model API Limit: %+v\n", summaries.ChangeMarginModelAPILimit)
	fmt.Printf("Number of currency summaries: %d\n", len(summaries.Summaries))

	// 打印前几个货币的信息
	for i, summary := range summaries.Summaries {
		if i >= 3 { // 只打印前3个
			break
		}
		fmt.Printf("Currency: %s, Balance: %.6f, Maintenance Margin: %.6f\n",
			summary.Currency, summary.Balance, summary.MaintenanceMargin)
	}
}

func TestGetAccountSummariesExtended(t *testing.T) {
	client := setupTestClient(t)
	summaries, err := client.GetAccountSummaries(true)

	assert.Nil(t, err)
	assert.NotNil(t, summaries)

	// 在extended模式下，应该有更多的账户信息
	fmt.Printf("Extended info - Username: %s, Email: %s\n", summaries.Username, summaries.Email)
	fmt.Printf("Number of currency summaries: %d\n", len(summaries.Summaries))
}
