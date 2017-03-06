package messagebird

// Balance object represents your balance at MessageBird.com
type Balance struct {
	Payment string
	Type    string
	Amount  float32
	Errors  []Error
}
