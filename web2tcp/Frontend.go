// Frontend.go
package web2tcp

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/igm/sockjs-go/sockjs"
)

type FrontendOptions struct {
	Endpoint string
	Route    string
	Tls      struct {
		Enabled     bool
		CertFile    string
		CertKeyFile string
	}
}

type Frontend struct {
	options *FrontendOptions
	router  *mux.Router
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var sockjsOptions = sockjs.DefaultOptions

func NewFrontend(options *FrontendOptions) *Frontend {
	frontend := &Frontend{
		options: options,
		router:  mux.NewRouter(),
	}

	return frontend
}

func (f Frontend) Listen() {
	var err error
	http.Handle("/", f.router)

	if f.options.Tls.Enabled {
		err = http.ListenAndServeTLS(f.options.Endpoint, f.options.Tls.CertFile, f.options.Tls.CertKeyFile, nil)
	} else {
		err = http.ListenAndServe(f.options.Endpoint, nil)
	}
	if err != nil {
		log.Panic(err)
	}
}

func (f *Frontend) RouteHandler(route string, handler func(Session)) {
	log.Printf("Registering SockJS Route [\\%s]\n", route)
	log.Printf("Registering WebSocket Route [\\%s\\websocket]\n", route)

	f.router.PathPrefix("/" + route + "/websocket").HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		wsConn, _ := wsUpgrader.Upgrade(resp, req, nil)
		xfHeader := req.Header.Get("X-Forwarded-For")
		var clientAddr string

		if len(xfHeader) > 0 {
			clientAddr = strings.Split(xfHeader, ",")[0]
		} else {
			clientAddr = wsConn.RemoteAddr().String()
		}

		go handler(newWSSession(wsConn, clientAddr))
	})

	sockjsHandler := sockjs.NewHandler("/"+route, sockjsOptions, func(session sockjs.Session) {
		var clientAddr string

		xfHeader := session.Request().Header.Get("X-Forwarded-For")

		if len(xfHeader) > 0 {
			clientAddr = strings.Split(xfHeader, ",")[0]
		} else {
			clientAddr = session.Request().RemoteAddr
		}

		go handler(newSockJSSession(session, clientAddr))
	})

	f.router.PathPrefix("/" + route + "/").Handler(sockjsHandler)
}
