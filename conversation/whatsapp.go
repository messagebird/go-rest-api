package conversation

// WhatsAppInteractiveType
// https://developers.messagebird.com/api/conversations/#whatsappinteractivetype-object
type WhatsAppInteractiveType string

const (
	WAITypeList        WhatsAppInteractiveType = "list"
	WAITypeButton      WhatsAppInteractiveType = "button"
	WAITypeProduct     WhatsAppInteractiveType = "product"
	WAITypeProductList WhatsAppInteractiveType = "product_list"
	WAITypeButtonReply WhatsAppInteractiveType = "button_reply"
)

// WhatsAppInteractive
// https://developers.messagebird.com/api/conversations/#whatsappinteractive-object
type WhatsAppInteractive struct {
	Type   WhatsAppInteractiveType    `json:"type"`
	Header *WhatsAppInteractiveHeader `json:"header"`
	Body   *WhatsAppInteractiveBody   `json:"body"`
	Action *WhatsAppInteractiveAction `json:"action"`
	Footer *WhatsAppInteractiveFooter `json:"footer,omitempty"`
	Reply  *WhatsAppInteractiveReply  `json:"reply,omitempty"`
}

// WhatsAppInteractiveHeader
// https://developers.messagebird.com/api/conversations/#whatsappinteractiveheader-object
type WhatsAppInteractiveHeader struct {
	Type     WhatsAppInteractiveHeaderType `json:"type"`
	Text     string                        `json:"text"`
	Video    *Media                        `json:"video"`
	Image    *Media                        `json:"image"`
	Document *Media                        `json:"document"`
}

// WhatsAppInteractiveHeaderType
// https://developers.messagebird.com/api/conversations/#whatsappinteractiveheadertype-object
type WhatsAppInteractiveHeaderType string

const (
	WAIHeaderTypeText     WhatsAppInteractiveHeaderType = "text"
	WAIHeaderTypeVideo    WhatsAppInteractiveHeaderType = "video"
	WAIHeaderTypeImage    WhatsAppInteractiveHeaderType = "image"
	WAIHeaderTypeDocument WhatsAppInteractiveHeaderType = "document"
)

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
	Title        string                           `json:"title"`
	Rows         []*WhatsAppInteractiveSectionRow `json:"rows"`
	ProductItems []*WhatsAppInteractiveProduct    `json:"product_items"`
}

type WhatsAppInteractiveSectionRow struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type WhatsAppInteractiveProduct struct {
	ProductRetailerId string `json:"product_retailer_id"`
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

// WhatsAppSticker
// URL of the sticker image. The format must be image/webp and the maximum size is 100 KB.
type WhatsAppSticker struct {
	Link string `json:"link"`
}

// WhatsAppOrder
// https://developers.messagebird.com/api/conversations/#whatsapporder-object
type WhatsAppOrder struct {
	CatalogId    string                  `json:"catalog_id"`
	ProductItems []*WhatsAppOrderProduct `json:"product_items"`
	Text         string                  `json:"text"`
}

// WhatsAppOrderProduct
// https://developers.messagebird.com/api/conversations/#whatsapporderproduct-object
type WhatsAppOrderProduct struct {
	ProductRetailerId string `json:"product_retailer_id"`
	Quantity          int    `json:"quantity"`
	ItemPrice         string `json:"item_price"`
	Currency          string `json:"currency"`
}

// WhatsAppText
// https://developers.messagebird.com/api/conversations/#whatsapptext-object
type WhatsAppText struct {
	Text    *WhatsAppTextBody `json:"text"`
	Context *WhatsAppContext  `json:"context"`
}

// WhatsAppTextBody
// https://developers.messagebird.com/api/conversations/#whatsapptextbody-object
type WhatsAppTextBody struct {
	Body string `json:"body"`
}

// WhatsAppContext
// https://developers.messagebird.com/api/conversations/#whatsappcontext-object
type WhatsAppContext struct {
	Id              string                   `json:"id"`
	From            string                   `json:"from"`
	ReferredProduct *WhatsAppReferredProduct `json:"referred_product"`
}

// WhatsAppReferredProduct
// https://developers.messagebird.com/api/conversations/#whatsappreferredproduct-object
type WhatsAppReferredProduct struct {
	CatalogId         string `json:"catalog_id"`
	ProductRetailerId string `json:"product_retailer_id"`
}
