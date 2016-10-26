// Backend.go
package web2tcp

import (
	"crypto/tls"
	"log"
	"net"
)

type BackendOptions struct {
	Endpoint      string
	ProxyProtocol bool
	Tls           struct {
		Enabled           bool
		SkipVerify        bool
		CertAuthorityFile string
	}
}

type Backend struct {
	options *BackendOptions
}

func NewBackend(options *BackendOptions) *Backend {
	backend := &Backend{
		options: options,
	}

	return backend
}

func (b *Backend) NewSession() (Session, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: b.options.Tls.SkipVerify,
	}

	var conn net.Conn
	var err error

	if b.options.Tls.Enabled {
		conn, err = tls.Dial("tcp", b.options.Endpoint, tlsConfig)
	} else {
		conn, err = net.Dial("tcp", b.options.Endpoint)
	}

	if err != nil {
		log.Println(err)
	}

	return newTcpSession(conn), err
}
