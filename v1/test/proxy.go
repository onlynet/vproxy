package main

import (
	"github.com/456vv/vproxy"
    "net"
    "net/http"
    "net/url"
    "time"
    "flag"
    "fmt"
)
var faddr = flag.String("addr", "0.0.0.0:8080", "`string`: 代理服务器地")
var fproxy = flag.String("proxy", "", "`string`: 代理服务器的上级代理IP地址")
var fmaxIdleConnsPerHost = flag.Int("maxIdleConnsPerHost", 512, "`int`: 保持空闲连接数量")
var fdisableKeepAlives = flag.Bool("disableKeepAlives", false, "`bool`: 禁止长连接 (default false)")
var fdisableCompression = flag.Bool("disableCompression", false, "`bool`: 禁止传送数据时候进行压缩 (default false)")
var ftlsHandshakeTimeout = flag.Int64("tlsHandshakeTimeout", 10000, "`int64`: SSL握手超时，单位毫秒")
var fexpectContinueTimeout = flag.Int64("expectContinueTimeout", 1000, "`int64`: http1.1过度到http2的等待超时，单位毫秒")
var fresponseHeaderTimeout = flag.Int64("responseHeaderTimeout", 0, "`int64`: 读取服务器发来的文件标头超时，单位毫秒")
var fdataBufioSize = flag.Int("dataBufioSize", 1024*10, "`int`: 代理数据交换缓冲区大小，单位字节")

func main(){
    flag.Parse()
    //flag.Usage()
    //return
	p := &vproxy.Proxy{
        Addr        : *faddr,
        Transport   : &http.Transport{
            //Proxy: func(*Request) (*url.URL, error),
            //Dial: func(network, addr string) (net.Conn, error),
            //DialTLS: func(network, addr string) (net.Conn, error),
            //TLSClientConfig: *tls.Config,
            DisableKeepAlives: *fdisableKeepAlives,
            DisableCompression: *fdisableCompression,
            MaxIdleConnsPerHost: *fmaxIdleConnsPerHost,
            ResponseHeaderTimeout: time.Duration(*fresponseHeaderTimeout) * time.Millisecond,
            //TLSNextProto: map[string]func(authority string, c *tls.Conn),
            TLSHandshakeTimeout:   time.Duration(*ftlsHandshakeTimeout) * time.Millisecond,
            ExpectContinueTimeout: time.Duration(*fexpectContinueTimeout) * time.Millisecond,
        },
        Config      : &vproxy.Config{
            DataBufioSize: *fdataBufioSize,
        },
    }
    if tr, ok := p.Transport.(*http.Transport); ok && *fproxy != "" {
        tr.Proxy = func(r *http.Request) (*url.URL, error){
            return r.URL, nil
        }
        tr.Dial = func(network, addr string) (net.Conn, error){
        	return net.Dial(network, *fproxy)
        }
    }
    go func(){
        time.Sleep(time.Second)
        fmt.Printf("vproxy-IP: %s\r\n", p.Addr)
    }()
    err := p.ListenAndServ()
    fmt.Println("vproxy-error: ", err)
}





























