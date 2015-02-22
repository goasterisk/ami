package ami

import (
	"log"
	"net"
	"net/textproto"
)

type Manager struct {
	hostname string
	port     string
	username string
	secret   string

	conn *textproto.Conn
}

func NewManager(hostname string, port string, username string, secret string) *Manager {
	var mngr = Manager{hostname, port, username, secret, nil}
	return &mngr
}

func (m *Manager) Connect() (err error) {
	m.conn, err = textproto.Dial("tcp", net.JoinHostPort(m.hostname, m.port))

	if err != nil {
		return
	}
	log.Printf("Connection successed\n")

	// Get AMI connection header + version
	header, err := m.conn.ReadLine()
	log.Printf("Header: %s", header)

	// Send login action
	var params = map[string]string{
		"Username": m.username,
		"Secret":   m.secret,
	}
	_, err = m.Execute(NewAction("Login", params))

	return
}

func (m *Manager) Execute(action *Action) (response string, err error) {

	log.Printf("%s", action)
	err = m.conn.Writer.PrintfLine("%s", action)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	// Handle response
	headers, err := m.conn.ReadMIMEHeader()
	//response, err = m.readMessage()
	log.Printf("%s", headers)

	return
}

func (m *Manager) ListenForEvents(c chan *Event) {
	go func() {
		for {
			response, err := m.conn.ReadMIMEHeader()
			if err != nil {
				break
			}
			event, err := newEvent(&response)
			c <- event
		}
	}()
}

func (m *Manager) Disconnect() (err error) {

	_, err = m.Execute(NewAction("Logoff", nil))
	err = m.conn.Close()
	return
}
