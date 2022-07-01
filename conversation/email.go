package conversation

// Email
// https://developers.messagebird.com/api/conversations/#email-object
type Email struct {
	Id                   string              `json:"id"`
	To                   []*EmailRecipient   `json:"to"`
	From                 *EmailRecipient     `json:"from"`
	Subject              string              `json:"subject"`
	Content              *EmailContent       `json:"content"`
	ReplyTo              string              `json:"replyTo,omitempty"`
	ReturnPath           string              `json:"returnPath,omitempty"`
	Headers              interface{}         `json:"headers,omitempty"`
	Tracking             *EmailTracking      `json:"tracking,omitempty"`
	ReportUrl            string              `json:"reportUrl,omitempty"`
	PerformSubstitutions bool                `json:"performSubstitutions,omitempty"`
	Attachments          []*EmailAttachment  `json:"attachments,omitempty"`
	InlineImages         []*EmailInlineImage `json:"inlineImages,omitempty"`
}

// EmailRecipient
// https://developers.messagebird.com/api/conversations/#emailrecipient-object
type EmailRecipient struct {
	Address   string                   `json:"address"`
	Name      string                   `json:"name,omitempty"`
	Variables *EmailRecipientVariables `json:"variables,omitempty"`
}

type EmailRecipientVariables struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// EmailContent
// https://developers.messagebird.com/api/conversations/#emailcontent-object
type EmailContent struct {
	Html string `json:"html,omitempty"`
	Text string `json:"text,omitempty"`
}

// EmailTracking
// https://developers.messagebird.com/api/conversations/#emailtracking-object
type EmailTracking struct {
	Open  bool `json:"open"`
	Click bool `json:"click"`
}

// EmailAttachment
// The Attachment object represents a file attached to a particular message. The maximum attachment size is 20 MB.
// https://developers.messagebird.com/api/conversations/#emailattachment-object
type EmailAttachment struct {
	Id     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	URL    string `json:"URL,omitempty"`
	Length int    `json:"length,omitempty"`
}

// EmailInlineImage
// https://developers.messagebird.com/api/conversations/#emailinlineimage-object
type EmailInlineImage struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	URL       string `json:"URL,omitempty"`
	Length    int    `json:"length,omitempty"`
	ContentId string `json:"contentId,omitempty"`
}
