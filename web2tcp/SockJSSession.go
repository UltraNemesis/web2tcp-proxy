// SockJSSession.go
package web2tcp

import (
	//"log"

	"github.com/igm/sockjs-go/sockjs"
)

type sjsSession struct {
	conn   sockjs.Session
	active bool
}

func newSockJSSession(conn sockjs.Session) *sjsSession {
	s := &sjsSession{
		conn:   conn,
		active: true,
	}

	return s
}

func (s *sjsSession) Type() string {
	return "sockjs"
}

func (s *sjsSession) ID() string {
	return ""
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
