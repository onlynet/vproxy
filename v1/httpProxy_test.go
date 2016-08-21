package vproxy

import (
	"testing"
    "net"
    //"io/ioutil"
    //"io"
    //"net/url"
    "net/http"
    "bufio"
    //"bytes"
    "time"
    "fmt"
    "crypto/tls"
)
func Test_httpProxy_ServeHTTP(t *testing.T) {
    tests := []struct{
        req string
        statusCode int
    }{
        {req:"GET http://www.baidu.com:80/index.html HTTP/1.0\r\nHost:abcdef\r\nConnection:Close\r\n\r\n", statusCode:200},
        {req:"GET http://www.baidu.com:80/ HTTP/1.1\r\nHost:abcdef\r\nConnection:Keep-Alive\r\n\r\n", statusCode:200},
        {req:"GET /index.html?123 HTTP/1.1\r\nHost:www.baidu.com:80\r\nConnection:Keep-Alive\r\n\r\n", statusCode:200},
        {req:"GET https://translate.google.com:443/?124 HTTP/1.1\r\nHost:translate.google.com\r\nConnection:Keep-Alive\r\n\r\n", statusCode:200},
        {req:"GET https://kyfw.12306.cn/otn/regist/init HTTP/1.1\r\nHost:kyfw.12306.cn\r\nConnection:Keep-Alive\r\n\r\n", statusCode:502},
    }

    //服务器
    c := &Config{
        DataBufioSize:1024,
    }
    tlsConfig := &tls.Config{}
    hp := &httpProxy{
        config  : c,
        transport       : &http.Transport{
        //Proxy: func(*Request) (*url.URL, error), Dial: func(network, addr string) (net.Conn, error),
            DialTLS: func(network, addr string) (net.Conn, error){
    			rc, err := net.Dial("tcp", addr)
    			if err != nil {
    				return nil, err
    			}
                host, _, _ := net.SplitHostPort(addr)
                tlsConfig.ServerName = host
    			return tls.Client(rc, tlsConfig), nil
            },
            TLSClientConfig: tlsConfig,
        //DisableKeepAlives: bool, DisableCompression: bool, MaxIdleConnsPerHost: int, ResponseHeaderTimeout: time.Duration, TLSHandshakeTimeout: time.Duration, ExpectContinueTimeout: time.Duration, TLSNextProto: map[string]func(authority string, c *tls.Conn),}{

        },
    }
    srv := &http.Server{
        Handler: http.HandlerFunc(hp.ServeHTTP),
    }
    l, err := net.Listen("tcp", "127.0.0.1:1210")
    if err != nil {
        t.Fatal(err)
    }
    laddr := l.Addr().String()
    fmt.Println("服务器IP: ", laddr)

    go srv.Serve(l)
    time.Sleep(time.Second*2)

    for _, test := range tests {
        netConn, err := net.Dial("tcp", laddr)
        if err != nil {
            t.Fatalf("连接 %s 地址发生错误：%s", laddr, err)
        }
        netConn.Write([]byte(test.req))
        httpResponse, err := http.ReadResponse(bufio.NewReader(netConn), &http.Request{})
        if err != nil {
            t.Fatalf("连接 %s 地址，返回内容发生错误：%s", test.req, err)
        }
        if httpResponse.StatusCode != test.statusCode {
            t.Fatalf("连接 %s 地址，返回状态码不是200，是：%d", test.req, httpResponse.StatusCode)
        }
    }


}