package portal

import (
	"fmt"
)

type TrustLevel string

func (t TrustLevel) String() string {
	return string(t)
}

// Set updates the value of the trust level.
func (t *TrustLevel) Set(s string) error {
	switch s {
	case "proxy":
		(*t) = ProxyTrustLevel
	case "client":
		(*t) = ClientTrustLevel
	case "connection":
		(*t) = ConnectionTrustLevel
	default:
		return fmt.Errorf("unrecognized trust level %q", s)
	}

	return nil
}

const (
	ProxyTrustLevel      TrustLevel = "proxy"
	ClientTrustLevel     TrustLevel = "client"
	ConnectionTrustLevel TrustLevel = "connection"
)

var TrustLevelValues = []string{"proxy", "client", "connection"}
