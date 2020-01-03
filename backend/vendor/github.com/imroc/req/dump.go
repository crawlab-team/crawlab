package req

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Debug enable debug mode if set to true
var Debug bool

// dumpConn is a net.Conn which writes to Writer and reads from Reader
type dumpConn struct {
	io.Writer
	io.Reader
}

func (c *dumpConn) Close() error                       { return nil }
func (c *dumpConn) LocalAddr() net.Addr                { return nil }
func (c *dumpConn) RemoteAddr() net.Addr               { return nil }
func (c *dumpConn) SetDeadline(t time.Time) error      { return nil }
func (c *dumpConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *dumpConn) SetWriteDeadline(t time.Time) error { return nil }

// delegateReader is a reader that delegates to another reader,
// once it arrives on a channel.
type delegateReader struct {
	c chan io.Reader
	r io.Reader // nil until received from c
}

func (r *delegateReader) Read(p []byte) (int, error) {
	if r.r == nil {
		r.r = <-r.c
	}
	return r.r.Read(p)
}

type dummyBody struct {
	N   int
	off int
}

func (d *dummyBody) Read(p []byte) (n int, err error) {
	if d.N <= 0 {
		err = io.EOF
		return
	}
	left := d.N - d.off
	if left <= 0 {
		err = io.EOF
		return
	}

	if l := len(p); l > 0 {
		if l >= left {
			n = left
			err = io.EOF
		} else {
			n = l
		}
		d.off += n
		for i := 0; i < n; i++ {
			p[i] = '*'
		}
	}

	return
}

func (d *dummyBody) Close() error {
	return nil
}

type dumpBuffer struct {
	bytes.Buffer
}

func (b *dumpBuffer) Write(p []byte) {
	if b.Len() > 0 {
		b.Buffer.WriteString("\r\n\r\n")
	}
	b.Buffer.Write(p)
}

func (b *dumpBuffer) WriteString(s string) {
	b.Write([]byte(s))
}

func (r *Resp) dumpRequest(dump *dumpBuffer) {
	head := r.r.flag&LreqHead != 0
	body := r.r.flag&LreqBody != 0

	if head {
		r.dumpReqHead(dump)
	}
	if body {
		if r.multipartHelper != nil {
			dump.Write(r.multipartHelper.Dump())
		} else if len(r.reqBody) > 0 {
			dump.Write(r.reqBody)
		}
	}
}

func (r *Resp) dumpReqHead(dump *dumpBuffer) {
	reqSend := new(http.Request)
	*reqSend = *r.req
	if reqSend.URL.Scheme == "https" {
		reqSend.URL = new(url.URL)
		*reqSend.URL = *r.req.URL
		reqSend.URL.Scheme = "http"
	}

	if reqSend.ContentLength > 0 {
		reqSend.Body = &dummyBody{N: int(reqSend.ContentLength)}
	} else {
		reqSend.Body = &dummyBody{N: 1}
	}

	// Use the actual Transport code to record what we would send
	// on the wire, but not using TCP.  Use a Transport with a
	// custom dialer that returns a fake net.Conn that waits
	// for the full input (and recording it), and then responds
	// with a dummy response.
	var buf bytes.Buffer // records the output
	pr, pw := io.Pipe()
	defer pw.Close()
	dr := &delegateReader{c: make(chan io.Reader)}

	t := &http.Transport{
		Dial: func(net, addr string) (net.Conn, error) {
			return &dumpConn{io.MultiWriter(&buf, pw), dr}, nil
		},
	}
	defer t.CloseIdleConnections()

	client := new(http.Client)
	*client = *r.client
	client.Transport = t

	// Wait for the request before replying with a dummy response:
	go func() {
		req, err := http.ReadRequest(bufio.NewReader(pr))
		if err == nil {
			// Ensure all the body is read; otherwise
			// we'll get a partial dump.
			io.Copy(ioutil.Discard, req.Body)
			req.Body.Close()
		}

		dr.c <- strings.NewReader("HTTP/1.1 204 No Content\r\nConnection: close\r\n\r\n")
		pr.Close()
	}()

	_, err := client.Do(reqSend)
	if err != nil {
		dump.WriteString(err.Error())
	} else {
		reqDump := buf.Bytes()
		if i := bytes.Index(reqDump, []byte("\r\n\r\n")); i >= 0 {
			reqDump = reqDump[:i]
		}
		dump.Write(reqDump)
	}
}

func (r *Resp) dumpResponse(dump *dumpBuffer) {
	head := r.r.flag&LrespHead != 0
	body := r.r.flag&LrespBody != 0
	if head {
		respDump, err := httputil.DumpResponse(r.resp, false)
		if err != nil {
			dump.WriteString(err.Error())
		} else {
			if i := bytes.Index(respDump, []byte("\r\n\r\n")); i >= 0 {
				respDump = respDump[:i]
			}
			dump.Write(respDump)
		}
	}
	if body && len(r.Bytes()) > 0 {
		dump.Write(r.Bytes())
	}
}

// Cost return the time cost of the request
func (r *Resp) Cost() time.Duration {
	return r.cost
}

// Dump dump the request
func (r *Resp) Dump() string {
	dump := new(dumpBuffer)
	if r.r.flag&Lcost != 0 {
		dump.WriteString(fmt.Sprint(r.cost))
	}
	r.dumpRequest(dump)
	l := dump.Len()
	if l > 0 {
		dump.WriteString("=================================")
		l = dump.Len()
	}

	r.dumpResponse(dump)

	return dump.String()
}
