package pkg

type MethodError struct {
	Code          int64         `json:"error_code"`
	Message       string        `json:"error_msg"`
	RequestParams []interface{} `json:"request_params"`
}

func (err *MethodError) Error() string {
	return err.Message
}
