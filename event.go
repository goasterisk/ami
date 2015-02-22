package ami

import (
	"errors"
	"fmt"
	"net/textproto"
)

type Event struct {
	Event    string
	ActionID string
	fields   map[string]string
}

func newEvent(data *textproto.MIMEHeader) (*Event, error) {

	if "" == data.Get("Event") {
		return nil, errors.New("Not a valid event")
	}

	event := &Event{data.Get("Event"), data.Get("ActionID"), make(map[string]string)}
	data.Del("Event")
	data.Del("ActionID")
	data.Del("Privilege")

	for key, value := range *data {
		event.fields[key] = value[0]
	}

	return event, nil
}

func (e *Event) String() string {
	return fmt.Sprintf("[Event] %s - Fields: %s", e.Event, e.fields)
}
