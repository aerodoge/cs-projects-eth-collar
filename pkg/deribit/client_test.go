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
