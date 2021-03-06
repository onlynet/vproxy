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
    "io"
    "os"
    "context"
)

var (
    flog = flag.String("log", "", "日志文件(默认留空在控制台显示日志)  (format \"./vproxy.txt\")")
    fuser = flag.String("user", "", "用户名")
    fpwd = flag.String("pwd", "", "密码")
    flogLevel = flag.Int("logLevel", 0, "日志级别，0)不记录 1)客户端IP 2)认证 3)访问的Host地址 4)路径 5)请求 6)响应 7)错误 (default 0)")
    faddr = flag.String("addr", "", "代理服务器地 (format \"0.0.0.0:8080\")")
    fproxy = flag.String("proxy", "", "代理服务器的上级代理IP地址 (format \"11.22.33.44:8888\")")
    fmaxIdleConns = flag.Int("maxIdleConns", 500, "保持空闲连接(TCP)数量")
    fmaxIdleConnsPerHost = flag.Int("maxIdleConnsPerHost", 500, "保持空闲连接(Host)数量")
    fdisableKeepAlives = flag.Bool("disableKeepAlives", false, "禁止长连接 (default false)")
    fdisableCompression = flag.Bool("disableCompression", false, "禁止传送数据时候进行压缩 (default false)")
    ftlsHandshakeTimeout = flag.Int64("tlsHandshakeTimeout", 10000, "SSL握手超时，单位毫秒")
    fidleConnTimeout = flag.Int64("idleConnTimeout", 0, "空闲连接超时时，单位毫秒 (default 0)")
    fexpectContinueTimeout = flag.Int64("expectContinueTimeout", 1000, "http1.1过度到http2的等待超时，单位毫秒")
    fresponseHeaderTimeout = flag.Int64("responseHeaderTimeout", 0, "读取服务器发来的文件标头超时，单位毫秒 (default 0)")
    fdataBufioSize = flag.Int("dataBufioSize", 1024*10, "代理数据交换缓冲区大小，单位字节")
    ftimeout = flag.Int64("timeout", 5000, "转发连接请求超时，单位毫秒")
)

func main(){
    flag.Parse()
    if flag.NFlag() == 0 {
        flag.PrintDefaults()
        fmt.Println("\r\n命令行例子：vproxy -addr 0.0.0.0:8080")
        return
    }
    var out io.Writer = os.Stdout
    if *flog != "" {
        file, err := os.OpenFile(*flog, os.O_CREATE | os.O_RDWR, 0777)
        if err != nil {
            fmt.Println(err)
            return
        }
        out = file
    }
	p := &vproxy.Proxy{
        Config      : &vproxy.Config{
            DataBufioSize: *fdataBufioSize,
            Timeout: time.Duration(*ftimeout) * time.Millisecond,
        },
        Addr        : *faddr,
        Transport   : &http.Transport{
            DisableKeepAlives: *fdisableKeepAlives,
            DisableCompression: *fdisableCompression,
            MaxIdleConns: *fmaxIdleConns,
            MaxIdleConnsPerHost: *fmaxIdleConnsPerHost,
            IdleConnTimeout: time.Duration(*fidleConnTimeout) * time.Millisecond,
            ResponseHeaderTimeout: time.Duration(*fresponseHeaderTimeout) * time.Millisecond,
            TLSHandshakeTimeout:   time.Duration(*ftlsHandshakeTimeout) * time.Millisecond,
            ExpectContinueTimeout: time.Duration(*fexpectContinueTimeout) * time.Millisecond,
        },
        ErrorLogLevel: vproxy.LogLevel(*flogLevel),
    }
    p.ErrorLog = log.New(out, "", log.LstdFlags)
    if *fuser != "" {
        p.Config.Auth = func(username, password string) bool {
            return username == *fuser && password == *fpwd
        }
     }

    if tr, ok := p.Transport.(*http.Transport); ok && *fproxy != "" {
        tr.Proxy = func(r *http.Request) (*url.URL, error){
            return r.URL, nil
        }
        tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error){
            ctx, cancel := context.WithTimeout(ctx, p.Config.Timeout)
            defer cancel()
        	return new(net.Dialer).DialContext(ctx, network, *fproxy)
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
