package email

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
)

// Example 1: Send Welcome Email to New User
func sendWelcomeEmail(ctx context.Context) error {
	// HTML template for welcome email
	welcomeTemplate := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f4f4f4; }
        .button { background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block; margin-top: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Our Platform!</h1>
        </div>
        <div class="content">
            <h2>Hello {{.Name}},</h2>
            <p>Thank you for joining us! Your account has been successfully created.</p>
            <p><strong>Your Details:</strong></p>
            <ul>
                <li>Username: {{.Username}}</li>
                <li>Email: {{.Email}}</li>
                <li>Account Type: {{.AccountType}}</li>
            </ul>
            <p>To get started, please verify your email address by clicking the button below:</p>
            <a href="{{.VerificationLink}}" class="button">Verify Email</a>
            <p>If you have any questions, feel free to reply to this </p>
            <p>Best regards,<br>The Support Team</p>
        </div>
    </div>
</body>
</html>
`

	// Parse template
	tmpl, err := template.New("welcome").Parse(welcomeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// SMTP configuration (usually from environment variables)
	smtpConfig := SMTP{
		Email:    os.Getenv("SMTP_EMAIL"),    // e.g., "noreply@company.com"
		Password: os.Getenv("SMTP_PASSWORD"), // e.g., "your-app-password"
		Host:     os.Getenv("SMTP_HOST"),     // e.g., "smtp.gmail.com"
		Port:     os.Getenv("SMTP_PORT"),     // e.g., "587"
	}

	// Create email payload
	emailPayload := NewSendEmail(
		ctx,
		smtpConfig,
		"MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8",
		"Welcome to Our Platform - Please Verify Your Email",
		tmpl,
	)

	// Add primary recipient
	err = emailPayload.AddReceiver(ctx, Receiver{
		Email: "afandywibowo979@gmail.com",
		Data: map[string]interface{}{
			"Name":             "John Doe",
			"Username":         "johndoe",
			"Email":            "afandywibowo979@gmail.com",
			"AccountType":      "Premium",
			"VerificationLink": "https://example.com/verify?token=abc123xyz",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add receiver: %w", err)
	}

	// Add CC to admin for monitoring new registrations
	err = emailPayload.AddCC(ctx, "afandywibowo2000@gmail.com")
	if err != nil {
		return fmt.Errorf("failed to add CC: %w", err)
	}

	// Send the email
	emailPayload.SendEmail(ctx)

	return nil
}

// Example 2: Send Order Confirmation with Multiple CC
func sendOrderConfirmation(ctx context.Context, orderID string) error {
	// Order confirmation template
	orderTemplate := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .order-box { border: 1px solid #ddd; padding: 20px; margin: 20px 0; }
        .order-header { background-color: #f8f9fa; padding: 10px; margin: -20px -20px 20px -20px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
        .total { font-weight: bold; font-size: 1.2em; }
    </style>
</head>
<body>
    <div class="order-box">
        <div class="order-header">
            <h2>Order Confirmation #{{.OrderID}}</h2>
        </div>
        <p>Dear {{.CustomerName}},</p>
        <p>Thank you for your order! We've received your payment and will begin processing your order shortly.</p>
        
        <h3>Order Details:</h3>
        <table>
            <tr>
                <th>Item</th>
                <th>Quantity</th>
                <th>Price</th>
                <th>Subtotal</th>
            </tr>
            {{range .Items}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.Quantity}}</td>
                <td>${{.Price}}</td>
                <td>${{.Subtotal}}</td>
            </tr>
            {{end}}
            <tr class="total">
                <td colspan="3">Total:</td>
                <td>${{.Total}}</td>
            </tr>
        </table>
        
        <h3>Shipping Information:</h3>
        <p>{{.ShippingAddress}}</p>
        <p>Estimated Delivery: {{.EstimatedDelivery}}</p>
        
        <p>You can track your order status at: <a href="{{.TrackingLink}}">Track Order</a></p>
    </div>
</body>
</html>
`

	tmpl, err := template.New("order").Parse(orderTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// SMTP configuration
	smtpConfig := SMTP{
		Email:    os.Getenv("SMTP_EMAIL"),    // e.g., "noreply@company.com"
		Password: os.Getenv("SMTP_PASSWORD"), // e.g., "your-app-password"
		Host:     os.Getenv("SMTP_HOST"),     // e.g., "smtp.gmail.com"
		Port:     os.Getenv("SMTP_PORT"),     // e.g., "587"
	}

	// Create email payload
	emailPayload := NewSendEmail(
		ctx,
		smtpConfig,
		"MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8",
		fmt.Sprintf("Order Confirmation #%s", orderID),
		tmpl,
	)

	// Order data
	orderData := map[string]interface{}{
		"OrderID":      orderID,
		"CustomerName": "Jane Smith",
		"Items": []map[string]interface{}{
			{"Name": "Product A", "Quantity": 2, "Price": 29.99, "Subtotal": 59.98},
			{"Name": "Product B", "Quantity": 1, "Price": 49.99, "Subtotal": 49.99},
			{"Name": "Product C", "Quantity": 3, "Price": 9.99, "Subtotal": 29.97},
		},
		"Total":             139.94,
		"ShippingAddress":   "123 Main St, Apt 4B, New York, NY 10001",
		"EstimatedDelivery": "3-5 business days",
		"TrackingLink":      fmt.Sprintf("https://example.com/track/%s", orderID),
	}

	// Add customer as primary recipient
	err = emailPayload.AddReceiver(ctx, Receiver{
		Email: "customer@example.com",
		Data:  orderData,
	})
	if err != nil {
		return fmt.Errorf("failed to add receiver: %w", err)
	}

	// Add multiple CC recipients
	ccRecipients := []string{
		"sales@company.com",      // Sales team
		"warehouse@company.com",  // Warehouse for fulfillment
		"accounting@company.com", // Accounting for records
	}

	err = emailPayload.AddMultipleCC(ctx, ccRecipients)
	if err != nil {
		return fmt.Errorf("failed to add CC recipients: %w", err)
	}

	// Send the email
	emailPayload.SendEmail(ctx)

	return nil
}

// Example 3: Send Bulk Newsletter to Multiple Recipients
func sendNewsletter(ctx context.Context) error {
	// Newsletter template
	newsletterTemplate := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Georgia, serif; }
        .newsletter { max-width: 600px; margin: 0 auto; }
        .header { text-align: center; padding: 30px 0; background-color: #2c3e50; color: white; }
        .content { padding: 30px; }
        .article { margin-bottom: 30px; padding-bottom: 20px; border-bottom: 1px solid #eee; }
        .footer { background-color: #ecf0f1; padding: 20px; text-align: center; }
    </style>
</head>
<body>
    <div class="newsletter">
        <div class="header">
            <h1>{{.Title}}</h1>
            <p>{{.Date}}</p>
        </div>
        <div class="content">
            <p>Hello {{.Name}},</p>
            {{range .Articles}}
            <div class="article">
                <h2>{{.Title}}</h2>
                <p>{{.Summary}}</p>
                <a href="{{.Link}}">Read more →</a>
            </div>
            {{end}}
        </div>
        <div class="footer">
            <p>You're receiving this because you're subscribed to our newsletter.</p>
            <p><a href="{{.UnsubscribeLink}}">Unsubscribe</a> | <a href="{{.PreferencesLink}}">Update Preferences</a></p>
        </div>
    </div>
</body>
</html>
`

	tmpl, err := template.New("newsletter").Parse(newsletterTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// SMTP configuration
	smtpConfig := SMTP{
		Email:    os.Getenv("SMTP_EMAIL"),    // e.g., "noreply@company.com"
		Password: os.Getenv("SMTP_PASSWORD"), // e.g., "your-app-password"
		Host:     os.Getenv("SMTP_HOST"),     // e.g., "smtp.gmail.com"
		Port:     os.Getenv("SMTP_PORT"),     // e.g., "587"
	}

	// Create email payload
	emailPayload := NewSendEmail(
		ctx,
		smtpConfig,
		"MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8",
		"Monthly Newsletter - December 2024",
		tmpl,
	)

	// Newsletter data (same for all recipients)
	newsletterData := map[string]interface{}{
		"Title": "Monthly Newsletter",
		"Date":  "December 2024",
		"Articles": []map[string]interface{}{
			{
				"Title":   "New Feature Release",
				"Summary": "We're excited to announce our latest feature that will help you...",
				"Link":    "https://example.com/blog/new-feature",
			},
			{
				"Title":   "Customer Success Story",
				"Summary": "Learn how Company XYZ increased their productivity by 50%...",
				"Link":    "https://example.com/blog/success-story",
			},
			{
				"Title":   "Upcoming Webinar",
				"Summary": "Join us for a deep dive into best practices...",
				"Link":    "https://example.com/webinar/register",
			},
		},
	}

	// Add multiple subscribers
	subscribers := []struct {
		Email string
		Name  string
	}{
		{"subscriber1@example.com", "Alice Johnson"},
		{"subscriber2@example.com", "Bob Williams"},
		{"subscriber3@example.com", "Carol Davis"},
	}

	for _, subscriber := range subscribers {
		// Create personalized data for each subscriber
		personalizedData := make(map[string]interface{})
		for k, v := range newsletterData {
			personalizedData[k] = v
		}
		personalizedData["Name"] = subscriber.Name
		personalizedData["UnsubscribeLink"] = fmt.Sprintf("https://example.com/unsubscribe?email=%s", subscriber.Email)
		personalizedData["PreferencesLink"] = fmt.Sprintf("https://example.com/preferences?email=%s", subscriber.Email)

		err = emailPayload.AddReceiver(ctx, Receiver{
			Email: subscriber.Email,
			Data:  personalizedData,
		})
		if err != nil {
			log.Printf("Failed to add subscriber %s: %v", subscriber.Email, err)
			continue
		}
	}

	// Add marketing team as CC
	err = emailPayload.AddCC(ctx, "marketing@company.com")
	if err != nil {
		log.Printf("Failed to add CC: %v", err)
	}

	// Send the newsletter (individual emails to each subscriber)
	emailPayload.SendEmail(ctx)

	// Alternative: Use SendBulkEmail for single email to all recipients
	// Note: This sends one email with all recipients visible to each other
	// emailPayload.SendBulkEmail(ctx)

	return nil
}

// Example 4: Send Password Reset Email
func sendPasswordReset(ctx context.Context, userEmail, resetToken string) error {
	resetTemplate := `
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2>Password Reset Request</h2>
        <p>Hello,</p>
        <p>We received a request to reset your password. Click the button below to create a new password:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.ResetLink}}" style="background-color: #007bff; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
        </div>
        <p>This link will expire in {{.ExpiryTime}}.</p>
        <p>If you didn't request this, please ignore this email and your password will remain unchanged.</p>
        <hr style="margin: 30px 0;">
        <p style="color: #666; font-size: 0.9em;">For security reasons, this link can only be used once.</p>
    </div>
</body>
</html>
`

	tmpl, err := template.New("reset").Parse(resetTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	smtpConfig := SMTP{
		Email:    os.Getenv("SMTP_EMAIL"),    // e.g., "noreply@company.com"
		Password: os.Getenv("SMTP_PASSWORD"), // e.g., "your-app-password"
		Host:     os.Getenv("SMTP_HOST"),     // e.g., "smtp.gmail.com"
		Port:     os.Getenv("SMTP_PORT"),     // e.g., "587"
	}

	emailPayload := NewSendEmail(
		ctx,
		smtpConfig,
		"MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8",
		"Password Reset Request",
		tmpl,
	)

	err = emailPayload.AddReceiver(ctx, Receiver{
		Email: userEmail,
		Data: map[string]interface{}{
			"ResetLink":  fmt.Sprintf("https://example.com/reset-password?token=%s", resetToken),
			"ExpiryTime": "24 hours",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add receiver: %w", err)
	}

	// CC to security team for monitoring
	err = emailPayload.AddCC(ctx, "security-alerts@company.com")
	if err != nil {
		log.Printf("Failed to add security CC: %v", err)
	}

	emailPayload.SendEmail(ctx)
	return nil
}

// Main function showing how to use the examples
func main() {
	ctx := context.Background()

	// Set up environment variables (in production, these would be set externally)
	os.Setenv("SMTP_EMAIL", "tested@gmail.com")
	os.Setenv("SMTP_PASSWORD", "1234")
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "587")

	// Example 1: Send welcome email
	if err := sendWelcomeEmail(ctx); err != nil {
		log.Printf("Failed to send welcome email: %v", err)
	}

	// // Example 2: Send order confirmation
	// if err := sendOrderConfirmation(ctx, "ORD-2024-12345"); err != nil {
	// 	log.Printf("Failed to send order confirmation: %v", err)
	// }

	// // Example 3: Send newsletter
	// if err := sendNewsletter(ctx); err != nil {
	// 	log.Printf("Failed to send newsletter: %v", err)
	// }

	// // Example 4: Send password reset
	// if err := sendPasswordReset(ctx, "user@example.com", "reset-token-xyz"); err != nil {
	// 	log.Printf("Failed to send password reset: %v", err)
	// }
}
