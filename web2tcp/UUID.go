// UUID.go
package web2tcp

import (
	"github.com/nats-io/nuid"
)

func newUUID() string {
	return nuid.Next()
}
