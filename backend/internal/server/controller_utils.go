package server

import (
	"net"
	"net/http"
	"net/netip"
)

func ipFromReq(req *http.Request) (*netip.Addr, error) {
	ipAddrString, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, err
	}

	if value := req.Header.Get("X-Forwarded-For"); value != "" {
		ipAddrString = value
	}

	ipAddr, err := netip.ParseAddr(ipAddrString)
	if err != nil {
		return nil, err
	}

	return &ipAddr, nil
}
