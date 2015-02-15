package ami

import (
	"bytes"
	"fmt"
)

type Action struct {
	action string
	id     string
	params map[string]string
}

func NewAction(action string, params map[string]string) Action {
	return Action{action, "", params}
}

func (a Action) Build() string {
	var request bytes.Buffer
	request.WriteString(fmt.Sprintf("Action: %s\n", a.action))

	// Process each action param
	for key, value := range a.params {
		request.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	request.WriteString("\r\n")

	return request.String()
}
