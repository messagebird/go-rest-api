package messagebird

type Formats struct {
	E164          string
	International string
	National      string
	Rfc3966       string
}

type Lookup struct {
	Href          string
	CountryCode   string
	CountryPrefix int
	PhoneNumber   int
	Type          string
	Formats       *Formats
	HLR           *HLR
}
