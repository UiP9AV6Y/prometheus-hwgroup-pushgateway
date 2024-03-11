package portal

import (
	"context"
	"crypto/subtle"
	"net/netip"
	"net/url"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

type Request struct {
	Credentials *url.Userinfo
	Payload     *request.Root

	client *netip.Addr
	ctx    context.Context
}

func NewRequest(ctx context.Context, client *netip.Addr) *Request {
	result := &Request{
		ctx:    ctx,
		client: client,
	}

	return result
}

func (r *Request) Client() *netip.Addr {
	return r.client
}

func (r *Request) Context() context.Context {
	return r.ctx
}

func (r *Request) CompareCredentials(c *url.Userinfo) bool {
	if c == nil {
		return true
	} else if r.Credentials == nil {
		return false
	}

	au := r.Credentials.Username()
	ap, _ := r.Credentials.Password()
	gu := c.Username()
	gp, _ := c.Password()

	return SecureCompare(gu, au) && SecureCompare(gp, ap)
}

// SecureCompare performs a constant time compare of two strings to limit timing attacks.
// https://go.dev/play/p/NU5uTaB-sp
func SecureCompare(want, have string) bool {
	if subtle.ConstantTimeEq(int32(len(want)), int32(len(have))) == 1 {
		return subtle.ConstantTimeCompare([]byte(want), []byte(have)) == 1
	}

	/* Securely compare have to itself to keep constant time, but always return false */
	return subtle.ConstantTimeCompare([]byte(have), []byte(have)) == 1 && false
}
