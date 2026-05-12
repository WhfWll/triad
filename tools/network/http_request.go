package network

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/textproto"
	"strings"
)

var commonHeader map[string]string

func ReadHTTPRequestEx(reader *bufio.Reader, loadbody bool) (*http.Request, error) {
	var buf bytes.Buffer
	req, err := http.ReadRequest(bufio.NewReader(io.TeeReader(reader, &buf)))
	if err != nil {
		return nil, err
	}

	if loadbody {
		var finalBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(finalBody))
	}

	var cache = make(map[string]string)
	// 这里用来恢复 Req 的大小写
	SplitHTTPHeadersAndBodyFromPacket(buf.Bytes(), func(line string) {
		if index := strings.Index(line, ":"); index > 0 {
			key := line[:index]
			ckey := textproto.CanonicalMIMEHeaderKey(key)
			_, ok := commonHeader[ckey]
			// 大小写发生了变化，并且不是常见公共头，则说明需要恢复一下
			if ckey != key && !ok {
				cache[ckey] = key
			}
		}
	})

	for ckey, key := range cache {
		values, ok := req.Header[ckey]
		if ok {
			req.Header[key] = values
			delete(req.Header, ckey)
		}
	}

	//black magic fix when browser use http proxy the RequestURI is not canonical
	if strings.HasPrefix(req.RequestURI, "http://") || strings.HasPrefix(req.RequestURI, "https://") {
		if req.Header.Get("Host") == "" {
			req.Header.Add("Host", req.URL.Host)
		}
		req.RequestURI = req.URL.RequestURI()
	}

	return req, nil
}
