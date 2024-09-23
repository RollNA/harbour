package httpUtil

type ResponseDto struct {
	Body []byte
	Err  error
}

func (r *ResponseDto) Result() ([]byte, error) {
	return r.Body, r.Err
}
