package req

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// create a default client
func newClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   2 * time.Minute,
	}
}

// Client return the default underlying http client
func (r *Req) Client() *http.Client {
	if r.client == nil {
		r.client = newClient()
	}
	return r.client
}

// Client return the default underlying http client
func Client() *http.Client {
	return std.Client()
}

// SetClient sets the underlying http.Client.
func (r *Req) SetClient(client *http.Client) {
	r.client = client // use default if client == nil
}

// SetClient sets the default http.Client for requests.
func SetClient(client *http.Client) {
	std.SetClient(client)
}

// SetFlags control display format of *Resp
func (r *Req) SetFlags(flags int) {
	r.flag = flags
}

// SetFlags control display format of *Resp
func SetFlags(flags int) {
	std.SetFlags(flags)
}

// Flags return output format for the *Resp
func (r *Req) Flags() int {
	return r.flag
}

// Flags return output format for the *Resp
func Flags() int {
	return std.Flags()
}

func (r *Req) getTransport() *http.Transport {
	trans, _ := r.Client().Transport.(*http.Transport)
	return trans
}

// EnableInsecureTLS allows insecure https
func (r *Req) EnableInsecureTLS(enable bool) {
	trans := r.getTransport()
	if trans == nil {
		return
	}
	if trans.TLSClientConfig == nil {
		trans.TLSClientConfig = &tls.Config{}
	}
	trans.TLSClientConfig.InsecureSkipVerify = enable
}

func EnableInsecureTLS(enable bool) {
	std.EnableInsecureTLS(enable)
}

// EnableCookieenable or disable cookie manager
func (r *Req) EnableCookie(enable bool) {
	if enable {
		jar, _ := cookiejar.New(nil)
		r.Client().Jar = jar
	} else {
		r.Client().Jar = nil
	}
}

// EnableCookieenable or disable cookie manager
func EnableCookie(enable bool) {
	std.EnableCookie(enable)
}

// SetTimeout sets the timeout for every request
func (r *Req) SetTimeout(d time.Duration) {
	r.Client().Timeout = d
}

// SetTimeout sets the timeout for every request
func SetTimeout(d time.Duration) {
	std.SetTimeout(d)
}

// SetProxyUrl set the simple proxy with fixed proxy url
func (r *Req) SetProxyUrl(rawurl string) error {
	trans := r.getTransport()
	if trans == nil {
		return errors.New("req: no transport")
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	trans.Proxy = http.ProxyURL(u)
	return nil
}

// SetProxyUrl set the simple proxy with fixed proxy url
func SetProxyUrl(rawurl string) error {
	return std.SetProxyUrl(rawurl)
}

// SetProxy sets the proxy for every request
func (r *Req) SetProxy(proxy func(*http.Request) (*url.URL, error)) error {
	trans := r.getTransport()
	if trans == nil {
		return errors.New("req: no transport")
	}
	trans.Proxy = proxy
	return nil
}

// SetProxy sets the proxy for every request
func SetProxy(proxy func(*http.Request) (*url.URL, error)) error {
	return std.SetProxy(proxy)
}

type jsonEncOpts struct {
	indentPrefix string
	indentValue  string
	escapeHTML   bool
}

func (r *Req) getJSONEncOpts() *jsonEncOpts {
	if r.jsonEncOpts == nil {
		r.jsonEncOpts = &jsonEncOpts{escapeHTML: true}
	}
	return r.jsonEncOpts
}

// SetJSONEscapeHTML specifies whether problematic HTML characters
// should be escaped inside JSON quoted strings.
// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
// to avoid certain safety problems that can arise when embedding JSON in HTML.
//
// In non-HTML settings where the escaping interferes with the readability
// of the output, SetEscapeHTML(false) disables this behavior.
func (r *Req) SetJSONEscapeHTML(escape bool) {
	opts := r.getJSONEncOpts()
	opts.escapeHTML = escape
}

// SetJSONEscapeHTML specifies whether problematic HTML characters
// should be escaped inside JSON quoted strings.
// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
// to avoid certain safety problems that can arise when embedding JSON in HTML.
//
// In non-HTML settings where the escaping interferes with the readability
// of the output, SetEscapeHTML(false) disables this behavior.
func SetJSONEscapeHTML(escape bool) {
	std.SetJSONEscapeHTML(escape)
}

// SetJSONIndent instructs the encoder to format each subsequent encoded
// value as if indented by the package-level function Indent(dst, src, prefix, indent).
// Calling SetIndent("", "") disables indentation.
func (r *Req) SetJSONIndent(prefix, indent string) {
	opts := r.getJSONEncOpts()
	opts.indentPrefix = prefix
	opts.indentValue = indent
}

// SetJSONIndent instructs the encoder to format each subsequent encoded
// value as if indented by the package-level function Indent(dst, src, prefix, indent).
// Calling SetIndent("", "") disables indentation.
func SetJSONIndent(prefix, indent string) {
	std.SetJSONIndent(prefix, indent)
}

type xmlEncOpts struct {
	prefix string
	indent string
}

func (r *Req) getXMLEncOpts() *xmlEncOpts {
	if r.xmlEncOpts == nil {
		r.xmlEncOpts = &xmlEncOpts{}
	}
	return r.xmlEncOpts
}

// SetXMLIndent sets the encoder to generate XML in which each element
// begins on a new indented line that starts with prefix and is followed by
// one or more copies of indent according to the nesting depth.
func (r *Req) SetXMLIndent(prefix, indent string) {
	opts := r.getXMLEncOpts()
	opts.prefix = prefix
	opts.indent = indent
}

// SetXMLIndent sets the encoder to generate XML in which each element
// begins on a new indented line that starts with prefix and is followed by
// one or more copies of indent according to the nesting depth.
func SetXMLIndent(prefix, indent string) {
	std.SetXMLIndent(prefix, indent)
}
