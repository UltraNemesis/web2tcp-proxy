// tcpSession.go
package web2tcp

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
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

func (s *tcpSession) WriteProxyHeader(feSession Session) {
	srcIP, srcPort, proto := parseNetAddr(feSession.RemoteAddr())
	dstIP, dstPort, _ := parseNetAddr(feSession.LocalAddr())

	_, err := fmt.Fprintf(s.conn, "PROXY %s %s %s %d %d\r\n", proto, srcIP, dstIP, srcPort, dstPort)

	if err != nil {
		log.Println(err)
	}
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

func (s *tcpSession) RemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *tcpSession) LocalAddr() string {
	return s.conn.LocalAddr().String()
}

func parseNetAddr(addrStr string) (string, int, string) {
	addr, err := net.ResolveTCPAddr("tcp", addrStr)

	var protocol string = "UNKNOWN"
	var ip string = "0.0.0.0"
	var port int = 0

	if err != nil {
		fmt.Println(err)

		return ip, port, protocol
	}

	protocol = "TCP6"
	port = addr.Port

	if addr.IP.To4() != nil {
		protocol = "TCP4"
		ip = addr.IP.To4().String()
	} else {
		ip = addr.IP.String()
	}

	log.Println("IP=", ip, "Port=", port, "Protocol=", protocol)

	return ip, port, protocol
}
