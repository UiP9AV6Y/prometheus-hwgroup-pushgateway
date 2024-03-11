package handler

import (
	"encoding/xml"
	"errors"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

var errAddressMissing = errors.New("request is missing client address information")

func Push(p *portal.Portal, logger log.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := request.NewRoot()
		if err := xml.NewDecoder(r.Body).Decode(body); err != nil {
			level.Debug(logger).Log("msg", "failed to decode request", "err", err)
			http.Error(w, "Malformed request payload", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		var client netip.Addr
		var err error
		switch p.TrustLevel() {
		case portal.ProxyTrustLevel:
			client, err = netip.ParseAddr(ParseHeaderAddress(r))
		case portal.ClientTrustLevel:
			client, err = ParseBodyAddress(body)
		//case portal.ConnectionTrustLevel:
		default:
			client, err = netip.ParseAddr(ParseRemoteAddress(r))
		}

		if err != nil {
			level.Debug(logger).Log("msg", "failed to parse client address", "err", err)
			http.Error(w, "Invalid request information", http.StatusBadRequest)
			return
		}

		req := portal.NewRequest(r.Context(), &client)
		req.Credentials = ParseUserInfo(r)
		req.Payload = body
		rsp, err := p.Process(req)
		if err != nil {
			level.Error(logger).Log("msg", "failed to process request", "err", err)
			http.Error(w, "Failed to process request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		if err := xml.NewEncoder(w).Encode(rsp); err != nil {
			level.Error(logger).Log("msg", "failed to encode response", "err", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})
}

func ParseUserInfo(r *http.Request) *url.Userinfo {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil
	}

	return url.UserPassword(username, password)
}

func ParseBodyAddress(r *request.Root) (a netip.Addr, err error) {
	if r.Network != nil {
		a = *r.Network.IPAddr
	} else {
		err = errAddressMissing
	}

	return
}

func ParseRemoteAddress(r *http.Request) string {
	addr, _, _ := net.SplitHostPort(r.RemoteAddr)
	return addr
}

func ParseHeaderAddress(r *http.Request) string {
	addr := r.Header.Get("X-Real-Ip")
	if addr != "" {
		return addr
	}

	fwd := r.Header.Get("X-Forwarded-For")
	if fwd != "" {
		ips := strings.Split(fwd, ",")
		addr = strings.TrimSpace(ips[len(ips)-1])
	}

	return addr
}
