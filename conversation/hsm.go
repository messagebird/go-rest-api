package conversation

import "time"

// HSM is a pre-approved, reusable message template required when messaging
// over WhatsApp. It allows you to just send the required parameter values
// instead of the full message. It also allows for localization of the message
// and decreases the possibility of being blocked on the first contact as the
// message is pre-approved by WhatsApp.
type HSM struct {
	Namespace             string                     `json:"namespace"`
	TemplateName          string                     `json:"templateName"`
	Language              *HSMLanguage               `json:"language"`
	LocalizableParameters []*HSMLocalizableParameter `json:"params"`
}

// HSMLanguage is used to set the message's locale.
type HSMLanguage struct {
	Policy HSMLanguagePolicy `json:"policy"`

	// Code can be both language and language_locale formats (e.g. en and
	// en_US).
	Code string `json:"code"`
}

// HSMLanguagePolicy sets how the provided language is enforced.
type HSMLanguagePolicy string

const (
	// HSMLanguagePolicyFallback will deliver the message template in the
	// user's device language. If the settings can't be found on the user's
	// device the fallback language is used.
	HSMLanguagePolicyFallback HSMLanguagePolicy = "fallback"

	// HSMLanguagePolicyDeterministic will deliver the message template
	// exactly in the language and locale asked for.
	HSMLanguagePolicyDeterministic HSMLanguagePolicy = "deterministic"
)

// HSMLocalizableParameter are used to replace the placeholders in the message
// template. They will be localized by WhatsApp. Default values are used when
// localization fails. Default is required. Additionally, currency OR DateTime
// may be present in a request.
type HSMLocalizableParameter struct {
	Default  string                           `json:"default"`
	Currency *HSMLocalizableParameterCurrency `json:"currency,omitempty"`
	DateTime *time.Time                       `json:"dateTime,omitempty"`
}

type HSMLocalizableParameterCurrency struct {
	// Code is the currency code in ISO 4217 format.
	Code string `json:"currencyCode"`

	// Amount is the total amount, including cents, multiplied by 1000. E.g.
	// 12.34 become 12340.
	Amount int64 `json:"amount"`
}

// DefaultLocalizableHSMParameter gets a simple parameter with a default value
// that will do a simple string replacement.
func DefaultLocalizableHSMParameter(d string) *HSMLocalizableParameter {
	return &HSMLocalizableParameter{
		Default: d,
	}
}

// CurrencyLocalizableHSMParameter gets a parameter that localizes a currency.
// Code is the currency code in ISO 4217 format and amount is the total amount,
// including cents, multiplied by 1000. E.g. 12.34 becomes 12340.
func CurrencyLocalizableHSMParameter(d string, code string, amount int64) *HSMLocalizableParameter {
	return &HSMLocalizableParameter{
		Default: d,
		Currency: &HSMLocalizableParameterCurrency{
			Code:   code,
			Amount: amount,
		},
	}
}

// DateTimeLocalizableHSMParameter gets a parameter that localizes a DateTime.
func DateTimeLocalizableHSMParameter(d string, dateTime time.Time) *HSMLocalizableParameter {
	return &HSMLocalizableParameter{
		Default:  d,
		DateTime: &dateTime,
	}
}
