package main

import (
	"github.com/456vv/vproxy/v1"
    "net"
    "net/http"
    "net/url"
    "time"
    "flag"
    "fmt"
    "log"
)
var fshow = flag.Bool("show", false, "在控制台显示请求网址")
var faddr = flag.String("addr", "0.0.0.0:8080", "代理服务器地")
var fproxy = flag.String("proxy", "", "代理服务器的上级代理IP地址 (format \"11.22.33.44:8888\")")
var fmaxIdleConnsPerHost = flag.Int("maxIdleConnsPerHost", 500, "保持空闲连接数量")
var fdisableKeepAlives = flag.Bool("disableKeepAlives", false, "禁止长连接 (default false)")
var fdisableCompression = flag.Bool("disableCompression", false, "禁止传送数据时候进行压缩 (default false)")
var ftlsHandshakeTimeout = flag.Int64("tlsHandshakeTimeout", 10000, "SSL握手超时，单位毫秒")
var fexpectContinueTimeout = flag.Int64("expectContinueTimeout", 1000, "http1.1过度到http2的等待超时，单位毫秒")
var fresponseHeaderTimeout = flag.Int64("responseHeaderTimeout", 0, "读取服务器发来的文件标头超时，单位毫秒 (default 0)")
var fdataBufioSize = flag.Int("dataBufioSize", 1024*10, "代理数据交换缓冲区大小，单位字节")

func main(){
    flag.Parse()
    //if flag.NFlag() == 0 {
    //    flag.PrintDefaults()
    //    return
    //}
	p := &vproxy.Proxy{
        Addr        : *faddr,
        Transport   : &http.Transport{
            DisableKeepAlives: *fdisableKeepAlives,
            DisableCompression: *fdisableCompression,
            MaxIdleConnsPerHost: *fmaxIdleConnsPerHost,
            ResponseHeaderTimeout: time.Duration(*fresponseHeaderTimeout) * time.Millisecond,
            TLSHandshakeTimeout:   time.Duration(*ftlsHandshakeTimeout) * time.Millisecond,
            ExpectContinueTimeout: time.Duration(*fexpectContinueTimeout) * time.Millisecond,
        },
        Config      : &vproxy.Config{
            DataBufioSize: *fdataBufioSize,
        },
    }
    if *fshow {
        p.Transport.(*http.Transport).Proxy = func(req *http.Request) (*url.URL, error){
            log.Println(req.URL)
            return nil, nil
        }
    }
    if tr, ok := p.Transport.(*http.Transport); ok && *fproxy != "" {
        tr.Proxy = func(r *http.Request) (*url.URL, error){
            return r.URL, nil
        }
        tr.Dial = func(network, addr string) (net.Conn, error){
        	return net.Dial(network, *fproxy)
        }
    }

    var err error
    exit := make(chan bool, 1)
    go func(){
        defer func(){
            p.Close()
            exit <- true
            close(exit)
        }()
        time.Sleep(time.Second)
        log.Printf("vproxy-IP: %s\r\n", p.Addr)

        var in0 string
        for err == nil  {
            log.Println("输入任何字符，并回车可以退出vproxy!")
            fmt.Scan(&in0)
            if in0 != "" {
                return
            }
        }
    }()
    defer p.Close()
    err = p.ListenAndServ()
    if err != nil {
        log.Printf("vproxy-Error：%s", err)
    }
    <-exit

}
