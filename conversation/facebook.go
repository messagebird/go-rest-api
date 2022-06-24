package conversation

// FacebookMessage
// https://developers.messagebird.com/api/conversations/#facebookmessage-object
type FacebookMessage struct {
	Text         string                `json:"text,omitempty"`
	Attachment   *FacebookAttachment   `json:"attachment,omitempty"`
	QuickReplies []*FacebookQuickReply `json:"quick_replies,omitempty"`
}

type FacebookAttachmentType string

const (
	FBAttachmentTypeImage    FacebookAttachmentType = "image"
	FBAttachmentTypeAudio    FacebookAttachmentType = "audio"
	FBAttachmentTypeVideo    FacebookAttachmentType = "video"
	FBAttachmentTypeFile     FacebookAttachmentType = "file"
	FBAttachmentTypeLocation FacebookAttachmentType = "location"
	FBAttachmentTypeFallback FacebookAttachmentType = "fallback"
	FBAttachmentTypeTemplate FacebookAttachmentType = "template"
)

// FacebookAttachment
// https://developers.messagebird.com/api/conversations/#facebookattachment-object
type FacebookAttachment struct {
	Type    FacebookAttachmentType     `json:"type"`
	Payload *FacebookAttachmentPayload `json:"payload"`
}

// FacebookAttachmentPayload
// https://developers.messagebird.com/api/conversations/#facebookattachmentpayload-object
type FacebookAttachmentPayload struct {
	Url              string                    `json:"url,omitempty"`
	IsReusable       bool                      `json:"is_reusable"`
	AttachmentId     string                    `json:"attachment_id,omitempty"`
	TemplateType     FacebookTemplateType      `json:"template_type,omitempty"`
	Elements         []*FacebookElement        `json:"elements,omitempty"`
	ImageAspectRatio *FacebookImageAspectRatio `json:"image_aspect_ratio,omitempty"`
}

type FacebookTemplateType string

const (
	FBTemplateTypeMedia   FacebookTemplateType = "media"
	FBTemplateTypeGeneric FacebookTemplateType = "generic"
)

// FacebookElement
// https://developers.messagebird.com/api/conversations/#facebookelement-object
type FacebookElement struct {
	MediaType     FacebookElementMediaType `json:"media_type"`
	AttachmentId  string                   `json:"attachment_id,omitempty"`
	MediaUrl      string                   `json:"media_url,omitempty"`
	Buttons       []*FacebookButton        `json:"buttons,omitempty"`
	Title         string                   `json:"title,omitempty"`
	Subtitle      string                   `json:"subtitle,omitempty"`
	DefaultAction *FacebookButton          `json:"default_action,omitempty"`
	ImageUrl      string                   `json:"image_url,omitempty"`
}

type FacebookElementMediaType string

const (
	FBElementMediaTypeVideo FacebookElementMediaType = "video"
	FBElementMediaTypeImage FacebookElementMediaType = "image"
)

type FacebookButtonType string

const (
	FBButtonTypeWebUrl      FacebookButtonType = "web_url"
	FBButtonTypePhoneNumber FacebookButtonType = "phone_number"
	FBButtonTypePostback    FacebookButtonType = "postback"
)

// FacebookButton
// https://developers.messagebird.com/api/conversations/#facebookbutton-object
type FacebookButton struct {
	Type    FacebookButtonType `json:"type"`
	Url     string             `json:"url"`
	Title   string             `json:"title"`
	Payload string             `json:"payload"`
}

// FacebookImageAspectRatio
// https://developers.messagebird.com/api/conversations/#facebookimageaspectratio-object
type FacebookImageAspectRatio string

const (
	FBImageAspectRatioHorizontal FacebookImageAspectRatio = "horizontal"
	FBImageAspectRatioSquare     FacebookImageAspectRatio = "square"
)

// FacebookQuickReply
// https://developers.messagebird.com/api/conversations/#facebookquickreply-object
type FacebookQuickReply struct {
	ContentType FacebookQuickReplyContentType `json:"content_type"`
	Title       string                        `json:"title"`
	Payload     string                        `json:"payload"`
	ImageUrl    string                        `json:"image_url"`
}

// FacebookQuickReplyContentType
// https://developers.messagebird.com/api/conversations/#facebookquickreplycontenttype-object
type FacebookQuickReplyContentType string

const (
	FBQuickReplyContentTypeText            FacebookQuickReplyContentType = "text"
	FBQuickReplyContentTypeUserPhoneNumber FacebookQuickReplyContentType = "user_phone_number"
	FBQuickReplyContentTypeUserEmail       FacebookQuickReplyContentType = "user_email"
)
