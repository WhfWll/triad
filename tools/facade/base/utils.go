// Package base
// @Author bcy2007  2025/12/23 14:48
package base

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gobwas/glob"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"smart/tools/ctxio"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	CRLF                                        = "\r\n"
	CommonHeaderStat                     string = "common-header"
	HeaderCheckStat                             = "header-Check"
	REQUEST_CONTEXT_KEY_RequestBareBytes        = "requestBareBytes"
	REQUEST_CONTEXT_INFOMAP                     = "InfoMap"
)

var commonHeader = map[string]string{
	"Accept":                    "Accept",
	"Accept-Charset":            "Accept-Charset",
	"Accept-Encoding":           "Accept-Encoding",
	"Accept-Language":           "Accept-Language",
	"Accept-Ranges":             "Accept-Ranges",
	"Cache-Control":             "Cache-Control",
	"Cc":                        "Cc",
	"Connection":                "Connection",
	"Content-Id":                "Content-Id",
	"Content-Language":          "Content-Language",
	"Content-Length":            "Content-Length",
	"Content-Transfer-Encoding": "Content-Transfer-Encoding",
	"Content-Type":              "Content-Type",
	"Cookie":                    "Cookie",
	"Date":                      "Date",
	"Etag":                      "Etag",
	"Expires":                   "Expires",
	"From":                      "From",
	"Host":                      "Host",
	"If-Modified-Since":         "If-Modified-Since",
	"If-None-Match":             "If-None-Match",
	"In-Reply-To":               "In-Reply-To",
	"Last-Modified":             "Last-Modified",
	"Location":                  "Location",
	"Message-Id":                "Message-Id",
	"Mime-Version":              "Mime-Version",
	"Pragma":                    "Pragma",
	"Received":                  "Received",
	"Return-Path":               "Return-Path",
	"Server":                    "Server",
	"Set-Cookie":                "Set-Cookie",
	"Subject":                   "Subject",
	"To":                        "To",
	"User-Agent":                "User-Agent",
	"X-Forwarded-For":           "X-Forwarded-For",
	"X-Powered-By":              "X-Powered-By",
}

func ParseStringToHostPort(raw string) (host string, port int, err error) {
	if strings.Contains(raw, "://") {
		urlObject, _ := url.Parse(raw)
		if urlObject != nil {
			// 处理 URL
			portRaw := urlObject.Port()
			portInt64, err := strconv.ParseInt(portRaw, 10, 32)
			if err != nil || portInt64 <= 0 {
				switch urlObject.Scheme {
				case "http", "ws":
					port = 80
				case "https", "wss":
					port = 443
				}
			} else {
				port = int(portInt64)
			}

			host = urlObject.Hostname()
			err = nil
			return host, port, err
		}
	}
	// 这里需要处理ipv6的情况，如果是ipv6的话，直接返回
	if ip := net.ParseIP(raw); ip != nil {
		return raw, 0, errors.Errorf("unknown port for [%s]", raw)
	}

	host = stripPort(raw)
	portStr := portOnly(raw)
	if len(portStr) <= 0 {
		return host, 0, errors.Errorf("unknown port for [%s]", raw)
	}

	portStr = strings.TrimSpace(portStr)
	portInt64, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		return host, 0, errors.Errorf("%s parse port(%s) failed: %s", raw, portStr, err)
	}

	port = int(portInt64)
	err = nil
	return
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

func portOnly(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return ""
	}
	if i := strings.Index(hostport, "]:"); i != -1 {
		return hostport[i+len("]:"):]
	}
	if strings.Contains(hostport, "]") {
		return ""
	}
	return hostport[colon+len(":"):]
}

func HostPort(host string, port interface{}) string {
	return fmt.Sprintf("%v:%v", ParseHostToAddrString(host), port)
}

func ParseHostToAddrString(host string) string {
	ip := net.ParseIP(host)
	if ip == nil {
		return host
	}

	if ret := ip.To4(); ret == nil {
		return fmt.Sprintf("[%v]", ip.String())
	}

	return host
}

func StableReaderEx(conn net.Conn, timeout time.Duration, maxSize int) []byte {
	var mu sync.Mutex
	buffer := bytes.NewBuffer(nil)
	readTimeout := 1000 * time.Millisecond
	readAsyncTimeout := 250 * time.Millisecond
	readGapTimeout := 350 * time.Millisecond

	defer conn.SetDeadline(time.Now().Add(3 * time.Minute))

	ddlCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan bool)
	go func() {
		defer close(done)
		ch := make([]byte, 1)
		for {
			if err := conn.SetDeadline(time.Now().Add(readAsyncTimeout)); err != nil {
				log.Debugf("SetDeadline failed: %v", err)
				return
			}

			n, err := conn.Read(ch)
			if n > 0 {
				mu.Lock()
				buffer.Write(ch)
				currentLen := buffer.Len()
				mu.Unlock()

				if currentLen >= maxSize {
					done <- true
					return
				}
			}

			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					done <- true
					return
				}
				if conn.RemoteAddr() != nil {
					log.Debugf("conn[%s] met error: %v", conn.RemoteAddr().String(), err)
				}
			}

			select {
			case <-ddlCtx.Done():
				done <- false
				return
			default:
			}
		}
	}()

	var lastLen int
	timer := time.NewTimer(readTimeout)
	defer timer.Stop()

	for {
		<-timer.C
		timer.Reset(readTimeout)

		mu.Lock()
		currentLen := buffer.Len()
		mu.Unlock()

		if currentLen == 0 || currentLen == lastLen {
			break
		}
		lastLen = currentLen
	}

	if success := <-done; success {
		time.Sleep(readGapTimeout)
	}

	mu.Lock()
	defer mu.Unlock()
	return buffer.Bytes()
}

func Md5(i interface{}) string {
	raw := md5.Sum(interfaceToBytes(i))
	return EncodeToHex(raw[:])
}

func interfaceToBytes(i interface{}) []byte {
	var bytes []byte

	switch ret := i.(type) {
	case string:
		bytes = []byte(ret)
	case []byte:
		bytes = ret
	case io.Reader:
		bytes, _ = ioutil.ReadAll(ret)
	default:
		bytes = []byte(fmt.Sprint(i))
	}

	return bytes
}

func EncodeToHex(i interface{}) string {
	raw := interfaceToBytes(i)
	return hex.EncodeToString(raw)
}

func ReadHTTPRequestFromBufioReader(reader *bufio.Reader) (*http.Request, error) {
	return readHTTPRequestFromBufioReader(reader, false, nil)
}

func readHTTPRequestFromBufioReader(reader *bufio.Reader, fixContentLength bool, onFirstLine func(string)) (*http.Request, error) {
	rawPacket := new(bytes.Buffer)

	req := &http.Request{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Form:       nil,
		Body:       http.NoBody,
		RequestURI: "", // do not handle it as client
		TLS:        nil,
	}

	defer func() {
		if req != nil && req.URL != nil {
			req.URL.Opaque = ""
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}

		if err := recover(); err != nil {
			log.Errorf("ReadHTTPRequestEx panic: %v", err)
			PrintCurrentGoroutineRuntimeStack()
		}
	}()

	// parse first line
	firstLine, err := BufioReadLine(reader)
	if err != nil {
		return nil, fmt.Errorf(`Read Request FirstLine Failed: %s`, err)
	}
	if onFirstLine != nil {
		onFirstLine(string(firstLine))
	}
	rawPacket.Write(firstLine)
	rawPacket.WriteString(CRLF)

	// handle proto
	perfix, firstLine, _ := CutBytesPrefixFunc(firstLine, NotSpaceRune)
	method, uri, proto, ok := ParseHTTPRequestLine(string(firstLine))
	if ok {
		req.Method = method
		req.RequestURI = uri
		req.Proto = proto
		_, after, ok := strings.Cut(proto, "/")
		if ok {
			major, minor, _ := strings.Cut(after, ".")
			req.ProtoMajor, _ = strconv.Atoi(major)
			req.ProtoMinor, _ = strconv.Atoi(minor)
		}
	} else {
		return nil, fmt.Errorf(`Parse Request FirstLine(%v) Failed: %s`, strconv.Quote(string(firstLine)), firstLine)
	}

	var (
		// RequestURI > URL > Host in header
		hostInURL    string
		hostInHeader string
	)

	/*
		handle headers
			1. keep gzip
			2. keep chunked if have
		    3. smuggle use max(chunked, contentLength)

		if smuggle { keep cl and te }
		if not smuggle { if te keep te }
	*/
	// close is default in 0.9 or 1.0
	defaultClose := (req.ProtoMajor == 1 && req.ProtoMinor == 0) || req.ProtoMajor < 1
	header := make(http.Header)
	useContentLength := false
	contentLengthInt := 0
	useTransferEncodingChunked := false

	_ = ScanHTTPHeader(reader, func(lineBytes []byte) {
		if len(lineBytes) == 0 {
			rawPacket.WriteString(CRLF)
			return
		}
		rawPacket.Write(lineBytes)
		rawPacket.WriteString(CRLF)

		before, after, _ := bytes.Cut(lineBytes, []byte{':'})
		keyStr := string(before)
		valStr := strings.TrimLeftFunc(string(after), unicode.IsSpace)

		if _, isCommonHeader := commonHeader[keyStr]; isCommonHeader {
			keyStr = http.CanonicalHeaderKey(keyStr)
		}

		isSingletonHeader := false
		switch strings.ToLower(keyStr) {
		case "content-length":
			useContentLength = true
			contentLengthInt = Atoi(valStr)
			if contentLengthInt != 0 || !ShouldRemoveZeroContentLengthHeader(method) {
				header[keyStr] = []string{valStr}
				req.ContentLength = int64(contentLengthInt)
			}
		case "host":
			hostInHeader = valStr
		case "content-type":
			isSingletonHeader = true
		case `transfer-encoding`:
			req.TransferEncoding = []string{valStr}
			if IContains(valStr, "chunked") {
				useTransferEncodingChunked = true
			}
		case "connection":
			if strings.EqualFold(valStr, "close") {
				defaultClose = true
			} else if strings.EqualFold(valStr, "keep-alive") {
				defaultClose = false
			}
		}

		// add header
		if keyStr == "" {
			return
		}
		if isSingletonHeader {
			header[keyStr] = append(header[keyStr], valStr)
			return
		}
		header[keyStr] = append(header[keyStr], valStr)
	}, perfix, false)

	// uri is very complex
	// utf8 valid or not
	before, fragment, haveFragment := strings.Cut(req.RequestURI, "#")
	var urlIns *url.URL
	if method == "CONNECT" {
		urlIns = new(url.URL)
		// if connect, the uri should be host:port
		host, port, _ := ParseStringToHostPort(before)
		if port > 0 {
			urlIns.Host = HostPort(host, port)
		} else {
			if strings.HasPrefix(hostInHeader, ":") {
				port := Atoi(hostInHeader[1:])
				if port > 0 && port <= 65535 {
					urlIns.Host = HostPort(host, port)
				} else {
					urlIns.Host = strings.Trim(host, "/")
				}
			} else {
				urlIns.Host = strings.Trim(host, "/")
			}
		}
	} else if urlIns, _ = url.ParseRequestURI(before); urlIns == nil {
		// remove : begin
		// utf8 invalid
		urlIns = new(url.URL)
		if IsHttpOrHttpsUrl(req.RequestURI) {
			urlIns, err = url.Parse(req.RequestURI)
			if err != nil {
				return nil, fmt.Errorf("parse uri-url (%v) failed: %s", req.RequestURI, err)
			}
		} else {
			urlIns.Path, urlIns.RawQuery, _ = strings.Cut(req.RequestURI, "?")
		}
	}

	if urlIns != nil && haveFragment {
		urlIns.Fragment = fragment
	}
	req.URL = urlIns

	// handle host
	hostInURL = req.URL.Host
	if ret := strings.LastIndex(hostInURL, ":"); ret >= 0 {
		hostname, portStr := strings.TrimSpace(hostInURL[:ret]), Atoi(hostInURL[ret+1:])
		if hostname == "" || portStr == 0 {
			req.URL.Host = ""
			hostInURL = ""
		}
	}

	req.Close = defaultClose
	req.Header = header

	// handling host
	if hostInHeader == "" && hostInURL == "" && method == "CONNECT" {
		return nil, errors.New(`Host(inHeader/inURL) is empty in CONNECT method`)
	}

	var host string
	if hostInURL != "" {
		host = hostInURL
	} else {
		host = hostInHeader
	}
	req.Host = host
	if req.URL.Host == "" {
		req.URL.Host = hostInHeader
	}
	bodyRawBuf := new(bytes.Buffer)
	if fixContentLength {
		// by reader
		raw, _ := io.ReadAll(reader)
		rawPacket.Write(raw)
		if useContentLength && !useTransferEncodingChunked {
			req.ContentLength = int64(len(raw))
			shrinkHeader(req.Header, "content-length")
			req.Header.Set("Content-Length", strconv.Itoa(len(raw)))
		}
		bodyRawBuf.Write(raw)
	} else {
		// by header
		if useContentLength && useTransferEncodingChunked {
			log.Warn("content-length and transfer-encoding chunked both exist, try smuggle? use content-length first!")
			if contentLengthInt > 0 {
				// smuggle
				bodyRaw, _ := io.ReadAll(io.NopCloser(io.LimitReader(reader, int64(contentLengthInt))))
				rawPacket.Write(bodyRaw)
				bodyRawBuf.Write(bodyRaw)
				if ret := contentLengthInt - len(bodyRaw); ret > 0 {
					bodyRawBuf.WriteString(strings.Repeat("\n", ret))
				}
			} else {
				// chunked
				_, fixed, _, err := HTTPChunkedDecoderWithRestBytes(reader)
				rawPacket.Write(fixed)
				if err != nil {
					return nil, fmt.Errorf("chunked decoder error: %v", err)
				}
				bodyRawBuf.Write(fixed)
			}
		} else if !useContentLength && useTransferEncodingChunked {
			// handle chunked
			_, fixed, _, err := HTTPChunkedDecoderWithRestBytes(reader)
			rawPacket.Write(fixed)
			if err != nil {
				return nil, fmt.Errorf("chunked decoder error: %v", err)
			}
			if len(fixed) > 0 {
				bodyRawBuf.Write(fixed)
			}
		} else {
			// handle content-length as default
			bodyRaw, err := io.ReadAll(io.NopCloser(io.LimitReader(reader, int64(contentLengthInt))))
			rawPacket.Write(bodyRaw)
			if err != nil && err != io.EOF {
				if !errors.Is(err, io.ErrUnexpectedEOF) {
					return nil, fmt.Errorf("read body error: %v", err)
				}
				log.Warnf("read body error: %v", err)
			}
			bodyLen := len(bodyRaw)
			bodyRawBuf.Write(bodyRaw)
			bodyRawBuf.WriteString(strings.Repeat("\n", contentLengthInt-bodyLen))
		}
	}
	if bodyRawBuf.Len() == 0 {
		req.Body = http.NoBody
	} else {
		req.Body = io.NopCloser(bodyRawBuf)
	}
	if req.URL != nil && req.URL.Host != "" {
		req.Host = req.URL.Host
	}
	SetBareRequestBytes(req, rawPacket.Bytes())
	return req, nil
}

func PrintCurrentGoroutineRuntimeStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("Current goroutine call stack:\n%s\n", buf[:n])
}

func BufioReadLine(reader *bufio.Reader) ([]byte, error) {
	if reader == nil {
		return nil, errors.New("empty reader(bufio)")
	}

	var buf bytes.Buffer
	for {
		tmp, isPrefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		buf.Write(tmp)
		if !isPrefix {
			return buf.Bytes(), nil
		}
	}
}

func CutBytesPrefixFunc(raw []byte, handle func(rune) bool) ([]byte, []byte, bool) {
	index := bytes.IndexFunc(raw, handle)
	if index < 0 {
		return nil, raw, false
	}
	return raw[:index], raw[index:], true
}

func ParseHTTPRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.LastIndex(line[s1+1:], " ")
	if s1 < 0 {
		return
	}

	httpVersion := "HTTP/1.1"
	if s2 < 0 {
		return line[:s1], line[s1+1:], httpVersion, true
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

func NotSpaceRune(r rune) bool {
	return !unicode.IsSpace(r)
}

func ScanHTTPHeader(reader io.Reader, headerCallback func(rawHeader []byte), prefix []byte, isResp bool) error {
	if isResp {
		return ScanHTTPHeaderWithHeaderFolding(reader, headerCallback, prefix)
	}
	return ScanHTTPHeaderSimple(reader, headerCallback, prefix)
}

func ScanHTTPHeaderWithHeaderFolding(reader io.Reader, headerCallback func(rawHeader []byte), prefix []byte) error {
	var headerRawCache []byte
	var currentSata = CommonHeaderStat
	var headerFoldingPrefix = make([]byte, 0)

	setHeaderFoldingPrefix := func(foldingPrefix []byte) {
		headerFoldingPrefix = foldingPrefix
	}

	setCurrentStat := func(stat string) {
		currentSata = stat
	}

	pushHeaderRawData := func(raw []byte) {
		headerRawCache = append(headerRawCache, raw...)
	}

	emitHeaderRaw := func() {
		if headerCallback != nil {
			headerCallback(headerRawCache)
		}
		headerRawCache = make([]byte, 0)
	}

	defer emitHeaderRaw()

	trimPrefix := func(raw []byte) []byte {
		minLen := Min(len(prefix), len(raw))
		i := 0
		for ; i < minLen; i++ {
			if raw[i] != prefix[i] {
				break
			}
		}
		return raw[i:]
	}

	for {
		lineBytes, err := ReadLine(reader)
		if err != nil && err != io.EOF {
			return errors.Wrap(err, "read HTTPResponse header failed")
		}
		lineBytes = trimPrefix(lineBytes)
	Retry:
		switch currentSata {
		case CommonHeaderStat:
			if len(lineBytes) == 0 {
				return nil
			}
			for i, b := range lineBytes {
				if b != ' ' && b != '\t' {
					setHeaderFoldingPrefix(lineBytes[:i])
					break
				}
			}
			pushHeaderRawData(lineBytes)
			setCurrentStat(HeaderCheckStat)
		case HeaderCheckStat:
			checkLine := bytes.TrimPrefix(lineBytes, headerFoldingPrefix)
			if len(checkLine) > 0 && (checkLine[0] == ' ' || checkLine[0] == '\t') {
				pushHeaderRawData(append([]byte(CRLF), checkLine...))
			} else {
				emitHeaderRaw()
				setCurrentStat(CommonHeaderStat)
				goto Retry
			}
		}
	}
}

func ScanHTTPHeaderSimple(reader io.Reader, headerCallback func(rawHeader []byte), prefix []byte) error {
	emitHeaderRaw := func(raw []byte) {
		if headerCallback != nil {
			headerCallback(raw)
		}
	}
	trimPrefix := func(raw []byte) []byte {
		minLen := Min(len(prefix), len(raw))
		i := 0
		for ; i < minLen; i++ {
			if raw[i] != prefix[i] {
				break
			}
		}
		return raw[i:]
	}

	for {
		lineBytes, err := ReadLine(reader)
		if err != nil && err != io.EOF {
			return errors.Wrap(err, "read HTTPResponse header failed")
		}
		lineBytes = trimPrefix(lineBytes)
		if len(bytes.TrimSpace(lineBytes)) == 0 {
			emitHeaderRaw(nil)
			return nil
		}
		emitHeaderRaw(lineBytes)
	}
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func ReadLine(reader io.Reader) ([]byte, error) {
	lineRaw, err := ReadUntilStableEx(reader, true, nil, 0, 0, '\n')
	if err != nil {
		return lineRaw, err
	}
	return bytes.TrimRight(lineRaw, "\r\n"), nil
}

// ReadUntilStableEx allow skip timeout, read until stop word or timeout
func ReadUntilStableEx(reader io.Reader, noTimeout bool, conn net.Conn, timeout time.Duration, stableTimeout time.Duration, sep ...byte) ([]byte, error) {
	buf := make([]byte, 1)
	var result bytes.Buffer

	var ctx context.Context
	var cancel context.CancelFunc
	if noTimeout {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	}
	defer cancel()
	stopStep := 0

	wrapperTimeout := func(originReader io.Reader) io.Reader {
		if noTimeout {
			return originReader
		}

		if conn != nil {
			_ = conn.SetReadDeadline(time.Now().Add(stableTimeout))
			return originReader
		} else {
			return ctxio.NewReader(TimeoutContext(stableTimeout), originReader)
		}
	}
	recoverTimeout := func() {
		if noTimeout {
			return
		}

		if conn != nil {
			_ = conn.SetReadDeadline(time.Time{})
		}
	}

	for {
		n, err := io.ReadFull(wrapperTimeout(reader), buf)
		recoverTimeout()

		if err != nil {
			var netOpError interface{ Timeout() bool }
			if errors.As(err, &netOpError) && netOpError != nil && netOpError.Timeout() {
				if result.Len() > 0 {
					return result.Bytes(), nil
				} else {
					return nil, err
				}
			}
			return result.Bytes(), err
		}
		if n > 0 {
			result.Write(buf[:n])
		}
		select {
		case <-ctx.Done():
			if result.Len() > 0 {
				return result.Bytes(), nil
			}
			return nil, errors.New("i/o timeout")
		default:
		}
		if n == 1 && stopStep < len(sep) {
			if buf[0] == sep[stopStep] {
				stopStep++
				if stopStep == len(sep) {
					return result.Bytes(), nil
				}
			} else {
				stopStep = 0
			}
		}
	}
}

func TimeoutContext(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}

var _noContentLengthHeader = map[string]struct{}{
	"GET":     {},
	"HEAD":    {},
	"DELETE":  {},
	"OPTIONS": {},
	"CONNECT": {},
	"get":     {},
	"head":    {},
	"delete":  {},
	"options": {},
	"connect": {},
}

func IContains(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

func IHasPrefix(s, sub string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(sub))
}

func IsHttpOrHttpsUrl(raw string) bool {
	return strings.HasPrefix(strings.TrimSpace(raw), "http://") || strings.HasPrefix(strings.TrimSpace(raw), "https://")
}

func shrinkHeader(header http.Header, key string) {
	values := getHeaderValueList(header, key)
	delete(header, key)
	delete(header, http.CanonicalHeaderKey(key))
	if len(values) > 0 {
		header[http.CanonicalHeaderKey(key)] = values
	}
}

func getHeaderValueList(header http.Header, key string) []string {
	if header == nil {
		return nil
	}
	cKey := http.CanonicalHeaderKey(key)
	if key == cKey {
		if raw, ok := header[key]; ok {
			return raw
		}
		return []string{}
	}

	v1, _ := header[key]
	v2, _ := header[cKey]
	vals := make([]string, 0, len(v1)+len(v2))
	var m = map[string]any{}
	for _, v := range [][]string{v1, v2} {
		for _, i := range v {
			if i == "" {
				continue
			}
			if _, ok := m[i]; ok {
				continue
			}
			m[i] = i
			vals = append(vals, i)
		}
	}
	return vals
}

func SetContextValueInfoFromRequest(r *http.Request, key string, value any) {
	infoMap := GetContextInfoMap(r)
	infoMap.Store(key, value)
}

func SetBareRequestBytes(r *http.Request, bytes []byte) {
	//if len(GetBareRequestBytes(r)) != 0 {
	//	log.Debug("SetBareRequestBytes: bare request bytes already set, ignore")
	//	return
	//}
	SetContextValueInfoFromRequest(r, REQUEST_CONTEXT_KEY_RequestBareBytes, string(bytes))
}

func GetContextInfoMap(r *http.Request) *sync.Map {
	if r == nil {
		return new(sync.Map)
	}
	if r.Context() == nil {
		*r = *r.WithContext(context.Background())
	}
	raw := r.Context().Value(REQUEST_CONTEXT_INFOMAP)
	if raw == nil {
		return _getContextInfoMap(r)
	}
	result, tOk := raw.(*sync.Map)
	if !tOk {
		return _getContextInfoMap(r)
	}
	return result
}

func _getContextInfoMap(r *http.Request) *sync.Map {
	value := r.Context().Value(REQUEST_CONTEXT_INFOMAP)
	var infoMap *sync.Map
	var uid string
	if value == nil {
		uid = uuid.New().String()
		ret := new(sync.Map)
		ret.Store("uuid", uid)
		*r = *r.WithContext(context.WithValue(r.Context(), REQUEST_CONTEXT_INFOMAP, ret))
		value = ret
		infoMap = ret
	} else {
		var ok bool
		infoMap, ok = value.(*sync.Map)
		if !ok {
			return nil
		}
	}
	if uid == "" {
		uidRaw, ok := infoMap.Load("uuid")
		if ok {
			uid = uidRaw.(string)
		}
	}
	return infoMap
}

func HTTPChunkedDecoderWithRestBytes(raw io.Reader) ([]byte, []byte, io.Reader, error) {
	return readChunkedDataFromReader(raw)
}

func InterfaceToBytes(i interface{}) (result []byte) {
	return AnyToBytes(i)
}

func AnyToBytes(i interface{}) (result []byte) {
	var b []byte
	defer func() {
		if err := recover(); err != nil {
			result = []byte(fmt.Sprintf("%v", i))
		}
	}()

	if i == nil {
		return []byte{}
	}

	switch s := i.(type) {
	case nil:
		return []byte{}
	case string:
		b = []byte(s)
	case []byte:
		b = s[0:]
	case bool:
		b = []byte(strconv.FormatBool(s))
	case float64:
		return []byte(strconv.FormatFloat(s, 'f', -1, 64))
	case float32:
		return []byte(strconv.FormatFloat(float64(s), 'f', -1, 32))
	case int:
		return []byte(strconv.Itoa(s))
	case int64:
		return []byte(strconv.FormatInt(s, 10))
	case int32:
		return []byte(strconv.Itoa(int(s)))
	case int16:
		return []byte(strconv.FormatInt(int64(s), 10))
	case int8:
		return []byte(strconv.FormatInt(int64(s), 10))
	case uint:
		return []byte(strconv.FormatUint(uint64(s), 10))
	case uint64:
		return []byte(strconv.FormatUint(s, 10))
	case uint32:
		return []byte(strconv.FormatUint(uint64(s), 10))
	case uint16:
		return []byte(strconv.FormatUint(uint64(s), 10))
	case uint8:
		return []byte(strconv.FormatUint(uint64(s), 10))
	case fmt.Stringer:
		return []byte(s.String())
	case error:
		return []byte(s.Error())
	// case io.Reader:
	//	if ret != nil && ret.Read != nil {
	//		bytes, _ = ioutil.ReadAll(ret)
	//		return bytes
	//	}
	//	return []byte(fmt.Sprintf("%v", i))
	default:
		// 尝试将i作为map转换成JSON
		if jsonBytes, err := json.Marshal(i); err == nil {
			b = jsonBytes
		} else {
			// 如果转换失败，则回退到使用fmt.Sprintf
			b = []byte(fmt.Sprintf("%v", i))
		}
	}

	return b
}

func MatchAllOfGlob(
	i interface{}, re ...string) bool {
	if len(re) <= 0 {
		return false
	}

	raw := interfaceToStr(i)
	for _, r := range re {
		if !glob.MustCompile(r).Match(raw) {
			return false
		}
	}
	return true
}

func interfaceToStr(i interface{}) string {
	return InterfaceToString(i)
}

func InterfaceToString(i interface{}) string {
	if a, ok := i.(interface{ String() string }); ok {
		return a.String()
	}
	return AnyToString(i)
}

func AnyToString(i interface{}) (result string) {
	defer func() {
		if err := recover(); err != nil {
			result = string(fmt.Sprintf("%v", i))
		}
	}()

	if i == nil {
		return ""
	}

	switch s := i.(type) {
	case nil:
		return ""
	case string:
		result = s
	case []byte:
		result = string(s[0:])
	case bool:
		result = strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		return strconv.Itoa(s)
	case int64:
		return strconv.FormatInt(s, 10)
	case int32:
		return strconv.Itoa(int(s))
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case uint:
		return strconv.FormatUint(uint64(s), 10)
	case uint64:
		return strconv.FormatUint(s, 10)
	case uint32:
		return strconv.FormatUint(uint64(s), 10)
	case uint16:
		return strconv.FormatUint(uint64(s), 10)
	case uint8:
		return strconv.FormatUint(uint64(s), 10)
	case fmt.Stringer:
		return s.String()
	case error:
		return s.Error()
	default:
		// 尝试将i作为map转换成JSON
		if jsonBytes, err := json.Marshal(i); err == nil {
			result = string(jsonBytes)
		} else {
			// 如果转换失败，则回退到使用fmt.Sprintf
			result = fmt.Sprintf("%v", i)
		}
	}

	return
}
