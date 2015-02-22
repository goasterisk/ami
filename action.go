package ami

import (
	"bytes"
	"fmt"
)

type Action struct {
	Action   string
	ActionID string
	Params   map[string]string
}

func NewAction(action string, params map[string]string) *Action {
	if "" == params["ActionID"] {
		params["ActionID"] = "nmartin"
	}
	return &Action{action, params["ActionID"], params}
}

func (a *Action) String() string {
	var request bytes.Buffer
	request.WriteString(fmt.Sprintf("Action: %s\n", a.Action))

	// Process each action param
	for key, value := range a.Params {
		request.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	return request.String()
}
