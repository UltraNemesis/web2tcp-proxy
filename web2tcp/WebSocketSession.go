// WebSocketSession.go
package web2tcp

import (
	//"log"

	"sync"

	"github.com/gorilla/websocket"
)

type wsSession struct {
	SessionBase
	conn       *websocket.Conn
	clientAddr string
	readLock   sync.Mutex
	writeLock  sync.Mutex
	active     bool
}

func newWSSession(conn *websocket.Conn, clientAddr string) *wsSession {
	s := &wsSession{
		SessionBase: newSessionBase("ws"),
		conn:        conn,
		clientAddr:  clientAddr,
		active:      true,
	}

	return s
}

func (s *wsSession) Read() (string, error) {
	s.readLock.Lock()
	defer s.readLock.Unlock()
	_, msg, readErr := s.conn.ReadMessage()

	return string(msg), readErr
}

func (s *wsSession) Write(data string) error {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()
	return s.conn.WriteMessage(websocket.TextMessage, []byte(data))
}

func (s *wsSession) Close() error {
	var err error = nil

	if s.active {
		err = s.conn.Close()
		s.active = false
	}

	return err
}

func (s *wsSession) IsActive() bool {
	return s.active
}

func (s *wsSession) RemoteAddr() string {
	return s.clientAddr
}

func (s *wsSession) LocalAddr() string {
	return s.conn.LocalAddr().String()
}
