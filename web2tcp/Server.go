// Server.go
package web2tcp

import (
	"log"
)

type Web2TCPServer struct {
	frontend *Frontend
	backend  *Backend
}

type Configuration struct {
	Frontend FrontendOptions
	Backend  BackendOptions
}

func NewServer(conf Configuration) *Web2TCPServer {
	frontend := NewFrontend(&conf.Frontend)
	backend := NewBackend(&conf.Backend)

	setupTunnel(conf.Frontend.Route, frontend, backend)

	server := &Web2TCPServer{
		frontend: frontend,
		backend:  backend,
	}

	return server
}

func setupTunnel(route string, frontend *Frontend, backend *Backend) {
	frontend.RouteHandler(route, func(feSession Session) {
		//log.Println("New Session, Type : ", feSession.Type())

		beSession, err := backend.NewSession()

		if err == nil {
			go pipeHandler(feSession, beSession)
			go pipeHandler(beSession, feSession)

			feSession.Write("STATUS:CONNECTED")
		}
	})
}

func pipeHandler(source, target Session) {
	defer closeSessions(source, target)
	var msg string
	var err error

	for {
		if msg, err = source.Read(); err == nil && len(msg) > 0 {
			if err = target.Write(msg); err == nil {
				continue
			} else {
				log.Println(source.Type(), "->", target.Type(), " : Writing Error")
				break
			}
		} else {
			log.Println(source.Type(), "->", target.Type(), " : Reading Error")
			break
		}
	}
}

func closeSessions(source, target Session) {
	if source.IsActive() {
		log.Println(source.Type(), " closed")
		source.Close()
	}

	if target.IsActive() {
		log.Println(target.Type(), " closed")
		target.Close()
	}
}

func (s *Web2TCPServer) Start() {
	s.frontend.Listen()
}
