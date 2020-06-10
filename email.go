package postmark

// https://github.com/keighl/postmark
// https://github.com/mattevans/postmark-go

import (
	"errors"
	"fmt"
	"time"
)

// Email is exactly what it sounds like
type Email struct {
	// From - REQUIRED The sender email address. Must have a registered and confirmed Sender Signature.
	From string `json:"From"`
	// To - REQUIRED Recipient email address. Multiple addresses are comma separated. Max 50.
	To string `json:"To"`
	// Cc - recipient email address. Multiple addresses are comma separated. Max 50.
	Cc string `json:"Cc,omitempty"`
	// Bcc - recipient email address. Multiple addresses are comma separated. Max 50.
	Bcc string `json:"Bcc,omitempty"`
	// Subject - Email subject
	Subject string `json:"Subject,omitempty"`
	// Tag - Email tag that allows you to categorize outgoing emails and get detailed statistics.
	Tag string `json:"Tag,omitempty"`
	// HTMLBody - HTML email message. REQUIRED, If no TextBody specified
	HTMLBody string `json:"HtmlBody,omitempty"`
	// TextBody - Plain text email message. REQUIRED, If no HtmlBody specified
	TextBody string `json:"TextBody,omitempty"`
	// ReplyTo - Reply To override email address. Defaults to the Reply To set in the sender signature.
	ReplyTo string `json:"ReplyTo,omitempty"`
	// Headers - List of custom headers to include.
	Headers []Header `json:"Headers,omitempty"`
	// TrackOpens - Activate open tracking for this email.
	TrackOpens bool `json:"TrackOpens,omitempty"`
	// Attachments - List of attachments
	Attachments []Attachment `json:"Attachments,omitempty"`
	// Metadata - metadata
	Metadata map[string]string `json:"Metadata,omitempty"`
}

// Header - an email header
type Header struct {
	// Name: header name
	Name string
	// Value: header value
	Value string
}

// Attachment is an optional encoded file to send along with an email
type Attachment struct {
	// Name: attachment name
	Name string
	// Content: Base64 encoded attachment data
	Content string
	// ContentType: attachment MIME type
	ContentType string
	// ContentId: populate for inlining images with the images cid
	ContentID string `json:",omitempty"`
}

// EmailResponse holds info in response to a send/send-batch request
// Even if API request comes back successful, check the ErrorCode to see if there might be a delivery problem
type EmailResponse struct {
	// To: Recipient email address
	To string
	// SubmittedAt: Timestamp
	SubmittedAt time.Time
	// MessageID: ID of message
	MessageID string
	// ErrorCode: API Error Codes
	ErrorCode int64
	// Message: Response message
	Message string
}

// SendEmail sends, well, an email.
func (client *Client) SendEmail(email *Email) (EmailResponse, error) {
	res := EmailResponse{}

	if email == nil {
		return res, errors.New("The email object is not set")
	}

	err := client.doRequest(parameters{
		Method:    "POST",
		Path:      "email",
		Payload:   email,
		TokenType: serverToken,
	}, &res)

	if res.ErrorCode != 0 {
		return res, fmt.Errorf(`%v %s`, res.ErrorCode, res.Message)
	}

	return res, err
}

// SendEmailBatch sends multiple emails together
// Note, individual emails in the batch can error, so it would be wise to
// range over the responses and sniff for errors
func (client *Client) SendEmailBatch(emails []Email) ([]EmailResponse, error) {
	res := []EmailResponse{}
	err := client.doRequest(parameters{
		Method:    "POST",
		Path:      "email/batch",
		Payload:   emails,
		TokenType: serverToken,
	}, &res)
	return res, err
}
