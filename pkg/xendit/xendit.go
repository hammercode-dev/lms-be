package xendit

import (
	"context"
	"fmt"

	xendit "github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/invoice"
)

type Client struct {
	APIClient *xendit.APIClient
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIClient: xendit.NewClient(apiKey),
	}
}

func (c *Client) CreateInvoice(ctx context.Context, orderNo, email string, amount float64, description string) (string, error) {
	if c.APIClient == nil {
		return "", fmt.Errorf("xendit client not initialized")
	}

	currency := "IDR"


	// Create request using NewCreateInvoiceRequest constructor
	req := *invoice.NewCreateInvoiceRequest(orderNo, amount)
	req.PayerEmail = &email
	req.Description = &description
	req.Currency = &currency

	resp, httpResp, err := c.APIClient.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(req).
		Execute()

	if err != nil {
		status := 0
		if httpResp != nil {
			status = httpResp.StatusCode
		}
		return "", fmt.Errorf("xendit api error (status=%d): %v", status, err)
	}

	if resp == nil {
		status := 0
		if httpResp != nil {
			status = httpResp.StatusCode
		}
		return "", fmt.Errorf("xendit returned nil response (status=%d)", status)
	}

	return resp.InvoiceUrl, nil
}
