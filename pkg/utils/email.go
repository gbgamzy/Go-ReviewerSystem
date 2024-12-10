package utils

import (
	"log"
)

// SendEmail simulates sending an email to the specified recipient.
func SendEmail(to, subject, body string) error {
	// Simulate sending email (mock implementation)
	log.Printf("Sending email to %s\nSubject: %s\nBody: %s\n", to, subject, body)

	// In a real application, you can integrate an email service like SendGrid, AWS SES, etc.
	// Example:
	// err := smtp.SendMail(...)
	// if err != nil {
	//     return fmt.Errorf("failed to send email: %v", err)
	// }

	return nil // No error in the mock implementation
}
