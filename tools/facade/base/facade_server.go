// Package base
// @Author bcy2007  2025/12/23 14:43
package base

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

const (
	LDAPMsgFlag         = "ldap_flag"
	RMIMsgFlag          = "rmi"
	RMIHandshakeMsgFlag = "rmi-handshake"
)

const (
	emptyVerbose         = "<empty>"
	getInfoFailedVerbose = "<get info failed>"
)

type FacadeResourceType interface {
	[]byte | map[string]any | *HttpResource
}
type resourceAndVerbose[T FacadeResourceType] struct {
	Resource T
	Verbose  string
}
type FacadeServerResource[T FacadeResourceType] struct {
	lock      sync.Mutex
	Resources map[string]*resourceAndVerbose[T]
}

func NewFacadeServerResource[T FacadeResourceType]() *FacadeServerResource[T] {
	return &FacadeServerResource[T]{
		Resources: map[string]*resourceAndVerbose[T]{},
	}
}

func (f *FacadeServerResource[T]) ForEachResource(fun func(token string, resource T, verbose string) error) error {
	f.lock.Lock()
	defer func() {
		f.lock.Unlock()
	}()
	for k, r := range f.Resources {
		err := fun(k, r.Resource, r.Verbose)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FacadeServerResource[T]) DeleteResource(token string) {
	f.lock.Lock()
	defer func() {
		f.lock.Unlock()
	}()
	delete(f.Resources, token)
}

func (f *FacadeServerResource[T]) GetResource(token string) (T, string, bool) {
	f.lock.Lock()
	defer func() {
		f.lock.Unlock()
	}()
	if v, ok := f.Resources[token]; ok {
		verbose := v.Verbose
		if v.Verbose == "" {
			verbose = getInfoFailedVerbose
		}
		return v.Resource, verbose, true
	}
	return nil, "", false
}

func (f *FacadeServerResource[T]) SetResource(token string, data T, verbose string) {
	f.lock.Lock()
	defer func() {
		f.lock.Unlock()
	}()
	f.Resources[token] = &resourceAndVerbose[T]{
		Resource: data,
		Verbose:  verbose,
	}
}

type FacadeServer struct {
	cancel func()

	Host         string
	Port         int
	ExternalHost string
	// 反连地址
	ReverseAddr string

	rmiResourceAddrs           *FacadeServerResource[[]byte]
	ldapResourceAddrs          *FacadeServerResource[map[string]any]
	httpResource               *FacadeServerResource[*HttpResource]
	handlers                   []func(notification *Notification)
	RemoteAddrConvertorHandler func(string) string
	// resourceName               string
	ldapEntry map[string]interface{}
	httpMux   *sync.Mutex
}

type ResourcesInfo struct {
	Protocol    string
	Url         string
	Data        any
	DataVerbose string
}

func (f *FacadeServer) GetAllResourcesInfo() []*ResourcesInfo {
	var res []*ResourcesInfo
	f.rmiResourceAddrs.ForEachResource(func(token string, resource []byte, verbose string) error {
		res = append(res, &ResourcesInfo{
			Protocol:    "rmi",
			Url:         fmt.Sprintf("rmi://%s/%s", f.ReverseAddr, token),
			Data:        resource,
			DataVerbose: verbose,
		})
		return nil
	})

	f.ldapResourceAddrs.ForEachResource(func(token string, resource map[string]any, verbose string) error {
		res = append(res, &ResourcesInfo{
			Protocol:    "ldap",
			Url:         fmt.Sprintf("ldap://%s/%s", f.ReverseAddr, token),
			Data:        resource,
			DataVerbose: verbose,
		})
		return nil
	})

	f.httpResource.ForEachResource(func(token string, resource *HttpResource, verbose string) error {
		res = append(res, &ResourcesInfo{
			Protocol:    "http",
			Url:         fmt.Sprintf("http://%s/%s", f.ReverseAddr, token),
			Data:        resource,
			DataVerbose: verbose,
		})
		return nil
	})
	return res
}

type (
	FactoryFun         func() string
	FacadeServerConfig func(f *FacadeServer)
)

func (f *FacadeServer) Config(configs ...FacadeServerConfig) {
	for _, config := range configs {
		config(f)
	}
}

func (f *FacadeServer) ConvertRemoteAddr(addr string) string {
	if f.RemoteAddrConvertorHandler != nil {
		return f.RemoteAddrConvertorHandler(addr)
	}
	return addr
}

func Serve(host string, port int, configs ...FacadeServerConfig) error {
	server := NewFacadeServer(host, port, configs...)
	err := server.ServeWithContext(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func NewFacadeServer(host string, port int, configs ...FacadeServerConfig) *FacadeServer {
	facadeServer := &FacadeServer{
		Host:              host,
		Port:              port,
		ldapEntry:         make(map[string]interface{}),
		httpResource:      NewFacadeServerResource[*HttpResource](),
		rmiResourceAddrs:  NewFacadeServerResource[[]byte](),
		ldapResourceAddrs: NewFacadeServerResource[map[string]any](),
		httpMux:           &sync.Mutex{},
	}
	facadeServer.rmiResourceAddrs.SetResource("", nil, emptyVerbose)
	facadeServer.httpResource.SetResource("", NewHttpRawResource([]byte(defaultHTTPFallback)), emptyVerbose)
	facadeServer.ldapResourceAddrs.SetResource("", nil, emptyVerbose)
	for _, config := range configs {
		config(facadeServer)
	}
	return facadeServer
}

func (f *FacadeServer) GetAddr() string {
	return fmt.Sprintf("%s:%d", f.Host, f.Port)
}

func (f *FacadeServer) OnHandle(h func(n *Notification)) {
	f.handlers = append(f.handlers, h)
}

func (f *FacadeServer) triggerNotification(t string, conn net.Conn, token string, raw []byte) {
	f.triggerNotificationEx(t, conn, token, raw, "")
}

func (f *FacadeServer) triggerNotificationEx(t string, conn net.Conn, token string, raw []byte, responseInfo string) {
	remoteAddr := f.ConvertRemoteAddr(conn.RemoteAddr().String())

	if token == "" {
		log.Infof("trigger %v from %v", t, f.ConvertRemoteAddr(conn.RemoteAddr().String()))
	} else {
		log.Infof("trigger %v[%v] from %v", t, token, f.ConvertRemoteAddr(conn.RemoteAddr().String()))
	}

	notif := NewNotification(t, remoteAddr, raw, token)
	// 通过conn的地址计算hash（因为每次连接的conn都是独立的对象，所以可以用conn地址的hash区分不同连接）
	// 通过这个hash去判断是否是同一个连接，如果是同一个连接，则在原通知基础上进行更新，否则新增通知
	notif.ConnectHash = Md5(fmt.Sprintf("%p", conn))
	// 响应内容
	notif.ResponseInfo = responseInfo
	if len(f.handlers) <= 0 {
		// spew.Dump(notif)
	}

	for _, handle := range f.handlers {
		if handle == nil {
			return
		}
		handle(notif)
	}
}

func (f *FacadeServer) Serve() error {
	return f.ServeWithContext(context.Background())
}

func (f *FacadeServer) CancelServe() {
	if f.cancel != nil {
		f.cancel()
	}
}

func (f *FacadeServer) ServeWithContext(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	f.cancel = cancel
	lis, err := net.Listen("tcp", HostPort(f.Host, f.Port))
	if f.Port == 0 {
		f.Port = lis.Addr().(*net.TCPAddr).Port
	}
	log.Infof("start to listen reverse(facade) on: %s:%d", f.Host, f.Port)

	if err != nil {
		return fmt.Errorf("create listen failed: %s", err)
	}

	go func() {
		<-ctx.Done()
		lis.Close()
	}()

	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}

		go func() {
			defer conn.Close()
			defer func() {
				if err := recover(); err != nil {
					log.Error(err)
					return
				}
			}()
			f.triggerNotification("tcp", conn, "", nil)
			log.Infof("recv conn from: %s", conn.RemoteAddr())
			f.handleConn(conn)
		}()
	}
}

func (f *FacadeServer) handleConn(conn net.Conn) {
	peekableConn := NewPeekableNetConn(conn)
	log.Infof("start to fallback http handlers for: %s", conn.RemoteAddr())
	err := f.getHTTPHandler(false)(peekableConn)
	if err != nil {
		log.Errorf("handle http failed: %s", err)
		return
	}

}
