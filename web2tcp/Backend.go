// Backend.go
package web2tcp

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
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
	options   *BackendOptions
	tlsConfig *tls.Config
}

func NewBackend(options *BackendOptions) *Backend {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: options.Tls.SkipVerify,
	}

	if !options.Tls.SkipVerify {
		bytes, err := ioutil.ReadFile(options.Tls.CertAuthorityFile)

		if err != nil {
			log.Fatalln("Unable to read Certificate Authorities file.")
		}

		rootCAs := x509.NewCertPool()

		ok := rootCAs.AppendCertsFromPEM(bytes)

		if !ok {
			panic("Failed to parse Certificate Authorities from file.")
		}

		tlsConfig.RootCAs = rootCAs
	}

	backend := &Backend{
		options:   options,
		tlsConfig: tlsConfig,
	}

	return backend
}

func (b *Backend) NewProxySession(feSession Session) (Session, error) {
	var conn net.Conn
	var err error

	if b.options.Tls.Enabled {
		conn, err = tls.Dial("tcp4", b.options.Endpoint, b.tlsConfig)
	} else {
		conn, err = net.Dial("tcp4", b.options.Endpoint)
	}

	if err != nil {
		log.Println(err)
	}

	beSession := newTcpSession(conn)

	if b.options.ProxyProtocol {
		beSession.WriteProxyHeader(feSession)
	}

	return beSession, err
}
