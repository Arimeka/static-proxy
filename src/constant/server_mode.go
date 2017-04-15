package constant

import "fmt"

type ServerMode int

const (
	Development ServerMode = iota
	Staging
	Production
)

var serverMods = [3]string{"development", "staging", "production"}

func (m ServerMode) String() string {
	if m < 0 || m+1 > ServerMode(len(serverMods)) {
		return fmt.Sprintf("ServerMode(%d)", m)
	}
	return serverMods[m]
}

func ParseServerMode(s string) (ServerMode, error) {
	for i, v := range serverMods {
		if v == s {
			return ServerMode(i), nil
		}
	}
	return Development, fmt.Errorf("Server mode %s not included in allowed array %v", s, serverMods)
}
