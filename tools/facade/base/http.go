// Package base
// @Author bcy2007  2025/12/23 15:15
package base

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func HttpDumpWithBody(i interface{}, body bool) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("HttpDumpWithBody panic: %v", err)
			PrintCurrentGoroutineRuntimeStack()
		}
	}()
	switch ret := i.(type) {
	case *http.Request:
		ret.Close = false
		return DumpHTTPRequest(ret, body)
	case http.Request:
		return HttpDumpWithBody(&ret, body)
	case *http.Response:
		return DumpHTTPResponse(ret, body)
	case http.Response:
		return HttpDumpWithBody(&ret, body)
	default:
		return nil, fmt.Errorf("error type for http.dump, Type: [%v]", reflect.TypeOf(i))
	}
}

func DumpHTTPRequest(req *http.Request, loadBody bool) ([]byte, error) {
	if req == nil {
		return nil, errors.New("nil request")
	}
	var (
		h2                      bool
		transferEncodingChunked bool
		contentLengthExisted    bool
		contentLengthInt        int64
	)

	header := make(http.Header)
	for k, v := range req.Header {
		header[k] = v
	}

	_ = contentLengthInt
	if len(req.TransferEncoding) > 0 {
		for _, v := range req.TransferEncoding {
			if IContains(v, "chunked") {
				transferEncodingChunked = true
				break
			}
		}
	}
	if !transferEncodingChunked {
		if ret := getHeaderValue(header, "transfer-encoding"); ret != "" {
			if IContains(ret, "chunked") {
				transferEncodingChunked = true
			}
		}
	}

	if req.ProtoMajor == 2 || strings.HasPrefix(req.Proto, "HTTP/2") {
		h2 = true
	}

	if ret := getHeaderValue(header, "content-length"); ret != "" || req.ContentLength > 0 {
		contentLengthExisted = true
		if ret != "" {
			contentLengthInt = int64(Atoi(ret))
		} else {
			contentLengthInt = req.ContentLength
		}
	}

	var buf bytes.Buffer
	buf.WriteString(req.Method)
	buf.WriteString(" ")
	if req.Method == "CONNECT" ||
		(len(req.RequestURI) != 0 && !strings.HasPrefix(req.RequestURI, "http")) {
		buf.WriteString(req.RequestURI)
	} else {
		uri := req.URL.RequestURI()
		buf.WriteString(uri)
	}
	buf.WriteString(" ")
	if h2 {
		req.Proto = "HTTP/2.0"
	} else {
		req.Proto = fmt.Sprint("HTTP/", req.ProtoMajor, ".", req.ProtoMinor)
	}
	buf.WriteString(req.Proto)
	buf.WriteString(CRLF)

	// handle host
	buf.WriteString("Host: ")
	if ret := getHeaderValue(header, "host"); ret == "" {
		if req.Host != "" {
			buf.WriteString(req.Host)
		} else if req.URL.Host != "" {
			buf.WriteString(req.URL.Host)
		}
	} else {
		buf.WriteString(ret)
	}
	buf.WriteString(CRLF)
	//shrinkHeader(header, "content-type")

	for k := range header {
		switch strings.ToLower(k) {
		case "host", "content-length", "transfer-encoding":
			continue
		}
		vals, ok := header[k]
		if !ok {
			continue
		}
		for _, v := range vals {
			buf.WriteString(k)
			buf.WriteString(": ")
			buf.WriteString(v)
			buf.WriteString(CRLF)
		}

		//cKey := http.CanonicalHeaderKey(k)
		//if cKey != k {
		//	vals, ok = header[cKey]
		//	if !ok {
		//		continue
		//	}
		//	for _, v := range vals {
		//		buf.WriteString(k)
		//		buf.WriteString(": ")
		//		buf.WriteString(v)
		//		buf.WriteString(CRLF)
		//	}
		//}
	}

	if req.Body == nil {
		req.Body = http.NoBody
	}
	rawBody, _ := io.ReadAll(req.Body)
	var backupBody = io.NopCloser(bytes.NewReader(rawBody))
	defer func() {
		req.Body = backupBody
	}()

	haveBody := len(rawBody) > 0
	// handle cl / te
	if transferEncodingChunked {
		buf.WriteString("Transfer-Encoding: chunked\r\n")
		// check body is chunked or not
		// if not, encode it
		if haveBody {
			decoded, fixed, _ := ReadHTTPChunkedDataWithFixed(rawBody)
			if len(decoded) == 0 {
				rawBody = HTTPChunkedEncode(rawBody)
			} else {
				rawBody = fixed
			}
		}
	} else {
		if haveBody || !ShouldRemoveZeroContentLengthHeader(req.Method) || contentLengthExisted {
			buf.WriteString("Content-Length: ")
			buf.WriteString(fmt.Sprint(len(rawBody)))
			buf.WriteString(CRLF)
		}
	}

	buf.WriteString(CRLF)
	if loadBody {
		buf.Write(rawBody)
	}
	return buf.Bytes(), nil
}

func getHeaderValue(header http.Header, key string) string {
	vals := getHeaderValueList(header, key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func ShouldRemoveZeroContentLengthHeader(s string) bool {
	_, ok := _noContentLengthHeader[s]
	return ok
}

func DumpHTTPResponse(rsp *http.Response, loadBody bool, wr ...io.Writer) ([]byte, error) {
	if rsp == nil {
		return nil, errors.New("nil response")
	}

	var (
		transferEncodingChunked bool
		contentLengthExisted    bool
		contentLengthInt        int64
	)

	header := make(http.Header)
	for k, v := range rsp.Header {
		header[k] = v
	}

	// handle transfer-encoding
	if len(rsp.TransferEncoding) > 0 {
		for _, v := range rsp.TransferEncoding {
			if IContains(v, "chunked") {
				transferEncodingChunked = true
				break
			}
		}
	}
	if !transferEncodingChunked {
		if ret := getHeaderValue(header, "transfer-encoding"); ret != "" {
			if IContains(ret, "chunked") {
				transferEncodingChunked = true
			}
		}
	}

	// handle content-length
	if rsp.ContentLength > 0 {
		contentLengthExisted = true
		contentLengthInt = rsp.ContentLength
	} else {
		if ret := getHeaderValue(header, "content-length"); ret != "" {
			contentLengthExisted = true
			rsp.ContentLength = int64(Atoi(ret))
			contentLengthInt = rsp.ContentLength
		}
	}

	var cacheBuf = new(bytes.Buffer)
	var wrs = make([]io.Writer, 0, len(wr)+1)
	wrs = append(wrs, cacheBuf)
	wrs = append(wrs, wr...)

	var buf = bufio.NewWriter(io.MultiWriter(wrs...))

	// handle proto
	protoRaw := rsp.Proto
	if rsp.ProtoMajor <= 0 && rsp.ProtoMinor <= 0 {
		rsp.ProtoMajor = 1
		rsp.ProtoMinor = 1
	}
	if protoRaw == "" {
		protoRaw = fmt.Sprintf("HTTP/%d.%d", rsp.ProtoMajor, rsp.ProtoMinor)
	}
	buf.WriteString(protoRaw)
	buf.WriteString(" ")
	if rsp.Status == "" {
		if rsp.StatusCode <= 0 {
			rsp.StatusCode = 200
			rsp.Status = "200 OK"
		} else {
			rsp.Status = fmt.Sprintf("%d %s", rsp.StatusCode, http.StatusText(rsp.StatusCode))
		}
	}
	buf.WriteString(rsp.Status)
	buf.WriteString(CRLF)
	buf.Flush()

	// handle server first
	shrinkHeader(header, "server")
	if ret := header.Get("server"); ret != "" {
		header.Set("Server", ret)
		buf.WriteString("Server: ")
		buf.WriteString(ret)
		buf.WriteString(CRLF)
		buf.Flush()
	}

	if rsp.Close {
		header.Set("connection", "close")
	}
	shrinkHeader(header, "connection") // just one connection header
	if ret := header.Get("connection"); ret != "" {
		header.Set("connection", ret)
		buf.WriteString("Connection: ")
		buf.WriteString(ret)
		buf.WriteString(CRLF)
		buf.Flush()
	}

	shrinkHeader(header, "content-length")
	for k := range header {
		switch strings.ToLower(k) {
		case "transfer-encoding", "content-length", "server", "connection":
			continue
		}

		vals, ok := header[k]
		if !ok {
			continue
		}
		for _, v := range vals {
			buf.WriteString(k)
			buf.WriteString(": ")
			buf.WriteString(v)
			buf.WriteString(CRLF)
		}

		cKey := http.CanonicalHeaderKey(k)
		if cKey != k {
			vals, ok = header[cKey]
			if !ok {
				continue
			}
			for _, v := range vals {
				buf.WriteString(k)
				buf.WriteString(": ")
				buf.WriteString(v)
				buf.WriteString(CRLF)
			}
		}
	}

	buf.Flush()
	if rsp.Body == nil {
		rsp.Body = http.NoBody
	}

	rawBody, _ := io.ReadAll(rsp.Body)
	var backupBody = io.NopCloser(bytes.NewReader(rawBody))
	defer func() {
		rsp.Body = backupBody
	}()
	haveBody := len(rawBody) > 0
	if transferEncodingChunked {
		rsp.ContentLength = -1 // unknown
		buf.WriteString("Transfer-Encoding: chunked\r\n")
		buf.Flush()
		if haveBody {
			decode, fixed, _ := ReadHTTPChunkedDataWithFixed(rawBody)
			if len(decode) == 0 {
				rawBody = HTTPChunkedEncode(rawBody)
			} else {
				rawBody = fixed
			}
		}
	} else {
		// handle content-length
		if rsp.StatusCode == 204 || rsp.StatusCode == 304 || (rsp.StatusCode >= 100 && rsp.StatusCode < 200) {
			// omit content-length for 204/304 and 1xx response codes without body
		} else if rsp.Request != nil && rsp.Request.Method == http.MethodConnect && (rsp.StatusCode >= 200 && rsp.StatusCode < 300) {
			// For CONNECT method, omit content-length
		} else {
			if haveBody || contentLengthExisted {
				rsp.ContentLength = int64(len(rawBody))
				contentLengthInt = rsp.ContentLength
				buf.WriteString("Content-Length: ")
				buf.WriteString(strconv.FormatInt(contentLengthInt, 10))
				buf.WriteString(CRLF)
				buf.Flush()
			} else {
				buf.WriteString("Content-Length: 0\r\n")
				buf.Flush()
			}
		}

	}

	buf.WriteString(CRLF)
	if loadBody {
		buf.Write(rawBody)
	}
	buf.Flush()
	return cacheBuf.Bytes(), nil
}

func Atoi(i string) int {
	raw, _ := strconv.Atoi(i)
	return raw
}

func ReadHTTPChunkedDataWithFixed(raw []byte) (data []byte, fixedChunked []byte, rest []byte) {
	blocks, fixed, rest, _ := ReadHTTPChunkedDataWithFixedError(raw)
	return blocks, fixed, rest
}

func ReadHTTPChunkedDataWithFixedError(raw []byte) (data []byte, fixedChunked []byte, rest []byte, _ error) {
	blocks, fixed, reader, err := readChunkedDataFromReader(bytes.NewReader(raw))
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, nil, rest, err
	}
	rest, err = io.ReadAll(reader)
	if err != nil {
		log.Errorf("read chunked data error: %v", err)
	}

	return blocks, fixed, rest, nil
}

func readChunkedDataFromReader(r io.Reader) ([]byte, []byte, io.Reader, error) {
	haveRead := new(bytes.Buffer)
	var reader *bufio.Reader
	switch r.(type) {
	case *bufio.Reader:
		reader = r.(*bufio.Reader)
	default:
		reader = bufio.NewReader(r)
	}
	// read until space
	for {
		spaceByte, err := reader.ReadByte()
		if err != nil {
			return nil, nil, io.MultiReader(bytes.NewReader(haveRead.Bytes()), reader), fmt.Errorf("read chunked (strip left space) data error: %v", err)
		}
		if unicode.IsSpace(rune(spaceByte)) {
			continue
		} else {
			err := reader.UnreadByte()
			if err != nil {
				return nil, nil, io.MultiReader(bytes.NewReader(haveRead.Bytes()), reader), fmt.Errorf("read chunked (strip left space) data error: %v", err)
			}
			break
		}
	}

	var results bytes.Buffer
	var fixedResults bytes.Buffer
	for {
		lineBytes, delim, err := bufioReadLine(reader)
		haveRead.Write(lineBytes)
		haveRead.Write(delim)

		if err != nil && len(lineBytes) > 0 {
			return nil, nil, io.MultiReader(bytes.NewReader(haveRead.Bytes()), reader), err
		}

		var comment []byte
		var commentExisted bool
		handledLineBytes, comment, commentExisted := bytes.Cut(lineBytes, []byte{';'})
		handledLineBytes = bytes.TrimSpace(handledLineBytes)
		size, err := strconv.ParseInt(string(handledLineBytes), 16, 64)
		if err != nil && len(handledLineBytes) > 0 {
			return nil, nil, io.MultiReader(bytes.NewReader(haveRead.Bytes()), reader), err
		}

		if size == 0 {
			lastLine, delim, err := bufioReadLine(reader)
			haveRead.Write(lastLine)
			haveRead.Write(delim)
			if len(lastLine) == 0 {
				fixedResults.WriteString("0\r\n\r\n")
			} else {
				return nil, nil, io.MultiReader(bytes.NewReader(haveRead.Bytes()), reader), fmt.Errorf("last line of chunked data is not empty: %s", lastLine)
			}

			if err != nil {
				if err == io.EOF {
					return results.Bytes(), fixedResults.Bytes(), reader, nil
				}
				return nil, nil, io.MultiReader(bytes.NewReader(haveRead.Bytes()), reader), err
			}
			return results.Bytes(), fixedResults.Bytes(), reader, nil
		}

		buf := make([]byte, size)
		blockN, err := io.ReadFull(reader, buf)
		results.Write(buf[:blockN])
		haveRead.Write(buf[:blockN])

		fixedResults.Write(lineBytes)
		if commentExisted {
			fixedResults.WriteByte(';')
			fixedResults.Write(comment)
		}
		fixedResults.WriteString("\r\n")
		fixedResults.Write(buf[:blockN])
		fixedResults.WriteString("\r\n")
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				return results.Bytes(), bytes.TrimSpace(fixedResults.Bytes()), reader, err
			} else {
				return nil, nil, io.MultiReader(bytes.NewReader(haveRead.Bytes()), reader), fmt.Errorf("read chunked data error: %v", err)
			}
		}

		endBlock, delim, _ := bufioReadLine(reader)
		haveRead.Write(endBlock)
		haveRead.Write(delim)
		if len(endBlock) != 0 {
			return nil, nil, io.MultiReader(bytes.NewReader(haveRead.Bytes()), reader), fmt.Errorf("read chunked data error: %v", err)
		}
	}
}

func bufioReadLine(reader *bufio.Reader) ([]byte, []byte, error) {
	if reader == nil {
		return nil, nil, errors.New("empty reader(bufio)")
	}

	var lineBuffer bytes.Buffer
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, nil, err
		}
		lineBuffer.WriteByte(b)
		if b == '\n' {
			break
		}
	}

	lines := lineBuffer.Bytes()
	if bytes.HasSuffix(lines, []byte{'\r', '\n'}) {
		return lines[:len(lines)-2], []byte{'\r', '\n'}, nil
	}
	return lines[:len(lines)-1], []byte{'\n'}, nil
}

func HTTPChunkedEncode(raw []byte) []byte {
	var buf bytes.Buffer
	writer := httputil.NewChunkedWriter(&buf)

	maxSplit := len(raw) / 2
	if maxSplit <= 0 {
		maxSplit = 47
	}

	offset := 0
	maxBuffer := 3 + rand.Intn(maxSplit)
	for offset < len(raw) {
		end := offset + maxBuffer
		if end > len(raw) {
			end = len(raw)
		}
		chunk := raw[offset:end]
		writer.Write(chunk)
		offset = end
		maxBuffer = 3 + rand.Intn(maxSplit)
	}

	writer.Close()
	buf.WriteString("\r\n")
	return buf.Bytes()
}
