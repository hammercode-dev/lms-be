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

func (c *Client) CreateInvoice(ctx context.Context, orderNo, email string, amount float64, description string, successRedirectURL string, failureRedirectURL string) (invoiceURL string, invoiceID string, err error) {
	if c.APIClient == nil {
		return "", "", fmt.Errorf("xendit client not initialized")
	}

	currency := "IDR"

	// Create request using NewCreateInvoiceRequest constructor
	req := *invoice.NewCreateInvoiceRequest(orderNo, amount)
	req.PayerEmail = &email
	req.Description = &description
	req.Currency = &currency

	// Set redirect URLs
	if successRedirectURL != "" {
		req.SuccessRedirectUrl = &successRedirectURL
		fmt.Printf("[XENDIT DEBUG] Setting success redirect URL: %s\n", successRedirectURL)
	}
	if failureRedirectURL != "" {
		req.FailureRedirectUrl = &failureRedirectURL
		fmt.Printf("[XENDIT DEBUG] Setting failure redirect URL: %s\n", failureRedirectURL)
	}

	resp, httpResp, err := c.APIClient.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(req).
		Execute()

	// Xendit SDK sometimes returns err even on success (false positive)
	// Priority: Check response first, ignore error if response is valid
	if resp != nil && resp.InvoiceUrl != "" {
		// Success - we have valid response
		var invID string
		if resp.Id != nil {
			invID = *resp.Id
		}
		return resp.InvoiceUrl, invID, nil
	}

	// Only treat as error if we don't have valid response
	if err != nil {
		status := 0
		if httpResp != nil {
			status = httpResp.StatusCode
		}
		return "", "", fmt.Errorf("xendit api error (status=%d): %v", status, err)
	}

	// No response and no error - shouldn't happen
	status := 0
	if httpResp != nil {
		status = httpResp.StatusCode
	}
	return "", "", fmt.Errorf("xendit returned nil response (status=%d)", status)
}

// GetInvoiceStatus retrieves invoice status from Xendit by invoice ID
func (c *Client) GetInvoiceStatus(ctx context.Context, invoiceID string) (status string, paidAt *string, paymentMethod *string, err error) {
	if c.APIClient == nil {
		return "", nil, nil, fmt.Errorf("xendit client not initialized")
	}

	if invoiceID == "" {
		return "", nil, nil, fmt.Errorf("invoice_id is empty")
	}

	// Debug log
	fmt.Printf("[XENDIT DEBUG] GetInvoiceStatus called with invoice_id: %s\n", invoiceID)

	resp, httpResp, err := c.APIClient.InvoiceApi.GetInvoiceById(ctx, invoiceID).Execute()

	// Debug response
	fmt.Printf("[XENDIT DEBUG] Response: resp=%v, httpStatus=%d, err=%v\n",
		resp != nil,
		func() int {
			if httpResp != nil {
				return httpResp.StatusCode
			}
			return 0
		}(),
		err)

	// Prioritize valid response over error (Xendit SDK issue)
	// Xendit SDK returns false positive errors even on HTTP 200 success
	if resp != nil {
		// Convert Xendit status to our status
		// Xendit statuses: PENDING, PAID, EXPIRED, SETTLED
		var finalStatus string
		switch resp.GetStatus() {
		case "PAID", "SETTLED":
			finalStatus = "paid"
		case "EXPIRED":
			finalStatus = "expired"
		default:
			finalStatus = "pending"
		}

		// Get payment details if paid
		var paidAtStr *string
		var payMethod *string

		if finalStatus == "paid" {
			// Use Updated field as paid_at (when invoice was last updated/paid)
			updatedTime := resp.Updated.Format("2006-01-02 15:04:05")
			paidAtStr = &updatedTime

			// Get payment method if available
			if resp.PaymentMethod != nil {
				method := string(*resp.PaymentMethod)
				payMethod = &method
			}
		}

		return finalStatus, paidAtStr, payMethod, nil
	}

	// No response - check error
	if err != nil {
		statusCode := 0
		if httpResp != nil {
			statusCode = httpResp.StatusCode
		}
		return "", nil, nil, fmt.Errorf("xendit api error (status=%d): %v", statusCode, err)
	}

	// No response and no error
	return "", nil, nil, fmt.Errorf("xendit returned nil response")
}
