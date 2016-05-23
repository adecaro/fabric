package fabric

type errorHandlerImpl struct {
	errors []error
}

func (eh *errorHandlerImpl) pushError(err error) {
	eh.errors = append(eh.errors, err)
}

func (eh *errorHandlerImpl) Flush() error {
	l := len(eh.errors)
	if l == 0 {
		return nil
	}
	res := eh.errors[l-1]
	eh.errors = eh.errors[:l-1]
	return res
}
