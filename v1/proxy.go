package vproxy

import (
	"net/http"
    "net"
    "time"
    "io"
)


const defaultDataBufioSize    = 1<<20                                                       // 默认数据缓冲1MB

//Config 配置
type Config struct {
    DataBufioSize int                                                                       // 缓冲区大小
}

type Proxy struct {
    Addr        string                                                                      // 代理IP地址
    Server      http.Server                                                                 // 服务器
    Transport   http.RoundTripper                                                           // 代理
    *Config                                                                                 // 配置
    l           net.Listener                                                                // 连接对象
}

//setDefault 设置默认
func (p *Proxy) setDefault(){
    if p.Transport == nil {
        p.Transport = http.DefaultTransport
    }
}

//initServer 初始化服务器
func (p *Proxy) initServer() *http.Server {
    srv := &p.Server
    if srv.Handler == nil {
        srv.Handler = http.HandlerFunc(p.ServeHTTP)
    }
    return srv
}

//ServeHTTP 处理服务
//  参：
//      rw http.ResponseWriter  响应
//      req *http.Request       请求
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request){
    switch req.Method {
    	case "CONNECT":
            cp := &connectProxy{
                config      : p.Config,
                transport   : p.Transport,
            }
            cp.ServeHTTP(rw, req)
    	default:
            hp := &httpProxy{
                config      : p.Config,
                transport   : p.Transport,
            }
            hp.ServeHTTP(rw, req)
    }
}

//ListenAndServ 开启监听
//  返：
//      error       错误
func (p *Proxy) ListenAndServ() error {
    p.setDefault()
    srv := p.initServer()
    addr := p.Addr
    if addr == "" {
        addr = ":0"
    }
    l, err := net.Listen("tcp", addr)
    if err != nil {
    	return err
    }
    p.l = l
    p.Addr = l.Addr().String()
    srv.Addr = p.Addr
    return srv.Serve(tcpKeepAliveListener{l.(*net.TCPListener)})
}

//Serve 开启监听
//  参：
//      l net.Listener  监听对象
//  返：
//      error           错误
func (p *Proxy) Serve(l net.Listener) error{
    p.setDefault()
    srv := p.initServer()
    p.l = l
    p.Addr = l.Addr().String()
    srv.Addr = p.Addr
    return p.Server.Serve(l)
}

//Close 关闭
//  返：
//      error       错误
func (p *Proxy) Close() error {
    if tr, ok := p.Transport.(*http.Transport); ok {
        tr.CloseIdleConnections()
    }
    return p.l.Close()
}

func copyDate(dst io.Writer, src io.ReadCloser, bufSize int) (n int64, err error){
    defer src.Close()
    buf := make([]byte, bufSize)
    return io.CopyBuffer(dst, src, buf)
}



type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}