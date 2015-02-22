package ami

import (
	"log"
	"net"
	"net/textproto"
	"strings"
)

type Manager struct {
	hostname string
	port     string
	username string
	secret   string

	conn *textproto.Conn

	Version string

	eventListeners    map[string]chan *Event
	responseListeners map[string]chan *Response
}

func NewManager(hostname string, port string, username string, secret string) *Manager {
	var mngr = Manager{hostname,
		port,
		username,
		secret,
		nil,
		"",
		make(map[string]chan *Event),
		make(map[string]chan *Response),
	}
	return &mngr
}

func (m *Manager) Connect() (err error) {
	m.conn, err = textproto.Dial("tcp", net.JoinHostPort(m.hostname, m.port))

	if err != nil {
		return
	}
	log.Print("Connection success")

	// Get AMI connection welcome message
	welcome, err := m.conn.ReadLine()
	if nil != err {
		return
	}
	m.Version = strings.Split(welcome, "/")[1]

	// Run main loop
	m.run()

	// Send login action
	var params = map[string]string{
		"Username": m.username,
		"Secret":   m.secret,
	}
	_, err = m.Execute(NewAction("Login", params))

	return err
}

func (m *Manager) Execute(a *Action) (*Response, error) {

	responseChan, err := m.ExcecuteAsync(a)
	if nil != err {
		return nil, err
	}

	// Handle response
	response := <-responseChan
	log.Printf("Response: %s", response)

	return response, nil
}

func (m *Manager) ExcecuteAsync(a *Action) (<-chan *Response, error) {

	log.Printf("%s", a)
	err := m.conn.Writer.PrintfLine("%s", a)
	if err != nil {
		return nil, err
	}

	m.responseListeners[a.ActionID] = make(chan *Response, 1)

	return m.responseListeners[a.ActionID], nil
}

func (m *Manager) AddEventListener(id string, c chan *Event) {
	m.eventListeners[id] = c
}

func (m *Manager) RemoveEventListener(id string) {
	delete(m.eventListeners, id)
}

func (m *Manager) Disconnect() (err error) {

	_, err = m.Execute(NewAction("Logoff", nil))
	err = m.conn.Close()
	return
}

func (m *Manager) run() {
	go func() {
		log.Print("[Loop] Launch...")
		for {
			message, err := m.conn.ReadMIMEHeader()
			if err != nil {
				break
			}

			// Is message an event
			if event, err := newEvent(&message); nil == err {
				for _, channel := range m.eventListeners {
					channel <- event
				}
				continue
			}

			// Is message a response to a previous action
			if response, err := newResponse(&message); nil == err {
				m.responseListeners[response.ActionID] <- response
				continue
			}

			log.Print("[Loop] Message is for nobody")
		}
	}()
}
