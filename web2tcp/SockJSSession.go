// SockJSSession.go
package web2tcp

import (
	//"log"

	"github.com/igm/sockjs-go/sockjs"
)

type sjsSession struct {
	SessionBase
	conn       sockjs.Session
	clientAddr string
	active     bool
}

func newSockJSSession(conn sockjs.Session, clientAddr string) *sjsSession {
	s := &sjsSession{
		SessionBase: newSessionBase("sockjs"),
		conn:        conn,
		clientAddr:  clientAddr,
		active:      true,
	}

	return s
}

func (s *sjsSession) Read() (string, error) {
	return s.conn.Recv()
}

func (s *sjsSession) Write(data string) error {
	return s.conn.Send(data)
}

func (s *sjsSession) Close() error {
	var err error = nil

	if s.active {
		err = s.conn.Close(0, "")
		s.active = false
	}

	return err
}

func (s *sjsSession) IsActive() bool {
	return s.active
}

func (s *sjsSession) RemoteAddr() string {
	return s.clientAddr
}

func (s *sjsSession) LocalAddr() string {
	//return s.conn.Request().
	return ""
}
