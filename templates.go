package postmark

import (
	"errors"
	"fmt"
)

// EmailWithTemplate is used to send an email via a template
type EmailWithTemplate struct {
	// TemplateId: REQUIRED if TemplateAlias is not specified. - The template id to use when sending this message.
	TemplateID int64 `json:"TemplateId,omitempty"`
	// TemplateAlias: REQUIRED if TemplateId is not specified. - The template alias to use when sending this message.
	TemplateAlias string `json:"TemplateAlias,omitempty"`
	// TemplateModel: The model to be applied to the specified template to generate HtmlBody, TextBody, and Subject.
	TemplateModel map[string]interface{} `json:"TemplateModel,omitempty"`
	// InlineCss: By default, if the specified template contains an HTMLBody, we will apply the style blocks as inline attributes to the rendered HTML content. You may opt-out of this behavior by passing false for this request field.
	InlineCCC bool `json:"InlineCss,omitempty"`
	// From: The sender email address. Must have a registered and confirmed Sender Signature.
	From     string  `json:"From,omitempty"`
	FromName *string `json:"-"`
	// To: REQUIRED Recipient email address. Multiple addresses are comma separated. Max 50.
	To     string  `json:"To,omitempty"`
	ToName *string `json:"-"`
	// Cc recipient email address. Multiple addresses are comma separated. Max 50.
	Cc string `json:"Cc,omitempty"`
	// Bcc recipient email address. Multiple addresses are comma separated. Max 50.
	Bcc string `json:"Bcc,omitempty"`
	// Tag: Email tag that allows you to categorize outgoing emails and get detailed statistics.
	Tag string `json:"Tag,omitempty"`
	// Reply To override email address. Defaults to the Reply To set in the sender signature.
	ReplyTo string `json:"ReplyTo,omitempty"`
	// Headers: List of custom headers to include.
	Headers []Header `json:"Headers,omitempty"`
	// TrackOpens: Activate open tracking for this email.
	TrackOpens bool `json:"TrackOpens,omitempty"`
	// Attachments: List of attachments
	Attachments []Attachment `json:"Attachments,omitempty"`
}

// SendEmailWithTemplate sends an email using a template (TemplateId)
func (client *Client) SendEmailWithTemplate(email *EmailWithTemplate) (*EmailResponse, error) {
	res := &EmailResponse{}

	if email == nil {
		return res, errors.New("The email object is not set")
	}

	// For email addresses that have names or titles with punctuation, you should quote them as such: "To" : "\"Joe Receiver, jr\" <receiver@example.com>"
	if email.FromName != nil {
		email.From = fmt.Sprintf(`"%v" %v`, *email.FromName, email.From)
	}

	if email.ToName != nil {
		email.To = fmt.Sprintf(`"%v" %v`, *email.ToName, email.To)
	}

	err := client.doRequest(parameters{
		Method:    "POST",
		Path:      "email/withTemplate",
		Payload:   email,
		TokenType: serverToken,
	}, &res)
	return res, err
}

// SendBatchEmailWithTemplate sends batch email using a template (TemplateId)
func (client *Client) SendBatchEmailWithTemplate(emails *[]EmailWithTemplate) (*[]EmailResponse, error) {
	res := &[]EmailResponse{}

	if emails == nil {
		return res, errors.New("The emails object is not set")
	}

	editedEmails := make([]EmailWithTemplate, len(*emails))
	for i, e := range *emails {
		// For email addresses that have names or titles with punctuation, you should quote them as such: "To" : "\"Joe Receiver, jr\" <receiver@example.com>"
		if e.FromName != nil {
			e.From = fmt.Sprintf(`"%v" %v`, *e.FromName, e.From)
		}

		if e.ToName != nil {
			e.To = fmt.Sprintf(`"%v" %v`, *e.ToName, e.To)
		}

		editedEmails[i] = e
	}

	var formatEmails map[string]interface{} = map[string]interface{}{
		"Messages": editedEmails,
	}
	err := client.doRequest(parameters{
		Method:    "POST",
		Path:      "email/batchWithTemplates",
		Payload:   formatEmails,
		TokenType: serverToken,
	}, &res)
	return res, err
}
