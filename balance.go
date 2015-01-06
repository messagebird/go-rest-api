package messagebird

type Balance struct {
	Payment string
	Type    string
	Amount  int
	Errors  []Error
}
