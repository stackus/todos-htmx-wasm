package domain

type ErrMarshaling struct {
	Err error
}

func (e ErrMarshaling) Error() string {
	return "failed to marshal request: " + e.Err.Error()
}

type ErrUnmarshaling struct {
	Err error
}

func (e ErrUnmarshaling) Error() string {
	return "failed to unmarshal response: " + e.Err.Error()
}

type ErrCreateRequest struct {
	Err error
}

func (e ErrCreateRequest) Error() string {
	return "failed to create request: " + e.Err.Error()
}

type ErrMakeRequest struct {
	Err error
}

func (e ErrMakeRequest) Error() string {
	return "failed to make request: " + e.Err.Error()
}
