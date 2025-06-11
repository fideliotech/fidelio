package main

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendConfirmationEmail(name, email string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		logger.Error("[Resend] RESEND_API_KEY environment variable is not set!")
		return fmt.Errorf("missing RESEND_API_KEY")
	}

	client := resend.NewClient(apiKey)

	emailBody := fmt.Sprintf(`
		<h2>Thank you for contacting fidelio</h2>
		<p>Dear %s,</p>
		<p>Thank you for reaching out to us. We have received your message and will get back to you as soon as possible.</p>
		<p>Best regards,<br>fidelio team</p>
	`, name)

	params := &resend.SendEmailRequest{
		From:    "fidelio <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Thank you for contacting fidelio",
		Html:    emailBody,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		logger.Error("[Resend] Failed to send confirmation email to %s: %v", email, err)
		return err
	}
	logger.Info("[Resend] Confirmation email sent successfully to %s, ID: %v", email, sent.Id)
	return nil
}

func SendContactNotificationEmail(name, email, company, message, industry, interest string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	to := "technology@konargo.com"

	if apiKey == "" {
		logger.Error("[Resend] RESEND_API_KEY environment variable is not set!")
		return fmt.Errorf("missing RESEND_API_KEY")
	}

	client := resend.NewClient(apiKey)

	// Send notification to admin
	emailBody := fmt.Sprintf(`
		<h2>New contact/quote request</h2>
		<p><strong>Name:</strong> %s</p>
		<p><strong>Email:</strong> %s</p>
		<p><strong>Company:</strong> %s</p>
		<p><strong>Industry:</strong> %s</p>
		<p><strong>Interest:</strong> %s</p>
		<p><strong>Message:</strong> %s</p>
	`, name, email, company, industry, interest, message)

	params := &resend.SendEmailRequest{
		From:    "fidelio <onboarding@resend.dev>",
		To:      []string{to},
		Subject: "New contact/quote request from fidelio's website",
		Html:    emailBody,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		logger.Error("[Resend] Failed to send notification email: %v", err)
		return err
	}
	logger.Info("[Resend] Notification email sent successfully, ID: %v", sent.Id)

	// Send confirmation to user
	err = SendConfirmationEmail(name, email)
	if err != nil {
		logger.Error("[Resend] Failed to send confirmation email: %v", err)
		// Don't return error here as the main notification was sent successfully
	}

	return nil
}
