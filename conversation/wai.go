package conversation

// WhatsAppInteractive
// https://developers.messagebird.com/api/conversations/#whatsappinteractive-object
type WhatsAppInteractive struct {
	Type   *WhatsAppInteractiveType   `json:"type"`
	Header *WhatsAppInteractiveHeader `json:"header"`
	Body   *WhatsAppInteractiveBody   `json:"body"`
	Action *WhatsAppInteractiveAction `json:"action"`
	Footer *WhatsAppInteractiveFooter `json:"footer,omitempty"`
	Reply  *WhatsAppInteractiveReply  `json:"reply,omitempty"`
}

// WhatsAppInteractiveType
// https://developers.messagebird.com/api/conversations/#whatsappinteractivetype-object
type WhatsAppInteractiveType struct {
	List        []string `json:"list"`
	Button      []string `json:"button"`
	Product     string   `json:"product"`
	ProductList []string `json:"product_list,omitempty"`
	ButtonReply []string `json:"button_reply,omitempty"`
}

// WhatsAppInteractiveHeader
// https://developers.messagebird.com/api/conversations/#whatsappinteractiveheader-object
type WhatsAppInteractiveHeader struct {
	Type     *WhatsAppInteractiveHeaderType `json:"type"`
	Text     string                         `json:"text"`
	Video    *Media                         `json:"video"`
	Image    *Media                         `json:"image"`
	Document *Media                         `json:"document"`
}

// WhatsAppInteractiveHeaderType
// https://developers.messagebird.com/api/conversations/#whatsappinteractiveheadertype-object
type WhatsAppInteractiveHeaderType struct {
	Text     string `json:"text,omitempty"`
	Video    string `json:"video,omitempty"`
	Image    string `json:"image,omitempty"`
	Document string `json:"document,omitempty"`
}

// WhatsAppInteractiveBody
// https://developers.messagebird.com/api/conversations/#whatsappinteractivebody-object
type WhatsAppInteractiveBody struct {
	Text string `json:"text"`
}

// WhatsAppInteractiveAction
// https://developers.messagebird.com/api/conversations/#whatsappinteractiveaction-object
type WhatsAppInteractiveAction struct {
	CatalogId         string                        `json:"catalog_id"`
	ProductRetailerId string                        `json:"product_retailer_id"`
	Sections          []*WhatsAppInteractiveSection `json:"sections"`
	Button            string                        `json:"button"`
	Buttons           *WhatsAppInteractiveButton    `json:"buttons"`
}

// WhatsAppInteractiveSection
// https://developers.messagebird.com/api/conversations/#whatsappinteractivesection-object
type WhatsAppInteractiveSection struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	ImageUrl string `json:"image_url,omitempty"`
}

// WhatsAppInteractiveButton
// https://developers.messagebird.com/api/conversations/#whatsappinteractivebutton-object
type WhatsAppInteractiveButton struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	ImageUrl string `json:"image_url,omitempty"`
}

// WhatsAppInteractiveFooter
// https://developers.messagebird.com/api/conversations/#whatsappinteractivefooter-object
type WhatsAppInteractiveFooter struct {
	Text string `json:"text"`
}

// WhatsAppInteractiveReply
// https://developers.messagebird.com/api/conversations/#whatsappinteractivereply-object
type WhatsAppInteractiveReply struct {
	Id          string `json:"id"`
	Text        string `json:"text"`
	Description string `json:"description,omitempty"`
}
