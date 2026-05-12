// Package base
// @Author bcy2007  2025/12/23 14:39
package base

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type BufferedPeekableConn struct {
	net.Conn
	buf []byte
}

func (b *BufferedPeekableConn) GetOriginConn() net.Conn {
	return b.Conn
}

func (b *BufferedPeekableConn) Peek(i int) ([]byte, error) {
	return _peekablePeek(b, i)
}

func (b *BufferedPeekableConn) PeekByte() (byte, error) {
	buf, err := b.Peek(1)
	if err != nil {
		return 0, err
	}
	if len(buf) != 1 {
		return 0, io.EOF
	}
	return buf[0], nil
}

func (b *BufferedPeekableConn) PeekUint16() uint16 {
	buf, err := b.Peek(2)
	if err != nil {
		return 0
	}
	if len(buf) != 2 {
		return 0
	}
	return uint16(buf[0])<<8 | uint16(buf[1])
}

func (b *BufferedPeekableConn) Read(buf []byte) (int, error) {
	return _peekableRead(b, buf)
}

func (b *BufferedPeekableConn) GetReader() io.Reader {
	return b.Conn
}

func (b *BufferedPeekableConn) SetBuf(buf []byte) {
	b.buf = buf
}

func (b *BufferedPeekableConn) GetBuf() []byte {
	return b.buf
}

type bufferable interface {
	GetBuf() []byte
	GetReader() io.Reader
	SetBuf([]byte)
}

func _peekableRead(p bufferable, b []byte) (int, error) {
	l := len(p.GetBuf())
	if l <= 0 {
		return p.GetReader().Read(b)
	}
	rl := len(b)
	if rl <= l {
		copy(b, p.GetBuf()[:rl])
		p.SetBuf(p.GetBuf()[rl:])
		return rl, nil
	}
	if l > 0 {
		n1 := copy(b, p.GetBuf())
		p.SetBuf(p.GetBuf()[n1:])
		n2, err := p.GetReader().Read(b[n1:])
		return n1 + n2, err
	}
	return p.GetReader().Read(b)
}

func _peekablePeek(r bufferable, i int) (_ []byte, fErr error) {
	defer func() {
		if err := recover(); err != nil {
			log.Infof("peekable failed: %s", err)
			fErr = io.EOF
		}
	}()

	var buf = make([]byte, i)
	l := len(r.GetBuf())
	if i <= l {
		copy(buf, r.GetBuf()[:i])
		return buf, nil
	} else {
		copy(buf, r.GetBuf())
		var a = r.GetReader()
		if a == nil {
			return nil, io.EOF
		}
		n, err := a.Read(buf[l:])
		r.SetBuf(buf[:l+n])
		return r.GetBuf(), err
	}
}

func NewPeekableNetConn(r net.Conn) *BufferedPeekableConn {
	return &BufferedPeekableConn{
		Conn: r,
	}
}
