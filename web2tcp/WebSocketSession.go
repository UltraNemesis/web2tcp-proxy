// WebSocketSession.go
package web2tcp

import (
	//"log"

	"sync"

	"github.com/gorilla/websocket"
)

type wsSession struct {
	conn      *websocket.Conn
	readLock  sync.Mutex
	writeLock sync.Mutex
	active    bool
}

func newWSSession(conn *websocket.Conn) *wsSession {
	s := &wsSession{
		conn:   conn,
		active: true,
	}

	return s
}

func (s *wsSession) Type() string {
	return "ws"
}

func (s *wsSession) ID() string {
	return ""
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
