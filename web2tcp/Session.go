// Configuration.go
package web2tcp

type Session interface {
	Type() string
	ID() string
	Read() (string, error)
	Write(data string) error
	Close() error
	IsActive() bool
	RemoteAddr() string
	LocalAddr() string
}

type SessionBase struct {
	sid   string
	stype string
}

func newSessionBase(stype string) SessionBase {
	return SessionBase{
		sid:   newUUID(),
		stype: stype,
	}
}

func (s SessionBase) ID() string {
	return s.sid
}

func (s SessionBase) Type() string {
	return s.stype
}
