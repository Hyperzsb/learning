package model

type EmptyQueryError struct {
	Msg string
}

func (eqe *EmptyQueryError) Error() string {
	return eqe.Msg
}
