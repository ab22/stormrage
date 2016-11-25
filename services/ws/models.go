package ws

import "net"

type requestOption int

const (
	START_PING requestOption = iota
	STOP_PING
)

type request struct {
	Option requestOption `json:"option"`
	IP     string        `json:"ip"`
}

func (r *request) IsValidIP() bool {
	ip := net.ParseIP(r.IP)

	return ip != nil
}
