// tcpSession.go
package web2tcp

import (
	"bufio"
	"encoding/base64"
	"net"
)

type tcpSession struct {
	SessionBase
	conn    net.Conn
	bufConn *bufio.Reader
	buf     []byte
	active  bool
}

func newTcpSession(conn net.Conn) *tcpSession {
	s := &tcpSession{
		SessionBase: newSessionBase("tcp"),
		conn:        conn,
		bufConn:     bufio.NewReader(conn),
		buf:         make([]byte, 3072),
		active:      true,
	}

	return s
}

func (s *tcpSession) Read() (string, error) {
	var data string
	var count int
	var err error

	if count, err = s.bufConn.Read(s.buf); err == nil {
		data = base64.StdEncoding.EncodeToString(s.buf[0:count])
	}

	return data, err
}

func (s *tcpSession) Write(data string) error {
	var bytes []byte
	var err error

	if bytes, err = base64.StdEncoding.DecodeString(data); err == nil {
		_, err = s.conn.Write(bytes)
	}

	return err
}

func (s *tcpSession) Close() error {
	var err error

	if s.active {
		err = s.conn.Close()
		s.active = false
	}

	return err
}

func (s *tcpSession) IsActive() bool {
	return s.active
}
