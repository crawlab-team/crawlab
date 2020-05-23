package req

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// default *Req
var std = New()

// flags to decide which part can be outputed
const (
	LreqHead  = 1 << iota // output request head (request line and request header)
	LreqBody              // output request body
	LrespHead             // output response head (response line and response header)
	LrespBody             // output response body
	Lcost                 // output time costed by the request
	LstdFlags = LreqHead | LreqBody | LrespHead | LrespBody
)

// Param represents  http request param
type Param map[string]interface{}

// QueryParam is used to force append http request param to the uri
type QueryParam map[string]interface{}

// Host is used for set request's Host
type Host string

// FileUpload represents a file to upload
type FileUpload struct {
	// filename in multipart form.
	FileName string
	// form field name
	FieldName string
	// file to uplaod, required
	File io.ReadCloser
}

type DownloadProgress func(current, total int64)

type UploadProgress func(current, total int64)

// File upload files matching the name pattern such as
// /usr/*/bin/go* (assuming the Separator is '/')
func File(patterns ...string) interface{} {
	matches := []string{}
	for _, pattern := range patterns {
		m, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		matches = append(matches, m...)
	}
	if len(matches) == 0 {
		return errors.New("req: no file have been matched")
	}
	uploads := []FileUpload{}
	for _, match := range matches {
		if s, e := os.Stat(match); e != nil || s.IsDir() {
			continue
		}
		file, _ := os.Open(match)
		uploads = append(uploads, FileUpload{
			File:      file,
			FileName:  filepath.Base(match),
			FieldName: "media",
		})
	}

	return uploads
}

type bodyJson struct {
	v interface{}
}

type bodyXml struct {
	v interface{}
}

// BodyJSON make the object be encoded in json format and set it to the request body
func BodyJSON(v interface{}) *bodyJson {
	return &bodyJson{v: v}
}

// BodyXML make the object be encoded in xml format and set it to the request body
func BodyXML(v interface{}) *bodyXml {
	return &bodyXml{v: v}
}

// Req is a convenient client for initiating requests
type Req struct {
	client           *http.Client
	jsonEncOpts      *jsonEncOpts
	xmlEncOpts       *xmlEncOpts
	flag             int
	progressInterval time.Duration
}

// New create a new *Req
func New() *Req {
	// default progress reporting interval is 200 milliseconds
	return &Req{flag: LstdFlags, progressInterval: 200 * time.Millisecond}
}

type param struct {
	url.Values
}

func (p *param) getValues() url.Values {
	if p.Values == nil {
		p.Values = make(url.Values)
	}
	return p.Values
}

func (p *param) Copy(pp param) {
	if pp.Values == nil {
		return
	}
	vs := p.getValues()
	for key, values := range pp.Values {
		for _, value := range values {
			vs.Add(key, value)
		}
	}
}
func (p *param) Adds(m map[string]interface{}) {
	if len(m) == 0 {
		return
	}
	vs := p.getValues()
	for k, v := range m {
		vs.Add(k, fmt.Sprint(v))
	}
}

func (p *param) Empty() bool {
	return p.Values == nil
}

// Do execute a http request with sepecify method and url,
// and it can also have some optional params, depending on your needs.
func (r *Req) Do(method, rawurl string, vs ...interface{}) (resp *Resp, err error) {
	if rawurl == "" {
		return nil, errors.New("req: url not specified")
	}
	req := &http.Request{
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	resp = &Resp{req: req, r: r}

	var queryParam param
	var formParam param
	var uploads []FileUpload
	var uploadProgress UploadProgress
	var progress func(int64, int64)
	var delayedFunc []func()
	var lastFunc []func()

	for _, v := range vs {
		switch vv := v.(type) {
		case Header:
			for key, value := range vv {
				req.Header.Add(key, value)
			}
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		case *bodyJson:
			fn, err := setBodyJson(req, resp, r.jsonEncOpts, vv.v)
			if err != nil {
				return nil, err
			}
			delayedFunc = append(delayedFunc, fn)
		case *bodyXml:
			fn, err := setBodyXml(req, resp, r.xmlEncOpts, vv.v)
			if err != nil {
				return nil, err
			}
			delayedFunc = append(delayedFunc, fn)
		case url.Values:
			p := param{vv}
			if method == "GET" || method == "HEAD" {
				queryParam.Copy(p)
			} else {
				formParam.Copy(p)
			}
		case Param:
			if method == "GET" || method == "HEAD" {
				queryParam.Adds(vv)
			} else {
				formParam.Adds(vv)
			}
		case QueryParam:
			queryParam.Adds(vv)
		case string:
			setBodyBytes(req, resp, []byte(vv))
		case []byte:
			setBodyBytes(req, resp, vv)
		case bytes.Buffer:
			setBodyBytes(req, resp, vv.Bytes())
		case *http.Client:
			resp.client = vv
		case FileUpload:
			uploads = append(uploads, vv)
		case []FileUpload:
			uploads = append(uploads, vv...)
		case *http.Cookie:
			req.AddCookie(vv)
		case Host:
			req.Host = string(vv)
		case io.Reader:
			fn := setBodyReader(req, resp, vv)
			lastFunc = append(lastFunc, fn)
		case UploadProgress:
			uploadProgress = vv
		case DownloadProgress:
			resp.downloadProgress = vv
		case func(int64, int64):
			progress = vv
		case context.Context:
			req = req.WithContext(vv)
			resp.req = req
		case error:
			return nil, vv
		}
	}

	if length := req.Header.Get("Content-Length"); length != "" {
		if l, err := strconv.ParseInt(length, 10, 64); err == nil {
			req.ContentLength = l
		}
	}

	if len(uploads) > 0 && (req.Method == "POST" || req.Method == "PUT") { // multipart
		var up UploadProgress
		if uploadProgress != nil {
			up = uploadProgress
		} else if progress != nil {
			up = UploadProgress(progress)
		}
		multipartHelper := &multipartHelper{
			form:             formParam.Values,
			uploads:          uploads,
			uploadProgress:   up,
			progressInterval: resp.r.progressInterval,
		}
		multipartHelper.Upload(req)
		resp.multipartHelper = multipartHelper
	} else {
		if progress != nil {
			resp.downloadProgress = DownloadProgress(progress)
		}
		if !formParam.Empty() {
			if req.Body != nil {
				queryParam.Copy(formParam)
			} else {
				setBodyBytes(req, resp, []byte(formParam.Encode()))
				setContentType(req, "application/x-www-form-urlencoded; charset=UTF-8")
			}
		}
	}

	if !queryParam.Empty() {
		paramStr := queryParam.Encode()
		if strings.IndexByte(rawurl, '?') == -1 {
			rawurl = rawurl + "?" + paramStr
		} else {
			rawurl = rawurl + "&" + paramStr
		}
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req.URL = u

	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}

	for _, fn := range delayedFunc {
		fn()
	}

	if resp.client == nil {
		resp.client = r.Client()
	}

	var response *http.Response
	if r.flag&Lcost != 0 {
		before := time.Now()
		response, err = resp.client.Do(req)
		after := time.Now()
		resp.cost = after.Sub(before)
	} else {
		response, err = resp.client.Do(req)
	}
	if err != nil {
		return nil, err
	}

	for _, fn := range lastFunc {
		fn()
	}

	resp.resp = response

	if _, ok := resp.client.Transport.(*http.Transport); ok && response.Header.Get("Content-Encoding") == "gzip" && req.Header.Get("Accept-Encoding") != "" {
		body, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		response.Body = body
	}

	// output detail if Debug is enabled
	if Debug {
		fmt.Println(resp.Dump())
	}
	return
}

func setBodyBytes(req *http.Request, resp *Resp, data []byte) {
	resp.reqBody = data
	req.Body = ioutil.NopCloser(bytes.NewReader(data))
	req.ContentLength = int64(len(data))
}

func setBodyJson(req *http.Request, resp *Resp, opts *jsonEncOpts, v interface{}) (func(), error) {
	var data []byte
	switch vv := v.(type) {
	case string:
		data = []byte(vv)
	case []byte:
		data = vv
	case *bytes.Buffer:
		data = vv.Bytes()
	default:
		if opts != nil {
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			enc.SetIndent(opts.indentPrefix, opts.indentValue)
			enc.SetEscapeHTML(opts.escapeHTML)
			err := enc.Encode(v)
			if err != nil {
				return nil, err
			}
			data = buf.Bytes()
		} else {
			var err error
			data, err = json.Marshal(v)
			if err != nil {
				return nil, err
			}
		}
	}
	setBodyBytes(req, resp, data)
	delayedFunc := func() {
		setContentType(req, "application/json; charset=UTF-8")
	}
	return delayedFunc, nil
}

func setBodyXml(req *http.Request, resp *Resp, opts *xmlEncOpts, v interface{}) (func(), error) {
	var data []byte
	switch vv := v.(type) {
	case string:
		data = []byte(vv)
	case []byte:
		data = vv
	case *bytes.Buffer:
		data = vv.Bytes()
	default:
		if opts != nil {
			var buf bytes.Buffer
			enc := xml.NewEncoder(&buf)
			enc.Indent(opts.prefix, opts.indent)
			err := enc.Encode(v)
			if err != nil {
				return nil, err
			}
			data = buf.Bytes()
		} else {
			var err error
			data, err = xml.Marshal(v)
			if err != nil {
				return nil, err
			}
		}
	}
	setBodyBytes(req, resp, data)
	delayedFunc := func() {
		setContentType(req, "application/xml; charset=UTF-8")
	}
	return delayedFunc, nil
}

func setContentType(req *http.Request, contentType string) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}
}

func setBodyReader(req *http.Request, resp *Resp, rd io.Reader) func() {
	var rc io.ReadCloser
	switch r := rd.(type) {
	case *os.File:
		stat, err := r.Stat()
		if err == nil {
			req.ContentLength = stat.Size()
		}
		rc = r

	case io.ReadCloser:
		rc = r
	default:
		rc = ioutil.NopCloser(rd)
	}
	bw := &bodyWrapper{
		ReadCloser: rc,
		limit:      102400,
	}
	req.Body = bw
	lastFunc := func() {
		resp.reqBody = bw.buf.Bytes()
	}
	return lastFunc
}

type bodyWrapper struct {
	io.ReadCloser
	buf   bytes.Buffer
	limit int
}

func (b *bodyWrapper) Read(p []byte) (n int, err error) {
	n, err = b.ReadCloser.Read(p)
	if left := b.limit - b.buf.Len(); left > 0 && n > 0 {
		if n <= left {
			b.buf.Write(p[:n])
		} else {
			b.buf.Write(p[:left])
		}
	}
	return
}

type multipartHelper struct {
	form             url.Values
	uploads          []FileUpload
	dump             []byte
	uploadProgress   UploadProgress
	progressInterval time.Duration
}

func (m *multipartHelper) Upload(req *http.Request) {
	pr, pw := io.Pipe()
	bodyWriter := multipart.NewWriter(pw)
	go func() {
		for key, values := range m.form {
			for _, value := range values {
				bodyWriter.WriteField(key, value)
			}
		}
		var upload func(io.Writer, io.Reader) error
		if m.uploadProgress != nil {
			var total int64
			for _, up := range m.uploads {
				if file, ok := up.File.(*os.File); ok {
					stat, err := file.Stat()
					if err != nil {
						continue
					}
					total += stat.Size()
				}
			}
			var current int64
			buf := make([]byte, 1024)
			var lastTime time.Time

			defer func() {
				m.uploadProgress(current, total)
			}()

			upload = func(w io.Writer, r io.Reader) error {
				for {
					n, err := r.Read(buf)
					if n > 0 {
						_, _err := w.Write(buf[:n])
						if _err != nil {
							return _err
						}
						current += int64(n)
						if now := time.Now(); now.Sub(lastTime) > m.progressInterval {
							lastTime = now
							m.uploadProgress(current, total)
						}
					}
					if err == io.EOF {
						return nil
					}
					if err != nil {
						return err
					}
				}
			}
		}

		i := 0
		for _, up := range m.uploads {
			if up.FieldName == "" {
				i++
				up.FieldName = "file" + strconv.Itoa(i)
			}
			fileWriter, err := bodyWriter.CreateFormFile(up.FieldName, up.FileName)
			if err != nil {
				continue
			}
			//iocopy
			if upload == nil {
				io.Copy(fileWriter, up.File)
			} else {
				if _, ok := up.File.(*os.File); ok {
					upload(fileWriter, up.File)
				} else {
					io.Copy(fileWriter, up.File)
				}
			}
			up.File.Close()
		}
		bodyWriter.Close()
		pw.Close()
	}()
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.Body = ioutil.NopCloser(pr)
}

func (m *multipartHelper) Dump() []byte {
	if m.dump != nil {
		return m.dump
	}
	var buf bytes.Buffer
	bodyWriter := multipart.NewWriter(&buf)
	for key, values := range m.form {
		for _, value := range values {
			m.writeField(bodyWriter, key, value)
		}
	}
	for _, up := range m.uploads {
		m.writeFile(bodyWriter, up.FieldName, up.FileName)
	}
	bodyWriter.Close()
	m.dump = buf.Bytes()
	return m.dump
}

func (m *multipartHelper) writeField(w *multipart.Writer, fieldname, value string) error {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"`, fieldname))
	p, err := w.CreatePart(h)
	if err != nil {
		return err
	}
	_, err = p.Write([]byte(value))
	return err
}

func (m *multipartHelper) writeFile(w *multipart.Writer, fieldname, filename string) error {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			fieldname, filename))
	h.Set("Content-Type", "application/octet-stream")
	p, err := w.CreatePart(h)
	if err != nil {
		return err
	}
	_, err = p.Write([]byte("******"))
	return err
}

// Get execute a http GET request
func (r *Req) Get(url string, v ...interface{}) (*Resp, error) {
	return r.Do("GET", url, v...)
}

// Post execute a http POST request
func (r *Req) Post(url string, v ...interface{}) (*Resp, error) {
	return r.Do("POST", url, v...)
}

// Put execute a http PUT request
func (r *Req) Put(url string, v ...interface{}) (*Resp, error) {
	return r.Do("PUT", url, v...)
}

// Patch execute a http PATCH request
func (r *Req) Patch(url string, v ...interface{}) (*Resp, error) {
	return r.Do("PATCH", url, v...)
}

// Delete execute a http DELETE request
func (r *Req) Delete(url string, v ...interface{}) (*Resp, error) {
	return r.Do("DELETE", url, v...)
}

// Head execute a http HEAD request
func (r *Req) Head(url string, v ...interface{}) (*Resp, error) {
	return r.Do("HEAD", url, v...)
}

// Options execute a http OPTIONS request
func (r *Req) Options(url string, v ...interface{}) (*Resp, error) {
	return r.Do("OPTIONS", url, v...)
}

// Get execute a http GET request
func Get(url string, v ...interface{}) (*Resp, error) {
	return std.Get(url, v...)
}

// Post execute a http POST request
func Post(url string, v ...interface{}) (*Resp, error) {
	return std.Post(url, v...)
}

// Put execute a http PUT request
func Put(url string, v ...interface{}) (*Resp, error) {
	return std.Put(url, v...)
}

// Head execute a http HEAD request
func Head(url string, v ...interface{}) (*Resp, error) {
	return std.Head(url, v...)
}

// Options execute a http OPTIONS request
func Options(url string, v ...interface{}) (*Resp, error) {
	return std.Options(url, v...)
}

// Delete execute a http DELETE request
func Delete(url string, v ...interface{}) (*Resp, error) {
	return std.Delete(url, v...)
}

// Patch execute a http PATCH request
func Patch(url string, v ...interface{}) (*Resp, error) {
	return std.Patch(url, v...)
}

// Do execute request.
func Do(method, url string, v ...interface{}) (*Resp, error) {
	return std.Do(method, url, v...)
}
