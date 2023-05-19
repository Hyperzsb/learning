package payment

type ChargeError struct {
	Msg string
}

func (ce *ChargeError) Error() string {
	return ce.Msg
}
