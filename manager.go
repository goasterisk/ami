package ami

import (
	"io"
	"log"
	"net"
	"strings"
)

type Manager struct {
	hostname   string
	port       string
	username   string
	secret     string
	connection net.Conn
}

func NewManager(hostname string, port string, username string, secret string) *Manager {
	var mngr = Manager{hostname, port, username, secret, nil}
	return &mngr
}

func (m *Manager) Connect() (err error) {
	m.connection, err = net.Dial("tcp", net.JoinHostPort(m.hostname, m.port))

	if err != nil {
		return
	} else {
		log.Printf("Connection successed\n")
	}

	// Get AMI connection headers
	header, err := m.read()
	if err != nil {
		return
	}
	log.Printf("Headers: %v", header)

	// Send login action
	var params = map[string]string{
		"ActionID": "monid",
		"Username": m.username,
		"Secret":   m.secret,
	}
	m.Execute(NewAction("Login", params))

	return
}

func (m *Manager) Execute(action Action) (response string, err error) {

	request := action.Build()
	log.Printf("%s", request)
	_, err = m.connection.Write([]byte(request))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	// Handle response
	response, err = m.read()
	log.Printf("%v", response)

	return
}

func (m *Manager) Disconnect() (err error) {

	_, err = m.Execute(NewAction("Logoff", nil))
	err = m.connection.Close()
	return
}

func (m *Manager) read() (response string, err error) {

	for {
		rawPart := make([]byte, 500)
		bytesRead, err := m.connection.Read(rawPart)
		if bytesRead == 0 && err == io.EOF {
			break
		}

		part := string(rawPart[:bytesRead])
		response += part
		if strings.HasSuffix(part, "\r\n") || err != nil {
			break
		}
	}

	return
}
