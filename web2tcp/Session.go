// Configuration.go
package web2tcp

type Session interface {
	Type() string
	ID() string
	Read() (string, error)
	Write(data string) error
	Close() error
	IsActive() bool
}
