package ami

type ParsingError struct {
	Err error
}

func (e *ParsingError) Error() string { return e.Err.Error() }
