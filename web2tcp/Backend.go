// Backend.go
package web2tcp

import (
	"crypto/tls"
	"log"
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

	conn, err := tls.Dial("tcp", b.options.Endpoint, tlsConfig)

	if err != nil {
		log.Println(err)
	}

	return newTcpSession(conn), err
}
