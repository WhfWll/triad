// Package base
// @Author bcy2007  2025/12/23 15:31
package base

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gobwas/glob"
	"github.com/klauspost/compress/zstd"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"io"
	"io/ioutil"
	"mime"
	"regexp"
	"smart/tools/go-funk"
	"smart/tools/mimetype"
	"smart/tools/mimetype/mimeutil/mimecharset"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var (
	expect100continue     = []byte("HTTP/1.1 100 Continue\r\n\r\n")
	contentEncodingRegexp = regexp.MustCompile(`(?i)content-encoding:\s*\w*\r?\n`)
	isChunkedBytes        = []byte("\r\n0\r\n\r\n")
	charsetRegexp         = regexp.MustCompile(`(?i)charset\s*=\s*"?\s*([^\s;\n\r"]+)`)
	contentTypeRegexp     = regexp.MustCompile(`(?i)content-type:\s*([^\r\n]*)`)
)

func FixHTTPResponse(raw []byte) (rsp []byte, body []byte, _ error) {
	raw, _ = bytes.CutPrefix(raw, expect100continue)

	isChunked := false
	// 这两个用来处理编码特殊情况
	var contentEncoding string
	var contentType string
	noContentTypeSet := true
	headers, body := SplitHTTPHeadersAndBodyFromPacket(raw, func(line string) {
		if strings.HasPrefix(strings.ToLower(line), "content-type:") {
			_, contentType = SplitHTTPHeader(line)
			noContentTypeSet = false
		}
		// 判断内容
		line = strings.ToLower(line)
		if strings.HasPrefix(line, "transfer-encoding:") && IContains(line, "chunked") {
			isChunked = true
		}
		if strings.HasPrefix(line, "content-encoding:") {
			contentEncoding = line
		}
	})
	if headers == "" {
		return nil, nil, fmt.Errorf("error for parsing http response")
	}
	headerBytes := []byte(headers)

	bodyRaw := body
	if bodyRaw != nil && isChunked {
		unchunked, chunkErr := HTTPChunkedDecode(bodyRaw)
		if unchunked != nil {
			bodyRaw = unchunked
		} else {
			if chunkErr == nil {
				bodyRaw = []byte{}
			}
		}
	}
	if contentEncoding != "" {
		decodedBodyRaw, fixed := ContentEncodingDecode(contentEncoding, bodyRaw)
		if decodedBodyRaw != nil && fixed {
			// contents get decoded
			headerBytes = RemoveCEHeaders(headerBytes)
			bodyRaw = decodedBodyRaw
		}
	}

	if len(bodyRaw) == 0 {
		return ReplaceHTTPPacketBodyEx(headerBytes, bodyRaw, false, true), bodyRaw, nil
	}
	mimeResult, err := MatchMIMEType(bodyRaw)
	if err != nil {
		log.Warnf("match mime type failed: %v", err)
		return ReplaceHTTPPacketBodyEx(headerBytes, bodyRaw, false, true), bodyRaw, nil
	}

	// 记录原始 contentType
	originContentType := contentType

	var bodyChanged bool
RetryContentType:
	switch {
	case IsTextPlainMIMEType(contentType):
		fallthrough
	case IsJsonMIMEType(contentType):
		fallthrough
	case IsJavaScriptMIMEType(contentType):
		bodyRaw, bodyChanged = mimeResult.TryUTF8Convertor(bodyRaw)
		if bodyChanged {
			if strings.Contains(strings.ToLower(originContentType), "charset=") {
				newContentType := charsetRegexp.ReplaceAllString(originContentType, "charset=utf-8")
				if strings.ToLower(newContentType) != strings.ToLower(originContentType) {
					log.Infof("auto fix content-type via utf convertor auto, from %#v to %#v", originContentType, newContentType)
					headerBytes = ReplaceMIMEType(headerBytes, newContentType)
				}
			}
		}
		return ReplaceHTTPPacketBodyEx(headerBytes, bodyRaw, false, true), bodyRaw, nil
	case IsHtmlOrXmlMIMEType(contentType):
		// body is not text, but content-type is ...
		// fix content-type header
		if !IsHtmlOrXmlMIMEType(mimeResult.MIMEType) && !IsTextPlainMIMEType(mimeResult.MIMEType) && !mimeResult.IsText && mimeResult.MIMEType != "application/octet-stream" {
			log.Warnf("origin content-type: %v(%v), fix new content-type: %v, reason: the actually body is not text...", contentType, originContentType, mimeResult.MIMEType)
			contentType = mimeResult.MIMEType
			goto RetryContentType
		}

		var after []byte
		var containsUTF8 bool
		after, bodyChanged = mimeResult.TryUTF8Convertor(bodyRaw)
		if bodyChanged {
			containsUTF8 = true
			log.Infof("HtmlOrXmlMIMEType(%#v) auto fix body, origin: len(%v) -> len(%v)", originContentType, len(bodyRaw), len(after))
			bodyRaw = after
		}

		var newContentType string
		origin, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			newContentType = contentType
		} else {
			// 如果服务端返回了 charset，并且转换utf-8成功，就直接覆盖设置 charset=utf-8，否则使用服务端设置的 charset
			if containsUTF8 {
				params["charset"] = "utf-8"
			}
			newContentType = mime.FormatMediaType(origin, params)
		}
		headerBytes = ReplaceMIMEType(headerBytes, newContentType)
		return ReplaceHTTPPacketBodyEx(headerBytes, bodyRaw, false, true), bodyRaw, nil
	default:
		if mimeResult == nil || mimeResult.MIMEType == "" {
			return ReplaceHTTPPacketBodyEx(headerBytes, bodyRaw, false, true), bodyRaw, nil
		}

		if contentType == "" && noContentTypeSet {
			contentType = mimeResult.MIMEType
			goto RetryContentType
		}

		if strings.HasPrefix(strings.ToLower(contentType), "text/") {
			bodyRaw, bodyChanged = mimeResult.TryUTF8Convertor(bodyRaw)
			if bodyChanged {
				withoutCharset, _, err := mime.ParseMediaType(contentType)
				if err != nil {
					withoutCharset = contentType
				} else {
					contentType = withoutCharset
				}
				headerBytes = ReplaceMIMEType(headerBytes, mime.FormatMediaType(contentType, map[string]string{"charset": "utf-8"}))
			}
			return ReplaceHTTPPacketBodyEx(headerBytes, bodyRaw, false, true), bodyRaw, nil
		} else {
			if !mimeResult.IsText {
				headerBytes = ReplaceMIMEType(headerBytes, mimeResult.MIMEType)
				return ReplaceHTTPPacketBodyEx(headerBytes, bodyRaw, false, true), bodyRaw, nil
			}
			return ReplaceHTTPPacketBodyEx(headerBytes, bodyRaw, false, true), bodyRaw, nil
		}
	}
}

func SplitHTTPHeadersAndBodyFromPacket(raw []byte, hook ...func(line string)) (headers string, body []byte) {
	return SplitHTTPHeadersAndBodyFromPacketEx(raw, nil, hook...)
}

func SplitHTTPHeadersAndBodyFromPacketEx(raw []byte, mf func(method string, requestUri string, proto string) error, hook ...func(line string)) (string, []byte) {
	if len(hook) > 0 {
		return SplitHTTPPacket(raw, mf, nil, func(line string) (ret string) {
			ret = line
			defer func() {
				if err := recover(); err != nil {
					PrintCurrentGoroutineRuntimeStack()
				}
				ret = line
			}()
			for _, h := range hook {
				h(line)
			}
			return ret
		})
	}
	return SplitHTTPPacket(raw, mf, nil)
}

func SplitHTTPPacket(
	raw []byte,
	reqFirstLine func(method string, requestUri string, proto string) error,
	rspFirstLine func(proto string, code int, codeMsg string) error,
	hook ...func(line string) string,
) (string, []byte) {
	return SplitHTTPPacketEx(raw, reqFirstLine, rspFirstLine, nil, hook...)
}

func SplitHTTPPacketEx(
	raw []byte,
	reqFirstLine func(method string, requestUri string, proto string) error,
	rspFirstLine func(proto string, code int, codeMsg string) error,
	rawFistLine func(string) error,
	hook ...func(line string) string,
) (string, []byte) {
	reader := bufio.NewReader(bytes.NewBuffer(raw))
	firstLineBytes, err := BufioReadLine(reader)
	if err != nil {
		return "", nil
	}
	prefix, firstLineBytes, _ := CutBytesPrefixFunc(firstLineBytes, NotSpaceRune)
	firstLineBytes = TrimSpaceHTTPPacket(firstLineBytes)
	if rawFistLine != nil {
		err := rawFistLine(string(firstLineBytes))
		if err != nil {
			log.Debugf("rawFistLine error: %s", err)
			return "", nil
		}
	}
	var isResp = bytes.HasPrefix(firstLineBytes, []byte("HTTP/")) || bytes.HasPrefix(firstLineBytes, []byte("RTSP/"))
	if isResp {
		// rsp
		if rspFirstLine != nil {
			proto, code, codeMsg, _ := ParseHTTPResponseLine(string(firstLineBytes))
			err := rspFirstLine(proto, code, codeMsg)
			if err != nil {
				log.Debugf("rspHeader error: %s", err)
				return "", nil
			}
		}
	} else {
		// req
		if reqFirstLine != nil {
			method, requestURI, proto, _ := ParseHTTPRequestLine(string(firstLineBytes))
			err := reqFirstLine(method, requestURI, proto)
			if err != nil && err.Error() != "normal abort" {
				log.Debugf("reqHeader error: %s", err)
				return "", nil
			}
		}
	}

	var headers []string
	headers = append(headers, string(firstLineBytes))
	haveCl := false
	err = ScanHTTPHeader(reader, func(rawHeader []byte) {
		if len(rawHeader) == 0 {
			return
		}
		line := string(rawHeader)
		skipHeader := false
		for _, h := range hook {
			hooked := h(line)
			if hooked == "" {
				skipHeader = true
			}
			if skipHeader {
				break
			}
			line = hooked
		}
		if skipHeader {
			return
		}
		k, _ := SplitHTTPHeader(line)
		if strings.ToLower(k) == "content-length" {
			haveCl = true
		}
		headers = append(headers, line)
	}, prefix, isResp)
	headersRaw := strings.Join(headers, CRLF) + CRLF + CRLF
	bodyRaw, _ := ioutil.ReadAll(reader)
	if bodyRaw == nil {
		return headersRaw, nil
	}

	if len(bytes.TrimSpace(bodyRaw)) == 0 && !haveCl {
		bodyRaw = nil
	}

	return headersRaw, bodyRaw
}

func TrimSpaceHTTPPacket(raw []byte) []byte {
	return bytes.TrimFunc(raw, unicode.IsSpace)
}

func SplitHTTPHeader(i string) (string, string) {
	if ret := strings.Index(i, ":"); ret < 0 {
		return i, ""
	} else {
		key := i[:ret]
		value := strings.TrimSpace(i[ret+1:])
		return key, value
	}
}

func HTTPChunkedDecode(raw []byte) ([]byte, error) {
	if ret := string(raw); ret == "" {
		return nil, errors.New("empty input")
	} else if ret == "0\r\n\r\n" {
		return nil, nil
	}

	results, _, rest, err := ReadHTTPChunkedDataWithFixedError(raw)
	_ = rest
	if len(results) > 0 {
		return results, nil
	}
	if len(raw) > 128 {
		raw = append(raw[:128], []byte("...")...)
	}
	return nil, errors.Errorf("parse %v to http chunked failed: %v", strconv.Quote(string(raw)), err)
}

// ParseHTTPResponseLine parses `HTTP/1.1 200 OK` into its ports
func ParseHTTPResponseLine(line string) (string, int, string, bool) {
	line = strings.TrimSpace(line)

	var proto string
	var code int
	var status string

	blocks := strings.SplitN(line, " ", 3)
	lenOfBlocks := len(blocks)
	if lenOfBlocks > 0 {
		proto = blocks[0]
	}
	if lenOfBlocks > 1 {
		code = Atoi(blocks[1])
	}
	if lenOfBlocks > 2 {
		status = blocks[2]
	}
	return proto, code, status, code != 0
}

func ContentEncodingDecode(contentEncoding string, bodyRaw []byte) (finalResult []byte, fixed bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("handle content-encoding decode failed! reason: %s", err)
			finalResult = bodyRaw
			fixed = false
		}
	}()

	switch true {
	case IContains(contentEncoding, "gzip"):
		// 假设在这里已经把 chunked 解决了
		if bytes.HasPrefix(bodyRaw, []byte{0x1f, 0x8b, 0x08}) {
			ungzipedRaw, err := gzip.NewReader(bytes.NewBuffer(bodyRaw[:]))
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				log.Warnf("uncompressed gzip failed: %s", err)
			}
			if ungzipedRaw != nil {
				raw, err := ioutil.ReadAll(ungzipedRaw)
				if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
					log.Errorf("read ungzip reader failed: %s", err)
				}
				if raw != nil {
					return raw, true
				}
			}
		}
		return bodyRaw, false
	case IContains(contentEncoding, "br"):
		raw, err := ioutil.ReadAll(brotli.NewReader(bytes.NewBuffer(bodyRaw)))
		if err != nil {
			log.Errorf("read[brotli] decode failed: %s", err)
			return bodyRaw, false
		}
		return raw, true
	case IContains(contentEncoding, "compress"):
		log.Errorf("Content-Encoding: compress is not supported")
		return bodyRaw, false
	case IContains(contentEncoding, "deflate"):
		rawReader, err := zlib.NewReader(bytes.NewBuffer(bodyRaw))
		if err != nil {
			decodedBody, _ := io.ReadAll(flate.NewReader(bytes.NewBuffer(bodyRaw)))
			if decodedBody != nil {
				return decodedBody, true
			}
			return bodyRaw, false
		}
		raw, err := ioutil.ReadAll(rawReader)
		if err != nil {
			return bodyRaw, false
		}
		return raw, true
	case IContains(contentEncoding, "zstd"):
		reader, err := zstd.NewReader(bytes.NewBuffer(bodyRaw))
		if err != nil {
			log.Errorf("read[zstd] new reader failed: %s", err)
		}
		raw, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Errorf("read[zstd] decode failed: %s", err)
			log.Infof("bodyRaw: %v", bodyRaw)
			return bodyRaw, false
		}
		return raw, true
	case IContains(contentEncoding, "identity"):
		fallthrough
	case IContains(contentEncoding, "*"):
		fallthrough
	default:
		return bodyRaw, false
	}
}

func RemoveCEHeaders(headerBytes []byte) []byte {
	return contentEncodingRegexp.ReplaceAll(headerBytes, []byte{})
}

func ReplaceHTTPPacketBodyEx(raw []byte, body []byte, chunk bool, forceCL bool) []byte {
	isChunked := false
	var firstLine string
	var headers []string
	_, _ = SplitHTTPPacketEx(raw, nil, nil, func(rawLine string) error {
		firstLine = rawLine
		return nil
	}, func(line string) string {
		if IHasPrefix(line, "transfer-encoding:") && IContains(line, "chunked") {
			isChunked = true
			return line
		}

		if IHasPrefix(line, "content-length") {
			return line
		}
		headers = append(headers, line)
		return line
	})
	headers = append([]string{firstLine}, headers...)
	var buf bytes.Buffer
	// 空 body
	if body == nil {
		buf.WriteString(strings.Join(headers, CRLF) + CRLF + CRLF)
		return buf.Bytes()
	}

	// 只有包含了Transfer-Encoding: chunked，以及body符合chunked格式，才认为已经是chunked
	if isChunked {
		isChunked = bytes.Contains(body, isChunkedBytes)
	}

	// chunked
	if chunk {
		headers = append(headers, "Transfer-Encoding: chunked")
		if !isChunked {
			body = HTTPChunkedEncode(body)
		}
	} else if isChunked {
		newBody, err := HTTPChunkedDecode(body)
		if err == nil {
			body = newBody
		}
	}
	if !chunk && (len(body) > 0 || forceCL) {
		headers = append(headers, fmt.Sprintf("Content-Length: %d", len(body)))
	}
	buf.WriteString(strings.Join(headers, CRLF))
	buf.WriteString(CRLF + CRLF)
	buf.Write(body)
	return buf.Bytes()
}

func MatchMIMEType(raw any) (*MIMEResult, error) {
	r := mimetype.Detect(interfaceToBytes(raw))
	if r == nil {
		return nil, fmt.Errorf("match(detect) mime type failed, check: %v", ShrinkString(fmt.Sprintf("%#v", raw), 64))
	}
	var result = &MIMEResult{
		MIMEType:    r.String(),
		IsText:      mimeIsText(r),
		NeedCharset: r.NeedCharset(),
		Charset:     r.Charset(),
	}
	return result, nil
}

func _mimeIsText(depth int, t *mimetype.MIME) bool {
	if depth > 20 || t == nil {
		return false
	}
	if strings.HasPrefix(t.String(), "text/") || strings.HasPrefix(t.String(), "Text/") {
		return true
	}
	return _mimeIsText(depth+1, t.Parent())
}

func mimeIsText(t *mimetype.MIME) bool {
	return _mimeIsText(0, t)
}

type MIMEResult struct {
	MIMEType    string
	IsText      bool
	NeedCharset bool
	Charset     string
}

type PrescanResult struct {
	Encoding encoding.Encoding
	Name     string
}

func (t *MIMEResult) IsChineseCharset() bool {
	switch strings.ToLower(t.Charset) {
	case "gb18030", "gb-18030", "gbk", "gb2312", "gb-2312":
		return true
	}
	return false
}

func (t *MIMEResult) TryUTF8Convertor(raw []byte) ([]byte, bool) {
	result, ok := t._tryUTF8Convertor(raw)
	if ok {
		if bytes.Contains(result, []byte{'\xef', '\xbf', '\xbd'}) {
			return raw, false
		}
		return result, true
	}
	return raw, false
}

func (t *MIMEResult) _tryUTF8Convertor(raw []byte) ([]byte, bool) {
	if strings.Contains(t.MIMEType, "/html") || strings.Contains(t.MIMEType, "/xhtml+xml") {
		result := raw
		// <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
		// <meta charset="UTF-8">
		// <meta http-equiv="Content-Type" content="text/html; charset=gb2312">
		// <meta charset="gb2312">
		newBuffer := new(bytes.Buffer)
		lastStart := -1
		var encodings []PrescanResult
		var set = make(map[string]struct{})
		enc, origin := HtmlCharsetPrescan(result, func(start, end int, match PrescanResult) {
			if _, ok := set[match.Name]; !ok {
				encodings = append(encodings, match)
				set[match.Name] = struct{}{}
			}
			if lastStart < 0 {
				newBuffer.Write(result[:start])
			} else {
				newBuffer.Write(result[lastStart:start])
			}
			newBuffer.WriteString("utf-8")
			lastStart = end
		})
		if strings.ToLower(origin) != "utf-8" && lastStart >= 0 {
			newBuffer.Write(result[lastStart:])
			result = newBuffer.Bytes()
		}

		if len(encodings) == 1 {
			if encodings[0].Name == "utf-8" {
				return result, false
			}

			decodedResult, err := enc.NewDecoder().Bytes(result)
			if err != nil {
				return result, false
			}
			return decodedResult, true
		} else if len(encodings) > 1 {
			log.Warnf("WARNING: ATTENTION multiple encodings [%v], try the best", funk.Keys(set))
			for _, v := range encodings {
				if v.Encoding != nil {
					decodeResult, err := v.Encoding.NewDecoder().Bytes(result)
					if err != nil {
						log.Infof("try encoding %#v failed: %v", v.Name, err)
						continue
					}
					return decodeResult, true
				}
			}
			return result, false
		} else {
			// no meta encoding, treat like plain text
			charsetFallback := mimecharset.FromPlain(result)
			enc, charsetFallback := charset.Lookup(charsetFallback)
			if !lo.Contains([]string{
				"utf-8", "utf8", "windows-1252", "iso-8859-1",
			}, charsetFallback) && enc != nil {
				decodedResult, err := enc.NewDecoder().Bytes(result)
				if err == nil {
					return decodedResult, true
				}
			}
		}
	}

	switch charsetLower := strings.ToLower(t.Charset); charsetLower {
	case "gb18030", "gb-18030", "gbk", "gb2312", "gb-2312":
		result, err := GB18030ToUtf8(raw)
		if err != nil {
			return raw, false
		}
		return result, true
	default:
		if t.MIMEType == "application/octet-stream" {
			// application/octet-stream is not text, but binary
			return raw, false
		}

		if charsetLower != "" && charsetLower != "utf-8" {
			log.Warnf("TBD: charset %#v not supported yet, use origin raw input", t.Charset)
		}

		if charsetLower == "" && t.IsText {
			charsetLower = mimecharset.FromPlain(raw)
			enc, _ := charset.Lookup(charsetLower)
			if enc != nil {
				fixed, err := enc.NewDecoder().Bytes(raw)
				if err == nil {
					return fixed, true
				}
			}
		}
	}
	return raw, false
}

var (
	gb18030encoding      encoding.Encoding
	gb18030encodingMutex = new(sync.Mutex)
)

func GB18030ToUtf8(s []byte) ([]byte, error) {
	if gb18030encoding != nil {
		return gb18030encoding.NewDecoder().Bytes(s)
	}

	gb18030encodingMutex.Lock()
	defer gb18030encodingMutex.Unlock()

	if gb18030encoding != nil {
		return gb18030encoding.NewDecoder().Bytes(s)
	}
	var name string
	gb18030encoding, name = charset.Lookup("gb18030")
	if gb18030encoding == nil {
		return nil, fmt.Errorf("failed to lookup gb18030 encoding: %s", name)
	}
	return gb18030encoding.NewDecoder().Bytes(s)
}

var metaCharset = regexp.MustCompile(`(?i)<\s*?meta[^>]*?charset\s*=\s*['"]?\s*([a-z-0-9]*)['"]?`)

func HtmlCharsetPrescan(content []byte, callback ...func(start, end int, matched PrescanResult)) (e encoding.Encoding, name string) {
	reader := bytes.NewReader(content)

	z := html.NewTokenizer(reader)
	var count = 0
	var tagCount = 0
	var startOffset int64 = 0
	var endOffset int64 = 0

	var finalResults []PrescanResult

	fallback := func() (encoding.Encoding, string) {
		if len(finalResults) == 0 {
			return nil, ""
		}
		last := finalResults[0]
		return last.Encoding, last.Name
	}

	for {
		count++
		if count > 800 {
			return fallback()
		}
		ret := z.Next()
		if tagCount > 200 {
			return fallback()
		}
		switch ret {
		case html.ErrorToken:
			return fallback()
		case html.EndTagToken:
			name, _ := z.TagName()
			if bytes.Equal(bytes.ToLower(name), []byte("head")) {
				return fallback()
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			tagCount++
			endOffset, _ = reader.Seek(0, io.SeekCurrent)
			endOffset -= int64(len(z.Buffered()))

			tagName, hasAttr := z.TagName()
			tagName = bytes.ToLower(tagName)
			if !bytes.Equal(tagName, []byte("meta")) {
				if bytes.Equal(tagName, []byte("head")) {
					startOffset, _ = reader.Seek(0, io.SeekCurrent)
					startOffset -= int64(len(z.Buffered()))
					endOffset = 0
				}
				continue
			}

			attrList := make(map[string]struct{})
			gotPragma := false

			const (
				dontKnow = iota
				doNeedPragma
				doNotNeedPragma
			)
			needPragma := dontKnow
			name = ""
			e = nil
			for hasAttr {
				var key, val []byte
				key, val, hasAttr = z.TagAttr()
				ks := string(key)
				if _, ok := attrList[ks]; ok {
					continue
				}

				if bytes.EqualFold(val, []byte("gb-18030")) {
					val = []byte("gb18030")
				}
				attrList[ks] = struct{}{}
				for i, c := range val {
					if 'A' <= c && c <= 'Z' {
						val[i] = c + 0x20
					}
				}

				switch ks {
				case "http-equiv":
					if bytes.Equal(bytes.ToLower(val), []byte("content-type")) {
						gotPragma = true
					}
				case "content":
					if e == nil {
						name = charsetFromMetaElement(string(val))
						if name != "" {
							e, name = charset.Lookup(name)
							if e != nil {
								needPragma = doNeedPragma
							}
						}
					}

				case "charset":
					valname := string(val)
					e, name = charset.Lookup(valname)
					needPragma = doNotNeedPragma
				}
			}

			if needPragma == dontKnow || needPragma == doNeedPragma && !gotPragma {
				continue
			}

			if strings.HasPrefix(name, "utf-16") {
				name = "utf-8"
				e = encoding.Nop
			}

			if e != nil {
				pRes := PrescanResult{
					Encoding: e, Name: name,
				}

				// gbk -> gb18030
				// gb2312 -> gb18030
				if strings.EqualFold(name, "gbk") || strings.EqualFold(name, "gb2312") {
					pRes.Encoding, pRes.Name = charset.Lookup("gb18030")
				}

				if endOffset > startOffset && startOffset >= 0 {
					for _, cb := range callback {
						result := metaCharset.FindSubmatchIndex(content[startOffset:endOffset])
						if len(result) > 3 {
							offset := int(startOffset)
							cb(offset+result[2], offset+result[3], pRes)
						}
					}
					startOffset = endOffset
				}
				finalResults = append(finalResults, pRes)
			}
		}
	}
}

func charsetFromMetaElement(s string) string {
	for s != "" {
		csLoc := strings.Index(s, "charset")
		if csLoc == -1 {
			return ""
		}
		s = s[csLoc+len("charset"):]
		s = strings.TrimLeft(s, " \t\n\f\r")
		if !strings.HasPrefix(s, "=") {
			continue
		}
		s = s[1:]
		s = strings.TrimLeft(s, " \t\n\f\r")
		if s == "" {
			return ""
		}
		if q := s[0]; q == '"' || q == '\'' {
			s = s[1:]
			closeQuote := strings.IndexRune(s, rune(q))
			if closeQuote == -1 {
				return ""
			}
			return s[:closeQuote]
		}

		end := strings.IndexAny(s, "; \t\n\f\r")
		if end == -1 {
			end = len(s)
		}
		return s[:end]
	}
	return ""
}

func ShrinkString(r any, size int) string {
	return shrinkStringWithMultiLine(r, size, false)
}

func ShrinkTextBlock(r any, size int) string {
	return shrinkStringWithMultiLine(r, size, true)
}

func shrinkStringWithMultiLine(r any, size int, multiline bool) string {
	if size <= 6 {
		size = 10
	}

	half := size / 2

	verbose := AnyToString(r)
	verbose = strings.TrimSpace(verbose)
	runes := []rune(verbose)
	if len(runes) > size {
		runes = append(runes[:half], append([]rune("..."), runes[len(runes)-half:]...)...)
		verbose = string(runes)
	}
	if !multiline {
		verbose = strconv.Quote(verbose)
		verbose = verbose[1:]
		verbose = verbose[:len(verbose)-1]
		verbose = strings.ReplaceAll(verbose, `\r`, " ")
		verbose = strings.ReplaceAll(verbose, `\n`, " ")
		verbose = strings.ReplaceAll(verbose, `\t`, " ")
		verbose = strings.ReplaceAll(verbose, `\"`, "\"")
	}
	return verbose
}

var (
	textPlainMIMEGlob = []glob.Glob{
		glob.MustCompile(`text/plain`),
	}
	jsonMIMEGlobs = []glob.Glob{
		glob.MustCompile(`application/json`),
		glob.MustCompile(`application/*json*`),
		glob.MustCompile(`text/*json*`),
	}
	jsMIMEGlobs = []glob.Glob{
		glob.MustCompile(`application/*javascript*`),
		glob.MustCompile(`text/*javascript*`),
		glob.MustCompile(`application/*ecmascript*`),
		glob.MustCompile(`text/*ecmascript*`),
		glob.MustCompile(`text/jscript`),
	}
	htmlMIMEGlob = []glob.Glob{
		glob.MustCompile(`text/html`),
		glob.MustCompile(`application/xhtml+xml`),
		glob.MustCompile(`application/html`),
		glob.MustCompile(`text/x-html`),
		glob.MustCompile(`application/xml`),
		glob.MustCompile(`text/xml`),
		glob.MustCompile(`application/xhtml`),
		glob.MustCompile(`application/*html*`),
		glob.MustCompile(`text/*html*`),
	}
)

func IsTextPlainMIMEType(s string) bool {
	if s == "" {
		return false
	}

	if strings.Contains(strings.ToLower(s), "charset=") {
		lake, _, err := mime.ParseMediaType(s)
		if err == nil {
			s = lake
		}
	}

	for _, g := range textPlainMIMEGlob {
		if g.Match(s) {
			return true
		}
	}
	return false
}

func IsJsonMIMEType(s string) bool {
	if s == "" {
		return false
	}
	for _, g := range jsonMIMEGlobs {
		if g.Match(s) {
			return true
		}
	}
	return false
}

func IsJavaScriptMIMEType(s string) bool {
	if s == "" {
		return false
	}
	for _, g := range jsMIMEGlobs {
		if g.Match(s) {
			return true
		}
	}
	return false
}

func IsHtmlOrXmlMIMEType(s string) bool {
	if s == "" {
		return false
	}

	if strings.Contains(strings.ToLower(s), "charset=") {
		lake, _, err := mime.ParseMediaType(s)
		if err == nil {
			s = lake
		}
	}

	for _, g := range htmlMIMEGlob {
		if g.Match(s) {
			return true
		}
	}
	return false
}

func ReplaceMIMEType(headerBytes []byte, mimeType string) []byte {
	if mimeType == "" {
		return headerBytes
	}

	idxs := contentTypeRegexp.FindSubmatchIndex(headerBytes)
	if len(idxs) > 3 {
		buf := bytes.NewBuffer(nil)
		buf.Write(headerBytes[:idxs[2]])
		buf.WriteString(mimeType)
		buf.Write(headerBytes[idxs[3]:])
		return buf.Bytes()
	} else {
		return AppendHeaderToHTTPPacket(headerBytes, "Content-Type: "+mimeType)
	}
}

func AppendHeaderToHTTPPacket(raw []byte, line string) []byte {
	header, body := SplitHTTPHeadersAndBodyFromPacket(raw)
	header = strings.TrimRight(header, "\r\n") + CRLF + strings.TrimSpace(line) + CRLF + CRLF
	return []byte(header + string(body))
}

func ReplaceHTTPPacketBody(raw []byte, body []byte, chunk bool) (newHTTPRequest []byte) {
	return ReplaceHTTPPacketBodyEx(raw, body, chunk, false)
}
