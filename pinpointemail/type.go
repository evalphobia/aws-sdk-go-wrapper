package pinpointemail

import (
	SDK "github.com/aws/aws-sdk-go/service/pinpointemail"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	defaultCharset = "UTF-8"
)

type EmailInput struct {
	From    string
	ReplyTo []string

	// destination
	To  []string
	Cc  []string
	Bcc []string

	// message
	Content Content

	// email address for bounce
	FeedbackForwardingEmail string

	ConfigurationSetName string
	Tags                 []Tag
}

func (in EmailInput) ToInput() *SDK.SendEmailInput {
	input := &SDK.SendEmailInput{
		Content:          in.Content.ToContent(),
		ReplyToAddresses: toSlicePointer(in.ReplyTo),
		EmailTags:        toTags(in.Tags),
	}

	if in.From != "" {
		input.FromEmailAddress = pointers.String(in.From)
	}

	switch {
	case len(in.To) != 0,
		len(in.Cc) != 0,
		len(in.Bcc) != 0:
		input.Destination = &SDK.Destination{
			ToAddresses:  toSlicePointer(in.To),
			CcAddresses:  toSlicePointer(in.Cc),
			BccAddresses: toSlicePointer(in.Bcc),
		}
	}

	if in.FeedbackForwardingEmail != "" {
		input.FeedbackForwardingEmailAddress = pointers.String(in.FeedbackForwardingEmail)
	}
	if in.ConfigurationSetName != "" {
		input.ConfigurationSetName = pointers.String(in.ConfigurationSetName)
	}

	return input
}

// Content has the entire content of the email.
type Content struct {
	// If use RawMessage, below fields are ignored.
	// see criteria: https://github.com/aws/aws-sdk-go/blob/master/service/pinpointemail/api.go#L5877
	RawMessage []byte

	// Simple messages
	SubjectCharset string
	Subject        string
	BodyCharset    string
	Body           string
	HTML           bool
}

func (c Content) ToContent() *SDK.EmailContent {
	ec := &SDK.EmailContent{}
	if len(c.RawMessage) != 0 {
		// Use raw message
		ec.Raw = &SDK.RawMessage{
			Data: c.RawMessage,
		}
		return ec
	}

	// Use simple message
	if c.SubjectCharset == "" {
		c.SubjectCharset = defaultCharset
	}
	if c.BodyCharset == "" {
		c.BodyCharset = defaultCharset
	}

	content := &SDK.Content{
		Charset: pointers.String(c.BodyCharset),
		Data:    pointers.String(c.Body),
	}
	body := &SDK.Body{}
	switch {
	case c.HTML:
		body.Html = content
	default:
		body.Text = content
	}

	ec.Simple = &SDK.Message{
		Subject: &SDK.Content{
			Charset: pointers.String(c.SubjectCharset),
			Data:    pointers.String(c.Subject),
		},
		Body: body,
	}
	return ec
}

type Tag struct {
	Name  string
	Value string
}

func (t Tag) ToTag() *SDK.MessageTag {
	return &SDK.MessageTag{
		Name:  pointers.String(t.Name),
		Value: pointers.String(t.Value),
	}
}

func toTags(list []Tag) []*SDK.MessageTag {
	if len(list) == 0 {
		return nil
	}

	result := make([]*SDK.MessageTag, len(list))
	for i, tag := range list {
		result[i] = tag.ToTag()
	}
	return result
}

func toSlicePointer(list []string) []*string {
	if len(list) == 0 {
		return nil
	}

	result := make([]*string, len(list))
	for i, v := range list {
		result[i] = pointers.String(v)
	}
	return result
}
