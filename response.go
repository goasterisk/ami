package ami

import (
	"errors"
	"net/textproto"
)

type Response struct {
	Status   string
	ActionID string
	fields   map[string]string
}

func newResponse(data *textproto.MIMEHeader) (*Response, error) {

	if "" == data.Get("Response") {
		return nil, errors.New("Not a valid response")
	}

	response := &Response{data.Get("Response"), data.Get("ActionID"), make(map[string]string)}
	data.Del("Response")
	data.Del("ActionID")

	for key, value := range *data {
		response.fields[key] = value[0]
	}

	return response, nil
}
